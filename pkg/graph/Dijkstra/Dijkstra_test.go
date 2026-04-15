package dijkstra

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_dijkstra(t *testing.T) {
	tests := []struct {
		name     string
		graph    [][][2]int
		src      int
		expected []int
	}{
		{
			name: "single_node",
			graph: [][][2]int{
				{}, // node 0: no edges
			},
			src:      0,
			expected: []int{0},
		},
		{
			name: "two_nodes_direct",
			graph: [][][2]int{
				{{1, 5}}, // 0→1 cost 5
				{},       // 1: no outgoing
			},
			src:      0,
			expected: []int{0, 5},
		},
		{
			name: "diamond_graph",
			// 0→1 cost 1, 0→2 cost 4, 1→3 cost 2, 2→3 cost 1
			// shortest: 0→1→3 = 3
			graph: [][][2]int{
				{{1, 1}, {2, 4}}, // 0
				{{3, 2}},         // 1
				{{3, 1}},         // 2
				{},               // 3
			},
			src:      0,
			expected: []int{0, 1, 4, 3},
		},
		{
			name: "stale_entry_pruned",
			// 0→1 via two paths: cost 1 and cost 10
			// stale entry (cost 10) should be skipped
			graph: [][][2]int{
				{{1, 1}, {1, 10}}, // 0→1 twice
				{},
			},
			src:      0,
			expected: []int{0, 1},
		},
		{
			name: "unreachable_node",
			graph: [][][2]int{
				{{1, 3}}, // 0→1
				{},       // 1
				{},       // 2: unreachable
			},
			src:      0,
			expected: []int{0, 3, math.MaxInt},
		},
		{
			name: "non_zero_src",
			// src=2, shortest from node 2: 2→0=10, 2→1=4, 2→2=0, 2→3=6
			graph: [][][2]int{
				{{1, 7}},          // 0→1
				{{3, 2}},          // 1→3
				{{0, 10}, {1, 4}}, // 2→0, 2→1
				{},                // 3
			},
			src:      2,
			expected: []int{10, 4, 0, 6},
		},
		{
			name: "cycle_no_infinite_loop",
			// 0→1→2→0 forms a cycle; Dijkstra must not loop
			// shortest: 0=0, 1=1, 2=3
			graph: [][][2]int{
				{{1, 1}}, // 0→1
				{{2, 2}}, // 1→2
				{{0, 1}}, // 2→0 (cycle back)
			},
			src:      0,
			expected: []int{0, 1, 3},
		},
		{
			name: "equal_cost_paths",
			// Two paths 0→2: via node 1 (1+2=3) and direct (3).
			// Both cost 3; either is valid; result must be 3.
			graph: [][][2]int{
				{{1, 1}, {2, 3}}, // 0→1, 0→2
				{{2, 2}},         // 1→2
				{},               // 2
			},
			src:      0,
			expected: []int{0, 1, 3},
		},
		{
			name: "linear_chain",
			// 0→1→2→3→4, each weight 5
			graph: [][][2]int{
				{{1, 5}}, // 0
				{{2, 5}}, // 1
				{{3, 5}}, // 2
				{{4, 5}}, // 3
				{},       // 4
			},
			src:      0,
			expected: []int{0, 5, 10, 15, 20},
		},
		{
			name: "dense_graph_multiple_relaxations",
			// 0→1=10, 0→2=3, 2→1=4, 1→3=2, 2→3=8
			// shortest: 0=0, 1=7(via 2), 2=3, 3=9
			graph: [][][2]int{
				{{1, 10}, {2, 3}}, // 0
				{{3, 2}},          // 1
				{{1, 4}, {3, 8}},  // 2
				{},                // 3
			},
			src:      0,
			expected: []int{0, 7, 3, 9},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, dijkstra(tt.graph, tt.src))
		})
	}
}

