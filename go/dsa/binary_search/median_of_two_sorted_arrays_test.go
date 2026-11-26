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
            name:     "both odd length",
            nums1:    []int{1, 3},
            nums2:    []int{2},
            expected: 2.0,
        },
        {
            name:     "both even length",
            nums1:    []int{1, 2},
            nums2:    []int{3, 4},
            expected: 2.5,
        },
        {
            name:     "nums1 empty",
            nums1:    []int{},
            nums2:    []int{2, 3},
            expected: 2.5,
        },
        {
            name:     "nums2 empty",
            nums1:    []int{1, 2, 3},
            nums2:    []int{},
            expected: 2.0,
        },
        {
            name:     "large difference in lengths",
            nums1:    []int{1, 2, 3, 4, 5},
            nums2:    []int{100},
            expected: 3.5,
        },
        {
            name:     "negative numbers",
            nums1:    []int{-5, -3, -1},
            nums2:    []int{-2, 4, 6},
            expected: -1.5,
        },
        {
            name:     "duplicates",
            nums1:    []int{1, 1},
            nums2:    []int{1, 1},
            expected: 1.0,
        },
        {
            name:     "median in second array",
            nums1:    []int{1, 2},
            nums2:    []int{3, 4, 5},
            expected: 3.0,
        },
        {
            name:     "median in first array",
            nums1:    []int{10, 20, 30},
            nums2:    []int{5, 6, 7},
            expected: 8.5,
        },
        {
            name:     "single element each",
            nums1:    []int{1},
            nums2:    []int{2},
            expected: 1.5,
        },
        {
            name:     "single element and multi array",
            nums1:    []int{3},
            nums2:    []int{1, 2, 4, 5},
            expected: 3.0,
        },
        {
            name:     "both arrays empty (edge case)",
            nums1:    []int{},
            nums2:    []int{},
            expected: 0, // your implementation may choose panic or return 0
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
