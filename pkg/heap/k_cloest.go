package heap

import "container/heap"

/**
 * 973. K Closest Points to Origin
 *
 */

func distance(point []int) int { return point[0]*point[0] + point[1]*point[1] }

type PointMaxHeap [][]int

func (h PointMaxHeap) Len() int            { return len(h) }
func (h PointMaxHeap) Less(i, j int) bool  { return distance(h[i]) > distance(h[j]) }
func (h PointMaxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *PointMaxHeap) Push(x interface{}) { *h = append(*h, x.([]int)) }
func (h *PointMaxHeap) Pop() interface{} {
	res := (*h)[len(*h)-1]
	*h = (*h)[:len(*h)-1]
	return res
}

func kClosest(points [][]int, k int) [][]int {
	max := PointMaxHeap{}
	for _, point := range points {
		heap.Push(&max, point)
		if len(max) > k {
			heap.Pop(&max)
		}
	}
	return max
}
