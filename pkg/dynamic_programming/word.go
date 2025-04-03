package dynamic_programming

/**
 * 139. Word Break
 *
 * Given a non-empty string s and a dictionary wordDict containing a list of non-empty words,
 * determine if s can be segmented into a space-separated sequence of one or more dictionary words.
 *
 * DynamicProgramming, Time: O(n^2), Space: O(n)
 *     dp[i] stores if substring s[0, i - 1] is a valid sequence
 *     dp[i] == true if (1). s[0, i - 1] in dict OR
 *                      (2). dp[j] == true AND s[j, i - 1] in dict
 */
func wordBreak(s string, wordDict []string) bool {
    wordMap := make(map[string]bool)
    for _, word := range wordDict {
        wordMap[word] = true
    }

    n := len(s)
    dp := make([]bool, n+1)
    for i := 1; i <= n; i += 1 {
        substr := s[0:i]
        if _, exist := wordMap[substr]; exist {
            dp[i] = true
        } else {
            for j := 0; j < i; j += 1 {
                substr = s[j:i]
                if _, exist := wordMap[substr]; exist && dp[j] {
                    dp[i] = true
                }
            }
        }
    }
    return dp[n]
}

/**
 * 140. Word Break II
 *
 */
func wordBreak2(s string, wordDict []string) []string {
    wordMap := make(map[string]bool)
    for _, word := range wordDict {
        wordMap[word] = true
    }

    n := len(s)
    dp := make([]bool, n+1)
    table := make([][]string, n+1)
    for i := 1; i <= n; i += 1 {
        substr := s[0:i]
        if _, exist := wordMap[substr]; exist {
            dp[i] = true

            table[i] = append(table[i], substr)
        } else {
            for j := 0; j < i; j += 1 {
                substr = s[j:i]
                if _, exist := wordMap[substr]; exist && dp[j] {
                    dp[i] = true

                    for _, str := range table[j] {
                        table[i] = append(table[i], str+" "+substr)
                    }
                }
            }
        }
    }
    return table[n]
}
