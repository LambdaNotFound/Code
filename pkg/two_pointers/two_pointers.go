package two_pointers

import (
	"sort"
	"strings"
)

/** Two Pointers
 *
 * 1). Opposite-end pointers (meet in the middle)
 *
 * 2). Same-direction (fast/slow pointers) (Linked List)
 *
 * 3). Sliding window (special case of two pointers) -> Sliding Window Problems
 *
 * 4). Partitioning with two pointers
 */

/**
 * 11. Container With Most Water
 *
 *    Always move the shorter side
 */
func maxArea(height []int) int {
    res := 0
    for i, j := 0, len(height)-1; i < j; {
        area := 0
        if height[i] <= height[j] {
            area = height[i] * (j - i)
            i += 1
        } else {
            area = height[j] * (j - i)
            j -= 1
        }
        res = max(res, area)
    }
    return res
}

/**
 * 15. 3Sum
 * Given an array nums of n integers, are there elements a, b, c in nums
 * such that a + b + c = 0?
 * Find all unique triplets in the array which gives the sum of zero.
 *
 */
func threeSum(nums []int) [][]int {
    var res [][]int
    if len(nums) < 3 {
        return res
    }

    sort.Ints(nums) // asc ordering
    for i := 0; i < len(nums)-2; i += 1 {
        if nums[i] > 0 {
            break
        }
        if i > 0 && nums[i] == nums[i-1] {
            continue
        }

        target, j, k := 0-nums[i], i+1, len(nums)-1
        for j < k {
            if nums[j]+nums[k] == target {
                res = append(res, []int{nums[i], nums[j], nums[k]})
                for {
                    j += 1
                    if !(j < k && nums[j] == nums[j-1]) {
                        break
                    }
                }
                for {
                    k -= 1
                    if !(j < k && nums[k] == nums[k+1]) {
                        break
                    }
                }
            } else if nums[j]+nums[k] < target {
                j += 1
            } else {
                k -= 1
            }
        }
    }

    return res
}

/**
 * 167. Two Sum II - Input Array Is Sorted
 * T: O(n)
 */
func twoSum(numbers []int, target int) []int {
    left := 0
    right := len(numbers) - 1
    // sort.Ints(nums)
    for left < right {
        sum := numbers[left] + numbers[right]

        if sum < target {
            left += 1
        } else if sum == target {
            return []int{left + 1, right + 1}
        } else {
            right -= 1
        }
    }

    return nil
}

/**
 * 75. Sort Colors
 *
 * Given an array nums with n objects colored red, white, or blue,
 * sort them in-place so that objects of the same color are adjacent,
 * with the colors in the order red, white, and blue.
 *                              0,   1,         2
 *  [0, red,...  i,... blue] red <= i,
 *      next red pos,  next blue pos
 */
func sortColors(nums []int) {
    n := len(nums)
    for red, blue, i := 0, n-1, 0; i <= blue; {
        if nums[i] == 0 { // red
            nums[i], nums[red] = nums[red], nums[i]
            red += 1
            i += 1 // advance both of the pointers
        } else if nums[i] == 2 { // blue
            nums[i], nums[blue] = nums[blue], nums[i]
            blue -= 1
        } else { // white
            i += 1
        }
    }
}

/**
 * 283. Move Zeroes
 */
func moveZeroes(nums []int) {
    for i, zeroIdx := 0, 0; i < len(nums); i++ {
        if nums[i] != 0 {
            nums[zeroIdx], nums[i] = nums[i], nums[zeroIdx]
            zeroIdx += 1
        }
    }
}

/**
 * 125. Valid Palindrome
 */
func isPalindrome(s string) bool {
    var isValidChar = func(s1 string) bool {
        if (s1 >= "A" && s1 <= "Z") ||
            (s1 >= "a" && s1 <= "z") ||
            (s1 >= "0" && s1 <= "9") {
            return true
        }
        return false
    }

    for idx, jdx := 0, len(s)-1; idx <= jdx; {
        left, right := string(s[idx]), string(s[jdx])
        if isValidChar(left) && isValidChar(right) {
            if strings.ToLower(left) != strings.ToLower(right) {
                return false
            }
            idx++
            jdx--
        } else if !isValidChar(left) {
            idx++
        } else if !isValidChar(right) {
            jdx--
        }
    }
    return true
}
