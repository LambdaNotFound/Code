package heap

import "container/heap"

/**
 * 973. K Closest Points to Origin
 *
 */

func distance(point []int) int { return point[0]*point[0] + point[1]*point[1] }

type MaxHeap [][]int

func (h MaxHeap) Len() int            { return len(h) }
func (h MaxHeap) Less(i, j int) bool  { return distance(h[i]) > distance(h[j]) }
func (h MaxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) { *h = append(*h, x.([]int)) }
func (h *MaxHeap) Pop() interface{} {
    res := (*h)[len(*h)-1]
    *h = (*h)[:len(*h)-1]
    return res
}

func kClosest(points [][]int, k int) [][]int {
    max := MaxHeap{}
    for _, point := range points {
        heap.Push(&max, point)
        if len(max) > k {
            heap.Pop(&max)
        }
    }
    return max
}
