package dynamic_programming

/**
 * 1143. Longest Common Subsequence
 */
func longestCommonSubsequence(text1 string, text2 string) int {
	if len(text1) > len(text2) {
		return longestCommonSubsequence(text2, text1)
	}

	m, n := len(text1), len(text2)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	res := 0
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if text1[i-1] == text2[j-1] {
				dp[i][j] = 1 + dp[i-1][j-1]
			} else {
				dp[i][j] = max(dp[i-1][j], dp[i][j-1])
			}

			res = max(res, dp[i][j])
		}
	}

	return res
}
