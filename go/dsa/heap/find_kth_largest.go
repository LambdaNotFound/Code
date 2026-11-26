package heap

/**
 * 215. Kth Largest Element in an Array
 */
import "container/heap"

func findKthLargest(nums []int, k int) int {
	minHeap := &Heap[int]{
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
