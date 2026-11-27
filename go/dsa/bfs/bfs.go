package bfs

/**
 * 994. Rotting Oranges
 *
 * 0 representing an empty cell,
 * 1 representing a fresh orange, or
 * 2 representing a rotten orange.
 *
 */
func orangesRotting(grid [][]int) int {
	const (
		fresh   = 1
		rotten  = 2
		visited = -1
	)

	type cell struct{ row, col int }

	rows, cols := len(grid), len(grid[0])
	queue := []cell{}

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			if grid[row][col] == rotten {
				queue = append(queue, cell{row, col})
				grid[row][col] = visited
			}
		}
	}

	dirs := [][2]int{{0, 1}, {1, 0}, {-1, 0}, {0, -1}}
	minutesElapsed := 0

	for len(queue) > 0 {
		size := len(queue)
		for i := 0; i < size; i++ {
			cur := queue[0]
			queue = queue[1:]

			for _, dir := range dirs {
				nextRow, nextCol := cur.row+dir[0], cur.col+dir[1]
				if nextRow < 0 || nextRow >= rows || nextCol < 0 || nextCol >= cols {
					continue
				}
				if grid[nextRow][nextCol] != fresh {
					continue
				}
				queue = append(queue, cell{nextRow, nextCol})
				grid[nextRow][nextCol] = visited
			}
		}
		if len(queue) > 0 {
			minutesElapsed++
		}
	}

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			if grid[row][col] == fresh {
				return -1
			}
		}
	}

	return minutesElapsed
}
