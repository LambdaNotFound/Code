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
