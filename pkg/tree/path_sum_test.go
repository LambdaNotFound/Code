package tree

import (
	. "gocode/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_maxPathSum(t *testing.T) {
	tests := []struct {
		name     string
		root     *TreeNode
		expected int
	}{
		{
			"leetcode_1",
			node(1, &TreeNode{Val: 2}, &TreeNode{Val: 3}),
			6,
		},
		{
			"leetcode_2",
			node(-10, node(9, nil, nil), node(20, &TreeNode{Val: 15}, &TreeNode{Val: 7})),
			42,
		},
		{
			"single_node",
			&TreeNode{Val: 5},
			5,
		},
		{
			"all_negative",
			node(-1, &TreeNode{Val: -2}, &TreeNode{Val: -3}),
			-1,
		},
		{
			"one_positive_child",
			node(-3, nil, &TreeNode{Val: 1}),
			1,
		},
		{
			"left_right_path",
			node(5, node(4, node(11, &TreeNode{Val: 2}, &TreeNode{Val: 7}), nil), node(8, &TreeNode{Val: 13}, node(4, nil, &TreeNode{Val: 1}))),
			48,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, maxPathSum(tt.root))
		})
	}
}

func Test_pathSum_rootToLeaf(t *testing.T) {
	tests := []struct {
		name      string
		root      *TreeNode
		targetSum int
		expected  [][]int
	}{
		{
			"nil_root",
			nil, 5,
			[][]int{},
		},
		{
			"single_match",
			&TreeNode{Val: 5},
			5,
			[][]int{{5}},
		},
		{
			"single_no_match",
			&TreeNode{Val: 1},
			5,
			[][]int{},
		},
		{
			"leetcode_example",
			// Tree:   5
			//        / \
			//       4   8
			//      /   / \
			//     11  13   4
			//    /  \       \
			//   7    2       1
			// Paths summing to 22: [5,4,11,2] and [5,8,4,1... wait no
			// let me recalculate: 5+4+11+2=22, 5+8+4+1=18 no...
			// Actually [5,4,11,2] = 22, [5,8,9] not valid
			// LeetCode 113 example: targetSum=22 → [[5,4,11,2],[5,8,4,5]]? No...
			// Let me use LeetCode example correctly:
			//   5+4+11+2 = 22 ✓
			//   5+8+4+1 = 18 ✗
			// Actually second leaf with 4 is wrong. The tree has root=5, left=4 right=8
			// left.left=11, right.left=13 right.right=4(leaf with right child 1)
			// No, the official LeetCode example:
			//  root=[5,4,8,11,null,13,4,7,2,null,null,5,1], targetSum=22
			//  Paths: [5,4,11,2] and [5,8,4,5] -- hmm that's the LeetCode 437 (path sum III)?
			// For LeetCode 113: root=[5,4,8,11,null,13,4,7,2,null,null,null,1], targetSum=22
			//   [5,4,11,2] = 22 ✓
			node(5,
				node(4, node(11, &TreeNode{Val: 7}, &TreeNode{Val: 2}), nil),
				node(8, &TreeNode{Val: 13}, node(4, nil, &TreeNode{Val: 1})),
			),
			22,
			[][]int{{5, 4, 11, 2}},
		},
		{
			"two_paths",
			// Tree: 3
			//      / \
			//     2   2
			// Both leaves sum to 5 with root: [3,2], [3,2]
			node(3, &TreeNode{Val: 2}, &TreeNode{Val: 2}),
			5,
			[][]int{{3, 2}, {3, 2}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pathSum(tt.root, tt.targetSum)
			if len(tt.expected) == 0 {
				assert.Empty(t, got)
			} else {
				assert.ElementsMatch(t, tt.expected, got)
			}
		})
	}
}