func Test_findCheapestPrice(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		flights  [][]int
		src, dst int
		k        int
		expected int
	}{
		{
			name:    "leetcode_example1",
			n:       4,
			flights: [][]int{{0, 1, 100}, {1, 2, 100}, {2, 0, 100}, {1, 3, 600}, {2, 3, 200}},
			src:     0, dst: 3, k: 1,
			expected: 700,
		},
		{
			name:    "leetcode_example2",
			n:       3,
			flights: [][]int{{0, 1, 100}, {1, 2, 100}, {0, 2, 500}},
			src:     0, dst: 2, k: 1,
			expected: 200,
		},
		{
			name:    "leetcode_example3_direct_only",
			n:       3,
			flights: [][]int{{0, 1, 100}, {1, 2, 100}, {0, 2, 500}},
			src:     0, dst: 2, k: 0,
			expected: 500,
		},
		{
			name:    "no_path",
			n:       3,
			flights: [][]int{{0, 1, 100}},
			src:     0, dst: 2, k: 1,
			expected: -1,
		},
		{
			name:    "direct_flight",
			n:       2,
			flights: [][]int{{0, 1, 50}},
			src:     0, dst: 1, k: 0,
			expected: 50,
		},
		{
			name:    "k_limits_cheaper_route",
			n:       4,
			flights: [][]int{{0, 1, 1}, {1, 2, 1}, {2, 3, 1}, {0, 3, 100}},
			src:     0, dst: 3, k: 1,
			// 0→1→2→3 needs 2 stops (k=1 allows 1 stop), only 0→3=100 works
			expected: 100,
		},
		{
			name: "pruning_same_node_same_stops",
			// Two parallel 0→1 flights (costs 1 and 2). When the second pops,
			// visited[1] <= cur.stops triggers the pruning branch.
			n:       3,
			flights: [][]int{{0, 1, 1}, {0, 1, 2}, {1, 2, 1}},
			src:     0, dst: 2, k: 2,
			expected: 2, // 0→1 (cost 1) → 2 (cost 1) = 2
		},
		{
			name:    "src_equals_dst",
			n:       3,
			flights: [][]int{{0, 1, 100}, {1, 2, 100}},
			src:     1, dst: 1, k: 1,
			expected: 0, // already at destination
		},
		{
			name: "k_larger_than_needed",
			// Cheapest is 0→1→2 (2 stops = k=1), k=5 shouldn't change the answer.
			n:       3,
			flights: [][]int{{0, 1, 100}, {1, 2, 100}, {0, 2, 500}},
			src:     0, dst: 2, k: 5,
			expected: 200,
		},
		{
			name: "cycle_in_flight_graph",
			// 0→1→2→1 forms a cycle; must terminate via stop-count pruning.
			n:       3,
			flights: [][]int{{0, 1, 10}, {1, 2, 20}, {2, 1, 5}},
			src:     0, dst: 2, k: 3,
			expected: 30, // 0→1→2
		},
		{
			name: "cheaper_path_needs_one_more_stop",
			// Cheap path 0→1→2→3 costs 10 (2 stops); direct 0→3 costs 100 (0 stops).
			// With k=1, cheap path needs 2 stops so only direct works.
			n:       4,
			flights: [][]int{{0, 1, 1}, {1, 2, 4}, {2, 3, 5}, {0, 3, 100}},
			src:     0, dst: 3, k: 1,
			expected: 100,
		},
		{
			name: "k_exactly_enough_for_cheapest",
			// Cheap path 0→1→2→3 costs 10 (2 stops); k=2 allows exactly that.
			n:       4,
			flights: [][]int{{0, 1, 1}, {1, 2, 4}, {2, 3, 5}, {0, 3, 100}},
			src:     0, dst: 3, k: 2,
			expected: 10,
		},
		{
			name:    "k_zero_no_direct_flight",
			n:       3,
			flights: [][]int{{0, 1, 50}, {1, 2, 50}},
			src:     0, dst: 2, k: 0,
			expected: -1, // direct flight doesn't exist, can't use any stop
		},
		{
			name: "multiple_routes_pick_cheapest_within_k",
			// 0→3: via A (0→1→3, cost=15, 1 stop) or via B (0→2→3, cost=8, 1 stop).
			// Both within k=1; should pick cheaper route via B.
			n:       4,
			flights: [][]int{{0, 1, 10}, {1, 3, 5}, {0, 2, 3}, {2, 3, 5}},
			src:     0, dst: 3, k: 1,
			expected: 8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, findCheapestPrice(tt.n, tt.flights, tt.src, tt.dst, tt.k))
		})
	}
}

