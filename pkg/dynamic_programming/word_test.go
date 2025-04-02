package dynamic_programming

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_wordBreak(t *testing.T) {
    testCases := []struct {
        name           string
        input_string   string
        input_wordDict []string
        expected       bool
    }{
        {
            "case 1",
            "leetcode",
            []string{"leet", "code"},
            true,
        },
        {
            "case 2",
            "applepenapple",
            []string{"apple", "pen"},
            true,
        },
        {
            "case 3",
            "catsandog",
            []string{"cats", "dog", "sand", "and", "cat"},
            false,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := wordBreak(tc.input_string, tc.input_wordDict)
            assert.Equal(t, tc.expected, result)
        })
    }
}
