package interview

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ── helpers ──────────────────────────────────────────────────────────────────

func newPlayer(gems map[Color]int, hand []*Card) *Player {
	return &Player{gems: gems, hand: hand}
}

func newCard(color Color, cost map[Color]int) *Card {
	return &Card{color: color, cost: cost}
}

// ── Color.String ─────────────────────────────────────────────────────────────

func Test_Color_String(t *testing.T) {
	testCases := []struct {
		name     string
		color    Color
		expected string
	}{
		{name: "blue", color: Blue, expected: "B"},
		{name: "white", color: White, expected: "W"},
		{name: "green", color: Green, expected: "G"},
		{name: "red", color: Red, expected: "R"},
		{name: "yellow", color: Yellow, expected: "Y"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.color.String())
		})
	}
}

// ── discountsOf ──────────────────────────────────────────────────────────────

func Test_discountsOf(t *testing.T) {
	testCases := []struct {
		name     string
		hand     []*Card
		expected map[Color]int
	}{
		{
			name:     "empty_hand",
			hand:     []*Card{},
			expected: map[Color]int{},
		},
		{
			name:     "single_blue_card",
			hand:     []*Card{newCard(Blue, nil)},
			expected: map[Color]int{Blue: 1},
		},
		{
			name:     "two_blue_cards",
			hand:     []*Card{newCard(Blue, nil), newCard(Blue, nil)},
			expected: map[Color]int{Blue: 2},
		},
		{
			name: "mixed_colors",
			hand: []*Card{
				newCard(Blue, nil),
				newCard(Red, nil),
				newCard(Blue, nil),
				newCard(Green, nil),
			},
			expected: map[Color]int{Blue: 2, Red: 1, Green: 1},
		},
		{
			name: "all_five_colors",
			hand: []*Card{
				newCard(Blue, nil), newCard(White, nil), newCard(Green, nil),
				newCard(Red, nil), newCard(Yellow, nil),
			},
			expected: map[Color]int{Blue: 1, White: 1, Green: 1, Red: 1, Yellow: 1},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := newPlayer(nil, tc.hand)
			assert.Equal(t, tc.expected, discountsOf(p))
		})
	}
}

// ── effectiveCost ─────────────────────────────────────────────────────────────

func Test_effectiveCost(t *testing.T) {
	testCases := []struct {
		name     string
		cardCost map[Color]int
		hand     []*Card // cards player already holds (source of discounts)
		expected map[Color]int
	}{
		{
			name:     "no_discount_no_hand",
			cardCost: map[Color]int{Blue: 3, Red: 2},
			hand:     []*Card{},
			expected: map[Color]int{Blue: 3, Red: 2},
		},
		{
			name:     "partial_discount_reduces_cost",
			cardCost: map[Color]int{Blue: 3, Red: 2},
			hand:     []*Card{newCard(Blue, nil), newCard(Blue, nil)},
			expected: map[Color]int{Blue: 1, Red: 2},
		},
		{
			name:     "discount_equals_cost_gives_zero",
			cardCost: map[Color]int{Blue: 2},
			hand:     []*Card{newCard(Blue, nil), newCard(Blue, nil)},
			expected: map[Color]int{Blue: 0},
		},
		{
			name:     "discount_exceeds_cost_clamped_at_zero",
			cardCost: map[Color]int{Blue: 1},
			hand:     []*Card{newCard(Blue, nil), newCard(Blue, nil), newCard(Blue, nil)},
			expected: map[Color]int{Blue: 0},
		},
		{
			name:     "discount_on_unrelated_color_has_no_effect",
			cardCost: map[Color]int{Red: 3},
			hand:     []*Card{newCard(Blue, nil), newCard(Blue, nil)},
			expected: map[Color]int{Red: 3},
		},
		{
			name:     "full_discount_all_colors",
			cardCost: map[Color]int{Blue: 1, White: 1, Green: 1},
			hand: []*Card{
				newCard(Blue, nil), newCard(White, nil), newCard(Green, nil),
			},
			expected: map[Color]int{Blue: 0, White: 0, Green: 0},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			card := newCard(Blue, tc.cardCost) // card color irrelevant for effectiveCost
			p := newPlayer(nil, tc.hand)
			assert.Equal(t, tc.expected, effectiveCost(card, p))
		})
	}
}

// ── canPurchase ───────────────────────────────────────────────────────────────

