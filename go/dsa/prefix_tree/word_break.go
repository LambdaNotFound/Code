package prefixtree

/**
 * 139. Word Break
 *
 * Given a non-empty string s and a dictionary wordDict containing a list of non-empty words,
 * determine if s can be segmented into a space-separated sequence of one or more dictionary words.
 *
 */
func wordBreakTrie(s string, wordDict []string) bool {
	root := NewTrie()
	for _, word := range wordDict {
		root.Insert(word)
	}

	dp := make([]bool, len(s)+1)
	dp[0] = true // dp[length]
	for i := 0; i < len(s); i++ {
		if !dp[i] {
			continue
		}

		// walk trie from position i character by character
		node := root
		for j := i; j < len(s); j++ {
			if _, found := node.nodes[rune(s[j])]; !found {
				break
			}
			node = node.nodes[rune(s[j])]

			if node.endOfWord {
				dp[j+1] = true
			}
		}
	}

	return dp[len(s)]
}
