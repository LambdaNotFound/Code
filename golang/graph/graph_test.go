package graph

import (
	. "gocode/types"
	"gocode/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_BFS_DFS(t *testing.T) {
	t.Run("BFS_nil_start_no_panic", func(t *testing.T) {
		BFS(nil) // must not panic
	})

	t.Run("BFS_single_node", func(t *testing.T) {
		n := &Node{Val: 1}
		BFS(n) // must not panic, visits exactly n
	})

	t.Run("BFS_with_cycle_terminates", func(t *testing.T) {
		n1 := &Node{Val: 1}
		n2 := &Node{Val: 2}
		n1.Neighbors = []*Node{n2}
		n2.Neighbors = []*Node{n1}
		BFS(n1) // must not infinite-loop
	})

	t.Run("DFS_nil_start_no_panic", func(t *testing.T) {
		DFS(nil, map[*Node]bool{})
	})

	t.Run("DFS_single_node_marks_visited", func(t *testing.T) {
		n := &Node{Val: 1}
		visited := map[*Node]bool{}
		DFS(n, visited)
		assert.True(t, visited[n])
	})

	t.Run("DFS_already_visited_skips_subtree", func(t *testing.T) {
		n1 := &Node{Val: 1}
		n2 := &Node{Val: 2}
		n1.Neighbors = []*Node{n2}
		visited := map[*Node]bool{n1: true} // n1 pre-marked
		DFS(n1, visited)
		assert.False(t, visited[n2]) // n2 never reached
	})

	t.Run("DFS_with_cycle_terminates", func(t *testing.T) {
		n1 := &Node{Val: 1}
		n2 := &Node{Val: 2}
		n1.Neighbors = []*Node{n2}
		n2.Neighbors = []*Node{n1}
		visited := map[*Node]bool{}
		DFS(n1, visited)
		assert.True(t, visited[n1])
		assert.True(t, visited[n2])
	})
}

func Test_cloneGraph(t *testing.T) {
    tests := []struct {
        name string
        root *Node
    }{
        {
            name: "single node",
            root: &Node{Val: 1},
        },
        {
            name: "two connected nodes",
            root: func() *Node {
                n1 := &Node{Val: 1}
                n2 := &Node{Val: 2}
                n1.Neighbors = []*Node{n2}
                n2.Neighbors = []*Node{n1}
                return n1
            }(),
        },
        {
            name: "triangle cycle",
            root: func() *Node {
                n1 := &Node{Val: 1}
                n2 := &Node{Val: 2}
                n3 := &Node{Val: 3}
                n1.Neighbors = []*Node{n2, n3}
                n2.Neighbors = []*Node{n1, n3}
                n3.Neighbors = []*Node{n1, n2}
                return n1
            }(),
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            cloned := cloneGraphBFS(tc.root)
            if tc.root == nil {
                if cloned != nil {
                    t.Errorf("Expected nil, got %v", cloned)
                }
                return
            }

            if cloned == tc.root {
                t.Errorf("Expected a deep clone, got same pointer")
            }

            // A simple check: the structure should match via traversal
            if !utils.GraphsEqual(tc.root, cloned, map[*Node]*Node{}) {
                t.Errorf("Cloned graph does not match original for test %s", tc.name)
            }

            cloned = cloneGraphDFS(tc.root)
            if tc.root == nil {
                if cloned != nil {
                    t.Errorf("Expected nil, got %v", cloned)
                }
                return
            }

            if cloned == tc.root {
                t.Errorf("Expected a deep clone, got same pointer")
            }

            // A simple check: the structure should match via traversal
            if !utils.GraphsEqual(tc.root, cloned, map[*Node]*Node{}) {
                t.Errorf("Cloned graph does not match original for test %s", tc.name)
            }
        })
    }
}
