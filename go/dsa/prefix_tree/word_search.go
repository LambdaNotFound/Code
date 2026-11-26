package prefixtree

/**
 * 212. Word Search II
 *
 * Backtracking w/ Trie
 */
func findWords(board [][]byte, words []string) []string {
	trie := NewTrie()
	for _, word := range words {
		trie.Insert(word)
	}

	m, n := len(board), len(board[0])
	dirs := [][]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
	res := make([]string, 0)

	var backtrack func(row, col int, node *Trie)
	backtrack = func(row, col int, node *Trie) {
		if row < 0 || row >= m || col < 0 || col >= n {
			return
		}

		c := rune(board[row][col])

		// single trie lookup — replaces both Search and StartsWith
		nextNode, found := node.nodes[c]
		if !found {
			return // prefix not in trie — prune this path
		}

		// word found — add to result and delete from trie to prevent duplicates
		if nextNode.word != "" {
			res = append(res, nextNode.word)
			nextNode.word = "" // trie deletion — prevent duplicate matches
		}

		// mark visited via board modification — saves O(m×n) visited matrix
		board[row][col] = '#'
		for _, d := range dirs {
			backtrack(row+d[0], col+d[1], nextNode)
		}
		board[row][col] = byte(c) // restore cell

		// prune dead trie branch so future DFS calls skip exhausted subtrees
		if len(nextNode.nodes) == 0 && nextNode.word == "" {
			delete(node.nodes, c)
		}
	}

	for row := 0; row < m; row++ {
		for col := 0; col < n; col++ {
			backtrack(row, col, trie)
		}
	}

	return res
}
