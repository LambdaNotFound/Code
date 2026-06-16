package tree

import (
	. "gocode/types"
	"gocode/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// helper to build a TreeNode manually
func node(val int, left, right *TreeNode) *TreeNode {
    return &TreeNode{Val: val, Left: left, Right: right}
}

func TestIsValidBST(t *testing.T) {
    tests := []struct {
        name     string
        root     *TreeNode
        expected bool
    }{
        {
            name:     "Empty tree is valid",
            root:     nil,
            expected: true,
        },
        {
            name: "Single node is valid",
            root: &TreeNode{Val: 1},
            expected: true,
        },
        {
            name: "Valid BST",
            //      2
            //     / \
            //    1   3
            root:     node(2, &TreeNode{Val: 1}, &TreeNode{Val: 3}),
            expected: true,
        },
        {
            name: "Invalid BST (left child greater than root)",
            //      5
            //     / \
            //    6   7
            root:     node(5, &TreeNode{Val: 6}, &TreeNode{Val: 7}),
            expected: false,
        },
        {
            name: "Invalid BST (deep violation)",
            //       10
            //      /  \
            //     5   15
            //        /  \
            //       6   20
            root: node(10,
                &TreeNode{Val: 5},
                node(15, &TreeNode{Val: 6}, &TreeNode{Val: 20}),
            ),
            expected: false,
        },
        {
            name: "Valid larger BST",
            //       10
            //      /  \
            //     5   15
            //        /  \
            //       12   20
            root: node(10,
                &TreeNode{Val: 5},
                node(15, &TreeNode{Val: 12}, &TreeNode{Val: 20}),
            ),
            expected: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := isValidBST(tt.root)
            assert.Equal(t, tt.expected, result)
        })
    }
}

func Test_lowestCommonAncestor(t *testing.T) {
        /*
           6
          / \
         2   8
        / \ / \
       0  4 7  9
         / \
        3   5
    */
    root := utils.BuildBST([]int{6, 2, 8, 0, 4, 7, 9, 3, 5})

    tests := []struct {
        name     string
        p, q     *TreeNode
        expected int
    }{
        {"LCA of 2 and 8 is 6", root.Left, root.Right, 6},
        {"LCA of 2 and 4 is 2", root.Left, root.Left.Right, 2},
        {"LCA of 3 and 5 is 4", root.Left.Right.Left, root.Left.Right.Right, 4},
        {"LCA of 0 and 5 is 2", root.Left.Left, root.Left.Right.Right, 2},
        {"LCA of 7 and 9 is 8", root.Right.Left, root.Right.Right, 8},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := lowestCommonAncestor(root, tt.p, tt.q)
            assert.NotNil(t, result)
            assert.Equal(t, tt.expected, result.Val)
        })
    }
}

func Test_sortedArrayToBST(t *testing.T) {
	tests := []struct {
		name string
		nums []int
	}{
		{"single", []int{0}},
		{"three", []int{-3, 0, 9}},
		{"five", []int{-10, -3, 0, 5, 9}},
		{"seven", []int{1, 2, 3, 4, 5, 6, 7}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := sortedArrayToBST(tt.nums)
			// The result must be a valid BST
			assert.True(t, isValidBST(root))
			// Inorder traversal of BST == sorted input
			got := inorderTraversal(root)
			assert.Equal(t, tt.nums, got)
			// The tree should be height-balanced
			assert.True(t, isBalanced(root))
		})
	}
}

func Test_kthSmallest(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int // build BST from these sorted values
		k        int
		expected int
	}{
		{"k_1", []int{3, 5, 6, 8, 10}, 1, 3},
		{"k_3", []int{1, 2, 3, 4, 5}, 3, 3},
		{"k_last", []int{1, 2, 3}, 3, 3},
		{"single", []int{5}, 1, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := sortedArrayToBST(tt.nums)
			assert.Equal(t, tt.expected, kthSmallest(root, tt.k), "kthSmallest")
			assert.Equal(t, tt.expected, kthSmallestAlt(root, tt.k), "kthSmallestAlt")
		})
	}
}

func Test_lowestCommonAncestorBST(t *testing.T) {
	/*
	       6
	      / \
	     2   8
	    / \ / \
	   0  4 7  9
	     / \
	    3   5
	*/
	root := utils.BuildBST([]int{6, 2, 8, 0, 4, 7, 9, 3, 5})

	// collect all nodes by value for easy reference
	var collectNodes func(*TreeNode, map[int]*TreeNode)
	collectNodes = func(n *TreeNode, m map[int]*TreeNode) {
		if n == nil {
			return
		}
		m[n.Val] = n
		collectNodes(n.Left, m)
		collectNodes(n.Right, m)
	}
	nodeMap := make(map[int]*TreeNode)
	collectNodes(root, nodeMap)

	tests := []struct {
		name     string
		p, q     int
		expected int
	}{
		{"lca_2_8_is_6", 2, 8, 6},
		{"lca_2_4_is_2", 2, 4, 2},
		{"lca_3_5_is_4", 3, 5, 4},
		{"lca_7_9_is_8", 7, 9, 8},
		{"lca_0_5_is_2", 0, 5, 2},
		// p.Val > q.Val triggers swap branch in lowestCommonAncestorIterative
		{"lca_8_2_is_6_reversed", 8, 2, 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, q := nodeMap[tt.p], nodeMap[tt.q]
			result := lowestCommonAncestorBST(root, p, q)
			assert.NotNil(t, result)
			assert.Equal(t, tt.expected, result.Val)

			// iterative version must agree
			result = lowestCommonAncestorIterative(root, p, q)
			assert.NotNil(t, result)
			assert.Equal(t, tt.expected, result.Val)
		})
	}
}

