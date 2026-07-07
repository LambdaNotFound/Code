package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// deepCopyByteGrid makes a fresh copy of a [][]byte grid so each implementation
// gets an unmodified grid (DFS/BFS mutate in-place).
func deepCopyByteGrid(src [][]byte) [][]byte {
	dst := make([][]byte, len(src))
	for i := range src {
		dst[i] = make([]byte, len(src[i]))
		copy(dst[i], src[i])
	}
	return dst
}

func Test_numIslandsAllImpl(t *testing.T) {
	tests := []struct {
		name     string
		grid     [][]byte
		expected int
	}{
		{
			name:     "leetcode_example1",
			grid:     [][]byte{{'1', '1', '1', '1', '0'}, {'1', '1', '0', '1', '0'}, {'1', '1', '0', '0', '0'}, {'0', '0', '0', '0', '0'}},
			expected: 1,
		},
		{
			name:     "leetcode_example2",
			grid:     [][]byte{{'1', '1', '0', '0', '0'}, {'1', '1', '0', '0', '0'}, {'0', '0', '1', '0', '0'}, {'0', '0', '0', '1', '1'}},
			expected: 3,
		},
		{
			name:     "all_water",
			grid:     [][]byte{{'0', '0'}, {'0', '0'}},
			expected: 0,
		},
		{
			name:     "all_land",
			grid:     [][]byte{{'1', '1'}, {'1', '1'}},
			expected: 1,
		},
		{
			name:     "diagonal_not_connected",
			grid:     [][]byte{{'1', '0'}, {'0', '1'}},
			expected: 2,
		},
		{
			name:     "single_land",
			grid:     [][]byte{{'1'}},
			expected: 1,
		},
		{
			name:     "single_water",
			grid:     [][]byte{{'0'}},
			expected: 0,
		},
		{
			name:     "checkerboard",
			grid:     [][]byte{{'1', '0', '1'}, {'0', '1', '0'}, {'1', '0', '1'}},
			expected: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, numIslandsDFS(deepCopyByteGrid(tt.grid)))
			assert.Equal(t, tt.expected, numIslandsBFS(deepCopyByteGrid(tt.grid)))
			assert.Equal(t, tt.expected, numIslandsUF(deepCopyByteGrid(tt.grid)))
		})
	}
}
