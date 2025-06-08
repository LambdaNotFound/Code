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
 * 200. Number of Islands
 *
 * Given an m x n 2D binary grid grid which represents a map of '1's (land) and '0's (water),
 * return the number of islands.
 *
 * DFS: Time: O(m x n), Space: O(m x n)
 */
func numIslands(grid [][]byte) int {
    var dfsHelper func(int, int)
    dfsHelper = func(i, j int) {
        if grid[i][j] != '1' {
            return
        }
        grid[i][j] = '2'

        if i != 0 {
            dfsHelper(i-1, j)
        }
        if i != len(grid)-1 {
            dfsHelper(i+1, j)
        }
        if j != 0 {
            dfsHelper(i, j-1)
        }
        if j != len(grid[0])-1 {
            dfsHelper(i, j+1)
        }
    }

    res := 0
    for i := 0; i < len(grid); i++ {
        for j := 0; j < len(grid[0]); j++ {
            if grid[i][j] == '1' {
                res++
                dfsHelper(i, j)
            }
        }
    }
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
        if x < 0 || x >= len(image) || y < 0 || y >= len(image[0]) || image[x][y] != sourceColor {
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
