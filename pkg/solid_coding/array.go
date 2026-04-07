package solid_coding

/**
 * 41. First Missing Positive
 *
 * in O(n) time and uses O(1) auxiliary space.
 *
 * Input: nums = [1,2,0]
 * Output: 3
 *
 * Input: nums = [3,4,-1,1]
 * Output: 2
 *
 * goal: nums[i] == i+1
 * move all nums[i] to nums[nums[i]-1]
 *
 */
func firstMissingPositive(nums []int) int {
	for i := range nums { // nums[i] != nums[nums[i]-1], Input: [1,1] => infinite loop
		for nums[i]-1 >= 0 && nums[i]-1 < len(nums) && nums[i] != i+1 && nums[i] != nums[nums[i]-1] {
			nums[i], nums[nums[i]-1] = nums[nums[i]-1], nums[i]
		}
	}

	for i := range nums {
		if nums[i] != i+1 {
			return i + 1
		}
	}
	return len(nums) + 1
}

/**
 * 80. Remove Duplicates from Sorted Array II
 *
 * each unique element appears at most twice
 * Input: nums = [1,1,1,2,2,3]
 * Output: 5, nums = [1,1,2,2,3,_]
 */
func removeDuplicatesFromSortedArray(nums []int) int {
	index, count := 1, 1
	for i := 1; i < len(nums); i++ {
		if nums[i] == nums[index-1] {
			count++
		} else {
			count = 1
		}

		if count <= 2 {
			nums[index] = nums[i]
			index++
		}
	}
	return index
}
