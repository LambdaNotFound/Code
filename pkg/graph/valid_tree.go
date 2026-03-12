package graph

/**
 * Graph Valid Tree
 *
 *
 *	// build undirected adjacency list — each edge goes both ways
 *	adjList := make(map[int][]int)
 *	for _, edge := range edges {
 *		src, dst := edge[0], edge[1]
 *		adjList[src] = append(adjList[src], dst)
 *		adjList[dst] = append(adjList[dst], src) // undirected
 *	}
 *
 */
func validTreeBFS(n int, edges [][]int) bool {
	// a valid tree must have exactly n-1 edges
	// more edges = cycle, fewer edges = disconnected
	if len(edges) != n-1 {
		return false
	}

	// build undirected adjacency list — each edge goes both ways
	adjList := make(map[int][]int)
	for _, edge := range edges {
		src, dst := edge[0], edge[1]
		adjList[src] = append(adjList[src], dst)
		adjList[dst] = append(adjList[dst], src) // undirected
	}

	// BFS from node 0 — track visited and parent to avoid
	// misidentifying parent edge as a cycle
	visited := make([]bool, n)
	parent := make([]int, n)
	for i := range parent {
		parent[i] = -1 // no parent initially
	}

	queue := []int{0} // always safe to start from 0 since edges = n-1 >= 0
	visited[0] = true

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		for _, neighbor := range adjList[cur] {
			if neighbor == parent[cur] {
				// skip parent edge — not a cycle in undirected graph
				continue
			}
			if visited[neighbor] {
				// visited non-parent node = cycle detected
				return false
			}
			visited[neighbor] = true
			parent[neighbor] = cur
			queue = append(queue, neighbor)
		}
	}

	// check all nodes were visited — ensures full connectivity
	for _, wasVisited := range visited {
		if !wasVisited {
			return false
		}
	}

	return true
}
