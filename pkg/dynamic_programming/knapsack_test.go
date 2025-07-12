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
        {
            "case 1",
            []int{-2, 1, -3, 4, -1, 2, 1, -5, 4},
            6,
        },
        {
            "case 2",
            []int{1},
            1,
        },
        {
            "case 3",
            []int{5, 4, -1, 7, 8},
            23,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := maxSubArray(tc.nums)
            assert.Equal(t, tc.expected, result)
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
        {
            "case 1",
            []int{1, 2, 5},
            11,
            3,
        },
        {
            "case 2",
            []int{2},
            3,
            -1,
        },
        {
            "case 3",
            []int{1},
            0,
            0,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := coinChange(tc.coins, tc.amount)
            assert.Equal(t, tc.expected, result)
        })
    }
}

func Test_canPartition(t *testing.T) {
    testCases := []struct {
        name     string
        nums     []int
        expected bool
    }{
        {
            "case 1",
            []int{1, 5, 11, 5},
            true,
        },
        {
            "case 2",
            []int{1, 2, 3, 5},
            false,
        },
        {
            "case 3",
            []int{1, 2, 5},
            false,
        },
        /*
           {
               "case 4",
               []int{1, -2, -2, 1},
               true,
           },
        */
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := canPartition(tc.nums)
            assert.Equal(t, tc.expected, result)
        })
    }
}
