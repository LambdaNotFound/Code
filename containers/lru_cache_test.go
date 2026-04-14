package containers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// checkLRULinks walks the doubly-linked list and cross-checks it against the
// map: every list node must be in the map, and the map and list lengths must
// match. It also verifies that each node's prev/next pointers are consistent.
func checkLRULinks[K comparable, V any](t *testing.T, c *LRUCache[K, V]) {
	t.Helper()
	count := 0
	n := c.head.next
	for n != c.tail {
		assert.NotNil(t, n.prev, "node.prev is nil")
		assert.NotNil(t, n.next, "node.next is nil")
		assert.Equal(t, n, n.prev.next, "n.prev.next != n (broken forward link)")
		assert.Equal(t, n, n.next.prev, "n.next.prev != n (broken backward link)")
		_, ok := c.items[n.key]
		assert.True(t, ok, "list node key %v not found in map", n.key)
		count++
		n = n.next
	}
	assert.Equal(t, len(c.items), count, "map and list length mismatch")
	assert.Equal(t, c.tail, c.tail.prev.next, "tail.prev.next != tail (broken tail link)")
}

// ---------------------------------------------------------------------------
// Group 1: Basic operations
// ---------------------------------------------------------------------------

func Test_LRUCache_EmptyGet(t *testing.T) {
	c := NewLRUCache[int, string](3)
	v, ok := c.Get(1)
	assert.False(t, ok)
	assert.Equal(t, "", v)
	assert.Equal(t, 0, c.Len())
	assert.Equal(t, 3, c.Cap())
}

func Test_LRUCache_PutGet(t *testing.T) {
	c := NewLRUCache[int, string](3)
	c.Put(1, "one")
	v, ok := c.Get(1)
	assert.True(t, ok)
	assert.Equal(t, "one", v)
	assert.Equal(t, 1, c.Len())
	checkLRULinks(t, c)
}

func Test_LRUCache_CacheMiss(t *testing.T) {
	c := NewLRUCache[int, string](3)
	c.Put(1, "one")
	v, ok := c.Get(99)
	assert.False(t, ok)
	assert.Equal(t, "", v)
}

func Test_LRUCache_Contains(t *testing.T) {
	c := NewLRUCache[string, int](2)
	c.Put("a", 1)
	assert.True(t, c.Contains("a"))
	assert.False(t, c.Contains("b"))
}

func Test_LRUCache_ContainsDoesNotPromote(t *testing.T) {
	// Fill [a, b, c] at cap=3. Contains(a) must not promote a.
	// Inserting d should still evict a (the LRU), not b.
	c := NewLRUCache[string, int](3)
	c.Put("a", 1)
	c.Put("b", 2)
	c.Put("c", 3)
	_ = c.Contains("a") // must not promote
	c.Put("d", 4)       // evicts a
	assert.False(t, c.Contains("a"))
	assert.True(t, c.Contains("b"))
	checkLRULinks(t, c)
}

func Test_LRUCache_LenCap(t *testing.T) {
	c := NewLRUCache[int, int](5)
	for i := 0; i < 5; i++ {
		c.Put(i, i)
		assert.Equal(t, i+1, c.Len())
		assert.Equal(t, 5, c.Cap())
	}
	checkLRULinks(t, c)
}

// ---------------------------------------------------------------------------
// Group 2: Recency promotion
// ---------------------------------------------------------------------------

func Test_LRUCache_GetPromotesToMRU(t *testing.T) {
	// Insert 1, 2, 3 → order MRU→LRU: [3, 2, 1]
	// Get(1) promotes 1   → order: [1, 3, 2]
	// Put(4) evicts 2 (LRU), not 1
	c := NewLRUCache[int, int](3)
	c.Put(1, 1)
	c.Put(2, 2)
	c.Put(3, 3)
	_, _ = c.Get(1) // promote 1
	c.Put(4, 4)     // evicts 2

	assert.True(t, c.Contains(1), "1 should still be present after promotion")
	assert.False(t, c.Contains(2), "2 should have been evicted as LRU")
	assert.True(t, c.Contains(3))
	assert.True(t, c.Contains(4))
	checkLRULinks(t, c)
}

