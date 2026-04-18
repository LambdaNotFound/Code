package tree

import (
	"math"
	"strings"
	"testing"

	. "gocode/types"

	"github.com/stretchr/testify/assert"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// vals collects Inorder() values — used for most structural assertions because
// a correct BST always produces a sorted sequence on inorder traversal.
func vals(b *BST) []int { return b.Inorder() }

// nodeVal safely dereferences a *TreeNode result; panics if nil (test bug, not code bug).
func nodeVal(n *TreeNode) int {
	if n == nil {
		panic("unexpected nil TreeNode in test")
	}
	return n.Val
}

// ---------------------------------------------------------------------------
// Insert
// ---------------------------------------------------------------------------

func Test_BST_Insert(t *testing.T) {
	t.Run("empty tree", func(t *testing.T) {
		b := NewBST()
		b.Insert(5)
		assert.Equal(t, []int{5}, vals(b))
		assert.Equal(t, 1, b.Size())
	})

	t.Run("duplicate ignored", func(t *testing.T) {
		b := NewBST()
		b.Insert(5)
		b.Insert(5)
		assert.Equal(t, []int{5}, vals(b))
		assert.Equal(t, 1, b.Size())
	})

	t.Run("left spine (decreasing inserts)", func(t *testing.T) {
		b := NewBST()
		for _, v := range []int{5, 4, 3, 2, 1} {
			b.Insert(v)
		}
		assert.Equal(t, []int{1, 2, 3, 4, 5}, vals(b))
		assert.True(t, b.IsValid())
	})

	t.Run("right spine (increasing inserts)", func(t *testing.T) {
		b := NewBST()
		for _, v := range []int{1, 2, 3, 4, 5} {
			b.Insert(v)
		}
		assert.Equal(t, []int{1, 2, 3, 4, 5}, vals(b))
		assert.True(t, b.IsValid())
	})

	t.Run("arbitrary order stays valid BST", func(t *testing.T) {
		b := NewBST()
		for _, v := range []int{5, 3, 7, 1, 4, 6, 8} {
			b.Insert(v)
		}
		assert.Equal(t, []int{1, 3, 4, 5, 6, 7, 8}, vals(b))
		assert.True(t, b.IsValid())
	})
}

// ---------------------------------------------------------------------------
// Delete
// ---------------------------------------------------------------------------

func Test_BST_Delete(t *testing.T) {
	// Shared helper: build the same 7-node tree for each sub-test.
	//       4
	//      / \
	//     2   6
	//    / \ / \
	//   1  3 5  7
	build := func() *BST {
		b := NewBST()
		for _, v := range []int{4, 2, 6, 1, 3, 5, 7} {
			b.Insert(v)
		}
		return b
	}

	t.Run("delete leaf — leftmost", func(t *testing.T) {
		b := build()
		b.Delete(1)
		assert.Equal(t, []int{2, 3, 4, 5, 6, 7}, vals(b))
		assert.True(t, b.IsValid())
	})

	t.Run("delete leaf — rightmost", func(t *testing.T) {
		b := build()
		b.Delete(7)
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, vals(b))
		assert.True(t, b.IsValid())
	})

	t.Run("delete leaf — interior leaf", func(t *testing.T) {
		b := build()
		b.Delete(3)
		assert.Equal(t, []int{1, 2, 4, 5, 6, 7}, vals(b))
		assert.True(t, b.IsValid())
	})

	t.Run("delete node with only left child", func(t *testing.T) {
		// Build tree where node 6 has only a left child (5).
		//       4
		//      / \
		//     2   6
		//    / \ /
		//   1  3 5
		b := NewBST()
		for _, v := range []int{4, 2, 6, 1, 3, 5} {
			b.Insert(v)
		}
		b.Delete(6)
		assert.Equal(t, []int{1, 2, 3, 4, 5}, vals(b))
		assert.True(t, b.IsValid())
	})

	t.Run("delete node with only right child", func(t *testing.T) {
		// Build tree where node 2 has only a right child (3).
		//       4
		//      / \
		//     2   6
		//      \ / \
		//      3 5  7
		b := NewBST()
		for _, v := range []int{4, 2, 6, 3, 5, 7} {
			b.Insert(v)
		}
		b.Delete(2)
		assert.Equal(t, []int{3, 4, 5, 6, 7}, vals(b))
		assert.True(t, b.IsValid())
	})

	t.Run("delete node with two children — interior", func(t *testing.T) {
		b := build()
		b.Delete(2) // successor of 2 is 3; 3 is promoted, then 3 is removed from right
		assert.Equal(t, []int{1, 3, 4, 5, 6, 7}, vals(b))
		assert.True(t, b.IsValid())
	})

	t.Run("delete root with two children", func(t *testing.T) {
		b := build()
		b.Delete(4) // successor of 4 is 5; 4→5, remove old 5
		assert.Equal(t, []int{1, 2, 3, 5, 6, 7}, vals(b))
		assert.True(t, b.IsValid())
	})

	t.Run("delete successor that itself has a right child", func(t *testing.T) {
		// Tree: 10 → left:5, right:20 → left:15 → right:17
		//   10
		//  /  \
		// 5   20
		//    /
		//   15
		//     \
		//     17
		b := NewBST()
		for _, v := range []int{10, 5, 20, 15, 17} {
			b.Insert(v)
		}
		b.Delete(10) // successor is 15; 15 has right child 17 which must be preserved
		assert.Equal(t, []int{5, 15, 17, 20}, vals(b))
		assert.True(t, b.IsValid())
		assert.Equal(t, 15, b.Root().Val) // 15 is new root value
	})

	t.Run("delete non-existent value — tree unchanged", func(t *testing.T) {
		b := build()
		b.Delete(99)
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7}, vals(b))
	})

	t.Run("delete from empty tree — no panic", func(t *testing.T) {
		b := NewBST()
		assert.NotPanics(t, func() { b.Delete(5) })
		assert.True(t, b.IsEmpty())
	})

	t.Run("delete until empty", func(t *testing.T) {
		b := build()
		for _, v := range []int{4, 2, 6, 1, 3, 5, 7} {
			b.Delete(v)
		}
		assert.True(t, b.IsEmpty())
		assert.Equal(t, 0, b.Size())
	})
}

