package interview

import (
	"strconv"
	"strings"
)

/**
 * Suspicious Activities — three implementations
 *
 * Each activity is a []string of F fields (e.g. [name, city, action]).
 * k = minimum number of fields that must match for two activities to be "similar".
 *
 * The result is every entry in newActivities reachable from any seed via BFS
 * over the similarity graph (transitivity applies).
 */

// ── Approach 1: BFS O(N² · F) ────────────────────────────────────────────────
//
// For every item popped from the queue, scan all N newActivities and call
// isSimilar (O(F) each). Simple but quadratic in the number of activities.

func findSuspiciousActivities(suspiciousActivities, newActivities [][]string, k int) [][]string {
	isSimilar := func(a, b []string) bool {
		matches := 0
		for i := range a {
			if a[i] == b[i] {
				matches++
			}
		}
		return matches >= k
	}

	n := len(newActivities)
	visited := make([]bool, n)
	var res [][]string

	queue := make([][]string, len(suspiciousActivities))
	copy(queue, suspiciousActivities)

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		for i := 0; i < n; i++ {
			if !visited[i] && isSimilar(cur, newActivities[i]) {
				visited[i] = true
				res = append(res, newActivities[i])
				queue = append(queue, newActivities[i])
			}
		}
	}
	return res
}

// ── Approach 2: BFS with inverted index O(N · C(F,k)) ────────────────────────
//
// Pre-index newActivities by every C(F,k) field-value combination.
// A frontier item now looks up C(F,k) map keys instead of scanning N activities.
// Each newActivity is visited at most once, so total work is O(N · C(F,k)).
//
// C(F,k) is a small constant for typical field counts (C(3,2) = 3).

func findSuspiciousActivitiesOpt(suspiciousActivities, newActivities [][]string, k int) [][]string {
	if len(newActivities) == 0 {
		return nil
	}

	combos := fieldCombinations(len(newActivities[0]), k)

	// Build inverted index: combo_key → []index into newActivities
	index := make(map[string][]int)
	for i, act := range newActivities {
		for _, combo := range combos {
			index[comboKey(act, combo)] = append(index[comboKey(act, combo)], i)
		}
	}

	n := len(newActivities)
	visited := make([]bool, n)
	var res [][]string

	queue := make([][]string, len(suspiciousActivities))
	copy(queue, suspiciousActivities)

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		for _, combo := range combos {
			for _, idx := range index[comboKey(cur, combo)] {
				if !visited[idx] {
					visited[idx] = true
					res = append(res, newActivities[idx])
					queue = append(queue, newActivities[idx])
				}
			}
		}
	}
	return res
}

// ── Approach 3: Union-Find with inverted index ────────────────────────────────
//
// Treat every activity (seeds + newActivities) as a node.
// Build the same inverted index over all nodes and union every pair that shares
// a k-field combination — this establishes connected components in one pass,
// no explicit BFS required.
//
// Then collect every newActivity whose component contains at least one seed.
//
// Seeds occupy indices 0..m-1; newActivities occupy indices m..m+n-1.

func findSuspiciousActivitiesUF(suspiciousActivities, newActivities [][]string, k int) [][]string {
	m, n := len(suspiciousActivities), len(newActivities)
	if n == 0 {
		return nil
	}

	// Merge seeds and newActivities into one slice for unified indexing
	all := make([][]string, m+n)
	copy(all[:m], suspiciousActivities)
	copy(all[m:], newActivities)

	combos := fieldCombinations(len(all[0]), k)

	// Build inverted index over ALL activities
	index := make(map[string][]int)
	for i, act := range all {
		for _, combo := range combos {
			index[comboKey(act, combo)] = append(index[comboKey(act, combo)], i)
		}
	}

	// Union-Find with path compression
	parent := make([]int, m+n)
	for i := range parent {
		parent[i] = i
	}
	var find func(int) int
	find = func(x int) int {
		if parent[x] != x {
			parent[x] = find(parent[x])
		}
		return parent[x]
	}
	union := func(x, y int) {
		if rx, ry := find(x), find(y); rx != ry {
			parent[rx] = ry
		}
	}

	// Union all activities that share a k-field combination
	for _, indices := range index {
		for j := 1; j < len(indices); j++ {
			union(indices[0], indices[j])
		}
	}

	// Collect the root of every seed component
	seedRoots := make(map[int]bool, m)
	for i := 0; i < m; i++ {
		seedRoots[find(i)] = true
	}

	// Every newActivity whose component root is a seed root is suspicious
	var res [][]string
	for i := m; i < m+n; i++ {
		if seedRoots[find(i)] {
			res = append(res, all[i])
		}
	}
	return res
}

// ── Shared helpers ────────────────────────────────────────────────────────────

// fieldCombinations returns all C(f, k) subsets of field indices 0..f-1.
func fieldCombinations(f, k int) [][]int {
	var result [][]int
	combo := make([]int, k)
	var generate func(start, depth int)
	generate = func(start, depth int) {
		if depth == k {
			c := make([]int, k)
			copy(c, combo)
			result = append(result, c)
			return
		}
		for i := start; i <= f-k+depth; i++ {
			combo[depth] = i
			generate(i+1, depth+1)
		}
	}
	generate(0, 0)
	return result
}

// comboKey builds a collision-safe map key for a subset of an activity's fields.
// Each value is prefixed with its field index so "city:London" and
// "name:London" never produce the same key.
func comboKey(activity []string, indices []int) string {
	var sb strings.Builder
	for i, fi := range indices {
		if i > 0 {
			sb.WriteByte('|')
		}
		sb.WriteString(strconv.Itoa(fi))
		sb.WriteByte(':')
		sb.WriteString(activity[fi])
	}
	return sb.String()
}
