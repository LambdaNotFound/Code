package backtracking

/**
 * 79. Word Search
 *
 * Given an m x n grid of characters board and a string word, return true if word exists in the grid.
 *
 * Time: O(mn x 4^len) Space: O(L)
 */
func exist(board [][]byte, word string) bool {
    var dfs func(int, int, int) bool
    dfs = func(i, j, wordIdx int) bool {
        if wordIdx == len(word) {
            return true
        }

        if i < 0 || i >= len(board) || 
           j < 0 || j >= len(board[0]) || 
           board[i][j] != word[wordIdx] {
            return false
        }

        temp := board[i][j]
        board[i][j] = '*'
        found := dfs(i+1, j, wordIdx+1) ||
                 dfs(i-1, j, wordIdx+1) ||
                 dfs(i, j+1, wordIdx+1) ||
                 dfs(i, j-1, wordIdx+1)
        board[i][j] = temp
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
