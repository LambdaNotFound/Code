package two_pointers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_lengthOfLongestSubstring(t *testing.T) {
    testCases := []struct {
        name     string
        str      string
        expected int
    }{
        {
            "case 1",
            "abcabcbb",
            3,
        },
        {
            "case 2",
            "bbbbb",
            1,
        },
        {
            "case 3",
            "pwwkew",
            3,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := lengthOfLongestSubstring(tc.str)
            assert.Equal(t, tc.expected, result)

            result = lengthOfLongestSubstring_optimized(tc.str)
            assert.Equal(t, tc.expected, result)
        })
    }
}
