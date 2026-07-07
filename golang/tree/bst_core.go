package tree

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"strings"

	. "gocode/types"
)

/**
 * BST — Binary Search Tree
 *
 * A BST is a binary tree satisfying the invariant:
 *   left.Val < node.Val < right.Val   (strict; no duplicates)
 *
 * This implementation uses the exact same node type as LeetCode:
 *
 *   type TreeNode struct {
 *       Val   int
 *       Left  *TreeNode
 *       Right *TreeNode
 *   }
 *
 * All operations are implemented as functions over *TreeNode and wrapped by
 * a thin BST struct. Callers never touch the root pointer directly — this
 * keeps the invariants intact across Insert and Delete.
 *
 * Complexity (h = height; O(log n) avg, O(n) worst on skewed trees):
 *   Insert, Delete, Search          O(h) time, O(h) stack
 *   Min, Max                        O(h) time, O(1) space
 *   Floor, Ceiling, Predecessor,
 *     Successor                     O(h) time, O(1) space
 *   Inorder, Preorder, Postorder    O(n) time, O(h) space
 *   LevelOrder                      O(n) time, O(n) space
 *   Height, Size                    O(n) time, O(h) space
 */

// LevelOrderNull is the sentinel used in NewBSTFromLevelOrder to represent an
// absent ("null") node — identical to the "null" token in LeetCode's JSON input.
// It is set to math.MinInt so it never collides with real node values in practice.
const LevelOrderNull = math.MinInt

// BST is a Binary Search Tree backed by *TreeNode (the LeetCode node type).
// The zero value is an empty, ready-to-use tree.
type BST struct {
	root *TreeNode
}

// ---------------------------------------------------------------------------
// Constructors
// ---------------------------------------------------------------------------

// NewBST returns an empty BST.
func NewBST() *BST { return &BST{} }

// NewBSTFromSorted builds a height-balanced BST from a sorted slice of
// distinct integers. Reuses the existing sortedArrayToBST implementation.
//
//   NewBSTFromSorted([]int{1, 2, 3, 4, 5, 6, 7})
//       4
//      / \
//     2   6
//    / \ / \
//   1  3 5  7
func NewBSTFromSorted(vals []int) *BST {
	return &BST{root: sortedArrayToBST(vals)}
}

// NewBSTFromLevelOrder builds a binary tree from a LeetCode-style level-order
// slice. Use LevelOrderNull (math.MinInt) to represent absent nodes.
//
//   NewBSTFromLevelOrder([]int{4, 2, 7, 1, 3})
//       4
//      / \
//     2   7
//    / \
//   1   3
//
// Note: this function builds any binary tree, not only valid BSTs.
// Call IsValid() afterward to verify the BST property if needed.
func NewBSTFromLevelOrder(vals []int) *BST {
	if len(vals) == 0 || vals[0] == LevelOrderNull {
		return &BST{}
	}
	root := &TreeNode{Val: vals[0]}
	queue := []*TreeNode{root}
	i := 1
	for len(queue) > 0 && i < len(vals) {
		n := queue[0]
		queue = queue[1:]

		if i < len(vals) {
			if vals[i] != LevelOrderNull {
				n.Left = &TreeNode{Val: vals[i]}
				queue = append(queue, n.Left)
			}
			i++
		}
		if i < len(vals) {
			if vals[i] != LevelOrderNull {
				n.Right = &TreeNode{Val: vals[i]}
				queue = append(queue, n.Right)
			}
			i++
		}
	}
	return &BST{root: root}
}

// Root returns the root node (read-only view; do not modify through this pointer).
func (b *BST) Root() *TreeNode { return b.root }

// ---------------------------------------------------------------------------
// Insert
// ---------------------------------------------------------------------------

// Insert adds val to the BST. Duplicates are silently ignored (BST invariant
// requires strict inequality; equal values are not stored twice).
func (b *BST) Insert(val int) { b.root = bstInsert(b.root, val) }

