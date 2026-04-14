package containers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MinMaxStack(t *testing.T) {
	// ---------------------------------------------------------------------------
	// Empty stack — all methods return safe zero values
	// ---------------------------------------------------------------------------
	t.Run("empty stack", func(t *testing.T) {
		s := NewMinMaxStack[int]()
		assert.True(t, s.IsEmpty())
		assert.Equal(t, 0, s.Size())

		val, ok := s.GetMin()
		assert.False(t, ok)
		assert.Equal(t, 0, val)

		val, ok = s.GetMax()
		assert.False(t, ok)
		assert.Equal(t, 0, val)

		val, ok = s.Pop()
		assert.False(t, ok)
		assert.Equal(t, 0, val)

		val, ok = s.Peek()
		assert.False(t, ok)
		assert.Equal(t, 0, val)
	})

	// ---------------------------------------------------------------------------
	// Single element: GetMin == GetMax == the element
	// ---------------------------------------------------------------------------
	t.Run("single element", func(t *testing.T) {
		s := NewMinMaxStack[int]()
		s.Push(7)
		assert.Equal(t, 1, s.Size())
		assert.False(t, s.IsEmpty())

		top, ok := s.Peek()
		assert.True(t, ok)
		assert.Equal(t, 7, top)
		assert.Equal(t, 1, s.Size()) // Peek must not remove

		mn, ok := s.GetMin()
		assert.True(t, ok)
		assert.Equal(t, 7, mn)

		mx, ok := s.GetMax()
		assert.True(t, ok)
		assert.Equal(t, 7, mx)
	})

	// ---------------------------------------------------------------------------
	// Push increasing sequence — min anchored at bottom, max tracks top
	// ---------------------------------------------------------------------------
	t.Run("push increasing sequence", func(t *testing.T) {
		s := NewMinMaxStack[int]()
		for i, push := range []int{1, 2, 3, 4} {
			s.Push(push)
			mn, _ := s.GetMin()
			mx, _ := s.GetMax()
			assert.Equal(t, 1, mn, "min should stay 1 after push %d", i)
			assert.Equal(t, push, mx, "max should equal latest push %d", push)
		}
	})

	// ---------------------------------------------------------------------------
	// Push decreasing sequence — max anchored at bottom, min tracks top
	// ---------------------------------------------------------------------------
	t.Run("push decreasing sequence", func(t *testing.T) {
		s := NewMinMaxStack[int]()
		for i, push := range []int{4, 3, 2, 1} {
			s.Push(push)
			mn, _ := s.GetMin()
			mx, _ := s.GetMax()
			assert.Equal(t, push, mn, "min should equal latest push %d", push)
			assert.Equal(t, 4, mx, "max should stay 4 after push %d", i)
		}
	})

	// ---------------------------------------------------------------------------
	// Pop the current minimum — GetMin restores to previous min
	// ---------------------------------------------------------------------------
	t.Run("pop restores previous min", func(t *testing.T) {
		s := NewMinMaxStack[int]()
		s.Push(5)
		s.Push(2) // new min
		s.Push(8)

		mn, _ := s.GetMin()
		assert.Equal(t, 2, mn)

		s.Pop() // remove 8; min still 2
		mn, _ = s.GetMin()
		assert.Equal(t, 2, mn)

		s.Pop() // remove 2; min restores to 5
		mn, ok := s.GetMin()
		assert.True(t, ok)
		assert.Equal(t, 5, mn)
	})

	// ---------------------------------------------------------------------------
	// Pop the current maximum — GetMax restores to previous max
	// ---------------------------------------------------------------------------
	t.Run("pop restores previous max", func(t *testing.T) {
		s := NewMinMaxStack[int]()
		s.Push(3)
		s.Push(9) // new max
		s.Push(1)

		mx, _ := s.GetMax()
		assert.Equal(t, 9, mx)

		s.Pop() // remove 1; max still 9
		mx, _ = s.GetMax()
		assert.Equal(t, 9, mx)

		s.Pop() // remove 9; max restores to 3
		mx, ok := s.GetMax()
		assert.True(t, ok)
		assert.Equal(t, 3, mx)
	})

	// ---------------------------------------------------------------------------
	// Duplicate elements
	// ---------------------------------------------------------------------------
	t.Run("duplicate elements", func(t *testing.T) {
		s := NewMinMaxStack[int]()
		s.Push(4)
		s.Push(4)

		mn, _ := s.GetMin()
		mx, _ := s.GetMax()
		assert.Equal(t, 4, mn)
		assert.Equal(t, 4, mx)

		s.Pop() // one 4 removed — the other remains
		mn, ok := s.GetMin()
		assert.True(t, ok)
		assert.Equal(t, 4, mn)
		mx, ok = s.GetMax()
		assert.True(t, ok)
		assert.Equal(t, 4, mx)

		s.Pop() // stack now empty
		_, ok = s.GetMin()
		assert.False(t, ok)
		_, ok = s.GetMax()
		assert.False(t, ok)
	})

	// ---------------------------------------------------------------------------
	// Interleaved push/pop — verify min/max after every operation
	// ---------------------------------------------------------------------------
	t.Run("interleaved push pop", func(t *testing.T) {
		s := NewMinMaxStack[int]()

		s.Push(10)
		mn, _ := s.GetMin()
		mx, _ := s.GetMax()
		assert.Equal(t, 10, mn)
		assert.Equal(t, 10, mx)

		s.Push(3)
		mn, _ = s.GetMin()
		mx, _ = s.GetMax()
		assert.Equal(t, 3, mn)
		assert.Equal(t, 10, mx)

		s.Push(7)
		mn, _ = s.GetMin()
		mx, _ = s.GetMax()
		assert.Equal(t, 3, mn)
		assert.Equal(t, 10, mx)

		v, _ := s.Pop() // removes 7
		assert.Equal(t, 7, v)
		mn, _ = s.GetMin()
		mx, _ = s.GetMax()
		assert.Equal(t, 3, mn)
		assert.Equal(t, 10, mx)

		s.Push(1)
		mn, _ = s.GetMin()
		mx, _ = s.GetMax()
		assert.Equal(t, 1, mn)
		assert.Equal(t, 10, mx)

		v, _ = s.Pop() // removes 1
		assert.Equal(t, 1, v)
		mn, _ = s.GetMin()
		mx, _ = s.GetMax()
		assert.Equal(t, 3, mn)
		assert.Equal(t, 10, mx)

		v, _ = s.Pop() // removes 3
		assert.Equal(t, 3, v)
		mn, _ = s.GetMin()
		mx, _ = s.GetMax()
		assert.Equal(t, 10, mn)
		assert.Equal(t, 10, mx)
	})

	// ---------------------------------------------------------------------------
	// LeetCode 155 canonical example
	// ---------------------------------------------------------------------------
	t.Run("LeetCode 155", func(t *testing.T) {
		s := NewMinMaxStack[int]()
		s.Push(-2)
		s.Push(0)
		s.Push(-3)

		mn, ok := s.GetMin()
		assert.True(t, ok)
		assert.Equal(t, -3, mn)

		s.Pop() // removes -3

		mn, ok = s.GetMin()
		assert.True(t, ok)
		assert.Equal(t, -2, mn)

		top, ok := s.Peek()
		assert.True(t, ok)
		assert.Equal(t, 0, top)
	})

	// ---------------------------------------------------------------------------
	// Size and IsEmpty invariants — using string type
	// ---------------------------------------------------------------------------
	t.Run("size and IsEmpty with string type", func(t *testing.T) {
		s := NewMinMaxStack[string]()
		assert.True(t, s.IsEmpty())
		assert.Equal(t, 0, s.Size())

		s.Push("b")
		assert.False(t, s.IsEmpty())
		assert.Equal(t, 1, s.Size())

		s.Push("a")
		s.Push("c")
		assert.Equal(t, 3, s.Size())

		mn, _ := s.GetMin()
		mx, _ := s.GetMax()
		assert.Equal(t, "a", mn)
		assert.Equal(t, "c", mx)

		s.Pop()
		s.Pop()
		s.Pop()
		assert.True(t, s.IsEmpty())
		assert.Equal(t, 0, s.Size())
	})

	// ---------------------------------------------------------------------------
	// Float64 type — exercises the constraint on non-integer ordered type
	// ---------------------------------------------------------------------------
	t.Run("float64 type", func(t *testing.T) {
		s := NewMinMaxStack[float64]()
		s.Push(1.5)
		s.Push(2.5)
		s.Push(0.5)

		mn, ok := s.GetMin()
		assert.True(t, ok)
		assert.Equal(t, 0.5, mn)

		mx, ok := s.GetMax()
		assert.True(t, ok)
		assert.Equal(t, 2.5, mx)
	})
}
