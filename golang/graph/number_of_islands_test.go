package graph

import (
	"testing"

	"gocode/utils"

	"github.com/stretchr/testify/assert"
)

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
			assert.Equal(t, tt.expected, numIslandsDFS(utils.DeepCopyMatrix(tt.grid)))
			assert.Equal(t, tt.expected, numIslandsBFS(utils.DeepCopyMatrix(tt.grid)))
			assert.Equal(t, tt.expected, numIslandsUF(utils.DeepCopyMatrix(tt.grid)))
		})
	}
}
