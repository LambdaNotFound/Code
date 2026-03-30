package prefixtree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWordTrie_InsertSearch(t *testing.T) {
	nilWords := ([]string)(nil)

	testCases := []struct {
		name    string
		inserts [][]string
		queries []struct {
			words    []string
			expected bool
		}
	}{
		{
			name:    "single path exact match",
			inserts: [][]string{{"leet", "code"}},
			queries: []struct {
				words    []string
				expected bool
			}{
				{words: []string{"leet", "code"}, expected: true},
				{words: []string{"leet"}, expected: false},        // prefix only
				{words: []string{"leet", "cod"}, expected: false}, // wrong token
				{words: nilWords, expected: false},                // empty query not inserted
				{words: []string{}, expected: false},              // empty query not inserted
			},
		},
		{
			name:    "empty sequence insertion",
			inserts: [][]string{nilWords},
			queries: []struct {
				words    []string
				expected bool
			}{
				{words: nilWords, expected: true},
				{words: []string{}, expected: true},
				{words: []string{"a"}, expected: false},
			},
		},
		{
			name:    "multiple branches",
			inserts: [][]string{{"a", "b"}, {"a", "c"}, {"d"}},
			queries: []struct {
				words    []string
				expected bool
			}{
				{words: []string{"a", "b"}, expected: true},
				{words: []string{"a", "c"}, expected: true},
				{words: []string{"d"}, expected: true},
				{words: []string{"a"}, expected: false},           // not inserted as terminal
				{words: []string{"a", "b", "x"}, expected: false}, // longer than inserted
				{words: []string{"x"}, expected: false},
			},
		},
		{
			name:    "overlapping sequences",
			inserts: [][]string{{"a", "b"}, {"a"}},
			queries: []struct {
				words    []string
				expected bool
			}{
				{words: []string{"a"}, expected: true},
				{words: []string{"a", "b"}, expected: true},
				{words: []string{"a", "c"}, expected: false},
				{words: []string{"b", "a"}, expected: false},
			},
		},
		{
			name:    "repeated insert does not break",
			inserts: [][]string{{"x", "y"}, {"x", "y"}},
			queries: []struct {
				words    []string
				expected bool
			}{
				{words: []string{"x", "y"}, expected: true},
				{words: []string{"x"}, expected: false},
				{words: []string{"x", "y", "z"}, expected: false},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			trie := ConstructorWordTrie()
			for _, seq := range tc.inserts {
				trie.Insert(seq)
			}

			for _, q := range tc.queries {
				assert.Equal(t, q.expected, trie.Search(q.words))
			}
		})
	}
}
