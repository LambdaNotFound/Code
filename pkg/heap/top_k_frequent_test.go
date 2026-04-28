package heap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_topKFrequent(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		k        int
		expected []int
	}{
		{name: "leetcode_example1", nums: []int{1, 1, 1, 2, 2, 3}, k: 2, expected: []int{1, 2}},
		{name: "leetcode_example2", nums: []int{1}, k: 1, expected: []int{1}},
		{name: "all_same", nums: []int{5, 5, 5, 5}, k: 1, expected: []int{5}},
		{name: "k_equals_len", nums: []int{1, 2, 3}, k: 3, expected: []int{1, 2, 3}},
		{name: "negatives", nums: []int{-1, -1, -2, -2, -2}, k: 1, expected: []int{-2}},
		{name: "two_tied_pick_both", nums: []int{1, 1, 2, 2}, k: 2, expected: []int{1, 2}},
		// Distinct frequencies for every element so the result is deterministic:
		// 4→freq3, 3→freq2, 2→freq1. No ties at the k-th boundary.
		{name: "large_k", nums: []int{4, 4, 4, 3, 3, 2}, k: 3, expected: []int{4, 3, 2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := topKFrequent(tt.nums, tt.k)
			assert.ElementsMatch(t, tt.expected, got)
		})
	}
}
