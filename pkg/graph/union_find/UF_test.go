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

// ---------------------------------------------------------------------------
// Template functions (unionFind / unionFindByRank / unionFindBySize)
//
// These three functions are pedagogical templates — they don't return a value.
// The tests below call each one with graphs that exercise every internal branch:
//   • a simple path (exercises the basic union + path-compression find)
//   • a triangle (exercises the "already same root" early-return in union)
//   • a chain that forces rank promotions (unionFindByRank equal-rank branch)
//   • a skewed tree (unionFindBySize both size-comparison branches)
// ---------------------------------------------------------------------------
func Test_unionFind_templates(t *testing.T) {
	t.Run("basic_path_graph", func(t *testing.T) {
		// 0-1-2-3: each union merges two new components
		unionFind(4, [][2]int{{0, 1}, {1, 2}, {2, 3}})
	})

	t.Run("triangle_has_redundant_edge", func(t *testing.T) {
		// 0-1, 1-2, 0-2: the third edge hits the "already same root" branch
		unionFind(3, [][2]int{{0, 1}, {1, 2}, {0, 2}})
	})

	t.Run("isolated_nodes_no_edges", func(t *testing.T) {
		unionFind(5, [][2]int{})
	})
}

func Test_unionFindByRank_templates(t *testing.T) {
	t.Run("chain_promotes_rank", func(t *testing.T) {
		// 0-1 → rank[root]=1; 2-3 → rank[root]=1; then union of both subtrees
		// triggers the rank[rootX]==rank[rootY] branch and increments rank.
		unionFindByRank(4, [][2]int{{0, 1}, {2, 3}, {0, 2}})
	})

	t.Run("deep_chain_path_compression", func(t *testing.T) {
		// Linear chain: 0-1-2-3-4 exercises path compression through several levels.
		unionFindByRank(5, [][2]int{{0, 1}, {1, 2}, {2, 3}, {3, 4}})
	})

	t.Run("redundant_edge_same_root", func(t *testing.T) {
		// The last edge 0-3 finds both nodes already in the same component.
		unionFindByRank(4, [][2]int{{0, 1}, {1, 2}, {2, 3}, {0, 3}})
	})

	t.Run("lower_rank_rootX_attaches_to_higher_rank_rootY", func(t *testing.T) {
		// After union(0,1): root=0 gets rank 1.
		// union(2,0): find(2)=2 (rank 0) < find(0)=0 (rank 1)
		// → hits the rank[rootX] < rank[rootY] branch → parent[2] = 0.
		unionFindByRank(3, [][2]int{{0, 1}, {2, 0}})
	})
}

func Test_unionFindBySize_templates(t *testing.T) {
	t.Run("size_guided_merge_large_into_small", func(t *testing.T) {
		// Build a 3-node component {0,1,2} then attach isolated node 3.
		// size[root of {0,1,2}]=3 > size[3]=1: exercises the else-branch
		// (parent[rootY]=rootX, size[rootX]+=size[rootY]).
		unionFindBySize(4, [][2]int{{0, 1}, {0, 2}, {0, 3}})
	})

	t.Run("size_guided_merge_equal_sizes", func(t *testing.T) {
		// {0,1} and {2,3} both have size 2 — equal sizes exercise the else-branch.
		unionFindBySize(4, [][2]int{{0, 1}, {2, 3}, {0, 2}})
	})

	t.Run("redundant_union_same_component", func(t *testing.T) {
		// After 0-1-2 are merged, unioning 0-2 again hits "rootX==rootY" return false.
		unionFindBySize(3, [][2]int{{0, 1}, {1, 2}, {0, 2}})
	})

	t.Run("small_rootX_merges_into_large_rootY", func(t *testing.T) {
		// Build {0,1,2} (size 3 at root 0), then union(3,1):
		// find(3)=3 (size 1), find(1)=0 (size 3) → size[rootX=3] < size[rootY=0]
		// → hits the if-branch: parent[3]=0, size[0]+=size[3].
		unionFindBySize(4, [][2]int{{0, 1}, {0, 2}, {3, 1}})
	})
}

