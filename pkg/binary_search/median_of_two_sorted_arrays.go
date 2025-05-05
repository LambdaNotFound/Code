package binarysearch

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
