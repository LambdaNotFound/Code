package backtracking

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ladderLength(t *testing.T) {
    testCases := []struct {
        name      string
        board [][]byte
        word string
        expected  bool
    }{
        {
            "case 1",
            [][]byte{
                []byte{'A','B','C','E'},
                []byte{'S','F','C','S'},
                []byte{'A','D','E','E'},
            },
            "ABCCED",
            true,
        },
        {
            "case 2",
            [][]byte{
                []byte{'A','B','C','E'},
                []byte{'S','F','C','S'},
                []byte{'A','D','E','E'},
            },
            "SEE",
            true,
        },
        {
            "case 3",
            [][]byte{
                []byte{'A','B','C','E'},
                []byte{'S','F','C','S'},
                []byte{'A','D','E','E'},
            },
            "ABCB",
            false,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := exist(tc.board, tc.word)
            assert.Equal(t, tc.expected, result)
        })
    }
}
