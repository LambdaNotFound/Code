package concurrency

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLockFreeStack_EmptyPop(t *testing.T) {
	s := NewLockFreeStack[int]()
	val, ok := s.Pop()
	assert.False(t, ok)
	assert.Equal(t, 0, val)
	assert.True(t, s.IsEmpty())
	assert.Equal(t, 0, s.Size())
}

func TestLockFreeStack_EmptyPeek(t *testing.T) {
	s := NewLockFreeStack[string]()
	val, ok := s.Peek()
	assert.False(t, ok)
	assert.Equal(t, "", val)
}

func TestLockFreeStack_SinglePushPop(t *testing.T) {
	s := NewLockFreeStack[string]()
	s.Push("hello")
	assert.Equal(t, 1, s.Size())
	assert.False(t, s.IsEmpty())

	val, ok := s.Pop()
	assert.True(t, ok)
	assert.Equal(t, "hello", val)
	assert.Equal(t, 0, s.Size())
	assert.True(t, s.IsEmpty())
}

func TestLockFreeStack_LIFOOrder(t *testing.T) {
	s := NewLockFreeStack[int]()
	for i := 0; i < 5; i++ {
		s.Push(i)
	}
	assert.Equal(t, 5, s.Size())

	// LIFO: last pushed (4) should come out first
	for i := 4; i >= 0; i-- {
		val, ok := s.Pop()
		assert.True(t, ok)
		assert.Equal(t, i, val)
	}

	_, ok := s.Pop()
	assert.False(t, ok)
}

func TestLockFreeStack_Peek(t *testing.T) {
	s := NewLockFreeStack[int]()
	s.Push(10)
	s.Push(20)

	val, ok := s.Peek()
	assert.True(t, ok)
	assert.Equal(t, 20, val)
	assert.Equal(t, 2, s.Size()) // Peek must not remove the element
}

func TestLockFreeStack_ConcurrentPush(t *testing.T) {
	s := NewLockFreeStack[int]()
	const n = 1000
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			s.Push(v)
		}(i)
	}
	wg.Wait()

	assert.Equal(t, n, s.Size())
}

func TestLockFreeStack_ConcurrentPop(t *testing.T) {
	s := NewLockFreeStack[int]()
	const n = 1000
	for i := 0; i < n; i++ {
		s.Push(i)
	}

	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.Pop()
		}()
	}
	wg.Wait()

	assert.Equal(t, 0, s.Size())
}

func TestLockFreeStack_ConcurrentPushPop(t *testing.T) {
	s := NewLockFreeStack[int]()
	const goroutines = 50
	const itemsEach = 100
	total := goroutines * itemsEach

	var wg sync.WaitGroup

	// producers
	for p := 0; p < goroutines; p++ {
		wg.Add(1)
		go func(base int) {
			defer wg.Done()
			for i := 0; i < itemsEach; i++ {
				s.Push(base + i)
			}
		}(p * itemsEach)
	}

	// consumers
	results := make(chan int, total)
	var consumerWg sync.WaitGroup
	for c := 0; c < goroutines; c++ {
		consumerWg.Add(1)
		go func() {
			defer consumerWg.Done()
			collected := 0
			for collected < itemsEach {
				if val, ok := s.Pop(); ok {
					results <- val
					collected++
				}
			}
		}()
	}

	wg.Wait()
	consumerWg.Wait()
	close(results)

	assert.Equal(t, 0, s.Size())
	assert.Equal(t, total, len(results))
}

func TestLockFreeStack_RaceDetector(t *testing.T) {
	s := NewLockFreeStack[int]()
	var wg sync.WaitGroup
	const goroutines = 50

	for i := 0; i < goroutines; i++ {
		wg.Add(3)
		go func(v int) {
			defer wg.Done()
			s.Push(v)
		}(i)
		go func() {
			defer wg.Done()
			s.Pop()
		}()
		go func() {
			defer wg.Done()
			s.Peek()
		}()
	}
	wg.Wait()
}
