package dynamic_programming

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_uniquePaths(t *testing.T) {
	tests := []struct {
		name     string
		m, n     int
		expected int
	}{
		{name: "example1", m: 3, n: 7, expected: 28},
		{name: "example2", m: 3, n: 2, expected: 3},
		{name: "single_row", m: 1, n: 5, expected: 1},
		{name: "single_col", m: 5, n: 1, expected: 1},
		{name: "one_by_one", m: 1, n: 1, expected: 1},
		{name: "square_2x2", m: 2, n: 2, expected: 2},
		{name: "square_3x3", m: 3, n: 3, expected: 6},
		{name: "large", m: 7, n: 3, expected: 28},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, uniquePaths(tt.m, tt.n))
			assert.Equal(t, tt.expected, uniquePathsNaive(tt.m, tt.n))
		})
	}
}

func Test_minPathSum(t *testing.T) {
	tests := []struct {
		name     string
		grid     [][]int
		expected int
	}{
		{
			name:     "example1",
			grid:     [][]int{{1, 3, 1}, {1, 5, 1}, {4, 2, 1}},
			expected: 7,
		},
		{
			name:     "example2",
			grid:     [][]int{{1, 2, 3}, {4, 5, 6}},
			expected: 12,
		},
		{
			name:     "one_by_one",
			grid:     [][]int{{5}},
			expected: 5,
		},
		{
			name:     "single_row",
			grid:     [][]int{{1, 2, 3}},
			expected: 6,
		},
		{
			name:     "single_col",
			grid:     [][]int{{1}, {2}, {3}},
			expected: 6,
		},
		{
			name:     "all_zeros",
			grid:     [][]int{{0, 0}, {0, 0}},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// minPathSumTopDown does not mutate the grid
			assert.Equal(t, tt.expected, minPathSumTopDown(tt.grid))

			// minPathSumBottomUp mutates the grid in-place; deep-copy first
			gridCopy := make([][]int, len(tt.grid))
			for i, row := range tt.grid {
				gridCopy[i] = make([]int, len(row))
				copy(gridCopy[i], row)
			}
			assert.Equal(t, tt.expected, minPathSumBottomUp(gridCopy))
		})
	}
}
