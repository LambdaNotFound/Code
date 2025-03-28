package graph

import . "gocode/types"

/**
 * 133. Clone Graph
 *
 * BFS + HashMap
 */
func cloneGraph(node *Node) *Node {
    if node == nil {
        return nil
    }

    queue := []*Node{node}
    mapToCloned := make(map[*Node]*Node)
    mapToCloned[node] = &Node{Val: node.Val}

    for len(queue) > 0 {
        curr := queue[0]
        queue = queue[1:]

        currCloned := mapToCloned[curr]
        for _, neighbor := range curr.Neighbors {
            if _, exist := mapToCloned[neighbor]; !exist {
                mapToCloned[neighbor] = &Node{Val: neighbor.Val}
                queue = append(queue, neighbor)
            }

            currCloned.Neighbors = append(currCloned.Neighbors, mapToCloned[neighbor])
        }
    }

    return mapToCloned[node]
}
