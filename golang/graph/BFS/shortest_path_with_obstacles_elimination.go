package bfs

/**
 * 1293. Shortest Path in a Grid with Obstacles Elimination
 *
 * BFS state (row, col, eliminationsRemaining)
 *
 * Queue:  starts with (0, 0, k)           — top-left, full eliminations
 * Visited: visited[row][col][elim] = true  — 3D boolean array
 *
 * For each state (row, col, elim) popped from queue:
 *    if (row, col) == destination → return steps
 *
 *    for each neighbor (nextRow, nextCol):
 *        if neighbor is in bounds:
 *            if grid[nextRow][nextCol] == 0:          // open cell
 *                nextState = (nextRow, nextCol, elim)
 *            elif elim > 0:                           // obstacle, can eliminate
 *                nextState = (nextRow, nextCol, elim-1)
 *            else:
 *                skip                                 // obstacle, no eliminations left
 *
 *            if not visited[nextState]:
 *                visited[nextState] = true
 *                enqueue nextState
 */
func shortestPath(grid [][]int, k int) int {
	m, n := len(grid), len(grid[0])

	// optimization: if k is large enough to eliminate every obstacle on any path
	if k >= m+n-2 {
		return m + n - 2
	}

	type State struct {
		row, col, elim int
	}

	visited := make([][][]bool, m)
	for row := range visited {
		visited[row] = make([][]bool, n)
		for col := range visited[row] {
			visited[row][col] = make([]bool, k+1)
		}
	}

	queue := []State{{row: 0, col: 0, elim: k}}
	visited[0][0][k] = true
	dirs := [][2]int{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}
	steps := 0

	for len(queue) > 0 {
		steps++
		size := len(queue)
		for i := 0; i < size; i++ {
			cur := queue[0]
			queue = queue[1:]

			for _, dir := range dirs {
				nextRow, nextCol := cur.row+dir[0], cur.col+dir[1]

				if nextRow < 0 || nextRow >= m || nextCol < 0 || nextCol >= n {
					continue
				}
				if nextRow == m-1 && nextCol == n-1 {
					return steps
				}

				nextElim := cur.elim
				if grid[nextRow][nextCol] == 1 {
					// obstacle: consume one elimination
					if cur.elim == 0 {
						continue // can't eliminate, skip
					}
					nextElim--
				}

				if visited[nextRow][nextCol][nextElim] {
					continue
				}
				visited[nextRow][nextCol][nextElim] = true
				queue = append(queue, State{row: nextRow, col: nextCol, elim: nextElim})
			}
		}
	}

	return -1
}
