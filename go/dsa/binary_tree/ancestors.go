package binarytree

import (
	. "gocode/dsa/types"
)

/**
 * 236. Lowest Common Ancestor of a Binary Tree
 */
func lowestCommonAncestor(root, p, q *TreeNode) *TreeNode {
	if root == nil || root == p || root == q {
		return root
	}

	// Process left and then right nodes
	left := lowestCommonAncestor(root.Left, p, q)
	right := lowestCommonAncestor(root.Right, p, q)

	// if both left and right aren't null, that means we found the targets on both sides of trees, means we need to return root
	if left != nil && right != nil {
		return root
	}
	// if we find in left, return left
	if left != nil {
		return left
	}
	return right
}

/**
 * 1026. Maximum Difference Between Node and Ancestor
 */
func maxAncestorDiff(root *TreeNode) int {
	res := 0
	if root == nil {
		return res
	}

	var dfs func(node *TreeNode, minVal, maxVal int)
	dfs = func(node *TreeNode, minVal, maxVal int) {
		if node == nil {
			return
		}

		res = max(res, node.Val-minVal)
		res = max(res, maxVal-node.Val)

		minVal = min(minVal, node.Val)
		maxVal = max(maxVal, node.Val)

		dfs(node.Left, minVal, maxVal)
		dfs(node.Right, minVal, maxVal)
	}

	dfs(root, root.Val, root.Val)
	return res
}
