package containers

import (
	"container/heap"
)

// Heap struct
type Heap[T comparable] struct {
    items []T
    less  func(a, b T) bool
}

// Implement sort.Interface
func (h *Heap[T]) Len() int           { return len(h.items) }
func (h *Heap[T]) Less(i, j int) bool { return h.less(h.items[i], h.items[j]) }
func (h *Heap[T]) Swap(i, j int)      { h.items[i], h.items[j] = h.items[j], h.items[i] }

// Implement heap.Interface methods
func (h *Heap[T]) Push(x interface{}) {
    h.items = append(h.items, x.(T))
}

func (h *Heap[T]) Pop() interface{} {
    n := len(h.items)
    x := h.items[n-1]
    h.items = h.items[:n-1]
    return x
}

// Public API
func (h *Heap[T]) PushItem(x T) {
    heap.Push(h, x)
}

func (h *Heap[T]) PopItem() T {
    return heap.Pop(h).(T)
}

func (h *Heap[T]) Peek() T {
    return h.items[0]
}
