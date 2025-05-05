package binarysearch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_findMedianSortedArrays(t *testing.T) {
    testCases := []struct {
        name     string
        nums1    []int
        nums2    []int
        expected float64
    }{
        {
            "case 1",
            []int{1, 3},
            []int{2},
            2.00000,
        },
        {
            "case 2",
            []int{1, 2},
            []int{3, 4},
            2.50000,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := findMedianSortedArrays(tc.nums1, tc.nums2)
            assert.Equal(t, tc.expected, result)

            result = findMedianSortedArrays_TwoPointersMerging(tc.nums1, tc.nums2)
            assert.Equal(t, tc.expected, result)
        })
    }
}
