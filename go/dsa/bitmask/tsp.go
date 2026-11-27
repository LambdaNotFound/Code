package bitmask

import "math"

/*
 * Variant 1: Hamiltonian Path (open) — Visit every city exactly once, minimize total distance. You stop at the last city. [COMMON]
 *
 */
func tspOpen(dist [][]int) int {
	n := len(dist)
	full := 1 << n

	dp := make([][]int, full)
	for m := range dp {
		dp[m] = make([]int, n)
		for i := range dp[m] {
			dp[m][i] = math.MaxInt64
		}
	}
	dp[1<<0][0] = 0

	for mask := 1; mask < full; mask++ {
		for i := 0; i < n; i++ {
			if dp[mask][i] == math.MaxInt64 || mask&(1<<i) == 0 {
				continue
			}
			for j := 0; j < n; j++ {
				if mask&(1<<j) != 0 {
					continue
				}
				newMask := mask | (1 << j)
				cost := dp[mask][i] + dist[i][j]
				if cost < dp[newMask][j] {
					dp[newMask][j] = cost
				}
			}
		}
	}

	// OPEN: just find the best ending city
	ans := math.MaxInt64
	allVisited := full - 1
	for i := 0; i < n; i++ {
		if dp[allVisited][i] < ans {
			ans = dp[allVisited][i]
		}
	}
	return ans
}

/*
 * Variant 2: Hamiltonian Cycle (closed) — Visit every city exactly once and return to the starting city, minimize total distance. This is the classical TSP. [COMMON]
 *
 */
func tspClosed(dist [][]int) int {
	n := len(dist)
	full := 1 << n

	dp := make([][]int, full)
	for m := range dp {
		dp[m] = make([]int, n)
		for i := range dp[m] {
			dp[m][i] = math.MaxInt64
		}
	}
	dp[1<<0][0] = 0

	for mask := 1; mask < full; mask++ {
		for i := 0; i < n; i++ {
			if dp[mask][i] == math.MaxInt64 || mask&(1<<i) == 0 {
				continue
			}
			for j := 0; j < n; j++ {
				if mask&(1<<j) != 0 {
					continue
				}
				newMask := mask | (1 << j)
				cost := dp[mask][i] + dist[i][j]
				if cost < dp[newMask][j] {
					dp[newMask][j] = cost
				}
			}
		}
	}

	// CLOSED: best ending city + cost to return to city 0
	ans := math.MaxInt64
	allVisited := full - 1
	for i := 0; i < n; i++ {
		if dp[allVisited][i] == math.MaxInt64 {
			continue
		}
		cost := dp[allVisited][i] + dist[i][0]
		if cost < ans {
			ans = cost
		}
	}
	return ans
}
