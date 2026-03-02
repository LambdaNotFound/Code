package heap

import "container/heap"

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
	if (sizeLargeHeap+sizeSmallHeap)%2 == 0 {
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
