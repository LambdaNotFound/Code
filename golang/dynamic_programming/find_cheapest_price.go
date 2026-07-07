package dynamic_programming

import "math"

/**
 * 787. Cheapest Flights Within K Stops
 *
 * Bellman-Ford is essentially DP.
 */
func findCheapestPrice(n int, flights [][]int, src int, dst int, k int) int {
	const INF = math.MaxInt / 2

	// dp[i][v] = min cost to reach v using at most i edges
	dp := make([][]int, k+2)
	for i := range dp {
		dp[i] = make([]int, n)
		for j := range dp[i] {
			dp[i][j] = INF
		}
	}
	dp[0][src] = 0

	for i := 1; i <= k+1; i++ {
		// carry forward: can always stay at same cost from previous round
		copy(dp[i], dp[i-1])

		for _, flight := range flights {
			from, to, price := flight[0], flight[1], flight[2]
			if dp[i-1][from] == INF {
				continue
			}
			// relax edge using PREVIOUS round's values
			dp[i][to] = min(dp[i][to], dp[i-1][from]+price)
		}
	}

	if dp[k+1][dst] == INF {
		return -1
	}
	return dp[k+1][dst]
}
