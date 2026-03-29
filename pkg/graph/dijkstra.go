package graph

import (
	"container/heap"
	"math"
)

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
	heap.Push(h, State{0, src, 0})

	// visited[node] = minimum stops to reach node at lowest cost
	visited := make([]int, n)
	for i := range visited {
		visited[i] = math.MaxInt
	}

	for h.Len() > 0 {
		cur := heap.Pop(h).(State)

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
			heap.Push(h, State{
				cost:  cur.cost + nextPrice,
				node:  nextNode,
				stops: cur.stops + 1,
			})
		}
	}
	return -1
}

type MinHeap []State

type State struct{ cost, node, stops int }

func (h MinHeap) Len() int            { return len(h) }
func (h MinHeap) Less(i, j int) bool  { return h[i].cost < h[j].cost }
func (h MinHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MinHeap) Push(x interface{}) { *h = append(*h, x.(State)) }
func (h *MinHeap) Pop() interface{} {
	x := (*h)[len(*h)-1]
	*h = (*h)[:len(*h)-1]
	return x
}
