package prefixtree

/**
 * 3076. Shortest Uncommon Substring in an Array
 *
 * You are given an array arr of size n consisting of non-empty strings.
 *
 * Find a string array answer of size n such that:
 *
 * answer[i] is the shortest substring of arr[i] that does not occur as a substring in any other string in arr.
 * If multiple such substrings exist, answer[i] should be the lexicographically smallest.
 * And if no such substring exists, answer[i] should be an empty string.
 *
 * Return the array answer.
 */
type TrieNode struct {
	children map[rune]*TrieNode

	count    int // distinct strings with a substring reaching this node
	lastSeen int // index of last string that incremented count (dedup guard)
}

func newTrieNode() *TrieNode {
	return &TrieNode{
		children: make(map[rune]*TrieNode),
		lastSeen: -1,
	}
}

func shortestSubstrings(arr []string) []string {
	root := newTrieNode()

	// Build: insert every substring of every string into the trie.
	// For each string, walk from every starting position so that every
	// root-to-node path represents a substring, not just a prefix.
	for strIdx, s := range arr {
		for start := range s {
			node := root
			for end := start; end < len(s); end++ {
				ch := rune(s[end])
				if node.children[ch] == nil {
					node.children[ch] = newTrieNode()
				}
				node = node.children[ch]

				// Only count each string once per node, for deduplicate
				// e.g. abcabc
				if node.lastSeen != strIdx {
					node.count++
					node.lastSeen = strIdx
				}
			}
		}
	}

	// Query: for each string, find the shortest substring whose trie node
	// has count == 1 (unique to this string). Ties broken lexicographically.
	result := make([]string, len(arr))
	for strIdx, s := range arr {
		best := ""
		for start := range s {
			node := root
			for end := start; end < len(s); end++ {
				// Path is guaranteed to exist — we built it from this string.
				node = node.children[rune(s[end])]

				if node.count == 1 {
					candidate := s[start : end+1]
					if best == "" || len(candidate) < len(best) ||
						(len(candidate) == len(best) && candidate < best) {
						best = candidate
					}
					break // longer extensions are still unique but can't be shorter
				}
			}
		}
		result[strIdx] = best
	}

	return result
}
