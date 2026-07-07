package heap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_smallestRange(t *testing.T) {
	tests := []struct {
		name     string
		nums     [][]int
		expected []int
	}{
		{
			name:     "leetcode_example",
			nums:     [][]int{{4, 10, 15, 24, 26}, {0, 9, 12, 20}, {5, 18, 22, 30}},
			expected: []int{20, 24},
		},
		{
			name:     "single_list",
			nums:     [][]int{{1, 2, 3}},
			expected: []int{1, 1},
		},
		{
			name:     "two_lists_one_element_each",
			nums:     [][]int{{2}, {4}},
			expected: []int{2, 4},
		},
		{
			name:     "two_lists_same_values",
			nums:     [][]int{{1}, {1}},
			expected: []int{1, 1},
		},
		{
			name:     "two_lists_overlapping",
			nums:     [][]int{{1, 5}, {3, 7}},
			expected: []int{1, 3},
		},
		{
			name:     "three_lists_all_same_start",
			nums:     [][]int{{1, 2, 3}, {1, 4, 5}, {1, 6, 7}},
			expected: []int{1, 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := smallestRange(tt.nums)
			assert.Equal(t, tt.expected, result)
		})
	}
}
