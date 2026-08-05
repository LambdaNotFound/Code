package graph

import (
	"testing"

	"gocode/utils"

	"github.com/stretchr/testify/assert"
)

func Test_islandPerimeter(t *testing.T) {
	tests := []struct {
		name     string
		grid     [][]int
		expected int
	}{
		{
			name:     "leetcode_example",
			grid:     [][]int{{0, 1, 0, 0}, {1, 1, 1, 0}, {0, 1, 0, 0}, {1, 1, 0, 0}},
			expected: 16,
		},
		{
			name:     "single_land_cell",
			grid:     [][]int{{1}},
			expected: 4,
		},
		{
			name:     "all_water",
			grid:     [][]int{{0, 0}, {0, 0}},
			expected: 0,
		},
		{
			name:     "row_of_three",
			grid:     [][]int{{1, 1, 1}},
			expected: 8,
		},
		{
			name:     "2x2_solid",
			grid:     [][]int{{1, 1}, {1, 1}},
			expected: 8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g1 := utils.DeepCopyMatrix(tt.grid)
			assert.Equal(t, tt.expected, islandPerimeter(g1))
			assert.Equal(t, tt.expected, islandPerimeterClaude(tt.grid))
		})
	}
}
