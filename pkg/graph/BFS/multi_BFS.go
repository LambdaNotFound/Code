package bfs

import "math"

/**
 * 542. 01 Matrix
 * a multi-source BFS approach
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
	queue := make([][]int, 0)
	for row := 0; row < m; row++ {
		for col := 0; col < n; col++ {
			if mat[row][col] == 0 {
				queue = append(queue, []int{row, col}) // seed
				dist[row][col] = 0
			} else {
				dist[row][col] = math.MaxInt // unvisited sentinel
			}
		}
	}

	// BFS outward from all zero cells
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		for _, dir := range dirs {
			r, c := cur[0]+dir[0], cur[1]+dir[1]
			if r < 0 || r >= m || c < 0 || c >= n {
				continue
			}

			// only update if we found a shorter distance
			if dist[cur[0]][cur[1]]+1 < dist[r][c] {
				dist[r][c] = dist[cur[0]][cur[1]] + 1
				queue = append(queue, []int{r, c})
			}
		}
	}

	return dist
}
