package containers

/**
 * LRUCache[K, V] — Least Recently Used Cache
 *
 * Evicts the least recently used entry when the cache is full.
 * Backed by a doubly-linked list and a hash map, giving O(1) for all ops.
 *
 * Structural invariant:
 *   The list is ordered Most Recently Used (MRU) → Least Recently Used (LRU):
 *
 *     head ↔ [MRU node] ↔ ... ↔ [LRU node] ↔ tail
 *
 *   Every Get and Put that touches a key must promote its node to head.next.
 *   The eviction victim is always tail.prev.
 *
 * Implementation notes:
 *   - Key is stored on each node so eviction can call delete(map, victim.key)
 *     without an O(n) reverse scan of the map.
 *   - Sentinel head and tail nodes are allocated at construction and never
 *     hold real data. They make removeNode and insertFront branch-free:
 *     n.prev and n.next are always valid pointers (at worst a sentinel).
 *   - K must be comparable (usable as a Go map key). No ordering is needed,
 *     so cmp.Ordered would be too restrictive.
 *   - The map is pre-allocated at capacity to avoid rehashing during the
 *     initial fill.
 *
 * Complexity:
 *   Get, Put, Contains, Len, Cap   O(1)
 */

// lruNode is a single entry in the doubly-linked list.
type lruNode[K comparable, V any] struct {
	key        K
	val        V
	prev, next *lruNode[K, V]
}

// LRUCache is a generic Least Recently Used cache with O(1) Get and Put.
// K must be comparable (usable as a Go map key). V may be any type.
// The zero value is not usable; use NewLRUCache.
type LRUCache[K comparable, V any] struct {
	cap   int
	items map[K]*lruNode[K, V]
	head  *lruNode[K, V] // sentinel: MRU end (head.next is the hottest entry)
	tail  *lruNode[K, V] // sentinel: LRU end (tail.prev is the coldest entry)
}

// NewLRUCache returns an empty LRUCache with the given capacity.
// Panics if capacity < 1.
func NewLRUCache[K comparable, V any](capacity int) *LRUCache[K, V] {
	if capacity < 1 {
		panic("lrucache: capacity must be >= 1")
	}
	head := &lruNode[K, V]{}
	tail := &lruNode[K, V]{}
	head.next = tail
	tail.prev = head
	return &LRUCache[K, V]{
		cap:   capacity,
		items: make(map[K]*lruNode[K, V], capacity),
		head:  head,
		tail:  tail,
	}
}

// Get returns the value associated with k and true if k is present.
// A cache hit promotes k to the most recently used position.
// Returns (zero, false) on a cache miss.
// Time: O(1).
func (c *LRUCache[K, V]) Get(k K) (V, bool) {
	node, ok := c.items[k]
	if !ok {
		var zero V
		return zero, false
	}
	c.moveToFront(node)
	return node.val, true
}

// Put inserts or updates the value for key k.
// If k is already present, its value is updated and it is promoted to MRU.
// If the cache is at capacity and k is new, the LRU entry is evicted first.
// Time: O(1).
func (c *LRUCache[K, V]) Put(k K, v V) {
	if node, ok := c.items[k]; ok {
		node.val = v
		c.moveToFront(node)
		return
	}
	if len(c.items) == c.cap {
		victim := c.tail.prev
		c.removeNode(victim)
		delete(c.items, victim.key)
	}
	node := &lruNode[K, V]{key: k, val: v}
	c.insertFront(node)
	c.items[k] = node
}

// Contains reports whether k is present in the cache.
// Does NOT affect recency order.
// Time: O(1).
func (c *LRUCache[K, V]) Contains(k K) bool {
	_, ok := c.items[k]
	return ok
}

// Len returns the number of entries currently in the cache.
// Time: O(1).
func (c *LRUCache[K, V]) Len() int {
	return len(c.items)
}

// Cap returns the maximum number of entries the cache will hold.
// Time: O(1).
func (c *LRUCache[K, V]) Cap() int {
	return c.cap
}

// ---------------------------------------------------------------------------
// Internal list helpers — made branch-free by the sentinel head/tail nodes.
// ---------------------------------------------------------------------------

// removeNode splices n out of the doubly-linked list.
func (c *LRUCache[K, V]) removeNode(n *lruNode[K, V]) {
	n.prev.next = n.next
	n.next.prev = n.prev
}

// insertFront inserts n immediately after the head sentinel (MRU position).
func (c *LRUCache[K, V]) insertFront(n *lruNode[K, V]) {
	n.next = c.head.next
	n.prev = c.head
	c.head.next.prev = n
	c.head.next = n
}

// moveToFront removes n from its current position and re-inserts it at front.
func (c *LRUCache[K, V]) moveToFront(n *lruNode[K, V]) {
	c.removeNode(n)
	c.insertFront(n)
}
