package heap

import "container/heap"

/**
 * 295. Find Median from Data Stream
 */
type MedianFinder struct {
	largeHeap *Heap // MinHeap, intended to hold the larger half
	smallHeap *Heap // MaxHeap, intended to hold the smaller half
}

func Constructor() MedianFinder {
	return MedianFinder{
		largeHeap: NewHeap(func(a, b int) bool {
			return a < b
		}),
		smallHeap: NewHeap(func(a, b int) bool {
			return a > b
		}),
	}
}

// AddNum: Time: O(log n), Space: O(n)
func (this *MedianFinder) AddNum(num int) {
	sizeLargeHeap, sizeSmallHeap := this.largeHeap.Len(), this.smallHeap.Len()
	if (sizeLargeHeap+sizeSmallHeap)%2 == 0 { // Odd/Even case logic can be swapped
		heap.Push(this.largeHeap, num)
		item := heap.Pop(this.largeHeap)
		heap.Push(this.smallHeap, item)
	} else {
		heap.Push(this.smallHeap, num)
		item := heap.Pop(this.smallHeap)
		heap.Push(this.largeHeap, item)
	}
}

// FindMedian: O(1)
func (this *MedianFinder) FindMedian() float64 {
	total := this.largeHeap.Len() + this.smallHeap.Len()
	if total%2 == 0 {
		return (float64(this.smallHeap.Peek()) + float64(this.largeHeap.Peek())) / 2
	}
	return float64(this.smallHeap.Peek())
}

/** to support delete number

type MedianFinder struct {
    maxHeap  *Heap[int]
    minHeap  *Heap[int]
    deleted  map[int]int  // num → count of pending deletions
    maxCount int          // effective size of maxHeap
    minCount int          // effective size of minHeap
}

func (mf *MedianFinder) Delete(num int) {
    mf.deleted[num]++

    // determine which heap it logically belongs to and decrement count
    if num <= mf.maxHeap.items[0] {
        mf.maxCount--
    } else {
        mf.minCount--
    }

    mf.rebalance()
}

// prune removes logically deleted tops from heaps
func (mf *MedianFinder) prune(h *Heap[int]) {
    for h.Len() > 0 && mf.deleted[h.items[0]] > 0 {
        mf.deleted[h.items[0]]--
        heap.Pop(h)
    }
}

func (mf *MedianFinder) rebalance() {
    // prune stale tops before rebalancing
    mf.prune(mf.maxHeap)
    mf.prune(mf.minHeap)

    if mf.maxCount > mf.minCount+1 {
        mf.minCount++
        mf.maxCount--
        heap.Push(mf.minHeap, heap.Pop(mf.maxHeap))
        mf.prune(mf.maxHeap)
    } else if mf.minCount > mf.maxCount {
        mf.maxCount++
        mf.minCount--
        heap.Push(mf.maxHeap, heap.Pop(mf.minHeap))
        mf.prune(mf.minHeap)
    }
}
*/
