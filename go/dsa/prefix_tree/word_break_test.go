package prefixtree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_wordBreakTrie(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		wordDict []string
		expected bool
	}{
		{
			name:     "leetcode_example1",
			s:        "leetcode",
			wordDict: []string{"leet", "code"},
			expected: true,
		},
		{
			name:     "leetcode_example2",
			s:        "applepenapple",
			wordDict: []string{"apple", "pen"},
			expected: true,
		},
		{
			name:     "leetcode_example3",
			s:        "catsandog",
			wordDict: []string{"cats", "dog", "sand", "and", "cat"},
			expected: false,
		},
		{
			name:     "single_word_match",
			s:        "hello",
			wordDict: []string{"hello"},
			expected: true,
		},
		{
			name:     "single_word_no_match",
			s:        "hello",
			wordDict: []string{"world"},
			expected: false,
		},
		{
			name:     "overlapping_prefixes",
			s:        "aaab",
			wordDict: []string{"a", "aa", "aaa"},
			expected: false,
		},
		{
			name:     "repeated_word",
			s:        "aaa",
			wordDict: []string{"a", "aa"},
			expected: true,
		},
		{
			name:     "prefix_not_in_dict",
			s:        "ab",
			wordDict: []string{"a", "b"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, wordBreakTrie(tt.s, tt.wordDict))
		})
	}
}
