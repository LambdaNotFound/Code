package heap

import (
	. "gocode/types"
)

// heap.Interface = sort.Interface + Push, Pop
type Heap struct {
	items []int
	less  func(int, int) bool
}

func (h *Heap) Less(i, j int) bool { return h.less(h.items[i], h.items[j]) }
func (h *Heap) Swap(i, j int)      { h.items[i], h.items[j] = h.items[j], h.items[i] }
func (h *Heap) Len() int           { return len(h.items) }
func (h *Heap) Pop() (v interface{}) {
	h.items, v = h.items[:h.Len()-1], h.items[h.Len()-1]
	return v
}
func (h *Heap) Push(v interface{}) { h.items = append(h.items, v.(int)) }
func (h *Heap) Peek() int          { return h.items[0] }

func NewHeap(less func(int, int) bool) *Heap {
	return &Heap{less: less}
}

// MinHeap struct
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

type MaxHeap []int

func (h MaxHeap) Len() int           { return len(h) }
func (h MaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h MaxHeap) Less(i, j int) bool { return h[i] > h[j] }

func (h *MaxHeap) Push(n interface{}) {
	*h = append(*h, n.(int))
}

func (h *MaxHeap) Pop() interface{} {
	top := (*h)[len(*h)-1]
	*h = (*h)[:len(*h)-1]
	return top
}

type MinHeap []int

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h MinHeap) Less(i, j int) bool { return h[i] < h[j] }

func (h *MinHeap) Push(n interface{}) {
	*h = append(*h, n.(int))
}

func (h *MinHeap) Pop() interface{} {
	top := (*h)[len(*h)-1]
	*h = (*h)[:len(*h)-1]
	return top
}
