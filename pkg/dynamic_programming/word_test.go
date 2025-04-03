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

func Test_wordBreak2(t *testing.T) {
    testCases := []struct {
        name           string
        input_string   string
        input_wordDict []string
        expected       []string
    }{
        {
            "case 1",
            "catsanddog",
            []string{"cat", "cats", "and", "sand", "dog"},
            []string{"cat sand dog", "cats and dog"},
        },
        {
            "case 2",
            "pineapplepenapple",
            []string{"apple","pen","applepen","pine","pineapple"},
            []string{"pine apple pen apple","pineapple pen apple","pine applepen apple"},
        },
        {
            "case 3",
            "catsandog",
            []string{"cats","dog","sand","and","cat"},
            []string{},
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := wordBreak2(tc.input_string, tc.input_wordDict)
            assert.ElementsMatch(t, tc.expected, result)
        })
    }
}
