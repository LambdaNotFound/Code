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
