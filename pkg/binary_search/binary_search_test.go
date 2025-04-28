package binarysearch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_bsearchRotatedSortedArray(t *testing.T) {
    testCases := []struct {
        name     string
        nums     []int
        target   int
        expected int
    }{
        {
            "case 1",
            []int{4,5,6,7,0,1,2},
            0,
            4,
        },
        {
            "case 2",
            []int{4,5,6,7,0,1,2},
            3,
            -1,
        },
        {
            "case 3",
            []int{1},
            0,
            -1,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := searchRotatedSortedArray(tc.nums, tc.target)
            assert.Equal(t, tc.expected, result)
        })
    }
}

func Test_binarySearch(t *testing.T) {
    testCases := []struct {
        name     string
        nums     []int
        target   int
        expected int
    }{
        {
            "case 1",
            []int{-1, 0, 3, 5, 9, 12},
            9,
            4,
        },
        {
            "case 2",
            []int{-1, 0, 3, 5, 9, 12},
            2,
            -1,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := binarySearch(tc.nums, tc.target)
            assert.Equal(t, tc.expected, result)
        })
    }
}
