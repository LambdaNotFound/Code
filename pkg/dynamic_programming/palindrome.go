package palindrome

func longestPalindrome(s string) string {
	start, length, size := 0, 1, len(s)
	dp := make([][]bool, size)
	for i := 0; i < size; i++ {
		dp[i] = make([]bool, size)
	}

	for j := 0; j < size; j++ {
		for i := 0; i <= j; i++ {
			dp[i][j] = (s[i] == s[j]) && (j-i <= 2 || dp[i+1][j-1])

			if dp[i][j] && (length < j-i+1) {
				length = j - i + 1
				start = i
			}
		}
	}

	return s[start : start+length]
}
