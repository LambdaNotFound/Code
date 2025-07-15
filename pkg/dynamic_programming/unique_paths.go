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
    for j := 0; j < n; j++ {
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
    dp := make([][]int, m)
    for i := range dp {
        dp[i] = make([]int, n)
    }
    for i := 0; i < m; i++ {
        dp[i][0] = 1
    }
    for j := 0; j < n; j++ {
        dp[0][j] = 1
    }

    for i := 1; i < m; i++ {
        for j := 1; j < n; j++ {
            dp[i][j] = dp[i-1][j] + dp[i][j-1]
        }
    }

    return dp[m-1][n-1]
}
