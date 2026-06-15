package interview

import (
	"errors"
	"math/rand"
	"sort"
	"testing"
)

// fakeRng returns scripted Intn values so a test can pin the exact selection
// the algorithm makes, independent of any real distribution.
type fakeRng struct {
	vals []int
	i    int
}

func (f *fakeRng) Intn(int) int {
	v := f.vals[f.i]
	f.i++
	return v
}

// ---- Deterministic: without replacement + exhaustion (requirements 1 & 3) ----

func TestDrawsEveryCardExactlyOnce(t *testing.T) {
	d := NewFiniteDeck(StandardCards(1), rand.New(rand.NewSource(1)))
	got := make([]int, 0, 52)
	for i := 0; i < 52; i++ {
		c, err := d.Draw()
		if err != nil {
			t.Fatalf("unexpected error on draw %d: %v", i, err)
		}
		got = append(got, int(c))
	}
	sort.Ints(got)
	for i := 0; i < 52; i++ {
		if got[i] != i { // sorted == 0..51 proves complete AND no duplicates
			t.Fatalf("expected card %d at position %d, got %d", i, i, got[i])
		}
	}
	if d.Remaining() != 0 {
		t.Fatalf("expected 0 remaining, got %d", d.Remaining())
	}
}

func TestEmptyAfterAllDrawn(t *testing.T) {
	d := NewFiniteDeck(StandardCards(1), rand.New(rand.NewSource(1)))
	for i := 0; i < 52; i++ {
		if _, err := d.Draw(); err != nil {
			t.Fatalf("unexpected error on draw %d: %v", i, err)
		}
	}
	if _, err := d.Draw(); !errors.Is(err, ErrEmptyDeck) {
		t.Fatalf("expected ErrEmptyDeck on 53rd draw, got %v", err)
	}
}

// ---- Deterministic: reproducibility + mechanism (via injected RNG) ----

func TestReproducibleWithSeed(t *testing.T) {
	seq := func() []Card {
		d := NewFiniteDeck(StandardCards(1), rand.New(rand.NewSource(42)))
		out := make([]Card, 52)
		for i := range out {
			out[i], _ = d.Draw()
		}
		return out
	}
	a, b := seq(), seq()
	for i := range a {
		if a[i] != b[i] {
			t.Fatalf("same seed diverged at %d: %d vs %d", i, a[i], b[i])
		}
	}
}

func TestSelectionFollowsRng(t *testing.T) {
	// Intn always returns 0 -> always pick the front of the remaining region.
	// cards=[10,20,30]: draw1 returns 10, draw2 returns 30, draw3 returns 20.
	d := NewFiniteDeck([]Card{10, 20, 30}, &fakeRng{vals: []int{0, 0, 0}})
	want := []Card{10, 30, 20}
	for i, w := range want {
		got, _ := d.Draw()
		if got != w {
			t.Fatalf("draw %d: want %d, got %d", i, w, got)
		}
	}
}

// ---- Statistical: uniformity (requirement 2) ----

func TestFirstDrawIsUniform(t *testing.T) {
	const trials = 52 * 2000
	rng := rand.New(rand.NewSource(0)) // seeded -> chi2 is deterministic, non-flaky
	counts := make([]int, 52)
	for i := 0; i < trials; i++ {
		d := NewFiniteDeck(StandardCards(1), rng)
		c, _ := d.Draw()
		counts[c]++
	}
	expected := float64(trials) / 52.0
	chi2 := 0.0
	for _, c := range counts {
		diff := float64(c) - expected
		chi2 += diff * diff / expected
	}
	const critical = 90.6 // 51 dof, ~0.1% significance
	if chi2 >= critical {
		t.Fatalf("chi2=%.2f exceeds critical %.2f (non-uniform)", chi2, critical)
	}
	t.Logf("chi2=%.2f (< %.2f)", chi2, critical)
}

// ---- N decks: set-equality becomes multiset-equality ----

func TestThreeDecks(t *testing.T) {
	const k = 3
	d := MakeDeck(k, rand.New(rand.NewSource(7)))
	counts := make(map[Card]int)
	total := 0
	for {
		c, err := d.Draw()
		if errors.Is(err, ErrEmptyDeck) {
			break
		}
		counts[c]++
		total++
	}
	if total != 52*k {
		t.Fatalf("expected %d total draws, got %d", 52*k, total)
	}
	for v := Card(0); v < 52; v++ {
		if counts[v] != k { // each value appears exactly k times
			t.Fatalf("value %d drawn %d times, want %d", v, counts[v], k)
		}
	}
}

