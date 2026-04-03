package dynamic_programming

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_lengthOfLIS(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected int
	}{
		{name: "example1", nums: []int{10, 9, 2, 5, 3, 7, 101, 18}, expected: 4},
		{name: "example2", nums: []int{0, 1, 0, 3, 2, 3}, expected: 4},
		{name: "example3", nums: []int{7, 7, 7, 7, 7}, expected: 1},
		{name: "single", nums: []int{5}, expected: 1},
		{name: "strictly_increasing", nums: []int{1, 2, 3, 4, 5}, expected: 5},
		{name: "strictly_decreasing", nums: []int{5, 4, 3, 2, 1}, expected: 1},
		{name: "two_elements_increasing", nums: []int{1, 2}, expected: 2},
		{name: "two_elements_equal", nums: []int{3, 3}, expected: 1},
		{name: "mixed", nums: []int{3, 5, 6, 2, 5, 4, 19, 5, 6, 7, 12}, expected: 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, lengthOfLIS(tt.nums))
			assert.Equal(t, tt.expected, lengthOfLISBinarySearch(tt.nums))
		})
	}
}
