package binarysearch

/**
 * for rotated sorted array
 *     compare nums[mid] with nums[right] to determine
 *     which half is strictly sorted
 */

/**
 * 33. Search in Rotated Sorted Array
 */
func searchRotatedSortedArray(nums []int, target int) int {
	left, right := 0, len(nums)-1
	for left <= right {
		mid := left + (right-left)/2
		if nums[mid] == target {
			return mid
		}

		if nums[mid] < nums[right] {
			if nums[mid] < target && target <= nums[right] {
				left = mid + 1
			} else {
				right = mid - 1 // nums[mid] == target
			}
		} else {
			if nums[left] <= target && target < nums[mid] {
				right = mid - 1 // nums[mid] == target
			} else {
				left = mid + 1
			}
		}
	}
	return -1
}

/**
 * 153. Find Minimum in Rotated Sorted Array
 */
func findMinInRotatedSortedArray(nums []int) int {
	left, right := 0, len(nums)-1
	for left < right {
		mid := left + (right-left)/2
		if nums[mid] < nums[right] { // min is in [left, mid)
			right = mid
		} else { // min is in (mid, right)
			left = mid + 1
		}
	}
	return nums[left] // return nums[right]
}

/**
 * 154. Find Minimum in Rotated Sorted Array II
 *         (not necessarily with distinct values)
 *
 * shrink right when nums[mid] == nums[right], narrow the search by
 * moving right: nums[mid] > nums[right]
 *
 */
func findMinInRotatedSortedArrayII(nums []int) int {
	left, right := 0, len(nums)-1
	for left < right {
		mid := left + (right-left)/2
		if nums[mid] < nums[right] { // min is in [left, mid)
			right = mid
		} else if nums[mid] == nums[right] {
			right -= 1
		} else { // nums[mid] > nums[right]
			left = mid + 1
		}
	}
	return nums[left] // return nums[right]
}

/**
 * 81. Search in Rotated Sorted Array II
 *     (not necessarily with distinct values)
 * Why right-- works: when they’re equal,
 * we know the min is still somewhere in [left, right],
 * and removing nums[right] doesn’t remove the minimum
 * because nums[mid] has the same value.
 * So the invariant that the min is in [left, right] is preserved.
 */
func searchRotatedSortedArrayII(nums []int, target int) bool {
	left, right := 0, len(nums)-1
	for left <= right {
		mid := left + (right-left)/2
		if nums[mid] == target {
			return true
		}

		if nums[mid] < nums[right] {
			if nums[mid] < target && target <= nums[right] {
				left = mid + 1
			} else {
				right = mid - 1 // nums[mid] == target
			}
		} else if nums[mid] == nums[right] {
			right--
		} else {
			if nums[left] <= target && target < nums[mid] {
				right = mid - 1 // nums[mid] == target
			} else {
				left = mid + 1
			}
		}
	}
	return false
}

/*
 * 34. Find First and Last Position of Element in Sorted Array
 */
func searchRange(nums []int, target int) []int {
	var left, right = 0, len(nums) // [0, len)
	var first, last = -1, -1

	for left < right {
		mid := left + (right-left)/2
		if target > nums[mid] { // lower_bound(), left is the first element <= target
			left = mid + 1
		} else {
			right = mid
		}
	}
	if left < len(nums) && nums[left] == target {
		first = left
	}

	left, right = 0, len(nums)
	for left < right {
		mid := left + (right-left)/2
		if target >= nums[mid] { // upper_bound(), left is the first element > target
			left = mid + 1
		} else {
			right = mid
		}
	}
	if left > 0 && nums[left-1] == target {
		last = left - 1
	}

	return []int{first, last}
}

/**
 * 278. First Bad Version
 */
func isBadVersion(version int) bool { return true }
func firstBadVersion(n int) int {
	left, right := 0, n
	for left < right {
		mid := left + (right-left)/2
		if isBadVersion(mid) == true {
			right = mid
		} else {
			left = mid + 1
		}
	}
	return left
}

/**
 * 704. Binary Search
 */
func binarySearch(nums []int, target int) int {
	left, right := 0, len(nums)
	for left < right {
		mid := left + (right-left)/2
		if nums[mid] == target {
			return mid
		} else if nums[mid] < target {
			left = mid + 1
		} else {
			right = mid
		}
	}
	return -1
}
