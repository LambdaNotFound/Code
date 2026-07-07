package prefixsum

import (
	. "gocode/types"
)

/**
 * 437. Path Sum III
 *
 * Given the root of a binary tree and an integer targetSum,
 * return the number of paths where the sum of the values along the path equals targetSum.
 *
 * "traveling only from parent nodes to child nodes)""
 *
 * Time: O(n), Space: O(n) — prefix map
 */
func pathSum(root *TreeNode, targetSum int) int {
	// prefixCount[sum] = number of paths from root with this prefix sum
	prefixCount := map[int]int{0: 1} // base case: empty path
	res := 0

	var dfs func(node *TreeNode, currSum int)
	dfs = func(node *TreeNode, currSum int) {
		if node == nil {
			return
		}

		currSum += node.Val

		// how many paths ending at current node sum to targetSum?
		// currSum - targetSum = prefix sum we need to have seen before
		res += prefixCount[currSum-targetSum]

		// record current prefix sum
		prefixCount[currSum]++

		dfs(node.Left, currSum)
		dfs(node.Right, currSum)

		// undo — remove current node's prefix sum when backtracking
		prefixCount[currSum]--
	}

	dfs(root, 0)
	return res
}
