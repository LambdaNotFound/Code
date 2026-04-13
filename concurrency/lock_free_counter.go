package concurrency

import "sync/atomic"

/**
 * Lock-Free Counter
 *
 * A concurrent counter backed by a single int64 operated on exclusively
 * through CPU atomic instructions (LOCK XADD / LOCK CMPXCHG on x86).
 * No mutex, no blocking — goroutines never wait for each other.
 *
 * Two flavours are provided:
 *
 *   LockFreeCounter  — simple add/subtract, built on atomic.Int64.
 *                      Add and subtract are single instructions; no CAS loop needed.
 *
 *   BoundedCounter   — enforces a [0, max] range using a CAS loop.
 *                      Demonstrates the compare-and-swap retry pattern when
 *                      the operation is conditional (can't be expressed as plain Add).
 *
 * When to use atomic vs mutex:
 *   - Use atomics for a single integer that only needs Add/Load/Store/CAS.
 *   - Use a mutex when you need to protect a multi-field invariant or a
 *     data structure that requires more than one memory operation to update.
 */

// -----------------------------------------------------------------------------
// LockFreeCounter — unbounded, signed counter
// -----------------------------------------------------------------------------

// LockFreeCounter is a goroutine-safe integer counter.
// The zero value is ready to use.
type LockFreeCounter struct {
	value atomic.Int64
}

// Increment adds 1 to the counter and returns the new value.
func (c *LockFreeCounter) Increment() int64 {
	return c.value.Add(1)
}

// Decrement subtracts 1 from the counter and returns the new value.
func (c *LockFreeCounter) Decrement() int64 {
	return c.value.Add(-1)
}

// Add adds delta to the counter and returns the new value.
func (c *LockFreeCounter) Add(delta int64) int64 {
	return c.value.Add(delta)
}

// Reset sets the counter to zero and returns the previous value.
func (c *LockFreeCounter) Reset() int64 {
	return c.value.Swap(0)
}

// Load returns the current value without modifying it.
func (c *LockFreeCounter) Load() int64 {
	return c.value.Load()
}

// -----------------------------------------------------------------------------
// BoundedCounter — counter clamped to [0, max] via CAS loop
//
// Plain atomic.Add cannot enforce a ceiling — by the time we check the value,
// another goroutine may have already added past the limit. We need a
// read-modify-write that is atomic as a whole: load the current value, decide
// whether to proceed, then CAS the new value in one shot.
// -----------------------------------------------------------------------------

// BoundedCounter is a goroutine-safe counter clamped to [0, max].
type BoundedCounter struct {
	value atomic.Int64
	max   int64
}

// NewBoundedCounter returns a counter that never exceeds max.
func NewBoundedCounter(max int64) *BoundedCounter {
	return &BoundedCounter{max: max}
}

// Increment adds 1 if the counter is below max.
// Returns (newValue, true) on success, or (currentValue, false) if already at max.
func (c *BoundedCounter) Increment() (int64, bool) {
	for {
		cur := c.value.Load()
		if cur >= c.max {
			return cur, false
		}
		next := cur + 1
		if c.value.CompareAndSwap(cur, next) {
			return next, true
		}
		// Another goroutine changed the value between Load and CAS — retry.
	}
}

// Decrement subtracts 1 if the counter is above 0.
// Returns (newValue, true) on success, or (currentValue, false) if already at 0.
func (c *BoundedCounter) Decrement() (int64, bool) {
	for {
		cur := c.value.Load()
		if cur <= 0 {
			return cur, false
		}
		next := cur - 1
		if c.value.CompareAndSwap(cur, next) {
			return next, true
		}
	}
}

// Load returns the current value.
func (c *BoundedCounter) Load() int64 {
	return c.value.Load()
}
