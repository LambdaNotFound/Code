package tree

import (
	. "gocode/types"
	"math"
)

/**
 * 124. Binary Tree Maximum Path Sum
 *
 * 1. root.Val + leftMaxSum + rightMaxSum (as a whole tree)
 * 2. root.Val + leftMaxSum (as a left branch)
 * 3. root.Val + rightMaxSum (as a right branch)
 *
 * need to compute leftMaxSum and rightMaxSum first => Postorder traversal
 *
 * note that the path does not need to pass through the root
 */
func maxPathSum(root *TreeNode) int {
	maxSum := math.MinInt

	var postOrderTraversal func(*TreeNode) int
	postOrderTraversal = func(root *TreeNode) int {
		if root == nil {
			return 0
		}

		leftMaxSum := max(postOrderTraversal(root.Left), 0)   // exclude negative sum
		rightMaxSum := max(postOrderTraversal(root.Right), 0) // exclude negative sum
		pathSum := leftMaxSum + rightMaxSum + root.Val
		maxSum = max(maxSum, pathSum)

		return root.Val + max(leftMaxSum, rightMaxSum)
	}

	postOrderTraversal(root)
	return maxSum
}

/**
 * 113. Path Sum II
 *
 * Given the root of a binary tree and an integer targetSum,
 * return all root-to-leaf paths where the sum of the node values in the path equals targetSum.
 *
 * A root-to-leaf path is a path starting from the root and ending at any leaf node. A leaf is a node with no children.
 */
func pathSum(root *TreeNode, targetSum int) [][]int {
	res := make([][]int, 0)
	path := make([]int, 0)

	var dfs func(node *TreeNode, remaining int)
	dfs = func(node *TreeNode, remaining int) {
		if node == nil {
			return
		}

		path = append(path, node.Val)
		remaining -= node.Val

		if node.Left == nil && node.Right == nil && remaining == 0 {
			res = append(res, append([]int{}, path...)) // deep copy
		} else {
			dfs(node.Left, remaining)
			dfs(node.Right, remaining)
		}

		path = path[:len(path)-1]
	}

	dfs(root, targetSum)
	return res
}
