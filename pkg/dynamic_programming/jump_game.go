package dynamic_programming

import "math"

/**
 * 55. Jump Game
 *
 * Greedy
 *
 * bool canJump(vector<int>& nums) {
 *     int maxReach = 0;
 *     for (int i = 0; i < nums.size(); ++i) {
 *         if (i > maxReach) return false;
 *         maxReach = max(maxReach, i + nums[i]);
 *     }
 *     return true;
 * }
 */
func canJump(nums []int) bool {
	table := make([]bool, len(nums))
	table[0] = true
	for i := 0; i < len(nums) && table[i] == true; i++ {
		for j := 1; j <= nums[i] && i+j < len(nums); j++ {
			table[i+j] = true
		}
	}
	return table[len(nums)-1]
}

/**
 * 45. Jump Game II
 */
func jump(nums []int) int {
	n := len(nums)
	dp := make([]int, n)
	for i := range dp {
		dp[i] = math.MaxInt
	}
	dp[0] = 0

	for i := 0; i < n; i++ {
		for j := 1; j <= nums[i] && i+j < n; j++ {
			dp[i+j] = min(dp[i+j], dp[i]+1)

		}
	}
	return dp[n-1]
}
