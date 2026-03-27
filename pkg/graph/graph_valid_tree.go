package graph

/**
 * Graph Valid Tree
 *
 *  1. a valid tree must have exactly n-1 edges
 *  2. no cycle in graph
 *  3. full connectivity
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
func validTree(n int, edges [][]int) bool {
	if len(edges) != n-1 {
		return false
	}

	parent := make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
	}

	var find func(int) int
	find = func(a int) int {
		if parent[a] != a {
			return find(parent[a])
		}
		return parent[a]
	}

	union := func(a, b int) bool {
		rootA, rootB := find(a), find(b)
		if rootA == rootB {
			return false
		}

		parent[rootA] = rootB
		return true
	}

	for _, edge := range edges {
		if !union(edge[0], edge[1]) {
			return false
		}
	}
	return true
}

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

func validTreeDFS(n int, edges [][]int) bool {
	if len(edges) != n-1 {
		return false
	}

	// build undirected adjacency list
	adjList := make(map[int][]int)
	for _, edge := range edges {
		src, dst := edge[0], edge[1]
		adjList[src] = append(adjList[src], dst)
		adjList[dst] = append(adjList[dst], src)
	}

	visited := make([]bool, n)

	var dfs func(node, parent int) bool
	dfs = func(node, parent int) bool {
		visited[node] = true

		for _, neighbor := range adjList[node] {
			if neighbor == parent {
				// skip parent edge — not a cycle in undirected graph
				continue
			}
			if visited[neighbor] {
				// visited non-parent node = cycle detected
				return false
			}
			if !dfs(neighbor, node) {
				return false
			}
		}
		return true
	}

	// check no cycle from node 0
	if !dfs(0, -1) {
		return false
	}

	// check full connectivity
	for _, wasVisited := range visited {
		if !wasVisited {
			return false
		}
	}

	return true
}

func validTreeUF(n int, edges [][]int) bool {
	if len(edges) != n-1 {
		return false
	}

	// initialise each node as its own parent
	parent := make([]int, n)
	rank := make([]int, n)
	for i := range parent {
		parent[i] = i
	}

	// find root of node with path compression
	var find func(node int) int
	find = func(node int) int {
		if parent[node] != node {
			parent[node] = find(parent[node]) // path compression
		}
		return parent[node]
	}

	// union by rank — attach smaller tree under larger tree
	union := func(a, b int) bool {
		rootA, rootB := find(a), find(b)

		if rootA == rootB {
			// same root = cycle detected
			return false
		}

		// attach smaller rank tree under larger rank tree
		switch {
		case rank[rootA] > rank[rootB]:
			parent[rootB] = rootA
		case rank[rootA] < rank[rootB]:
			parent[rootA] = rootB
		default:
			parent[rootB] = rootA
			rank[rootA]++
		}
		return true
	}

	for _, edge := range edges {
		if !union(edge[0], edge[1]) {
			return false
		}
	}

	return true
}
