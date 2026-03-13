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
	dp[0] = true // base case: empty string is always valid

	for i := 1; i <= n; i++ {
		for j := 0; j < i; j++ {
			if dp[j] && wordMap[s[j:i]] {
				dp[i] = true
				break // no need to check other split points
			}
		}
	}

	return dp[n]
}

type Trie struct {
	nodes     map[rune]*Trie
	endOfWord bool
}

func NewTrie() *Trie {
	return &Trie{nodes: make(map[rune]*Trie)}
}

func (t *Trie) Insert(word string) {
	for _, c := range word {
		if _, found := t.nodes[c]; !found {
			t.nodes[c] = NewTrie()
		}
		t = t.nodes[c]
	}
	t.endOfWord = true
}

func wordBreakTrie(s string, wordDict []string) bool {
	// build trie from dictionary
	root := NewTrie()
	for _, word := range wordDict {
		root.Insert(word)
	}

	n := len(s)
	dp := make([]bool, n+1)
	dp[0] = true // base case: empty string is always valid

	for i := 0; i < n; i++ {
		// only expand from positions that are reachable
		if !dp[i] {
			continue
		}

		// walk trie from position i character by character
		// this finds all words in dict that start at position i
		node := root
		for j := i; j < n; j++ {
			c := rune(s[j])

			// no word in dict starts with this prefix — prune
			if _, found := node.nodes[c]; !found {
				break
			}
			node = node.nodes[c]

			// found a complete word s[i..j] in dictionary
			// mark dp[j+1] as reachable
			if node.endOfWord {
				dp[j+1] = true
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
		}

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
	return table[n]
}
