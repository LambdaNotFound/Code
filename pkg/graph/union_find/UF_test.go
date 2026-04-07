package unionfind

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_longestConsecutive(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected int
	}{
		{name: "leetcode_example1", nums: []int{100, 4, 200, 1, 3, 2}, expected: 4},
		{name: "leetcode_example2", nums: []int{0, 3, 7, 2, 5, 8, 4, 6, 0, 1}, expected: 9},
		{name: "empty", nums: []int{}, expected: 0},
		{name: "single_element", nums: []int{5}, expected: 1},
		{name: "all_same", nums: []int{2, 2, 2}, expected: 1},
		{name: "no_consecutive", nums: []int{1, 3, 5, 7}, expected: 1},
		{name: "two_separate_sequences", nums: []int{1, 2, 10, 11, 12}, expected: 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, longestConsecutive(tt.nums))
		})
	}
}

func Test_UnionFind_map(t *testing.T) {
	t.Run("new_union_find_is_empty", func(t *testing.T) {
		uf := NewUnionFind()
		assert.Equal(t, 0, uf.MaxSize())
	})

	t.Run("add_single_element", func(t *testing.T) {
		uf := NewUnionFind()
		uf.Add(42)
		assert.Equal(t, 42, uf.Find(42))
		assert.Equal(t, 1, uf.MaxSize())
	})

	t.Run("add_is_idempotent", func(t *testing.T) {
		uf := NewUnionFind()
		uf.Add(1)
		uf.Add(1) // second Add should be a no-op
		assert.Equal(t, 1, uf.MaxSize())
	})

	t.Run("union_two_elements_merges_sets", func(t *testing.T) {
		uf := NewUnionFind()
		uf.Add(1)
		uf.Add(2)
		assert.Equal(t, 1, uf.MaxSize())
		uf.Union(1, 2)
		assert.Equal(t, 2, uf.MaxSize())
	})

	t.Run("union_same_root_is_noop", func(t *testing.T) {
		uf := NewUnionFind()
		uf.Add(1)
		uf.Add(2)
		uf.Union(1, 2)
		uf.Union(1, 2) // redundant
		assert.Equal(t, 2, uf.MaxSize())
	})

	t.Run("path_compression_via_find", func(t *testing.T) {
		uf := NewUnionFind()
		uf.Add(1)
		uf.Add(2)
		uf.Add(3)
		uf.Union(1, 2)
		uf.Union(2, 3)
		root := uf.Find(1)
		assert.Equal(t, root, uf.Find(2))
		assert.Equal(t, root, uf.Find(3))
		assert.Equal(t, 3, uf.MaxSize())
	})

	t.Run("union_by_rank_higher_rank_wins", func(t *testing.T) {
		uf := NewUnionFind()
		// Build a two-node component (raises rank of root)
		uf.Add(10)
		uf.Add(11)
		uf.Union(10, 11) // root of 10-11 gets rank 1
		// Add an isolated node and union it in — triggers rank[px] > rank[py] or <
		uf.Add(12)
		uf.Union(12, 10)
		assert.Equal(t, uf.Find(10), uf.Find(12))
		assert.Equal(t, 3, uf.MaxSize())
	})
}

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
