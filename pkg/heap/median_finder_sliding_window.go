package heap

import (
	"container/heap"
	"sort"
)

/**
 * 480. Sliding Window Median
 *
 * Two Heaps + Lazy Deletion for Sliding Window Median
 *
 */
func medianSlidingWindow(nums []int, k int) []float64 {
	maxHeap := &HeapLazy{
		less: func(a, b int) bool {
			return a > b
		},
		deleted: make(map[int]int),
	}
	minHeap := &HeapLazy{
		less: func(a, b int) bool {
			return a < b
		},
		deleted: make(map[int]int),
	}

	maxSize, minSize := 0, 0
	rebalance := func() {
		// invariant: maxSize == minSize or maxSize == minSize + 1
		if maxSize > minSize+1 {
			// lower half too large → move top of maxHeap to minHeap
			maxHeap.Purge()
			val := heap.Pop(maxHeap).(int)
			maxSize--
			heap.Push(minHeap, val)
			minSize++
		} else if minSize > maxSize {
			// upper half too large → move top of minHeap to maxHeap
			minHeap.Purge()
			val := heap.Pop(minHeap).(int)
			minSize--
			heap.Push(maxHeap, val)
			maxSize++
		}
	}

	add := func(val int) {
		// route to correct half based on value
		maxHeap.Purge()
		if maxHeap.Len() == 0 || val <= maxHeap.Top() {
			heap.Push(maxHeap, val)
			maxSize++
		} else {
			heap.Push(minHeap, val)
			minSize++
		}
		rebalance()
	}

	remove := func(val int) {
		// tombstone on the heap that logically owns this value
		maxHeap.Purge()
		if val <= maxHeap.Top() {
			maxHeap.Delete(val)
			maxSize--
		} else {
			minHeap.Delete(val)
			minSize--
		}
		rebalance()
	}

	getMedian := func() float64 {
		maxHeap.Purge()
		minHeap.Purge()
		if k%2 == 1 {
			return float64(maxHeap.Top())
		}
		return float64(maxHeap.Top()+minHeap.Top()) / 2.0
	}

	// seed the first window
	for i := 0; i < k; i++ {
		add(nums[i])
	}

	result := make([]float64, 0, len(nums)-k+1)
	result = append(result, getMedian())

	// slide the window
	for i := k; i < len(nums); i++ {
		add(nums[i])
		remove(nums[i-k])
		result = append(result, getMedian())
	}

	return result
}

type HeapLazy struct {
	items   []int
	less    func(a, b int) bool
	deleted map[int]int
}

func (h *HeapLazy) Len() int           { return len(h.items) }
func (h *HeapLazy) Swap(i, j int)      { h.items[i], h.items[j] = h.items[j], h.items[i] }
func (h *HeapLazy) Less(i, j int) bool { return h.less(h.items[i], h.items[j]) }

func (h *HeapLazy) Push(item any) { h.items = append(h.items, item.(int)) }

func (h *HeapLazy) Pop() any {
	item := h.items[len(h.items)-1]
	h.items = h.items[:len(h.items)-1]
	return item
}

func (h *HeapLazy) Top() int { return h.items[0] }

// Purge(), Amortized O(log n)
func (h *HeapLazy) Purge() {
	for h.Len() > 0 && h.deleted[h.Top()] > 0 {
		h.deleted[h.Top()]--
		heap.Pop(h)
	}
}

// Delete(), O(1) — just increments a counter in the map:
func (h *HeapLazy) Delete(val int) {
	h.deleted[val]++
}

// naive
func medianSlidingWindowInsertion(nums []int, k int) []float64 {
	window := make([]int, k)
	copy(window, nums[:k])
	sort.Ints(window)

	getMedian := func() float64 {
		if k%2 == 1 {
			return float64(window[k/2])
		}
		return float64(window[k/2-1]+window[k/2]) / 2.0
	}

	result := make([]float64, 0, len(nums)-k+1)
	result = append(result, getMedian())

	for i := k; i < len(nums); i++ {
		// insert nums[i] in sorted position
		insertPos := sort.SearchInts(window, nums[i])
		window = append(window, 0)
		copy(window[insertPos+1:], window[insertPos:])
		window[insertPos] = nums[i]

		// remove nums[i-k] from sorted position
		removePos := sort.SearchInts(window, nums[i-k])
		window = append(window[:removePos], window[removePos+1:]...)

		result = append(result, getMedian())
	}

	return result
}
