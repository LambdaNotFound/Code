package tree

import (
	. "gocode/types"
	"math"
)

/**
 * 98. Validate Binary Search Tree
 */
func isValidBST(root *TreeNode) bool {
    var validate func(*TreeNode, int, int) bool
    validate = func(node *TreeNode, lower int, upper int) bool {
        if node == nil {
            return true
        }

        if (lower < node.Val) && (node.Val < upper) {
            return validate(node.Left, lower, node.Val) && validate(node.Right, node.Val, upper)
        } else {
            return false
        }
    }

    return validate(root, math.MinInt, math.MaxInt)
}

/**
 * 230. Kth Smallest Element in a BST
 */
func kthSmallest(root *TreeNode, k int) int {
    stack := []*TreeNode{}
    curr := root
    for root != nil {
        for curr != nil {
            stack = append(stack, curr)
            curr = curr.Left
        }
        curr = stack[len(stack)-1]
        stack = stack[:len(stack)-1]
        k--
        if k == 0 {
            return curr.Val
        }

        curr = curr.Right
    }
    return -1
}
