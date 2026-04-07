package heap

import (
	. "gocode/types"
)

/*
 * Heap implementation
 *    heap.Interface = sort.Interface + Push, Pop
 *
 *    less func(T, T) bool
 *    ... return h.less(h.items[i], h.items[j])
 */
type Heap struct {
	items []int // items[0] is the MAX or MIN
	less  func(int, int) bool
}

func (h *Heap) Less(i, j int) bool { return h.less(h.items[i], h.items[j]) }
func (h *Heap) Swap(i, j int)      { h.items[i], h.items[j] = h.items[j], h.items[i] }
func (h *Heap) Len() int           { return len(h.items) }

/*
 * container/heap calls Swap(0, n-1) before calling your Pop()
 * — so by the time your Pop() runs, the minimum has already been swapped to the last position.
 * Your Pop() just removes it from there.
 *
 * The heap then re-heapifies items[0] downward to restore the invariant.
 */
func (h *Heap) Pop() interface{} {
	v := h.items[h.Len()-1]
	h.items = h.items[:h.Len()-1]
	return v
}
func (h *Heap) Push(v interface{}) { h.items = append(h.items, v.(int)) }

func (h *Heap) Peek() int { return h.items[0] }

func NewHeap(less func(int, int) bool) *Heap {
	return &Heap{less: less}
}

// Heap w/ ListNode
type ListNodeMinHeap []*ListNode

// Implement heap.Interface methods
// sort.Interface: Len(), Less() and Swap()
func (h ListNodeMinHeap) Len() int { return len(h) }

// Less(i, j) answers: “Should i be closer to the root than j?”
// Min heap: h[i] < h[j] → i stays higher
// Max heap: h[i] > h[j] → i stays higher
func (h ListNodeMinHeap) Less(i, j int) bool { return h[i].Val < h[j].Val } // if true, move to last
func (h ListNodeMinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }    // swap

// Push, Pop
// With a pointer receiver (h *MinHeap), the method can change the underlying value.
// With a value receiver (h MinHeap), it cannot.
// When you need to change the slice itself (append, reassign, reslice),
// pass *[]T so the caller gets the new header.
func (h *ListNodeMinHeap) Push(x interface{}) {
	*h = append(*h, x.(*ListNode)) // append expects a slice, not a pointer
}

func (h *ListNodeMinHeap) Pop() interface{} { // Slice header is passed by value
	top := (*h)[len(*h)-1] // Remove last element
	*h = (*h)[:len(*h)-1]
	return top
}
