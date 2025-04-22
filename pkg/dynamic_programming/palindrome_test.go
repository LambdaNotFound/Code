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
        {
            "case 1",
            "babad",
            "bab",
        },
        {
            "case 2",
            "cbbd",
            "bb",
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := longestPalindrome(tc.input)
            assert.Equal(t, tc.expected, result)

            result = longestPalindrome_optimized(tc.input)
            assert.Equal(t, tc.expected, result)
        })
    }
}

func Test_longestPalindromeLength(t *testing.T) {
    testCases := []struct {
        name     string
        input    string
        expected int
    }{
        {
            "case 1",
            "abccccdd",
            7,
        },
        {
            "case 2",
            "a",
            1,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := longestPalindromeLength(tc.input)
            assert.Equal(t, tc.expected, result)
        })
    }
}
