package graph

/*
 * 323. Number of Connected Components in an Undirected Graph
 *
 * Time: O(n + m * α(n)) w/ path compression
 * Space: O(n)
 */
func countComponents(n int, edges [][]int) int {
	parent := make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
	}

	var find func(int) int
	find = func(x int) int {
		if parent[x] != x {
			parent[x] = find(parent[x])
		}
		return parent[x]
	}

	union := func(x, y int) bool {
		rootX, rootY := find(x), find(y)
		if rootX == rootY {
			return false
		}
		parent[rootX] = rootY
		return true
	}

	count := n
	for _, e := range edges {
		if union(e[0], e[1]) { // connect all the vertices via edges
			count--
		}
	}
	return count
}

func countComponentsBFS(n int, edges [][]int) int {
	// build adjacency list
	adjList := make(map[int][]int)
	for _, edge := range edges {
		adjList[edge[0]] = append(adjList[edge[0]], edge[1])
		adjList[edge[1]] = append(adjList[edge[1]], edge[0])
	}

	visited := make([]bool, n)
	components := 0

	for node := 0; node < n; node++ {
		if visited[node] {
			continue
		}
		// BFS from unvisited node — explores entire component
		components++
		queue := []int{node}
		visited[node] = true

		for len(queue) > 0 {
			cur := queue[0]
			queue = queue[1:]
			for _, neighbor := range adjList[cur] {
				if !visited[neighbor] {
					visited[neighbor] = true
					queue = append(queue, neighbor)
				}
			}
		}
	}

	return components
}

func countComponentsDFS(n int, edges [][]int) int {
	adjList := make(map[int][]int)
	for _, edge := range edges {
		adjList[edge[0]] = append(adjList[edge[0]], edge[1])
		adjList[edge[1]] = append(adjList[edge[1]], edge[0])
	}

	visited := make([]bool, n)

	var dfs func(node int)
	dfs = func(node int) {
		visited[node] = true
		for _, neighbor := range adjList[node] {
			if !visited[neighbor] {
				dfs(neighbor)
			}
		}
	}

	components := 0
	for node := 0; node < n; node++ {
		if !visited[node] {
			components++
			dfs(node)
		}
	}

	return components
}
