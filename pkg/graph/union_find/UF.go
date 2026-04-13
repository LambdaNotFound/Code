package unionfind

/**
 * Union Find Template
 */
func unionFind(n int, graph [][2]int) {
	parent := make([]int, n)
	for i := range parent {
		parent[i] = i
	}

	var find func(x int) int
	find = func(x int) int {
		if x != parent[x] {
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

	for _, edge := range graph {
		union(edge[0], edge[1])
	}
}

func unionFindByRank(n int, graph [][2]int) {
	parent := make([]int, n)
	for i := range parent {
		parent[i] = i
	}

	rank := make([]int, n)

	var find func(x int) int
	find = func(x int) int {
		if x != parent[x] {
			parent[x] = find(parent[x])
		}
		return parent[x]
	}

	union := func(x, y int) bool {
		rootX, rootY := find(x), find(y)
		if rootX == rootY {
			return false
		}
		if rank[rootX] < rank[rootY] {
			parent[rootX] = rootY
		} else if rank[rootX] > rank[rootY] {
			parent[rootY] = rootX
		} else {
			parent[rootY] = rootX
			rank[rootX]++
		}
		return true
	}

	for _, edge := range graph {
		union(edge[0], edge[1])
	}
}

func unionFindBySize(n int, graph [][2]int) {
	parent := make([]int, n)
	for i := range parent {
		parent[i] = i
	}

	size := make([]int, n)
	for i := range size {
		size[i] = 1
	}

	var find func(x int) int
	find = func(x int) int {
		if x != parent[x] {
			parent[x] = find(parent[x])
		}
		return parent[x]
	}

	union := func(x, y int) bool {
		rootX, rootY := find(x), find(y)
		if rootX == rootY {
			return false
		}
		if size[rootX] < size[rootY] {
			parent[rootX] = rootY
			size[rootY] += size[rootX]
		} else {
			parent[rootY] = rootX
			size[rootX] += size[rootY]
		}
		return true
	}

	for _, edge := range graph {
		union(edge[0], edge[1])
	}
}

/*
 * 323. Number of Connected Components in an Undirected Graph
 *
 * Space: O(n)
 * initialize parent/size:   O(n)
 * process each edge:        O(e · α(n))   ← e unions, each costs α(n) amortized
 *
 * total:                    O(n + e · α(n)), in practice this is O(n + e).
 *
 * Amortized find cost, union by rank/size onlyO(log n)
 */
func countComponents(n int, edges [][]int) int {
	parent := make([]int, n)
	for i := range parent {
		parent[i] = i
	}
	size := make([]int, n)
	for i := range size {
		size[i] = 1
	}

	var find func(x int) int
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
		if size[rootX] < size[rootY] {
			parent[rootX] = rootY
			size[rootY] += size[rootX]
		} else {
			parent[rootY] = rootX
			size[rootX] += size[rootY]
		}
		return true
	}

	components := n
	for _, edge := range edges {
		if union(edge[0], edge[1]) {
			components--
		}
	}
	return components
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
