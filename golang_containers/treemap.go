package containers

import "cmp"

/**
 * TreeMap[K, V] — Ordered Map backed by a Red-Black Tree
 *
 * A generic key-value map that keeps keys in sorted order at all times.
 * Backed by a Red-Black tree (Bayer 1972 / Guibas & Sedgewick 1978),
 * which is a self-balancing BST guaranteeing O(log n) for all mutations
 * and lookups by maintaining five structural invariants:
 *
 *   1. Every node is red or black.
 *   2. The root is black.
 *   3. Every nil leaf is treated as black.
 *   4. A red node's children are both black (no consecutive reds).
 *   5. Every root→nil path passes through the same number of black nodes.
 *
 * Implementation notes:
 *   - Nodes carry a parent pointer, which simplifies deleteFixup (avoids
 *     threading an ancestor stack through the deletion path).
 *   - nil is treated as black everywhere (isRed(nil) == false), the standard
 *     Go alternative to CLRS's sentinel node — no extra allocation needed.
 *   - All key comparisons use cmp.Compare, which is generic-safe and returns
 *     a canonical -1/0/+1, making each branch direction visually obvious.
 *   - Keys must satisfy cmp.Ordered (all numeric types and string).
 *
 * Complexity:
 *   Put, Get, Delete, Contains, Floor, Ceiling, Min, Max  O(log n)
 *   Keys, Values, All                                     O(n)
 *   Size, IsEmpty                                         O(1)
 */

// TreeMap is a generic ordered map backed by a Red-Black tree.
// Keys must satisfy cmp.Ordered. The zero value is not usable; use NewTreeMap.
type TreeMap[K cmp.Ordered, V any] struct {
	root *rbNode[K, V]
	size int
}

// NewTreeMap returns an empty, ready-to-use TreeMap.
func NewTreeMap[K cmp.Ordered, V any]() *TreeMap[K, V] {
	return &TreeMap[K, V]{}
}

// Put inserts or updates the value for key k.
// If k already exists its value is replaced with v.
// Time: O(log n).
func (m *TreeMap[K, V]) Put(k K, v V) {
	// Standard BST insert — find the parent of the new leaf.
	var parent *rbNode[K, V]
	cur := m.root
	for cur != nil {
		parent = cur
		switch cmp.Compare(k, cur.key) {
		case -1:
			cur = cur.left
		case 1:
			cur = cur.right
		default:
			cur.val = v // key exists — update in place, no size change
			return
		}
	}

	z := &rbNode[K, V]{key: k, val: v, color: red, parent: parent}
	if parent == nil {
		m.root = z
	} else if cmp.Compare(k, parent.key) < 0 {
		parent.left = z
	} else {
		parent.right = z
	}
	m.size++
	m.insertFixup(z)
}

// Get returns the value associated with k.
// Returns (value, true) if found, or (zero, false) if not present.
// Time: O(log n).
func (m *TreeMap[K, V]) Get(k K) (V, bool) {
	n := m.search(k)
	if n == nil {
		var zero V
		return zero, false
	}
	return n.val, true
}

// Delete removes k and its value from the map.
// If k is not present, Delete is a no-op.
// Time: O(log n).
func (m *TreeMap[K, V]) Delete(k K) {
	z := m.search(k)
	if z == nil {
		return
	}
	m.delete(z)
	m.size--
}

// Contains reports whether k is present in the map.
// Time: O(log n).
func (m *TreeMap[K, V]) Contains(k K) bool {
	return m.search(k) != nil
}

// Size returns the number of key-value pairs in the map.
// Time: O(1).
func (m *TreeMap[K, V]) Size() int {
	return m.size
}

// IsEmpty reports whether the map has no entries.
func (m *TreeMap[K, V]) IsEmpty() bool {
	return m.size == 0
}

// Min returns the smallest key and its value.
// Returns (zero, zero, false) if the map is empty.
// Time: O(log n).
func (m *TreeMap[K, V]) Min() (K, V, bool) {
	if m.root == nil {
		var zk K
		var zv V
		return zk, zv, false
	}
	n := minimum(m.root)
	return n.key, n.val, true
}

