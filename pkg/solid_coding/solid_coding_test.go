package solid_coding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_spiralOrder(t *testing.T) {
    testCases := []struct {
        name     string
        matrix   [][]int
        expected []int
    }{
        {
            "case 1",
            [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
            []int{1, 2, 3, 6, 9, 8, 7, 4, 5},
        },
        {
            "case 2",
            [][]int{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}},
            []int{1, 2, 3, 4, 8, 12, 11, 10, 9, 5, 6, 7},
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := spiralOrder(tc.matrix)
            assert.Equal(t, tc.expected, result)
        })
    }
}

func Test_merge(t *testing.T) {
    testCases := []struct {
        name      string
        intervals [][]int
        expected  [][]int
    }{
        {
            "case 1",
            [][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}},
            [][]int{{1, 6}, {8, 10}, {15, 18}},
        },
        {
            "case 2",
            [][]int{{1, 4}, {4, 5}},
            [][]int{{1, 5}},
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := merge(tc.intervals)
            assert.Equal(t, tc.expected, result)
        })
    }
}

func Test_insert(t *testing.T) {
    testCases := []struct {
        name        string
        intervals   [][]int
        newInterval []int
        expected    [][]int
    }{
        {
            "case 1",
            [][]int{{1, 3}, {6, 9}},
            []int{2, 5},
            [][]int{{1, 5}, {6, 9}},
        },
        {
            "case 2",
            [][]int{{1, 2}, {3, 5}, {6, 7}, {8, 10}, {12, 16}},
            []int{4, 8},
            [][]int{{1, 2}, {3, 10}, {12, 16}},
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := insert(tc.intervals, tc.newInterval)
            assert.Equal(t, tc.expected, result)
        })
    }
}
