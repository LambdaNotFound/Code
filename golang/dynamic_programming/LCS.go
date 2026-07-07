package dynamic_programming

import "sort"

/**
 * 1143. Longest Common Subsequence
 *
 * dp[i][j] = length of LCS of text1[0..i-1] and text2[0..j-1]
 *
 * if s1[i] == s2[j]:
 *   dp[i][j] = dp[i-1][j-1] + 1
 * else:
 *   dp[i][j] = max(dp[i-1][j],   // skip s1[i]
 *                  dp[i][j-1])   // skip s2[j]
 *
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
 *
 * if nums1[i] == nums2[j]:
 *    dp[i][j] = dp[i-1][j-1] + 1   // extend the streak
 * else:
 *    dp[i][j] = 0                  // streak broken — reset
 *
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

/**
 * 1048. Longest String Chain
 *
 * Input: words = ["a","b","ba","bca","bda","bdca"]
 * Output: 4
 * Explanation: One of the longest word chains is ["a","ba","bda","bdca"].
 *
 */
func longestStrChain(words []string) int {
	sort.Slice(words, func(i, j int) bool {
		return len(words[i]) < len(words[j])
	})

	dp := make(map[string]int)
	res := 0
	for _, word := range words {
		dp[word] = 1
		for i := 0; i < len(word); i++ {
			prev_word := word[:i] + word[i+1:]
			if val, exists := dp[prev_word]; exists {
				dp[word] = max(dp[word], val+1)
			}
		}
		res = max(res, dp[word])
	}
	return res
}
