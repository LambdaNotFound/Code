package dfs

/**
 * 827. Making A Large Island
 *
 * Time: Overall: O(N²)
 * Space: Space: O(N²)
 */
func largestIsland(grid [][]int) int {
	n := len(grid)
	directions := [4][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}

	inBounds := func(row, col int) bool {
		return row >= 0 && col >= 0 && row < n && col < n
	}

	// floodFill stamps every cell of one island with `id` and returns its size.
	// Overwriting the 1 with `id` doubles as the visited marker.
	var floodFill func(row, col, id int) int
	floodFill = func(row, col, id int) int {
		if !inBounds(row, col) || grid[row][col] != 1 {
			return 0
		}
		grid[row][col] = id
		size := 1
		for _, dir := range directions {
			size += floodFill(row+dir[0], col+dir[1], id)
		}
		return size
	}

	// Pass 1: label each island with a unique id (>= 2) and record its size.
	islandSize := map[int]int{}
	nextID := 2
	maxIsland := 0
	for row := range grid {
		for col := range grid[row] {
			if grid[row][col] == 1 {
				size := floodFill(row, col, nextID)
				islandSize[nextID] = size
				maxIsland = max(maxIsland, size)
				nextID++
			}
		}
	}

	// Pass 2: flipping a water cell fuses all DISTINCT islands around it.
	// `seen` dedupes an island that touches this cell from two sides.
	hasZero := false
	for row := range grid {
		for col := range grid[row] {
			if grid[row][col] != 0 {
				continue
			}
			hasZero = true

			visited := map[int]bool{}
			mergedSize := 1 // the flipped cell itself
			for _, dir := range directions {
				neiRow, neiCol := row+dir[0], col+dir[1]
				if !inBounds(neiRow, neiCol) {
					continue
				}
				neighborID := grid[neiRow][neiCol]
				if neighborID < 2 || visited[neighborID] {
					continue // water, or an island already counted
				}
				visited[neighborID] = true
				mergedSize += islandSize[neighborID]
			}
			maxIsland = max(maxIsland, mergedSize)
		}
	}

	if !hasZero {
		return n * n // all 1s; nothing to flip
	}
	return maxIsland
}
