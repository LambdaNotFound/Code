package dynamic_programming

/**
 * 62. Unique Paths
 *
 * The robot tries to move to the bottom-right corner (i.e., grid[m - 1][n - 1]).
 *
 * DynamicProgramming, Time: O(m x n), Space: O(m x n) => O(n)
 *     dp[i, j] stores the steps to reach [i, j]:
 *
 *     dp[i, j] = dp[i-1][j] + dp[i][j-1]
 */
func uniquePaths(m int, n int) int {
	dp := make([]int, n)
	for j := 0; j < n; j++ { // optimized for space
		dp[j] = 1
	}

	for i := 1; i < m; i++ {
		for j := 1; j < n; j++ {
			dp[j] += dp[j-1]
		}
	}

	return dp[n-1]
}

func uniquePathsNaive(m int, n int) int {
	ways := make([][]int, m)
	for r := range ways {
		ways[r] = make([]int, n)
		ways[r][0] = 1
	}
	for c := range ways[0] {
		ways[0][c] = 1
	}

	for r := 1; r < m; r++ {
		for c := 1; c < n; c++ {
			ways[r][c] = ways[r-1][c] + ways[r][c-1]
		}
	}
	return ways[m-1][n-1]
}
