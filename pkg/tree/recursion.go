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
