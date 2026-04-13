package containers

import (
	"cmp"
	"iter"
)

// All returns an iterator over all key-value pairs in ascending key order.
// It uses the Go 1.23 range-over-func protocol (iter.Seq2), so it integrates
// directly with for-range and supports early exit:
//
//	for k, v := range m.All() { ... }
//
// The traversal is iterative using the parent pointer, so it allocates no
// auxiliary stack. The yield return value is respected: if the caller breaks
// out of the range loop, iteration stops immediately.
//
// Time: O(n) total; O(1) amortized per step.
func (m *TreeMap[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if m.root == nil {
			return
		}
		// Start at the leftmost node (the minimum).
		n := minimum(m.root)
		for n != nil {
			if !yield(n.key, n.val) {
				return // caller broke out of the range loop
			}
			n = successor(n)
		}
	}
}

// successor returns the in-order successor of n using the parent pointer.
// Returns nil when n is the maximum node.
func successor[K cmp.Ordered, V any](n *rbNode[K, V]) *rbNode[K, V] {
	if n.right != nil {
		// Successor is the minimum of the right subtree.
		return minimum(n.right)
	}
	// Walk up until we come from a left child.
	for n.parent != nil && n == n.parent.right {
		n = n.parent
	}
	return n.parent
}
