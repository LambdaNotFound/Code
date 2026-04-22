package dijkstra

import (
	"container/heap"
	"math"
)

/**
 * Dijkstra, heap + cost function
 *
 * min-heap, minimizing the overall maximum
 *
 * max-heap, maximizing the overall minimum
 *    Path With Maximum Minimum Value (LC 1102)
 *    Maximum Probability Path (LC 1514)
 *
 * In standard Dijkstra, the state is just node —
 *    once the optimal cost to a node is finalized,
 *    any other entry for that node is guaranteed stale.
 *
 * prune both before & after push
 *
 * Space complexity:
 *    In the standard lazy deletion implementation (push duplicates, skip stale on pop),
 *    nodes can be pushed multiple times — once per incoming edge — giving O(E log E).
 *
 */
func dijkstra(graph [][][2]int, src int) []int {
	n := len(graph)
	dist := make([]int, n)
	for i := range dist {
		dist[i] = math.MaxInt
	}
	dist[src] = 0 // initialize dist[k] = 0 explicitly

	type state struct {
		cost, node int
	}

	// minHeap: [cost, node]
	h := &Heap[state]{
		less: func(a, b state) bool {
			return a.cost < b.cost
		},
	}
	heap.Push(h, state{0, src})

	for h.Len() > 0 {
		cur := heap.Pop(h).(state)
		cost, node := cur.cost, cur.node

		// prune on pop: stale entry — already found better path
		if cost > dist[node] { // strict >, to include first edge
			continue
		}

		for _, edge := range graph[node] {
			nei, weight := edge[0], edge[1]
			newCost := dist[node] + weight // ← sum of weights

			if newCost < dist[nei] { // prune on push: only push if improvement found
				dist[nei] = newCost
				heap.Push(h, state{newCost, nei})
			}
		}
	}
	return dist
}

/**
 * 787. Cheapest Flights Within K Stops
 *
 */
func findCheapestPrice(n int, flights [][]int, src int, dst int, k int) int {
	graph := make([][][2]int, n)
	for _, flight := range flights {
		from, to, price := flight[0], flight[1], flight[2]
		graph[from] = append(graph[from], [2]int{to, price})
	}

	type state struct{ cost, node, stops int }

	h := &Heap[state]{
		less: func(a, b state) bool {
			return a.cost < b.cost
		},
	}
	heap.Push(h, state{0, src, 0})

	// visited[node] = minimum stops to reach node at lowest cost
	visited := make([]int, n)
	for i := range visited {
		visited[i] = math.MaxInt
	}

	for h.Len() > 0 {
		cur := heap.Pop(h).(state)

		if cur.node == dst {
			return cur.cost
		}

		if cur.stops > k {
			continue
		}

		visited[cur.node] = cur.stops

		for _, nei := range graph[cur.node] {
			nextNode, nextPrice := nei[0], nei[1]
			newPrice, newStops := cur.cost+nextPrice, cur.stops+1

			if visited[nextNode] <= newStops {
				continue // pruning: already visited with fewer stops
			}

			heap.Push(h, state{
				cost:  newPrice,
				node:  nextNode,
				stops: newStops,
			})
		}
	}
	return -1
}

/**
 * 743. Network Delay Time
 *
 * Given a network of n nodes (1-indexed) and directed weighted edges in
 * `times` ([src, dst, weight]), find the minimum time for a signal sent
 * from node k to reach ALL nodes. Return -1 if any node is unreachable.
 *
 * Time:
 *   build adjList:    O(E)
 *   heap operations:  each edge can trigger at most one push → at most E items in heap
 *   each push/pop:    O(log E)
 *   ─────────────────────────────────────────
 *   total:            O(E log E)
 *
 *
 * Space:
 *   dist array:    O(V)
 *   adjList:       O(V + E)
 *   heap:          O(E)   — at most E entries
 *   ─────────────────────────────────────────
 *   total:         O(V + E)
 */
