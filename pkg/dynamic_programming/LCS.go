package dynamic_programming

/**
 * 1143. Longest Common Subsequence
 *
 * dp[i][j] = length of LCS of text1[0..i-1] and text2[0..j-1]
 */
func longestCommonSubsequence(text1, text2 string) int {
	m, n := len(text1), len(text2)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if text1[i-1] == text2[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				dp[i][j] = max(dp[i-1][j], dp[i][j-1])
			}
		}
	}

	return dp[m][n]
}

/**
 * 718. Maximum Length of Repeated Subarray
 * Longest Common Continuous Subarray
 *
 * Given two integer arrays nums1 and nums2, return the maximum length of a subarray that appears
 * in both arrays.
 *
 * dp[i][j] where dp[i][j] represents the length of the longest common suffix of A[0...i-1] and B[0...j-1]
 *                                                                      ending at i, j
 */
func findLength(nums1 []int, nums2 []int) int {
	m, n := len(nums1), len(nums2)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	res := 0
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if nums1[i-1] == nums2[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			}
			res = max(res, dp[i][j])
		}
	}

	return res
}
