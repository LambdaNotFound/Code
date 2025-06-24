package dynamic_programming

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := canPartition(tc.nums)
            assert.Equal(t, tc.expected, result)
        })
    }
}
