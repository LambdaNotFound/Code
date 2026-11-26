package dp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_combinationSum4Memoization(t *testing.T) {
	testCases := []struct {
		name     string
		nums     []int
		target   int
		expected int
	}{
		{"leetcode_example_1", []int{1, 2, 3}, 4, 7},
		{"leetcode_example_2", []int{9}, 3, 0},
		{"single_num_exact", []int{3}, 3, 1},
		{"target_zero", []int{1, 2}, 0, 1}, // one way: pick nothing
		{"multiple_ways", []int{1, 2, 3}, 3, 4},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := combinationSum4Memoization(tc.nums, tc.target)
			assert.Equal(t, tc.expected, result)
		})
	}
}