// ---------------------------------------------------------------------------
// Search
// ---------------------------------------------------------------------------

func Test_BST_Search(t *testing.T) {
	b := NewBSTFromSorted([]int{1, 3, 5, 7, 9})

	t.Run("found", func(t *testing.T) {
		n := b.Search(5)
		assert.NotNil(t, n)
		assert.Equal(t, 5, n.Val)
	})

	t.Run("not found", func(t *testing.T) {
		assert.Nil(t, b.Search(4))
	})

	t.Run("search minimum", func(t *testing.T) {
		assert.Equal(t, 1, nodeVal(b.Search(1)))
	})

	t.Run("search maximum", func(t *testing.T) {
		assert.Equal(t, 9, nodeVal(b.Search(9)))
	})

	t.Run("empty tree", func(t *testing.T) {
		assert.Nil(t, NewBST().Search(1))
	})
}

// ---------------------------------------------------------------------------
// Min / Max
// ---------------------------------------------------------------------------

func Test_BST_MinMax(t *testing.T) {
	t.Run("nil tree", func(t *testing.T) {
		b := NewBST()
		assert.Nil(t, b.Min())
		assert.Nil(t, b.Max())
	})

	t.Run("single node", func(t *testing.T) {
		b := NewBST()
		b.Insert(42)
		assert.Equal(t, 42, nodeVal(b.Min()))
		assert.Equal(t, 42, nodeVal(b.Max()))
	})

	t.Run("balanced tree", func(t *testing.T) {
		b := NewBSTFromSorted([]int{1, 2, 3, 4, 5, 6, 7})
		assert.Equal(t, 1, nodeVal(b.Min()))
		assert.Equal(t, 7, nodeVal(b.Max()))
	})

	t.Run("right-skewed (min is root)", func(t *testing.T) {
		b := NewBST()
		for _, v := range []int{1, 2, 3} {
			b.Insert(v)
		}
		assert.Equal(t, 1, nodeVal(b.Min()))
		assert.Equal(t, 3, nodeVal(b.Max()))
	})
}

// ---------------------------------------------------------------------------
// Floor & Ceiling
// ---------------------------------------------------------------------------

