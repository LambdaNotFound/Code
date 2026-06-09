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
