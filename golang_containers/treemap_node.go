package containers

import "cmp"

// color is the Red-Black node color.
type color bool

const (
	red   color = true
	black color = false
)

// rbNode is a single node in the Red-Black tree.
type rbNode[K cmp.Ordered, V any] struct {
	key         K
	val         V
	color       color
	left, right *rbNode[K, V]
	parent      *rbNode[K, V]
}

// isRed reports whether n is a non-nil red node.
// Nil is treated as black — the standard RB-tree convention that eliminates
// nil checks in every rotation and fixup path.
func isRed[K cmp.Ordered, V any](n *rbNode[K, V]) bool {
	return n != nil && n.color == red
}

// minimum returns the leftmost (smallest-key) node in the subtree rooted at n.
func minimum[K cmp.Ordered, V any](n *rbNode[K, V]) *rbNode[K, V] {
	for n.left != nil {
		n = n.left
	}
	return n
}

// maximum returns the rightmost (largest-key) node in the subtree rooted at n.
func maximum[K cmp.Ordered, V any](n *rbNode[K, V]) *rbNode[K, V] {
	for n.right != nil {
		n = n.right
	}
	return n
}

// floorNode returns the node with the largest key ≤ k in the subtree rooted at n,
// or nil if no such node exists.
func floorNode[K cmp.Ordered, V any](n *rbNode[K, V], k K) *rbNode[K, V] {
	if n == nil {
		return nil
	}
	switch cmp.Compare(k, n.key) {
	case 0:
		return n
	case -1:
		return floorNode(n.left, k)
	default: // k > n.key: n is a candidate; a better one may exist to the right
		if better := floorNode(n.right, k); better != nil {
			return better
		}
		return n
	}
}

// ceilingNode returns the node with the smallest key ≥ k in the subtree rooted at n,
// or nil if no such node exists.
func ceilingNode[K cmp.Ordered, V any](n *rbNode[K, V], k K) *rbNode[K, V] {
	if n == nil {
		return nil
	}
	switch cmp.Compare(k, n.key) {
	case 0:
		return n
	case 1:
		return ceilingNode(n.right, k)
	default: // k < n.key: n is a candidate; a better one may exist to the left
		if better := ceilingNode(n.left, k); better != nil {
			return better
		}
		return n
	}
}

// ---------------------------------------------------------------------------
// Rotations
//
// rotateLeft pivots x down-left and y (x.right) up:
//
//       x                y
//      / \              / \
//     a   y    →       x   c
//        / \          / \
//       b   c        a   b
//
// rotateRight is the mirror image.
// Both operations preserve BST ordering and carefully rewire parent pointers.
// ---------------------------------------------------------------------------

func (m *TreeMap[K, V]) rotateLeft(x *rbNode[K, V]) {
	y := x.right
	x.right = y.left
	if y.left != nil {
		y.left.parent = x
	}
	y.parent = x.parent
	if x.parent == nil {
		m.root = y
	} else if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}
	y.left = x
	x.parent = y
}

func (m *TreeMap[K, V]) rotateRight(x *rbNode[K, V]) {
	y := x.left
	x.left = y.right
	if y.right != nil {
		y.right.parent = x
	}
	y.parent = x.parent
	if x.parent == nil {
		m.root = y
	} else if x == x.parent.right {
		x.parent.right = y
	} else {
		x.parent.left = y
	}
	y.right = x
	x.parent = y
}

// ---------------------------------------------------------------------------
// Insert fixup  (CLRS §13.3)
//
// After a plain BST insert of a red node z, at most one RB invariant is
// violated: invariant 4 (no consecutive reds).  insertFixup restores it by
// handling three uncle cases, looping upward until the violation is resolved.
// ---------------------------------------------------------------------------

