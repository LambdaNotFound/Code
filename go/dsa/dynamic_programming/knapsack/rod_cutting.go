package knapsack

import "slices"

/**
 * Rod Cutting problem <= unbounded knapsack shape, 1-D Tabulation
 *
 */
func maxProfit(length int, salePrice []int) int {
	dp := make([]int, length+1)
	for i := 1; i <= length; i++ {
		for l := 1; l <= i; l++ {
			dp[i] = max(dp[i], dp[i-l]+salePrice[l-1])
		}
	}
	return dp[length]
}

/**
 *
 * You are given a list of rod lengths, and you may cut any rod into smaller pieces to produce rods of a single chosen integer length (saleLength).
 * Each cut costs costPerCut, and each resulting rod of length saleLength earns revenue equal to saleLength × salePrice
 *
 * pieces = L / saleLength          // floor division
 * cuts = pieces if L % saleLength != 0 else pieces - 1
 * revenue = pieces * saleLength * salePrice
 * cost = cuts * costPerCut
 * rodProfit = revenue - cost
 * // discard if rodProfit < 0, so max(0, rodProfit)
 */
func maxTotalProfit(rods []int, salePrice, costPerCut float64) float64 {
	maxSaleLength := slices.Max(rods)
	best := 0.0

	for saleLength := 1; saleLength <= maxSaleLength; saleLength++ {
		total := 0.0
		for _, L := range rods {
			pieces := L / saleLength
			if pieces == 0 {
				continue
			}
			cuts := pieces
			if L%saleLength == 0 {
				cuts = pieces - 1
			}
			profit := float64(pieces*saleLength)*salePrice - float64(cuts)*costPerCut
			if profit > 0 {
				total += profit
			}
		}
		best = max(best, total)
	}
	return best
}
