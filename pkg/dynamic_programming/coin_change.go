package dynamic_programming

import "math"

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

func coinChangeRecursion(coins []int, amount int) int {
	var dfs func(remain int) int
	dfs = func(remain int) int {
		if remain == 0 {
			return 0
		}
		if remain < 0 {
			return math.MaxInt32
		}
		numberOfCoins := math.MaxInt32
		for _, c := range coins {
			if c > remain {
				continue
			}
			numberOfCoins = min(numberOfCoins, 1+dfs(remain-c))
		}
		return numberOfCoins
	}
	numberOfCoins := dfs(amount)
	if numberOfCoins == math.MaxInt32 {
		return -1
	}
	return numberOfCoins
}

func coinChangeRecursionnMemoization(coins []int, amount int) int {
	memo := make([]*int, amount+1)
	var dfs func(remain int) int
	dfs = func(remain int) int {
		if memo[remain] != nil {
			return *memo[remain]
		}
		if remain == 0 {
			return 0
		}
		if remain < 0 {
			return math.MaxInt32
		}
		numberOfCoins := math.MaxInt32
		for _, c := range coins {
			if c > remain {
				continue
			}
			numberOfCoins = min(numberOfCoins, 1+dfs(remain-c))
		}
		memo[remain] = &numberOfCoins
		return numberOfCoins
	}
	numberOfCoins := dfs(amount)
	if numberOfCoins == math.MaxInt32 {
		return -1
	}
	return numberOfCoins
}

func coinChange2DDP(coins []int, amount int) int {
	dp := make([][]int, len(coins)+1)
	for i := range dp {
		dp[i] = make([]int, amount+1)
	}
	for i := 1; i <= amount; i++ {
		dp[0][i] = math.MaxInt32
	}

	for c := 1; c <= len(coins); c++ {
		for a := 1; a <= amount; a++ {
			if a >= coins[c-1] {
				dp[c][a] = min(dp[c-1][a], 1+dp[c][a-coins[c-1]])
			} else {
				dp[c][a] = dp[c-1][a]
			}
		}
	}

	if dp[len(coins)][amount] == math.MaxInt32 {
		return -1
	}
	return dp[len(coins)][amount]
}

func coinChange1DDP(coins []int, amount int) int {
	dp := make([]int, amount+1)
	for i := 1; i <= amount; i++ {
		dp[i] = math.MaxInt32
	}

	for _, c := range coins {
		for a := 1; a <= amount; a++ {
			if a >= c {
				dp[a] = min(dp[a], 1+dp[a-c])
			}
		}
	}

	if dp[amount] == math.MaxInt32 {
		return -1
	}
	return dp[amount]
}
