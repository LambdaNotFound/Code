package dynamic_programming

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_maxSubArray(t *testing.T) {
    testCases := []struct {
        name     string
        nums     []int
        expected int
    }{
        {"case 1", []int{-2, 1, -3, 4, -1, 2, 1, -5, 4}, 6},
        {"case 2", []int{1}, 1},
        {"case 3", []int{5, 4, -1, 7, 8}, 23},
        {name: "all_negative", nums: []int{-3, -1, -2}, expected: -1},
        {name: "single_negative", nums: []int{-5}, expected: -5},
        {name: "two_elements", nums: []int{-1, 2}, expected: 2},
        {name: "all_same_positive", nums: []int{3, 3, 3}, expected: 9},
        {name: "alternating_signs", nums: []int{2, -1, 2, -1, 3}, expected: 5},
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            assert.Equal(t, tc.expected, maxSubArray(tc.nums))
            assert.Equal(t, tc.expected, maxSubArrayAlt(tc.nums))
        })
    }
}

func Test_maxProduct(t *testing.T) {
    tests := []struct {
        name     string
        nums     []int
        expected int
    }{
        {name: "example1", nums: []int{2, 3, -2, 4}, expected: 6},
        {name: "example2", nums: []int{-2, 0, -1}, expected: 0},
        {name: "single", nums: []int{5}, expected: 5},
        {name: "single_negative", nums: []int{-3}, expected: -3},
        {name: "all_positive", nums: []int{1, 2, 3, 4}, expected: 24},
        {name: "two_negatives", nums: []int{-2, -3}, expected: 6},
        {name: "zero_in_middle", nums: []int{3, -1, 4, 0, 2}, expected: 4},
        {name: "all_negative_even", nums: []int{-1, -2, -3, -4}, expected: 24},
        {name: "all_negative_odd", nums: []int{-1, -2, -3}, expected: 6},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            assert.Equal(t, tt.expected, maxProduct(tt.nums))
            assert.Equal(t, tt.expected, maxProduct2(tt.nums))
        })
    }
}

func Test_coinChange(t *testing.T) {
    testCases := []struct {
        name     string
        coins    []int
        amount   int
        expected int
    }{
        {"case 1", []int{1, 2, 5}, 11, 3},
        {"case 2", []int{2}, 3, -1},
        {"case 3", []int{1}, 0, 0},
        {name: "exact_coin", coins: []int{5}, amount: 5, expected: 1},
        {name: "multiple_ways_pick_min", coins: []int{1, 3, 4}, amount: 6, expected: 2},
        {name: "large_coin_only", coins: []int{10}, amount: 3, expected: -1},
        {name: "single_coin_many", coins: []int{2}, amount: 10, expected: 5},
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            assert.Equal(t, tc.expected, coinChange(tc.coins, tc.amount))
            assert.Equal(t, tc.expected, coinChange1DDP(tc.coins, tc.amount))
            assert.Equal(t, tc.expected, coinChange2DDP(tc.coins, tc.amount))
            assert.Equal(t, tc.expected, coinChangeRecursionMemoization(tc.coins, tc.amount))
        })
    }
}

func Test_change(t *testing.T) {
    tests := []struct {
        name     string
        amount   int
        coins    []int
        expected int
    }{
        {name: "example1", amount: 5, coins: []int{1, 2, 5}, expected: 4},
        {name: "example2", amount: 3, coins: []int{2}, expected: 0},
        {name: "example3", amount: 10, coins: []int{10}, expected: 1},
        {name: "zero_amount", amount: 0, coins: []int{1, 2, 3}, expected: 1},
        {name: "single_coin_exact", amount: 4, coins: []int{2}, expected: 1},
        {name: "many_ways", amount: 4, coins: []int{1, 2, 3}, expected: 4},
        {name: "no_way", amount: 5, coins: []int{2, 4}, expected: 0},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            assert.Equal(t, tt.expected, change(tt.amount, tt.coins))
            assert.Equal(t, tt.expected, change2DDP(tt.amount, tt.coins))
            assert.Equal(t, tt.expected, changeRecursionMemoization(tt.amount, tt.coins))
        })
    }
}

func Test_canPartition(t *testing.T) {
    testCases := []struct {
        name     string
        nums     []int
        expected bool
    }{
        {"case 1", []int{1, 5, 11, 5}, true},
        {"case 2", []int{1, 2, 3, 5}, false},
        {"case 3", []int{1, 2, 5}, false},
        {name: "two_equal", nums: []int{4, 4}, expected: true},
        {name: "two_unequal", nums: []int{1, 3}, expected: false},
        {name: "odd_sum", nums: []int{1, 2, 3, 4, 5}, expected: false},
        {name: "all_same_even_count", nums: []int{2, 2, 2, 2}, expected: true},
        {name: "large_partition", nums: []int{3, 3, 3, 4, 5}, expected: true},
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            assert.Equal(t, tc.expected, canPartition(tc.nums))
            assert.Equal(t, tc.expected, canPartitionMemoization(tc.nums))
        })
    }
}

func Test_findTargetSumWays(t *testing.T) {
    testCases := []struct {
        name     string
        nums     []int
        target   int
        expected int
    }{
        {"case 1", []int{1, 1, 1, 1, 1}, 3, 5},
        {"case 2", []int{1}, 1, 1},
        {name: "single_negative_target", nums: []int{1}, target: -1, expected: 1},
        {name: "no_way", nums: []int{1}, target: 2, expected: 0},
        {name: "target_zero", nums: []int{1, 1}, target: 0, expected: 2},
        {name: "all_plus", nums: []int{2, 3}, target: 5, expected: 1},
        {name: "all_minus", nums: []int{2, 3}, target: -5, expected: 1},
        {name: "multiple_ways", nums: []int{1, 2, 1}, target: 2, expected: 2},
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            assert.Equal(t, tc.expected, findTargetSumWays(tc.nums, tc.target))
            assert.Equal(t, tc.expected, findTargetSumWaysMemoization(tc.nums, tc.target))
        })
    }
}
