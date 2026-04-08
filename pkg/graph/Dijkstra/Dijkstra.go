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

	// minHeap: [cost, node]
	h := &Heap{
		less: func(a, b Item) bool {
			return a.cost < b.cost
		},
	}
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

type Item struct {
	cost, node int
}

type Heap struct {
	items []Item
	less  func(a, b Item) bool
}

func (h *Heap) Len() int           { return len(h.items) }
func (h *Heap) Swap(i, j int)      { h.items[i], h.items[j] = h.items[j], h.items[i] }
func (h *Heap) Less(i, j int) bool { return h.less(h.items[i], h.items[j]) }

func (h *Heap) Pop() interface{} {
	item := h.items[len(h.items)-1]
	h.items = h.items[:len(h.items)-1]
	return item
}

func (h *Heap) Push(item interface{}) {
	h.items = append(h.items, item.(Item))
}

/**
 * 787. Cheapest Flights Within K Stops
 */
func findCheapestPrice(n int, flights [][]int, src int, dst int, k int) int {
	// build adjacency list: graph[u] = [(v, price)]
	graph := make([][][2]int, n)
	for _, flight := range flights {
		from, to, price := flight[0], flight[1], flight[2]
		graph[from] = append(graph[from], [2]int{to, price})
	}

	// minHeap entry: [cost, node, stops]
	h := &MinHeap{}
	heap.Push(h, StateFlight{0, src, 0})

	// visited[node] = minimum stops to reach node at lowest cost
	visited := make([]int, n)
	for i := range visited {
		visited[i] = math.MaxInt
	}

	for h.Len() > 0 {
		cur := heap.Pop(h).(StateFlight)

		// reached destination
		if cur.node == dst {
			return cur.cost
		}

		// pruning: too many stops
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
			heap.Push(h, StateFlight{
				cost:  cur.cost + nextPrice,
				node:  nextNode,
				stops: cur.stops + 1,
			})
		}
	}
	return -1
}

type MinHeapFlights []StateFlight

type StateFlight struct{ cost, node, stops int }

func (h MinHeapFlights) Len() int            { return len(h) }
func (h MinHeapFlights) Less(i, j int) bool  { return h[i].cost < h[j].cost }
func (h MinHeapFlights) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MinHeapFlights) Push(x interface{}) { *h = append(*h, x.(StateFlight)) }
func (h *MinHeapFlights) Pop() interface{} {
	x := (*h)[len(*h)-1]
	*h = (*h)[:len(*h)-1]
	return x
}
