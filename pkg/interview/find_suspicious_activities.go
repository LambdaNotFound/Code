package interview

/**
 * Suspicious Activites
 *
 */
func findSuspiciousActivities(suspiciousActivites, newActivities [][]string, k int) [][]string {
	isSimilar := func(a, b []string, k int) bool {
		matches := 0
		for i := 0; i < len(a); i++ {
			if a[i] == b[i] {
				matches++
			}
		}
		return matches >= k
	}

	res := [][]string{}
	n := len(newActivities)
	visited := make([]bool, n)

	// Seeds drive BFS propagation but are NOT added to res —
	// only newActivities discovered during BFS belong in the output.
	queue := make([][]string, len(suspiciousActivites))
	copy(queue, suspiciousActivites)

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		for i := 0; i < n; i++ {
			if !visited[i] && isSimilar(cur, newActivities[i], k) {
				visited[i] = true
				res = append(res, newActivities[i])
				queue = append(queue, newActivities[i])
			}
		}
	}

	return res
}
