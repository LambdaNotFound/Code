package binarytree

import (
	. "gocode/dsa/types"
	"gocode/dsa/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// helper to build a TreeNode manually
func Test_lowestCommonAncestor(t *testing.T) {
        /*
           6
          / \
         2   8
        / \ / \
       0  4 7  9
         / \
        3   5
    */
    root := utils.BuildBST([]int{6, 2, 8, 0, 4, 7, 9, 3, 5})

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

func Test_maxAncestorDiff(t *testing.T) {
	tests := []struct {
		name     string
		root     *TreeNode
		expected int
	}{
		{
			name:     "nil_root",
			root:     nil,
			expected: 0,
		},
		{
			name:     "single_node",
			root:     &TreeNode{Val: 5},
			expected: 0,
		},
		{
			// path 8→3→1: |8-1|=7
			name:     "left_skewed",
			root:     node(8, node(3, node(1, nil, nil), nil), nil),
			expected: 7,
		},
		{
			// LC example: max diff is |8-1|=7
			name: "balanced",
			root: node(8,
				node(3, node(1, nil, nil), node(6, node(4, nil, nil), node(7, nil, nil))),
				node(10, nil, node(14, node(13, nil, nil), nil)),
			),
			expected: 7,
		},
		{
			// ascending right spine 1→2→3: max ancestor diff = |3-1|=2
			name:     "right_skewed",
			root:     node(1, nil, node(2, nil, node(3, nil, nil))),
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, maxAncestorDiff(tt.root))
		})
	}
}

// node builds a tree node; copied from bst/BST_test.go when 236 and 1026 moved here.
func node(val int, left, right *TreeNode) *TreeNode {
    return &TreeNode{Val: val, Left: left, Right: right}
}
