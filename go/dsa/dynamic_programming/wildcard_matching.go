package dynamic_programming

/**
 * 44. Wildcard Matching
 *
 * Given an input string (s) and a pattern (p), implement wildcard pattern matching with support for '?' and '*' where:
 *
 * '?' Matches any single character.
 * '*' Matches any sequence of characters (including the empty sequence).
 * The matching should cover the entire input string (not partial).
 *
 */

/*
 State: dp[i][j] = s[0..i-1] matches p[0..j-1]

 Transitions:

 p[j-1] == '*' → dp[i][j-1] (star matches empty) or dp[i-1][j] (star consumes s[i-1], stays at same * to potentially consume more)
 p[j-1] == '?' or p[j-1] == s[i-1] → dp[i-1][j-1]
 Base: dp[0][0] = true; dp[0][j] = true only if p[0..j-1] is all *s.

 T/S: O(m·n) time, O(m·n) space. Space can be reduced to O(n) with a rolling 1D array if needed.
*/

// dp[i][j] = s[0..i-1] matches p[0..j-1]
// '*' can match zero chars (dp[i][j-1]) or one more char (dp[i-1][j])

func isMatch(s string, p string) bool {
	m, n := len(s), len(p)
	dp := make([][]bool, m+1)
	for i := range dp {
		dp[i] = make([]bool, n+1)
	}
	dp[0][0] = true
	for j := 1; j <= n; j++ {
		if p[j-1] == '*' {
			dp[0][j] = dp[0][j-1]
		}
	}
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if p[j-1] == '*' {
				dp[i][j] = dp[i][j-1] || dp[i-1][j]
			} else if p[j-1] == '?' || p[j-1] == s[i-1] {
				dp[i][j] = dp[i-1][j-1]
			}
		}
	}
	return dp[m][n]
}

/* bool isMatch(string s, string p) {
    int m = s.size(), n = p.size();
    vector<vector<bool>> dp(m + 1, vector<bool>(n + 1));
    dp[0][0] = true;

    for (int i = 0; i <= m; ++i)
        for (int j = 1; j <= n; ++j)
            if (p[j - 1] == '*')
                dp[i][j] = dp[i][j - 1] || // '*' matches empty sequence
                           (i > 0 && dp[i - 1][j]); // '*' matches 1 preceding char
            else
                dp[i][j] = i > 0 && (s[i - 1] == p[j - 1] || p[j - 1] == '?') && dp[i - 1][j - 1];

    return dp[m][n];
} */
