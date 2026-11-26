package apidesign

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MyQueue(t *testing.T) {
	t.Run("empty_on_construction", func(t *testing.T) {
		q := ConstructorMyQueue()
		assert.True(t, q.Empty())
	})

	t.Run("not_empty_after_push", func(t *testing.T) {
		q := ConstructorMyQueue()
		q.Push(1)
		assert.False(t, q.Empty())
	})

	t.Run("fifo_order_pop", func(t *testing.T) {
		q := ConstructorMyQueue()
		q.Push(1)
		q.Push(2)
		q.Push(3)
		assert.Equal(t, 1, q.Pop())
		assert.Equal(t, 2, q.Pop())
		assert.Equal(t, 3, q.Pop())
	})

	t.Run("peek_does_not_remove", func(t *testing.T) {
		q := ConstructorMyQueue()
		q.Push(10)
		q.Push(20)
		assert.Equal(t, 10, q.Peek())
		assert.Equal(t, 10, q.Peek()) // still there
		assert.Equal(t, 10, q.Pop())  // now removed
	})

	t.Run("empty_after_all_popped", func(t *testing.T) {
		q := ConstructorMyQueue()
		q.Push(1)
		q.Pop()
		assert.True(t, q.Empty())
	})

	t.Run("interleaved_push_pop", func(t *testing.T) {
		q := ConstructorMyQueue()
		q.Push(1)
		q.Push(2)
		assert.Equal(t, 1, q.Pop())
		q.Push(3)
		assert.Equal(t, 2, q.Pop())
		assert.Equal(t, 3, q.Pop())
		assert.True(t, q.Empty())
	})

	t.Run("leetcode_example", func(t *testing.T) {
		q := ConstructorMyQueue()
		q.Push(1)
		q.Push(2)
		assert.Equal(t, 1, q.Peek())
		assert.Equal(t, 1, q.Pop())
		assert.False(t, q.Empty())
	})
}
