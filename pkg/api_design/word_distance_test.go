package apidesign

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_WordDistance(t *testing.T) {
	wordsDict := []string{"practice", "makes", "perfect", "coding", "makes"}
	wd := NewWordDistance(wordsDict)

	assert.Equal(t, 3, wd.shortest("coding", "practice"))
	assert.Equal(t, 1, wd.shortest("makes", "coding"))
}

func Test_WordDistance_adjacentDuplicates(t *testing.T) {
	// word1 and word2 appear right next to each other
	wordsDict := []string{"a", "b", "a", "b"}
	wd := NewWordDistance(wordsDict)

	assert.Equal(t, 1, wd.shortest("a", "b"))
}

func Test_WordDistance_repeatedQueries(t *testing.T) {
	// calling shortest multiple times on the same instance
	wordsDict := []string{"a", "c", "b", "b", "a"}
	wd := NewWordDistance(wordsDict)

	assert.Equal(t, 1, wd.shortest("a", "b"))
	assert.Equal(t, 1, wd.shortest("a", "b")) // same query again
	assert.Equal(t, 1, wd.shortest("a", "c"))
}

func Test_WordDistance_singleOccurrence(t *testing.T) {
	// each word appears exactly once
	wordsDict := []string{"x", "y", "z"}
	wd := NewWordDistance(wordsDict)

	assert.Equal(t, 1, wd.shortest("x", "y"))
	assert.Equal(t, 2, wd.shortest("x", "z"))
	assert.Equal(t, 1, wd.shortest("y", "z"))
}