func Test_BST_FloorCeiling(t *testing.T) {
	// BST contains: 1 3 5 7 9
	b := NewBSTFromSorted([]int{1, 3, 5, 7, 9})

	// --- Floor ---
	t.Run("floor exact match", func(t *testing.T) {
		assert.Equal(t, 5, nodeVal(b.Floor(5)))
	})
	t.Run("floor between two nodes", func(t *testing.T) {
		assert.Equal(t, 3, nodeVal(b.Floor(4)))
	})
	t.Run("floor equals minimum", func(t *testing.T) {
		assert.Equal(t, 1, nodeVal(b.Floor(1)))
	})
	t.Run("floor equals maximum", func(t *testing.T) {
		assert.Equal(t, 9, nodeVal(b.Floor(9)))
	})
	t.Run("floor above maximum returns maximum", func(t *testing.T) {
		assert.Equal(t, 9, nodeVal(b.Floor(100)))
	})
	t.Run("floor below minimum returns nil", func(t *testing.T) {
		assert.Nil(t, b.Floor(0))
	})

	// --- Ceiling ---
	t.Run("ceiling exact match", func(t *testing.T) {
		assert.Equal(t, 5, nodeVal(b.Ceiling(5)))
	})
	t.Run("ceiling between two nodes", func(t *testing.T) {
		assert.Equal(t, 5, nodeVal(b.Ceiling(4)))
	})
	t.Run("ceiling equals minimum", func(t *testing.T) {
		assert.Equal(t, 1, nodeVal(b.Ceiling(1)))
	})
	t.Run("ceiling below minimum returns minimum", func(t *testing.T) {
		assert.Equal(t, 1, nodeVal(b.Ceiling(0)))
	})
	t.Run("ceiling above maximum returns nil", func(t *testing.T) {
		assert.Nil(t, b.Ceiling(10))
	})
}

// ---------------------------------------------------------------------------
// Predecessor & Successor
// ---------------------------------------------------------------------------

func Test_BST_PredecessorSuccessor(t *testing.T) {
	// BST contains: 1 3 5 7 9
	b := NewBSTFromSorted([]int{1, 3, 5, 7, 9})

	// --- Predecessor ---
	t.Run("predecessor of interior node", func(t *testing.T) {
		assert.Equal(t, 3, nodeVal(b.Predecessor(5)))
	})
	t.Run("predecessor of maximum", func(t *testing.T) {
		assert.Equal(t, 7, nodeVal(b.Predecessor(9)))
	})
	t.Run("predecessor of minimum returns nil", func(t *testing.T) {
		assert.Nil(t, b.Predecessor(1))
	})
	t.Run("predecessor of value not in tree", func(t *testing.T) {
		// 4 is not in tree; predecessor is the largest value < 4, which is 3.
		assert.Equal(t, 3, nodeVal(b.Predecessor(4)))
	})

	// --- Successor ---
	t.Run("successor of interior node", func(t *testing.T) {
		assert.Equal(t, 7, nodeVal(b.Successor(5)))
	})
	t.Run("successor of minimum", func(t *testing.T) {
		assert.Equal(t, 3, nodeVal(b.Successor(1)))
	})
	t.Run("successor of maximum returns nil", func(t *testing.T) {
		assert.Nil(t, b.Successor(9))
	})
	t.Run("successor of value not in tree", func(t *testing.T) {
		// 4 is not in tree; successor is smallest value > 4, which is 5.
		assert.Equal(t, 5, nodeVal(b.Successor(4)))
	})

	// Predecessor and Successor are symmetric: pred(succ(x)) == x when x is in the tree.
	t.Run("predecessor and successor are inverses", func(t *testing.T) {
		for _, v := range []int{3, 5, 7} {
			succ := b.Successor(v)
			assert.NotNil(t, succ)
			pred := b.Predecessor(succ.Val)
			assert.NotNil(t, pred)
			assert.Equal(t, v, pred.Val)
		}
	})
}

// ---------------------------------------------------------------------------
// Traversals
// ---------------------------------------------------------------------------

