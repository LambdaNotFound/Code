package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// isValidAlienOrder checks whether `result` is a valid topological order
// consistent with the ordering constraints implied by `words`.
// This lets us test both BFS and DFS implementations, which may produce
// different valid orderings.
func isValidAlienOrder(result string, words []string) bool {
	if result == "" {
		return false
	}
	pos := make(map[byte]int, len(result))
	for i := 0; i < len(result); i++ {
		pos[result[i]] = i
	}
	for i := 0; i+1 < len(words); i++ {
		w1, w2 := words[i], words[i+1]
		minLen := len(w1)
		if len(w2) < minLen {
			minLen = len(w2)
		}
		for j := 0; j < minLen; j++ {
			if w1[j] != w2[j] {
				if pos[w1[j]] >= pos[w2[j]] {
					return false
				}
				break
			}
		}
	}
	return true
}

func Test_foreignDictionary(t *testing.T) {
	tests := []struct {
		name        string
		words       []string
		expectEmpty bool // true when the function must return ""
	}{
		{name: "leetcode_example1", words: []string{"wrt", "wrf", "er", "ett", "rftt"}, expectEmpty: false},
		{name: "leetcode_example2", words: []string{"z", "x"}, expectEmpty: false},
		{name: "cycle_invalid", words: []string{"z", "x", "z"}, expectEmpty: true},
		{name: "prefix_invalid", words: []string{"abc", "ab"}, expectEmpty: true},
		{name: "single_word", words: []string{"abc"}, expectEmpty: false},
		{name: "all_same", words: []string{"aa", "aa"}, expectEmpty: false},
		{name: "no_ordering_needed", words: []string{"a", "b", "c"}, expectEmpty: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBFS := foreignDictionary(tt.words)
			gotDFS := foreignDictionaryDFS(tt.words)

			if tt.expectEmpty {
				assert.Equal(t, "", gotBFS)
				assert.Equal(t, "", gotDFS)
			} else {
				assert.True(t, isValidAlienOrder(gotBFS, tt.words), "BFS result %q is not a valid order", gotBFS)
				assert.True(t, isValidAlienOrder(gotDFS, tt.words), "DFS result %q is not a valid order", gotDFS)
			}
		})
	}
}
