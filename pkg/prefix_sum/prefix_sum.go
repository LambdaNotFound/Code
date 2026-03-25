package math

import (
	. "gocode/types"
)

/**
 * 238. Product of Array Except Self
 *
 * [a, b, c, d]
 *
 * LEFT  [1      a     ab     abc]    prefix
 * RIGHT [bcd    cd    d       1 ]    suffix
 * ANS   [bcd    acd   abd    abc]
 */
func productExceptSelf(nums []int) []int {
	res := make([]int, len(nums))

	multiplier := 1
	for i := 0; i < len(nums); i++ {
		res[i] = multiplier
		multiplier *= nums[i]
	}

	multiplier = 1
	for i := len(res) - 1; i >= 0; i-- {
		res[i] *= multiplier
		multiplier *= nums[i]
	}

	return res
}

func productExceptSelfClaude(nums []int) []int {
	n := len(nums)
	res := make([]int, n)
	// forward pass: res[i] = product of all elements to the LEFT of i
	res[0] = 1
	for i := 1; i < n; i++ {
		res[i] = res[i-1] * nums[i-1]
	}
	// backward pass: multiply res[i] by product of all elements to the RIGHT of i
	suffix := 1
	for i := n - 2; i >= 0; i-- {
		suffix *= nums[i+1]
		res[i] *= suffix
	}
	return res
}

/**
 * 525. Contiguous Array
 */
func findMaxLength(nums []int) int {
	// balance: +1 for 0, -1 for 1
	// equal 0s and 1s → balance returns to same value
	// store first occurrence of each balance value
	balanceToIndex := map[int]int{0: -1}

	res, balance := 0, 0
	for i, num := range nums {
		if num == 0 {
			balance++
		} else {
			balance--
		}

		if idx, found := balanceToIndex[balance]; found {
			res = max(res, i-idx)
		} else {
			balanceToIndex[balance] = i
		}
	}

	return res
}

func findMaxLengthPrefixSum(nums []int) int {
	// firstSeen[prefixSum] = earliest index where this prefix sum occurred
	firstSeen := map[int]int{0: -1} // base case: sum=0 seen before index 0

	prefixSum := 0
	maxLen := 0

	for i, num := range nums {
		// transform: 0 → -1, 1 → +1
		if num == 0 {
			prefixSum--
		} else {
			prefixSum++
		}

		if firstIdx, seen := firstSeen[prefixSum]; seen {
			// same prefix sum seen before → subarray [firstIdx+1..i] has sum 0
			maxLen = max(maxLen, i-firstIdx)
		} else {
			// first time seeing this prefix sum — record earliest index
			firstSeen[prefixSum] = i
		}
	}

	return maxLen
}

/**
 * 560. Subarray Sum Equals K
 */
func subarraySum(nums []int, k int) int {
	n := len(nums)
	prefixSum := make([]int, n+1)
	for i, num := range nums {
		prefixSum[i+1] = prefixSum[i] + num
	}

	// check all pairs (i, j) — O(n²)
	count := 0
	for i := 0; i < n; i++ {
		for j := i + 1; j <= n; j++ {
			if prefixSum[j]-prefixSum[i] == k {
				count++
			}
		}
	}

	return count
}

func subarraySumWithHashmap(nums []int, k int) int {
	// countMap[sum] = number of times this prefix sum has been seen
	countMap := map[int]int{0: 1} // base case: empty prefix has sum 0
	prefixSum := 0
	count := 0

	// Instead of checking all pairs, use a hashmap to look up
	// the complementc in O(1):
	// prefixSum[j] - prefixSum[i] = k
	// => prefixSum[i] = prefixSum[j] - k
	for _, num := range nums {
		prefixSum += num

		// how many previous prefixSums equal prefixSum-k?
		// each one forms a valid subarray ending at current index
		count += countMap[prefixSum-k]

		// record current prefix sum
		countMap[prefixSum]++
	}

	return count
}

/**
 * 930. Binary Subarrays With Sum
 *
 * [0, ..., i]            prefixSum[i]
 * [0, ..., ..., j]       prefixSum[j]
 *            delta       subarray's sum
 */
func numSubarraysWithSum(nums []int, goal int) int {
	countMap, prefixSum := map[int]int{0: 1}, 0
	count := 0
	for _, num := range nums {
		prefixSum += num

		count += countMap[prefixSum-goal]
		countMap[prefixSum]++
	}
	return count
}

/**
 * 437. Path Sum III
 *
 * Given the root of a binary tree and an integer targetSum,
 * return the number of paths where the sum of the values along the path equals targetSum.
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
