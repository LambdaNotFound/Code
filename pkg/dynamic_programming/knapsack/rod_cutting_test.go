package knapsack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_maxProfit(t *testing.T) {
	tests := []struct {
		name     string
		length   int
		price    []int
		expected int
	}{
		{
			name:     "example_cut_2_and_6",
			length:   8,
			price:    []int{1, 5, 8, 9, 10, 17, 17, 20},
			expected: 22, // cut into 2+6: 5+17=22
		},
		{
			name:     "no_cut_full_rod",
			length:   1,
			price:    []int{10},
			expected: 10,
		},
		{
			name:     "uniform_price_no_cut_needed",
			length:   4,
			price:    []int{1, 2, 3, 10},
			expected: 10, // length 4 unsplit earns 10
		},
		{
			name:     "all_equal_prices_many_cuts",
			length:   4,
			price:    []int{3, 3, 3, 3},
			expected: 12, // 4 cuts of length 1: 4*3=12
		},
		{
			name:     "cut_into_equal_halves",
			length:   4,
			price:    []int{2, 5, 6, 7},
			expected: 10, // cut into 2+2: 5+5=10
		},
		{
			name:     "length_zero",
			length:   0,
			price:    []int{1, 5, 8, 9},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, maxProfit(tt.length, tt.price))
		})
	}
}

func Test_maxTotalProfit(t *testing.T) {
	tests := []struct {
		name        string
		rods        []int
		salePrice   float64
		costPerCut  float64
		expected    float64
	}{
		{
			// saleLength=6: 1 piece, 0 cuts, revenue=18, profit=18 (best)
			name:       "single_rod_best_is_full_length",
			rods:       []int{6},
			salePrice:  3.0,
			costPerCut: 4.0,
			expected:   18.0,
		},
		{
			// saleLength=7 for rod 7: 1 piece, 0 cuts, profit=7.0
			// saleLength=3 for rods [3,7]: rod3→profit=3, rod7→pieces=2,cuts=2,profit=2 → total=5
			// best across all lengths is 7.0 (saleLength=7 skips rod 3, profits 7 from rod 7)
			name:       "two_rods_optimal_uses_full_rod",
			rods:       []int{3, 7},
			salePrice:  1.0,
			costPerCut: 2.0,
			expected:   7.0,
		},
		{
			// high cut cost → only profitable at full rod length; saleLength=5: profit=5
			name:       "high_cut_cost_favors_no_cuts",
			rods:       []int{5},
			salePrice:  1.0,
			costPerCut: 10.0,
			expected:   5.0,
		},
		{
			// saleLength=2: rod4→pieces=2,cuts=1,profit=7; rod6→pieces=3,cuts=2,profit=10 → total=17
			name:       "two_rods_optimal_mid_length",
			rods:       []int{4, 6},
			salePrice:  2.0,
			costPerCut: 1.0,
			expected:   17.0,
		},
		{
			// rod shorter than every saleLength > rod: pieces=0 → skip (covers continue branch)
			// rod=[2], salePrice=1, costPerCut=0: saleLength=2 → pieces=1,cuts=0,profit=2
			name:       "short_rod_skip_longer_sale_lengths",
			rods:       []int{2},
			salePrice:  1.0,
			costPerCut: 0.0,
			expected:   2.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.InDelta(t, tt.expected, maxTotalProfit(tt.rods, tt.salePrice, tt.costPerCut), 1e-9)
		})
	}
}
