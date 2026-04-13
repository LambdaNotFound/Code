package concurrency

import "sync/atomic"

/**
 * Lock-Free Stack — Treiber Stack (1986)
 *
 * A concurrent LIFO stack that never blocks. Goroutines spin and retry
 * using Compare-And-Swap (CAS) instead of acquiring a mutex.
 *
 * Structure: a singly-linked list with a single top pointer.
 *
 *   top → [C] → [B] → [A] → nil
 *
 * Push algorithm:
 *   1. Create a new node.
 *   2. Set new.next = top  (snapshot the current top).
 *   3. CAS top: nil→new or old→new.
 *      If another goroutine pushed between steps 2 and 3, the CAS fails
 *      and we reload top and retry.
 *
 * Pop algorithm:
 *   1. Load top. If nil, stack is empty.
 *   2. CAS top: old→old.next.
 *      If another goroutine popped between steps 1 and 2, the CAS fails
 *      and we reload and retry.
 *
 * Why no ABA problem here?
 *   ABA happens when a pointer changes A→B→A between a Load and a CAS,
 *   making the CAS succeed on a stale value. Go's GC prevents this: as long
 *   as we hold a reference to a node, its memory cannot be reclaimed and
 *   reused at the same address. No hazard pointers or epoch counters needed.
 *
 * Comparison with Michael-Scott Queue:
 *   The Treiber Stack is simpler — one atomic pointer vs. two (head+tail).
 *   The queue needs a tail pointer to allow O(1) enqueue without traversal.
 *   The stack always operates at the top, so one pointer suffices.
 */

// stackNode is a single element in the linked list.
type stackNode[T any] struct {
	val  T
	next atomic.Pointer[stackNode[T]]
}

// LockFreeStack is a generic, goroutine-safe LIFO stack.
// The zero value is not usable; use NewLockFreeStack.
type LockFreeStack[T any] struct {
	top  atomic.Pointer[stackNode[T]]
	size atomic.Int64
}

// NewLockFreeStack returns an empty, ready-to-use LockFreeStack.
func NewLockFreeStack[T any]() *LockFreeStack[T] {
	return &LockFreeStack[T]{}
}

// Push adds val to the top of the stack.
func (s *LockFreeStack[T]) Push(val T) {
	newNode := &stackNode[T]{val: val}

	for {
		top := s.top.Load()
		// Point new node at the current top before attempting the swing.
		newNode.next.Store(top)

		if s.top.CompareAndSwap(top, newNode) {
			s.size.Add(1)
			return
		}
		// Another goroutine pushed first — reload and retry.
	}
}

// Pop removes and returns the top value of the stack.
// Returns (value, true) on success, or (zero, false) if the stack is empty.
func (s *LockFreeStack[T]) Pop() (T, bool) {
	for {
		top := s.top.Load()
		if top == nil {
			var zero T
			return zero, false
		}

		next := top.next.Load()
		if s.top.CompareAndSwap(top, next) {
			s.size.Add(-1)
			return top.val, true
		}
		// Another goroutine popped first — reload and retry.
	}
}

// Peek returns the top value without removing it.
// Returns (value, true) on success, or (zero, false) if the stack is empty.
// Note: the returned value may be stale by the time the caller uses it.
func (s *LockFreeStack[T]) Peek() (T, bool) {
	top := s.top.Load()
	if top == nil {
		var zero T
		return zero, false
	}
	return top.val, true
}

// Size returns the number of elements currently in the stack.
func (s *LockFreeStack[T]) Size() int {
	return int(s.size.Load())
}

// IsEmpty reports whether the stack has no elements.
func (s *LockFreeStack[T]) IsEmpty() bool {
	return s.top.Load() == nil
}
