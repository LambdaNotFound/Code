package concurrency

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Waiter_WaitBlocksUntilClose(t *testing.T) {
	w := New()

	done := make(chan error, 1)
	go func() {
		done <- w.Wait(context.Background())
	}()

	select {
	case <-done:
		t.Fatal("Wait returned before Close was called")
	case <-time.After(50 * time.Millisecond):
		// expected: still waiting
	}

	w.Close()

	select {
	case err := <-done:
		assert.NoError(t, err)
	case <-time.After(time.Second):
		t.Fatal("Wait should have returned once Close was called")
	}
}

func Test_Waiter_WaitReturnsImmediatelyIfAlreadyClosed(t *testing.T) {
	w := New()
	w.Close()

	err := w.Wait(context.Background())

	assert.NoError(t, err)
}

func Test_Waiter_WaitReturnsContextErrorOnCancellation(t *testing.T) {
	w := New() // never closed

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	err := w.Wait(ctx)

	assert.ErrorIs(t, err, context.DeadlineExceeded)
}

func Test_Waiter_CloseIsIdempotent(t *testing.T) {
	w := New()

	assert.NotPanics(t, func() {
		w.Close()
		w.Close()
		w.Close()
	})
}

func Test_Waiter_CloseIsSafeForConcurrentCallers(t *testing.T) {
	w := New()
	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			w.Close()
		}()
	}

	assert.NotPanics(t, func() { wg.Wait() })
}

func Test_Waiter_NilWaiterWaitReturnsNil(t *testing.T) {
	var w *Waiter

	err := w.Wait(context.Background())

	assert.NoError(t, err)
}

func Test_Waiter_BroadcastsToAllWaiters(t *testing.T) {
	w := New()
	const numWaiters = 5

	var wg sync.WaitGroup
	results := make([]error, numWaiters)

	for i := 0; i < numWaiters; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			results[idx] = w.Wait(context.Background())
		}(i)
	}

	w.Close() // broadcasts to every waiter above, however far each has gotten

	wg.Wait()

	for _, err := range results {
		assert.NoError(t, err)
	}
}

func Test_Waiter_ChUsableInSelect(t *testing.T) {
	w := New()

	select {
	case <-w.Ch():
		t.Fatal("Ch() should not be ready before Close")
	default:
		// expected: not ready yet
	}

	w.Close()

	select {
	case <-w.Ch():
		// expected: ready after Close
	case <-time.After(time.Second):
		t.Fatal("Ch() should be ready after Close")
	}
}
