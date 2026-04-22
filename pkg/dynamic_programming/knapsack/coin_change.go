package knapsack

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
	for i := 1; i < amount+1; i++ {
		dp[i] = math.MaxInt
	}

	for _, c := range coins { //
		for a := 1; a < amount+1; a++ {
			if a-c >= 0 && dp[a-c] != math.MaxInt {
				dp[a] = min(dp[a], 1+dp[a-c])
			}
		}
	}

	if dp[amount] == math.MaxInt {
		return -1
	} else {
		return dp[amount]
	}
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

func coinChangeRecursionMemoization(coins []int, amount int) int {
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

/**
 * 518. Coin Change II
 *
 * Return the number of combinations that make up that amount.
 * If that amount of money cannot be made up by any combination of the coins, return 0.
 */
func change(amount int, coins []int) int {
	dp := make([]int, amount+1)
	dp[0] = 1
	for _, c := range coins { // coins outer → each coin considered in fixed order
		for a := 1; a <= amount; a++ {
			if a >= c {
				dp[a] += dp[a-c]
			}
		}
	}
	return dp[amount]
}

func changeRecursion(amount int, coins []int) int {
	var dfs func(amount, pos int) int
	dfs = func(amount, pos int) int {
		if amount == 0 {
			return 1
		}
		if amount < 0 {
			return 0
		}
		if pos == len(coins) {
			return 0
		}
		return dfs(amount-coins[pos], pos) + dfs(amount, pos+1)
	}
	return dfs(amount, 0)
}

func changeRecursionMemoization(amount int, coins []int) int {
	memo := make([][]*int, len(coins))
	for i := range memo {
		memo[i] = make([]*int, amount+1)
	}
	var dfs func(amount, pos int) int
	dfs = func(amount, pos int) int {
		if amount == 0 {
			return 1
		}
		if amount < 0 {
			return 0
		}
		if pos == len(coins) {
			return 0
		}
		if memo[pos][amount] != nil {
			return *memo[pos][amount]
		}
		res := dfs(amount-coins[pos], pos) + dfs(amount, pos+1)
		memo[pos][amount] = &res
		return res
	}
	return dfs(amount, 0)
}

func change2DDP(amount int, coins []int) int {
	dp := make([][]int, len(coins)+1)
	for i := range dp {
		dp[i] = make([]int, amount+1)
		dp[i][0] = 1
	}

	for c := 1; c <= len(coins); c++ {
		for a := 1; a <= amount; a++ {
			if a >= coins[c-1] {
				dp[c][a] = dp[c-1][a] + dp[c][a-coins[c-1]]
			} else {
				dp[c][a] = dp[c-1][a]
			}
		}
	}
	return dp[len(coins)][amount]
}
