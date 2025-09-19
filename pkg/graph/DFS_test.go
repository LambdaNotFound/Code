package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNumIslands(t *testing.T) {
    tests := []struct {
        name     string
        grid     [][]byte
        expected int
    }{
        {
            name:     "empty grid",
            grid:     [][]byte{},
            expected: 0,
        },
        {
            name: "single island",
            grid: [][]byte{
                {'1', '1', '0', '0', '0'},
                {'1', '1', '0', '0', '0'},
                {'0', '0', '0', '1', '1'},
                {'0', '0', '0', '1', '1'},
            },
            expected: 2,
        },
        {
            name: "all water",
            grid: [][]byte{
                {'0', '0'},
                {'0', '0'},
            },
            expected: 0,
        },
        {
            name: "all land",
            grid: [][]byte{
                {'1', '1'},
                {'1', '1'},
            },
            expected: 1,
        },
        {
            name: "diagonal not connected",
            grid: [][]byte{
                {'1', '0'},
                {'0', '1'},
            },
            expected: 2, // diagonals don’t connect
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            got := numIslands(tc.grid)
            assert.Equal(t, tc.expected, got)
        })
    }
}
