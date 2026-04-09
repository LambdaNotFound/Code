package dynamic_programming

/**
 * 198. House Robber
 */
func rob(nums []int) int {
	if len(nums) <= 1 {
		return nums[0]
	}

	dp := make([]int, len(nums))
	dp[0], dp[1] = nums[0], max(nums[0], nums[1])

	for i := 2; i < len(nums); i++ {
		dp[i] = max(dp[i-1], dp[i-2]+nums[i])
	}

	return dp[len(nums)-1]
}

/**
 * 213. House Robber II
 *
 * All houses at this place are arranged in a circle.
 */
func robII(nums []int) int {
	if len(nums) == 1 {
		return nums[0]
	}

	var rob func([]int) int
	rob = func(nums []int) int {
		if len(nums) == 0 {
			return 0
		}
		if len(nums) == 1 {
			return nums[0]
		}

		dp := make([]int, len(nums))
		dp[0], dp[1] = nums[0], max(nums[0], nums[1])

		for i := 2; i < len(nums); i++ {
			dp[i] = max(dp[i-1], dp[i-2]+nums[i])
		}

		return dp[len(nums)-1]
	}

	nums1 := nums[:len(nums)-1]
	nums2 := nums[1:]

	return max(rob(nums1), rob(nums2))
}