func Test_networkDelayTime(t *testing.T) {
	tests := []struct {
		name     string
		times    [][]int
		n, k     int
		expected int
	}{
		{
			name:     "leetcode_example1",
			times:    [][]int{{2, 1, 1}, {2, 3, 1}, {3, 4, 1}},
			n:        4,
			k:        2,
			expected: 2, // signal reaches 1 in 1, 3 in 1, 4 in 2 → max = 2
		},
		{
			name:     "leetcode_example2",
			times:    [][]int{{1, 2, 1}},
			n:        2,
			k:        1,
			expected: 1,
		},
		{
			name:     "single_node",
			times:    [][]int{},
			n:        1,
			k:        1,
			expected: 0, // only source exists, already reached
		},
		{
			name:     "unreachable_node",
			times:    [][]int{{1, 2, 1}},
			n:        3,
			k:        1,
			expected: -1, // node 3 is unreachable from 1
		},
		{
			name:     "source_has_no_outgoing_edges",
			times:    [][]int{{2, 3, 5}},
			n:        3,
			k:        1,
			expected: -1, // k=1 has no outgoing edges; nodes 2 and 3 unreachable
		},
		{
			name: "star_topology",
			// k=1 broadcasts directly to 2,3,4 with different delays
			times:    [][]int{{1, 2, 5}, {1, 3, 3}, {1, 4, 8}},
			n:        4,
			k:        1,
			expected: 8, // bottleneck is node 4 at delay 8
		},
		{
			name: "longer_path_beats_direct",
			// 1→3 direct costs 10; 1→2→3 costs 1+1=2.
			// All nodes reachable; bottleneck = max(dist[2]=1, dist[3]=2) = 2.
			times:    [][]int{{1, 2, 1}, {2, 3, 1}, {1, 3, 10}},
			n:        3,
			k:        1,
			expected: 2,
		},
		{
			name: "linear_chain",
			// 1→2→3→4→5, each weight 1; delay to reach all = 4
			times:    [][]int{{1, 2, 1}, {2, 3, 1}, {3, 4, 1}, {4, 5, 1}},
			n:        5,
			k:        1,
			expected: 4,
		},
		{
			name: "k_is_not_node_1",
			// Source is node 3; nodes 1 and 2 are unreachable from 3.
			times:    [][]int{{1, 2, 1}, {2, 3, 1}},
			n:        3,
			k:        3,
			expected: -1,
		},
		{
			name: "cycle_in_graph",
			// 1→2→3→1 cycle; signal reaches 2 in 1, 3 in 2 → max = 2, no infinite loop
			times:    [][]int{{1, 2, 1}, {2, 3, 1}, {3, 1, 1}},
			n:        3,
			k:        1,
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, networkDelayTime(tt.times, tt.n, tt.k))
		})
	}
}

// ---------------------------------------------------------------------------
// 1631. Path With Minimum Effort  (TDD — implementation has a bug)
//
// Effort of a path = max absolute difference between adjacent cells' heights.
// Find the minimum effort path from (0,0) to (m-1,n-1).
//
// Cost function must be: max(cur.effort, abs(heights[nr][nc] - heights[r][c]))
// NOT a sum — the current implementation sums differences and skips abs, so
// all tests except the trivial ones are expected to fail until it is fixed.
// ---------------------------------------------------------------------------

