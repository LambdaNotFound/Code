package tree

import (
	. "gocode/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Example function signature:
// func lowestCommonAncestor(root, p, q *TreeNode) *TreeNode

func buildBST() *TreeNode {
    /*
           6
          / \
         2   8
        / \ / \
       0  4 7  9
         / \
        3   5
    */
    root := &TreeNode{Val: 6}
    root.Left = &TreeNode{Val: 2}
    root.Right = &TreeNode{Val: 8}
    root.Left.Left = &TreeNode{Val: 0}
    root.Left.Right = &TreeNode{Val: 4}
    root.Left.Right.Left = &TreeNode{Val: 3}
    root.Left.Right.Right = &TreeNode{Val: 5}
    root.Right.Left = &TreeNode{Val: 7}
    root.Right.Right = &TreeNode{Val: 9}
    return root
}

func TestLowestCommonAncestor(t *testing.T) {
    root := buildBST()

    tests := []struct {
        name     string
        p, q     *TreeNode
        expected int
    }{
        {"LCA of 2 and 8 is 6", root.Left, root.Right, 6},
        {"LCA of 2 and 4 is 2", root.Left, root.Left.Right, 2},
        {"LCA of 3 and 5 is 4", root.Left.Right.Left, root.Left.Right.Right, 4},
        {"LCA of 0 and 5 is 2", root.Left.Left, root.Left.Right.Right, 2},
        {"LCA of 7 and 9 is 8", root.Right.Left, root.Right.Right, 8},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := lowestCommonAncestor(root, tt.p, tt.q)
            assert.NotNil(t, result)
            assert.Equal(t, tt.expected, result.Val)
        })
    }
}
