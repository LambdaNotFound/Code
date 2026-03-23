package dynamic_programming

import "math"

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

/**
 * 64. Minimum Path Sum
 */
func minPathSumBottomUp(grid [][]int) int {
	m, n := len(grid), len(grid[0])
	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			switch {
			case r == 0 && c == 0:
				continue
			case r == 0:
				grid[r][c] += grid[r][c-1]
			case c == 0:
				grid[r][c] += grid[r-1][c]
			default:
				grid[r][c] += min(grid[r][c-1], grid[r-1][c])
			}
		}
	}
	return grid[m-1][n-1]
}

func minPathSumTopDown(grid [][]int) int {
	m, n := len(grid), len(grid[0])

	var dfs func(r, c int) int
	dfs = func(r, c int) int {
		if r == m || c == n {
			return math.MaxInt
		}
		if r == m-1 && c == n-1 {
			return grid[r][c]
		}
		return grid[r][c] + min(dfs(r, c+1), dfs(r+1, c))
	}

	return dfs(0, 0)
}
