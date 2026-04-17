package knapsack

/**
 * 39. Combination Sum
 *
 * Given an array of distinct integers candidates and a target integer target,
 * return a list of all unique combinations of candidates where the chosen numbers sum to target.
 *
 * Time O(candidates × target × avg combinations)
 * Space O(target × total combinations stored)
 *
 *
 * candidate=2:
 *  amount=2: dp[2-2]=dp[0]={{}}      → dp[2] = [[2]]
 *  amount=3: dp[3-2]=dp[1]={}        → dp[3] = []       (nothing sums to 1)
 *  amount=4: dp[4-2]=dp[2]=[[2]]     → dp[4] = [[2,2]]  ← reuses 2!
 *
 * candidate=3:
 *  amount=3: dp[3-3]=dp[0]={{}}      → dp[3] = [[3]]
 *  amount=4: dp[4-3]=dp[1]={}        → dp[4] unchanged  = [[2,2]]
 *
 */
func combinationSum(candidates []int, target int) [][]int {
	// dp[i] = all combinations that sum to i
	dp := make([][][]int, target+1)
	dp[0] = [][]int{{}} // base case: one way to sum to 0 — empty combination

	for _, candidate := range candidates {
		for amount := candidate; amount <= target; amount++ {
			// unbounded: start from candidate, not from target downward
			for _, combo := range dp[amount-candidate] {
				// deep copy existing combo, append current candidate
				newCombo := append([]int{}, combo...)
				newCombo = append(newCombo, candidate)
				dp[amount] = append(dp[amount], newCombo)
			}
		}
	}

	return dp[target]
}

/**
 * 377. Combination Sum IV
 *
 * dp[amount] += dp[amount - num]
 */
func combinationSum4(nums []int, target int) int {
	dp := make([]int, target+1)
	dp[0] = 1 // base case: one way to sum to 0 — pick nothing

	// outer loop: amounts — counts permutations (order matters)
	for amount := 1; amount <= target; amount++ {
		for _, num := range nums {
			if num <= amount {
				dp[amount] += dp[amount-num]
			}
		}
	}

	return dp[target]
}
