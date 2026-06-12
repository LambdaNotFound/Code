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

/**
 * 101. Symmetric Tree
 */
func isSymmetric(root *TreeNode) bool {
	var dfs func(p, q *TreeNode) bool
	dfs = func(p, q *TreeNode) bool {
		if p == nil || q == nil {
			return p == q
		}
		return p.Val == q.Val && dfs(p.Left, q.Right) && dfs(p.Right, q.Left)
	}

	return dfs(root.Left, root.Right)
}

/**
 * 110. Balanced Binary Tree
 *
 * Time: O(n), each node visited once
 */
func isBalanced(root *TreeNode) bool {
	var dfs func(*TreeNode) int
	// dfs returns height if subtree is balanced, -1 if not
	dfs = func(root *TreeNode) int {
		if root == nil {
			return 0
		}

		leftHeight := dfs(root.Left)
		if leftHeight == -1 {
			return -1 // left subtree already unbalanced — early exit
		}

		rightHeight := dfs(root.Right)
		if rightHeight == -1 {
			return -1 // right subtree already unbalanced — early exit
		}

		if leftHeight-rightHeight > 1 || rightHeight-leftHeight > 1 {
			return -1 // current node unbalanced
		}

		return max(leftHeight, rightHeight) + 1 // return height if balanced
	}

	return dfs(root) != -1
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

/**
 * 1145. Binary Tree Coloring Game
 *
 * In each turn, that player chooses a node of their color (red if player 1, blue if player 2) and colors an uncolored neighbor of the chosen node (either the left child, right child, or parent of the chosen node.)
 * BLOCKS the subtree
 *
 */

/**
 * Time: O(n) — single DFS traversal visiting every node once.
 * Space: O(h) — call stack depth equals tree height h. O(log n) for a balanced tree, O(n) worst case for a skewed tree.
 */
func btreeGameWinningMove(root *TreeNode, n int, x int) bool {
	leftCount, rightCount, totalCount := 0, 0, n
	var dfs func(node *TreeNode) int
	dfs = func(node *TreeNode) int {
		if node == nil {
			return 0
		}

		l, r := dfs(node.Left), dfs(node.Right)
		if node.Val == x {
			leftCount, rightCount = l, r
		}

		return l + r + 1
	}

	dfs(root)
	return max(max(leftCount, rightCount), totalCount-1-leftCount-rightCount) > n/2
}
