package array

import "slices"

/**
 * 189. Rotate Array
 *
 * [1,2,3,4,5,6,7], k = 3 => [5,6,7,1,2,3,4]
 *
 * reverse subarray + full reverse
 * [0:n-k] [n-k,n]
 * 4,3,2,1, 7,6,5
 */
func rotateArray(nums []int, k int) {
	n := len(nums)
	k = k % n
	slices.Reverse(nums[:n-k])
	slices.Reverse(nums[n-k : n])
	slices.Reverse(nums)
}