func (m *TreeMap[K, V]) insertFixup(z *rbNode[K, V]) {
	for isRed(z.parent) {
		if z.parent == z.parent.parent.left {
			uncle := z.parent.parent.right
			if isRed(uncle) {
				// Case 1: uncle is red — recolor and move violation up.
				z.parent.color = black
				uncle.color = black
				z.parent.parent.color = red
				z = z.parent.parent
			} else {
				if z == z.parent.right {
					// Case 2: uncle is black, z is right child — rotate to convert to case 3.
					z = z.parent
					m.rotateLeft(z)
				}
				// Case 3: uncle is black, z is left child — rotate at grandparent.
				z.parent.color = black
				z.parent.parent.color = red
				m.rotateRight(z.parent.parent)
			}
		} else {
			// Mirror: parent is right child of grandparent.
			uncle := z.parent.parent.left
			if isRed(uncle) {
				z.parent.color = black
				uncle.color = black
				z.parent.parent.color = red
				z = z.parent.parent
			} else {
				if z == z.parent.left {
					z = z.parent
					m.rotateRight(z)
				}
				z.parent.color = black
				z.parent.parent.color = red
				m.rotateLeft(z.parent.parent)
			}
		}
	}
	m.root.color = black
}

// ---------------------------------------------------------------------------
// Transplant  (CLRS §13.4)
//
// transplant replaces the subtree rooted at u with the subtree rooted at v,
// updating u's parent to point to v.  It does NOT update v.left or v.right.
// ---------------------------------------------------------------------------

func (m *TreeMap[K, V]) transplant(u, v *rbNode[K, V]) {
	if u.parent == nil {
		m.root = v
	} else if u == u.parent.left {
		u.parent.left = v
	} else {
		u.parent.right = v
	}
	if v != nil {
		v.parent = u.parent
	}
}

// ---------------------------------------------------------------------------
// Delete fixup  (CLRS §13.4)
//
// deleteFixup is called when the node spliced out was black, creating a
// "double-black" deficit at x.  The loop propagates or resolves the deficit
// upward through four sibling cases until x reaches the root or becomes red.
// ---------------------------------------------------------------------------

func (m *TreeMap[K, V]) deleteFixup(x, xParent *rbNode[K, V]) {
	for x != m.root && !isRed(x) {
		if x == xParent.left {
			sibling := xParent.right
			if isRed(sibling) {
				// Case 1: sibling is red — recolor and rotate to convert to cases 2–4.
				sibling.color = black
				xParent.color = red
				m.rotateLeft(xParent)
				sibling = xParent.right
			}
			if !isRed(sibling.left) && !isRed(sibling.right) {
				// Case 2: sibling's children are both black — push deficit up.
				sibling.color = red
				x = xParent
				xParent = x.parent
			} else {
				if !isRed(sibling.right) {
					// Case 3: sibling's right child is black — rotate sibling right to get case 4.
					if sibling.left != nil {
						sibling.left.color = black
					}
					sibling.color = red
					m.rotateRight(sibling)
					sibling = xParent.right
				}
				// Case 4: sibling's right child is red — rotate left at xParent and recolor.
				sibling.color = xParent.color
				xParent.color = black
				if sibling.right != nil {
					sibling.right.color = black
				}
				m.rotateLeft(xParent)
				x = m.root
				xParent = nil
			}
		} else {
			// Mirror: x is right child.
			sibling := xParent.left
			if isRed(sibling) {
				sibling.color = black
				xParent.color = red
				m.rotateRight(xParent)
				sibling = xParent.left
			}
			if !isRed(sibling.right) && !isRed(sibling.left) {
				sibling.color = red
				x = xParent
				xParent = x.parent
			} else {
				if !isRed(sibling.left) {
					if sibling.right != nil {
						sibling.right.color = black
					}
					sibling.color = red
					m.rotateLeft(sibling)
					sibling = xParent.left
				}
				sibling.color = xParent.color
				xParent.color = black
				if sibling.left != nil {
					sibling.left.color = black
				}
				m.rotateRight(xParent)
				x = m.root
				xParent = nil
			}
		}
	}
	if x != nil {
		x.color = black
	}
}
