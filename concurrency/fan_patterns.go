package concurrency

import (
	"fmt"
	"sync"
)

/**
 * Fan-Out, Fan-In, and Worker Pool
 *
 * These three patterns scale work across multiple goroutines and collect results.
 *
 * Fan-Out   — one input channel, N workers each with their own result channel.
 *             Best when workers produce independent streams consumed separately.
 *
 * Fan-In    — N input channels merged into one output channel.
 *             The inverse of fan-out; commonly used after fan-out to recombine.
 *
 * Worker Pool — N workers all share one job channel AND one result channel.
 *               Best when results need to be collected in one place and order
 *               does not matter (results arrive as workers finish).
 *
 * Semaphore  — buffered channel used as a counting semaphore to cap the number
 *               of goroutines running concurrently without a fixed pool.
 */

// -----------------------------------------------------------------------------
// Pattern 5: Fan-Out
//
// A single source of jobs is consumed by N worker goroutines running in
// parallel. Each worker owns its result channel, which it closes when done.
//
//              ┌─ worker 0 ──→ resultCh[0]
//   jobs ch ──→├─ worker 1 ──→ resultCh[1]
//              └─ worker 2 ──→ resultCh[2]
//
// Use when:
//   - Workers produce independent streams that will be consumed separately,
//     or recombined with fan-in.
//   - You know the number of workers upfront.
//
// Caveat: if the shared jobs channel is closed, all workers exit cleanly via
// their range loops. Each worker closes its own result channel.
// -----------------------------------------------------------------------------

// fanOut distributes jobs across numWorkers goroutines.
// Each goroutine squares its input and sends results on its private channel.
// Returns a slice of per-worker result channels; each is closed when that
// worker finishes consuming all jobs.
//
// NOTE: all workers share the same jobs channel — whichever goroutine is free
// picks the next job, so work is distributed dynamically.
func fanOut(jobs <-chan int, numWorkers int) []<-chan int {
	results := make([]<-chan int, numWorkers)

	for i := range numWorkers {
		out := make(chan int)
		results[i] = out

		go func() {
			defer close(out)
			for j := range jobs {
				out <- j * j // worker's computation
			}
		}()
	}

	return results
}

func ExampleFanOut() {
	jobs := make(chan int, 6)
	for i := 1; i <= 6; i++ {
		jobs <- i
	}
	close(jobs)

	resultChans := fanOut(jobs, 3)

	// Drain each worker's channel sequentially (order within each is preserved).
	for i, ch := range resultChans {
		for v := range ch {
			fmt.Printf("worker %d: %d\n", i, v)
		}
	}
}

// -----------------------------------------------------------------------------
// Pattern 6: Fan-In (Merge)
//
// Multiple goroutines each produce values on their own channel. A merge goroutine
// drains all of them into a single output channel that consumers read from.
//
//   inCh[0] ─┐
//   inCh[1] ──┼──→ merge() ──→ merged
//   inCh[2] ─┘
//
// Implementation detail:
//   A sync.WaitGroup tracks when every input goroutine has finished. A separate
//   goroutine waits on the WaitGroup and then closes merged, signalling
//   downstream that all values have been produced.
//
// Use when:
//   - Results from fan-out workers need to be collected into one stream.
//   - Order of results across workers does not matter (they arrive as-ready).
// -----------------------------------------------------------------------------

// merge drains all channels in inputs into a single output channel.
// The output channel is closed exactly once, after all inputs are exhausted.
func merge(inputs ...<-chan int) <-chan int {
	merged := make(chan int)
	var wg sync.WaitGroup

	// For each input channel, launch a goroutine that forwards its values.
	relay := func(ch <-chan int) {
		defer wg.Done()
		for v := range ch {
			merged <- v
		}
	}

	wg.Add(len(inputs))
	for _, ch := range inputs {
		go relay(ch)
	}

	// Close merged only after every relay goroutine has exited.
	go func() {
		wg.Wait()
		close(merged)
	}()

	return merged
}

