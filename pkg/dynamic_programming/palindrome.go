package dynamic_programming

/**
 * 5. Longest Palindromic Substring
 *
 * Given a string s, find the longest palindromic substring in s. You may assume that the maximum length of s is 1000.
 *
 * DynamicProgramming, Time: O(n^2), Space: O(n^2)
 *     dp[i][j] stores if substring from i to j is palindrome, s[i, j]
 *     dp[i][j] == true if (s[i] == s[j] AND (j - i + 1 <= 2 OR table[i + 1][j - 1] == true)
 *                                            a, aa, aba        substring[i+1, j-1] is a palindrome
 */
func longestPalindrome(s string) string {
	start, length, size := 0, 1, len(s)
	dp := make([][]bool, size)
	for i := 0; i < size; i++ {
		dp[i] = make([]bool, size)
	}

	for j := 0; j < size; j++ {
		for i := 0; i <= j; i++ {
			dp[i][j] = (s[i] == s[j]) && (j-i+1 <= 2 || dp[i+1][j-1])

			if dp[i][j] && (length < j-i+1) {
				length = j - i + 1
				start = i
			}
		}
	}

	return s[start : start+length]
}

/**
 * 647. Palindromic Substrings
 */
func countSubstrings(s string) int {
	n := len(s)
	dp := make([][]bool, n)
	for i := range dp {
		dp[i] = make([]bool, n)
	}

	count := 0
	for right := 0; right < n; right++ {
		for left := 0; left <= right; left++ {
			isEdgeCase := right-left+1 <= 2 // single char or pair
			isInnerPalin := !isEdgeCase && dp[left+1][right-1]
			if s[left] == s[right] && (isEdgeCase || isInnerPalin) {
				dp[left][right] = true
				count++
			}
		}
	}

	return count
}

/**
 * Optimize space complexity, Space: O(n)
 *    ->
 *  |           [ i ][ j ]
 *  v [i+1][j-1]
 */
func longestPalindrome_optimized(s string) string {
	start, length, size := 0, 1, len(s)
	dp := make([]bool, size)

	for j := 0; j < size; j++ {
		for i := 0; i <= j; i++ {
			dp[i] = (s[i] == s[j]) && (j-i+1 <= 2 || dp[i+1])

			if dp[i] && (length < j-i+1) {
				length = j - i + 1
				start = i
			}
		}
	}

	return s[start : start+length]
}
