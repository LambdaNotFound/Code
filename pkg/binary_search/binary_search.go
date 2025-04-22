package binarysearch

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
func search(nums []int, target int) int {
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
