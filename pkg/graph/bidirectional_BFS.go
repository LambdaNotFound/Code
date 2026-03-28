package graph

/**
 * Minimum Knight Moves
 *
 * 1. TWO FRONTIERS    — current layer being expanded on each side
 * 2. TWO VISITED MAPS — all nodes seen so far + their distance
 * 3. EXPAND SMALLER   — keeps search balanced, prevents explosion
 * 4. MEETING CHECK    — when frontier A discovers a node in visitB
 */
func minKnightMoves(x, y int) int {
	// normalize to first quadrant — knight moves are symmetric
	x, y = abs(x), abs(y)

	if x == 0 && y == 0 {
		return 0
	}

	dirs := [][2]int{
		{1, 2}, {2, 1}, {-1, 2}, {-2, 1},
		{1, -2}, {2, -1}, {-1, -2}, {-2, -1},
	}

	// frontier: set of (x,y) positions
	// visited:  map of (x,y) → steps from that side
	frontF := map[[2]int]int{{0, 0}: 0} // forward frontier
	frontB := map[[2]int]int{{x, y}: 0} // backward frontier
	visitF := map[[2]int]int{{0, 0}: 0} // forward visited
	visitB := map[[2]int]int{{x, y}: 0} // backward visited

	steps := 0

	for len(frontF) > 0 && len(frontB) > 0 {
		steps++

		// always expand the smaller frontier
		if len(frontF) > len(frontB) {
			frontF, frontB = frontB, frontF
			visitF, visitB = visitB, visitF
		}

		nextFront := map[[2]int]int{}
		for pos := range frontF {
			for _, d := range dirs {
				next := [2]int{pos[0] + d[0], pos[1] + d[1]}

				if _, seen := visitF[next]; seen {
					continue
				}

				// meeting point found!
				if stepsB, found := visitB[next]; found {
					return visitF[pos] + 1 + stepsB
				}

				nextFront[next] = visitF[pos] + 1
				visitF[next] = visitF[pos] + 1
			}
		}
		frontF = nextFront
	}
	return -1
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
