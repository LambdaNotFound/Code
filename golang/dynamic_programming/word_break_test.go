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
        {"case 1", "leetcode", []string{"leet", "code"}, true},
        {"case 2", "applepenapple", []string{"apple", "pen"}, true},
        {"case 3", "catsandog", []string{"cats", "dog", "sand", "and", "cat"}, false},

        {
            name:           "single_word_match",
            input_string:   "hello",
            input_wordDict: []string{"hello"},
            expected:       true,
        },
        {
            name:           "no_match",
            input_string:   "hello",
            input_wordDict: []string{"world"},
            expected:       false,
        },
        {
            name:           "repeated_word",
            input_string:   "aaaa",
            input_wordDict: []string{"a", "aa"},
            expected:       true,
        },
        {
            name:           "word_used_multiple_times",
            input_string:   "appleapple",
            input_wordDict: []string{"apple"},
            expected:       true,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            assert.Equal(t, tc.expected, wordBreak(tc.input_string, tc.input_wordDict))
            assert.Equal(t, tc.expected, wordBreakTrie(tc.input_string, tc.input_wordDict))
            assert.Equal(t, tc.expected, wordBreakCharIndex(tc.input_string, tc.input_wordDict))
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
            []string{"apple", "pen", "applepen", "pine", "pineapple"},
            []string{"pine apple pen apple", "pineapple pen apple", "pine applepen apple"},
        },
        {
            "case 3",
            "catsandog",
            []string{"cats", "dog", "sand", "and", "cat"},
            []string{},
        },
        {
            name:           "single_word",
            input_string:   "hello",
            input_wordDict: []string{"hello"},
            expected:       []string{"hello"},
        },
        {
            name:           "repeated_word",
            input_string:   "aa",
            input_wordDict: []string{"a"},
            expected:       []string{"a a"},
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := wordBreak2(tc.input_string, tc.input_wordDict)
            assert.ElementsMatch(t, tc.expected, result)
        })
    }
}

func Test_wordBreakNoReuse(t *testing.T) {
    testCases := []struct {
        name           string
        input_string   string
        input_wordDict []string
        expected       bool
    }{
        {"case 1", "leetcode", []string{"leet", "code"}, true},
        {"case 2", "catsandog", []string{"cats", "dog", "sand", "and", "cat"}, false},
        {
            name:           "single_word_match",
            input_string:   "hello",
            input_wordDict: []string{"hello"},
            expected:       true,
        },
        {
            name:           "no_match",
            input_string:   "hello",
            input_wordDict: []string{"world"},
            expected:       false,
        },
        {
            name:           "word_used_multiple_times_unavailable",
            input_string:   "appleapple",
            input_wordDict: []string{"apple"},
            expected:       false,
        },
        {
            name:           "word_used_multiple_times_available",
            input_string:   "appleapple",
            input_wordDict: []string{"apple", "apple"},
            expected:       true,
        },
        {
            name:           "prefix_word_consumed_needed_later",
            input_string:   "abba",
            input_wordDict: []string{"ab", "a", "abb"},
            expected:       true,
        },
        {
            name:           "backback_single_back_unavailable",
            input_string:   "backback",
            input_wordDict: []string{"back"},
            expected:       false,
        },
        {
            name:           "backend_with_extra_dict_entries",
            input_string:   "backend",
            input_wordDict: []string{"back", "end", "backe", "front", "start", "back", "nds", "backend"},
            expected:       true,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            assert.Equal(t, tc.expected, wordBreakNoReuse(tc.input_string, tc.input_wordDict))
        })
    }
}
