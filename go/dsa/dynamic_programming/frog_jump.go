package dynamic_programming

/**
 * 403. Frog Jump
 *
 * Given a list of stones positions (in units) in sorted ascending order, determine if the frog can cross the river by landing on the last stone. Initially, the frog is on the first stone and assumes the first jump must be 1 unit.
 *
 * If the frog's last jump was k units, its next jump must be either k - 1, k, or k + 1 units. The frog can only jump in the forward direction.
 *
 */
func canCross(stones []int) bool {
	n := len(stones)
	if n == 2 {
		return stones[1]-stones[0] == 1
	}

	stoneIndex := make(map[int]int, n)
	for i, pos := range stones {
		stoneIndex[pos] = i
	}

	// dp[i] = set of jump sizes that can land you on stone i
	dp := make([]map[int]bool, n)
	for i := range dp {
		dp[i] = make(map[int]bool)
	}
	dp[0][0] = true

	for i := 0; i < n; i++ {
		for jump := range dp[i] {
			for _, delta := range []int{-1, 0, 1} {
				nextJump := jump + delta
				if nextJump <= 0 {
					continue
				}
				nextPos := stones[i] + nextJump
				if nextIdx, exists := stoneIndex[nextPos]; exists {
					dp[nextIdx][nextJump] = true
				}
			}
		}
	}

	return len(dp[n-1]) > 0
}

// Time: O(n²) Total states across all stones: Σ (i=1 to n) O(i) = O(1+2+3+...+n) = O(n²/2) = O(n²)
// Space: O(n²) (visited + queue)
func canCrossBFS(stones []int) bool {
	n := len(stones)
	if n == 2 {
		return stones[1]-stones[0] == 1
	}

	stoneIndex := make(map[int]int, n)
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
		cur := queue[0]
		queue = queue[1:]

		if cur.pos == target {
			return true
		}

		for _, delta := range []int{-1, 0, 1} {
			nextJump := cur.jump + delta
			if nextJump <= 0 {
				continue
			}

			nextPos := cur.pos + nextJump
			if _, exists := stoneIndex[nextPos]; !exists {
				continue
			}

			next := state{nextPos, nextJump}
			if visited[next] {
				continue
			}
			visited[next] = true
			queue = append(queue, next)
		}
	}

	return false
}

func canCrossDFS(stones []int) bool {
	n := len(stones)
	if n == 2 {
		return stones[1]-stones[0] == 1
	}

	stoneIndex := make(map[int]int, n)
	for i, pos := range stones {
		stoneIndex[pos] = i
	}

	target := stones[n-1]
	memo := make(map[[2]int]bool)

	var dfs func(pos, jump int) bool
	dfs = func(pos, jump int) bool {
		if pos == target {
			return true
		}

		key := [2]int{pos, jump}
		if val, ok := memo[key]; ok {
			return val
		}

		for _, delta := range []int{-1, 0, 1} {
			nextJump := jump + delta
			if nextJump <= 0 {
				continue
			}

			nextPos := pos + nextJump
			if _, exists := stoneIndex[nextPos]; !exists {
				continue
			}

			if dfs(nextPos, nextJump) {
				memo[key] = true
				return true
			}
		}

		memo[key] = false
		return false
	}

	return dfs(stones[0], 0)
}