func Test_canPurchase(t *testing.T) {
	testCases := []struct {
		name     string
		gems     map[Color]int
		hand     []*Card
		cardCost map[Color]int
		expected bool
	}{
		{
			name:     "exact_gems_no_discount",
			gems:     map[Color]int{Blue: 3},
			hand:     []*Card{},
			cardCost: map[Color]int{Blue: 3},
			expected: true,
		},
		{
			name:     "surplus_gems",
			gems:     map[Color]int{Blue: 5, Red: 4},
			hand:     []*Card{},
			cardCost: map[Color]int{Blue: 3, Red: 2},
			expected: true,
		},
		{
			name:     "one_gem_short",
			gems:     map[Color]int{Blue: 2},
			hand:     []*Card{},
			cardCost: map[Color]int{Blue: 3},
			expected: false,
		},
		{
			name:     "missing_color_entirely",
			gems:     map[Color]int{Red: 5},
			hand:     []*Card{},
			cardCost: map[Color]int{Blue: 1},
			expected: false,
		},
		{
			name:     "discount_makes_purchase_possible",
			gems:     map[Color]int{Blue: 1, Red: 2},
			hand:     []*Card{newCard(Blue, nil), newCard(Blue, nil)}, // 2 Blue discount
			cardCost: map[Color]int{Blue: 3, Red: 2},                 // effective Blue: 1
			expected: true,
		},
		{
			name:     "discount_covers_cost_fully_zero_gems_needed",
			gems:     map[Color]int{},
			hand:     []*Card{newCard(Blue, nil), newCard(Blue, nil), newCard(Blue, nil)},
			cardCost: map[Color]int{Blue: 2},
			expected: true,
		},
		{
			name:     "multi_color_one_color_short",
			gems:     map[Color]int{Blue: 3, Red: 1},
			hand:     []*Card{},
			cardCost: map[Color]int{Blue: 3, Red: 2},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := newPlayer(tc.gems, tc.hand)
			c := newCard(Blue, tc.cardCost)
			assert.Equal(t, tc.expected, canPurchase(p, c))
		})
	}
}

// ── purchase ──────────────────────────────────────────────────────────────────

func Test_purchase(t *testing.T) {
	testCases := []struct {
		name          string
		startGems     map[Color]int
		hand          []*Card
		cardCost      map[Color]int
		cardColor     Color
		wantOk        bool
		wantGems      map[Color]int // gems after the call
		wantHandLen   int
	}{
		{
			name:        "successful_purchase_deducts_gems",
			startGems:   map[Color]int{Blue: 5},
			hand:        []*Card{},
			cardCost:    map[Color]int{Blue: 3},
			cardColor:   Red,
			wantOk:      true,
			wantGems:    map[Color]int{Blue: 2},
			wantHandLen: 1,
		},
		{
			name:        "failed_purchase_leaves_state_unchanged",
			startGems:   map[Color]int{Blue: 2},
			hand:        []*Card{},
			cardCost:    map[Color]int{Blue: 3},
			cardColor:   Green,
			wantOk:      false,
			wantGems:    map[Color]int{Blue: 2},
			wantHandLen: 0,
		},
		{
			name:        "purchase_with_discount_deducts_reduced_amount",
			startGems:   map[Color]int{Blue: 1, Red: 2},
			hand:        []*Card{newCard(Blue, nil), newCard(Blue, nil)},
			cardCost:    map[Color]int{Blue: 3, Red: 2},
			cardColor:   White,
			wantOk:      true,
			wantGems:    map[Color]int{Blue: 0, Red: 0},
			wantHandLen: 3,
		},
		{
			name:        "purchase_when_discount_covers_full_cost",
			startGems:   map[Color]int{},
			hand:        []*Card{newCard(Blue, nil), newCard(Blue, nil)},
			cardCost:    map[Color]int{Blue: 2},
			cardColor:   Yellow,
			wantOk:      true,
			wantGems:    map[Color]int{},
			wantHandLen: 3,
		},
		{
			name:        "multi_color_purchase",
			startGems:   map[Color]int{Blue: 2, White: 1, Red: 3},
			hand:        []*Card{},
			cardCost:    map[Color]int{Blue: 2, White: 1, Red: 3},
			cardColor:   Green,
			wantOk:      true,
			wantGems:    map[Color]int{Blue: 0, White: 0, Red: 0},
			wantHandLen: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := newPlayer(tc.startGems, tc.hand)
			c := newCard(tc.cardColor, tc.cardCost)

			ok := purchase(p, c)

			assert.Equal(t, tc.wantOk, ok)
			assert.Equal(t, tc.wantHandLen, len(p.hand))
			for color, want := range tc.wantGems {
				assert.Equal(t, want, p.gems[color], "gems[%s]", color)
			}
		})
	}
}

// ── purchase concurrency ──────────────────────────────────────────────────────
//
// Fires N goroutines all trying to buy the same card simultaneously.
// The player only has enough gems for exactly one purchase, so exactly
// one goroutine must succeed and the rest must fail.  Run with -race to
// verify the mutex prevents data races.

func Test_purchase_concurrent(t *testing.T) {
	testCases := []struct {
		name        string
		gems        map[Color]int   // enough for exactly one purchase
		cardCost    map[Color]int
		goroutines  int
	}{
		{
			name:       "only_one_of_ten_succeeds",
			gems:       map[Color]int{Blue: 3},
			cardCost:   map[Color]int{Blue: 3},
			goroutines: 10,
		},
		{
			name:       "only_one_of_fifty_succeeds",
			gems:       map[Color]int{Blue: 2, Red: 1},
			cardCost:   map[Color]int{Blue: 2, Red: 1},
			goroutines: 50,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := newPlayer(tc.gems, []*Card{})
			c := newCard(Green, tc.cardCost)

			results := make([]bool, tc.goroutines)
			var wg sync.WaitGroup
			for i := 0; i < tc.goroutines; i++ {
				wg.Add(1)
				go func(idx int) {
					defer wg.Done()
					results[idx] = purchase(p, c)
				}(i)
			}
			wg.Wait()

			successes := 0
			for _, ok := range results {
				if ok {
					successes++
				}
			}
			assert.Equal(t, 1, successes, "exactly one goroutine should succeed")
			assert.Equal(t, 1, len(p.hand), "card added to hand exactly once")
		})
	}
}
