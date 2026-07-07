package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_validTree(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		edges    [][]int
		expected bool
	}{
		{name: "leetcode_example1", n: 5, edges: [][]int{{0, 1}, {0, 2}, {0, 3}, {1, 4}}, expected: true},
		{name: "leetcode_example2", n: 5, edges: [][]int{{0, 1}, {1, 2}, {2, 3}, {1, 3}, {1, 4}}, expected: false},
		{name: "single_node", n: 1, edges: [][]int{}, expected: true},
		{name: "two_nodes_connected", n: 2, edges: [][]int{{0, 1}}, expected: true},
		{name: "two_nodes_disconnected", n: 2, edges: [][]int{}, expected: false},
		{name: "cycle_triangle", n: 3, edges: [][]int{{0, 1}, {1, 2}, {2, 0}}, expected: false},
		{name: "disconnected_forest", n: 4, edges: [][]int{{0, 1}, {2, 3}}, expected: false},
		{name: "chain", n: 4, edges: [][]int{{0, 1}, {1, 2}, {2, 3}}, expected: true},
		{name: "star_graph", n: 5, edges: [][]int{{0, 1}, {0, 2}, {0, 3}, {0, 4}}, expected: true},
		// edge ordering that triggers rank[rootA] < rank[rootB] branch in validTreeUF:
		// union(1,0) → rank[1]=1; then union(2,0): rootA=2 (rank 0) < rootB=1 (rank 1)
		{name: "rank_lt_branch", n: 3, edges: [][]int{{1, 0}, {2, 0}}, expected: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, validTree(tt.n, tt.edges))
			assert.Equal(t, tt.expected, validTreeBFS(tt.n, tt.edges))
			assert.Equal(t, tt.expected, validTreeDFS(tt.n, tt.edges))
			assert.Equal(t, tt.expected, validTreeUF(tt.n, tt.edges))
		})
	}
}