func Test_minimumEffortPath(t *testing.T) {
	tests := []struct {
		name     string
		heights  [][]int
		expected int
	}{
		{
			// Single cell — already at destination.
			name:     "single_cell",
			heights:  [][]int{{5}},
			expected: 0,
		},
		{
			// 1×2: only one path, effort = |1-5| = 4.
			name:     "one_row_two_cols",
			heights:  [][]int{{1, 5}},
			expected: 4,
		},
		{
			// 2×1: only one path, effort = |1-5| = 4.
			name:     "two_rows_one_col",
			heights:  [][]int{{1}, {5}},
			expected: 4,
		},
		{
			// All cells same height — zero effort regardless of path.
			name:     "all_same_height",
			heights:  [][]int{{3, 3, 3}, {3, 3, 3}, {3, 3, 3}},
			expected: 0,
		},
		{
			// LeetCode example 1.
			// Optimal: (0,0)→(1,0)→(2,0)→(2,1)→(2,2)
			// diffs: |1-3|=2, |3-5|=2, |5-3|=2, |3-5|=2 → max = 2.
			name:     "leetcode_example1",
			heights:  [][]int{{1, 2, 2}, {3, 8, 2}, {5, 3, 5}},
			expected: 2,
		},
		{
			// LeetCode example 2.
			// Optimal: (0,0)→(0,1)→(0,2)→(1,2)→(2,2)
			// diffs: 1,1,1,1 → max = 1.
			name:     "leetcode_example2",
			heights:  [][]int{{1, 2, 3}, {3, 8, 4}, {5, 3, 5}},
			expected: 1,
		},
		{
			// LeetCode example 3.
			// A flat corridor exists along the bottom — effort = 0.
			name:    "leetcode_example3",
			heights: [][]int{{1, 2, 1, 1, 1}, {1, 2, 1, 2, 1}, {1, 2, 1, 2, 1}, {1, 2, 1, 2, 1}, {1, 1, 1, 2, 1}},
			expected: 0,
		},
		{
			// Two candidate paths; going down-then-right beats right-then-down.
			// heights = [[1,10],[2,3]]
			// Path A: (0,0)→(0,1)→(1,1): max(|1-10|,|10-3|) = max(9,7) = 9
			// Path B: (0,0)→(1,0)→(1,1): max(|1-2|,|2-3|)  = max(1,1)  = 1
			name:     "two_paths_pick_lower_effort",
			heights:  [][]int{{1, 10}, {2, 3}},
			expected: 1,
		},
		{
			// Effort is the max of the path, not the sum.
			// heights = [[1,2,10],[1,1,1]]
			// Path A (top row): max(1,8) = 8
			// Path B (down, across, up): (0,0)→(1,0)→(1,1)→(1,2)→(0,2)?
			//   Not allowed to go up to (0,2) from (1,2) then... actually can.
			//   diffs: 0,0,0,9 → max=9. Worse.
			// Path C: (0,0)→(1,0)→(1,1)→(1,2): diffs 0,0,1 → max=1. But doesn't reach (0,2).
			//   Destination is (1,2) for a 2-row grid. So path C is optimal with effort 1.
			name:     "effort_is_max_not_sum",
			heights:  [][]int{{1, 2, 10}, {1, 1, 1}},
			expected: 1,
		},
		{
			// 1×n strip — only one path, effort = max diff between any two adjacent cells.
			// heights = [1, 3, 1, 4]
			// diffs: 2, 2, 3 → max = 3.
			name:     "single_row_strip",
			heights:  [][]int{{1, 3, 1, 4}},
			expected: 3,
		},
		{
			// Uniform increase along the direct path, but a flat detour exists.
			// heights = [[1,1,1],[5,5,1],[5,5,1]]
			// Top row then right column: diffs 0,0,0,0,0 → effort 0.
			name:     "flat_detour_exists",
			heights:  [][]int{{1, 1, 1}, {5, 5, 1}, {5, 5, 1}},
			expected: 0,
		},
		{
			// Large cliff in every direction — can't avoid it, effort = cliff height.
			// heights = [[1,100],[100,1]]
			// Path A: (0,0)→(0,1)→(1,1): max(99, 99) = 99
			// Path B: (0,0)→(1,0)→(1,1): max(99, 99) = 99
			name:     "unavoidable_cliff",
			heights:  [][]int{{1, 100}, {100, 1}},
			expected: 99,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, minimumEffortPath(tt.heights))
		})
	}
}
