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

// Dijkstra
func swimInWaterDijkstra(grid [][]int) int {
	n := len(grid)
	dirs := [][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}

	// minHeap: [maxElevation, row, col]
	type state struct{ maxElev, row, col int }
	h := &MinHeap{}
	heap.Push(h, state{grid[0][0], 0, 0})

	visited := make([][]bool, n)
	for i := range visited {
		visited[i] = make([]bool, n)
	}
	visited[0][0] = true

	for h.Len() > 0 {
		cur := heap.Pop(h).(state)
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
			heap.Push(h, state{nextT, nr, nc})
		}
	}
	return -1
}

type MinHeap []struct{ maxElev, row, col int }

func (h MinHeap) Len() int            { return len(h) }
func (h MinHeap) Less(i, j int) bool  { return h[i].maxElev < h[j].maxElev }
func (h MinHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MinHeap) Push(x interface{}) { *h = append(*h, x.(struct{ maxElev, row, col int })) }
func (h *MinHeap) Pop() interface{} {
	x := (*h)[len(*h)-1]
	*h = (*h)[:len(*h)-1]
	return x
}
