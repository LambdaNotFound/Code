package divide_and_conquer

import (
	"testing"

	. "gocode/types"
	"gocode/utils"

	"github.com/stretchr/testify/assert"
)

func Test_sortList(t *testing.T) {
    testCases := []struct {
        name     string
        list     []int
        expected []int
    }{
        {
            "case 1",
            []int{4, 2, 1, 3},
            []int{1, 2, 3, 4},
        },
        {
            "case 2",
            []int{-1, 5, 3, 4, 0},
            []int{-1, 0, 3, 4, 5},
        },
        {
            "case 3",
            []int{},
            []int{},
        },
        {
            name:     "single_node",
            list:     []int{7},
            expected: []int{7},
        },
        {
            name:     "already_sorted",
            list:     []int{1, 2, 3, 4, 5},
            expected: []int{1, 2, 3, 4, 5},
        },
        {
            name:     "reverse_sorted",
            list:     []int{5, 4, 3, 2, 1},
            expected: []int{1, 2, 3, 4, 5},
        },
        {
            name:     "with_duplicates",
            list:     []int{3, 1, 2, 1, 3},
            expected: []int{1, 1, 2, 3, 3},
        },
        {
            name:     "all_same",
            list:     []int{4, 4, 4, 4},
            expected: []int{4, 4, 4, 4},
        },
        {
            name:     "mixed_negatives",
            list:     []int{10, -3, 7, -1, 2},
            expected: []int{-3, -1, 2, 7, 10},
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            list := utils.CreateLinkedList(tc.list)
            result := sortListMergeSort(list)
            expected := utils.CreateLinkedList(tc.expected)
            isEqual := utils.VerifyLinkedLists(result, expected)
            assert.Equal(t, true, isEqual)
        })
    }
}

func Test_MergeSort(t *testing.T) {
    tests := []struct {
        name     string
        input    []int
        expected []int
    }{
        {
            name:     "Empty slice",
            input:    []int{},
            expected: []int{},
        },
        {
            name:     "Single element",
            input:    []int{42},
            expected: []int{42},
        },
        {
            name:     "Already sorted",
            input:    []int{1, 2, 3, 4, 5},
            expected: []int{1, 2, 3, 4, 5},
        },
        {
            name:     "Reverse order",
            input:    []int{5, 4, 3, 2, 1},
            expected: []int{1, 2, 3, 4, 5},
        },
        {
            name:     "With duplicates",
            input:    []int{4, 2, 4, 1, 3, 2},
            expected: []int{1, 2, 2, 3, 4, 4},
        },
        {
            name:     "All same elements",
            input:    []int{7, 7, 7, 7},
            expected: []int{7, 7, 7, 7},
        },
        {
            name:     "Contains negative numbers",
            input:    []int{-3, -1, -7, 2, 0},
            expected: []int{-7, -3, -1, 0, 2},
        },
        {
            name:     "Random order",
            input:    []int{10, 1, 5, 8, 3, 2, 9, 4, 7, 6},
            expected: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            got := MergeSort(tc.input)
            assert.Equal(t, tc.expected, got)
        })
    }
}

func Test_sortedListToBST(t *testing.T) {
	t.Run("nil_returns_nil", func(t *testing.T) {
		assert.Nil(t, sortedListToBST(nil))
	})

	t.Run("single_node", func(t *testing.T) {
		head := &ListNode{Val: 5}
		root := sortedListToBST(head)
		assert.NotNil(t, root)
		assert.Equal(t, 5, root.Val)
		assert.Nil(t, root.Left)
		assert.Nil(t, root.Right)
	})

	t.Run("three_nodes_balanced", func(t *testing.T) {
		// 1 -> 2 -> 3  =>  root=2, left=1, right=3
		head := utils.CreateLinkedList([]int{1, 2, 3})
		root := sortedListToBST(head)
		assert.Equal(t, 2, root.Val)
		assert.Equal(t, 1, root.Left.Val)
		assert.Equal(t, 3, root.Right.Val)
	})

	t.Run("two_nodes", func(t *testing.T) {
		// upper-midpoint: [1,2] → root=2, left=1, right=nil
		head := utils.CreateLinkedList([]int{1, 2})
		root := sortedListToBST(head)
		assert.NotNil(t, root)
		assert.Equal(t, 2, root.Val)
		assert.Equal(t, 1, root.Left.Val)
		assert.Nil(t, root.Right)
	})

	t.Run("five_nodes_in_order_bst", func(t *testing.T) {
		head := utils.CreateLinkedList([]int{1, 2, 3, 4, 5})
		root := sortedListToBST(head)
		// verify BST property: in-order traversal must match [1,2,3,4,5]
		var inorder func(*TreeNode) []int
		inorder = func(n *TreeNode) []int {
			if n == nil {
				return nil
			}
			return append(append(inorder(n.Left), n.Val), inorder(n.Right)...)
		}
		assert.Equal(t, []int{1, 2, 3, 4, 5}, inorder(root))
	})
}