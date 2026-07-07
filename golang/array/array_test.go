package array

import (
	"testing"

	"github.com/stretchr/testify/assert"
)


func Test_firstMissingPositive(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected int
	}{
		{"leetcode_1", []int{1, 2, 0}, 3},
		{"leetcode_2", []int{3, 4, -1, 1}, 2},
		{"leetcode_3", []int{7, 8, 9, 11, 12}, 1},
		{"single_one", []int{1}, 2},
		{"single_two", []int{2}, 1},
		{"consecutive", []int{1, 2, 3, 4, 5}, 6},
		{"all_negative", []int{-1, -2, -3}, 1},
		{"with_duplicates", []int{1, 1, 2, 2}, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nums := append([]int(nil), tt.nums...)
			assert.Equal(t, tt.expected, firstMissingPositive(nums))
		})
	}
}

func Test_canArrange(t *testing.T) {
	tests := []struct {
		name  string
		team1 []int
		team2 []int
		want  bool
	}{
		{"shorter_team1_matchable", []int{1, 3}, []int{2, 4, 5}, true},
		{"shorter_team1_no_match", []int{5, 6}, []int{1, 2, 3}, false},
		{"single_front_covered", []int{3}, []int{1, 4, 2}, true},
		{"shorter_team2_matchable", []int{2, 4, 5}, []int{1, 3}, true},
		{"shorter_team2_no_match", []int{1, 2, 3}, []int{5, 6}, false},
		{"equal_size_team1_leads", []int{1, 2}, []int{3, 4}, true},
		{"equal_size_team2_leads", []int{3, 4}, []int{1, 2}, true},
		{"equal_size_partial_block", []int{1, 3}, []int{2, 3}, false},
		{"equal_heights_strict_fail", []int{1, 2}, []int{1, 2}, false},
		{"single_team2_leads", []int{2}, []int{1}, true},
		{"single_team1_leads", []int{1}, []int{2}, true},
		{"duplicates_front_covered", []int{1, 1}, []int{2, 3, 4}, true},
		{"duplicates_insufficient_back", []int{3, 3}, []int{1, 2, 4}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			team1 := append([]int(nil), tt.team1...)
			team2 := append([]int(nil), tt.team2...)
			assert.Equal(t, tt.want, canArrange(team1, team2))
		})
	}
}

func Test_removeDuplicatesFromSortedArray(t *testing.T) {
	tests := []struct {
		name         string
		nums         []int
		expectedLen  int
		expectedNums []int // first expectedLen elements
	}{
		{
			name:         "no_duplicates",
			nums:         []int{1, 2, 3, 4},
			expectedLen:  4,
			expectedNums: []int{1, 2, 3, 4},
		},
		{
			name:         "all_same",
			nums:         []int{2, 2, 2, 2},
			expectedLen:  2,
			expectedNums: []int{2, 2},
		},
		{
			name:         "each_appears_twice",
			nums:         []int{1, 1, 2, 2, 3, 3},
			expectedLen:  6,
			expectedNums: []int{1, 1, 2, 2, 3, 3},
		},
		{
			name:         "three_duplicates_trimmed",
			nums:         []int{1, 1, 1, 2, 2, 3},
			expectedLen:  5,
			expectedNums: []int{1, 1, 2, 2, 3},
		},
		{
			name:         "single_element",
			nums:         []int{7},
			expectedLen:  1,
			expectedNums: []int{7},
		},
		{
			name:         "two_same",
			nums:         []int{5, 5},
			expectedLen:  2,
			expectedNums: []int{5, 5},
		},
		{
			name:         "three_same_trimmed_to_two",
			nums:         []int{5, 5, 5},
			expectedLen:  2,
			expectedNums: []int{5, 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nums := append([]int(nil), tt.nums...)
			got := removeDuplicatesFromSortedArray(nums)
			assert.Equal(t, tt.expectedLen, got)
			assert.Equal(t, tt.expectedNums, nums[:got])
		})
	}
}
