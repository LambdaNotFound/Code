package dijkstra

import "container/heap"

/**
 * 778. Swim in Rising Water
 *
 * 1. Binary Search + DFS/BFS
 * 2. Dijkstra
 */
func swimInWater(grid [][]int) int {
	n := len(grid)

	// check if path exists using only cells with elevation <= t
	canReach := func(t int) bool {
		if grid[0][0] > t {
			return false
		}
		visited := make([][]bool, n)
		for i := range visited {
			visited[i] = make([]bool, n)
		}

		queue := [][2]int{{0, 0}}
		visited[0][0] = true
		dirs := [][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}

		for len(queue) > 0 {
			cur := queue[0]
			queue = queue[1:]
			row, col := cur[0], cur[1]

			if row == n-1 && col == n-1 {
				return true
			}

			for _, d := range dirs {
				nr, nc := row+d[0], col+d[1]
				if nr < 0 || nr >= n || nc < 0 || nc >= n {
					continue
				}
				if visited[nr][nc] || grid[nr][nc] > t {
					continue
				}
				visited[nr][nc] = true
				queue = append(queue, [2]int{nr, nc})
			}
		}
		return false
	}

	// binary search on answer t in [0, n*n-1]
	lo, hi := grid[0][0], n*n-1
	for lo < hi {
		mid := lo + (hi-lo)/2
		if canReach(mid) {
			hi = mid // try smaller t
		} else {
			lo = mid + 1 // need larger t
		}
	}
	return lo
}

func swimInWaterDFS(grid [][]int) int {
	n := len(grid)
	dirs := [][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
	visited := make(map[[2]int]bool)

	var dfs func(row, col, t int) bool
	dfs = func(row, col, t int) bool {
		if grid[row][col] > t {
			return false
		}
		if row == n-1 && col == n-1 {
			return true
		}
		for _, d := range dirs {
			r, c := row+d[0], col+d[1]
			if r < 0 || r >= n || c < 0 || c >= n {
				continue
			}
			if visited[[2]int{r, c}] {
				continue
			}
			visited[[2]int{r, c}] = true
			if dfs(r, c, t) {
				return true
			}
		}
		return false
	}

	left, right := grid[0][0], n*n-1
	for left < right {
		mid := left + (right-left)/2
		visited = make(map[[2]int]bool) // reset each iteration
		visited[[2]int{0, 0}] = true    // mark start
		if dfs(0, 0, mid) {
			right = mid
		} else {
			left = mid + 1
		}
	}
	return left
}

// swimState is the element type for MinHeap used by swimInWaterDijkstra.
type swimState struct{ maxElev, row, col int }

// Dijkstra
func swimInWaterDijkstra(grid [][]int) int {
	n := len(grid)
	dirs := [][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}

	// minHeap: [maxElevation, row, col]
	h := &Heap[swimState]{
		less: func(a, b swimState) bool {
			return a.maxElev < b.maxElev
		},
	}
	heap.Push(h, swimState{grid[0][0], 0, 0})

	visited := make([][]bool, n)
	for i := range visited {
		visited[i] = make([]bool, n)
	}
	visited[0][0] = true

	for h.Len() > 0 {
		cur := heap.Pop(h).(swimState)
		t, row, col := cur.maxElev, cur.row, cur.col

		if row == n-1 && col == n-1 {
			return t
		}

		for _, d := range dirs {
			nr, nc := row+d[0], col+d[1]
			if nr < 0 || nr >= n || nc < 0 || nc >= n {
				continue
			}
			if visited[nr][nc] {
				continue
			}
			visited[nr][nc] = true
			nextT := max(t, grid[nr][nc]) // bottleneck = max elevation
			heap.Push(h, swimState{nextT, nr, nc})
		}
	}
	return -1
}
