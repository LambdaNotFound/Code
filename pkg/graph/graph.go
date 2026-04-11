package graph

import (
	. "gocode/types"
)

/**
 * BFS
 */
func BFS(start *Node) {
	if start == nil {
		return
	}

	visited := make(map[*Node]bool)
	queue := []*Node{start}
	visited[start] = true

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		for _, neighbor := range curr.Neighbors {
			if !visited[neighbor] {
				visited[neighbor] = true
				queue = append(queue, neighbor)
			}
		}
	}
}

/**
 * DFS
 */
func DFS(start *Node, visited map[*Node]bool) {
	if start == nil {
		return
	}
	if visited[start] {
		return
	}

	visited[start] = true

	for _, neighbor := range start.Neighbors {
		DFS(neighbor, visited)
	}
}

/**
 * 133. Clone Graph - a connected undirected graph
 *
 * BFS + HashMap: queue, map[*Node]*Node
 */
func cloneGraphDFS(node *Node) *Node {
	cloned := make(map[*Node]*Node)

	var dfs func(*Node) *Node
	dfs = func(node *Node) *Node {
		if node == nil {
			return nil
		}

		if clone, exists := cloned[node]; exists {
			return clone
		}

		clone := &Node{Val: node.Val}
		cloned[node] = clone

		for _, neighbor := range node.Neighbors {
			clone.Neighbors = append(clone.Neighbors, dfs(neighbor))
		}
		return clone
	}

	return dfs(node)
}

func cloneGraphBFS(node *Node) *Node {
	if node == nil {
		return nil
	}

	queue := []*Node{node}
	cloned := make(map[*Node]*Node)
	cloned[node] = &Node{Val: node.Val}

	for len(queue) > 0 {
		original := queue[0]
		queue = queue[1:]
		clone := cloned[original]

		for _, neighbor := range original.Neighbors {
			if _, exists := cloned[neighbor]; !exists {
				cloned[neighbor] = &Node{Val: neighbor.Val}
				queue = append(queue, neighbor)
			}
			clone.Neighbors = append(clone.Neighbors, cloned[neighbor])
		}
	}

	return cloned[node]
}
