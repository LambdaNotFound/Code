package monoqueue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_maxSlidingWindow(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		k        int
		expected []int
	}{
		{"leetcode_1", []int{1, 3, -1, -3, 5, 3, 6, 7}, 3, []int{3, 3, 5, 5, 6, 7}},
		{"leetcode_2", []int{1}, 1, []int{1}},
		{"k_equals_len", []int{4, 2, 5, 1}, 4, []int{5}},
		{"k_one", []int{3, 1, 2}, 1, []int{3, 1, 2}},
		{"all_same", []int{5, 5, 5, 5}, 2, []int{5, 5, 5}},
		{"decreasing", []int{10, 9, 8, 7}, 2, []int{10, 9, 8}},
		{"increasing", []int{1, 2, 3, 4}, 2, []int{2, 3, 4}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input1 := append([]int(nil), tt.nums...)
			assert.Equal(t, tt.expected, maxSlidingWindow(input1, tt.k), "maxSlidingWindow")
		})
	}
}

func Test_longestSubarray(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		limit    int
		expected int
	}{
		{"leetcode_1", []int{8, 2, 4, 7}, 4, 2},
		{"leetcode_2", []int{10, 1, 2, 4, 7, 2}, 5, 4},
		{"leetcode_3", []int{4, 2, 2, 2, 4, 4, 2, 2}, 0, 3},
		{"single_element", []int{5}, 0, 1},
		{"all_same", []int{3, 3, 3, 3}, 0, 4},
		{"limit_zero_all_diff", []int{1, 2, 3}, 0, 1},
		{"entire_array_valid", []int{1, 3, 2, 4}, 10, 4},
		{"large_limit", []int{1, 100, 1, 100}, 99, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nums := append([]int(nil), tt.nums...)
			got := longestSubarray(nums, tt.limit)
			assert.Equal(t, tt.expected, got)

			// brute force must agree
			nums2 := append([]int(nil), tt.nums...)
			assert.Equal(t, tt.expected, longestSubarrayBruteForce(nums2, tt.limit))
		})
	}
}
