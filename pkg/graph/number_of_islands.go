package graph

/*
*

  - 200. Number of Islands
    *

  - Given an m x n 2D binary grid grid which represents a map of '1's (land) and '0's (water),

  - return the number of islands.
    *

  - DFS: Time: O(m x n), Space: O(m x n)

  - BFS: Time: O(m x n), Space: O(m x n)

  - UF:  Time: O(m x n), Space: O(m x n)
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
