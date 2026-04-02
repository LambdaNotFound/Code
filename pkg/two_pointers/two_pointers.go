package two_pointers

import (
	"sort"
	"strings"
)

/**
 * Two Pointers
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
	area := 0
	for left, right := 0, len(height)-1; left < right; {
		h := min(height[left], height[right])
		w := right - left
		area = max(area, h*w)
		if height[left] > height[right] {
			right--
		} else {
			left++
		}
	}
	return area
}

/**
 * 15. 3Sum
 * Given an array nums of n integers, are there elements a, b, c in nums
 * such that a + b + c = 0?
 * Find all unique triplets in the array which gives the sum of zero.
 *
 */
func threeSum(nums []int) [][]int {
	res := make([][]int, 0)
	sort.Ints(nums)

	for i := 0; i < len(nums); i++ {
		if i > 0 && nums[i] == nums[i-1] { // skip previously used num
			continue
		}

		for j, k := i+1, len(nums)-1; j < k; {
			sum := nums[i] + nums[j] + nums[k]
			if sum == 0 {
				res = append(res, []int{nums[i], nums[j], nums[k]})
				j++
				k-- // k + 1 < len(nums) unnecessary, k is decremented BEFORE the duplicate check
				for ; j < k && nums[j] == nums[j-1]; j++ {
				}
				for ; j < k && nums[k] == nums[k+1]; k-- {
				}
			} else if sum > 0 {
				k--
			} else {
				j++
			}
		}
	}
	return res
}

/**
 * 16. 3Sum Closest
 */
func threeSumClosest(nums []int, target int) int {
	sort.Ints(nums)
	closest := nums[0] + nums[1] + nums[2] // initialise with first triplet

	abs := func(x int) int {
		if x < 0 {
			return -x
		}
		return x
	}

	for i := 0; i < len(nums)-2; i++ {
		if i > 0 && nums[i] == nums[i-1] {
			continue
		}

		j, k := i+1, len(nums)-1
		for j < k {
			sum := nums[i] + nums[j] + nums[k]

			// update closest if this sum is nearer to target
			if abs(sum-target) < abs(closest-target) {
				closest = sum
			}

			if sum == target {
				return sum // can't get closer than exact match
			} else if sum < target {
				j++
			} else {
				k--
			}
		}
	}

	return closest
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
 *
 * [0, red)    → all 0s  (sorted)
 * [red, i)    → all 1s  (sorted)
 * [i, blue]   → unknown (unprocessed)
 * (blue, n-1] → all 2s  (sorted)
 *
 */
func sortColors(nums []int) {
	for red, blue, i := 0, len(nums)-1, 0; i <= blue; {
		switch nums[i] {
		case 0: // red
			nums[i], nums[red] = nums[red], nums[i]
			red++
			i++ // advance both of the pointers
		case 2: // blue
			nums[i], nums[blue] = nums[blue], nums[i]
			blue--
		default: // white
			i++
		}
	}
}

/**
 * 977. Squares of a Sorted Array
 *
 */
func sortedSquares(nums []int) []int {
	res := make([]int, len(nums))
	for left, right, i := 0, len(nums)-1, len(nums)-1; left <= right; i-- {
		a, b := nums[left]*nums[left], nums[right]*nums[right]
		if a > b {
			res[i] = a
			left++
		} else {
			res[i] = b
			right--
		}
	}
	return res
}

/**
 * 283. Move Zeroes
 */
func moveZeroes(nums []int) {
	for left, right := 0, 0; right < len(nums); right++ {
		if nums[right] != 0 {
			nums[left], nums[right] = nums[right], nums[left]
			left++
		}
	}
}

/**
 * 283. Remove Element
 */
func removeElement(nums []int, val int) int {
	pos := 0
	for i := 0; i < len(nums); i++ {
		if nums[i] != val {
			nums[pos] = nums[i]
			pos += 1
		}
	}

	return pos
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
			if !strings.EqualFold(left, right) {
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
