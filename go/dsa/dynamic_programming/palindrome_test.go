package dynamic_programming

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_longestPalindrome(t *testing.T) {
    testCases := []struct {
        name     string
        input    string
        expected string
    }{
        {"case 1", "babad", "bab"},
        {"case 2", "cbbd", "bb"},
        {name: "single_char", input: "a", expected: "a"},
        {name: "all_same", input: "aaaa", expected: "aaaa"},
        {name: "entire_string", input: "racecar", expected: "racecar"},
        {name: "no_palindrome_longer_than_1", input: "abcd", expected: "a"},
        {name: "even_palindrome", input: "aabb", expected: "aa"},
        // "xabacabay": "abacaba" is a 7-char palindrome in the middle
        {name: "palindrome_in_middle", input: "xabacabay", expected: "abacaba"},
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            assert.Equal(t, tc.expected, longestPalindrome(tc.input))
            assert.Equal(t, tc.expected, longestPalindrome_optimized(tc.input))
        })
    }
}

func Test_countSubstrings(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{name: "single_char", input: "a", expected: 1},
		{name: "two_same", input: "aa", expected: 3},
		{name: "two_diff", input: "ab", expected: 2},
		{name: "leetcode_example1", input: "abc", expected: 3},
		{name: "leetcode_example2", input: "aaa", expected: 6},
		{name: "palindrome_in_middle", input: "aba", expected: 4},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, countSubstrings(tc.input))
		})
	}
}
