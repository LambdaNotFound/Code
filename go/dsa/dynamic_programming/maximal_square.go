package dynamic_programming

/**
 * 221. Maximal Square
 *
 * Given an m x n binary matrix filled with 0's and 1's, find the largest square containing only 1's and return its area.
 *
 */
func maximalSquare(matrix [][]byte) int {
	m, n := len(matrix), len(matrix[0])
	dp := make([][]int, m+1)
	for r := range dp {
		dp[r] = make([]int, n+1)
	}

	side := 0
	for r := 1; r <= m; r++ {
		for c := 1; c <= n; c++ {
			if matrix[r-1][c-1] == '0' {
				continue
			}
			dp[r][c] = 1 + min(dp[r-1][c], min(dp[r-1][c-1], dp[r][c-1]))
			side = max(side, dp[r][c])
		}
	}
	return side * side
}
