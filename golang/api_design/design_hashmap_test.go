package apidesign

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MyHashSet(t *testing.T) {
	t.Run("contains returns false on empty set", func(t *testing.T) {
		s := NewMyHashSet()
		assert.False(t, s.Contains(1))
	})

	t.Run("add then contains returns true", func(t *testing.T) {
		s := NewMyHashSet()
		s.Add(42)
		assert.True(t, s.Contains(42))
	})

	t.Run("remove existing key", func(t *testing.T) {
		s := NewMyHashSet()
		s.Add(7)
		s.Remove(7)
		assert.False(t, s.Contains(7))
	})

	t.Run("remove non-existing key is a no-op", func(t *testing.T) {
		s := NewMyHashSet()
		s.Remove(99) // should not panic
		assert.False(t, s.Contains(99))
	})

	t.Run("add duplicate is idempotent", func(t *testing.T) {
		s := NewMyHashSet()
		s.Add(5)
		s.Add(5)
		s.Remove(5)
		assert.False(t, s.Contains(5))
	})

	t.Run("hash collision keys coexist", func(t *testing.T) {
		s := NewMyHashSet()
		s.Add(0)
		s.Add(1009) // same bucket as 0
		assert.True(t, s.Contains(0))
		assert.True(t, s.Contains(1009))
		s.Remove(0)
		assert.False(t, s.Contains(0))
		assert.True(t, s.Contains(1009))
	})

	t.Run("multiple keys independent", func(t *testing.T) {
		s := NewMyHashSet()
		for i := 0; i < 10; i++ {
			s.Add(i)
		}
		for i := 0; i < 10; i++ {
			assert.True(t, s.Contains(i))
		}
		s.Remove(5)
		assert.False(t, s.Contains(5))
		assert.True(t, s.Contains(4))
		assert.True(t, s.Contains(6))
	})
}

func Test_MyHashMap(t *testing.T) {
	t.Run("get on empty map returns -1", func(t *testing.T) {
		m := NewMyHashMap()
		assert.Equal(t, -1, m.Get(1))
	})

	t.Run("put then get returns value", func(t *testing.T) {
		m := NewMyHashMap()
		m.Put(1, 100)
		assert.Equal(t, 100, m.Get(1))
	})

	t.Run("update existing key changes value", func(t *testing.T) {
		m := NewMyHashMap()
		m.Put(1, 10)
		m.Put(1, 20)
		assert.Equal(t, 20, m.Get(1))
	})

	t.Run("remove existing key", func(t *testing.T) {
		m := NewMyHashMap()
		m.Put(2, 42)
		m.Remove(2)
		assert.Equal(t, -1, m.Get(2))
	})

	t.Run("remove non-existing key is a no-op", func(t *testing.T) {
		m := NewMyHashMap()
		m.Remove(99) // should not panic
		assert.Equal(t, -1, m.Get(99))
	})

	t.Run("hash collision keys coexist", func(t *testing.T) {
		m := NewMyHashMap()
		m.Put(0, 1)
		m.Put(1009, 2) // same bucket as 0
		assert.Equal(t, 1, m.Get(0))
		assert.Equal(t, 2, m.Get(1009))
		m.Remove(0)
		assert.Equal(t, -1, m.Get(0))
		assert.Equal(t, 2, m.Get(1009))
	})

	t.Run("multiple independent keys", func(t *testing.T) {
		m := NewMyHashMap()
		for i := 0; i < 10; i++ {
			m.Put(i, i*10)
		}
		for i := 0; i < 10; i++ {
			assert.Equal(t, i*10, m.Get(i))
		}
	})
}
