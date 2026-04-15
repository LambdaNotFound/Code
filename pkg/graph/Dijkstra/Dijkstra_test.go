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
			name:     "leetcode_example3",
			heights:  [][]int{{1, 2, 1, 1, 1}, {1, 2, 1, 2, 1}, {1, 2, 1, 2, 1}, {1, 2, 1, 2, 1}, {1, 1, 1, 2, 1}},
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
			name:     "effort_is_max_not_sum",
			heights:  [][]int{{1, 2, 10}, {1, 1, 1}},
			expected: 0, // bottom row is flat: 1→1→1, all diffs = 0
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

// ---------------------------------------------------------------------------
// 1514. Path with Maximum Probability
//
// Undirected weighted graph where each edge has a success probability.
// Find the path from start to end that maximises the product of probabilities.
// If no path exists return 0.  Uses a max-heap (Dijkstra on probability).
// ---------------------------------------------------------------------------

func Test_maxProbability(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		edges    [][]int
		succProb []float64
		start    int
		end      int
		expected float64
	}{
		{
			// LeetCode example 1.
			// 0→1→2: 0.5×0.5 = 0.25, direct 0→2: 0.2 → pick 0.25.
			name:     "leetcode_example1",
			n:        3,
			edges:    [][]int{{0, 1}, {1, 2}, {0, 2}},
			succProb: []float64{0.5, 0.5, 0.2},
			start:    0, end: 2,
			expected: 0.25,
		},
		{
			// LeetCode example 2.
			// 0→1→2: 0.25, direct 0→2: 0.3 → direct wins.
			name:     "leetcode_example2",
			n:        3,
			edges:    [][]int{{0, 1}, {1, 2}, {0, 2}},
			succProb: []float64{0.5, 0.5, 0.3},
			start:    0, end: 2,
			expected: 0.3,
		},
		{
			// LeetCode example 3.
			// Node 2 is reachable from 0 only through 1, but no edge from 1 to 2.
			name:     "leetcode_example3_no_path",
			n:        3,
			edges:    [][]int{{0, 1}},
			succProb: []float64{0.5},
			start:    0, end: 2,
			expected: 0.0,
		},
		{
			// Single direct edge — answer equals the edge probability.
			name:     "direct_single_edge",
			n:        2,
			edges:    [][]int{{0, 1}},
			succProb: []float64{0.8},
			start:    0, end: 1,
			expected: 0.8,
		},
		{
			// start == end: already at destination, probability = 1.
			name:     "start_equals_end",
			n:        2,
			edges:    [][]int{{0, 1}},
			succProb: []float64{0.5},
			start:    1, end: 1,
			expected: 1.0,
		},
		{
			// Indirect path beats direct.
			// 0→1→2: 0.9×0.9 = 0.81 vs direct 0→2: 0.5 → 0.81 wins.
			name:     "indirect_beats_direct",
			n:        3,
			edges:    [][]int{{0, 1}, {1, 2}, {0, 2}},
			succProb: []float64{0.9, 0.9, 0.5},
			start:    0, end: 2,
			expected: 0.81,
		},
		{
			// Disconnected graph: {0,1} and {2,3} are separate components.
			name:     "disconnected_components",
			n:        4,
			edges:    [][]int{{0, 1}, {2, 3}},
			succProb: []float64{0.5, 0.5},
			start:    0, end: 3,
			expected: 0.0,
		},
		{
			// Single node, no edges: start == end, probability = 1.
			name:     "single_node",
			n:        1,
			edges:    [][]int{},
			succProb: []float64{},
			start:    0, end: 0,
			expected: 1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := maxProbability(tt.n, tt.edges, tt.succProb, tt.start, tt.end)
			assert.InDelta(t, tt.expected, got, 1e-5)
		})
	}
}

// ---------------------------------------------------------------------------
// 505. The Maze II
//
// A ball rolls in one of four directions until it hits a wall or the boundary.
// The distance is the number of cells traversed.  Find the minimum distance
// to reach destination (a cell where the ball stops).  Return -1 if impossible.
// ---------------------------------------------------------------------------

