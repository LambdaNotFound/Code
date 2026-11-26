package topologicalsort

/**
 * 329. Longest Increasing Path in a Matrix
 *
 * 1. BFS w/ in-degrees (Kahn’s algo)
 * 2. DFS + memo: dp[curr] = max(dp[curr], dp[neighbor] + 1)
 */
func longestIncreasingPath(matrix [][]int) int {
	m, n := len(matrix), len(matrix[0])
	indegree := make(map[[2]int]int)

	dirs := [][]int{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}
	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			for _, dir := range dirs {
				row, col := r+dir[0], c+dir[1]
				if row >= 0 && row < m && col >= 0 && col < n &&
					matrix[row][col] > matrix[r][c] {
					indegree[[2]int{row, col}]++
				}
			}
		}
	}

	queue := [][2]int{}
	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			if indegree[[2]int{r, c}] == 0 {
				queue = append(queue, [2]int{r, c})
			}
		}
	}

	level := 0
	for len(queue) > 0 {
		size := len(queue)
		for i := 0; i < size; i++ {
			cell := queue[0]
			queue = queue[1:]

			for _, dir := range dirs {
				row, col := cell[0]+dir[0], cell[1]+dir[1]
				if row >= 0 && row < m && col >= 0 && col < n &&
					matrix[row][col] > matrix[cell[0]][cell[1]] {
					indegree[[2]int{row, col}]--
					if indegree[[2]int{row, col}] == 0 {
						queue = append(queue, [2]int{row, col})
					}
				}
			}
		}
		level++
	}

	return level
}

func longestIncreasingPathMemoization(matrix [][]int) int {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return 0
	}

	m, n := len(matrix), len(matrix[0])
	memo := make([][]int, m)
	for i := range memo {
		memo[i] = make([]int, n)
	}

	dirs := [][]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}

	var dfs func(int, int) int
	dfs = func(row, col int) int {
		if memo[row][col] != 0 {
			return memo[row][col]
		}
		maxLen := 1
		for _, dir := range dirs {
			nextRow, nextCol := row+dir[0], col+dir[1]
			if nextRow >= 0 && nextRow < m && nextCol >= 0 && nextCol < n &&
				matrix[nextRow][nextCol] > matrix[row][col] {
				length := 1 + dfs(nextRow, nextCol)
				if length > maxLen {
					maxLen = length
				}
			}
		}
		memo[row][col] = maxLen
		return maxLen
	}

	longest := 0
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			longest = max(longest, dfs(i, j))
		}
	}
	return longest
}
