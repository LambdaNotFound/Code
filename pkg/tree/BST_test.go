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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, q := nodeMap[tt.p], nodeMap[tt.q]
			result := lowestCommonAncestorBST(root, p, q)
			assert.NotNil(t, result)
			assert.Equal(t, tt.expected, result.Val)
		})
	}

}
