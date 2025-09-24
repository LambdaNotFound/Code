package graph

import (
	. "gocode/types"
	"gocode/utils"
	"testing"
)

func TestCloneGraph(t *testing.T) {
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
            cloned := cloneGraph(tc.root)

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
