package prefixtree

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_findWords(t *testing.T) {
	tests := []struct {
		name     string
		board    [][]byte
		words    []string
		expected []string
	}{
		{
			name: "leetcode_example1",
			board: [][]byte{
				{'o', 'a', 'a', 'n'},
				{'e', 't', 'a', 'e'},
				{'i', 'h', 'k', 'r'},
				{'i', 'f', 'l', 'v'},
			},
			words:    []string{"oath", "pea", "eat", "rain"},
			expected: []string{"eat", "oath"},
		},
		{
			name: "leetcode_example2",
			board: [][]byte{
				{'a', 'b'},
				{'c', 'd'},
			},
			words:    []string{"abcb"},
			expected: []string{},
		},
		{
			name: "single_cell_match",
			board: [][]byte{
				{'a'},
			},
			words:    []string{"a"},
			expected: []string{"a"},
		},
		{
			name: "single_cell_no_match",
			board: [][]byte{
				{'a'},
			},
			words:    []string{"b"},
			expected: []string{},
		},
		{
			name: "no_words_in_dict",
			board: [][]byte{
				{'a', 'b'},
				{'c', 'd'},
			},
			words:    []string{"xyz"},
			expected: []string{},
		},
		{
			name: "duplicate_word_only_once",
			board: [][]byte{
				{'a', 'a'},
			},
			words:    []string{"a"},
			expected: []string{"a"},
		},
		{
			name: "word_spans_full_board",
			board: [][]byte{
				{'a', 'b'},
				{'d', 'c'},
			},
			words:    []string{"abcd"},
			expected: []string{"abcd"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findWords(tt.board, tt.words)
			sort.Strings(result)
			sort.Strings(tt.expected)
			assert.Equal(t, tt.expected, result)
		})
	}
}
