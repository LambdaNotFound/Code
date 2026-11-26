package apidesign

import "math"

/*
 * 465. Optimal Account Balancing
 *
 * Greedy (sorted/heap) does not guarantee minimal transactions.
 *
 * DFS works because it explores all possible pairings.
 */
func minTransfers(transactions [][]int) int {
	balance := map[int]int{}

	for _, t := range transactions {
		from, to, amount := t[0], t[1], t[2]
		balance[from] -= amount
		balance[to] += amount
	}

	// collect non-zero balances
	debts := []int{}
	for _, v := range balance {
		if v != 0 {
			debts = append(debts, v)
		}
	}

	return dfs(debts, 0)
}

func dfs(debts []int, start int) int {
	// skip settled accounts
	for start < len(debts) && debts[start] == 0 {
		start++
	}

	if start == len(debts) {
		return 0
	}

	res := math.MaxInt32

	for i := start + 1; i < len(debts); i++ {
		// only settle opposite signs
		if debts[start]*debts[i] < 0 {

			// settle start with i
			debts[i] += debts[start]

			res = min(res, 1+dfs(debts, start+1))

			// backtrack
			debts[i] -= debts[start]

			// pruning: avoid duplicate states
			if debts[i]+debts[start] == 0 {
				break
			}
		}
	}

	return res
}
