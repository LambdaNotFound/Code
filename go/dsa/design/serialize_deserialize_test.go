package design

import (
	. "gocode/dsa/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

// treeEqual checks structural and value equality of two binary trees.
func treeEqual(a, b *TreeNode) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.Val == b.Val && treeEqual(a.Left, b.Left) && treeEqual(a.Right, b.Right)
}

func Test_Codec(t *testing.T) {
	testCases := []struct {
		name string
		root *TreeNode
	}{
		{
			name: "nil_tree",
			root: nil,
		},
		{
			name: "single_node",
			root: &TreeNode{Val: 1},
		},
		{
			name: "complete_tree",
			root: &TreeNode{
				Val:   1,
				Left:  &TreeNode{Val: 2},
				Right: &TreeNode{Val: 3},
			},
		},
		{
			name: "left_skewed",
			root: &TreeNode{
				Val:  1,
				Left: &TreeNode{Val: 2, Left: &TreeNode{Val: 3}},
			},
		},
		{
			name: "right_skewed",
			root: &TreeNode{
				Val:   1,
				Right: &TreeNode{Val: 2, Right: &TreeNode{Val: 3}},
			},
		},
		{
			name: "leetcode_example",
			root: &TreeNode{
				Val:   1,
				Left:  &TreeNode{Val: 2},
				Right: &TreeNode{Val: 3, Left: &TreeNode{Val: 4}, Right: &TreeNode{Val: 5}},
			},
		},
	}

	codec := Constructor()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			serialized := codec.serialize(tc.root)
			deserialized := codec.deserialize(serialized)
			assert.True(t, treeEqual(tc.root, deserialized))
		})
	}
}

// Test_CodecClaude exercises serializeClaude / deserializeClaude (BFS with "null" sentinel).
// Uses the same tree shapes as Test_Codec so round-trip correctness is consistent.
func Test_CodecClaude(t *testing.T) {
	testCases := []struct {
		name string
		root *TreeNode
	}{
		{name: "nil_tree", root: nil},
		{name: "single_node", root: &TreeNode{Val: 1}},
		{
			name: "complete_tree",
			root: &TreeNode{Val: 1, Left: &TreeNode{Val: 2}, Right: &TreeNode{Val: 3}},
		},
		{
			name: "left_skewed",
			root: &TreeNode{Val: 1, Left: &TreeNode{Val: 2, Left: &TreeNode{Val: 3}}},
		},
		{
			name: "right_skewed",
			root: &TreeNode{Val: 1, Right: &TreeNode{Val: 2, Right: &TreeNode{Val: 3}}},
		},
		{
			name: "leetcode_example",
			root: &TreeNode{
				Val:   1,
				Left:  &TreeNode{Val: 2},
				Right: &TreeNode{Val: 3, Left: &TreeNode{Val: 4}, Right: &TreeNode{Val: 5}},
			},
		},
		{
			// negative values exercise strconv paths
			name: "negative_values",
			root: &TreeNode{Val: -1, Left: &TreeNode{Val: -2}, Right: &TreeNode{Val: -3}},
		},
	}

	codec := Constructor()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			serialized := codec.serializeClaude(tc.root)
			deserialized := codec.deserializeClaude(serialized)
			assert.True(t, treeEqual(tc.root, deserialized))
		})
	}
}