func Test_LRUCache_PutUpdatePromotes(t *testing.T) {
	// Insert 1, 2, 3 → MRU→LRU: [3, 2, 1]
	// Put(1, new) promotes 1 → [1, 3, 2]
	// Put(4) evicts 2
	c := NewLRUCache[int, int](3)
	c.Put(1, 10)
	c.Put(2, 20)
	c.Put(3, 30)
	c.Put(1, 99) // update + promote
	c.Put(4, 40) // evicts 2

	assert.True(t, c.Contains(1))
	assert.False(t, c.Contains(2), "2 should have been evicted")
	v, ok := c.Get(1)
	assert.True(t, ok)
	assert.Equal(t, 99, v)
	checkLRULinks(t, c)
}

// ---------------------------------------------------------------------------
// Group 3: Eviction
// ---------------------------------------------------------------------------

func Test_LRUCache_EvictsLRUOnOverflow(t *testing.T) {
	c := NewLRUCache[int, int](3)
	c.Put(1, 1)
	c.Put(2, 2)
	c.Put(3, 3)
	c.Put(4, 4) // evicts 1
	assert.False(t, c.Contains(1), "1 should be evicted")
	assert.True(t, c.Contains(2))
	assert.True(t, c.Contains(3))
	assert.True(t, c.Contains(4))
	assert.Equal(t, 3, c.Len())
	checkLRULinks(t, c)
}

func Test_LRUCache_EvictionOrder(t *testing.T) {
	// Fill [1,2,3]. Insert 4 → evicts 1. Insert 5 → evicts 2.
	c := NewLRUCache[int, int](3)
	c.Put(1, 1)
	c.Put(2, 2)
	c.Put(3, 3)

	c.Put(4, 4)
	assert.False(t, c.Contains(1))
	checkLRULinks(t, c)

	c.Put(5, 5)
	assert.False(t, c.Contains(2))
	checkLRULinks(t, c)

	assert.True(t, c.Contains(3))
	assert.True(t, c.Contains(4))
	assert.True(t, c.Contains(5))
}

func Test_LRUCache_SequentialEviction(t *testing.T) {
	const cap = 4
	c := NewLRUCache[int, int](cap)
	const total = 20

	for i := 1; i <= total; i++ {
		c.Put(i, i)
		assert.Equal(t, min(i, cap), c.Len())

		if i > cap {
			evicted := i - cap
			assert.False(t, c.Contains(evicted), "key %d should have been evicted", evicted)
			for j := evicted + 1; j <= i; j++ {
				assert.True(t, c.Contains(j), "key %d should still be present", j)
			}
		}
		checkLRULinks(t, c)
	}
}

// ---------------------------------------------------------------------------
// Group 4: Update semantics
// ---------------------------------------------------------------------------

func Test_LRUCache_UpdateDoesNotGrowLen(t *testing.T) {
	c := NewLRUCache[int, int](3)
	c.Put(1, 10)
	c.Put(1, 20)
	assert.Equal(t, 1, c.Len())
}

func Test_LRUCache_UpdateChangesValue(t *testing.T) {
	c := NewLRUCache[int, int](3)
	c.Put(1, 10)
	c.Put(1, 99)
	v, ok := c.Get(1)
	assert.True(t, ok)
	assert.Equal(t, 99, v)
}

func Test_LRUCache_UpdateSavesFromEviction(t *testing.T) {
	// Insert A, B, C (cap=3). A is LRU.
	// Update A → A promoted to MRU, B becomes LRU.
	// Insert D → evicts B (not A).
	c := NewLRUCache[string, int](3)
	c.Put("A", 1)
	c.Put("B", 2)
	c.Put("C", 3)
	c.Put("A", 100) // update promotes A
	c.Put("D", 4)   // evicts B

	v, ok := c.Get("A")
	assert.True(t, ok)
	assert.Equal(t, 100, v)
	assert.False(t, c.Contains("B"), "B should have been evicted")
	checkLRULinks(t, c)
}

