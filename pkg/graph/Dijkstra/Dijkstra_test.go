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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, findCheapestPrice(tt.n, tt.flights, tt.src, tt.dst, tt.k))
		})
	}
}
