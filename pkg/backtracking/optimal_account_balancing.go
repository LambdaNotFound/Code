package backtracking

import (
	"math"
	"math/bits"
)

/**
 * 465. Optimal Account Balancing
 *
 * Given a list of transactions between a group of people,
 * return the minimum number of transactions required to settle the debt.
 *
 * Input: [[0,1,10], [2,0,5]]
 * Output: 2
 *
 * Explanation:
 * Person #0 gave person #1 $10.
 * Person #2 gave person #0 $5.
 *
 * Input: [[0,1,10], [1,0,1], [1,2,5], [2,0,5]]
 * Output: 1
 *
 * Explanation:
 * Person #0 gave person #1 $10.
 * Person #1 gave person #0 $1.
 * Person #1 gave person #2 $5.
 * Person #2 gave person #0 $5.
 *
 *  Therefore, person #1 only need to give person #0 $4, and all debt is settled.
 */
func minTransfers(transactions [][]int) int {
	// Step 1: compute net balance per person
	balanceMap := make(map[int]int)
	for _, t := range transactions {
		balanceMap[t[0]] -= t[2] // payer loses money
		balanceMap[t[1]] += t[2] // receiver gains money
	}

	// Step 2: collect non-zero balances only
	balances := make([]int, 0)
	for _, balance := range balanceMap {
		if balance != 0 {
			balances = append(balances, balance)
		}
	}

	// Step 3: backtrack to find minimum transactions
	var dfs func(start int) int
	dfs = func(start int) int {
		// skip already-settled balances
		for start < len(balances) && balances[start] == 0 {
			start++
		}

		// base case: all balances settled
		if start == len(balances) {
			return 0
		}

		minTx := math.MaxInt
		for j := start + 1; j < len(balances); j++ {
			// pruning: skip same sign — can't settle debt with debt
			if balances[start]*balances[j] > 0 {
				continue
			}

			// settle balances[start] against balances[j]
			balances[j] += balances[start] // apply
			minTx = min(minTx, 1+dfs(start+1))
			balances[j] -= balances[start] // undo
		}

		return minTx
	}

	return dfs(0)
}

func minTransfersDP(transactions [][]int) int {
	// Step 1: compute net balances
	balanceMap := make(map[int]int)
	for _, t := range transactions {
		balanceMap[t[0]] -= t[2]
		balanceMap[t[1]] += t[2]
	}

	balances := make([]int, 0)
	for _, b := range balanceMap {
		if b != 0 {
			balances = append(balances, b)
		}
	}

	n := len(balances)
	total := 1 << n // 2^n subsets

	// Step 2: precompute sum for each subset
	subsetSum := make([]int, total)
	for mask := 1; mask < total; mask++ {
		for i := 0; i < n; i++ {
			if mask&(1<<i) != 0 {
				subsetSum[mask] = subsetSum[mask^(1<<i)] + balances[i]
				break
			}
		}
	}

	// Step 3: bitmask DP
	// dp[mask] = min transactions to settle people in mask
	dp := make([]int, total)
	for mask := range dp {
		dp[mask] = math.MaxInt
	}
	dp[0] = 0

	for mask := 1; mask < total; mask++ {
		// enumerate all non-empty subsets of mask
		for sub := mask; sub > 0; sub = (sub - 1) & mask {
			// subset must sum to 0 to settle internally
			if subsetSum[sub] == 0 {
				prev := dp[mask^sub]
				if prev != math.MaxInt {
					// popcount(sub)-1 transactions to settle sub
					dp[mask] = min(dp[mask], prev+bits.OnesCount(uint(sub))-1)
				}
			}
		}
	}

	return dp[total-1]
}
