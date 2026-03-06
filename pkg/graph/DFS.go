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

	var dfs func(*TreeNode) int
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

/**
 * 104. Maximum Depth of Binary Tree
 */
func maxDepth(root *TreeNode) int {
	if root == nil {
		return 0
	}
	return 1 + max(maxDepth(root.Left), maxDepth(root.Right))
}

/**
 * 733. Flood Fill
 */
func floodFill(image [][]int, sr int, sc int, newColor int) [][]int {
	var dfs func(int, int, int, int)
	dfs = func(x int, y int, sourceColor int, newColor int) {
		if x < 0 || x >= len(image) || y < 0 || y >= len(image[0]) ||
			image[x][y] != sourceColor {
			return
		}
		image[x][y] = newColor
		dfs(x-1, y, sourceColor, newColor)
		dfs(x+1, y, sourceColor, newColor)
		dfs(x, y-1, sourceColor, newColor)
		dfs(x, y+1, sourceColor, newColor)
	}
	sourceColor := image[sr][sc]
	if sourceColor != newColor {
		dfs(sr, sc, sourceColor, newColor)
	}

	return image
}
