package concurrency

import (
	"context"
	"sync"
)

/**
 * Waiter — a reusable "done channel" (see channels.go's ExampleDoneSignal and
 * select_patterns.go Pattern 12) packaged as a type, adding two things a raw
 * chan struct{} doesn't give you for free:
 *
 *   - Idempotent close: closing an already-closed channel panics; sync.Once
 *     makes Close safe to call more than once, from one or many goroutines.
 *   - Context-aware waiting: Wait races the done signal against ctx, so a
 *     caller can bound how long it's willing to wait regardless of whether
 *     the Waiter ever completes.
 */

// Waiter broadcasts a one-time completion signal to any number of waiters.
// The zero value is not usable — construct with New.
type Waiter struct {
	c    chan struct{}
	once sync.Once
}

// WaitFn is the shape of a cancellable unit of work: run it with a context
// and it reports whether it completed or was cancelled/timed out.
type WaitFn func(ctx context.Context) error

// New returns a Waiter that has not yet completed.
func New() *Waiter {
	return &Waiter{c: make(chan struct{})}
}

// Ch exposes the underlying done channel so it can be used as a select case
// alongside other channels, instead of always blocking via Wait.
func (s *Waiter) Ch() <-chan struct{} {
	return s.c
}

// Close signals completion. Safe to call multiple times or concurrently —
// only the first call actually closes the channel.
func (s *Waiter) Close() {
	s.once.Do(func() { close(s.c) })
}

// Wait blocks until Close is called or ctx is done, whichever comes first.
// A nil *Waiter is treated as already complete, so callers holding an
// optional/unset Waiter field don't need a nil check before calling Wait.
func (s *Waiter) Wait(ctx context.Context) error {
	if s == nil {
		return nil
	}
	select {
	case <-s.c:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
