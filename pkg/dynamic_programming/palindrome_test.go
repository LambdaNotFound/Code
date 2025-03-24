package palindrome

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

            result = longestPalindrome_opt_space(tc.input)
            assert.Equal(t, tc.expected, result)
        })
    }
}
