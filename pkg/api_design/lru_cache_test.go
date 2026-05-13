package apidesign

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LRUCacheGet(t *testing.T) {
	t.Run("miss on empty cache returns -1", func(t *testing.T) {
		cache := ConstructorLRUCache(2)
		assert.Equal(t, -1, cache.Get(1))
	})

	t.Run("hit returns correct value", func(t *testing.T) {
		cache := ConstructorLRUCache(2)
		cache.Put(1, 10)
		assert.Equal(t, 10, cache.Get(1))
	})

	t.Run("miss on absent key returns -1", func(t *testing.T) {
		cache := ConstructorLRUCache(2)
		cache.Put(1, 10)
		assert.Equal(t, -1, cache.Get(99))
	})
}

func Test_LRUCachePut(t *testing.T) {
	t.Run("update existing key changes value without eviction", func(t *testing.T) {
		cache := ConstructorLRUCache(2)
		cache.Put(1, 10)
		cache.Put(2, 20)
		cache.Put(1, 100)
		assert.Equal(t, 100, cache.Get(1))
		assert.Equal(t, 20, cache.Get(2))
	})

	t.Run("update promotes key to most recently used", func(t *testing.T) {
		cache := ConstructorLRUCache(2)
		cache.Put(1, 1)
		cache.Put(2, 2)
		cache.Put(1, 11) // key 1 becomes MRU; key 2 is now LRU
		cache.Put(3, 3)  // evicts key 2
		assert.Equal(t, 11, cache.Get(1))
		assert.Equal(t, -1, cache.Get(2))
		assert.Equal(t, 3, cache.Get(3))
	})
}

func Test_LRUCacheEviction(t *testing.T) {
	t.Run("evicts least recently used on overflow", func(t *testing.T) {
		cache := ConstructorLRUCache(3)
		cache.Put(1, 1) // MRU→LRU: [1]
		cache.Put(2, 2) // [2, 1]
		cache.Get(1)    // [1, 2]
		cache.Put(3, 3) // [3, 1, 2]
		cache.Get(2)    // [2, 3, 1]
		cache.Put(4, 4) // evicts key 1 (LRU); [4, 2, 3]
		assert.Equal(t, -1, cache.Get(1))
		assert.Equal(t, 2, cache.Get(2))
		assert.Equal(t, 3, cache.Get(3))
		assert.Equal(t, 4, cache.Get(4))
	})

	t.Run("capacity 1 evicts previous key on each put", func(t *testing.T) {
		cache := ConstructorLRUCache(1)
		cache.Put(1, 1)
		cache.Put(2, 2)
		assert.Equal(t, -1, cache.Get(1))
		assert.Equal(t, 2, cache.Get(2))
	})

	t.Run("get refreshes key preventing its eviction", func(t *testing.T) {
		cache := ConstructorLRUCache(2)
		cache.Put(1, 1)
		cache.Put(2, 2)
		cache.Get(1)    // key 1 becomes MRU; key 2 is LRU
		cache.Put(3, 3) // evicts key 2, not key 1
		assert.Equal(t, 1, cache.Get(1))
		assert.Equal(t, -1, cache.Get(2))
		assert.Equal(t, 3, cache.Get(3))
	})
}
