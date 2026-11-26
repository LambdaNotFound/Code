package dynamic_programming

import "math"

/**
 * 1293. Shortest Path in a Grid with Obstacles Elimination
 */
func shortestPath(grid [][]int, k int) int {
	m, n := len(grid), len(grid[0])
	type State struct{ row, col, elim int }
	visited := make(map[State]bool)
	queue := []State{{0, 0, k}}
	visited[State{0, 0, k}] = true
	dirs := [][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
	steps := 0

	for len(queue) > 0 {
		size := len(queue)
		for i := 0; i < size; i++ {
			cur := queue[i]
			if cur.row == m-1 && cur.col == n-1 {
				return steps
			}
			for _, d := range dirs {
				r, c := cur.row+d[0], cur.col+d[1]
				if r < 0 || r >= m || c < 0 || c >= n {
					continue
				}
				elim := cur.elim - grid[r][c] // costs 1 elimination if obstacle
				if elim < 0 {
					continue
				}
				next := State{r, c, elim}
				if !visited[next] {
					visited[next] = true
					queue = append(queue, next)
				}
			}
		}
		queue = queue[size:]
		steps++
	}
	return -1
}

func shortestPathTopDown(grid [][]int, k int) int {
	m, n := len(grid), len(grid[0])
	dirs := [][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
	memo := make(map[[3]int]int)
	const INF = math.MaxInt

	var dfs func(row, col, elim int) int
	dfs = func(row, col, elim int) int {
		if row == m-1 && col == n-1 {
			return 0
		}
		if elim < 0 {
			return INF
		}
		key := [3]int{row, col, elim}
		if val, found := memo[key]; found {
			return val
		}

		best := INF
		for _, d := range dirs {
			r, c := row+d[0], col+d[1]
			if r < 0 || r >= m || c < 0 || c >= n {
				continue
			}
			sub := dfs(r, c, elim-grid[r][c])
			if sub != INF {
				best = min(best, sub+1)
			}
		}

		memo[key] = best
		return best
	}

	result := dfs(0, 0, k-grid[0][0])
	if result == INF {
		return -1
	}
	return result
}
