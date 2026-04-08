package solid_coding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_linearScan(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected [][]int
	}{
		{
			name:     "no_gaps",
			nums:     []int{1, 2, 3, 4},
			expected: [][]int{},
		},
		{
			name:     "single_gap",
			nums:     []int{1, 3},
			expected: [][]int{{2, 2}},
		},
		{
			name:     "multiple_gaps",
			nums:     []int{1, 4, 7, 10},
			expected: [][]int{{2, 3}, {5, 6}, {8, 9}},
		},
		{
			name:     "consecutive_and_gap_mixed",
			nums:     []int{0, 1, 3, 4, 7},
			expected: [][]int{{2, 2}, {5, 6}},
		},
		{
			name:     "single_element",
			nums:     []int{5},
			expected: [][]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, linearScan(tt.nums))
		})
	}
}

func Test_summaryRanges(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected []string
	}{
		{
			name:     "leetcode_example1",
			nums:     []int{0, 1, 2, 4, 5, 7},
			expected: []string{"0->2", "4->5", "7"},
		},
		{
			name:     "leetcode_example2",
			nums:     []int{0, 2, 3, 4, 6, 8, 9},
			expected: []string{"0", "2->4", "6", "8->9"},
		},
		{
			name:     "empty",
			nums:     []int{},
			expected: []string{},
		},
		{
			name:     "single_element",
			nums:     []int{42},
			expected: []string{"42"},
		},
		{
			name:     "all_consecutive",
			nums:     []int{1, 2, 3, 4, 5},
			expected: []string{"1->5"},
		},
		{
			name:     "no_consecutive",
			nums:     []int{1, 3, 5, 7},
			expected: []string{"1", "3", "5", "7"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, summaryRanges(tt.nums))
		})
	}
}

func Test_findMissingRanges(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		lower    int
		upper    int
		expected [][]int
	}{
		{
			name:     "leetcode_example",
			nums:     []int{0, 1, 3, 50, 75},
			lower:    0,
			upper:    99,
			expected: [][]int{{2, 2}, {4, 49}, {51, 74}, {76, 99}},
		},
		{
			name:     "empty_nums_full_range",
			nums:     []int{},
			lower:    1,
			upper:    5,
			expected: [][]int{{1, 5}},
		},
		{
			name:     "no_missing_ranges",
			nums:     []int{1, 2, 3},
			lower:    1,
			upper:    3,
			expected: [][]int{},
		},
		{
			name:     "missing_at_start",
			nums:     []int{3, 4, 5},
			lower:    1,
			upper:    5,
			expected: [][]int{{1, 2}},
		},
		{
			name:     "missing_at_end",
			nums:     []int{1, 2, 3},
			lower:    1,
			upper:    5,
			expected: [][]int{{4, 5}},
		},
		{
			name:     "single_missing",
			nums:     []int{1, 3},
			lower:    1,
			upper:    3,
			expected: [][]int{{2, 2}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, findMissingRanges(tt.nums, tt.lower, tt.upper))
		})
	}
}
