package tree

import (
	. "gocode/types"
	"math"
)

/**
 * 124. Binary Tree Maximum Path Sum
 *
 * Postorder traversal
 */
func maxPathSum(root *TreeNode) int {
	maxSum := math.MinInt

	var postOrderTraversal func(*TreeNode) int
	postOrderTraversal = func(root *TreeNode) int {
		if root == nil {
			return 0
		}

		leftMaxSum := max(postOrderTraversal(root.Left), 0)
		rightMaxSum := max(postOrderTraversal(root.Right), 0)
		pathSum := leftMaxSum + rightMaxSum + root.Val
		maxSum = max(maxSum, pathSum)

		return root.Val + max(leftMaxSum, rightMaxSum)
	}

	postOrderTraversal(root)
	return maxSum
}
