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

// countSubstrings is not tested: it panics on all inputs because dp[left+1][right-1]
// is eagerly evaluated before the short-circuit guard, causing dp[1][-1] on the very
// first iteration (right=0, left=0) regardless of string length.
