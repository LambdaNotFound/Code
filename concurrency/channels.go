package concurrency

import (
	"fmt"
	"sync"
)

/**
 * Goroutines and Channels — Core Mechanics
 *
 * Go's concurrency model is built on CSP (Communicating Sequential Processes):
 * goroutines are independent units of execution; channels are the typed conduits
 * through which they communicate.
 *
 * Guiding principle:
 *   "Do not communicate by sharing memory;
 *    share memory by communicating." — Rob Pike
 *
 * A goroutine is launched with the `go` keyword and runs concurrently with its
 * caller. It is cheap (~2 KB initial stack, grown on demand) — tens of thousands
 * can coexist. Goroutines are multiplexed onto OS threads by the Go scheduler.
 *
 * A channel is a typed, goroutine-safe pipe.
 *   Unbuffered:  make(chan T)     — sender blocks until receiver is ready (rendezvous)
 *   Buffered:    make(chan T, n)  — sender blocks only when the buffer is full
 *
 * Closing a channel:
 *   close(ch) signals "no more values". Only the sender should close.
 *   A range loop over a channel exits when the channel is closed.
 *   Sending on a closed channel panics.
 *   Receiving from a closed, empty channel returns (zero, false).
 */

// -----------------------------------------------------------------------------
// Pattern 1: Unbuffered channel — goroutine rendezvous
//
// An unbuffered channel has no buffer: the send blocks until a receiver is
// waiting and the receive blocks until a sender fires. Both goroutines meet
// at the channel — a synchronisation point sometimes called a rendezvous.
//
//   goroutine ──── ch <- v ────→ ch ──── v := <-ch ────→ caller
//                (blocks here)          (blocks here)
//
// Use when: the sender must not proceed until the receiver has the value
// (e.g. handing off ownership of a resource, signalling completion).
// -----------------------------------------------------------------------------

func ExampleUnbufferedChannel() {
	ch := make(chan string) // unbuffered

	go func() {
		ch <- "hello from goroutine" // blocks until main receives
	}()

	msg := <-ch // blocks until the goroutine sends
	fmt.Println(msg)
	// Output: hello from goroutine
}

// ExampleDoneSignal shows the common pattern of using a channel of struct{}
// (zero size) purely as a completion signal rather than a value carrier.
func ExampleDoneSignal() {
	done := make(chan struct{})

	go func() {
		fmt.Println("working...")
		close(done) // broadcast: any receiver unblocks
	}()

	<-done // wait for the goroutine to finish
	fmt.Println("goroutine finished")
}

// -----------------------------------------------------------------------------
// Pattern 2: Buffered channel — asynchronous handoff
//
// A buffered channel holds up to n values without a receiver being ready.
// The sender only blocks when the buffer is full; the receiver only blocks
// when the buffer is empty. This decouples producer and consumer speeds.
//
//   producer ─→ [ v1 | v2 | __ ] ch ─→ consumer
//               (sends until cap=3 is full, then blocks)
//
// Use when: bursts of work should be absorbed without stalling the producer,
// or when a known, fixed number of sends will happen before any receive.
// -----------------------------------------------------------------------------

func ExampleBufferedChannel() {
	ch := make(chan int, 3) // buffer holds up to 3 values

	// All three sends complete without a receiver — buffer absorbs them.
	ch <- 1
	ch <- 2
	ch <- 3

	fmt.Println(<-ch) // 1
	fmt.Println(<-ch) // 2
	fmt.Println(<-ch) // 3
}

// ExampleBufferedBackpressure demonstrates that a full buffer stalls the sender,
// providing natural flow-control (backpressure) between producer and consumer.
func ExampleBufferedBackpressure() {
	jobs := make(chan int, 2) // small buffer to keep the example short
	var wg sync.WaitGroup

	// Consumer: processes jobs as they arrive.
	wg.Add(1)
	go func() {
		defer wg.Done()
		for j := range jobs {
			fmt.Printf("processing job %d\n", j)
		}
	}()

	// Producer: sends 5 jobs; will stall whenever the buffer is full.
	for i := 1; i <= 5; i++ {
		jobs <- i // blocks when buffer holds 2 unprocessed jobs
	}
	close(jobs) // signal: no more jobs

	wg.Wait()
}

// -----------------------------------------------------------------------------
// Pattern 3: Range-over-close — draining a channel
//
// `for v := range ch` receives values until the channel is closed and empty.
// The loop body runs for every value; it exits automatically on close.
//
// Rules:
//   - Only the sender (the goroutine that owns the channel) should call close.
//   - Closing an already-closed channel panics.
//   - Closing a nil channel panics.
//   - Receiving from a closed, drained channel returns (zero, false) instantly.
// -----------------------------------------------------------------------------

// integers returns a channel that emits the given values then closes.
func integers(vals ...int) <-chan int {
	ch := make(chan int)
	go func() {
		for _, v := range vals {
			ch <- v
		}
		close(ch) // signal: no more values
	}()
	return ch
}

func ExampleRangeOverClose() {
	for v := range integers(10, 20, 30) {
		fmt.Println(v)
	}
	// prints 10, 20, 30 then the loop exits cleanly
}

// ExampleReceiveOkIdiom shows how to detect a closed channel without range.
// The two-value receive `v, ok := <-ch` returns ok=false when the channel
// is closed and drained — a safe way to break out of a select loop.
func ExampleReceiveOkIdiom() {
	ch := make(chan int, 2)
	ch <- 7
	close(ch)

	for {
		v, ok := <-ch
		if !ok {
			fmt.Println("channel closed")
			break
		}
		fmt.Println("received:", v)
	}
}

// -----------------------------------------------------------------------------
// Pattern 4: Pipeline — chained goroutine stages
//
// A pipeline is a series of goroutines connected by channels. Each stage:
//   1. Receives values from an inbound channel.
//   2. Applies a transformation.
//   3. Sends results to an outbound channel it owns and closes.
//
//   generate ─→ ch1 ─→ double ─→ ch2 ─→ filter ─→ ch3 ─→ sink (range)
//
// Benefits:
//   - Stages run concurrently: generate produces while double transforms.
//   - Each stage is independently testable.
//   - Memory stays bounded: values flow one-by-one, not collected into slices.
//   - Adding or removing a stage requires only changing which channels connect.
// -----------------------------------------------------------------------------

// pipelineGenerate emits each value in vals on a new channel then closes it.
func pipelineGenerate(vals ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, v := range vals {
			out <- v
		}
		close(out)
	}()
	return out
}

// pipelineDouble reads from in, doubles each value, and sends to a new channel.
func pipelineDouble(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for v := range in {
			out <- v * 2
		}
		close(out)
	}()
	return out
}

// pipelineFilter reads from in and forwards only values satisfying pred.
func pipelineFilter(in <-chan int, pred func(int) bool) <-chan int {
	out := make(chan int)
	go func() {
		for v := range in {
			if pred(v) {
				out <- v
			}
		}
		close(out)
	}()
	return out
}

// ExamplePipeline wires three stages together.
// Data flows: generate(1..5) → double → filter(>5) → print
//
// Concurrently:
//   stage 1 generates 1,2,3,4,5
//   stage 2 doubles to  2,4,6,8,10
//   stage 3 keeps only  6,8,10
func ExamplePipeline() {
	isGreaterThanFive := func(n int) bool { return n > 5 }

	src := pipelineGenerate(1, 2, 3, 4, 5)
	doubled := pipelineDouble(src)
	filtered := pipelineFilter(doubled, isGreaterThanFive)

	for v := range filtered {
		fmt.Println(v) // 6, 8, 10
	}
}
