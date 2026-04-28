package apidesign

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LRUCache(t *testing.T) {
    cache := ConstructorLRUCache(3)

    cache.Put(1, 1)
    cache.Put(2, 2)
    value := cache.Get(1)
    assert.Equal(t, 1, value)

    cache.Put(3, 3)
    value = cache.Get(2)
    assert.Equal(t, 2, value)

    cache.Put(4, 4)
    value = cache.Get(1)
    assert.Equal(t, -1, value)
    value = cache.Get(4)
    assert.Equal(t, 4, value)

    value = cache.Get(3)
    assert.Equal(t, 3, value)
}

// Test_LRUCachePutUpdate covers the Put branch where the key already exists.
// The node must be removed, its value updated, and re-inserted at the front
// (most-recently-used position) without consuming extra capacity.
func Test_LRUCachePutUpdate(t *testing.T) {
	t.Run("update_existing_key_changes_value", func(t *testing.T) {
		cache := ConstructorLRUCache(2)
		cache.Put(1, 10)
		cache.Put(2, 20)

		// Update key 1's value — should not evict anything (capacity still 2).
		cache.Put(1, 100)
		assert.Equal(t, 100, cache.Get(1)) // updated value
		assert.Equal(t, 20, cache.Get(2))  // key 2 still present
	})

	t.Run("update_moves_key_to_most_recently_used", func(t *testing.T) {
		cache := ConstructorLRUCache(2)
		cache.Put(1, 1)
		cache.Put(2, 2)
		// Key 1 was the LRU. Update it — now key 2 becomes LRU.
		cache.Put(1, 11)
		// Insert key 3: should evict key 2 (the new LRU), not key 1.
		cache.Put(3, 3)
		assert.Equal(t, 11, cache.Get(1)) // key 1 survived
		assert.Equal(t, -1, cache.Get(2)) // key 2 was evicted
		assert.Equal(t, 3, cache.Get(3))  // key 3 is present
	})
}

func Test_LRUCacheEviction(t *testing.T) {
    // identical eviction scenario using the concrete LRUCache API
    cache := ConstructorLRUCache(3)

    cache.Put(1, 1)
    cache.Put(2, 2)
    assert.Equal(t, 1, cache.Get(1)) // hit, refreshes key 1

    cache.Put(3, 3)
    assert.Equal(t, 2, cache.Get(2)) // hit, refreshes key 2

    cache.Put(4, 4)                   // evicts key 1 (LRU)
    assert.Equal(t, -1, cache.Get(1)) // miss
    assert.Equal(t, 4, cache.Get(4))  // hit
    assert.Equal(t, 3, cache.Get(3))  // hit
}
