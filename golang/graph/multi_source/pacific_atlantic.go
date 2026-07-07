package multisource

/**
 * 417. Pacific Atlantic Water Flow
 *
 * multi-source BFS, DFS
 */
func pacificAtlanticDFS(heights [][]int) [][]int {
	if len(heights) == 0 || len(heights[0]) == 0 {
		return [][]int{}
	}

	numRows, numCols := len(heights), len(heights[0])

	pacificReachable := make([][]bool, numRows)
	atlanticReachable := make([][]bool, numRows)
	for row := 0; row < numRows; row++ {
		pacificReachable[row] = make([]bool, numCols)
		atlanticReachable[row] = make([]bool, numCols)
	}

	directions := [][]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}

	// We pass reachable matrix by reference so all DFS calls share the same visited state (multi-source)
	var dfs func(row, col int, reachable [][]bool)
	dfs = func(row, col int, reachable [][]bool) {
		reachable[row][col] = true // Base case: already visited this cell in this ocean's traversal

		for _, dir := range directions {
			nextRow, nextCol := row+dir[0], col+dir[1]

			// Skip if: out of bounds, already visited, or lower height (can't flow uphill)
			isOutOfBounds := nextRow < 0 || nextRow >= numRows || nextCol < 0 || nextCol >= numCols
			if isOutOfBounds {
				continue
			}

			isAlreadyVisited := reachable[nextRow][nextCol]
			isLowerHeight := heights[nextRow][nextCol] < heights[row][col]

			if isAlreadyVisited || isLowerHeight {
				continue
			}

			dfs(nextRow, nextCol, reachable)
		}
	}

	// Seed Pacific DFS: top row + left column
	for col := 0; col < numCols; col++ {
		dfs(0, col, pacificReachable) // top row
	}
	for row := 0; row < numRows; row++ {
		dfs(row, 0, pacificReachable) // left column
	}

	// Seed Atlantic DFS: bottom row + right column
	for col := 0; col < numCols; col++ {
		dfs(numRows-1, col, atlanticReachable) // bottom row
	}
	for row := 0; row < numRows; row++ {
		dfs(row, numCols-1, atlanticReachable) // right column
	}

	// Collect cells reachable by BOTH oceans
	result := [][]int{}
	for row := 0; row < numRows; row++ {
		for col := 0; col < numCols; col++ {
			if pacificReachable[row][col] && atlanticReachable[row][col] {
				result = append(result, []int{row, col})
			}
		}
	}

	return result
}

func pacificAtlanticBFS(heights [][]int) [][]int {
	if len(heights) == 0 || len(heights[0]) == 0 {
		return [][]int{}
	}

	numRows, numCols := len(heights), len(heights[0])

	pacificReachable := make([][]bool, numRows)
	atlanticReachable := make([][]bool, numRows)
	for row := 0; row < numRows; row++ {
		pacificReachable[row] = make([]bool, numCols)
		atlanticReachable[row] = make([]bool, numCols)
	}

	directions := [][]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}

	// BFS helper: expands UPHILL from ocean borders
	// queue is pre-seeded with all border cells (multi-source BFS)
	bfs := func(queue [][2]int, reachable [][]bool) {
		for len(queue) > 0 {
			// Dequeue front cell
			curr := queue[0]
			queue = queue[1:]

			for _, dir := range directions {
				nextRow, nextCol := curr[0]+dir[0], curr[1]+dir[1]

				// Skip if: out of bounds, already visited, or lower height
				if nextRow < 0 || nextRow >= numRows || nextCol < 0 || nextCol >= numCols {
					continue
				}
				if reachable[nextRow][nextCol] {
					continue
				}
				if heights[nextRow][nextCol] < heights[curr[0]][curr[1]] {
					continue
				}

				reachable[nextRow][nextCol] = true
				queue = append(queue, [2]int{nextRow, nextCol})
			}
		}
	}

	// Seed Pacific queue: top row + left column
	pacificQueue := [][2]int{}
	for col := 0; col < numCols; col++ {
		pacificQueue = append(pacificQueue, [2]int{0, col})
		pacificReachable[0][col] = true // mark on seed, not on dequeue
	}
	for row := 1; row < numRows; row++ { // start at 1 to avoid double-marking (0,0)
		pacificQueue = append(pacificQueue, [2]int{row, 0})
		pacificReachable[row][0] = true
	}

	// Seed Atlantic queue: bottom row + right column
	atlanticQueue := [][2]int{}
	for col := 0; col < numCols; col++ {
		atlanticQueue = append(atlanticQueue, [2]int{numRows - 1, col})
		atlanticReachable[numRows-1][col] = true
	}
	for row := 0; row < numRows-1; row++ { // stop at numRows-1 to avoid double-marking corner
		atlanticQueue = append(atlanticQueue, [2]int{row, numCols - 1})
		atlanticReachable[row][numCols-1] = true
	}

	bfs(pacificQueue, pacificReachable)
	bfs(atlanticQueue, atlanticReachable)

	// Collect cells reachable by BOTH oceans
	result := [][]int{}
	for row := 0; row < numRows; row++ {
		for col := 0; col < numCols; col++ {
			if pacificReachable[row][col] && atlanticReachable[row][col] {
				result = append(result, []int{row, col})
			}
		}
	}

	return result
}
