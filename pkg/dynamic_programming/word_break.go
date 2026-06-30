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
 *
 *
 * Time: Naive (HashMap) TimeO(n² · L) <- Hash Computation: O(L)
 * Space: O(n + W)
 *     dp slice: O(n).
 *     wordMap: O(W) where W = total characters across all words in the dict.
 */
func wordBreak(s string, wordDict []string) bool {
	wordMap := make(map[string]bool)
	for _, word := range wordDict {
		wordMap[word] = true
	}

	dp := make([]bool, len(s)+1)
	dp[0] = true // base case: empty string, dp[length]
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
	dp[0] = true // base case: empty string, dp[length]
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

// Trie lookup approach:
// Time: O(n² + W), Space: O(n + W)
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
 *
 * dp[i]: substr of length i, can be break into the words in the dict
 * table[i]: sentences of substr of length i
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

/**
 * follow up: each word in the input list can only be used once
 *
 * par1 = ["back", "end", "front", "start"]
 * par2 = "backend" | "frontend" --> T ; backwards --> F
 *        "backback"  --> false,  word can not used multiple times
 *
 * par1 := []string{"back", "end", "backe", "front", "start", "back", "nds", "backend"}
 * par2 := "backend"
 * fmt.Println(wordMatch(par1, par2))
 *
 * par2 = "frontend"
 * fmt.Println(wordMatch(par1, par2))
 *
 * par2 = "backwards"
 * fmt.Println(wordMatch(par1, par2))
 *
 * par2 = "backe nds"
 * fmt.Println(wordMatch(par1, par2))
 *
 */

/**
 * s = "abba"
 * wordDict = ["ab", "a", "abb"]
 * Correct answer: true  ("abb" + "a")
 *
 * Trace:
 * i=1: "a" matches, freq["a"]-- → 0, dp[1]=true
 * i=3: "abb" matches, dp[3]=true
 * i=4: needs "a" at dp[3], but freq["a"]=0 already → returns false ✗
 */

// Time: O(2ⁿ · n) worst case, Branching ⇒ 2ⁿ⁻¹ paths & Per-node cost O(n)
// Space: O(n + W)
func wordBreakNoReuse(s string, words []string) bool {
	freq := map[string]int{}
	for _, w := range words {
		freq[w]++
	}

	var backtrack func(s string, start int) bool
	backtrack = func(s string, start int) bool {
		if start == len(s) {
			return true
		}
		for end := start + 1; end <= len(s); end++ {
			sub := s[start:end]
			if freq[sub] > 0 {
				freq[sub]--
				if backtrack(s, end) {
					return true
				}
				freq[sub]++ // restore on failure
			}
		}
		return false
	}

	return backtrack(s, 0)
}
