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
 */
func dijkstra(graph [][][2]int, src int) []int {
	n := len(graph)
	dist := make([]int, n)
	for i := range dist {
		dist[i] = math.MaxInt
	}
	dist[src] = 0

	type Item struct {
		cost, node int
	}

	// minHeap: [cost, node]
	h := NewHeap(func(a, b Item) bool {
		return a.cost < b.cost
	})
	heap.Push(h, Item{0, src})

	for h.Len() > 0 {
		cur := heap.Pop(h).(Item)
		cost, node := cur.cost, cur.node

		// stale entry — already found better path
		if cost > dist[node] {
			continue
		}

		for _, edge := range graph[node] {
			nei, weight := edge[0], edge[1]
			newCost := dist[node] + weight // ← sum of weights

			if newCost < dist[nei] {
				dist[nei] = newCost
				heap.Push(h, Item{newCost, nei})
			}
		}
	}
	return dist
}

/**
 * 787. Cheapest Flights Within K Stops
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

		// pruning: already visited with fewer stops
		if visited[cur.node] <= cur.stops {
			continue
		}
		visited[cur.node] = cur.stops

		for _, nei := range graph[cur.node] {
			nextNode, nextPrice := nei[0], nei[1]
			heap.Push(h, Flight{
				cost:  cur.cost + nextPrice,
				node:  nextNode,
				stops: cur.stops + 1,
			})
		}
	}
	return -1
}
