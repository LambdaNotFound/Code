package heap

import (
	"container/heap"
)

/*
 * Heap implementation
 *    heap.Interface = sort.Interface + Push, Pop
 *
 *    less func(i, j T) bool
 *    ... return h.less(h.items[i], h.items[j])
 */
type Heap[T comparable] struct {
	items   []T
	less    func(a, b T) bool
	deleted map[T]int
}

// container/heap interface
func (h *Heap[T]) Len() int           { return len(h.items) }
func (h *Heap[T]) Less(i, j int) bool { return h.less(h.items[i], h.items[j]) }
func (h *Heap[T]) Swap(i, j int)      { h.items[i], h.items[j] = h.items[j], h.items[i] }

/*
 * container/heap calls Swap(0, n-1) before calling your Pop()
 * — so by the time your Pop() runs, the minimum has already been swapped to the last position.
 * Your Pop() just removes it from there.
 *
 * The heap then re-heapifies items[0] downward to restore the invariant.
 */
func (h *Heap[T]) Pop() any {
	x := h.items[len(h.items)-1]
	h.items = h.items[:len(h.items)-1]
	return x
}

func (h *Heap[T]) Push(x any) { h.items = append(h.items, x.(T)) }

func (h *Heap[T]) Peek() T { return h.items[0] }

func NewHeap[T comparable](less func(a, b T) bool) *Heap[T] {
	h := &Heap[T]{less: less, deleted: make(map[T]int)}
	heap.Init(h)
	return h
}

// Delete marks val for lazy removal. O(1).
func (h *Heap[T]) Delete(val T) {
	if h.deleted == nil {
		h.deleted = make(map[T]int)
	}
	h.deleted[val]++
}

// Purge removes deleted elements from the top of the heap. Amortized O(log n).
func (h *Heap[T]) Purge() {
	for h.Len() > 0 && h.deleted[h.Peek()] > 0 {
		h.deleted[h.Peek()]--
		heap.Pop(h)
	}
}
