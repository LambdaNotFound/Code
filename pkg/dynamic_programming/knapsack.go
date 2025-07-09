package dynamic_programming

/**
 * 53. Maximum Subarray
 *
 * Given an integer array nums, find the subarray with the largest sum, and return its sum.
 *
 * Kadane's algorithm
 *
 * DynamicProgramming, Time: O(n * sum), Space: O(n)
 *     dp[i] stores the maximum subarray ending at i
 *     dp[i] = dp[i - 1] + nums[i] if dp[i - 1] > 0
 *                                 else dp[i] = nums[i]
 */
func maxSubArray(nums []int) int {
    globalMax, curSum := nums[0], 0 // opt space, replace dp[i] w/ a var
    for _, num := range nums {
        if curSum < 0 {
            curSum = 0
        }
        curSum += num
        if curSum > globalMax {
            globalMax = curSum
        }
    }
    return globalMax
}

/**
 * 416. Partition Equal Subset Sum
 *
 * Given an integer array nums, return true if you can partition the
 * array into two subsets such that the sum of the elements in both subsets is equal or false otherwise
 *
 * DynamicProgramming, Time: O(n * sum), Space: O(n)
 *     dp[i] stores if sum i can be partitioned to 2 equal subsets
 *     target is the sum of the nums in the array
 *     dp[i] == true if dp[target - i] is true
 */
func canPartition(nums []int) bool {
    totalSum := 0
    for _, num := range nums {
        totalSum += num
    }

    // If total sum is odd, can't partition into equal subsets
    if totalSum%2 != 0 {
        return false
    }

    target := totalSum / 2
    dp := make([]bool, target+1)
    dp[0] = true // Base case: sum 0 can always be achieved

    for _, num := range nums {
        // Iterate backwards to avoid overwriting values we need to check
        for j := target; j >= num; j-- {
            if dp[j-num] {
                dp[j] = true
            }
        }
    }

    return dp[target]
}
