package binarytree

import (
	. "gocode/dsa/types"
)

/**
 * 543. Diameter of Binary Tree
 */
func diameterOfBinaryTree(root *TreeNode) int {
	res := 0
	var dfs func(*TreeNode) int
	dfs = func(node *TreeNode) int {
		if node == nil {
			return 0
		}
		left, right := dfs(node.Left), dfs(node.Right)
		res = max(res, left+right)  // diameter = edges through node
		return 1 + max(left, right) // height = edges to deepest leaf
	}
	dfs(root)
	return res
}

/**
 * 104. Maximum Depth of Binary Tree
 */
func maxDepth(root *TreeNode) int {
	if root == nil {
		return 0
	}
	return 1 + max(maxDepth(root.Left), maxDepth(root.Right))
}
