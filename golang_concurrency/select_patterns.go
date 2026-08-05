package concurrency

import (
	"context"
	"fmt"
	"sync"
	"time"
)

/**
 * select Statement and Goroutine Lifecycle Control
 *
 * `select` is Go's channel-multiplexing construct. It waits until one of its
 * cases can proceed, then executes that case — chosen at random when multiple
 * are ready simultaneously.
 *
 *   select {
 *   case v := <-ch1:   // receives from ch1
 *   case ch2 <- x:     // sends to ch2
 *   case <-time.After(d): // deadline/timeout arm
 *   default:           // runs immediately if no other case is ready
 *   }
 *
 * Key facts:
 *   - A nil channel is never selected — disabling a case dynamically is as
 *     simple as setting that channel variable to nil.
 *   - `default` makes a select non-blocking (like a try-receive or try-send).
 *   - `select {}` (no cases) blocks forever — sometimes used to park main().
 *
 * Goroutine lifecycle:
 *   Every goroutine must have a clean exit path; leaked goroutines never release
 *   their memory or file descriptors. The patterns below show three mechanisms:
 *
 *     Done channel    — a channel closed by the caller to broadcast cancellation.
 *     context.Context — the standard library's generalisation of done channels,
 *                       carrying deadlines, cancellation, and request-scoped values.
 */

// -----------------------------------------------------------------------------
// Pattern 9: select to multiplex two channels
//
// select can read from multiple channels in a single statement, making progress
// as soon as any one of them is ready. This avoids the deadlock that would
// occur if you read them sequentially and the first channel is empty.
//
// Nil-channel trick:
//   Setting a channel variable to nil inside a select removes that case from
//   future iterations without needing a flag variable.
// -----------------------------------------------------------------------------

// multiplex drains two integer channels into one output channel.
// It exits — and closes output — only after both inputs are exhausted.
func multiplex(a, b <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for a != nil || b != nil { // loop until both channels are drained
			select {
			case v, ok := <-a:
				if !ok {
					a = nil // disable this case; never selected on a nil channel
					continue
				}
				out <- v
			case v, ok := <-b:
				if !ok {
					b = nil
					continue
				}
				out <- v
			}
		}
	}()
	return out
}

func ExampleMultiplex() {
	a := func() <-chan int {
		ch := make(chan int)
		go func() { ch <- 1; ch <- 3; close(ch) }()
		return ch
	}()

	b := func() <-chan int {
		ch := make(chan int)
		go func() { ch <- 2; ch <- 4; close(ch) }()
		return ch
	}()

	sum := 0
	for v := range multiplex(a, b) {
		sum += v
	}
	fmt.Println("sum:", sum) // always 10 (1+2+3+4), order non-deterministic
}

// -----------------------------------------------------------------------------
// Pattern 10: Timeout with time.After
//
// time.After(d) returns a <-chan time.Time that receives a value after duration d.
// Placing it as a select case creates a bounded wait: if the primary channel
// does not deliver within d, the timeout arm runs instead.
//
// Note: time.After leaks the underlying timer if the primary case fires first
// and the timer has not expired. For long-lived, high-frequency code, prefer
// time.NewTimer + t.Stop() to allow GC to reclaim the timer immediately.
// -----------------------------------------------------------------------------

// tryReceiveWithTimeout waits up to timeout for a value from ch.
// Returns (value, true) if a value arrived in time, (zero, false) otherwise.
func tryReceiveWithTimeout(ch <-chan int, timeout time.Duration) (int, bool) {
	select {
	case v := <-ch:
		return v, true
	case <-time.After(timeout):
		return 0, false
	}
}

func ExampleTimeout() {
	// Fast channel — value arrives before the deadline.
	fast := make(chan int, 1)
	fast <- 42
	if v, ok := tryReceiveWithTimeout(fast, time.Second); ok {
		fmt.Println("fast:", v) // fast: 42
	}

	// Slow channel — nothing arrives; timeout fires.
	slow := make(chan int) // unbuffered, nobody sends
	if _, ok := tryReceiveWithTimeout(slow, 10*time.Millisecond); !ok {
		fmt.Println("slow: timed out")
	}
}

// ExampleTimerStop shows the leak-free variant using time.NewTimer.
func ExampleTimerStop() {
	ch := make(chan int, 1)
	ch <- 99

	timer := time.NewTimer(time.Second)
	defer timer.Stop() // reclaim the timer even if ch fires first

	select {
	case v := <-ch:
		fmt.Println("received:", v) // received: 99
	case <-timer.C:
		fmt.Println("timed out")
	}
}

// -----------------------------------------------------------------------------
// Pattern 11: Non-blocking receive with default
//
// Adding `default` to a select makes it non-blocking: if no channel is ready,
// the default case runs immediately instead of waiting. Useful for polling
// or for trying to drain a channel without stalling the caller.
//
// Analogous to a try-lock in mutex terminology.
// -----------------------------------------------------------------------------

