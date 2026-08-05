package concurrency

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLockFreeQueue_EmptyDequeue(t *testing.T) {
	q := NewLockFreeQueue[int]()
	val, ok := q.Dequeue()
	assert.False(t, ok)
	assert.Equal(t, 0, val)
	assert.True(t, q.IsEmpty())
	assert.Equal(t, 0, q.Size())
}

func TestLockFreeQueue_SingleEnqueueDequeue(t *testing.T) {
	q := NewLockFreeQueue[string]()
	q.Enqueue("hello")
	assert.Equal(t, 1, q.Size())
	assert.False(t, q.IsEmpty())

	val, ok := q.Dequeue()
	assert.True(t, ok)
	assert.Equal(t, "hello", val)
	assert.Equal(t, 0, q.Size())
	assert.True(t, q.IsEmpty())
}

func TestLockFreeQueue_FIFOOrder(t *testing.T) {
	q := NewLockFreeQueue[int]()
	for i := 0; i < 5; i++ {
		q.Enqueue(i)
	}
	assert.Equal(t, 5, q.Size())

	for i := 0; i < 5; i++ {
		val, ok := q.Dequeue()
		assert.True(t, ok)
		assert.Equal(t, i, val)
	}

	_, ok := q.Dequeue()
	assert.False(t, ok)
}

func TestLockFreeQueue_ConcurrentEnqueue(t *testing.T) {
	q := NewLockFreeQueue[int]()
	const n = 1000
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			q.Enqueue(v)
		}(i)
	}
	wg.Wait()

	assert.Equal(t, n, q.Size())
}

func TestLockFreeQueue_ConcurrentDequeue(t *testing.T) {
	q := NewLockFreeQueue[int]()
	const n = 1000
	for i := 0; i < n; i++ {
		q.Enqueue(i)
	}

	var wg sync.WaitGroup
	dequeued := make([]int, n)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(slot int) {
			defer wg.Done()
			val, ok := q.Dequeue()
			if ok {
				dequeued[slot] = val + 1 // sentinel: 0 means "nothing dequeued"
			}
		}(i)
	}
	wg.Wait()

	assert.Equal(t, 0, q.Size())
}

func TestLockFreeQueue_ConcurrentEnqueueDequeue(t *testing.T) {
	q := NewLockFreeQueue[int]()
	const producers = 10
	const itemsPerProducer = 100
	total := producers * itemsPerProducer

	var wg sync.WaitGroup

	// producers
	for p := 0; p < producers; p++ {
		wg.Add(1)
		go func(base int) {
			defer wg.Done()
			for i := 0; i < itemsPerProducer; i++ {
				q.Enqueue(base + i)
			}
		}(p * itemsPerProducer)
	}

	// consumers — drain until we've collected every item
	results := make(chan int, total)
	var consumerWg sync.WaitGroup
	for c := 0; c < producers; c++ {
		consumerWg.Add(1)
		go func() {
			defer consumerWg.Done()
			collected := 0
			for collected < itemsPerProducer {
				if val, ok := q.Dequeue(); ok {
					results <- val
					collected++
				}
			}
		}()
	}

	wg.Wait()
	consumerWg.Wait()
	close(results)

	assert.Equal(t, 0, q.Size())
	assert.Equal(t, total, len(results))
}

func TestLockFreeQueue_RaceDetector(t *testing.T) {
	// Run with -race to verify no data races.
	q := NewLockFreeQueue[int]()
	var wg sync.WaitGroup
	const goroutines = 50

	for i := 0; i < goroutines; i++ {
		wg.Add(2)
		go func(v int) {
			defer wg.Done()
			q.Enqueue(v)
		}(i)
		go func() {
			defer wg.Done()
			q.Dequeue()
		}()
	}
	wg.Wait()
}
