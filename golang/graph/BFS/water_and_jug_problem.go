package bfs

/**
 * 365. Water and Jug Problem
 *
 */
func canMeasureWater(jug1Cap int, jug2Cap int, target int) bool {
	type state struct{ x, y int }

	visited := make(map[state]bool)
	queue := []state{{0, 0}}
	visited[state{0, 0}] = true

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		x, y := cur.x, cur.y

		if x == target || y == target || x+y == target {
			return true
		}

		pour1to2 := min(x, jug2Cap-y)
		pour2to1 := min(y, jug1Cap-x)

		candidates := []state{
			{jug1Cap, y},                 // fill jug1
			{x, jug2Cap},                 // fill jug2
			{0, y},                       // empty jug1  ← was missing
			{x, 0},                       // empty jug2  ← was missing
			{x - pour1to2, y + pour1to2}, // pour jug1 → jug2
			{x + pour2to1, y - pour2to1}, // pour jug2 → jug1
		}

		for _, next := range candidates {
			if !visited[next] {
				visited[next] = true
				queue = append(queue, next)
			}
		}
	}
	return false
}
