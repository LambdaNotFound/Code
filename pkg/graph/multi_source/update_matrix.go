package multisource

import "math"

/**
 * 542. 01 Matrix
 * a multi-source BFS approach
 *
 * start BFS from zero cells
 *
 * Complexity
 * Time O(m × n) — each cell enqueued once
 * Space O(m × n) — queue + dist matrix
 */
func updateMatrix(mat [][]int) [][]int {
	m, n := len(mat), len(mat[0])
	dirs := [][]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}

	dist := make([][]int, m)
	for i := range dist {
		dist[i] = make([]int, n)
	}

	// seed queue with ALL zero cells simultaneously
	// mark non-zero cells as unvisited with sentinel value
	queue := make([][2]int, 0)
	for row := 0; row < m; row++ {
		for col := 0; col < n; col++ {
			if mat[row][col] == 0 {
				queue = append(queue, [2]int{row, col}) // seed
				dist[row][col] = 0
			} else {
				dist[row][col] = math.MaxInt // unvisited sentinel
			}
		}
	}

	// BFS outward from all zero cells
	for len(queue) > 0 {
		row, col := queue[0][0], queue[0][1]
		queue = queue[1:]

		for _, dir := range dirs {
			r, c := row+dir[0], col+dir[1]
			if r < 0 || r >= m || c < 0 || c >= n {
				continue
			}

			// only update if we found a shorter distance
			if dist[row][col]+1 < dist[r][c] {
				dist[r][c] = dist[row][col] + 1
				queue = append(queue, [2]int{r, c})
			}
		}
	}

	return dist
}