func Test_shortestDistance(t *testing.T) {
	tests := []struct {
		name        string
		maze        [][]int
		start       []int
		destination []int
		expected    int
	}{
		{
			// LeetCode example 1.
			// Ball navigates from (0,4) to (4,4); minimum distance is 12.
			name: "leetcode_example1",
			maze: [][]int{
				{0, 0, 1, 0, 0},
				{0, 0, 0, 0, 0},
				{0, 0, 0, 1, 0},
				{1, 1, 0, 1, 1},
				{0, 0, 0, 0, 0},
			},
			start: []int{0, 4}, destination: []int{4, 4},
			expected: 12,
		},
		{
			// LeetCode example 2.
			// (3,2) is not a valid stopping point — ball always rolls through it.
			name: "leetcode_example2",
			maze: [][]int{
				{0, 0, 1, 0, 0},
				{0, 0, 0, 0, 0},
				{0, 0, 0, 1, 0},
				{1, 1, 0, 1, 1},
				{0, 0, 0, 0, 0},
			},
			start: []int{0, 4}, destination: []int{3, 2},
			expected: -1,
		},
		{
			// Ball is already at the destination — no movement needed.
			name:        "start_equals_destination",
			maze:        [][]int{{0, 0, 0}},
			start:       []int{0, 1},
			destination: []int{0, 1},
			expected:    0,
		},
		{
			// 1×5 corridor: ball rolls right from col 0 to col 4 in one move.
			// Distance = 4 cells traversed.
			name:        "single_row_roll_to_boundary",
			maze:        [][]int{{0, 0, 0, 0, 0}},
			start:       []int{0, 0},
			destination: []int{0, 4},
			expected:    4,
		},
		{
			// Destination (1,2) lies in an interior cell with open cells on both
			// sides (row 0 and row 2).  The ball always rolls through it without
			// stopping, so the destination is unreachable.
			//
			//   [0, 0, 0]
			//   [0, 1, 0]   ← (1,2) is open but never a stop point
			//   [0, 0, 0]
			name: "destination_never_a_stop_point",
			maze: [][]int{
				{0, 0, 0},
				{0, 1, 0},
				{0, 0, 0},
			},
			start:       []int{0, 0},
			destination: []int{1, 2},
			expected:    -1,
		},
		{
			// Two symmetric paths to the destination; both have distance 4.
			// Verifies that the algorithm finds the optimal (= only) answer
			// without getting trapped in cycles.
			//
			//   [0, 0, 0]   (0,0) → right (2) → (0,2) → down (2) → (2,2)
			//   [0, 1, 0]   (0,0) → down  (2) → (2,0) → right (2) → (2,2)
			//   [0, 0, 0]
			name: "two_symmetric_paths",
			maze: [][]int{
				{0, 0, 0},
				{0, 1, 0},
				{0, 0, 0},
			},
			start:       []int{0, 0},
			destination: []int{2, 2},
			expected:    4,
		},
		{
			// Interior walls force the ball to use the perimeter.
			// Both perimeter routes (top-right then down, or down then bottom-right)
			// cost 6, which is the minimum.
			//
			//   [0, 0, 0, 0, 0]   row 0
			//   [0, 1, 1, 1, 0]   row 1 — walls at cols 1-3
			//   [0, 0, 0, 0, 0]   row 2
			//
			// (0,0) → right (4) → (0,4) → down (2) → (2,4): total 6.
			// (0,0) → down  (2) → (2,0) → right (4) → (2,4): total 6.
			name: "walls_force_perimeter_path",
			maze: [][]int{
				{0, 0, 0, 0, 0},
				{0, 1, 1, 1, 0},
				{0, 0, 0, 0, 0},
			},
			start:       []int{0, 0},
			destination: []int{2, 4},
			expected:    6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, shortestDistance(tt.maze, tt.start, tt.destination))
		})
	}
}

// ---------------------------------------------------------------------------
// 499. The Maze III
//
// Same rolling mechanic as Maze II, but the ball drops into a hole if it
// passes through or lands on it during a roll.  Return the lexicographically
// smallest path string (using 'd','l','r','u') that minimises total distance,
// or "impossible" if the hole can never be reached.
// ---------------------------------------------------------------------------

func Test_findShortestWay(t *testing.T) {
	tests := []struct {
		name     string
		maze     [][]int
		ball     []int
		hole     []int
		expected string
	}{
		{
			// LeetCode example 1.
			// ball=[4,3], hole=[0,1] → "lul"
			name: "leetcode_example1",
			maze: [][]int{
				{0, 0, 0, 0, 0},
				{1, 1, 0, 0, 1},
				{0, 0, 0, 0, 0},
				{0, 1, 0, 0, 1},
				{0, 1, 0, 0, 0},
			},
			ball: []int{4, 3}, hole: []int{0, 1},
			expected: "lul",
		},
		{
			// LeetCode example 2.
			// hole=[3,0]: every path from [4,3] that heads left is blocked by
			// walls at col 1 in row 4; the hole is unreachable.
			name: "leetcode_example2",
			maze: [][]int{
				{0, 0, 0, 0, 0},
				{1, 1, 0, 0, 1},
				{0, 0, 0, 0, 0},
				{0, 1, 0, 0, 1},
				{0, 1, 0, 0, 0},
			},
			ball: []int{4, 3}, hole: []int{3, 0},
			expected: "impossible",
		},
		{
			// Ball rolls right in a 1×5 corridor and falls into the hole
			// at (0,2) before reaching the far wall.
			// Only valid move: "r", distance 2.
			name:     "hole_caught_mid_roll",
			maze:     [][]int{{0, 0, 0, 0, 0}},
			ball:     []int{0, 0},
			hole:     []int{0, 2},
			expected: "r",
		},
		{
			// Two paths of equal distance (4); lexicographically smaller wins.
			//
			//   [0, 0, 0]
			//   [0, 1, 0]
			//   [0, 0, 0]
			//
			// ball=[2,0], hole=[0,2].
			//
			// Path "ru": roll right (2,0)→(2,2), then up (2,2)→passes (0,2)=hole. dist=4.
			// Path "ur": roll up   (2,0)→(0,0), then right (0,0)→passes (0,2)=hole. dist=4.
			//
			// 'r' < 'u', so "ru" < "ur" → answer = "ru".
			name:     "lexicographic_tiebreak",
			maze:     [][]int{{0, 0, 0}, {0, 1, 0}, {0, 0, 0}},
			ball:     []int{2, 0},
			hole:     []int{0, 2},
			expected: "ru",
		},
		{
			// Ball must roll right then down to reach the hole; only one reachable path.
			//
			//   [0, 0, 0]
			//   [1, 0, 0]   ← (1,0) is a wall; blocks downward roll from (0,0)
			//   [0, 0, 0]
			//
			// ball=[0,0], hole=[2,2].
			// (0,0) → right (dist 2) → (0,2) → down (dist 2) → passes (2,2)=hole.
			// Downward roll from (0,0) is blocked immediately by the wall at (1,0).
			// Only valid path: "rd", distance 4.
			name: "single_valid_path",
			maze: [][]int{
				{0, 0, 0},
				{1, 0, 0},
				{0, 0, 0},
			},
			ball:     []int{0, 0},
			hole:     []int{2, 2},
			expected: "rd",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, findShortestWay(tt.maze, tt.ball, tt.hole))
		})
	}
}
