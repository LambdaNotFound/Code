package twopointers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_nextPermutation(t *testing.T) {
	testCases := []struct {
		name     string
		nums     []int
		expected []int
	}{
		{"leetcode_example_1", []int{1, 2, 3}, []int{1, 3, 2}},
		{"leetcode_example_2", []int{3, 2, 1}, []int{1, 2, 3}}, // last permutation wraps to first
		{"leetcode_example_3", []int{1, 1, 5}, []int{1, 5, 1}},
		{"single_element", []int{1}, []int{1}},
		{"two_elements_ascending", []int{1, 2}, []int{2, 1}},
		{"two_elements_descending", []int{2, 1}, []int{1, 2}},
		{"mid_sequence", []int{1, 3, 2}, []int{2, 1, 3}},
		{"larger_example", []int{1, 2, 3, 6, 5, 4}, []int{1, 2, 4, 3, 5, 6}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			nextPermutation(tc.nums)
			assert.Equal(t, tc.expected, tc.nums)
		})
	}
}
