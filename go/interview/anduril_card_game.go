package interview

import (
	"errors"
	"math/rand"
	"time"
)

// Card is modeled as 0..51 for this problem. In real code you'd likely use a
// richer type (e.g. struct{ Rank, Suit int }); nothing here depends on that.
type Card int

// ErrEmptyDeck is returned by Draw when a finite deck is exhausted.
var ErrEmptyDeck = errors.New("deck is empty")

// Unbounded is the sentinel Remaining() value for a source that never exhausts.
const Unbounded = -1

// Rng is the randomness contract. *rand.Rand satisfies it, and tests inject a
// scripted fake. Injecting this (rather than calling the global rand) is the
// key design decision that makes the draw logic testable.
type Rng interface {
	Intn(n int) int
}

// CardSource is the common contract. FiniteDeck and InfiniteDeck are the two
// interchangeable implementations (Strategy), selected by MakeDeck (Factory).
type CardSource interface {
	Draw() (Card, error)
	Remaining() int // Unbounded (-1) for an infinite source
}

func defaultRng(rng Rng) Rng {
	if rng != nil {
		return rng
	}
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

// FiniteDeck draws without replacement using an incremental Fisher-Yates
// shuffle: cards[0:n) are still drawable; each draw picks a uniform index in
// that region, swaps the chosen card to the boundary, and shrinks it. O(1) per
// draw, uniform, no repeats.
type FiniteDeck struct {
	cards []Card
	n     int // [0, n) is the drawable region
	rng   Rng
}

// NewFiniteDeck copies the provided card universe so callers can't mutate
// internal state, and so the same slice can seed several decks.
func NewFiniteDeck(cards []Card, rng Rng) *FiniteDeck {
	cp := make([]Card, len(cards))
	copy(cp, cards)
	return &FiniteDeck{cards: cp, n: len(cp), rng: defaultRng(rng)}
}

func (d *FiniteDeck) Draw() (Card, error) {
	if d.n == 0 {
		return 0, ErrEmptyDeck
	}
	i := d.rng.Intn(d.n) // uniform over the remaining region
	d.n--
	d.cards[i], d.cards[d.n] = d.cards[d.n], d.cards[i]
	return d.cards[d.n], nil
}

func (d *FiniteDeck) Remaining() int { return d.n }

// InfiniteDeck models infinitely many decks. With k copies of each value the
// probability of any value is k/(52k) = 1/52, and the depletion effect after m
// draws is O(m/k); as k -> infinity for any fixed m the draws converge to
// i.i.d. uniform. So the correct model is a memoryless sampler WITH
// replacement -- it carries no state and never exhausts.
type InfiniteDeck struct {
	values []Card
	rng    Rng
}

func NewInfiniteDeck(values []Card, rng Rng) *InfiniteDeck {
	cp := make([]Card, len(values))
	copy(cp, values)
	return &InfiniteDeck{values: cp, rng: defaultRng(rng)}
}

func (d *InfiniteDeck) Draw() (Card, error) {
	return d.values[d.rng.Intn(len(d.values))], nil // i.i.d. uniform, with replacement
}

func (d *InfiniteDeck) Remaining() int { return Unbounded }

// StandardCards returns numDecks copies of values 0..51 (156 cards for 3 decks).
func StandardCards(numDecks int) []Card {
	cards := make([]Card, 0, 52*numDecks)
	for c := 0; c < 52; c++ {
		for k := 0; k < numDecks; k++ {
			cards = append(cards, Card(c))
		}
	}
	return cards
}

// MakeDeck is the factory: pass Unbounded for an infinite source, otherwise the
// number of standard decks. Adding the infinite case did not require touching
// the finite logic (Open/Closed).
func MakeDeck(numDecks int, rng Rng) CardSource {
	if numDecks == Unbounded {
		return NewInfiniteDeck(StandardCards(1), rng)
	}
	return NewFiniteDeck(StandardCards(numDecks), rng)
}

// ShuffledDeck is the "shuffle the whole slice once, then hand out in order"
// approach. rand.Shuffle is itself Fisher-Yates, so this is distributionally
// identical to FiniteDeck -- the same algorithm scheduled eagerly instead of
// lazily. Note it takes a concrete *rand.Rand, not the tiny Rng interface,
// because Shuffle is a method on *rand.Rand (it needs more than Intn).
type ShuffledDeck struct {
	cards []Card
	next  int
}

func NewShuffledDeck(cards []Card, r *rand.Rand) *ShuffledDeck {
	cp := make([]Card, len(cards))
	copy(cp, cards)
	r.Shuffle(len(cp), func(i, j int) { cp[i], cp[j] = cp[j], cp[i] })
	return &ShuffledDeck{cards: cp}
}

func (d *ShuffledDeck) Draw() (Card, error) {
	if d.next == len(d.cards) {
		return 0, ErrEmptyDeck
	}
	c := d.cards[d.next]
	d.next++
	return c, nil
}

func (d *ShuffledDeck) Remaining() int { return len(d.cards) - d.next }
