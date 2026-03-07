package graph

import "sort"

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

/**
 * 721. Accounts Merge
 */
type UnionFind struct {
	parent []int
	// rank   []int // performance guarantees by keeping trees shallow
}

func NewUnionFind(n int) *UnionFind {
	// parent, rank := make([]int, n), make([]int, n)
	parent := make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
	}
	return &UnionFind{
		parent: parent,
		// rank:   rank,
	}
}

func (uf *UnionFind) Find(node int) int {
	if uf.parent[node] != node {
		uf.parent[node] = uf.Find(uf.parent[node])
	}
	return uf.parent[node]
}

func (uf *UnionFind) Union(p, q int) bool {
	rootP, rootQ := uf.Find(p), uf.Find(q)
	if rootP == rootQ {
		return false
	}

	/*
		if uf.rank[rootP] < uf.rank[rootQ] {
			uf.parent[rootP] = rootQ
		} else if uf.rank[rootP] > uf.rank[rootQ] {
			uf.parent[rootQ] = rootP
		} else {
			uf.parent[rootQ] = rootP
			uf.rank[rootP]++
		}
	*/

	uf.parent[uf.Find(p)] = uf.Find(q)

	return true
}

func accountsMerge(accounts [][]string) [][]string {
	uf := NewUnionFind(len(accounts))
	emailToAccount := make(map[string]int)
	for i, account := range accounts {
		for _, email := range account[1:] {
			if acctID, exists := emailToAccount[email]; exists {
				uf.Union(i, acctID)
			} else {
				emailToAccount[email] = i
			}
		}
	}
	emailGroups := make(map[int][]string)
	for email, accountID := range emailToAccount {
		root := uf.Find(accountID)
		emailGroups[root] = append(emailGroups[root], email)
	}
	ret := make([][]string, 0, len(emailGroups))
	for acctID, emails := range emailGroups {
		sort.Strings(emails)
		account := append([]string{accounts[acctID][0]}, emails...)
		ret = append(ret, account)
	}
	return ret
}
