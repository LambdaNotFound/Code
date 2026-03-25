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
	visited := make(map[*Node]*Node) // map[ptr]ptr

	var dfs func(*Node) *Node
	dfs = func(node *Node) *Node {
		if node == nil {
			return nil
		}

		// If the node is already cloned, return it
		if copy, exist := visited[node]; exist {
			return copy
		}

		// Create a new clone
		copy := &Node{Val: node.Val}
		visited[node] = copy // mark as cloned before visiting neighbors

		// Recurse on neighbors
		for _, neighbor := range node.Neighbors {
			neighborClone := dfs(neighbor)
			copy.Neighbors = append(copy.Neighbors, neighborClone)
			// clone.Neighbors = append(clone.Neighbors, dfs(neighbor))
		}
		return copy
	}

	return dfs(node)
}

func cloneGraphBFS(node *Node) *Node {
	if node == nil {
		return nil
	}

	queue := []*Node{node}
	mapToCopy := make(map[*Node]*Node)
	mapToCopy[node] = &Node{Val: node.Val} // tracks visited nodes

	for len(queue) > 0 {
		currNode := queue[0]
		queue = queue[1:]
		currNodeCopy := mapToCopy[currNode]

		for _, neighbor := range currNode.Neighbors {
			if _, exist := mapToCopy[neighbor]; !exist {
				mapToCopy[neighbor] = &Node{Val: neighbor.Val}
				queue = append(queue, neighbor) // only enqueue un-visited
			}

			neighborCopy := mapToCopy[neighbor]
			currNodeCopy.Neighbors = append(currNodeCopy.Neighbors, neighborCopy)
		}
	}

	return mapToCopy[node]
}
