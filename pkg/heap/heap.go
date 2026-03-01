package heap

import (
	. "gocode/types"
)

// heap.Interface = sort.Interface + Push, Pop

// MinHeap struct
type MinHeap []*ListNode

// Implement heap.Interface methods
// sort.Interface: Len(), Less() and Swap()
func (h MinHeap) Len() int { return len(h) }

// Less(i, j) answers: “Should i be closer to the root than j?”
// Min heap: h[i] < h[j] → i stays higher
// Max heap: h[i] > h[j] → i stays higher
func (h MinHeap) Less(i, j int) bool { return h[i].Val < h[j].Val } // if true, move to last
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }    // swap

// Push, Pop
// With a pointer receiver (h *MinHeap), the method can change the underlying value.
// With a value receiver (h MinHeap), it cannot.
// When you need to change the slice itself (append, reassign, reslice),
// pass *[]T so the caller gets the new header.
func (h *MinHeap) Push(x interface{}) {
	*h = append(*h, x.(*ListNode)) // append expects a slice, not a pointer
}

func (h *MinHeap) Pop() interface{} { // Slice header is passed by value
	top := (*h)[len(*h)-1] // Remove last element
	*h = (*h)[:len(*h)-1]
	return top
}
