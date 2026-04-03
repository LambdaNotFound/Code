package heap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_findKthLargest(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		k        int
		expected int
	}{
		{name: "leetcode_example1", nums: []int{3, 2, 1, 5, 6, 4}, k: 2, expected: 5},
		{name: "leetcode_example2", nums: []int{3, 2, 3, 1, 2, 4, 5, 5, 6}, k: 4, expected: 4},
		{name: "k_equals_1", nums: []int{1, 2, 3}, k: 1, expected: 3},
		{name: "k_equals_len", nums: []int{1, 2, 3}, k: 3, expected: 1},
		{name: "all_duplicates", nums: []int{5, 5, 5, 5}, k: 2, expected: 5},
		{name: "negatives", nums: []int{-1, -2, -3}, k: 1, expected: -1},
		{name: "single_element", nums: []int{7}, k: 1, expected: 7},
		{name: "mixed_negatives_positives", nums: []int{-5, 3, 0, -1, 2}, k: 2, expected: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findKthLargest(tt.nums, tt.k)
			assert.Equal(t, tt.expected, result)
		})
	}
}
