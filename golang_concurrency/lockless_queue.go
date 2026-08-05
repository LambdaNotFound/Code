package concurrency

import (
	"sync/atomic"
)

/**
 * Lock-Free Queue — Michael-Scott Algorithm (1996)
 *
 * A concurrent FIFO queue that never blocks. Goroutines spin and retry
 * using Compare-And-Swap (CAS) instead of acquiring a mutex.
 *
 * Structure: a singly-linked list with a sentinel (dummy) head node.
 *
 *   head → [dummy] → [A] → [B] → [C] ← tail
 *
 * Invariants:
 *   - head always points to the dummy node (the node before the first real value).
 *   - tail always points to the last node OR the second-to-last node
 *     (it can lag by one when a goroutine is preempted mid-enqueue).
 *   - Any goroutine that observes a lagging tail advances it before proceeding.
 *     This "helping" pattern keeps all goroutines making progress.
 *
 * Why no ABA problem here?
 *   ABA happens when a pointer changes A→B→A between a Load and a CAS,
 *   making the CAS succeed on a stale value. Go's GC prevents this: as long
 *   as we hold a reference to a node, its memory cannot be reclaimed and
 *   reused at the same address. No hazard pointers or epoch counters needed.
 */

// node is a single element in the linked list.
type node[T any] struct {
	val  T
	next atomic.Pointer[node[T]]
}

// LockFreeQueue is a generic, goroutine-safe FIFO queue.
// The zero value is not usable; use NewLockFreeQueue.
type LockFreeQueue[T any] struct {
	head atomic.Pointer[node[T]] // always points to the dummy sentinel
	tail atomic.Pointer[node[T]] // points to last node (may lag by one)
	size atomic.Int64
}

// NewLockFreeQueue returns an empty, ready-to-use LockFreeQueue.
// It seeds the list with a sentinel node so that head and tail
// always point to a valid node — simplifying the enqueue/dequeue logic.
func NewLockFreeQueue[T any]() *LockFreeQueue[T] {
	q := &LockFreeQueue[T]{}
	sentinel := &node[T]{}
	q.head.Store(sentinel)
	q.tail.Store(sentinel)
	return q
}

// Enqueue appends val to the back of the queue.
//
// Algorithm (two steps):
//  1. Link the new node onto tail.next via CAS.
//     If another goroutine wins the CAS first, reload tail and retry.
//  2. Swing q.tail forward to the new node via CAS.
//     This CAS may fail (another goroutine already did it), which is fine —
//     we still return successfully.
func (q *LockFreeQueue[T]) Enqueue(val T) {
	newNode := &node[T]{val: val}

	for {
		tail := q.tail.Load()
		next := tail.next.Load()

		if next == nil {
			// tail.next is empty — try to attach the new node.
			if tail.next.CompareAndSwap(nil, newNode) {
				// Linked successfully. Try to advance tail.
				// It is fine if this CAS fails: the next Enqueue (or Dequeue)
				// will observe the lag and advance tail itself.
				q.tail.CompareAndSwap(tail, newNode)
				q.size.Add(1)
				return
			}
			// Another goroutine linked a node first — retry.
		} else {
			// tail is lagging (a goroutine linked a node but hasn't swung tail yet).
			// Help by advancing tail, then retry.
			q.tail.CompareAndSwap(tail, next)
		}
	}
}

// Dequeue removes and returns the front value of the queue.
// Returns (value, true) on success, or (zero, false) if the queue is empty.
//
// Algorithm:
//  1. Read head (the sentinel) and head.next (the first real node).
//  2. If head.next is nil the queue is empty.
//  3. If tail == head, tail is lagging — advance it and retry.
//  4. Read the value from head.next, then CAS head forward to head.next.
//     head.next becomes the new sentinel; the old sentinel is unreachable.
func (q *LockFreeQueue[T]) Dequeue() (T, bool) {
	for {
		head := q.head.Load()
		tail := q.tail.Load()
		next := head.next.Load()

		if head == tail {
			if next == nil {
				// Queue is genuinely empty.
				var zero T
				return zero, false
			}
			// tail is lagging behind head — advance tail, then retry.
			q.tail.CompareAndSwap(tail, next)
			continue
		}

		// Read the value before the CAS: once we swing head, another goroutine
		// may enqueue onto the node we just claimed, overwriting next.val.
		val := next.val
		if q.head.CompareAndSwap(head, next) {
			// next is now the new sentinel; its val field is no longer meaningful.
			q.size.Add(-1)
			return val, true
		}
		// Another goroutine dequeued first — reload and retry.
	}
}

// Size returns the number of elements currently in the queue.
// Because enqueuers increment the counter after a successful CAS and
// dequeuers decrement it after theirs, the count is always consistent
// with the actual list length.
func (q *LockFreeQueue[T]) Size() int {
	return int(q.size.Load())
}

// IsEmpty reports whether the queue has no elements.
func (q *LockFreeQueue[T]) IsEmpty() bool {
	return q.head.Load().next.Load() == nil
}
