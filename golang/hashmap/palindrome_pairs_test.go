package hashmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_isPalindrome(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		expected bool
	}{
		{name: "empty", s: "", expected: true},
		{name: "single_char", s: "a", expected: true},
		{name: "two_same", s: "aa", expected: true},
		{name: "two_diff", s: "ab", expected: false},
		{name: "odd_palindrome", s: "racecar", expected: true},
		{name: "even_palindrome", s: "abba", expected: true},
		{name: "not_palindrome", s: "hello", expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, isPalindrome(tt.s))
		})
	}
}

func Test_reverse(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		expected string
	}{
		{name: "empty", s: "", expected: ""},
		{name: "single", s: "a", expected: "a"},
		{name: "two", s: "ab", expected: "ba"},
		{name: "word", s: "hello", expected: "olleh"},
		{name: "palindrome", s: "racecar", expected: "racecar"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, reverse(tt.s))
		})
	}
}

func Test_palindromePairs(t *testing.T) {
	tests := []struct {
		name     string
		words    []string
		expected [][]int
	}{
		{
			name:     "leetcode_example1",
			words:    []string{"abcd", "dcba", "lls", "s", "sssll"},
			expected: [][]int{{0, 1}, {1, 0}, {3, 2}, {2, 4}},
		},
		{
			name:     "leetcode_example2",
			words:    []string{"bat", "tab", "cat"},
			expected: [][]int{{0, 1}, {1, 0}},
		},
		{
			name:     "with_empty_string",
			words:    []string{"a", ""},
			expected: [][]int{{0, 1}, {1, 0}},
		},
		{
			name:     "all_same_palindromes",
			words:    []string{"aa", "bb"},
			expected: [][]int{},
		},
		{
			name:     "single_word",
			words:    []string{"abc"},
			expected: [][]int{},
		},
		{
			// Case 4: suffix is palindrome, reverse(prefix) found in list
			// "abcxyx": suffix "xyx" is palindrome, prefix "abc" → reverse is "cba" (index 1)
			// → "abcxyx"+"cba" = "abcxyxcba" ✓
			name:     "case4_palindrome_suffix",
			words:    []string{"abcxyx", "cba"},
			expected: [][]int{{0, 1}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := palindromePairs(tt.words)
			assert.ElementsMatch(t, tt.expected, got)
		})
	}
}
