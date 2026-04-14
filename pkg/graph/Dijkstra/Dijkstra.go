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
 */
func dijkstra(graph [][][2]int, src int) []int {
	n := len(graph)
	dist := make([]int, n)
	for i := range dist {
		dist[i] = math.MaxInt
	}
	dist[src] = 0 // initialize dist[k] = 0 explicitly

	type Item struct {
		cost, node int
	}

	// minHeap: [cost, node]
	h := &Heap[Item]{
		less: func(a, b Item) bool {
			return a.cost < b.cost
		},
	}
	heap.Push(h, Item{0, src})

	for h.Len() > 0 {
		cur := heap.Pop(h).(Item)
		cost, node := cur.cost, cur.node

		// prune on pop: stale entry — already found better path
		if cost > dist[node] {
			continue
		}

		for _, edge := range graph[node] {
			nei, weight := edge[0], edge[1]
			newCost := dist[node] + weight // ← sum of weights

			if newCost < dist[nei] { // prune on push: only push if improvement found
				dist[nei] = newCost
				heap.Push(h, Item{newCost, nei})
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

	type Flight struct{ cost, node, stops int }

	h := &Heap[Flight]{
		less: func(a, b Flight) bool {
			return a.cost < b.cost
		},
	}
	heap.Push(h, Flight{0, src, 0})

	// visited[node] = minimum stops to reach node at lowest cost
	visited := make([]int, n)
	for i := range visited {
		visited[i] = math.MaxInt
	}

	for h.Len() > 0 {
		cur := heap.Pop(h).(Flight)

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

			heap.Push(h, Flight{
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
