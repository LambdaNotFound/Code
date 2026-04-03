package backtracking

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_combineClaude(t *testing.T) {
	testCases := []struct {
		name     string
		n, k     int
		expected [][]int
	}{
		{
			name: "n4_k2",
			n:    4, k: 2,
			expected: [][]int{{1, 2}, {1, 3}, {1, 4}, {2, 3}, {2, 4}, {3, 4}},
		},
		{
			name: "n1_k1",
			n:    1, k: 1,
			expected: [][]int{{1}},
		},
		{
			name: "k_equals_n",
			n:    3, k: 3,
			expected: [][]int{{1, 2, 3}},
		},
		{
			name: "n5_k3",
			n:    5, k: 3,
			expected: [][]int{
				{1, 2, 3}, {1, 2, 4}, {1, 2, 5},
				{1, 3, 4}, {1, 3, 5}, {1, 4, 5},
				{2, 3, 4}, {2, 3, 5}, {2, 4, 5},
				{3, 4, 5},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := combineClaude(tc.n, tc.k)
			assert.ElementsMatch(t, tc.expected, result)
		})
	}
}

func Test_combinationSum4(t *testing.T) {
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
			result := combinationSum4(tc.nums, tc.target)
			assert.Equal(t, tc.expected, result)

			// both implementations must agree
			result2 := combinationSum4RecursionMemoization(tc.nums, tc.target)
			assert.Equal(t, tc.expected, result2)
		})
	}
}
