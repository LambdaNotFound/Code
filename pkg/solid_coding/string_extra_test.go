package solid_coding

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_addBinary(t *testing.T) {
	tests := []struct {
		name     string
		a, b     string
		expected string
	}{
		{"leetcode_1", "11", "1", "100"},
		{"leetcode_2", "1010", "1011", "10101"},
		{"both_zero", "0", "0", "0"},
		{"one_empty_bit", "1", "0", "1"},
		{"carry_chain", "111", "1", "1000"},
		{"different_lengths", "1", "111", "1000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, addBinary(tt.a, tt.b))
		})
	}
}

// sortAnagramGroups sorts each inner slice and the outer slice so that
// groupAnagrams output can be compared deterministically.
func sortAnagramGroups(groups [][]string) [][]string {
	for _, g := range groups {
		sort.Strings(g)
	}
	sort.Slice(groups, func(i, j int) bool {
		if len(groups[i]) != len(groups[j]) {
			return len(groups[i]) < len(groups[j])
		}
		return groups[i][0] < groups[j][0]
	})
	return groups
}

func Test_groupAnagrams(t *testing.T) {
	tests := []struct {
		name     string
		strs     []string
		expected [][]string
	}{
		{
			"leetcode_1",
			[]string{"eat", "tea", "tan", "ate", "nat", "bat"},
			[][]string{{"bat"}, {"nat", "tan"}, {"ate", "eat", "tea"}},
		},
		{
			"single_empty",
			[]string{""},
			[][]string{{""}},
		},
		{
			"single_char",
			[]string{"a"},
			[][]string{{"a"}},
		},
		{
			"all_same",
			[]string{"ab", "ba", "ab"},
			[][]string{{"ab", "ab", "ba"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sortAnagramGroups(groupAnagrams(tt.strs))
			want := sortAnagramGroups(tt.expected)
			assert.Equal(t, want, got)
		})
	}
}

func Test_validPalindrome(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		expected bool
	}{
		{"leetcode_1", "aba", true},
		{"leetcode_2", "abca", true},
		{"leetcode_3", "abc", false},
		{"empty", "", true},
		{"single", "a", true},
		{"already_palindrome", "racecar", true},
		{"one_delete_mid", "abcbba", true},
		{"delete_first", "xabcba", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, validPalindrome(tt.s))
		})
	}
}

func Test_longestPalindrome(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		expected int
	}{
		{"leetcode_1", "abccccdd", 7},
		{"leetcode_2", "a", 1},
		{"all_same", "aaaa", 4},
		{"two_chars", "ab", 1},
		{"mixed_case", "Aa", 1},
		{"odd_all", "abc", 1},
		{"even_plus_one", "aabbcc", 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, longestPalindromeLength(tt.s), "longestPalindromeLength")
			assert.Equal(t, tt.expected, longestPalindrome(tt.s), "longestPalindrome")
		})
	}
}

func Test_longestCommonPrefix(t *testing.T) {
	tests := []struct {
		name     string
		strs     []string
		expected string
	}{
		{"leetcode_1", []string{"flower", "flow", "flight"}, "fl"},
		{"leetcode_2", []string{"dog", "racecar", "car"}, ""},
		{"all_same", []string{"abc", "abc", "abc"}, "abc"},
		{"single_string", []string{"hello"}, "hello"},
		{"empty_prefix", []string{"a", "b"}, ""},
		{"one_empty", []string{"", "abc"}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, longestCommonPrefix(tt.strs))
		})
	}
}

func Test_largestNumber(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected string
	}{
		{"leetcode_1", []int{10, 2}, "210"},
		{"leetcode_2", []int{3, 30, 34, 5, 9}, "9534330"},
		{"all_zeros", []int{0, 0}, "0"},
		{"single", []int{1}, "1"},
		{"same_values", []int{1, 1, 1}, "111"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, largestNumber(tt.nums))
		})
	}
}

func Test_removeDuplicates(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		expected string
	}{
		{"leetcode_1", "abbaca", "ca"},
		{"no_duplicates", "abc", "abc"},
		{"all_duplicates", "aabbcc", ""},
		{"cascade", "aabaa", "b"},
		{"single", "a", "a"},
		{"empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, removeDuplicates(tt.s))
		})
	}
}

func Test_removeKDuplicates(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		k        int
		expected string
	}{
		{"leetcode_1", "abcd", 2, "abcd"},
		{"leetcode_2", "deeedbbcccbdaa", 3, "aa"},
		{"leetcode_3", "pbbcggttciiippooaais", 2, "ps"},
		{"all_removed", "aaa", 3, ""},
		{"no_removal", "abc", 3, "abc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, removeKDuplicates(tt.s, tt.k))
		})
	}
}
