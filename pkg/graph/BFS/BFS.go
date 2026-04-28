package bfs

import (
	. "gocode/containers"
	. "gocode/types"
	"slices"
)

/**
 * 994. Rotting Oranges
 *
 * 0 representing an empty cell,
 * 1 representing a fresh orange, or
 * 2 representing a rotten orange.
 *
 */
func orangesRotting(grid [][]int) int {
	m, n := len(grid), len(grid[0])
	queue := Queue[[2]int]{} // queue of [row, col]
	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			if grid[r][c] == 2 {
				queue.Enqueue([2]int{r, c})
			}
		}
	}

	directions := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

	minutes := 0
	for !queue.IsEmpty() {
		for i := 0; i < queue.Size(); i++ {
			orange, _ := queue.Peek()
			queue.Dequeue()
			for _, dir := range directions {
				r, c := orange[0]+dir[0], orange[1]+dir[1]
				if r < 0 || r >= m || c < 0 || c >= n {
					continue
				}
				if grid[r][c] != 1 {
					continue
				}
				grid[r][c] = 2
				queue.Enqueue([2]int{r, c})
			}
		}
		if !queue.IsEmpty() {
			minutes++
		}
	}

	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			if grid[r][c] == 1 {
				return -1
			}
		}
	}
	return minutes
}

func orangesRottingSlice(grid [][]int) int {
	m, n := len(grid), len(grid[0])
	queue := [][]int{}
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if grid[i][j] == 2 {
				queue = append(queue, []int{i, j})
			}
			if grid[i][j] == 1 {
			}
		}
	}

	dirs := [][]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
	minute := 0
	for len(queue) > 0 {
		size := len(queue)
		for i := 0; i < size; i++ {
			cell := queue[0]
			queue = queue[1:]

			for _, dir := range dirs {
				r, c := cell[0]+dir[0], cell[1]+dir[1]
				if r >= 0 && r < m && c >= 0 && c < n && grid[r][c] == 1 {
					grid[r][c] = 2
					queue = append(queue, []int{r, c})

				}
			}
		}
		if len(queue) > 0 {
			minute++
		}
	}

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if grid[i][j] == 1 {
				return -1
			}
		}
	}
	return minute
}

func orangesRottingSlice2(grid [][]int) int {
	m, n := len(grid), len(grid[0])
	queue := [][]int{}
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if grid[i][j] == 2 {
				queue = append(queue, []int{i, j})
			}
			if grid[i][j] == 1 {
			}
		}
	}

	dirs := [][]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
	minute := 0
	for len(queue) > 0 {
		for _, cell := range queue {
			queue = queue[1:]

			for _, dir := range dirs {
				r, c := cell[0]+dir[0], cell[1]+dir[1]
				if r >= 0 && r < m && c >= 0 && c < n && grid[r][c] == 1 {
					grid[r][c] = 2
					queue = append(queue, []int{r, c})
				}
			}
		}
		if len(queue) > 0 {
			minute++
		}
	}

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if grid[i][j] == 1 {
				return -1
			}
		}
	}
	return minute
}

/**
 * 199. Binary Tree Right Side View
 */
func rightSideView(root *TreeNode) []int {
	rightSide := make([]int, 0)
	if root == nil {
		return rightSide
	}

	queue := []*TreeNode{root}
	for len(queue) > 0 {
		levelSize := len(queue)

		for i := 0; i < levelSize; i++ {
			node := queue[0]
			queue = queue[1:]

			// last node in level is visible from right side
			if i == levelSize-1 {
				rightSide = append(rightSide, node.Val)
			}

			if node.Left != nil {
				queue = append(queue, node.Left)
			}
			if node.Right != nil {
				queue = append(queue, node.Right)
			}
		}
	}

	return rightSide
}

/**
 * 102. Binary Tree Level Order Traversal
 */
