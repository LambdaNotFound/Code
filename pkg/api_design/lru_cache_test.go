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
