package solid_coding

import (
	"fmt"
)

/**
 * Ranges template, very clean logic (2 branches)
 */
func linearScan(nums []int) [][]int {
	gaps := [][]int{}
	for i := 0; i < len(nums)-1; i++ {
		// bookmarking

		// check relationship between nums[i] and nums[i+1]
		if nums[i+1] > nums[i]+1 { // gap → missing range
			gaps = append(gaps, []int{nums[i] + 1, nums[i+1] - 1})
		}
		if nums[i+1] == nums[i]+1 { // consecutive → summary range
			// extend current range
		}
	}
	return gaps
}

/**
 * 228. Summary Ranges
 *
 * Input: nums = [0,2,3,4,6,8,9]
 * Output: ["0","2->4","6","8->9"]
 */
func summaryRanges(nums []int) []string {
	ranges := []string{}
	for i := 0; i < len(nums); i++ {
		start := i
		for i+1 < len(nums) && nums[i+1]-nums[i] == 1 {
			i++
		}
		if nums[start] == nums[i] {
			ranges = append(ranges, fmt.Sprintf("%d", nums[start]))
		} else {
			ranges = append(ranges, fmt.Sprintf("%d->%d", nums[start], nums[i]))
		}
	}
	return ranges
}

/**
 * 163. Missing Ranges
 *
 * Input:  nums=[0,1,3,50,75], lower=0, upper=99
 * Output: [[2,2], [4,49], [51,74], [76,99]]
 */
func findMissingRanges(nums []int, lower int, upper int) [][]int {
	nums = append([]int{lower - 1}, nums...)
	nums = append(nums, upper+1)

	missing := [][]int{}
	for i := 0; i < len(nums)-1; i++ {
		if nums[i+1]-nums[i] >= 2 {
			missing = append(missing, []int{nums[i] + 1, nums[i+1] - 1})
		}
	}
	return missing
}