// ---------------------------------------------------------------------------
// 839. Similar String Groups
// ---------------------------------------------------------------------------
func Test_numSimilarGroups(t *testing.T) {
	tests := []struct {
		name     string
		strs     []string
		expected int
	}{
		{
			// LeetCode example 1.
			// tars↔rats (pos 0,2 swap), rats↔arts (pos 0,1 swap) → {tars,rats,arts}, {star}
			name: "leetcode_example1",
			strs: []string{"tars", "rats", "arts", "star"},
			expected: 2,
		},
		{
			// LeetCode example 2.
			// omv↔ovm: positions 1,2 differ (2 diffs) → similar → 1 group.
			name: "leetcode_example2",
			strs: []string{"omv", "ovm"},
			expected: 1,
		},
		{
			// Single string — trivially one group.
			name:     "single_string",
			strs:     []string{"abc"},
			expected: 1,
		},
		{
			// All identical: diffs==0 for every pair → all similar → 1 group.
			// Exercises the diffs==0 branch of isSimilar.
			name:     "all_identical",
			strs:     []string{"abc", "abc", "abc"},
			expected: 1,
		},
		{
			// abcd vs dcba: positions 0,1,2,3 all differ — isSimilar exits early
			// at diffs>2, returns false → 2 separate groups.
			// Exercises the early-exit (diffs > 2) branch.
			name:     "four_diffs_triggers_early_exit",
			strs:     []string{"abcd", "dcba"},
			expected: 2,
		},
		{
			// Transitive chain: abc↔bac (2 diffs), bac↔bca (2 diffs), abc↔bca (3 diffs).
			// abc and bca are not directly similar but are connected via bac → 1 group.
			// Exercises BFS queue re-entry and the visited[i] skip branch.
			name:     "transitive_chain_one_group",
			strs:     []string{"abc", "bac", "bca"},
			expected: 1,
		},
		{
			// Two completely disjoint similar pairs — no cross-pair similarity.
			name:     "two_disjoint_groups",
			strs:     []string{"abc", "bac", "xyz", "yxz"},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, numSimilarGroups(tt.strs))
			assert.Equal(t, tt.expected, numSimilarGroupsUF(tt.strs))
		})
	}
}

// ---------------------------------------------------------------------------
// 1061. Lexicographically Smallest Equivalent String
// ---------------------------------------------------------------------------
func Test_smallestEquivalentString(t *testing.T) {
	tests := []struct {
		name     string
		s1       string
		s2       string
		baseStr  string
		expected string
	}{
		{
			// LeetCode example 1.
			// Equivalence classes after pairing parker/morris:
			//   {m,p}, {a,o}, {k,r,s}, {e,i}
			// baseStr "parser" → p→m, a→a, r→k, s→k, e→e, r→k → "makkek"
			name: "leetcode_example1",
			s1: "parker", s2: "morris", baseStr: "parser",
			expected: "makkek",
		},
		{
			// LeetCode example 2.
			// Classes after pairing hello/world: {h,w}, {d,e,o}, {l,r}
			// "hold" → h→h, o→d, l→l, d→d → "hdld"
			name: "leetcode_example2",
			s1: "hello", s2: "world", baseStr: "hold",
			expected: "hdld",
		},
		{
			// Single pair: a↔b. Any 'b' in baseStr maps to 'a'.
			name: "single_pair_maps_to_smaller",
			s1: "a", s2: "b", baseStr: "b",
			expected: "a",
		},
		{
			// Pair maps a character to itself — no-op union, identity result.
			name: "identity_pair_no_change",
			s1: "a", s2: "a", baseStr: "a",
			expected: "a",
		},
		{
			// Transitive chain: a=b (pair 1), b=c (pair 2) → {a,b,c}, all map to 'a'.
			name: "transitive_chain_all_map_to_smallest",
			s1: "ab", s2: "bc", baseStr: "abc",
			expected: "aaa",
		},
		{
			// baseStr contains a character not involved in any pair — returned unchanged.
			name: "baseStr_char_not_in_any_pair",
			s1: "a", s2: "b", baseStr: "c",
			expected: "c",
		},
		{
			// All three input chars are pairwise equivalent:
			// a=b (pair 1), b=c (pair 2), c=a (pair 3, redundant). "xyz" unchanged.
			name: "three_way_equivalence_baseStr_unrelated",
			s1: "abc", s2: "bca", baseStr: "xyz",
			expected: "xyz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, smallestEquivalentString(tt.s1, tt.s2, tt.baseStr))
		})
	}
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
		{
			// Build {0,1,2} (size 3), then union(3,1): find(3)=3 (size 1) < find(1)=0 (size 3)
			// → hits the size[rootX] < size[rootY] if-branch in countComponents.
			name: "small_component_merges_into_large", n: 4, edges: [][]int{{0, 1}, {1, 2}, {3, 1}}, expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, countComponents(tt.n, tt.edges))
			assert.Equal(t, tt.expected, countComponentsBFS(tt.n, tt.edges))
			assert.Equal(t, tt.expected, countComponentsDFS(tt.n, tt.edges))
		})
	}
}