func Test_BST_Traversals(t *testing.T) {
	//       4
	//      / \
	//     2   6
	//    / \ / \
	//   1  3 5  7
	b := NewBST()
	for _, v := range []int{4, 2, 6, 1, 3, 5, 7} {
		b.Insert(v)
	}

	t.Run("inorder is sorted", func(t *testing.T) {
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7}, b.Inorder())
	})

	t.Run("preorder is root-left-right", func(t *testing.T) {
		assert.Equal(t, []int{4, 2, 1, 3, 6, 5, 7}, b.Preorder())
	})

	t.Run("postorder is left-right-root", func(t *testing.T) {
		assert.Equal(t, []int{1, 3, 2, 5, 7, 6, 4}, b.Postorder())
	})

	t.Run("level order is BFS grouped by level", func(t *testing.T) {
		assert.Equal(t, [][]int{{4}, {2, 6}, {1, 3, 5, 7}}, b.LevelOrder())
	})

	t.Run("empty tree returns nil or empty", func(t *testing.T) {
		e := NewBST()
		assert.Empty(t, e.Inorder())   // existing impl returns []int{}, not nil
		assert.Nil(t, e.LevelOrder())  // bstLevelOrder explicitly returns nil
	})

	t.Run("single node", func(t *testing.T) {
		s := NewBST()
		s.Insert(99)
		assert.Equal(t, []int{99}, s.Inorder())
		assert.Equal(t, []int{99}, s.Preorder())
		assert.Equal(t, []int{99}, s.Postorder())
		assert.Equal(t, [][]int{{99}}, s.LevelOrder())
	})
}

// ---------------------------------------------------------------------------
// Height & Size
// ---------------------------------------------------------------------------

func Test_BST_HeightSize(t *testing.T) {
	t.Run("empty tree", func(t *testing.T) {
		b := NewBST()
		assert.Equal(t, 0, b.Height())
		assert.Equal(t, 0, b.Size())
	})

	t.Run("single node", func(t *testing.T) {
		b := NewBST()
		b.Insert(1)
		assert.Equal(t, 1, b.Height())
		assert.Equal(t, 1, b.Size())
	})

	t.Run("balanced 7-node tree has height 3", func(t *testing.T) {
		b := NewBSTFromSorted([]int{1, 2, 3, 4, 5, 6, 7})
		assert.Equal(t, 3, b.Height())
		assert.Equal(t, 7, b.Size())
	})

	t.Run("right-skewed n-node tree has height n", func(t *testing.T) {
		b := NewBST()
		for _, v := range []int{1, 2, 3, 4, 5} {
			b.Insert(v)
		}
		assert.Equal(t, 5, b.Height())
		assert.Equal(t, 5, b.Size())
	})

	t.Run("size decreases after delete", func(t *testing.T) {
		b := NewBSTFromSorted([]int{1, 2, 3, 4, 5})
		assert.Equal(t, 5, b.Size())
		b.Delete(3)
		assert.Equal(t, 4, b.Size())
		b.Delete(99) // non-existent
		assert.Equal(t, 4, b.Size())
	})
}

// ---------------------------------------------------------------------------
// NewBSTFromLevelOrder
// ---------------------------------------------------------------------------

func Test_NewBSTFromLevelOrder(t *testing.T) {
	N := math.MinInt // alias for LevelOrderNull

	t.Run("empty input", func(t *testing.T) {
		assert.True(t, NewBSTFromLevelOrder([]int{}).IsEmpty())
	})

	t.Run("null root", func(t *testing.T) {
		assert.True(t, NewBSTFromLevelOrder([]int{N}).IsEmpty())
	})

	t.Run("single node", func(t *testing.T) {
		b := NewBSTFromLevelOrder([]int{5})
		assert.Equal(t, []int{5}, b.Inorder())
	})

	t.Run("complete 3-node tree", func(t *testing.T) {
		// [2, 1, 3] →  2
		//             / \
		//            1   3
		b := NewBSTFromLevelOrder([]int{2, 1, 3})
		assert.True(t, b.IsValid())
		assert.Equal(t, []int{1, 2, 3}, b.Inorder())
		assert.Equal(t, 2, b.Height())
	})

	t.Run("absent left child via LevelOrderNull", func(t *testing.T) {
		// [5, N, 7] →  5
		//               \
		//                7
		b := NewBSTFromLevelOrder([]int{5, N, 7})
		assert.True(t, b.IsValid())
		assert.Equal(t, []int{5, 7}, b.Inorder())
		assert.Nil(t, b.Root().Left)
		assert.Equal(t, 7, b.Root().Right.Val)
	})

	t.Run("leetcode-style 7-node BST", func(t *testing.T) {
		// [4, 2, 6, 1, 3, 5, 7]
		b := NewBSTFromLevelOrder([]int{4, 2, 6, 1, 3, 5, 7})
		assert.True(t, b.IsValid())
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7}, b.Inorder())
		assert.Equal(t, 3, b.Height())
		assert.Equal(t, [][]int{{4}, {2, 6}, {1, 3, 5, 7}}, b.LevelOrder())
	})
}

