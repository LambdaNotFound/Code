package two_pointers

func twoSum(numbers []int, target int) []int {
	left := 0
	right := len(numbers) - 1

	for left < right {
		sum := numbers[left] + numbers[right]

		if sum < target {
			a
			left += 1
		} else if sum == target {
			return []int{left + 1, right + 1}
		} else {
			right -= 1
		}
	}

	return nil
}
