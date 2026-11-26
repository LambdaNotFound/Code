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
        {
            name:      "single_job",
            startTime: []int{5},
            endTime:   []int{10},
            profit:    []int{42},
            expected:  42,
        },
        {
            name:      "non_overlapping_take_all",
            startTime: []int{1, 3, 5},
            endTime:   []int{3, 5, 7},
            profit:    []int{10, 20, 30},
            expected:  60,
        },
        {
            name:      "all_overlapping_take_best",
            startTime: []int{1, 1, 1},
            endTime:   []int{10, 10, 10},
            profit:    []int{5, 15, 8},
            expected:  15,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := jobScheduling(tc.startTime, tc.endTime, tc.profit)
            assert.Equal(t, tc.expected, result)

            result = jobSchedulingMemoization(tc.startTime, tc.endTime, tc.profit)
            assert.Equal(t, tc.expected, result)
        })
    }
}