// tryReceive returns (value, true) if ch has a value ready, or (zero, false)
// without blocking.
func tryReceive(ch <-chan int) (int, bool) {
	select {
	case v := <-ch:
		return v, true
	default:
		return 0, false
	}
}

// trySend attempts to send v to ch without blocking.
// Returns true if the send succeeded (ch had room), false otherwise.
func trySend(ch chan<- int, v int) bool {
	select {
	case ch <- v:
		return true
	default:
		return false
	}
}

func ExampleNonBlocking() {
	ch := make(chan int, 1)

	// Channel is empty — tryReceive returns immediately with false.
	if _, ok := tryReceive(ch); !ok {
		fmt.Println("nothing ready")
	}

	// Send a value without blocking.
	trySend(ch, 7)

	// Now there is a value.
	if v, ok := tryReceive(ch); ok {
		fmt.Println("got:", v) // got: 7
	}
}

// -----------------------------------------------------------------------------
// Pattern 12: Done-channel cancellation
//
// A done channel (chan struct{}) is closed by the parent to broadcast a stop
// signal to all child goroutines. Goroutines select on both their work channel
// and done; they exit as soon as done is closed.
//
//   done := make(chan struct{})
//   close(done) // broadcasts to all goroutines selecting on done
//
// Why close instead of send?
//   close unblocks ALL receivers simultaneously. A send only unblocks one.
//   This makes close the right primitive for broadcasting cancellation.
//
// This pattern predates context.Context and is its conceptual basis. Prefer
// context.Context in new code (Pattern 13) — done channels are shown here to
// demystify how context works internally.
// -----------------------------------------------------------------------------

// infiniteCounter emits 0, 1, 2, … until done is closed.
func infiniteCounter(done <-chan struct{}) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := 0; ; i++ {
			select {
			case out <- i:
			case <-done: // parent closed done — exit cleanly
				return
			}
		}
	}()
	return out
}

func ExampleDoneChannel() {
	done := make(chan struct{})
	counter := infiniteCounter(done)

	// Read a few values then cancel.
	for i := 0; i < 5; i++ {
		fmt.Println(<-counter)
	}
	close(done) // signal the goroutine to stop

	// Drain the channel so the goroutine can exit and close out.
	for range counter {
	}
	fmt.Println("counter stopped")
}

// -----------------------------------------------------------------------------
// Pattern 13: context.Context cancellation (modern idiom)
//
// context.Context is the standard mechanism for goroutine lifecycle control.
// It carries:
//   - A cancellation signal (ctx.Done())
//   - An optional deadline or timeout
//   - Request-scoped key-value pairs
//
// context.WithCancel  — returns a child ctx + cancel(); call cancel() to stop.
// context.WithTimeout — automatically cancels after a duration.
// context.WithDeadline — cancels at a specific time.Time.
//
// Convention:
//   - Pass ctx as the first argument to any function that starts goroutines or
//     does I/O.
//   - Always call the cancel function (defer cancel()) to release resources even
//     if the context times out on its own.
//   - Never store ctx in a struct; pass it explicitly per-call.
// -----------------------------------------------------------------------------

// fetchWithContext simulates an HTTP request that honours cancellation.
// In real code this would call http.NewRequestWithContext or a DB driver's
// context-aware method; here we use time.Sleep to model latency.
func fetchWithContext(ctx context.Context, url string) (string, error) {
	done := make(chan string, 1)

	go func() {
		// Simulate variable network latency.
		time.Sleep(50 * time.Millisecond)
		done <- fmt.Sprintf("response from %s", url)
	}()

	select {
	case result := <-done:
		return result, nil
	case <-ctx.Done():
		return "", ctx.Err() // context.DeadlineExceeded or context.Canceled
	}
}

func ExampleContextWithTimeout() {
	// Give the fetch 200 ms — more than enough for the 50 ms simulated latency.
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	result, err := fetchWithContext(ctx, "https://example.com")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(result) // response from https://example.com
}

func ExampleContextCanceled() {
	// Give only 10 ms — less than the 50 ms simulated latency.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := fetchWithContext(ctx, "https://slow.example.com")
	fmt.Println(err) // context deadline exceeded
}

// ExampleContextManualCancel shows context.WithCancel — useful when the
// cancellation trigger is an event (e.g. first result found) rather than time.
func ExampleContextManualCancel() {
	ctx, cancel := context.WithCancel(context.Background())

	results := make(chan int, 10)
	var wg sync.WaitGroup

	// Launch 5 workers; all share the same context.
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			select {
			case results <- id * id:
			case <-ctx.Done(): // context was cancelled — exit without sending
			}
		}(i)
	}

	// Collect the first result, then cancel all remaining workers.
	first := <-results
	cancel() // signal all remaining goroutines to stop

	wg.Wait()
	close(results)

	fmt.Println("first result:", first)
}

