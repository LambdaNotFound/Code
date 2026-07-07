package bellmanford

import "math"

/**
 * Bellman-Ford, relax all edges V-1 times
 *
 * Handles negative edge weights — unlike Dijkstra which requires non-negative weights.
 * Detects negative-weight cycles: if a Vth relaxation round still improves any distance,
 * a negative cycle is reachable from src.
 *
 * Compared to Dijkstra:
 *    Dijkstra  O(E log E) time, O(V + E) space — greedy, non-negative weights only
 *    Bellman-Ford  O(V·E) time, O(V) space    — DP, handles negative weights + cycle detection
 *
 * When to prefer Bellman-Ford over Dijkstra:
 *    - Graph has negative edge weights
 *    - Need to detect a negative-weight cycle
 *    - K-step constraint (e.g. LC 787) — fix exactly K rounds of relaxation
 */
func bellmanFord(graph [][][2]int, src int) (dist []int, hasNegCycle bool) {
	n := len(graph)
	dist = make([]int, n)
	for i := range dist {
		dist[i] = math.MaxInt
	}
	dist[src] = 0

	// Relax all edges V-1 times.
	// After round i, dist[v] holds the shortest path using at most i+1 edges.
	for i := 0; i < n-1; i++ {
		updated := false
		for node := range graph {
			if dist[node] == math.MaxInt {
				continue // unreachable — skip to avoid overflow on addition
			}
			for _, edge := range graph[node] {
				nei, weight := edge[0], edge[1]
				if newCost := dist[node] + weight; newCost < dist[nei] {
					dist[nei] = newCost
					updated = true
				}
			}
		}
		if !updated {
			break // early exit: stable — no further improvement possible
		}
	}

	// Vth round: if any distance still improves, a negative cycle is reachable.
	for node := range graph {
		if dist[node] == math.MaxInt {
			continue
		}
		for _, edge := range graph[node] {
			nei, weight := edge[0], edge[1]
			if dist[node]+weight < dist[nei] {
				return nil, true
			}
		}
	}

	return dist, false
}

/**
 * 787. Cheapest Flights Within K Stops — Bellman-Ford variant
 *
 * Run exactly K+1 relaxation rounds (src counts as 0 stops, so K stops = K+1 edges).
 * Snapshot the previous round's dist before each pass so updates within a single round
 * don't chain — this enforces the "at most one new edge per round" invariant.
 *
 * Time:  O(K · E)   — K+1 rounds, each scanning all E edges
 * Space: O(V)       — dist + one snapshot copy
 */
func findCheapestPrice(n int, flights [][]int, src int, dst int, k int) int {
	dist := make([]int, n)
	for i := range dist {
		dist[i] = math.MaxInt
	}
	dist[src] = 0

	for i := 0; i <= k; i++ {
		prev := append([]int{}, dist...) // snapshot: only extend paths found in prior rounds
		for _, flight := range flights {
			from, to, price := flight[0], flight[1], flight[2]
			if prev[from] == math.MaxInt {
				continue
			}
			if newCost := prev[from] + price; newCost < dist[to] {
				dist[to] = newCost
			}
		}
	}

	if dist[dst] == math.MaxInt {
		return -1
	}
	return dist[dst]
}
