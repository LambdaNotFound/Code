package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ladderLength(t *testing.T) {
    testCases := []struct {
        name      string
        beginWord string
        endWord   string
        wordList  []string
        expected  int
    }{
        {
            "case 1",
            "hit",
            "cog",
            []string{"hot", "dot", "dog", "lot", "log", "cog"},
            5,
        },
        {
            "case 2",
            "hit",
            "cog",
            []string{"hot", "dot", "dog", "lot", "log"},
            0,
        },
        {
            "case 3",
            "ymain",
            "oecij",
            []string{"ymann", "yycrj", "oecij", "ymcnj", "yzcrj", "yycij", "xecij", "yecij", "ymanj", "yzcnj", "ymain"},
            10,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := ladderLength(tc.beginWord, tc.endWord, tc.wordList)
            assert.Equal(t, tc.expected, result)

            result = ladderLengthBidirectionalBFS(tc.beginWord, tc.endWord, tc.wordList)
            assert.Equal(t, tc.expected, result)
        })
    }
}


func TestLadderLength(t *testing.T) {
    testCases := []struct {
        name     string
        begin    string
        end      string
        wordList []string
        expected int
    }{
        {
            name:     "leetcode example 1",
            begin:    "hit",
            end:      "cog",
            wordList: []string{"hot", "dot", "dog", "lot", "log", "cog"},
            expected: 5, // hit -> hot -> dot -> dog -> cog
        },
        {
            name:     "leetcode example 2 (no path)",
            begin:    "hit",
            end:      "cog",
            wordList: []string{"hot", "dot", "dog", "lot", "log"}, // no "cog"
            expected: 0,
        },
        {
            name:     "begin equals end",
            begin:    "a",
            end:      "a",
            wordList: []string{"a"},
            expected: 1, // can consider 1-step sequence (problem treats begin=end case ambiguously but this works for testing)
        },
        {
            name:     "immediate neighbor",
            begin:    "a",
            end:      "c",
            wordList: []string{"a", "b", "c"},
            expected: 2, // a -> c
        },
        {
            name:     "unreachable end due to isolated word",
            begin:    "hit",
            end:      "cog",
            wordList: []string{"hot", "dot", "tod", "hog", "hop"}, // no chain reaching cog
            expected: 0,
        },
        {
            name:     "multi-branch shortest path",
            begin:    "talk",
            end:      "tail",
            wordList: []string{"tall", "tail", "balk", "tark"},
            expected: 3, // talk -> tall -> tail
        },
        {
            name:     "single-letter transformations",
            begin:    "a",
            end:      "c",
            wordList: []string{"a", "b", "c"},
            expected: 2, // a -> c
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            got := ladderLength(tc.begin, tc.end, tc.wordList)
            assert.Equal(t, tc.expected, got)
        })
    }
}