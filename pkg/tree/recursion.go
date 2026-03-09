package tree

import (
	. "gocode/types"
)

/*
 * 100. Same Tree
 */
func isSameTree(p *TreeNode, q *TreeNode) bool {
	if p == nil && q == nil {
		return true
	}

	if (p != nil && q == nil) || (p == nil && q != nil) {
		return false
	}

	if p.Val != q.Val {
		return false
	}

	return isSameTree(p.Left, q.Left) && isSameTree(p.Right, q.Right)
}

/*
 * 572. Subtree of Another Tree
 *
 */
func isSubtree(root *TreeNode, subRoot *TreeNode) bool {
	if root == nil {
		return false
	}
	if subRoot == nil {
		return true
	}
	if isSameTree(root, subRoot) {
		return true
	}
	return isSubtree(root.Left, subRoot) || isSubtree(root.Right, subRoot)
}
