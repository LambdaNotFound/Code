package containers

import "cmp"

/**
 * Heap[T] — Generic Binary Heap
 *
 * A min- or max-heap over any type T, ordered by a caller-supplied comparator.
 * Backed by a slice with the standard array-heap layout:
 *
 *   items[0]                       — root (top element)
 *   items[(i-1)/2]                 — parent of items[i]
 *   items[2*i+1], items[2*i+2]    — left and right children of items[i]
 *
 * Heap property: for every node i, less(items[i], items[child]) is false —
 * i.e. the root is always the "best" element according to the comparator.
 *
 * Implementation notes:
 *   - Implemented from scratch (siftUp + siftDown) rather than wrapping
 *     container/heap. This avoids the heap.Interface double-dispatch, removes
 *     the mandatory heap.Init call, eliminates interface{} casts, and makes
 *     the algorithm directly visible in this file.
 *   - T is constrained to `any` so struct elements such as
 *     {priority int, val string} are fully supported via a custom comparator.
 *   - NewMinHeap and NewMaxHeap accept T cmp.Ordered and wire in cmp.Compare
 *     automatically. NewHeap accepts an arbitrary less func for struct types.
 *   - Heapify replaces the heap's contents and establishes the heap property
 *     in O(n) via Floyd's bottom-up algorithm. A defensive copy of the input
 *     slice is made; the caller's slice is never modified.
 *   - Pop and Peek return (zero, false) on an empty heap and never panic.
 *
 * Complexity:
 *   Push, Pop           O(log n)
 *   Peek, Size, IsEmpty O(1)
 *   Heapify             O(n)
 */

// Heap is a generic binary heap ordered by a caller-supplied comparator.
// The zero value is not usable; construct with NewMinHeap, NewMaxHeap, or NewHeap.
type Heap[T any] struct {
	items []T
	less  func(a, b T) bool
}

// NewMinHeap returns a min-heap for any ordered type.
// The smallest element (by natural ordering) is always at the top.
func NewMinHeap[T cmp.Ordered]() *Heap[T] {
	return &Heap[T]{
		items: make([]T, 0),
		less:  func(a, b T) bool { return cmp.Compare(a, b) < 0 },
	}
}

// NewMaxHeap returns a max-heap for any ordered type.
// The largest element (by natural ordering) is always at the top.
func NewMaxHeap[T cmp.Ordered]() *Heap[T] {
	return &Heap[T]{
		items: make([]T, 0),
		less:  func(a, b T) bool { return cmp.Compare(a, b) > 0 },
	}
}

// NewHeap returns a heap ordered by the supplied comparator.
// less(a, b) must return true when a should appear above b in the heap.
// Use this constructor for struct elements or non-standard ordering.
func NewHeap[T any](less func(a, b T) bool) *Heap[T] {
	return &Heap[T]{
		items: make([]T, 0),
		less:  less,
	}
}

// Push adds x to the heap.
// Time: O(log n).
func (h *Heap[T]) Push(x T) {
	h.items = append(h.items, x)
	h.siftUp(len(h.items) - 1)
}

// Pop removes and returns the top element.
// Returns (zero, false) if the heap is empty — never panics.
// Time: O(log n).
func (h *Heap[T]) Pop() (T, bool) {
	if len(h.items) == 0 {
		var zero T
		return zero, false
	}
	top := h.items[0]
	n := len(h.items) - 1
	h.items[0] = h.items[n]
	h.items = h.items[:n]
	if n > 0 {
		h.siftDown(0)
	}
	return top, true
}

// Peek returns the top element without removing it.
// Returns (zero, false) if the heap is empty — never panics.
// Time: O(1).
func (h *Heap[T]) Peek() (T, bool) {
	if len(h.items) == 0 {
		var zero T
		return zero, false
	}
	return h.items[0], true
}

// Size returns the number of elements in the heap.
// Time: O(1).
func (h *Heap[T]) Size() int {
	return len(h.items)
}

// IsEmpty reports whether the heap contains no elements.
// Time: O(1).
func (h *Heap[T]) IsEmpty() bool {
	return len(h.items) == 0
}

// Heapify replaces the heap's contents with a copy of items and establishes
// the heap property in O(n) using Floyd's bottom-up siftDown algorithm.
// The caller's slice is not modified.
// Time: O(n).
func (h *Heap[T]) Heapify(items []T) {
	h.items = make([]T, len(items))
	copy(h.items, items)
	for i := len(h.items)/2 - 1; i >= 0; i-- {
		h.siftDown(i)
	}
}

// ---------------------------------------------------------------------------
// Internal helpers
// ---------------------------------------------------------------------------

// siftUp restores the heap property after a Push by walking element i
// toward the root, swapping with its parent while it is "better".
func (h *Heap[T]) siftUp(i int) {
	for i > 0 {
		parent := (i - 1) / 2
		if h.less(h.items[i], h.items[parent]) {
			h.items[i], h.items[parent] = h.items[parent], h.items[i]
			i = parent
		} else {
			break
		}
	}
}

// siftDown restores the heap property after a Pop or during Heapify by
// walking element i toward the leaves, swapping with the "better" child.
func (h *Heap[T]) siftDown(i int) {
	n := len(h.items)
	for {
		best := i
		left := 2*i + 1
		right := 2*i + 2
		if left < n && h.less(h.items[left], h.items[best]) {
			best = left
		}
		if right < n && h.less(h.items[right], h.items[best]) {
			best = right
		}
		if best == i {
			break
		}
		h.items[i], h.items[best] = h.items[best], h.items[i]
		i = best
	}
}
