package dynamic_programming

import "math"

/**
 * 121. Best Time to Buy and Sell Stock
 */
func maxProfit(prices []int) int {
	lowest, profit := math.MaxInt, 0
	for _, price := range prices {
		profit = max(profit, price-lowest)
		lowest = min(lowest, price)
	}
	return profit
}
