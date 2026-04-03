package dynamic_programming

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_longestCommonSubsequence(t *testing.T) {
	tests := []struct {
		name     string
		text1    string
		text2    string
		expected int
	}{
		{name: "example1", text1: "abcde", text2: "ace", expected: 3},
		{name: "example2", text1: "abc", text2: "abc", expected: 3},
		{name: "example3", text1: "abc", text2: "def", expected: 0},
		{name: "one_char_match", text1: "a", text2: "a", expected: 1},
		{name: "one_char_no_match", text1: "a", text2: "b", expected: 0},
		{name: "subsequence_interleaved", text1: "abcba", text2: "abcbcba", expected: 5},
		{name: "empty_like_single", text1: "z", text2: "z", expected: 1},
		{name: "long_lcs", text1: "AGGTAB", text2: "GXTXAYB", expected: 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, longestCommonSubsequence(tt.text1, tt.text2))
		})
	}
}

func Test_findLength(t *testing.T) {
	tests := []struct {
		name     string
		nums1    []int
		nums2    []int
		expected int
	}{
		{name: "example1", nums1: []int{1, 2, 3, 2, 1}, nums2: []int{3, 2, 1, 4, 7}, expected: 3},
		{name: "example2", nums1: []int{0, 0, 0, 0, 0}, nums2: []int{0, 0, 0, 0, 0}, expected: 5},
		{name: "no_common", nums1: []int{1, 2, 3}, nums2: []int{4, 5, 6}, expected: 0},
		{name: "single_match", nums1: []int{1, 2, 3}, nums2: []int{3, 4, 5}, expected: 1},
		{name: "identical_arrays", nums1: []int{1, 2, 3, 4}, nums2: []int{1, 2, 3, 4}, expected: 4},
		{name: "match_at_end", nums1: []int{1, 2, 3}, nums2: []int{4, 2, 3}, expected: 2},
		{name: "match_at_start", nums1: []int{1, 2, 3}, nums2: []int{1, 2, 5}, expected: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, findLength(tt.nums1, tt.nums2))
		})
	}
}
