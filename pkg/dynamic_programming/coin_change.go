package dynamic_programming

import (
	. "gocode/utils"
)

/**
 * 322. Coin Change
 * Return the fewest number of coins that you need to make up that amount.
 * If that amount of money cannot be made up by any combination of the coins, return -1.
 *
 * You may assume that you have an infinite number of each kind of coin.
 *
 *    dp[i] stores the minimum number of coins used at amount i:
 *
 *    dp[0] = 0
 *    dp[i] = min(dp[i], dp[i - coins[j]] + 1) if (i - coins[j] >= 0)
 */
func coinChange(coins []int, amount int) int {
    dp := make([]int, amount+1)
    for i := range dp {
        dp[i] = amount + 1
    }
    dp[0] = 0
    for i := 1; i <= amount; i++ {
        for _, coin := range coins {
            if coin <= i {
                dp[i] = Min(dp[i], dp[i-coin]+1)
            }
        }
    }
    if dp[amount] > amount {
        return -1
    }
    return dp[amount]
}
