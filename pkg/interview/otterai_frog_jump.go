package interview

/**
 * 403. Frog Jump
 */
// Time: O(n²) — the jump size after k steps is bounded by O(n) (it increases by at most 1 each stone-hop, and there are at most n stones)
// Space: O(n²) — for the visited map and queue
func canCross(stones []int) bool {
	n := len(stones)

	stoneIndex := map[int]int{}
	for i, pos := range stones {
		stoneIndex[pos] = i
	}
	target := stones[n-1]

	type state struct {
		pos, jump int
	}

	visited := make(map[state]bool)
	queue := []state{{stones[0], 0}}
	visited[state{stones[0], 0}] = true

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		for _, delta := range []int{-1, 0, 1} {
			nextJump := curr.jump + delta
			nextPos := curr.pos + nextJump

			if _, exists := stoneIndex[nextPos]; !exists {
				continue
			}

			if nextPos == target {
				return true
			}

			nextState := state{nextPos, nextJump}

			if visited[nextState] == true {
				continue
			}
			visited[nextState] = true
			queue = append(queue, nextState)
		}
	}

	return false
}