// ---------------------------------------------------------------------------
// Group 5: Edge cases
// ---------------------------------------------------------------------------

func Test_LRUCache_CapacityOne(t *testing.T) {
	c := NewLRUCache[int, int](1)
	c.Put(1, 1)
	assert.Equal(t, 1, c.Len())

	c.Put(2, 2) // evicts 1
	_, ok := c.Get(1)
	assert.False(t, ok)
	v, ok := c.Get(2)
	assert.True(t, ok)
	assert.Equal(t, 2, v)
	checkLRULinks(t, c)

	c.Put(3, 3) // evicts 2
	assert.False(t, c.Contains(2))
	assert.True(t, c.Contains(3))
	checkLRULinks(t, c)
}

func Test_LRUCache_CapacityOnePutSameKey(t *testing.T) {
	c := NewLRUCache[int, string](1)
	c.Put(1, "a")
	c.Put(1, "b")
	assert.Equal(t, 1, c.Len())
	v, ok := c.Get(1)
	assert.True(t, ok)
	assert.Equal(t, "b", v)
	checkLRULinks(t, c)
}

func Test_LRUCache_InvalidCapacityPanics(t *testing.T) {
	assert.Panics(t, func() { NewLRUCache[int, int](0) })
	assert.Panics(t, func() { NewLRUCache[int, int](-5) })
}

// ---------------------------------------------------------------------------
// Group 6: LeetCode 146 examples
// ---------------------------------------------------------------------------

func Test_LRUCache_LeetCode146_Example1(t *testing.T) {
	// capacity = 2
	// Put(1,1), Put(2,2), Get(1)→1, Put(3,3)[evicts 2],
	// Get(2)→miss, Put(4,4)[evicts 1], Get(1)→miss, Get(3)→3, Get(4)→4
	miss := -1
	getOrMiss := func(c *LRUCache[int, int], k int) int {
		if v, ok := c.Get(k); ok {
			return v
		}
		return miss
	}

	c := NewLRUCache[int, int](2)
	c.Put(1, 1)
	c.Put(2, 2)
	assert.Equal(t, 1, getOrMiss(c, 1))
	checkLRULinks(t, c)

	c.Put(3, 3) // evicts 2 (1 was just promoted to MRU)
	assert.Equal(t, miss, getOrMiss(c, 2))
	checkLRULinks(t, c)

	c.Put(4, 4) // evicts 1 (3 is MRU, 1 is LRU)
	assert.Equal(t, miss, getOrMiss(c, 1))
	assert.Equal(t, 3, getOrMiss(c, 3))
	assert.Equal(t, 4, getOrMiss(c, 4))
	checkLRULinks(t, c)
}

func Test_LRUCache_LeetCode146_Example2(t *testing.T) {
	// capacity = 1
	// Put(2,1), Get(2)→1, Put(3,2)[evicts 2], Get(2)→miss, Get(3)→2
	miss := -1
	getOrMiss := func(c *LRUCache[int, int], k int) int {
		if v, ok := c.Get(k); ok {
			return v
		}
		return miss
	}

	c := NewLRUCache[int, int](1)
	c.Put(2, 1)
	assert.Equal(t, 1, getOrMiss(c, 2))

	c.Put(3, 2) // evicts 2
	assert.Equal(t, miss, getOrMiss(c, 2))
	assert.Equal(t, 2, getOrMiss(c, 3))
	checkLRULinks(t, c)
}

// ---------------------------------------------------------------------------
// Group 7: String key type (generics)
// ---------------------------------------------------------------------------

func Test_LRUCache_StringKeys(t *testing.T) {
	c := NewLRUCache[string, int](2)
	c.Put("foo", 1)
	c.Put("bar", 2)

	v, ok := c.Get("foo") // promotes "foo"
	assert.True(t, ok)
	assert.Equal(t, 1, v)

	c.Put("baz", 3) // evicts "bar" (LRU)
	assert.False(t, c.Contains("bar"))
	assert.True(t, c.Contains("foo"))
	assert.True(t, c.Contains("baz"))
	checkLRULinks(t, c)
}
