package backtracking

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_exist(t *testing.T) {
    testCases := []struct {
        name     string
        board    [][]byte
        word     string
        expected bool
    }{
        {
            "case 1",
            [][]byte{
                []byte{'A', 'B', 'C', 'E'},
                []byte{'S', 'F', 'C', 'S'},
                []byte{'A', 'D', 'E', 'E'},
            },
            "ABCCED",
            true,
        },
        {
            "case 2",
            [][]byte{
                []byte{'A', 'B', 'C', 'E'},
                []byte{'S', 'F', 'C', 'S'},
                []byte{'A', 'D', 'E', 'E'},
            },
            "SEE",
            true,
        },
        {
            "case 3",
            [][]byte{
                []byte{'A', 'B', 'C', 'E'},
                []byte{'S', 'F', 'C', 'S'},
                []byte{'A', 'D', 'E', 'E'},
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

func Test_solveSudoku(t *testing.T) {
    testCases := []struct {
        name     string
        board    [][]byte
        expected [][]byte
    }{
        {
            "case 1",
            [][]byte{
                []byte{'5', '3', '.', '.', '7', '.', '.', '.', '.'},
                []byte{'6', '.', '.', '1', '9', '5', '.', '.', '.'},
                []byte{'.', '9', '8', '.', '.', '.', '.', '6', '.'},

                []byte{'8', '.', '.', '.', '6', '.', '.', '.', '3'},
                []byte{'4', '.', '.', '8', '.', '3', '.', '.', '1'},
                []byte{'7', '.', '.', '.', '2', '.', '.', '.', '6'},

                []byte{'.', '6', '.', '.', '.', '.', '2', '8', '.'},
                []byte{'.', '.', '.', '4', '1', '9', '.', '.', '5'},
                []byte{'.', '.', '.', '.', '8', '.', '.', '7', '9'},
            },
            [][]byte{
                []byte{'5', '3', '4', '6', '7', '8', '9', '1', '2'},
                []byte{'6', '7', '2', '1', '9', '5', '3', '4', '8'},
                []byte{'1', '9', '8', '3', '4', '2', '5', '6', '7'},

                []byte{'8', '5', '9', '7', '6', '1', '4', '2', '3'},
                []byte{'4', '2', '6', '8', '5', '3', '7', '9', '1'},
                []byte{'7', '1', '3', '9', '2', '4', '8', '5', '6'},

                []byte{'9', '6', '1', '5', '3', '7', '2', '8', '4'},
                []byte{'2', '8', '7', '4', '1', '9', '6', '3', '5'},
                []byte{'3', '4', '5', '2', '8', '6', '1', '7', '9'},
            },
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            solveSudoku(tc.board)
            assert.Equal(t, tc.expected, tc.board)
        })
    }
}

func Test_solveNQueens(t *testing.T) {
    testCases := []struct {
        name     string
        n        int
        expected [][]string
    }{
        {
            "case 1",
            4,
            [][]string{
                []string{".Q..", "...Q", "Q...", "..Q."},
                []string{"..Q.", "Q...", "...Q", ".Q.."},
            },
        },
        {
            "case 2",
            1,
            [][]string{
                []string{"Q"},
            },
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := solveNQueens(tc.n)
            assert.Equal(t, tc.expected, result)
        })
    }
}
