package two_pointers

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
