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
