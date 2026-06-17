package apidesign

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MaxStack(t *testing.T) {
	t.Run("push and pop lifo order", func(t *testing.T) {
		s := ConstructorMaxStack()
		s.Push(1)
		s.Push(2)
		s.Push(3)
		assert.Equal(t, 3, s.Pop())
		assert.Equal(t, 2, s.Pop())
		assert.Equal(t, 1, s.Pop())
	})

	t.Run("top does not remove element", func(t *testing.T) {
		s := ConstructorMaxStack()
		s.Push(5)
		assert.Equal(t, 5, s.Top())
		assert.Equal(t, 5, s.Top())
		assert.Equal(t, 5, s.Pop())
	})

	t.Run("peek max returns largest value", func(t *testing.T) {
		s := ConstructorMaxStack()
		s.Push(3)
		s.Push(1)
		s.Push(5)
		s.Push(2)
		assert.Equal(t, 5, s.PeekMax())
	})

	t.Run("pop max removes max from middle of stack", func(t *testing.T) {
		s := ConstructorMaxStack()
		s.Push(1)
		s.Push(5)
		s.Push(3)
		assert.Equal(t, 5, s.PopMax())
		assert.Equal(t, 3, s.Top())
		assert.Equal(t, 3, s.PeekMax())
	})

	t.Run("duplicate values: pop max removes most recently pushed", func(t *testing.T) {
		s := ConstructorMaxStack()
		s.Push(5)
		s.Push(1)
		s.Push(5) // second 5
		assert.Equal(t, 5, s.PopMax()) // removes the second (most recent) 5
		assert.Equal(t, 1, s.Top())    // 1 is now on top
		assert.Equal(t, 5, s.PeekMax())
		assert.Equal(t, 5, s.PopMax()) // removes the first 5
		assert.Equal(t, 1, s.Top())
	})

	t.Run("leetcode example", func(t *testing.T) {
		s := ConstructorMaxStack()
		s.Push(5)
		s.Push(1)
		s.Push(5)
		assert.Equal(t, 5, s.Top())
		assert.Equal(t, 5, s.PopMax())
		assert.Equal(t, 1, s.Top())
		assert.Equal(t, 5, s.PeekMax())
		assert.Equal(t, 1, s.Pop())
		assert.Equal(t, 5, s.Top())
	})

	t.Run("pop after pop max", func(t *testing.T) {
		s := ConstructorMaxStack()
		s.Push(2)
		s.Push(4)
		s.Push(1)
		assert.Equal(t, 4, s.PopMax())
		assert.Equal(t, 1, s.Pop())
		assert.Equal(t, 2, s.Top())
	})
}
