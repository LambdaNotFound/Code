package concurrency

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// ---------------------------------------------------------------------------
// Allow: burst handling
// ---------------------------------------------------------------------------

func TestLeakyBucket_Allow_BurstAboveCapacityRejected(t *testing.T) {
	tests := []struct {
		name         string
		capacity     int
		attempts     int
		wantAdmitted int
	}{
		{name: "exactly at capacity", capacity: 3, attempts: 3, wantAdmitted: 3},
		{name: "burst above capacity", capacity: 3, attempts: 10, wantAdmitted: 3},
		{name: "single slot bucket", capacity: 1, attempts: 5, wantAdmitted: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// leakInterval long enough that no tick fires during the test
			lb := NewLeakyBucket(tt.capacity, time.Hour)
			defer lb.Close()

			admitted := 0
			for i := 0; i < tt.attempts; i++ {
				if lb.Allow() {
					admitted++
				}
			}

			assert.Equal(t, tt.wantAdmitted, admitted)
		})
	}
}

func TestLeakyBucket_Allow_AdmitsAgainAfterLeak(t *testing.T) {
	const leakInterval = 30 * time.Millisecond
	lb := NewLeakyBucket(1, leakInterval)
	defer lb.Close()

	assert.True(t, lb.Allow(), "first request should fill the only slot")
	assert.False(t, lb.Allow(), "bucket is full, second request should be rejected")

	assert.Eventually(t, lb.Allow, time.Second, 5*time.Millisecond,
		"a slot should free up once the leak interval elapses")
}

// ---------------------------------------------------------------------------
// Steady-state rate
// ---------------------------------------------------------------------------

func TestLeakyBucket_SteadyStateRateRespectedOverTime(t *testing.T) {
	const (
		capacity     = 1
		leakInterval = 30 * time.Millisecond
		window       = 300 * time.Millisecond
	)
	lb := NewLeakyBucket(capacity, leakInterval)
	defer lb.Close()

	assert.True(t, lb.Allow(), "consume the initial burst allowance")

	admitted := 0
	deadline := time.Now().Add(window)
	for time.Now().Before(deadline) {
		if lb.Allow() {
			admitted++
		}
		time.Sleep(5 * time.Millisecond)
	}

	wantLeaks := int(window / leakInterval)
	assert.InDelta(t, wantLeaks, admitted, 3,
		"admitted count should track window/leakInterval within scheduler jitter")
}

// ---------------------------------------------------------------------------
// Concurrency safety
// ---------------------------------------------------------------------------

func TestLeakyBucket_ConcurrentAllowRespectsBound(t *testing.T) {
	const (
		capacity     = 5
		leakInterval = 10 * time.Millisecond
		numCallers   = 200
	)
	lb := NewLeakyBucket(capacity, leakInterval)
	defer lb.Close()

	var admitted atomic.Int64
	var wg sync.WaitGroup
	start := time.Now()
	for i := 0; i < numCallers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if lb.Allow() {
				admitted.Add(1)
			}
		}()
	}
	wg.Wait()
	elapsed := time.Since(start)

	// upper bound: initial capacity plus every leak that could have fired
	// while the goroutines were racing, with one tick of slack
	maxAdmitted := int64(capacity) + int64(elapsed/leakInterval) + 1
	assert.LessOrEqual(t, admitted.Load(), maxAdmitted)
	assert.GreaterOrEqual(t, admitted.Load(), int64(capacity),
		"at least the initial burst capacity should have been admitted")
}

func TestLeakyBucket_ConcurrentWaitAllSucceedEventually(t *testing.T) {
	const (
		capacity     = 2
		leakInterval = 5 * time.Millisecond
		numCallers   = 20
	)
	lb := NewLeakyBucket(capacity, leakInterval)
	defer lb.Close()

	var wg sync.WaitGroup
	errs := make([]error, numCallers)
	for i := 0; i < numCallers; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			errs[idx] = lb.Wait(ctx)
		}(i)
	}
	wg.Wait()

	for i, err := range errs {
		assert.NoErrorf(t, err, "caller %d should have been admitted before its context expired", i)
	}
}

// ---------------------------------------------------------------------------
// Wait: blocking and cancellation
// ---------------------------------------------------------------------------

func TestLeakyBucket_Wait_BlocksUntilSlotLeaks(t *testing.T) {
	const leakInterval = 40 * time.Millisecond
	lb := NewLeakyBucket(1, leakInterval)
	defer lb.Close()

	assert.True(t, lb.Allow(), "fill the only slot")

	done := make(chan error, 1)
	start := time.Now()
	go func() {
		done <- lb.Wait(context.Background())
	}()

	select {
	case err := <-done:
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, time.Since(start), leakInterval/2,
			"Wait should not return before a slot actually leaked")
	case <-time.After(time.Second):
		t.Fatal("Wait should have unblocked once a slot leaked")
	}
}

func TestLeakyBucket_Wait_ReturnsContextErrorOnCancellation(t *testing.T) {
	// leakInterval long enough that no slot frees up before the context expires
	lb := NewLeakyBucket(1, time.Hour)
	defer lb.Close()

	assert.True(t, lb.Allow(), "fill the bucket")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	err := lb.Wait(ctx)

	assert.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestLeakyBucket_Wait_ReturnsErrClosedOnClose(t *testing.T) {
	lb := NewLeakyBucket(1, time.Hour)
	assert.True(t, lb.Allow(), "fill the bucket")

	done := make(chan error, 1)
	go func() {
		done <- lb.Wait(context.Background())
	}()

	// give the goroutine a moment to actually reach the blocking send
	time.Sleep(20 * time.Millisecond)
	lb.Close()

	select {
	case err := <-done:
		assert.ErrorIs(t, err, ErrClosed)
	case <-time.After(time.Second):
		t.Fatal("Wait should have returned ErrClosed once the bucket was closed")
	}
}

// ---------------------------------------------------------------------------
// Close
// ---------------------------------------------------------------------------

func TestLeakyBucket_Close_IsIdempotent(t *testing.T) {
	lb := NewLeakyBucket(1, time.Millisecond)

	assert.NotPanics(t, func() {
		lb.Close()
		lb.Close()
		lb.Close()
	})
}

func TestLeakyBucket_Close_IsSafeForConcurrentCallers(t *testing.T) {
	lb := NewLeakyBucket(1, time.Millisecond)
	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			lb.Close()
		}()
	}

	assert.NotPanics(t, func() { wg.Wait() })
}

// ---------------------------------------------------------------------------
// Constructor validation
// ---------------------------------------------------------------------------

func TestNewLeakyBucket_PanicsOnNonPositiveCapacity(t *testing.T) {
	tests := []struct {
		name     string
		capacity int
	}{
		{name: "zero capacity", capacity: 0},
		{name: "negative capacity", capacity: -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Panics(t, func() {
				NewLeakyBucket(tt.capacity, time.Millisecond)
			})
		})
	}
}