func networkDelayTime(times [][]int, n int, k int) int {
	dist := make([]int, n+1)
	for i := range dist {
		dist[i] = math.MaxInt
	}
	dist[k] = 0

	adjList := make([][][2]int, n+1)
	for _, time := range times {
		src, dst, weight := time[0], time[1], time[2]
		adjList[src] = append(adjList[src], [2]int{dst, weight})
	}

	type state struct{ node, cost int }
	minHeap := &Heap[state]{
		less: func(a, b state) bool { return a.cost < b.cost },
	}
	heap.Push(minHeap, state{k, 0})

	for minHeap.Len() > 0 {
		top := heap.Pop(minHeap).(state)

		if top.cost > dist[top.node] {
			continue
		}

		for _, edge := range adjList[top.node] {
			nei, weight := edge[0], edge[1]
			newCost := top.cost + weight

			if newCost < dist[nei] {
				dist[nei] = newCost
				heap.Push(minHeap, state{nei, newCost})
			}
		}
	}

	res := 0
	for _, d := range dist[1:] {
		if d == math.MaxInt {
			return -1
		}
		res = max(res, d)
	}
	return res
}

/**
 * 1631. Path With Minimum Effort
 *
 * A route's effort is the maximum absolute difference in heights between two consecutive cells of the route.
 *
 * 1. adjGrid [][]int
 * 2. dist[]
 *
 */
func minimumEffortPath(heights [][]int) int {
	m, n := len(heights), len(heights[0])
	dist := make([][]int, m)
	for i := range dist {
		dist[i] = make([]int, n)

		for j := range dist[i] {
			dist[i][j] = math.MaxInt
		}
	}
	dist[0][0] = 0

	type state struct{ row, col, cost int }
	minHeap := &Heap[state]{
		less: func(a, b state) bool {
			return a.cost < b.cost
		},
	}
	heap.Push(minHeap, state{0, 0, 0})

	abs := func(a int) int {
		if a < 0 {
			return -a
		}
		return a
	}
	dirs := [][2]int{{0, 1}, {1, 0}, {-1, 0}, {0, -1}}
	for minHeap.Len() > 0 {
		cur := heap.Pop(minHeap).(state)
		row, col := cur.row, cur.col

		if cur.cost > dist[cur.row][cur.col] {
			continue
		}

		for _, dir := range dirs {
			r, c := row+dir[0], col+dir[1]
			if r < 0 || c < 0 || r >= m || c >= n {
				continue
			}

			newCost := max(cur.cost, abs(heights[r][c]-heights[row][col]))
			if newCost < dist[r][c] {
				dist[r][c] = newCost
				heap.Push(minHeap, state{r, c, newCost})
			}
		}
	}

	return dist[m-1][n-1]
}

/**
 * 1514. Path with Maximum Probability
 *
 * 1. adjList [][]edge
 * 2. prob[]
 *
 * Time: O((V + E) log V) — standard Dijkstra with a heap.
 * Space: O(V + E) — graph adjacency list + prob array + heap.
 *
 */
func maxProbability(n int, edges [][]int, succProb []float64, start_node int, end_node int) float64 {
	prob := make([]float64, n)
	prob[start_node] = 1

	type State struct {
		node int
		cost float64
	}

	maxHeap := &Heap[State]{
		less: func(a, b State) bool {
			return a.cost > b.cost
		},
	}
	heap.Push(maxHeap, State{node: start_node, cost: 1})

	type edge struct {
		node int
		cost float64
	}
	graph := make([][]edge, n)
	for i, e := range edges {
		v1, v2 := e[0], e[1]
		cost := succProb[i]
		graph[v1] = append(graph[v1], edge{v2, cost})
		graph[v2] = append(graph[v2], edge{v1, cost})
	}

	for maxHeap.Len() > 0 {
		cur := heap.Pop(maxHeap).(State)

		if cur.cost < prob[cur.node] {
			continue
		}
		prob[cur.node] = cur.cost

		for _, edge := range graph[cur.node] {
			newCost := cur.cost * edge.cost
			if newCost > prob[edge.node] {
				heap.Push(maxHeap, State{node: edge.node, cost: newCost})
			}
		}
	}

	return prob[end_node]
}
