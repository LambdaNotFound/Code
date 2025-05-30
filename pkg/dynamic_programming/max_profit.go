package dynamic_programming

import "math"

/**
 * 121. Best Time to Buy and Sell Stock
 */
func maxProfit(prices []int) int {
    minPrice := math.MaxInt32
    maxProfit := 0

    for _, currentPrice := range prices {
        minPrice = min(currentPrice, minPrice)
        maxProfit = max(maxProfit, currentPrice-minPrice)
    }

    return maxProfit
}
