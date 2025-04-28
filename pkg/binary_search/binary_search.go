package binarysearch

/**
 * 33. Search in Rotated Sorted Array
 */
func search(nums []int, target int) int {
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
                right = mid - 1
            }
        } else {
            if nums[left] <= target && target < nums[mid] {
                right = mid - 1
            } else {
                left = mid + 1
            }
        }
    }
    return -1
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
