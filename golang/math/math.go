package math

/**
 * 1884. Egg Drop With 2 Eggs and N Floors
 *
 * 2 eggs, Egg 1 is your "probe" — once it breaks, it's gone Egg 2 must be used linearly from the last safe floor upward
 * — because you can't afford to break it
 *
 */
func twoEggDrop(n int) int {
	// find minimum t such that t*(t+1)/2 >= n
	t := 1
	for t*(t+1)/2 < n {
		t++
	}
	return t
}

/**
 * 1828. Queries on Number of Points Inside a Circle
 *
 * Time: O(m · n)
 * Space: O(n)
 *
 * follow ups:
 *     m >> n (many points, few circles) — KD-tree
 *     m << n (few points, many circles) — sort + binary search on x
 */
func countPoints(points [][]int, queries [][]int) []int {
	answer := make([]int, len(queries))

	for queryIndex, query := range queries {
		centerX, centerY, radius := query[0], query[1], query[2]
		radiusSquared := radius * radius

		count := 0
		for _, point := range points {
			deltaX := point[0] - centerX
			deltaY := point[1] - centerY
			// Compare squared distance against squared radius so we stay in
			// exact integer math and avoid a floating-point sqrt. "<=" because
			// points on the boundary count as inside.
			if deltaX*deltaX+deltaY*deltaY <= radiusSquared {
				count++
			}
		}
		answer[queryIndex] = count
	}

	return answer
}