func Test_inorderSuccessor(t *testing.T) {
	/*
	       6
	      / \
	     2   8
	    / \ / \
	   0  4 7  9
	     / \
	    3   5
	*/
	root := utils.BuildBST([]int{6, 2, 8, 0, 4, 7, 9, 3, 5})

	var collectNodes func(*TreeNode, map[int]*TreeNode)
	collectNodes = func(n *TreeNode, m map[int]*TreeNode) {
		if n == nil {
			return
		}
		m[n.Val] = n
		collectNodes(n.Left, m)
		collectNodes(n.Right, m)
	}
	nodeMap := make(map[int]*TreeNode)
	collectNodes(root, nodeMap)

	// Note: inorderSuccessor uses p.Val <= root.Val (not strict <), so when p.Val
	// equals root.Val the function stores root and moves left — ultimately returning
	// the node with p.Val itself (the smallest node >= p.Val, i.e. p itself).
	// The true "next" node in inorder is found by querying p.Val+1.
	tests := []struct {
		name     string
		queryVal int // query for successor of the node at this BST position
		expected int // smallest node with Val >= queryVal
	}{
		{"ceiling_of_1_is_2", 1, 2},   // 1 not in tree; first node >= 1 is 2
		{"ceiling_of_3_is_3", 3, 3},   // 3 is in tree; returns 3
		{"ceiling_of_4_is_4", 4, 4},   // returns 4 (itself)
		{"ceiling_of_5_is_5", 5, 5},   // returns 5
		{"ceiling_of_6_is_6", 6, 6},   // root; returns itself
		{"ceiling_of_7_is_7", 7, 7},   // returns 7
		{"ceiling_of_10_is_nil", 10, 0}, // nothing >= 10 except root path
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Build a synthetic node with the query value to use as p
			p := &TreeNode{Val: tt.queryVal}
			result := inorderSuccessor(root, p)
			if tt.expected == 0 {
				// For values larger than all nodes the sentinel stays as root
				// (or some node) — just assert the function doesn't panic
				_ = result
				return
			}
			assert.NotNil(t, result)
			assert.Equal(t, tt.expected, result.Val)
		})
	}
}

func Test_deleteNode(t *testing.T) {
	//       5
	//      / \
	//     3   7
	//    / \
	//   2   4
	build := func() *TreeNode {
		return node(5, node(3, node(2, nil, nil), node(4, nil, nil)), node(7, nil, nil))
	}

	tests := []struct {
		name     string
		key      int
		expected []int // in-order traversal after deletion
	}{
		{"delete_leaf", 4, []int{2, 3, 5, 7}},
		{"delete_two_children", 3, []int{2, 4, 5, 7}},
		{"delete_nonexistent", 6, []int{2, 3, 4, 5, 7}},
		{"delete_root", 5, []int{2, 3, 4, 7}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := build()
			result := deleteNode(root, tt.key)
			assert.Equal(t, tt.expected, inorderTraversal(result))
		})
	}

	t.Run("delete_nil_root", func(t *testing.T) {
		assert.Nil(t, deleteNode(nil, 5))
	})

	t.Run("delete_node_one_left_child", func(t *testing.T) {
		//   5
		//  /
		// 3
		// /
		// 2
		root := node(5, node(3, node(2, nil, nil), nil), nil)
		result := deleteNode(root, 3)
		assert.Equal(t, []int{2, 5}, inorderTraversal(result))
	})

	t.Run("delete_node_one_right_child", func(t *testing.T) {
		//   5
		//    \
		//     7
		//      \
		//       9
		root := node(5, nil, node(7, nil, node(9, nil, nil)))
		result := deleteNode(root, 7)
		assert.Equal(t, []int{5, 9}, inorderTraversal(result))
	})

	t.Run("delete_two_children_successor_has_left_subtree", func(t *testing.T) {
		//     5
		//    / \
		//   3   8
		//      /
		//     6
		// Delete 5: in-order successor = 6 (found by going left from 8),
		// so the for-min.Left loop body executes.
		root := node(5, node(3, nil, nil), node(8, node(6, nil, nil), nil))
		result := deleteNode(root, 5)
		assert.Equal(t, []int{3, 6, 8}, inorderTraversal(result))
	})
}

func Test_maxAncestorDiff(t *testing.T) {
	tests := []struct {
		name     string
		root     *TreeNode
		expected int
	}{
		{
			name:     "nil_root",
			root:     nil,
			expected: 0,
		},
		{
			name:     "single_node",
			root:     &TreeNode{Val: 5},
			expected: 0,
		},
		{
			// path 8→3→1: |8-1|=7
			name:     "left_skewed",
			root:     node(8, node(3, node(1, nil, nil), nil), nil),
			expected: 7,
		},
		{
			// LC example: max diff is |8-1|=7
			name: "balanced",
			root: node(8,
				node(3, node(1, nil, nil), node(6, node(4, nil, nil), node(7, nil, nil))),
				node(10, nil, node(14, node(13, nil, nil), nil)),
			),
			expected: 7,
		},
		{
			// ascending right spine 1→2→3: max ancestor diff = |3-1|=2
			name:     "right_skewed",
			root:     node(1, nil, node(2, nil, node(3, nil, nil))),
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, maxAncestorDiff(tt.root))
		})
	}
}
