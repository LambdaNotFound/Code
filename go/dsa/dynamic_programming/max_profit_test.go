package dynamic_programming

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_maxProfit(t *testing.T) {
	tests := []struct {
		name     string
		prices   []int
		expected int
	}{
		{name: "example1", prices: []int{7, 1, 5, 3, 6, 4}, expected: 5},
		{name: "example2", prices: []int{7, 6, 4, 3, 1}, expected: 0},
		{name: "single_price", prices: []int{5}, expected: 0},
		{name: "two_prices_profit", prices: []int{1, 10}, expected: 9},
		{name: "two_prices_no_profit", prices: []int{10, 1}, expected: 0},
		{name: "buy_at_start", prices: []int{1, 2, 3, 4, 5}, expected: 4},
		{name: "sell_at_end", prices: []int{5, 4, 3, 2, 10}, expected: 8},
		{name: "all_same", prices: []int{3, 3, 3, 3}, expected: 0},
		{name: "valley_then_peak", prices: []int{3, 1, 4, 1, 5, 9, 2, 6}, expected: 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, maxProfit(tt.prices))
		})
	}
}