// ---------------------------------------------------------------------------
// IsValid
// ---------------------------------------------------------------------------

func Test_BST_IsValid(t *testing.T) {
	t.Run("empty tree is valid", func(t *testing.T) {
		assert.True(t, NewBST().IsValid())
	})

	t.Run("tree built by Insert is always valid", func(t *testing.T) {
		b := NewBST()
		for _, v := range []int{5, 3, 7, 1, 4, 6, 8} {
			b.Insert(v)
		}
		assert.True(t, b.IsValid())
	})

	t.Run("manually corrupted tree is invalid", func(t *testing.T) {
		b := NewBSTFromSorted([]int{1, 2, 3})
		// Corrupt: set left child's value greater than root (violates BST).
		b.Root().Left.Val = 99
		assert.False(t, b.IsValid())
	})
}

// ---------------------------------------------------------------------------
// SprintTree / PrintTree
// ---------------------------------------------------------------------------

func Test_BST_SprintTree(t *testing.T) {
	t.Run("empty tree", func(t *testing.T) {
		b := NewBST()
		assert.Equal(t, "<empty>", b.SprintTree())
	})

	t.Run("single node", func(t *testing.T) {
		b := NewBST()
		b.Insert(1)
		assert.Equal(t, "1\n", b.SprintTree())
	})

	t.Run("three-node tree — root with two children", func(t *testing.T) {
		// root=2, left=1, right=3
		//  2
		// / \
		//1   3
		b := NewBSTFromLevelOrder([]int{2, 1, 3})
		out := b.SprintTree()

		// Structure checks: all three values appear, slashes present.
		assert.Contains(t, out, "2")
		assert.Contains(t, out, "1")
		assert.Contains(t, out, "3")
		assert.Contains(t, out, "/")
		assert.Contains(t, out, `\`)

		lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
		assert.Equal(t, 3, len(lines)) // node, slash, node rows
	})

	t.Run("seven-node complete BST", func(t *testing.T) {
		//       4
		//      / \
		//     2   6
		//    / \ / \
		//   1  3 5  7
		b := NewBSTFromLevelOrder([]int{4, 2, 6, 1, 3, 5, 7})
		out := b.SprintTree()

		lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
		assert.Equal(t, 5, len(lines)) // 3 node rows + 2 slash rows

		// Root must be on the first line, centered.
		assert.Contains(t, lines[0], "4")
		// Slash row
		assert.Contains(t, lines[1], "/")
		assert.Contains(t, lines[1], `\`)
		// All leaf values present in the last row.
		for _, v := range []string{"1", "3", "5", "7"} {
			assert.Contains(t, lines[4], v)
		}
	})

	t.Run("right-skewed chain", func(t *testing.T) {
		// 1 → 2 → 3 (right spine)
		b := NewBST()
		for _, v := range []int{1, 2, 3} {
			b.Insert(v)
		}
		out := b.SprintTree()
		// Values appear in descending depth, each shifted right.
		assert.Contains(t, out, "1")
		assert.Contains(t, out, "2")
		assert.Contains(t, out, "3")
		// Backslashes connect the spine.
		assert.Equal(t, 2, strings.Count(out, `\`))
	})

	t.Run("left-skewed chain", func(t *testing.T) {
		// 3 → 2 → 1 (left spine)
		b := NewBST()
		for _, v := range []int{3, 2, 1} {
			b.Insert(v)
		}
		out := b.SprintTree()
		assert.Contains(t, out, "1")
		assert.Contains(t, out, "2")
		assert.Contains(t, out, "3")
		// Forward slashes connect the spine.
		assert.Equal(t, 2, strings.Count(out, "/"))
	})
}
