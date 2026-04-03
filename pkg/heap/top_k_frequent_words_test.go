package heap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_topKFrequentWords(t *testing.T) {
	tests := []struct {
		name     string
		words    []string
		k        int
		expected []string
	}{
		{
			name:     "leetcode_example1",
			words:    []string{"i", "love", "leetcode", "i", "love", "coding"},
			k:        2,
			expected: []string{"i", "love"},
		},
		{
			name:     "leetcode_example2",
			words:    []string{"the", "day", "is", "sunny", "the", "the", "the", "sunny", "is", "is"},
			k:        4,
			expected: []string{"the", "is", "sunny", "day"},
		},
		{
			name:     "single_word",
			words:    []string{"a"},
			k:        1,
			expected: []string{"a"},
		},
		{
			name:     "k_equals_1_clear_winner",
			words:    []string{"apple", "banana", "apple", "apple"},
			k:        1,
			expected: []string{"apple"},
		},
		{
			name:     "lex_tiebreak",
			words:    []string{"a", "b", "c", "a", "b", "c"},
			k:        2,
			expected: []string{"a", "b"},
		},
		{
			name:     "all_unique",
			words:    []string{"z", "m", "a"},
			k:        3,
			expected: []string{"a", "m", "z"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := topKFrequentWords(tt.words, tt.k)
			assert.Equal(t, tt.expected, result)
		})
	}
}
