package backtracking

/**
 * 79. Word Search
 *
 * Given an m x n grid of characters board and a string word, return true if word exists in the grid.
 *
 * Time: O(m * n * 3^len) Space: O(L) callstack
 */
func exist(board [][]byte, word string) bool {
	m, n := len(board), len(board[0])
	visited := make([][]bool, m)
	for i := range m {
		visited[i] = make([]bool, n)
	}

	var dfs func(int, int, int) bool
	dirs := [][]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
	dfs = func(row, col, pos int) bool {
		if pos == len(word) {
			return true
		}
		if row < 0 || row >= m || col < 0 || col >= n {
			return false
		}
		if visited[row][col] == true || board[row][col] != word[pos] {
			return false
		}

		visited[row][col] = true
		found := false
		for _, d := range dirs {
			r, c := row+d[0], col+d[1]
			if dfs(r, c, pos+1) { // explores all directions even after finding answer
				found = true
			}
		}
		visited[row][col] = false

		return found
	}

	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[0]); j++ {
			if board[i][j] == word[0] && dfs(i, j, 0) {
				return true
			}
		}
	}
	return false
}

func existClaude(board [][]byte, word string) bool {
	m, n := len(board), len(board[0])
	dirs := [][]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}

	// mark cell as visited by modifying board in-place
	// eliminates need for separate visited matrix → O(1) space
	var dfs func(row, col, pos int) bool
	dfs = func(row, col, pos int) bool {
		if pos == len(word) {
			return true
		}
		if row < 0 || row >= m || col < 0 || col >= n {
			return false
		}
		if board[row][col] != word[pos] {
			return false
		}

		// mark as visited by temporarily corrupting the cell
		temp := board[row][col]
		board[row][col] = '#'

		for _, d := range dirs {
			if dfs(row+d[0], col+d[1], pos+1) {
				board[row][col] = temp // restore before returning
				return true
			}
		}

		// restore cell for other DFS paths
		board[row][col] = temp
		return false
	}

	for row := 0; row < m; row++ {
		for col := 0; col < n; col++ {
			if dfs(row, col, 0) {
				return true
			}
		}
	}

	return false
}

/**
 * 212. Word Search II
 *
 * Backtracking w/ Trie
 */
type Trie struct {
	nodes     map[rune]*Trie
	endOfWord string // store word itself instead of bool — avoids rebuilding prefix
}

func Constructor() Trie {
	return Trie{nodes: make(map[rune]*Trie)}
}

func (t *Trie) Insert(word string) {
	for _, c := range word {
		if _, found := t.nodes[c]; !found {
			node := Constructor()
			t.nodes[c] = &node
		}
		t = t.nodes[c]
	}
	t.endOfWord = word // store word at terminal node — no need to track prefix in DFS
}

func findWords(board [][]byte, words []string) []string {
	trie := Constructor()
	for _, word := range words {
		trie.Insert(word)
	}

	m, n := len(board), len(board[0])
	dirs := [][]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
	res := make([]string, 0)

	var dfs func(row, col int, node *Trie)
	dfs = func(row, col int, node *Trie) {
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
		if nextNode.endOfWord != "" {
			res = append(res, nextNode.endOfWord)
			nextNode.endOfWord = "" // trie deletion — prevent duplicate matches
		}

		// mark visited via board modification — saves O(m×n) visited matrix
		board[row][col] = '#'
		for _, d := range dirs {
			dfs(row+d[0], col+d[1], nextNode)
		}
		board[row][col] = byte(c) // restore cell
	}

	for row := 0; row < m; row++ {
		for col := 0; col < n; col++ {
			dfs(row, col, &trie)
		}
	}

	return res
}
