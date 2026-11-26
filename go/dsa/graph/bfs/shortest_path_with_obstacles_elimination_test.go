package bfs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/**
 * 1293. Shortest Path in a Grid with Obstacles Elimination
 *
 * BFS state: (row, col, eliminationsRemaining)
 * Returns the minimum number of steps to reach (m-1, n-1) from (0, 0),
 * or -1 if no such path exists.
 *
 * Test strategy:
 *   - LeetCode canonical examples.
 *   - Degenerate inputs: single cell, obstacle-free grid, 2×2 grids.
 *   - k=0 paired with k=1 to isolate the impact of a single elimination.
 *   - Single-row and single-column grids (1-D BFS).
 *   - The k ≥ m+n-2 early-exit optimisation: when k is large enough to
 *     eliminate every cell on any Manhattan path, the answer is always m+n-2
 *     regardless of obstacle layout.
 *   - A case where k exactly matches the number of obstacles on the optimal
 *     path (tight budget).
 *   - Diagonal bottleneck: two obstacles block the top-left corner; BFS must
 *     route around even with k≥1 because both obstacles sit on all short paths.
 *   - Zigzag maze: the only k=0 path doubles back twice (16 steps on a 5×5
 *     grid); a single elimination breaks through one wall and restores the
 *     Manhattan-optimal 8-step path.
 */
