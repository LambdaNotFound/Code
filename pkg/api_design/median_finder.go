package apidesign

import "container/heap"

/*
 * 295. Find Median from Data Stream
 *
 * Min Heap + Max Heap
 */
type MedianFinder struct {
	minHeap *IntHeap
	maxHeap *IntHeap
}

func ConstructorMedianFinder() MedianFinder {
	return MedianFinder{
		minHeap: NewHeap(func(a, b int) bool {
			return a > b
		}),
		maxHeap: NewHeap(func(a, b int) bool {
			return a < b
		}),
	}
}

func (mf *MedianFinder) AddNum(num int) {
	if (mf.minHeap.Len()+mf.maxHeap.Len())%2 == 0 {
		heap.Push(mf.maxHeap, num)
		heap.Push(mf.minHeap, heap.Pop(mf.maxHeap))
	} else {
		heap.Push(mf.minHeap, num)
		heap.Push(mf.maxHeap, heap.Pop(mf.minHeap))
	}
}

func (mf *MedianFinder) FindMedian() float64 {
	if (mf.minHeap.Len()+mf.maxHeap.Len())%2 == 0 {
		return (float64(mf.minHeap.Peek()) + float64(mf.maxHeap.Peek())) / 2
	}
	return float64(mf.minHeap.Peek())
}

type IntHeap struct {
	items []int
	less  func(int, int) bool
}

func (h *IntHeap) Less(i, j int) bool { return h.less(h.items[i], h.items[j]) }
func (h *IntHeap) Swap(i, j int)      { h.items[i], h.items[j] = h.items[j], h.items[i] }
func (h *IntHeap) Len() int           { return len(h.items) }
func (h *IntHeap) Peek() int          { return h.items[0] }
func (h *IntHeap) Pop() (v interface{}) {
	h.items, v = h.items[:h.Len()-1], h.items[h.Len()-1]
	return v
}
func (h *IntHeap) Push(v interface{}) { h.items = append(h.items, v.(int)) }

func NewHeap(less func(int, int) bool) *IntHeap {
	return &IntHeap{less: less}
}
