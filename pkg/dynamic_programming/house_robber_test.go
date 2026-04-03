package dynamic_programming

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_rob(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected int
	}{
		{name: "example1", nums: []int{1, 2, 3, 1}, expected: 4},
		{name: "example2", nums: []int{2, 7, 9, 3, 1}, expected: 12},
		{name: "single", nums: []int{5}, expected: 5},
		{name: "two_pick_larger", nums: []int{1, 9}, expected: 9},
		{name: "two_pick_first", nums: []int{9, 1}, expected: 9},
		{name: "all_same", nums: []int{3, 3, 3, 3}, expected: 6},
		{name: "alternating", nums: []int{5, 1, 5, 1, 5}, expected: 15},
		{name: "increasing", nums: []int{1, 2, 3, 4, 5}, expected: 9},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, rob(tt.nums))
		})
	}
}

func Test_robII(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected int
	}{
		{name: "example1", nums: []int{2, 3, 2}, expected: 3},
		{name: "example2", nums: []int{1, 2, 3, 1}, expected: 4},
		{name: "example3", nums: []int{1, 2, 3}, expected: 3},
		{name: "two_equal", nums: []int{5, 5}, expected: 5},
		{name: "two_pick_larger", nums: []int{3, 8}, expected: 8},
		{name: "all_same", nums: []int{4, 4, 4, 4}, expected: 8},
		{name: "large_ring", nums: []int{1, 3, 1, 3, 100}, expected: 103},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, robII(tt.nums))
		})
	}
}
