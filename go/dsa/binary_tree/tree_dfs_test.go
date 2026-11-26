package binarytree

import (
	. "gocode/dsa/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_diameterOfBinaryTree(t *testing.T) {
	tests := []struct {
		name     string
		root     *TreeNode
		expected int
	}{
		{
			name:     "empty tree",
			root:     nil,
			expected: 0,
		},
		{
			name:     "single node",
			root:     &TreeNode{Val: 1},
			expected: 0,
		},
		{
			name: "simple tree",
			//    1
			//   / \
			//  2   3
			// / \
			//4   5
			root: &TreeNode{
				Val: 1,
				Left: &TreeNode{
					Val:   2,
					Left:  &TreeNode{Val: 4},
					Right: &TreeNode{Val: 5},
				},
				Right: &TreeNode{Val: 3},
			},
			expected: 3, // path [4 → 2 → 1 → 3]
		},
		{
			name: "skewed tree",
			// 1
			//  \
			//   2
			//    \
			//     3
			//      \
			//       4
			root: &TreeNode{
				Val: 1,
				Right: &TreeNode{
					Val: 2,
					Right: &TreeNode{
						Val:   3,
						Right: &TreeNode{Val: 4},
					},
				},
			},
			expected: 3, // edges along the chain
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := diameterOfBinaryTree(tc.root)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func Test_maxDepth(t *testing.T) {
	tests := []struct {
		name     string
		root     *TreeNode
		expected int
	}{
		{
			name:     "empty tree",
			root:     nil,
			expected: 0,
		},
		{
			name:     "single node",
			root:     &TreeNode{Val: 1},
			expected: 1,
		},
		{
			name: "balanced tree",
			//    1
			//   / \
			//  2   3
			// / \
			//4   5
			root: &TreeNode{
				Val: 1,
				Left: &TreeNode{
					Val:   2,
					Left:  &TreeNode{Val: 4},
					Right: &TreeNode{Val: 5},
				},
				Right: &TreeNode{Val: 3},
			},
			expected: 3,
		},
		{
			name: "skewed tree",
			// 1
			//  \
			//   2
			//    \
			//     3
			//      \
			//       4
			root: &TreeNode{
				Val: 1,
				Right: &TreeNode{
					Val: 2,
					Right: &TreeNode{
						Val:   3,
						Right: &TreeNode{Val: 4},
					},
				},
			},
			expected: 4,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := maxDepth(tc.root)
			assert.Equal(t, tc.expected, got)
		})
	}
}
