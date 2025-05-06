package binarysearch

import (
	. "gocode/utils"
	"math"
)

/**
 * 4. Median of Two Sorted Arrays
 *
 * Given two sorted arrays nums1 and nums2 of size m and n respectively, return the median of the two sorted arrays.
 * The overall run time complexity should be O(log (m+n)).
 *
 *
 * Leveraging the properties of sorted arrays, we can apply a binary search on the smaller array,
 * effectively partitioning both arrays. This method ensures we find the median without explicitly merging the arrays,
 * adhering to the desired logarithmic time complexity.
 */
func findMedianSortedArrays(nums1 []int, nums2 []int) float64 {
    if len(nums1) > len(nums2) {
        nums1, nums2 = nums2, nums1
    }

    m, n := len(nums1), len(nums2)
    low, high := 0, m
    for low <= high {
        pivotIdxNums1 := low + (high-low)/2
        pivotIdxNums2 := (m+n+1)/2 - pivotIdxNums1 // x + y = (m + n + 1) / 2

        leftToPivotNums1 := math.MinInt64
        if pivotIdxNums1 > 0 {
            leftToPivotNums1 = nums1[pivotIdxNums1-1]
        }

        pivotNums1 := math.MaxInt64
        if pivotIdxNums1 < m {
            pivotNums1 = nums1[pivotIdxNums1]
        }

        leftToPivotNums2 := math.MinInt64
        if pivotIdxNums2 > 0 {
            leftToPivotNums2 = nums2[pivotIdxNums2-1]
        }

        pivotNums2 := math.MaxInt64
        if pivotIdxNums2 < n {
            pivotNums2 = nums2[pivotIdxNums2]
        }

        if leftToPivotNums1 <= pivotNums2 && leftToPivotNums2 <= pivotNums1 {
            if (m+n)%2 == 0 {
                return (float64(Max(leftToPivotNums1, leftToPivotNums2)) + float64(Min(pivotNums1, pivotNums2))) / 2.0
            }
            return float64(Max(leftToPivotNums1, leftToPivotNums2))
        } else if leftToPivotNums1 > pivotNums2 {
            high = pivotIdxNums1 - 1
        } else {
            low = pivotIdxNums1 + 1
        }
    }

    return 0.0
}

func findMedianSortedArrays_TwoPointersMerging(nums1 []int, nums2 []int) float64 {
    merged := make([]int, 0, len(nums1)+len(nums2))
    i, j := 0, 0

    for i < len(nums1) && j < len(nums2) {
        if nums1[i] < nums2[j] {
            merged = append(merged, nums1[i])
            i++
        } else {
            merged = append(merged, nums2[j])
            j++
        }
    }

    for i < len(nums1) {
        merged = append(merged, nums1[i])
        i++
    }
    for j < len(nums2) {
        merged = append(merged, nums2[j])
        j++
    }

    mid := len(merged) / 2
    if len(merged)%2 == 0 {
        return (float64(merged[mid-1]) + float64(merged[mid])) / 2.0
    } else {
        return float64(merged[mid])
    }
}
