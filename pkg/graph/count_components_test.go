package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_countComponents(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		edges    [][]int
		expected int
	}{
		{name: "leetcode_example1", n: 5, edges: [][]int{{0, 1}, {1, 2}, {3, 4}}, expected: 2},
		{name: "leetcode_example2", n: 5, edges: [][]int{{0, 1}, {1, 2}, {2, 3}, {3, 4}}, expected: 1},
		{name: "no_edges", n: 4, edges: [][]int{}, expected: 4},
		{name: "fully_connected", n: 3, edges: [][]int{{0, 1}, {1, 2}, {0, 2}}, expected: 1},
		{name: "single_node", n: 1, edges: [][]int{}, expected: 1},
		{name: "two_nodes_connected", n: 2, edges: [][]int{{0, 1}}, expected: 1},
		{name: "two_nodes_disconnected", n: 2, edges: [][]int{}, expected: 2},
		{name: "star_plus_isolated", n: 5, edges: [][]int{{0, 1}, {0, 2}, {0, 3}}, expected: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, countComponents(tt.n, tt.edges))
			assert.Equal(t, tt.expected, countComponentsBFS(tt.n, tt.edges))
			assert.Equal(t, tt.expected, countComponentsDFS(tt.n, tt.edges))
		})
	}
}
