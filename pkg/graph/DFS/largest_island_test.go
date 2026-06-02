package dfs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_largestIsland(t *testing.T) {
	tests := []struct {
		name     string
		grid     [][]int
		expected int
	}{
		{
			name:     "single zero flip connects two islands",
			grid:     [][]int{{1, 0}, {0, 1}},
			expected: 3,
		},
		{
			name:     "no zero — entire grid is one island",
			grid:     [][]int{{1, 1}, {1, 1}},
			expected: 4,
		},
		{
			name:     "all zeros — flip one cell",
			grid:     [][]int{{0, 0}, {0, 0}},
			expected: 1,
		},
		{
			name: "flip fuses four distinct islands",
			grid: [][]int{
				{1, 0, 1},
				{0, 0, 0},
				{1, 0, 1},
			},
			expected: 3,
		},
		{
			name: "flip merges two large islands",
			grid: [][]int{
				{1, 1, 0, 1, 1},
				{1, 1, 0, 1, 1},
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
			},
			expected: 9,
		},
		{
			name:     "single cell island no zero",
			grid:     [][]int{{1}},
			expected: 1,
		},
		{
			name:     "single zero cell",
			grid:     [][]int{{0}},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, largestIsland(tt.grid))
		})
	}
}
