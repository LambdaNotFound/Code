package two_pointers

import "sort"

/**
 * 11. Container With Most Water
 *
 * always move the shorter side
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
        if res < area {
            res = area
        }
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

    sort.Ints(nums)
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
 */
func sortColors(nums []int) {
    n := len(nums)
    for red, blue, i := 0, n-1, 0; i <= blue; {
        if nums[i] == 0 {
            nums[i], nums[red] = nums[red], nums[i]
            red += 1
            i += 1 // advance i, so that only 0, 1 on the left of i
        } else if nums[i] == 2 {
            nums[i], nums[blue] = nums[blue], nums[i]
            blue -= 1
        } else {
            i += 1
        }
    }
}

/**
 * 1249. Minimum Remove to Make Valid Parentheses
 */
func minRemoveToMakeValid(s string) string {
    return ""
}
