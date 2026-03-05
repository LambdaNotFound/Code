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
 * BFS: Time: O(m x n), Space: O(m x n)
 */
func numIslands(grid [][]byte) int {
	m, n := len(grid), len(grid[0])
	directions := [][]int{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}

	var dfs func(int, int)
	dfs = func(i, j int) {
		if grid[i][j] != '1' {
			return
		}

		if grid[i][j] == '1' {
			grid[i][j] = 'X'
		}
		for _, dir := range directions {
			row, col := i+dir[0], j+dir[1]
			if row >= 0 && row < m && col >= 0 && col < n {
				dfs(row, col)
			}
		}
	}

	res := 0
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if grid[i][j] == '1' {
				res++
				dfs(i, j)
			}
		}
	}
	return res
}

func numIslandsBFS(grid [][]byte) int {
	m, n := len(grid), len(grid[0])
	directions := [][]int{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}

	var bfs func(int, int)
	bfs = func(i, j int) {
		queue := [][]int{[]int{i, j}}
		grid[i][j] = 'X'
		for len(queue) > 0 {
			cur := queue[0]
			queue = queue[1:]

			for _, dir := range directions {
				row, col := cur[0]+dir[0], cur[1]+dir[1]
				if (row >= 0 && row < m && col >= 0 && col < n) && grid[row][col] == '1' {
					grid[row][col] = 'X'
					queue = append(queue, []int{row, col})
				}
			}
		}
	}

	res := 0
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if grid[i][j] == '1' {
				res++
				bfs(i, j)
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
