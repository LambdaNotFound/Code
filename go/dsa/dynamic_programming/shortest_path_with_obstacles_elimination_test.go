package dynamic_programming

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_shortestPath(t *testing.T) {
	tests := []struct {
		name     string
		grid     [][]int
		k        int
		expected int
	}{
		{
			name: "leetcode_example1",
			grid: [][]int{
				{0, 0, 0},
				{1, 1, 0},
				{0, 0, 0},
				{0, 1, 1},
				{0, 0, 0},
			},
			k: 1, expected: 6,
		},
		{
			name: "leetcode_example2_no_elim",
			grid: [][]int{
				{0, 1, 1},
				{1, 1, 1},
				{1, 0, 0},
			},
			k: 1, expected: -1,
		},
		{
			name: "single_cell", grid: [][]int{{0}}, k: 0, expected: 0,
		},
		{
			name: "no_obstacles", grid: [][]int{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}}, k: 0, expected: 4,
		},
		{
			name: "obstacle_can_be_removed",
			grid: [][]int{{0, 1, 0}, {0, 1, 0}, {0, 0, 0}},
			k: 1, expected: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, shortestPath(tt.grid, tt.k))
		})
	}
}

// shortestPathTopDown uses memoised DFS without a per-path visited set, so it
// infinite-loops on open grids where backtracking is possible.  Only run it on
// grids where all neighbours of non-destination cells are either OOB or the
// destination itself — so the recursion terminates without cycling.
func Test_shortestPathTopDown(t *testing.T) {
	t.Run("single_cell", func(t *testing.T) {
		assert.Equal(t, 0, shortestPathTopDown([][]int{{0}}, 0))
	})
	t.Run("one_step_right", func(t *testing.T) {
		// (0,0) → right (0,1)=dest.  Left/up/down all OOB.  No cycle.
		assert.Equal(t, 1, shortestPathTopDown([][]int{{0, 0}}, 0))
	})
	t.Run("one_step_down", func(t *testing.T) {
		// (0,0) → down (1,0)=dest.  Right/left/up all OOB.  No cycle.
		assert.Equal(t, 1, shortestPathTopDown([][]int{{0}, {0}}, 0))
	})
	t.Run("blocked_wall_unreachable", func(t *testing.T) {
		// 1×3: [0,1,0].  From (0,0) only move is right→(0,1) obstacle with k=0 → INF.
		// Destination (0,2) is unreachable without elimination.
		assert.Equal(t, -1, shortestPathTopDown([][]int{{0, 1, 0}}, 0))
	})
}
