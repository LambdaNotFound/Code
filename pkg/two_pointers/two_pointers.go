package two_pointers

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