func Test_shortestPath(t *testing.T) {
	tests := []struct {
		name     string
		grid     [][]int
		k        int
		expected int
	}{
		// -----------------------------------------------------------------------
		// LeetCode canonical examples
		// -----------------------------------------------------------------------
		{
			// Optimal path eliminates one obstacle:
			// (0,0)→(0,1)→(0,2)→(1,2)→(2,2)→(3,2)[elim]→(4,2)  = 6 steps
			name: "leetcode_example1",
			grid: [][]int{
				{0, 0, 0},
				{1, 1, 0},
				{0, 0, 0},
				{0, 1, 1},
				{0, 0, 0},
			},
			k:        1,
			expected: 6,
		},
		{
			// Every path from (0,0) to (2,2) requires at least 2 eliminations;
			// with k=1 no path exists.
			name: "leetcode_example2",
			grid: [][]int{
				{0, 1, 1},
				{1, 1, 1},
				{1, 0, 0},
			},
			k:        1,
			expected: -1,
		},

		// -----------------------------------------------------------------------
		// Degenerate / edge inputs
		// -----------------------------------------------------------------------
		{
			// Start == destination; 0 steps required.
			name:     "single_cell",
			grid:     [][]int{{0}},
			k:        0,
			expected: 0,
		},
		{
			// Obstacle-free 3×3 grid with k=0: Manhattan distance = m+n-2 = 4.
			name: "no_obstacles_k0",
			grid: [][]int{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
			k:        0,
			expected: 4,
		},
		{
			// 2×2 open grid, k=0: two steps along any L-shaped path.
			name:     "open_2x2_k0",
			grid:     [][]int{{0, 0}, {0, 0}},
			k:        0,
			expected: 2,
		},
		{
			// Anti-diagonal obstacles block both exits from (0,0); k=0 → -1.
			name:     "anti_diagonal_blocked_k0",
			grid:     [][]int{{0, 1}, {1, 0}},
			k:        0,
			expected: -1,
		},

		// -----------------------------------------------------------------------
		// k ≥ m+n-2 early-exit optimisation
		// -----------------------------------------------------------------------
		{
			// 3×3 grid fully packed with interior obstacles.
			// k=4 ≥ m+n-2=4 → answer is immediately 4 without BFS.
			name: "k_geq_diagonal_optimization",
			grid: [][]int{
				{0, 1, 1},
				{1, 1, 1},
				{1, 1, 0},
			},
			k:        4,
			expected: 4,
		},

		// -----------------------------------------------------------------------
		// Surrounded start — can't move without an elimination
		// -----------------------------------------------------------------------
		{
			// Both neighbors of (0,0) are obstacles; k=0 means instant dead-end.
			name:     "surrounded_start_k0_impossible",
			grid:     [][]int{{0, 1}, {1, 0}},
			k:        0,
			expected: -1,
		},
		{
			// Same grid with k=1: eliminate one obstacle and reach (1,1) in 2 steps.
			name:     "surrounded_start_k1_possible",
			grid:     [][]int{{0, 1}, {1, 0}},
			k:        1,
			expected: 2,
		},

		// -----------------------------------------------------------------------
		// Single-row grids (1-D BFS)
		// -----------------------------------------------------------------------
		{
			// One obstacle in the middle of a single row; no way around with k=0.
			name:     "single_row_obstacle_k0",
			grid:     [][]int{{0, 1, 0, 0, 0}},
			k:        0,
			expected: -1,
		},
		{
			// Same row with k=1: eliminate the obstacle and walk straight through.
			name:     "single_row_obstacle_k1",
			grid:     [][]int{{0, 1, 0, 0, 0}},
			k:        1,
			expected: 4,
		},

		// -----------------------------------------------------------------------
		// Single-column grids (1-D BFS, vertical)
		// -----------------------------------------------------------------------
		{
			// Obstacle in the middle of a single column; k=0 → -1.
			name:     "single_col_obstacle_k0",
			grid:     [][]int{{0}, {1}, {0}, {0}, {0}},
			k:        0,
			expected: -1,
		},
		{
			// Same column with k=1: 4 steps straight down.
			name:     "single_col_obstacle_k1",
			grid:     [][]int{{0}, {1}, {0}, {0}, {0}},
			k:        1,
			expected: 4,
		},

		// -----------------------------------------------------------------------
		// Tight elimination budget
		// -----------------------------------------------------------------------
		{
			// Top-right and mid-left obstacles isolate (0,0); any shortest path
			// must eliminate at least one.  k=0 → impossible.
			name: "need_1_elim_k0_impossible",
			grid: [][]int{
				{0, 1, 0},
				{1, 1, 0},
				{0, 0, 0},
			},
			k:        0,
			expected: -1,
		},
		{
			// Same grid with k=1: eliminate (0,1) or (1,0) and reach (2,2) in 4 steps.
			name: "need_1_elim_k1_exact_budget",
			grid: [][]int{
				{0, 1, 0},
				{1, 1, 0},
				{0, 0, 0},
			},
			k:        1,
			expected: 4,
		},

		// -----------------------------------------------------------------------
		// Diagonal bottleneck
		// -----------------------------------------------------------------------
		{
			// Two obstacles block the top-left corner of a 4×4 grid.
			// BFS must route around the right side; k=1 doesn't help shorten it
			// because both obstacles sit on every sub-6 path.
			name: "diagonal_block_k1_forced_detour",
			grid: [][]int{
				{0, 1, 0, 0},
				{1, 1, 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
			},
			k:        1,
			expected: 6,
		},
		{
			// k=2 eliminates both blocking obstacles but the Manhattan distance
			// is already 6 = m+n-2, so the answer is unchanged.
			name: "diagonal_block_k2_same_answer",
			grid: [][]int{
				{0, 1, 0, 0},
				{1, 1, 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
			},
			k:        2,
			expected: 6,
		},

		// -----------------------------------------------------------------------
		// Zigzag maze — elimination dramatically shortens the path
		// -----------------------------------------------------------------------
		{
			// 5×5 grid with two horizontal walls forcing a full double-back:
			//   Row 1: 1 1 1 1 0   (gap only at right end)
			//   Row 3: 0 1 1 1 1   (gap only at left end)
			// k=0: must snake row 0 → gap at (1,4) → row 2 → gap at (3,0) →
			//       row 4: total 16 steps.
			// k=1: eliminate (1,0), cutting straight down rows 0-4 in 8 steps
			//       (= Manhattan distance m+n-2).
			name: "zigzag_maze_k0_long_detour",
			grid: [][]int{
				{0, 0, 0, 0, 0},
				{1, 1, 1, 1, 0},
				{0, 0, 0, 0, 0},
				{0, 1, 1, 1, 1},
				{0, 0, 0, 0, 0},
			},
			k:        0,
			expected: 16,
		},
		{
			// Same maze with k=1: one elimination restores the optimal 8-step path.
			name: "zigzag_maze_k1_cuts_through",
			grid: [][]int{
				{0, 0, 0, 0, 0},
				{1, 1, 1, 1, 0},
				{0, 0, 0, 0, 0},
				{0, 1, 1, 1, 1},
				{0, 0, 0, 0, 0},
			},
			k:        1,
			expected: 8,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, shortestPath(tc.grid, tc.k))
		})
	}
}
