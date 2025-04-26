package graph

import (
	. "gocode/types"
)

/**
 * 543. Diameter of Binary Tree
 */
func diameterOfBinaryTree(root *TreeNode) int {
    max := 0
    var dfs func(node *TreeNode) int
    dfs = func(node *TreeNode) int {
        if node == nil {
            return 0
        }
        left, right := dfs(node.Left), dfs(node.Right)
        max = Max(max, left+right)
        return 1 + Max(left, right)
    }
    dfs(root)
    return max
}

func Max(a, b int) int {
    if a > b {
        return a
    }
    return b
}
