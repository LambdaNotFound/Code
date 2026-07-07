package multisource

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

// sortCells normalises a [][]int of [row,col] pairs so comparison is order-independent.
func sortCells(cells [][]int) [][]int {
	sort.Slice(cells, func(i, j int) bool {
		if cells[i][0] != cells[j][0] {
			return cells[i][0] < cells[j][0]
		}
		return cells[i][1] < cells[j][1]
	})
	return cells
}

func Test_pacificAtlantic(t *testing.T) {
	tests := []struct {
		name     string
		heights  [][]int
		expected [][]int
	}{
		{
			name: "leetcode_example1",
			heights: [][]int{
				{1, 2, 2, 3, 5},
				{3, 2, 3, 4, 4},
				{2, 4, 5, 3, 1},
				{6, 7, 1, 4, 5},
				{5, 1, 1, 2, 4},
			},
			expected: [][]int{{0, 4}, {1, 3}, {1, 4}, {2, 2}, {3, 0}, {3, 1}, {4, 0}},
		},
		{
			name:     "single_cell",
			heights:  [][]int{{1}},
			expected: [][]int{{0, 0}},
		},
		{
			name:     "flat_grid_all_flow",
			heights:  [][]int{{1, 1}, {1, 1}},
			expected: [][]int{{0, 0}, {0, 1}, {1, 0}, {1, 1}},
		},
		{
			name: "increasing_from_top_left",
			heights: [][]int{
				{1, 2, 3},
				{4, 5, 6},
				{7, 8, 9},
			},
			// only bottom-right corner (2,2)=9 can reach both; and top row, right col corners
			expected: [][]int{{0, 2}, {1, 2}, {2, 0}, {2, 1}, {2, 2}},
		},
		{
			name:     "empty_grid",
			heights:  [][]int{},
			expected: [][]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDFS := sortCells(pacificAtlanticDFS(tt.heights))
			gotBFS := sortCells(pacificAtlanticBFS(tt.heights))
			want := sortCells(tt.expected)
			assert.Equal(t, want, gotDFS)
			assert.Equal(t, want, gotBFS)
		})
	}
}
