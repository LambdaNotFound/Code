package string

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
		{"sum3_case", "11", "11", "110"}, // bit-1: 1+1+carry(1)=3 → triggers case 3
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

func Test_isPalindrome(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		expected bool
	}{
		{"leetcode_example1", "A man, a plan, a canal: Panama", true},
		{"leetcode_example2", "race a car", false},
		{"empty_string", "", true},
		{"single_char", "a", true},
		{"spaces_only", "   ", true},
		{"digits_palindrome", "12321", true},
		{"mixed_alphanumeric", "0P", false},
		{"uppercase_lowercase", "AbBa", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, isPalindrome(tt.s))
		})
	}
}

func Test_validPalindromeAtMostK(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		k        int
		expected bool
	}{
		{"k0_palindrome", "racecar", 0, true},
		{"k0_not_palindrome", "abcd", 0, false},
		{"k1_one_delete", "abca", 1, true},
		{"k1_not_fixable", "abcd", 1, false},
		{"k2_two_deletes", "abcda", 2, true},
		{"k_ge_len_always_true", "xyz", 3, true},
		{"already_palindrome", "aba", 1, true},
		{"single_char", "a", 0, true},
		{"two_chars_same", "aa", 0, true},
		{"two_chars_diff_k1", "ab", 1, true},
		// longer string forces memo lookups (overlapping subproblems)
		{"long_needs_k2", "abcddcba", 0, true},
		{"long_not_fixable", "abcdefgh", 2, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, validPalindromeAtMostK(tt.s, tt.k))
		})
	}
}

func Test_romanToInt(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		expected int
	}{
		{"III", "III", 3},
		{"LVIII", "LVIII", 58},
		{"MCMXCIV", "MCMXCIV", 1994},
		{"IV", "IV", 4},
		{"IX", "IX", 9},
		{"XL", "XL", 40},
		{"XC", "XC", 90},
		{"CD", "CD", 400},
		{"CM", "CM", 900},
		{"single_M", "M", 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, romanToInt(tt.s))
		})
	}
}

// Test_myAtoi_extra adds cases that exercise the uncovered branches:
//   - plus sign prefix (s[0]=='+' branch)
//   - int32 overflow in both directions
func Test_myAtoi_extra(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		expected int
	}{
		{"plus_sign", "+100", 100},
		{"empty_string", "", 0},
		{"only_spaces", "   ", 0},
		{"int32_max_overflow", "9999999999", 2147483647},
		{"int32_min_overflow", "-9999999999", -2147483648},
		{"leading_zeros", "007", 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, myAtoi(tt.s))
		})
	}
}
