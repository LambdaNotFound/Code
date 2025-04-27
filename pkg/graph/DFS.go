package graph

import (
	. "gocode/types"
)

/**
 * 543. Diameter of Binary Tree
 */
func diameterOfBinaryTree(root *TreeNode) int {
    res := 0
    max := func(a, b int) int {
        if a > b {
            return a
        }
        return b
    }

    var dfs func(node *TreeNode) int
    dfs = func(node *TreeNode) int {
        if node == nil {
            return 0
        }
        left, right := dfs(node.Left), dfs(node.Right)
        res = max(res, left+right)
        return 1 + max(left, right)
    }
    dfs(root)
    return res
}
