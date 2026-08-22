package concurrency

import (
	"context"
	"errors"
	"sync"
	"time"
)

/**
 * Leaky Bucket — a bounded queue (the "bucket") that admits bursts up to its
 * capacity but drains ("leaks") at a fixed, constant rate, turning whatever
 * burst pattern callers arrive in into a steady output rate.
 *
 *   callers --> [ o o o o . . ]  --tick, every leakInterval-->  one leaks out
 *               ^^^^^^^^^^^^^^
 *               buffered chan struct{}, size = capacity
 *
 * Each arrival tries to place a token into a buffered channel sized to the
 * bucket's capacity. A single background goroutine (drain) wakes on a
 * time.Ticker and, on every tick, removes at most one token — that's the
 * "leak": output is paced at exactly 1/leakInterval regardless of how bursty
 * the input is.
 *
 * Two admission modes are exposed because "what happens when the bucket is
 * full" is a policy decision, not a mechanism decision:
 *
 *   Allow() bool     — non-blocking. Full bucket -> immediate reject. This is
 *                      the textbook leaky-bucket behavior: excess traffic is
 *                      dropped rather than queued a second time on top of the
 *                      bucket itself.
 *   Wait(ctx) error  — blocking. Full bucket -> park until a slot leaks, ctx
 *                      is cancelled, or the bucket is closed. Useful when a
 *                      caller would rather backpressure than drop, e.g. an
 *                      internal work queue rather than a public API edge.
 *
 * Concurrency invariant:
 *   - stop is closed exactly once, via sync.Once in Close(), by whichever
 *     goroutine calls Close() first — idempotent and safe to call
 *     concurrently or multiple times (same pattern as Waiter.Close).
 *   - queue itself is NEVER closed. Closing a buffered channel while other
 *     goroutines may still be sending to it (blocked in Wait) is a panic
 *     waiting to happen; stop is a separate signal that both Wait and the
 *     drain loop select on instead. queue is simply abandoned to the GC once
 *     the bucket becomes unreachable.
 *   - Only the drain goroutine ever receives from queue; only Allow/Wait
 *     callers ever send to it. That fixed single-reader/single-role-per-side
 *     split is what makes the channel itself race-free with no mutex.
 *   - Close does not force Allow to start rejecting: Allow only ever
 *     inspects queue, never stop, since it never blocks and a stopped drain
 *     loop just means the bucket stops leaking, not that it becomes invalid.
 *     Wait, which can block indefinitely, does watch stop and returns
 *     ErrClosed instead of hanging forever on a bucket that will never leak
 *     again.
 */

// ErrClosed is returned by Wait when the bucket is closed while the caller
// is blocked waiting for capacity to free up.
var ErrClosed = errors.New("leaky bucket: closed")

// LeakyBucket is a concurrency-safe rate limiter: a bounded queue drained at
// a constant rate. Construct with NewLeakyBucket; the zero value is not
// usable.
type LeakyBucket struct {
	queue  chan struct{}
	ticker *time.Ticker
	stop   chan struct{}
	once   sync.Once
}

// NewLeakyBucket returns a LeakyBucket that holds up to capacity requests and
// leaks (drains) one every leakInterval. It starts a background drain
// goroutine that runs until Close is called — callers must call Close once
// the bucket is no longer needed, or the goroutine and its ticker leak.
func NewLeakyBucket(capacity int, leakInterval time.Duration) *LeakyBucket {
	if capacity <= 0 {
		panic("concurrency: leaky bucket capacity must be positive")
	}

	lb := &LeakyBucket{
		queue:  make(chan struct{}, capacity),
		ticker: time.NewTicker(leakInterval),
		stop:   make(chan struct{}),
	}
	go lb.drain()
	return lb
}

// drain owns the ticker and the queue's only receive path: on every tick it
// removes at most one token, and it exits when stop is closed.
func (lb *LeakyBucket) drain() {
	defer lb.ticker.Stop()
	for {
		select {
		case <-lb.ticker.C:
			select {
			case <-lb.queue:
			default:
				// nothing queued this tick — the leak rate caps outflow, it
				// doesn't guarantee there's always something to drain
			}
		case <-lb.stop:
			return
		}
	}
}

// Allow reports whether a request may proceed right now, without blocking.
// It admits the request into the bucket if there is spare capacity, or
// rejects it immediately if the bucket is full.
func (lb *LeakyBucket) Allow() bool {
	select {
	case lb.queue <- struct{}{}:
		return true
	default:
		return false
	}
}

// Wait blocks until the request is admitted into the bucket, ctx is done, or
// the bucket is closed, whichever happens first.
func (lb *LeakyBucket) Wait(ctx context.Context) error {
	// A closed bucket with spare queue capacity leaves both branches of the
	// select below ready; select picks among ready cases at random, so
	// without this check a Wait call made after Close has already completed
	// could still be admitted instead of deterministically returning
	// ErrClosed. This non-blocking check makes the already-closed case
	// deterministic; a Close that races with the blocking select below is a
	// genuine, acceptable race between concurrent calls.
	select {
	case <-lb.stop:
		return ErrClosed
	default:
	}

	select {
	case lb.queue <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-lb.stop:
		return ErrClosed
	}
}

// Close stops the background drain goroutine. Safe to call multiple times or
// concurrently; only the first call has an effect. Callers currently blocked
// in Wait are released with ErrClosed.
func (lb *LeakyBucket) Close() {
	lb.once.Do(func() { close(lb.stop) })
}
