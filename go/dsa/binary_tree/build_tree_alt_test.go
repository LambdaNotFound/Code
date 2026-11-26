package binarytree

import (
	. "gocode/dsa/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

// treeEqual checks structural and value equality of two binary trees.
func Test_buildTreeAlt(t *testing.T) {
	testCases := []struct {
		name     string
		preorder []int
		inorder  []int
		expected *TreeNode
	}{
		{
			name:     "leetcode_example",
			preorder: []int{3, 9, 20, 15, 7},
			inorder:  []int{9, 3, 15, 20, 7},
			expected: &TreeNode{
				Val:  3,
				Left: &TreeNode{Val: 9},
				Right: &TreeNode{
					Val:   20,
					Left:  &TreeNode{Val: 15},
					Right: &TreeNode{Val: 7},
				},
			},
		},
		{
			name:     "single_node",
			preorder: []int{1},
			inorder:  []int{1},
			expected: &TreeNode{Val: 1},
		},
		{
			name:     "left_skewed",
			preorder: []int{1, 2, 3},
			inorder:  []int{3, 2, 1},
			expected: &TreeNode{
				Val: 1,
				Left: &TreeNode{
					Val:  2,
					Left: &TreeNode{Val: 3},
				},
			},
		},
		{
			name:     "right_skewed",
			preorder: []int{1, 2, 3},
			inorder:  []int{1, 2, 3},
			expected: &TreeNode{
				Val:   1,
				Right: &TreeNode{Val: 2, Right: &TreeNode{Val: 3}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := buildTreeAlt(tc.preorder, tc.inorder)
			assert.True(t, treeEqual(tc.expected, result))
		})
	}
}

// treeEqual is copied from design/serialize_deserialize_test.go, which keeps the original.
func treeEqual(a, b *TreeNode) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.Val == b.Val && treeEqual(a.Left, b.Left) && treeEqual(a.Right, b.Right)
}
