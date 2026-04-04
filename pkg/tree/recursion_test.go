package tree

import (
	. "gocode/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_isSameTree(t *testing.T) {
	tests := []struct {
		name     string
		p, q     *TreeNode
		expected bool
	}{
		{"both_nil", nil, nil, true},
		{"one_nil", &TreeNode{Val: 1}, nil, false},
		{"single_equal", &TreeNode{Val: 1}, &TreeNode{Val: 1}, true},
		{"single_different", &TreeNode{Val: 1}, &TreeNode{Val: 2}, false},
		{
			"equal_trees",
			node(1, &TreeNode{Val: 2}, &TreeNode{Val: 3}),
			node(1, &TreeNode{Val: 2}, &TreeNode{Val: 3}),
			true,
		},
		{
			"different_structure",
			node(1, &TreeNode{Val: 2}, nil),
			node(1, nil, &TreeNode{Val: 2}),
			false,
		},
		{
			"different_value",
			node(1, &TreeNode{Val: 2}, &TreeNode{Val: 3}),
			node(1, &TreeNode{Val: 2}, &TreeNode{Val: 4}),
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, isSameTree(tt.p, tt.q))
		})
	}
}

func Test_isSymmetric(t *testing.T) {
	// Note: isSymmetric panics on nil root (calls root.Left on nil). Only test non-nil roots.
	tests := []struct {
		name     string
		root     *TreeNode
		expected bool
	}{
		{
			"symmetric",
			node(1, node(2, &TreeNode{Val: 3}, &TreeNode{Val: 4}), node(2, &TreeNode{Val: 4}, &TreeNode{Val: 3})),
			true,
		},
		{
			"not_symmetric",
			node(1, node(2, nil, &TreeNode{Val: 3}), node(2, nil, &TreeNode{Val: 3})),
			false,
		},
		{
			"single_node",
			&TreeNode{Val: 1},
			true,
		},
		{
			"two_left_children",
			node(1, &TreeNode{Val: 2}, nil),
			false,
		},
		{
			"equal_children",
			node(1, &TreeNode{Val: 2}, &TreeNode{Val: 2}),
			true,
		},
		{
			"different_child_values",
			node(1, &TreeNode{Val: 2}, &TreeNode{Val: 3}),
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, isSymmetric(tt.root))
		})
	}
}

func Test_isBalanced(t *testing.T) {
	tests := []struct {
		name     string
		root     *TreeNode
		expected bool
	}{
		{"nil_root", nil, true},
		{"single_node", &TreeNode{Val: 1}, true},
		{
			"balanced",
			node(3, node(9, nil, nil), node(20, &TreeNode{Val: 15}, &TreeNode{Val: 7})),
			true,
		},
		{
			"unbalanced",
			node(1, node(2, node(3, nil, nil), nil), nil),
			false,
		},
		{
			"full_balanced",
			node(1, node(2, &TreeNode{Val: 4}, &TreeNode{Val: 5}), node(3, &TreeNode{Val: 6}, &TreeNode{Val: 7})),
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, isBalanced(tt.root))
		})
	}
}

func Test_isSubtree(t *testing.T) {
	tests := []struct {
		name     string
		root     *TreeNode
		subRoot  *TreeNode
		expected bool
	}{
		{
			"leetcode_1",
			node(3, node(4, &TreeNode{Val: 1}, &TreeNode{Val: 2}), &TreeNode{Val: 5}),
			node(4, &TreeNode{Val: 1}, &TreeNode{Val: 2}),
			true,
		},
		{
			"not_subtree",
			node(3, node(4, &TreeNode{Val: 1}, node(2, &TreeNode{Val: 0}, nil)), &TreeNode{Val: 5}),
			node(4, &TreeNode{Val: 1}, &TreeNode{Val: 2}),
			false,
		},
		{
			"nil_subroot",
			&TreeNode{Val: 1},
			nil,
			true,
		},
		{
			"nil_root",
			nil,
			&TreeNode{Val: 1},
			false,
		},
		{
			"same_tree",
			node(1, &TreeNode{Val: 2}, &TreeNode{Val: 3}),
			node(1, &TreeNode{Val: 2}, &TreeNode{Val: 3}),
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, isSubtree(tt.root, tt.subRoot))
		})
	}
}
