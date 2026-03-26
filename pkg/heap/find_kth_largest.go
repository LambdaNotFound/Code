package heap

/**
 * 215. Kth Largest Element in an Array
 */
import "container/heap"

func findKthLargest(nums []int, k int) int {
	minHeap := &IntHeap{
		less: func(i, j int) bool {
			return i < j
		},
	}

	for _, num := range nums {
		heap.Push(minHeap, num)
		if minHeap.Len() > k {
			heap.Pop(minHeap)
		}
	}

	return heap.Pop(minHeap).(int)
}

type IntHeap struct {
	items []int
	less  func(int, int) bool
}

func (h *IntHeap) Len() int           { return len(h.items) }
func (h *IntHeap) Swap(i, j int)      { h.items[i], h.items[j] = h.items[j], h.items[i] }
func (h *IntHeap) Less(i, j int) bool { return h.less(h.items[i], h.items[j]) }

func (h *IntHeap) Push(i interface{}) {
	h.items = append(h.items, i.(int))
}
func (h *IntHeap) Pop() interface{} {
	item := h.items[len(h.items)-1]
	h.items = h.items[:len(h.items)-1]
	return item
}
