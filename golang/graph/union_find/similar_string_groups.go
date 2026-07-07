package unionfind

/**
 * 839. Similar String Groups
 *
 * 1. BFS
 *     a). optimized isSimilar(a, b string)
 *     b). index based visited array, e.g. strs = ["abc","abc"]
 *
 * 2. union-find
 *
 */
func numSimilarGroups(strs []string) int {
	isSimilar := func(a, b string) bool {
		diffs := 0
		for i := 0; i < len(a); i++ {
			if a[i] != b[i] {
				diffs++
				if diffs > 2 {
					return false // early exit
				}
			}
		}
		return diffs == 0 || diffs == 2 // identical or exactly one swap apart
	}

	n := len(strs)
	visited := make([]bool, n)
	groups := 0

	for i := 0; i < n; i++ {
		if visited[i] {
			continue
		}
		groups++
		queue := []int{i}
		visited[i] = true

		for len(queue) > 0 {
			curIdx := queue[0]
			queue = queue[1:]

			for j := 0; j < n; j++ {
				if !visited[j] && isSimilar(strs[curIdx], strs[j]) {
					visited[j] = true
					queue = append(queue, j)
				}
			}
		}
	}
	return groups
}

func numSimilarGroupsUF(strs []string) int {
	isSimilar := func(a, b string) bool {
		diffs := 0
		for i := 0; i < len(a); i++ {
			if a[i] != b[i] {
				diffs++
				if diffs > 2 {
					return false
				}
			}
		}
		return diffs == 0 || diffs == 2
	}

	n := len(strs)
	parent := make([]int, n)

	for i := range parent {
		parent[i] = i
	}

	var find func(x int) int
	find = func(x int) int {
		if x != parent[x] {
			parent[x] = find(parent[x]) // path compression
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

	groups := n
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if isSimilar(strs[i], strs[j]) {
				if union(i, j) {
					groups-- // one fewer independent group after merging
				}
			}
		}
	}

	return groups
}
