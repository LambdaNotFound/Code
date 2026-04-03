package backtracking

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_existClaude(t *testing.T) {
	testCases := []struct {
		name     string
		board    [][]byte
		word     string
		expected bool
	}{
		{
			name: "found_zigzag",
			board: [][]byte{
				{'A', 'B', 'C', 'E'},
				{'S', 'F', 'C', 'S'},
				{'A', 'D', 'E', 'E'},
			},
			word:     "ABCCED",
			expected: true,
		},
		{
			name: "found_straight",
			board: [][]byte{
				{'A', 'B', 'C', 'E'},
				{'S', 'F', 'C', 'S'},
				{'A', 'D', 'E', 'E'},
			},
			word:     "SEE",
			expected: true,
		},
		{
			name: "not_found_reuse",
			board: [][]byte{
				{'A', 'B', 'C', 'E'},
				{'S', 'F', 'C', 'S'},
				{'A', 'D', 'E', 'E'},
			},
			word:     "ABCB",
			expected: false,
		},
		{
			name:     "single_cell_match",
			board:    [][]byte{{'Z'}},
			word:     "Z",
			expected: true,
		},
		{
			name:     "single_cell_no_match",
			board:    [][]byte{{'Z'}},
			word:     "A",
			expected: false,
		},
		{
			name: "vertical_path",
			board: [][]byte{
				{'A'},
				{'B'},
				{'C'},
			},
			word:     "ABC",
			expected: true,
		},
		{
			name:     "word_longer_than_grid",
			board:    [][]byte{{'A', 'B'}, {'C', 'D'}},
			word:     "ABCDE",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := existClaude(tc.board, tc.word)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func Test_findWords(t *testing.T) {
	testCases := []struct {
		name     string
		board    [][]byte
		words    []string
		expected []string
	}{
		{
			name: "leetcode_example_1",
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
			name: "leetcode_example_2",
			board: [][]byte{
				{'a', 'b'},
				{'c', 'd'},
			},
			words:    []string{"abcb"},
			expected: []string{},
		},
		{
			name: "single_word_found",
			board: [][]byte{
				{'a', 'b', 'c'},
				{'d', 'e', 'f'},
			},
			words:    []string{"abc"},
			expected: []string{"abc"},
		},
		{
			name: "duplicate_word_in_list_found_once",
			board: [][]byte{
				{'a', 'b'},
				{'c', 'd'},
			},
			words:    []string{"ab", "ab"},
			expected: []string{"ab"},
		},
		{
			name: "no_words_found",
			board: [][]byte{
				{'a', 'b'},
			},
			words:    []string{"xyz"},
			expected: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := findWords(tc.board, tc.words)
			assert.ElementsMatch(t, tc.expected, result)
		})
	}
}
