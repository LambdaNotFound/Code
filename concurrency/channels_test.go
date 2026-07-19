package concurrency

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_UnbufferedChannel(t *testing.T) {
	ch := make(chan string) // unbuffered

	go func() {
		ch <- "hello from goroutine" // blocks until main receives
	}()

	msg := <-ch // blocks until the goroutine sends

	assert.Equal(t, "hello from goroutine", msg)
}

func Test_BufferedBackpressure_processesAllJobsInOrder(t *testing.T) {
	jobs := make(chan int, 2)
	var processed []int
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for j := range jobs {
			processed = append(processed, j)
		}
	}()

	for i := 1; i <= 5; i++ {
		jobs <- i // blocks once the buffer holds 2 unprocessed jobs
	}
	close(jobs)
	wg.Wait()

	assert.Equal(t, []int{1, 2, 3, 4, 5}, processed)
}

func Test_BufferedBackpressure_blocksWhenBufferIsFull(t *testing.T) {
	jobs := make(chan int, 2)
	release := make(chan struct{})
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-release // consumer holds off until told to start draining
		for range jobs {
		}
	}()

	jobs <- 1 // fills the buffer
	jobs <- 2 // buffer now full

	sent := make(chan struct{})
	go func() {
		jobs <- 3 // buffer is full and nobody is draining yet — should block
		close(sent)
	}()

	select {
	case <-sent:
		t.Fatal("send on a full buffered channel should have blocked")
	case <-time.After(50 * time.Millisecond):
		// expected: send is still blocked
	}

	close(release) // let the consumer start draining, freeing up buffer space

	select {
	case <-sent:
		// send completed once the consumer made room
	case <-time.After(time.Second):
		t.Fatal("send should have unblocked once the consumer drained the buffer")
	}

	close(jobs)
	wg.Wait()
}