// bstInsert returns the new root after inserting val.
func bstInsert(root *TreeNode, val int) *TreeNode {
	if root == nil {
		return &TreeNode{Val: val}
	}
	switch {
	case val < root.Val:
		root.Left = bstInsert(root.Left, val)
	case val > root.Val:
		root.Right = bstInsert(root.Right, val)
	// val == root.Val: duplicate — no-op
	}
	return root
}

// ---------------------------------------------------------------------------
// Delete
// ---------------------------------------------------------------------------

// Delete removes val from the BST. If val is not present, the tree is unchanged.
//
// Three cases:
//   1. Leaf (no children)   → remove and return nil.
//   2. One child            → splice the node out, return its child.
//   3. Two children         → replace node's value with its inorder successor
//                             (minimum of right subtree), then delete that
//                             successor from the right subtree.
//
// Only the Val field of the node is replaced in case 3 — existing *TreeNode
// pointers into the tree remain valid at their original addresses.
func (b *BST) Delete(val int) { b.root = bstDelete(b.root, val) }

// bstDelete returns the new root after deleting val.
func bstDelete(root *TreeNode, val int) *TreeNode {
	if root == nil {
		return nil
	}
	switch {
	case val < root.Val:
		root.Left = bstDelete(root.Left, val)
	case val > root.Val:
		root.Right = bstDelete(root.Right, val)
	default:
		// Found the node to remove.
		if root.Left == nil {
			return root.Right // covers leaf (both nil) and right-only child
		}
		if root.Right == nil {
			return root.Left // left-only child
		}
		// Two children: promote inorder successor.
		succ := bstMin(root.Right)
		root.Val = succ.Val
		root.Right = bstDelete(root.Right, succ.Val)
	}
	return root
}

// ---------------------------------------------------------------------------
// Search
// ---------------------------------------------------------------------------

// Search returns the node whose Val equals target, or nil if not found.
func (b *BST) Search(target int) *TreeNode { return bstSearch(b.root, target) }

