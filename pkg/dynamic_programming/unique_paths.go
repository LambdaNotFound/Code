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
	table := make([][]int, m)
	for i := 0; i < m; i++ {
		table[i] = make([]int, n)
	}

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if i == 0 || j == 0 {
				table[i][j] = 1
			} else {
				table[i][j] = table[i-1][j] + table[i][j-1]
			}
		}
	}

	return table[m-1][n-1]
}
