package topologicalsort

/**
 * 310. Minimum Height Trees | NOT classic topo sort, Tree center = optimal root
 *
 * Given a tree of n nodes labelled from 0 to n - 1,
 * Return a list of all MHTs' root labels.
 *
 * The height of a rooted tree is the number of edges on the
 * longest downward path between the root and a leaf.
 *
 */
func findMinHeightTrees(n int, edges [][]int) []int {
	if n == 1 {
		return []int{0}
	}

	degree := make([]int, n)
	adjList := make([][]int, n)
	for _, edge := range edges {
		degree[edge[0]]++
		degree[edge[1]]++
		adjList[edge[0]] = append(adjList[edge[0]], edge[1])
		adjList[edge[1]] = append(adjList[edge[1]], edge[0])
	}

	queue := []int{}
	for i := range degree {
		if degree[i] == 1 {
			queue = append(queue, i)
		}
	}

	numOfNodesLeft := n
	for len(queue) < numOfNodesLeft {
		size := len(queue)
		numOfNodesLeft = numOfNodesLeft - size

		for i := 0; i < size; i++ {
			node := queue[0]
			queue = queue[1:]

			for _, neighbor := range adjList[node] {
				degree[neighbor]--
				if degree[neighbor] == 1 {
					queue = append(queue, neighbor)
				}
			}
		}
	}

	return queue
}

func findMinHeightTreesTwoPassBFS(n int, edges [][]int) []int {
	if n == 1 {
		return []int{0}
	}

	// build adjacency list
	graph := make([][]int, n)
	for _, edge := range edges {
		graph[edge[0]] = append(graph[edge[0]], edge[1])
		graph[edge[1]] = append(graph[edge[1]], edge[0])
	}

	// BFS from a given start, returns (farthest node, parent map)
	bfs := func(start int) (int, []int) {
		dist := make([]int, n)
		for i := range dist {
			dist[i] = -1
		}
		parent := make([]int, n)
		for i := range parent {
			parent[i] = -1
		}

		dist[start] = 0
		queue := []int{start}
		farthest := start

		for len(queue) > 0 {
			node := queue[0]
			queue = queue[1:]
			if dist[node] > dist[farthest] {
				farthest = node
			}
			for _, neighbor := range graph[node] {
				if dist[neighbor] == -1 {
					dist[neighbor] = dist[node] + 1
					parent[neighbor] = node
					queue = append(queue, neighbor)
				}
			}
		}
		return farthest, parent
	}

	// pass 1: find one endpoint of the diameter
	endpoint1, _ := bfs(0)

	// pass 2: find other endpoint + track path
	endpoint2, parent := bfs(endpoint1)

	// reconstruct diameter path from endpoint2 back to endpoint1
	path := []int{}
	for node := endpoint2; node != -1; node = parent[node] {
		path = append(path, node)
	}

	// centers are middle node(s) of diameter path
	pathLen := len(path)
	if pathLen%2 == 1 {
		return []int{path[pathLen/2]}
	}
	return []int{path[pathLen/2-1], path[pathLen/2]}
}
