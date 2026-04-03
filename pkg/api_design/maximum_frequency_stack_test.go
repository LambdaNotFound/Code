package apidesign

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FreqStack(t *testing.T) {
	t.Run("leetcode_example", func(t *testing.T) {
		// Push: 5,7,5,7,4,5 → freq: {5:3, 7:2, 4:1}
		// Pop returns most frequent; ties broken by most recently pushed
		fs := ConstructorFreqStack()
		fs.Push(5)
		fs.Push(7)
		fs.Push(5)
		fs.Push(7)
		fs.Push(4)
		fs.Push(5)
		assert.Equal(t, 5, fs.Pop()) // 5 appears 3 times (most frequent)
		assert.Equal(t, 7, fs.Pop()) // 5 and 7 both appear 2 times; 7 was pushed more recently at freq=2
		assert.Equal(t, 5, fs.Pop()) // 5 now most frequent at 2
		assert.Equal(t, 4, fs.Pop()) // 4, 5, 7 all at freq=1; 4 was pushed most recently
	})

	t.Run("single_element", func(t *testing.T) {
		fs := ConstructorFreqStack()
		fs.Push(42)
		assert.Equal(t, 42, fs.Pop())
	})

	t.Run("all_same_element_lifo", func(t *testing.T) {
		fs := ConstructorFreqStack()
		fs.Push(1)
		fs.Push(1)
		fs.Push(1)
		// all at different frequencies; most frequent first, then LIFO within freq
		assert.Equal(t, 1, fs.Pop())
		assert.Equal(t, 1, fs.Pop())
		assert.Equal(t, 1, fs.Pop())
	})

	t.Run("two_elements_tie_broken_by_recency", func(t *testing.T) {
		fs := ConstructorFreqStack()
		fs.Push(1)
		fs.Push(2)
		// both freq=1; 2 was pushed more recently
		assert.Equal(t, 2, fs.Pop())
		assert.Equal(t, 1, fs.Pop())
	})

	t.Run("push_after_pop_restores_frequency", func(t *testing.T) {
		fs := ConstructorFreqStack()
		fs.Push(1)
		fs.Push(1)
		fs.Pop()       // removes one 1, freq[1]=1
		fs.Push(2)     // freq[2]=1
		fs.Push(2)     // freq[2]=2, now most frequent
		assert.Equal(t, 2, fs.Pop())
		assert.Equal(t, 2, fs.Pop()) // tie at freq=1, 2 pushed more recently
	})
}