// Max returns the largest key and its value.
// Returns (zero, zero, false) if the map is empty.
// Time: O(log n).
func (m *TreeMap[K, V]) Max() (K, V, bool) {
	if m.root == nil {
		var zk K
		var zv V
		return zk, zv, false
	}
	n := maximum(m.root)
	return n.key, n.val, true
}

// Floor returns the largest key less than or equal to k, with its value.
// Returns (zero, zero, false) if no such key exists.
// Time: O(log n).
func (m *TreeMap[K, V]) Floor(k K) (K, V, bool) {
	n := floorNode(m.root, k)
	if n == nil {
		var zk K
		var zv V
		return zk, zv, false
	}
	return n.key, n.val, true
}

// Ceiling returns the smallest key greater than or equal to k, with its value.
// Returns (zero, zero, false) if no such key exists.
// Time: O(log n).
func (m *TreeMap[K, V]) Ceiling(k K) (K, V, bool) {
	n := ceilingNode(m.root, k)
	if n == nil {
		var zk K
		var zv V
		return zk, zv, false
	}
	return n.key, n.val, true
}

// Keys returns all keys in ascending order.
// Time: O(n).
func (m *TreeMap[K, V]) Keys() []K {
	out := make([]K, 0, m.size)
	m.inorderKeys(m.root, &out)
	return out
}

// Values returns all values in key-ascending order.
// Time: O(n).
func (m *TreeMap[K, V]) Values() []V {
	out := make([]V, 0, m.size)
	m.inorderValues(m.root, &out)
	return out
}

// ---------------------------------------------------------------------------
// Unexported helpers
// ---------------------------------------------------------------------------

// search returns the node with key k, or nil if not found.
func (m *TreeMap[K, V]) search(k K) *rbNode[K, V] {
	n := m.root
	for n != nil {
		switch cmp.Compare(k, n.key) {
		case -1:
			n = n.left
		case 1:
			n = n.right
		default:
			return n
		}
	}
	return nil
}

// delete removes node z from the tree, maintaining RB invariants.
// Implements CLRS §13.4 RB-DELETE.
func (m *TreeMap[K, V]) delete(z *rbNode[K, V]) {
	// y is the node that will actually be spliced out.
	// yOrigColor tracks whether a black node is removed (triggering fixup).
	y := z
	yOrigColor := y.color

	var x *rbNode[K, V]    // x replaces y's position; may be nil
	var xParent *rbNode[K, V] // x's parent (needed when x is nil)

	if z.left == nil {
		// Case A: no left child — replace z with its right child.
		x = z.right
		xParent = z.parent
		m.transplant(z, z.right)
	} else if z.right == nil {
		// Case B: no right child — replace z with its left child.
		x = z.left
		xParent = z.parent
		m.transplant(z, z.left)
	} else {
		// Case C: two children — replace z with its in-order successor y.
		y = minimum(z.right)
		yOrigColor = y.color
		x = y.right
		if y.parent == z {
			xParent = y
		} else {
			xParent = y.parent
			m.transplant(y, y.right)
			y.right = z.right
			y.right.parent = y
		}
		m.transplant(z, y)
		y.left = z.left
		y.left.parent = y
		y.color = z.color
	}

	if yOrigColor == black {
		m.deleteFixup(x, xParent)
	}
}

func (m *TreeMap[K, V]) inorderKeys(n *rbNode[K, V], out *[]K) {
	if n == nil {
		return
	}
	m.inorderKeys(n.left, out)
	*out = append(*out, n.key)
	m.inorderKeys(n.right, out)
}

func (m *TreeMap[K, V]) inorderValues(n *rbNode[K, V], out *[]V) {
	if n == nil {
		return
	}
	m.inorderValues(n.left, out)
	*out = append(*out, n.val)
	m.inorderValues(n.right, out)
}