func ExampleFanIn() {
	// Two independent producers.
	evens := func() <-chan int {
		ch := make(chan int)
		go func() {
			for _, v := range []int{2, 4, 6} {
				ch <- v
			}
			close(ch)
		}()
		return ch
	}()

	odds := func() <-chan int {
		ch := make(chan int)
		go func() {
			for _, v := range []int{1, 3, 5} {
				ch <- v
			}
			close(ch)
		}()
		return ch
	}()

	// Merge both streams; drain the combined output.
	// Values arrive in non-deterministic order (whichever producer is faster).
	sum := 0
	for v := range merge(evens, odds) {
		sum += v
	}
	fmt.Println("sum:", sum) // always 21
}

// -----------------------------------------------------------------------------
// Pattern 7: Worker Pool
//
// A fixed pool of N goroutines all share ONE jobs channel and ONE results
// channel. The scheduler dispatches each job to whichever worker is free.
//
//              ┌─ worker 0 ─┐
//   jobs ch ──→├─ worker 1 ──┼──→ results ch
//              └─ worker 2 ─┘
//
// Difference from fan-out:
//   Fan-out: each worker owns its result channel → N result channels.
//   Worker pool: single shared result channel → easier to aggregate.
//
// Backpressure: if the result channel fills up (when buffered) or no consumer
// is reading (when unbuffered), workers stall, which stalls producers — a
// natural flow-control mechanism requiring no extra code.
//
// Lifecycle:
//   1. Close jobs when all work is submitted.
//   2. Workers exit their range loops automatically.
//   3. A supervisor goroutine waits for all workers via WaitGroup, then closes results.
//   4. The consumer's range loop on results exits when results is closed.
// -----------------------------------------------------------------------------

// workerPool spawns numWorkers goroutines that consume from jobs and send
// squared values to the returned results channel.
func workerPool(numWorkers int, jobs <-chan int) <-chan int {
	results := make(chan int)
	var wg sync.WaitGroup

	worker := func() {
		defer wg.Done()
		for j := range jobs {
			results <- j * j
		}
	}

	wg.Add(numWorkers)
	for range numWorkers {
		go worker()
	}

	// Supervisor: close results once all workers are done.
	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}

func ExampleWorkerPool() {
	const numJobs = 9
	const numWorkers = 3

	jobs := make(chan int, numJobs)
	for i := 1; i <= numJobs; i++ {
		jobs <- i
	}
	close(jobs)

	sum := 0
	for r := range workerPool(numWorkers, jobs) {
		sum += r
	}
	// 1²+2²+…+9² = 285
	fmt.Println("sum of squares:", sum) // 285
}

// -----------------------------------------------------------------------------
// Pattern 8: Semaphore via buffered channel
//
// A buffered channel of capacity N acts as a counting semaphore that allows
// at most N goroutines to run a critical section concurrently. No fixed pool
// of goroutines is maintained — instead, each goroutine acquires a "token"
// from the semaphore before doing work and releases it afterward.
//
//   sem := make(chan struct{}, N)
//   sem <- struct{}{}   // acquire (blocks when N goroutines are already inside)
//   defer func() { <-sem }()    // release
//
// Use when:
//   - You want to launch many goroutines dynamically (e.g. one per request)
//     but cap concurrent resource usage (DB connections, file descriptors).
//   - A fixed worker pool would require pre-allocating workers for work that
//     may never arrive.
// -----------------------------------------------------------------------------

// concurrentFetch simulates fetching URLs with at most maxConcurrent
// goroutines running at the same time.
func concurrentFetch(urls []string, maxConcurrent int) []string {
	sem := make(chan struct{}, maxConcurrent) // semaphore
	results := make([]string, len(urls))
	var wg sync.WaitGroup

	for i, url := range urls {
		wg.Add(1)
		go func(idx int, u string) {
			defer wg.Done()

			sem <- struct{}{}        // acquire token (blocks if maxConcurrent reached)
			defer func() { <-sem }() // release token on exit

			// Simulate network I/O.
			results[idx] = fmt.Sprintf("fetched: %s", u)
		}(i, url)
	}

	wg.Wait()
	return results
}

func ExampleSemaphore() {
	urls := []string{"https://a.com", "https://b.com", "https://c.com", "https://d.com"}
	results := concurrentFetch(urls, 2) // at most 2 concurrent fetches

	for _, r := range results {
		fmt.Println(r)
	}
}
