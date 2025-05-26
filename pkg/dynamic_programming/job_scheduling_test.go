package dynamic_programming

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_jobScheduling(t *testing.T) {
    testCases := []struct {
        name      string
        startTime []int
        endTime   []int
        profit    []int
        expected  int
    }{
        {
            "case 1",
            []int{1, 2, 3, 3},
            []int{3, 4, 5, 6},
            []int{50, 10, 40, 70},
            120,
        },
        {
            "case 2",
            []int{1, 2, 3, 4, 6},
            []int{3, 5, 10, 6, 9},
            []int{20, 20, 100, 70, 60},
            150,
        },
        {
            "case 3",
            []int{1, 1, 1},
            []int{2, 3, 4},
            []int{5, 6, 4},
            6,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := jobScheduling_dfs(tc.startTime, tc.endTime, tc.profit)
            assert.Equal(t, tc.expected, result)
        })
    }
}
