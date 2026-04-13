package concurrency

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ---------------------------------------------------------------------------
// LockFreeCounter tests
// ---------------------------------------------------------------------------

func TestLockFreeCounter_ZeroValue(t *testing.T) {
	var c LockFreeCounter
	assert.Equal(t, int64(0), c.Load())
}

func TestLockFreeCounter_IncrementDecrement(t *testing.T) {
	var c LockFreeCounter
	assert.Equal(t, int64(1), c.Increment())
	assert.Equal(t, int64(2), c.Increment())
	assert.Equal(t, int64(1), c.Decrement())
	assert.Equal(t, int64(1), c.Load())
}

func TestLockFreeCounter_Add(t *testing.T) {
	var c LockFreeCounter
	c.Add(10)
	c.Add(-3)
	assert.Equal(t, int64(7), c.Load())
}

func TestLockFreeCounter_Reset(t *testing.T) {
	var c LockFreeCounter
	c.Add(42)
	prev := c.Reset()
	assert.Equal(t, int64(42), prev)
	assert.Equal(t, int64(0), c.Load())
}

func TestLockFreeCounter_ConcurrentIncrement(t *testing.T) {
	var c LockFreeCounter
	const n = 10_000
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Increment()
		}()
	}
	wg.Wait()

	assert.Equal(t, int64(n), c.Load())
}

func TestLockFreeCounter_ConcurrentIncrementDecrement(t *testing.T) {
	var c LockFreeCounter
	const n = 5_000
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			c.Increment()
		}()
		go func() {
			defer wg.Done()
			c.Decrement()
		}()
	}
	wg.Wait()

	// equal increments and decrements → net zero
	assert.Equal(t, int64(0), c.Load())
}

// ---------------------------------------------------------------------------
// BoundedCounter tests
// ---------------------------------------------------------------------------

func TestBoundedCounter_IncrementAtCeiling(t *testing.T) {
	c := NewBoundedCounter(3)

	for i := int64(1); i <= 3; i++ {
		val, ok := c.Increment()
		assert.True(t, ok)
		assert.Equal(t, i, val)
	}

	// already at max — should be rejected
	val, ok := c.Increment()
	assert.False(t, ok)
	assert.Equal(t, int64(3), val)
}

func TestBoundedCounter_DecrementAtFloor(t *testing.T) {
	c := NewBoundedCounter(5)

	// counter is at 0 — decrement should fail
	val, ok := c.Decrement()
	assert.False(t, ok)
	assert.Equal(t, int64(0), val)

	c.Increment()
	val, ok = c.Decrement()
	assert.True(t, ok)
	assert.Equal(t, int64(0), val)
}

func TestBoundedCounter_ConcurrentIncrement(t *testing.T) {
	const max = 100
	c := NewBoundedCounter(max)

	var wg sync.WaitGroup
	// 500 goroutines compete; only 100 should succeed
	for i := 0; i < 500; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Increment()
		}()
	}
	wg.Wait()

	assert.Equal(t, int64(max), c.Load())
}

func TestBoundedCounter_ConcurrentDecrementBelowZero(t *testing.T) {
	const max = 100
	c := NewBoundedCounter(max)
	for i := 0; i < max; i++ {
		c.Increment()
	}

	var wg sync.WaitGroup
	// 500 goroutines compete; only 100 should succeed
	for i := 0; i < 500; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Decrement()
		}()
	}
	wg.Wait()

	assert.Equal(t, int64(0), c.Load())
}

func TestBoundedCounter_RaceDetector(t *testing.T) {
	c := NewBoundedCounter(50)
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			c.Increment()
		}()
		go func() {
			defer wg.Done()
			c.Decrement()
		}()
	}
	wg.Wait()
}
