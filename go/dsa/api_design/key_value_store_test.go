package apidesign

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_KVStore_basic(t *testing.T) {
	t.Run("get missing key returns false", func(t *testing.T) {
		kv := NewKVStore()
		_, ok := kv.Get("x")
		assert.False(t, ok)
	})

	t.Run("set then get", func(t *testing.T) {
		kv := NewKVStore()
		kv.Set("x", "hello")
		v, ok := kv.Get("x")
		assert.True(t, ok)
		assert.Equal(t, "hello", v)
	})

	t.Run("delete writes sentinel empty string", func(t *testing.T) {
		kv := NewKVStore()
		kv.Set("x", "hello")
		kv.Delete("x")
		v, ok := kv.Get("x")
		assert.True(t, ok)
		assert.Equal(t, "", v)
	})

	t.Run("overwrite existing key", func(t *testing.T) {
		kv := NewKVStore()
		kv.Set("k", "first")
		kv.Set("k", "second")
		v, _ := kv.Get("k")
		assert.Equal(t, "second", v)
	})
}

func Test_KVStore_rollback(t *testing.T) {
	t.Run("rollback discards transaction changes", func(t *testing.T) {
		kv := NewKVStore()
		kv.Set("x", "base")
		kv.Begin()
		kv.Set("x", "txn")
		v, _ := kv.Get("x")
		assert.Equal(t, "txn", v)
		kv.Rollback()
		v, _ = kv.Get("x")
		assert.Equal(t, "base", v)
	})

	t.Run("rollback with no active transaction panics", func(t *testing.T) {
		kv := NewKVStore()
		assert.Panics(t, func() { kv.Rollback() })
	})
}

func Test_KVStore_commit(t *testing.T) {
	t.Run("commit merges transaction into base", func(t *testing.T) {
		kv := NewKVStore()
		kv.Set("x", "base")
		kv.Begin()
		kv.Set("x", "txn")
		kv.Commit()
		v, ok := kv.Get("x")
		assert.True(t, ok)
		assert.Equal(t, "txn", v)
		assert.Equal(t, 0, len(kv.layers))
	})

	t.Run("commit with no active transaction panics", func(t *testing.T) {
		kv := NewKVStore()
		assert.Panics(t, func() { kv.Commit() })
	})

	t.Run("transaction layer shadows base during txn", func(t *testing.T) {
		kv := NewKVStore()
		kv.Set("a", "1")
		kv.Begin()
		kv.Set("b", "2")
		a, _ := kv.Get("a")
		b, _ := kv.Get("b")
		assert.Equal(t, "1", a) // falls through to base
		assert.Equal(t, "2", b) // found in layer
	})
}

func Test_KVStore_nested(t *testing.T) {
	t.Run("nested commit merges inner into outer layer", func(t *testing.T) {
		kv := NewKVStore()
		kv.Set("a", "1")
		kv.Begin()
		kv.Set("a", "2")
		kv.Begin()
		kv.Set("a", "3")
		kv.Commit() // inner → outer layer
		v, _ := kv.Get("a")
		assert.Equal(t, "3", v)
		assert.Equal(t, 1, len(kv.layers))
		kv.Commit() // outer → base
		v, _ = kv.Get("a")
		assert.Equal(t, "3", v)
		assert.Equal(t, 0, len(kv.layers))
	})

	t.Run("rollback inner keeps outer layer", func(t *testing.T) {
		kv := NewKVStore()
		kv.Begin()
		kv.Set("x", "outer")
		kv.Begin()
		kv.Set("x", "inner")
		kv.Rollback()
		v, _ := kv.Get("x")
		assert.Equal(t, "outer", v)
		assert.Equal(t, 1, len(kv.layers))
	})
}
