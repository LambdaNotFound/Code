package graph

import (
	"testing"

	. "gocode/types"

	"github.com/stretchr/testify/assert"
)

func Test_distanceK(t *testing.T) {
	tests := []struct {
		name     string
		buildFn  func() (*TreeNode, *TreeNode) // returns (root, target)
		k        int
		expected []int
	}{
		{
			name: "leetcode_example",
			buildFn: func() (*TreeNode, *TreeNode) {
				//        3
				//       / \
				//      5   1
				//     / \ / \
				//    6  2 0   8
				//      / \
				//     7   4
				n7 := &TreeNode{Val: 7}
				n4 := &TreeNode{Val: 4}
				n6 := &TreeNode{Val: 6}
				n2 := &TreeNode{Val: 2, Left: n7, Right: n4}
				n0 := &TreeNode{Val: 0}
				n8 := &TreeNode{Val: 8}
				n5 := &TreeNode{Val: 5, Left: n6, Right: n2}
				n1 := &TreeNode{Val: 1, Left: n0, Right: n8}
				root := &TreeNode{Val: 3, Left: n5, Right: n1}
				return root, n5 // target = 5
			},
			k:        2,
			expected: []int{7, 4, 1},
		},
		{
			name: "k_zero_returns_target",
			buildFn: func() (*TreeNode, *TreeNode) {
				root := &TreeNode{Val: 1, Left: &TreeNode{Val: 2}, Right: &TreeNode{Val: 3}}
				return root, root // target = root
			},
			k:        0,
			expected: []int{1},
		},
		{
			name: "k_larger_than_depth",
			buildFn: func() (*TreeNode, *TreeNode) {
				root := &TreeNode{Val: 1, Left: &TreeNode{Val: 2}}
				return root, root
			},
			k:        5,
			expected: []int{},
		},
		{
			name: "target_is_leaf",
			buildFn: func() (*TreeNode, *TreeNode) {
				leaf := &TreeNode{Val: 4}
				root := &TreeNode{Val: 1, Left: &TreeNode{Val: 2, Left: &TreeNode{Val: 3, Left: leaf}}}
				return root, leaf
			},
			k:        2,
			expected: []int{2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root, target := tt.buildFn()
			got := distanceKClaude(root, target, tt.k)
			assert.ElementsMatch(t, tt.expected, got)
		})
	}
}

func Test_zigzagLevelOrder(t *testing.T) {
	tests := []struct {
		name     string
		input    []any
		expected [][]int
	}{
		{
			name:     "leetcode_example",
			input:    []any{3, 9, 20, nil, nil, 15, 7},
			expected: [][]int{{3}, {20, 9}, {15, 7}},
		},
		{
			name:     "single_node",
			input:    []any{1},
			expected: [][]int{{1}},
		},
		{
			name:     "empty_tree",
			input:    []any{},
			expected: [][]int{},
		},
		{
			name:     "complete_3_levels",
			input:    []any{1, 2, 3, 4, 5, 6, 7},
			expected: [][]int{{1}, {3, 2}, {4, 5, 6, 7}},
		},
		{
			name:     "two_levels",
			input:    []any{1, 2, 3},
			expected: [][]int{{1}, {3, 2}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := buildTree(tt.input)
			assert.Equal(t, tt.expected, zigzagLevelOrder(root))
		})
	}
}
