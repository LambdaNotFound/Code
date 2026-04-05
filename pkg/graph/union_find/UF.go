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
		// the union operation must update the root of one set to point to the root of the other set
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

/**
 * 128. Longest Consecutive Sequence
 *
 * HashSet is strictly simpler for this problem — Union Find shines when the problem has dynamic updates
 * (adding numbers one by one and querying max sequence length after each insertion).
 */
type UnionFind struct {
	parent map[int]int
	rank   map[int]int
	size   map[int]int // size of each component
}

func NewUnionFind() *UnionFind {
	return &UnionFind{
		parent: map[int]int{},
		rank:   map[int]int{},
		size:   map[int]int{},
	}
}

func (uf *UnionFind) Add(x int) {
	if _, exists := uf.parent[x]; !exists {
		uf.parent[x] = x
		uf.rank[x] = 0
		uf.size[x] = 1
	}
}

func (uf *UnionFind) Find(x int) int {
	if uf.parent[x] != x {
		uf.parent[x] = uf.Find(uf.parent[x]) // path compression
	}
	return uf.parent[x]
}

func (uf *UnionFind) Union(x, y int) {
	px, py := uf.Find(x), uf.Find(y)
	if px == py {
		return
	}
	// union by rank
	if uf.rank[px] < uf.rank[py] {
		px, py = py, px
	}
	uf.parent[py] = px
	uf.size[px] += uf.size[py]
	if uf.rank[px] == uf.rank[py] {
		uf.rank[px]++
	}
}

func (uf *UnionFind) MaxSize() int {
	best := 0
	for x, px := range uf.parent {
		if x == px { // only check roots
			best = max(best, uf.size[px])
		}
	}
	return best
}

func longestConsecutive(nums []int) int {
	if len(nums) == 0 {
		return 0
	}

	uf := NewUnionFind()
	numSet := map[int]bool{}

	for _, num := range nums {
		uf.Add(num)
		numSet[num] = true
	}

	for _, num := range nums {
		if numSet[num+1] {
			uf.Union(num, num+1) // union consecutive neighbors
		}
	}

	return uf.MaxSize()
}
