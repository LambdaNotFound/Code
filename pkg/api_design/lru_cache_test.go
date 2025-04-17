package apidesign

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LRUCache(t *testing.T) {
    cache := Constructor(3)

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

func Test_LRUCacheWithList(t *testing.T) {
    cache := ConstructorWithList[int](3)

    cache.Put(1, 1)
    cache.Put(2, 2)
    value := cache.Get(1)

    assert.Equal(t, 1, value.(int))
    cache.Put(3, 3)
    value = cache.Get(2)
    assert.Equal(t, 2, value.(int))

    cache.Put(4, 4)
    value = cache.Get(1)
    assert.Equal(t, nil, value)
    value = cache.Get(4)
    assert.Equal(t, 4, value.(int))

    value = cache.Get(3)
    assert.Equal(t, 3, value.(int))
}
