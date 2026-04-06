package solid_coding

import (
	"fmt"
	"strconv"
)

/**
 * Ranges template
 */
func linearScan(nums []int) [][]int {
	res := [][]int{}

	for i := 0; i < len(nums)-1; i++ {
		// check relationship between nums[i] and nums[i+1]
		if nums[i+1]-nums[i] >= 2 { // gap → missing range
			res = append(res, []int{nums[i] + 1, nums[i+1] - 1})
		}
		if nums[i+1] == nums[i]+1 { // consecutive → summary range
			// extend current range
		}
	}
	return res
}

/**
 * 228. Summary Ranges
 *
 * Input: nums = [0,2,3,4,6,8,9]
 * Output: ["0","2->4","6","8->9"]
 */
func summaryRanges(nums []int) []string {
	res := []string{}
	for i := 0; i < len(nums); {
		j := i + 1
		for j < len(nums) && nums[j] == nums[j-1]+1 { // reads more naturally as "consecutive"
			j++
		}
		if j == i+1 { // iff the inner loop never advanced
			res = append(res, strconv.Itoa(nums[i]))
		} else {
			res = append(res, fmt.Sprintf("%d->%d", nums[i], nums[j-1]))
		}
		i = j
	}
	return res
}

/**
 * 163. Missing Ranges
 *
 * Input:  nums=[0,1,3,50,75], lower=0, upper=99
 * Output: [[2,2], [4,49], [51,74], [76,99]]
 */
func findMissingRanges(nums []int, lower int, upper int) [][]int {
	res := [][]int{}

	nums = append([]int{lower - 1}, nums...) // ← lower-1
	nums = append(nums, upper+1)             // ← upper+1

	for i := 0; i < len(nums)-1; i++ {
		if nums[i+1]-nums[i] >= 2 { // gap exists
			res = append(res, []int{nums[i] + 1, nums[i+1] - 1})
		}
	}

	return res
}