func levelOrder(root *TreeNode) [][]int {
	if root == nil {
		return [][]int{}
	}

	var levels [][]int
	queue := []*TreeNode{root}
	for len(queue) > 0 {
		levelSize := len(queue)
		var level []int
		for i := 0; i < levelSize; i++ {
			node := queue[0]
			queue = queue[1:]

			level = append(level, node.Val)
			if node.Left != nil {
				queue = append(queue, node.Left)
			}
			if node.Right != nil {
				queue = append(queue, node.Right)
			}
		}
		levels = append(levels, level)
	}
	return levels
}

/**
 * 103. Binary Tree Zigzag Level Order Traversal
 */
func zigzagLevelOrder(root *TreeNode) [][]int {
	if root == nil {
		return [][]int{}
	}

	queue := []*TreeNode{root}
	levels := [][]int{}
	for len(queue) > 0 {
		size := len(queue)
		level := []int{}
		for i := 0; i < size; i++ {
			node := queue[0]
			queue = queue[1:]
			level = append(level, node.Val)
			if node.Left != nil {
				queue = append(queue, node.Left)
			}
			if node.Right != nil {
				queue = append(queue, node.Right)
			}
		}

		if len(levels)%2 == 1 {
			slices.Reverse(level)
		}
		levels = append(levels, level)
	}
	return levels
}

/**
 * 863. All Nodes Distance K in Binary Tree
 *
 * parentMap
 */
func distanceK(root *TreeNode, target *TreeNode, k int) []int {
	nodeToParent := make(map[*TreeNode]*TreeNode)

	var dfs func(node *TreeNode)
	dfs = func(node *TreeNode) {
		if node == nil {
			return
		}
		if node.Left != nil {
			nodeToParent[node.Left] = node
			dfs(node.Left)
		}
		if node.Right != nil {
			nodeToParent[node.Right] = node
			dfs(node.Right)
		}
	}
	dfs(root)

	queue := []*TreeNode{target}
	visited := map[*TreeNode]bool{target: true}
	stepsFromTarget := 0

	for len(queue) > 0 {
		if stepsFromTarget == k {
			break
		}
		size := len(queue)
		for i := 0; i < size; i++ {
			curNode := queue[0]
			queue = queue[1:]

			if curNode.Left != nil && !visited[curNode.Left] {
				visited[curNode.Left] = true
				queue = append(queue, curNode.Left)
			}
			if curNode.Right != nil && !visited[curNode.Right] {
				visited[curNode.Right] = true
				queue = append(queue, curNode.Right)
			}
			if parentNode, found := nodeToParent[curNode]; found && !visited[parentNode] {
				visited[parentNode] = true
				queue = append(queue, parentNode)
			}
		}
		stepsFromTarget++
	}

	result := make([]int, 0, len(queue))
	for _, node := range queue {
		result = append(result, node.Val)
	}
	return result
}

func distanceKClaude(root *TreeNode, target *TreeNode, k int) []int {
	// pass 1: build parent pointers
	parent := map[*TreeNode]*TreeNode{}
	var dfs func(*TreeNode, *TreeNode)
	dfs = func(node, p *TreeNode) {
		if node == nil {
			return
		}
		parent[node] = p
		dfs(node.Left, node)
		dfs(node.Right, node)
	}
	dfs(root, nil)

	// pass 2: BFS from target treating tree as undirected graph
	queue := []*TreeNode{target}
	visited := map[*TreeNode]bool{target: true}

	for dist := 0; len(queue) > 0; dist++ {
		if dist == k {
			res := []int{}
			for _, node := range queue {
				res = append(res, node.Val)
			}
			return res
		}

		next := []*TreeNode{}
		for _, node := range queue {
			for _, neighbor := range []*TreeNode{node.Left, node.Right, parent[node]} {
				if neighbor != nil && !visited[neighbor] {
					visited[neighbor] = true
					next = append(next, neighbor)
				}
			}
		}
		queue = next
	}
	return []int{}
}
