package bfs

import (
	"gocode/dsa/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_orangesRotting(t *testing.T) {
	testCases := []struct {
		name     string
		grid     [][]int
		expected int
	}{
		{
			"case 1",
			[][]int{{2, 1, 1}, {1, 1, 0}, {0, 1, 1}},
			4,
		},
		{
			"case 2",
			[][]int{{2, 1, 1}, {0, 1, 1}, {1, 0, 1}},
			-1,
		},
		{
			name:     "single_fresh_adjacent",
			grid:     [][]int{{2, 1}},
			expected: 1,
		},
		{
			name:     "isolated_fresh",
			grid:     [][]int{{1, 0, 2}},
			expected: -1,
		},
		{
			name:     "single_fresh_cell",
			grid:     [][]int{{1}},
			expected: -1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			grid := utils.DeepCopyMatrix(tc.grid)
			result := orangesRotting(grid)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func Test_canMeasureWater(t *testing.T) {
	tests := []struct {
		name     string
		jug1     int
		jug2     int
		target   int
		expected bool
	}{
		{"leetcode_example1", 3, 5, 4, true},
		{"leetcode_example2", 2, 6, 5, false},
		{"leetcode_example3", 1, 2, 3, true},
		{"target_zero", 3, 5, 0, true},
		{"target_equals_jug1", 3, 5, 3, true},
		{"target_equals_jug2", 3, 5, 5, true},
		{"target_exceeds_capacity", 3, 5, 9, false},
		{"both_jugs_same", 4, 4, 4, true},
		{"target_is_sum", 3, 5, 8, true},
		{"impossible", 2, 4, 3, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, canMeasureWater(tt.jug1, tt.jug2, tt.target))
		})
	}
}
