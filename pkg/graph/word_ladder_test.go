package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ladderLength(t *testing.T) {
    testCases := []struct {
        name     string
        beginWord string
        endWord string
        wordList    []string
        expected int
    }{
        {
            "case 1",
            "hit",
            "cog",
            []string{"hot","dot","dog","lot","log","cog"},
            5,
        },
        {
            "case 2",
            "hit",
            "cog",
            []string{"hot","dot","dog","lot","log"},
            0,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := ladderLength(tc.beginWord, tc.endWord, tc.wordList)
            assert.Equal(t, tc.expected, result)

            result = ladderLengthBiDirectionalBFS(tc.beginWord, tc.endWord, tc.wordList)
            assert.Equal(t, tc.expected, result)
        })
    }
}