// bstSearch is the iterative O(h) search. Avoids recursion overhead and
// cannot overflow the stack on skewed trees.
func bstSearch(root *TreeNode, target int) *TreeNode {
	for root != nil {
		switch {
		case target < root.Val:
			root = root.Left
		case target > root.Val:
			root = root.Right
		default:
			return root
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// Min / Max
// ---------------------------------------------------------------------------

// Min returns the node with the smallest value, or nil if the tree is empty.
func (b *BST) Min() *TreeNode { return bstMin(b.root) }

// Max returns the node with the largest value, or nil if the tree is empty.
func (b *BST) Max() *TreeNode { return bstMax(b.root) }

// bstMin walks left until it cannot go further.
func bstMin(root *TreeNode) *TreeNode {
	for root != nil && root.Left != nil {
		root = root.Left
	}
	return root
}

// bstMax walks right until it cannot go further.
func bstMax(root *TreeNode) *TreeNode {
	for root != nil && root.Right != nil {
		root = root.Right
	}
	return root
}

// ---------------------------------------------------------------------------
// Floor & Ceiling
// ---------------------------------------------------------------------------

// Floor returns the node with the largest value ≤ target, or nil if every
// node in the tree is strictly greater than target.
//
//   BST: 1 3 5 7 9
//   Floor(4) → node(3)
//   Floor(5) → node(5)   (exact match)
//   Floor(0) → nil       (nothing ≤ 0)
func (b *BST) Floor(target int) *TreeNode { return bstFloor(b.root, target) }

func bstFloor(root *TreeNode, target int) *TreeNode {
	if root == nil {
		return nil
	}
	switch {
	case root.Val == target:
		return root // exact match
	case root.Val > target:
		return bstFloor(root.Left, target) // floor must be smaller
	default:
		// root.Val < target: floor is root OR something larger in the right subtree.
		if right := bstFloor(root.Right, target); right != nil {
			return right
		}
		return root
	}
}

// Ceiling returns the node with the smallest value ≥ target, or nil if every
// node in the tree is strictly less than target.
//
//   BST: 1 3 5 7 9
//   Ceiling(4) → node(5)
//   Ceiling(5) → node(5)   (exact match)
//   Ceiling(10) → nil      (nothing ≥ 10)
func (b *BST) Ceiling(target int) *TreeNode { return bstCeiling(b.root, target) }

func bstCeiling(root *TreeNode, target int) *TreeNode {
	if root == nil {
		return nil
	}
	switch {
	case root.Val == target:
		return root // exact match
	case root.Val < target:
		return bstCeiling(root.Right, target) // ceiling must be larger
	default:
		// root.Val > target: ceiling is root OR something smaller in the left subtree.
		if left := bstCeiling(root.Left, target); left != nil {
			return left
		}
		return root
	}
}

// ---------------------------------------------------------------------------
// Predecessor & Successor
// ---------------------------------------------------------------------------

// Predecessor returns the node with the largest value strictly less than
// target, or nil if target is ≤ the minimum value in the tree.
//
// Iterative: record the last node where we went right; that is the best
// candidate seen so far that is still less than target.
func (b *BST) Predecessor(target int) *TreeNode { return bstPredecessor(b.root, target) }

func bstPredecessor(root *TreeNode, target int) *TreeNode {
	var res *TreeNode
	for root != nil {
		if root.Val < target {
			res = root    // best candidate so far
			root = root.Right
		} else {
			root = root.Left
		}
	}
	return res
}

// Successor returns the node with the smallest value strictly greater than
// target, or nil if target is ≥ the maximum value in the tree.
//
// Iterative: record the last node where we went left; that is the best
// candidate seen so far that is still greater than target.
func (b *BST) Successor(target int) *TreeNode { return bstSuccessor(b.root, target) }

func bstSuccessor(root *TreeNode, target int) *TreeNode {
	var res *TreeNode
	for root != nil {
		if root.Val > target {
			res = root    // best candidate so far
			root = root.Left
		} else {
			root = root.Right
		}
	}
	return res
}

// ---------------------------------------------------------------------------
// Traversals (iterative)
//
// The iterative forms are preferred in interviews: they cannot stack-overflow
// on skewed trees and make the order of operations explicit.
// The package already contains iterative inorder, preorder, and postorder in
// tree_traversal.go — the BST methods below delegate to those.
// ---------------------------------------------------------------------------

// Inorder returns node values in ascending sorted order (left → root → right).
func (b *BST) Inorder() []int { return inorderTraversal(b.root) }

// Preorder returns node values in root → left → right order.
func (b *BST) Preorder() []int { return preorderTraversal(b.root) }

// Postorder returns node values in left → right → root order.
func (b *BST) Postorder() []int { return postorderTraversal(b.root) }

// LevelOrder returns node values grouped by level (BFS order).
// Each inner slice contains the values of one level, left to right.
// Absent children are not represented — only present nodes appear.
func (b *BST) LevelOrder() [][]int { return bstLevelOrder(b.root) }

func bstLevelOrder(root *TreeNode) [][]int {
	if root == nil {
		return nil
	}
	var res [][]int
	queue := []*TreeNode{root}
	for len(queue) > 0 {
		level := make([]int, len(queue))
		next := queue[:0:0] // reuse backing array
		for i, n := range queue {
			level[i] = n.Val
			if n.Left != nil {
				next = append(next, n.Left)
			}
			if n.Right != nil {
				next = append(next, n.Right)
			}
		}
		res = append(res, level)
		queue = next
	}
	return res
}

// ---------------------------------------------------------------------------
// Height & Size
// ---------------------------------------------------------------------------

// Height returns the number of edges on the longest root-to-leaf path.
// An empty tree has height 0; a single-node tree has height 1.
func (b *BST) Height() int { return bstHeight(b.root) }

func bstHeight(root *TreeNode) int {
	if root == nil {
		return 0
	}
	return max(bstHeight(root.Left), bstHeight(root.Right)) + 1
}

// Size returns the total number of nodes in the tree.
func (b *BST) Size() int { return bstSize(b.root) }

func bstSize(root *TreeNode) int {
	if root == nil {
		return 0
	}
	return bstSize(root.Left) + bstSize(root.Right) + 1
}

// ---------------------------------------------------------------------------
// Validity
// ---------------------------------------------------------------------------

// IsValid reports whether the tree satisfies the BST property at every node.
// Delegates to the existing isValidBST implementation.
func (b *BST) IsValid() bool { return isValidBST(b.root) }

// IsEmpty reports whether the tree contains no nodes.
func (b *BST) IsEmpty() bool { return b.root == nil }

// ---------------------------------------------------------------------------
// PrintTree / SprintTree — top-down ASCII illustration
//
// The output matches the style used in LeetCode problem illustrations:
//
//    ┌── example (h = 3, 7 nodes) ──────────────────────┐
//    │       4                                           │
//    │      / \                                          │
//    │     2   6                                         │
//    │    / \ / \                                        │
//    │   1  3 5  7                                       │
//    └───────────────────────────────────────────────────┘
//
// Algorithm:
//   1. Compute tree height h.
//   2. Allocate a byte grid of (2h-1) rows × (2^h - 1) columns.
//   3. Recursively fill: for a node at depth d with column range [lo, hi],
//      place its value at col = (lo+hi)/2.  Place '/' at (2d+1, col-1) and
//      '\' at (2d+1, col+1).  Recurse left into [lo, col-1], right into
//      [col+1, hi].
//   4. Print rows, trimming trailing spaces.
//
// Best results: height ≤ 6, values ≤ 2 digits. For taller trees, output
// remains correct but grows exponentially wide (width = 2^h - 1 columns).
// ---------------------------------------------------------------------------

// SprintTree returns the top-down ASCII representation as a string.
// An empty tree returns "<empty>".
func (b *BST) SprintTree() string {
	return sprintBST(b.root)
}

// PrintTree prints the top-down ASCII representation to stdout.
func (b *BST) PrintTree() {
	fmt.Print(sprintBST(b.root))
}

// sprintBST builds and returns the formatted tree string.
func sprintBST(root *TreeNode) string {
	if root == nil {
		return "<empty>"
	}
	h := bstHeight(root)
	if h == 1 {
		return strconv.Itoa(root.Val) + "\n"
	}

	width := (1 << h) - 1 // 2^h − 1 columns
	numRows := 2*h - 1     // alternating node rows and slash rows

	// byte grid initialised to spaces
	grid := make([][]byte, numRows)
	for i := range grid {
		grid[i] = bytes.Repeat([]byte{' '}, width)
	}

	var fill func(n *TreeNode, depth, lo, hi int)
	fill = func(n *TreeNode, depth, lo, hi int) {
		if n == nil {
			return
		}
		nodeRow := depth * 2
		col := (lo + hi) / 2

		// Center the value string around col (left-biased for even-length strings).
		s := strconv.Itoa(n.Val)
		start := col - (len(s)-1)/2
		for i := 0; i < len(s); i++ {
			if idx := start + i; idx >= 0 && idx < width {
				grid[nodeRow][idx] = s[i]
			}
		}

		// Slash characters sit one column inside the parent on the slash row below.
		slashRow := nodeRow + 1
		if slashRow < numRows {
			if n.Left != nil && col-1 >= 0 {
				grid[slashRow][col-1] = '/'
				fill(n.Left, depth+1, lo, col-1)
			}
			if n.Right != nil && col+1 < width {
				grid[slashRow][col+1] = '\\'
				fill(n.Right, depth+1, col+1, hi)
			}
		}
	}

	fill(root, 0, 0, width-1)

	var sb strings.Builder
	for _, row := range grid {
		line := strings.TrimRight(string(row), " ")
		sb.WriteString(line)
		sb.WriteByte('\n')
	}
	return sb.String()
}
