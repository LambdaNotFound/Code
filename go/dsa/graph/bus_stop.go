package graph

/**
 * 815. Bus Routes
 */
func numBusesToDestination(routes [][]int, source int, target int) int {
	stopToBus := make(map[int][]int)
	for i, route := range routes {
		for _, stop := range route {
			stopToBus[stop] = append(stopToBus[stop], i)
		}
	}

	queue := []int{source}
	visited := map[int]int{source: 1} // visited + tracking num of stops
	for len(queue) > 0 {
		f := queue[0]
		queue = queue[1:]

		for _, bus := range stopToBus[f] {
			for _, stop := range routes[bus] {
				if _, found := visited[stop]; found {
					continue
				}
				visited[stop] = visited[bus] + 1
				if stop == target {
					return visited[stop]
				}

				queue = append(queue, stop)
			}
		}
	}

	return -1
}

func numBusesToDestinationClaude(routes [][]int, source int, target int) int {
	if source == target {
		return 0
	}

	// stop → list of buses serving it
	stopToBus := make(map[int][]int)
	for i, route := range routes {
		for _, stop := range route {
			stopToBus[stop] = append(stopToBus[stop], i)
		}
	}

	// seed: all buses reachable from source
	visitedBus := map[int]bool{}
	queue := []int{}
	for _, bus := range stopToBus[source] {
		if !visitedBus[bus] {
			visitedBus[bus] = true
			queue = append(queue, bus)
		}
	}

	buses := 1 // already boarded first bus
	for len(queue) > 0 {
		size := len(queue)

		for i := 0; i < size; i++ {
			bus := queue[0]
			queue = queue[1:]

			for _, stop := range routes[bus] {
				if stop == target {
					return buses
				}
				// board all unvisited buses at this stop
				for _, nextBus := range stopToBus[stop] {
					if !visitedBus[nextBus] {
						visitedBus[nextBus] = true
						queue = append(queue, nextBus)
					}
				}
			}
		}
		buses++
	}
	return -1
}
