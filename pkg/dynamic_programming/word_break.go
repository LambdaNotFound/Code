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
 *
 * Gap index:  0   1   2   3   4   5   6   7   8
 *             |   |   |   |   |   |   |   |   |
 * Character:    l   e   e   t   c   o   d   e
 *
 * The moment you see dp of size n+1 over a string, anchor your brain: i is a gap, dp[i] is about the prefix s[0:i].
 */
func wordBreak(s string, wordDict []string) bool {
	wordMap := make(map[string]bool)
	for _, word := range wordDict {
		wordMap[word] = true
	}

	dp := make([]bool, len(s)+1)
	dp[0] = true
	for i := 1; i <= len(s); i++ {
		for j := 0; j < i; j++ {
			substr := s[j:i] // i is the gap index
			if wordMap[substr] && dp[j] {
				dp[i] = true
				break // no need to check remaining splits
			}
		}
	}

	return dp[len(s)]
}

func wordBreakCharIndex(s string, wordDict []string) bool {
	wordMap := make(map[string]bool)
	for _, word := range wordDict {
		wordMap[word] = true
	}

	dp := make([]bool, len(s)+1)
	dp[0] = true // base case: empty string is always valid
	for i := 0; i < len(s); i++ {
		for j := 0; j <= i; j++ {
			substr := s[j : i+1] // i is the char index
			if wordMap[substr] && dp[j] {
				dp[i+1] = true
				break // no need to check remaining splits
			}
		}
	}

	return dp[len(s)]
}

func wordBreakTrie(s string, wordDict []string) bool {
	root := NewTrie()
	for _, word := range wordDict {
		root.Insert(word)
	}

	dp := make([]bool, len(s)+1)
	dp[0] = true // base case: empty string is always valid
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

/**
 * 140. Word Break II
 *
 * Return all such possible sentences in any order.
 */
func wordBreak2(s string, wordDict []string) []string {
	wordMap := make(map[string]bool)
	for _, word := range wordDict {
		wordMap[word] = true
	}

	n := len(s)
	dp := make([]bool, n+1)
	table := make([][]string, n+1)
	for i := 1; i <= n; i++ {
		substr := s[0:i]
		if _, exist := wordMap[substr]; exist {
			dp[i] = true

			table[i] = append(table[i], substr)
		}

		for j := 0; j < i; j += 1 {
			substr = s[j:i]
			if _, exist := wordMap[substr]; exist && dp[j] {
				dp[i] = true

				for _, str := range table[j] { // multipe sentences at s[:j]
					table[i] = append(table[i], str+" "+substr)
				}
			}
		}
	}
	return table[n]
}
