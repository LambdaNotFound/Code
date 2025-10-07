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
func cloneGraph(node *Node) *Node {
    if node == nil {
        return nil
    }

    queue := []*Node{node}
    mapToCloned := make(map[*Node]*Node)
    mapToCloned[node] = &Node{Val: node.Val} // tracks visited nodes

    for len(queue) > 0 {
        currNode := queue[0]
        queue = queue[1:]

        for _, neighbor := range currNode.Neighbors {
            if _, exist := mapToCloned[neighbor]; !exist {
                mapToCloned[neighbor] = &Node{Val: neighbor.Val}
                queue = append(queue, neighbor) // only enqueue un-visited
            }

            currNodeCopy := mapToCloned[currNode]
            neighborCopy := mapToCloned[neighbor]
            currNodeCopy.Neighbors = append(currNodeCopy.Neighbors, neighborCopy)
        }
    }

    return mapToCloned[node]
}

func cloneGraph_DFS(node *Node, visited map[*Node]*Node) *Node {
    // If the node is already cloned, return it
    if clone, found := visited[node]; found {
        return clone
    }

    // Create a new clone
    clone := &Node{Val: node.Val}
    visited[node] = clone // mark as cloned before visiting neighbors

    // Recurse on neighbors
    for _, neighbor := range node.Neighbors {
        neighborCopy := cloneGraph_DFS(neighbor, visited)
        clone.Neighbors = append(clone.Neighbors, neighborCopy)
    }

    return clone
}