package interview

import (
	"container/heap"
)

type Cell struct {
	height, x, y int
}

type MinHeap []Cell

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h MinHeap) Less(i, j int) bool { return h[i].height < h[j].height }

func (h *MinHeap) Pop() any {
	item := (*h)[len(*h)-1]
	*h = (*h)[:len(*h)-1]
	return item
}

func (h *MinHeap) Push(item any) {
	*h = append(*h, item.(Cell))
}

/**
 * 407. Trapping Rain Water II
 */
func trapRainWater(heightMap [][]int) int {
	m, n := len(heightMap), len(heightMap[0])

	minHeap := &MinHeap{}
	visited := map[[2]int]bool{}
	for i := 0; i < m; i++ {
		heap.Push(minHeap, Cell{heightMap[i][0], i, 0})
		heap.Push(minHeap, Cell{heightMap[i][n-1], i, n - 1})
		visited[[2]int{i, 0}] = true
		visited[[2]int{i, n - 1}] = true
	}
	for j := 0; j < n; j++ {
		heap.Push(minHeap, Cell{heightMap[0][j], 0, j})
		heap.Push(minHeap, Cell{heightMap[m-1][j], m - 1, j})
		visited[[2]int{0, j}] = true
		visited[[2]int{m - 1, j}] = true
	}

	result := 0
	dirs := [][]int{{0, 1}, {0, -1}, {-1, 0}, {1, 0}}
	for len(*minHeap) > 0 {
		curr := heap.Pop(minHeap).(Cell)

		for _, dir := range dirs {
			nx, ny := curr.x+dir[0], curr.y+dir[1]
			if nx < 0 || ny < 0 || nx >= m || ny >= n {
				continue
			}
			if visited[[2]int{nx, ny}] == true {
				continue
			}

			height := heightMap[nx][ny]
			result += max(0, curr.height-height)
			heap.Push(minHeap, Cell{max(curr.height, height), nx, ny})
			visited[[2]int{nx, ny}] = true
		}
	}
	return result
}
