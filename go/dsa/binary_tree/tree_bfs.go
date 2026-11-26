package binarytree

import (
	. "gocode/dsa/types"
	"slices"
)

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

/**
 * 958. Check Completeness of a Binary Tree
 *
 * Given the root of a binary tree, determine if it is a complete binary tree.
 */
func isCompleteTree(root *TreeNode) bool {
	queue := []*TreeNode{root}
	foundNil := false

	for len(queue) > 0 {
		curNode := queue[0]
		queue = queue[1:]

		if curNode == nil {
			foundNil = true
		} else {
			if foundNil {
				return false
			}
			queue = append(queue, curNode.Left)
			queue = append(queue, curNode.Right)
		}
	}

	return true
}
