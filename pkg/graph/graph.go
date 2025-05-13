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

        for _, neighbor := range curr.Neighbors {
            if _, exist := mapToCloned[neighbor]; !exist {
                mapToCloned[neighbor] = &Node{Val: neighbor.Val}
                queue = append(queue, neighbor)
            }

            mapToCloned[curr].Neighbors = append(mapToCloned[curr].Neighbors, mapToCloned[neighbor])
        }
    }

    return mapToCloned[node]
}

/**
 * 310. Minimum Height Trees
 *
 * The height of a rooted tree is the number of edges on the longest downward path between the root and a leaf,
 *
 * Adjacency List
 */
func findMinHeightTrees(n int, edges [][]int) []int {
    if n == 1 {
        return []int{0}
    }

    graph := map[int][]int{}

    for _, edge := range edges {
        graph[edge[0]] = append(graph[edge[0]], edge[1])
        graph[edge[1]] = append(graph[edge[1]], edge[0])
    }

    leaves := []int{}
    for k, v := range graph {
        if len(v) == 1 {
            leaves = append(leaves, k)
        }
    }

    for len(leaves) < n {
        n -= len(leaves)

        new_leaves := []int{}

        for _, leaf := range leaves {
            node := graph[leaf][0] // remove leaf from node
            for i := 0; i < len(graph[node]); i++ {
                if graph[node][i] == leaf {
                    graph[node][i] = graph[node][len(graph[node])-1] // copy/move to truncate
                    break
                }
            }
            graph[node] = graph[node][:len(graph[node])-1]

            if len(graph[node]) == 1 {
                new_leaves = append(new_leaves, node)
            }
        }

        leaves = new_leaves
    }

    return leaves
}
