package dp

/**
 * Backtracking (Recursive approach, DFS to try all possibilities)
 *
 *    for choice in choices {
 *        if notAllowed(choice) // Pruning => mark visited & skip if visited
 *            continue
 *
 *        makeChoice(choice)
 *        backtrack(path + choice, updatedChoices)
 *        undoChoice(choice) // backtrack
 *    }
 *
 * | Component              | Purpose                                       |
 * | ---------------------- | --------------------------------------------- |
 * | **Recursive function** | DFS-style traversal of the solution space     |
 * | **Base case**          | When a full solution is found, save it        |
 * | **Choices**            | The candidates to try at each step            |
 * | **Constraints check**  | Skip invalid paths early (pruning)            |
 * | **Backtrack (undo)**   | Remove the last choice before trying the next |
 *
 */

/**
 * 377. Combination Sum IV
 *
 * Given an array of distinct integers nums and a target integer target, return the number of possible combinations that add up to target.
 *
 * dp[amount] += dp[amount - num]
 */
func combinationSum4Memoization(nums []int, target int) int {
	memo := make([]int, target+1)
	for i := range memo {
		memo[i] = -1
	}

	var dfs func(remaining int) int
	dfs = func(remaining int) int {
		if remaining < 0 {
			return 0
		}
		if remaining == 0 {
			return 1
		}
		if memo[remaining] != -1 {
			return memo[remaining]
		}

		count := 0
		for _, num := range nums {
			count += dfs(remaining - num)
		}
		memo[remaining] = count
		return count
	}

	return dfs(target)
}
