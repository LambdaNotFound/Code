package dynamic_programming

import "strconv"

/**
 * 0/1 Knapsack, Unbounded Knapsack
 *
 * Loop Direction	Behavior
 * Forward (i++)	Reuses numbers (like "unlimited" items, a.k.a. unbounded knapsack)
 * Backward (i--)	Uses each number once, which is what we want for 0/1 knapsack
 *
 */

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
 * 322. Coin Change
 * Return the fewest number of coins that you need to make up that amount.
 * If that amount of money cannot be made up by any combination of the coins, return -1.
 *
 * You may assume that you have an infinite number of each kind of coin.
 *
 * DynamicProgramming, Time: O(n), Space: O(n)
 *     dp[i] stores the minimum number of coins used for amount i:
 *     coins[j] is the jth coin
 *
 *     dp[i] = min(dp[i], dp[i - coins[j]] + 1) if (i - coins[j] >= 0)
 *         base case: dp[0] = 0
 */
func coinChange(coins []int, amount int) int {
    dp := make([]int, amount+1)
    for i := range dp {
        dp[i] = amount + 1
    }
    dp[0] = 0
    for i := 1; i <= amount; i++ {
        for _, coin := range coins { // reuse coins of same value
            if coin <= i {
                dp[i] = min(dp[i], dp[i-coin]+1)
            }
        }
    }
    if dp[amount] > amount {
        return -1
    }
    return dp[amount]
}

/**
 * 416. Partition Equal Subset Sum
 *
 * Given an integer array nums, return true if you can partition the
 * array into two subsets such that the sum of the elements in both subsets is equal or false otherwise
 *
 * DynamicProgramming, Time: O(n * sum), Space: O(sum)
 *     dp[i] stores if sum i can be partitioned to 2 equal subsets
 *     target = sum of the nums in the array / 2
 *
 *     Sum1 [a, b, c...], Sum2 [x] => Sum1 == Sum2
 *
 *     dp[i] == true if dp[target - i] is true
 *         base case: dp[0] = 0
 */
func canPartition(nums []int) bool {
    totalSum := 0 // if total sum is odd, can't partition into equal subsets
    for _, num := range nums {
        totalSum += num
    }
    if totalSum%2 != 0 {
        return false
    }

    target := totalSum / 2
    dp := make([]bool, target+1)
    dp[0] = true // base case: sum 0 can always be achieved

    for _, num := range nums {
        // iterate backwards to ensure each number is only used once
        // same number used more than once violates the rules of 0/1 knapsack
        for j := target; j >= num; j-- { // [num, ..., target]
            dp[j] = dp[j] || dp[j-num] // if Sum == j-num then theres Sum == j
        }
    }
    return dp[target]
}

func canPartitionMemoization(nums []int) bool {
    sum := 0
    for _, num := range nums {
        sum += num
    }
    if sum%2 == 1 {
        return false
    }

    cache := make(map[string]bool)
    var validPartition func(int, int) bool
    validPartition = func(target, idx int) bool {
        if idx == len(nums) || target < 0 {
            return false
        }
        if target-nums[idx] == 0 {
            return true
        }

        key := strconv.Itoa(target) + "-" + strconv.Itoa(idx)
        if value, ok := cache[key]; ok {
            return value
        }

        cache[key] = validPartition(target, idx+1) || validPartition(target-nums[idx], idx+1)

        return cache[key]
    }

    return validPartition(sum/2, 0)
}

/**
 * 494. Target Sum
 *
 * Given an array of integers nums and an integer target,
 * return the number of ways to assign + and - signs to make the sum equal to target.
 *
 *     dp[index][value] stores num of ways to make the sum equal to value
 */
func findTargetSumWays(nums []int, target int) int {
    total := 0
    for _, num := range nums {
        total += num
    }
    if target > total || target < -total {
        return 0 // Out of possible range
    }

    offset := total // sum values can go negative, use an offset shift
    dp := make([][]int, len(nums)+1)
    for i := range dp {
        dp[i] = make([]int, 2*total+1) // range from -total to +total
    }

    dp[0][offset] = 1 // base case: 0 sum before starting
    // dp[0][0]         -total
    // dp[0][offset]    0
    // dp[0][2*offset]  total

    for i := 0; i < len(nums); i++ {
        for sum := -total; sum <= total; sum++ {
            curr := dp[i][sum+offset]
            if curr > 0 {
                dp[i+1][sum+nums[i]+offset] += curr
                dp[i+1][sum-nums[i]+offset] += curr
            }
        }
    }

    return dp[len(nums)][target+offset]
}

func findTargetSumWaysMemoization(nums []int, target int) int {
    memo := map[[2]int]int{}

    var dfs func(int, int) int
    dfs = func(idx int, sum int) int {
        if idx == len(nums) {
            if sum == target {
                return 1
            }
            return 0
        }

        key := [2]int{idx, sum}
        if val, ok := memo[key]; ok {
            return val
        }

        add := dfs(idx+1, sum+nums[idx])
        subtract := dfs(idx+1, sum-nums[idx])

        memo[key] = add + subtract
        return memo[key]
    }

    return dfs(0, 0)
}