// ---- Infinite deck: the no-replacement and exhaustion tests INVERT ----

func TestInfiniteNeverExhausts(t *testing.T) {
	d := MakeDeck(Unbounded, rand.New(rand.NewSource(0)))
	for i := 0; i < 100_000; i++ {
		if _, err := d.Draw(); err != nil {
			t.Fatalf("infinite deck errored at draw %d: %v", i, err)
		}
	}
	if d.Remaining() != Unbounded {
		t.Fatalf("expected Unbounded remaining, got %d", d.Remaining())
	}
}

func TestInfiniteAllowsRepeats(t *testing.T) {
	d := MakeDeck(Unbounded, rand.New(rand.NewSource(0)))
	seen := make(map[Card]bool)
	const draws = 1000
	for i := 0; i < draws; i++ {
		c, _ := d.Draw()
		seen[c] = true
	}
	if len(seen) >= draws { // repeats are now expected, not a bug
		t.Fatalf("expected repeats over %d draws, saw %d distinct", draws, len(seen))
	}
}

// Compile-time proof it's a drop-in for the same contract.
var _ CardSource = (*ShuffledDeck)(nil)

func TestShuffledDrawsEveryCardExactlyOnce(t *testing.T) {
	d := NewShuffledDeck(StandardCards(1), rand.New(rand.NewSource(1)))
	got := make([]int, 0, 52)
	for i := 0; i < 52; i++ {
		c, err := d.Draw()
		if err != nil {
			t.Fatalf("unexpected error on draw %d: %v", i, err)
		}
		got = append(got, int(c))
	}
	sort.Ints(got)
	for i := 0; i < 52; i++ {
		if got[i] != i {
			t.Fatalf("expected %d at position %d, got %d", i, i, got[i])
		}
	}
	if d.Remaining() != 0 {
		t.Fatalf("expected 0 remaining, got %d", d.Remaining())
	}
}

func TestShuffledEmptyAfterAllDrawn(t *testing.T) {
	d := NewShuffledDeck(StandardCards(1), rand.New(rand.NewSource(1)))
	for i := 0; i < 52; i++ {
		if _, err := d.Draw(); err != nil {
			t.Fatalf("unexpected error on draw %d: %v", i, err)
		}
	}
	if _, err := d.Draw(); !errors.Is(err, ErrEmptyDeck) {
		t.Fatalf("expected ErrEmptyDeck, got %v", err)
	}
}

func TestShuffledReproducibleWithSeed(t *testing.T) {
	seq := func() []Card {
		d := NewShuffledDeck(StandardCards(1), rand.New(rand.NewSource(42)))
		out := make([]Card, 52)
		for i := range out {
			out[i], _ = d.Draw()
		}
		return out
	}
	a, b := seq(), seq()
	for i := range a {
		if a[i] != b[i] {
			t.Fatalf("same seed diverged at %d: %d vs %d", i, a[i], b[i])
		}
	}
}

// Same chi-square uniformity check, now over the FIRST card of an upfront shuffle.
func TestShuffledFirstDrawIsUniform(t *testing.T) {
	const trials = 52 * 2000
	rng := rand.New(rand.NewSource(0))
	counts := make([]int, 52)
	for i := 0; i < trials; i++ {
		d := NewShuffledDeck(StandardCards(1), rng)
		c, _ := d.Draw()
		counts[c]++
	}
	expected := float64(trials) / 52.0
	chi2 := 0.0
	for _, c := range counts {
		diff := float64(c) - expected
		chi2 += diff * diff / expected
	}
	const critical = 90.6
	if chi2 >= critical {
		t.Fatalf("chi2=%.2f exceeds critical %.2f", chi2, critical)
	}
	t.Logf("chi2=%.2f (< %.2f)", chi2, critical)
}

func TestShuffledThreeDecks(t *testing.T) {
	const k = 3
	d := NewShuffledDeck(StandardCards(k), rand.New(rand.NewSource(7)))
	counts := make(map[Card]int)
	total := 0
	for {
		c, err := d.Draw()
		if errors.Is(err, ErrEmptyDeck) {
			break
		}
		counts[c]++
		total++
	}
	if total != 52*k {
		t.Fatalf("expected %d total, got %d", 52*k, total)
	}
	for v := Card(0); v < 52; v++ {
		if counts[v] != k {
			t.Fatalf("value %d drawn %d times, want %d", v, counts[v], k)
		}
	}
}
