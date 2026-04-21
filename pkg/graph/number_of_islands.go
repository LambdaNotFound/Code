package graph

/*
 * 200. Number of Islands
 *
 * - Given an m x n 2D binary grid grid which represents a map of '1's (land) and '0's (water),
 * - return the number of islands.
 *
 * DFS: Time: O(m x n), Space: O(m x n)
 * BFS: Time: O(m x n), Space: O(m x n)
 * UF:  Time: O(m x n), Space: O(m x n)
 *
 * Complexity
 *   Time O(m × n)
 *   Space O(m × n) — call stack
 */
func numIslandsDFS(grid [][]byte) int {
	if len(grid) == 0 || len(grid[0]) == 0 {
		return 0
	}

	m, n := len(grid), len(grid[0])
	directions := [][]int{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}

	var dfs func(int, int)
	dfs = func(i, j int) {
		if grid[i][j] != '1' {
			return
		}

		if grid[i][j] == '1' {
			grid[i][j] = 'X'
		}
		for _, dir := range directions {
			row, col := i+dir[0], j+dir[1]
			if row >= 0 && row < m && col >= 0 && col < n {
				dfs(row, col)
			}
		}
	}

	res := 0
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if grid[i][j] == '1' {
				res++
				dfs(i, j)
			}
		}
	}
	return res
}

func numIslandsBFS(grid [][]byte) int {
	if len(grid) == 0 || len(grid[0]) == 0 {
		return 0
	}

	m, n := len(grid), len(grid[0])
	directions := [][]int{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}

	var bfs func(int, int)
	bfs = func(i, j int) {
		queue := [][]int{{i, j}}
		grid[i][j] = 'X'
		for len(queue) > 0 {
			cur := queue[0]
			queue = queue[1:]

			for _, dir := range directions {
				row, col := cur[0]+dir[0], cur[1]+dir[1]
				if (row >= 0 && row < m && col >= 0 && col < n) && grid[row][col] == '1' {
					grid[row][col] = 'X'
					queue = append(queue, []int{row, col})
				}
			}
		}
	}

	res := 0
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if grid[i][j] == '1' {
				res++
				bfs(i, j)
			}
		}
	}
	return res
}

func numIslandsUF(grid [][]byte) int {
	if len(grid) == 0 || len(grid[0]) == 0 {
		return 0
	}
	m, n := len(grid), len(grid[0])

	parent := make([]int, m*n)
	for i := range parent {
		parent[i] = i // parent array only (no rank/size)
	}

	var find func(int) int
	find = func(x int) int {
		if parent[x] != x {
			parent[x] = find(parent[x]) // path compression only
		}
		return parent[x]
	}

	union := func(a, b int) bool {
		ra, rb := find(a), find(b)
		if ra == rb {
			return false
		}
		parent[rb] = ra // no rank: arbitrarily attach rb under ra
		return true
	}

	dirs := [][]int{{1, 0}, {0, 1}} // Only check down and right to avoid double-union of same edge
	count := 0
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if grid[i][j] != '1' {
				continue
			}
			count++
			idx := i*n + j // flat 2D grid
			for _, d := range dirs {
				row, col := i+d[0], j+d[1]
				if row >= 0 && row < m && col >= 0 && col < n && grid[row][col] == '1' {
					if union(idx, row*n+col) {
						count--
					}
				}
			}
		}
	}
	return count
}

/**
 * 463. Island Perimeter
 */
func islandPerimeter(grid [][]int) int {
	directions := [][2]int{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}
	rows, cols := len(grid), len(grid[0])

	const (
		water   = 0
		land    = 1
		visited = 2
	)

	var dfs func(row, col int) int
	dfs = func(row, col int) int {
		// Out of bounds → this edge is on the grid boundary, contributes 1 to perimeter
		if row < 0 || row >= rows || col < 0 || col >= cols {
			return 1
		}
		// Water → this edge borders water, contributes 1 to perimeter
		if grid[row][col] == water {
			return 1
		}
		// Already counted this land cell's contributions
		if grid[row][col] == visited {
			return 0
		}

		grid[row][col] = visited

		perimeter := 0
		for _, dir := range directions {
			nextRow, nextCol := row+dir[0], col+dir[1]
			perimeter += dfs(nextRow, nextCol)
		}
		return perimeter
	}

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			if grid[row][col] == land {
				return dfs(row, col)
			}
		}
	}
	return 0
}

func islandPerimeterClaude(grid [][]int) int {
	dirs := [][2]int{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}
	m, n := len(grid), len(grid[0])
	perimeter := 0

	for row := 0; row < m; row++ {
		for col := 0; col < n; col++ {
			if grid[row][col] != 1 {
				continue
			}
			for _, dir := range dirs {
				nextRow, nextCol := row+dir[0], col+dir[1]
				// Out of bounds OR water → this side contributes to perimeter
				if nextRow < 0 || nextRow >= m ||
					nextCol < 0 || nextCol >= n ||
					grid[nextRow][nextCol] == 0 {
					perimeter++
				}
			}
		}
	}
	return perimeter
}
