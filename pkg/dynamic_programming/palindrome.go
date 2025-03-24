package palindrome

/**
 * 5. Longest Palindromic Substring
 *
 * Given a string s, find the longest palindromic substring in s. You may assume that the maximum length of s is 1000.
 *
 * DynamicProgramming <= Time: O(n^2), Space: O(n^2)
 *     dp[i][j] marks if string from i to j is palindrome
 *     dp[i][j] == true if (s[i] == s[j] && (j - i + 1 <= 2 || table[i + 1][j - 1))
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
 * Space optimize to O(n)
 *    ->
 *  |           [ i ][ j ]
 *  v [i+1][j-1]
 */
func longestPalindrome_opt_space(s string) string {
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
