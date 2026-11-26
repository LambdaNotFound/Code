package math

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_majorityElement(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected int
	}{
		{name: "leetcode_example1", nums: []int{3, 2, 3}, expected: 3},
		{name: "leetcode_example2", nums: []int{2, 2, 1, 1, 1, 2, 2}, expected: 2},
		{name: "single_element", nums: []int{1}, expected: 1},
		{name: "all_same", nums: []int{7, 7, 7}, expected: 7},
		{name: "two_elements", nums: []int{1, 2, 2}, expected: 2},
		{name: "majority_at_start", nums: []int{3, 3, 3, 1, 2}, expected: 3},
		{name: "majority_at_end", nums: []int{1, 2, 4, 4, 4}, expected: 4},
		{name: "large_majority", nums: []int{3, 3, 4, 2, 4, 4, 2, 4, 4}, expected: 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := majorityElement(tt.nums)
			assert.Equal(t, tt.expected, result)
		})
	}
}
