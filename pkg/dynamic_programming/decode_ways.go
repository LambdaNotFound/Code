package dynamic_programming

import "strconv"

/**
 * 91. Decode Ways
 *
 * dp[i] is the number of decode ways for string of length == i
 * goal: dp[len(s)]
 *
 * dp[i] = dp[i-1] if one digit > 0 + dp[i-2] if two digit >= 10 && <= 26
 *
 * strconv.Atoi(s[i-2 : i])
 */
func numDecodings(s string) int {
	dp := make([]int, len(s)+1)
	dp[0] = 1
	if s[0] != '0' {
		dp[1] = 1
	}
	for i := 2; i <= len(s); i++ {
		if oneDigit, _ := strconv.Atoi(s[i-1 : i]); oneDigit > 0 {
			dp[i] += dp[i-1]
		}
		if twoDigits, _ := strconv.Atoi(s[i-2 : i]); twoDigits >= 10 && twoDigits <= 26 {
			dp[i] += dp[i-2]
		}
	}
	return dp[len(s)]
}

func numDecodingsTopDown(s string) int {
	var dfs func(pos int) int
	dfs = func(pos int) int {
		if pos == len(s) {
			return 1
		}
		numWays := 0
		if pos+1 <= len(s) {
			if oneDigit, _ := strconv.Atoi(s[pos : pos+1]); oneDigit != 0 {
				numWays += dfs(pos + 1)
			}
		}
		if pos+2 <= len(s) {
			if twoDigits, _ := strconv.Atoi(s[pos : pos+2]); twoDigits >= 10 && twoDigits <= 26 {
				numWays += dfs(pos + 2)
			}
		}
		return numWays
	}
	return dfs(0)
}

func numDecodingsTopDownMemo(s string) int {
	memo := map[int]int{}
	var dfs func(pos int) int
	dfs = func(pos int) int {
		if pos == len(s) {
			return 1
		}
		if numWays, found := memo[pos]; found {
			return numWays
		}
		numWays := 0
		if pos+1 <= len(s) {
			if oneDigit, _ := strconv.Atoi(s[pos : pos+1]); oneDigit != 0 {
				numWays += dfs(pos + 1)
			}
		}
		if pos+2 <= len(s) {
			if twoDigits, _ := strconv.Atoi(s[pos : pos+2]); twoDigits >= 10 && twoDigits <= 26 {
				numWays += dfs(pos + 2)
			}
		}
		memo[pos] = numWays
		return numWays
	}
	return dfs(0)
}
