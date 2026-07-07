package hashmap

/**
 * 128. Longest Consecutive Sequence
 *
 * Hash Set: map[int]bool
 *     iterate through keys
 */
func longestConsecutive(nums []int) int {
	numSet := make(map[int]bool)
	for _, num := range nums {
		numSet[num] = true
	}

	longestStreak := 0
	for num := range numSet {
		if numSet[num-1] { // only start sequence from the smallest number
			continue
		}

		currentStreak := 1
		for numSet[num+currentStreak] {
			currentStreak++
		}

		longestStreak = max(longestStreak, currentStreak)
	}

	return longestStreak
}

/**
 * 217. Contains Duplicate
 */
func containsDuplicate(nums []int) bool {
	seen := make(map[int]struct{})
	for _, num := range nums {
		if _, found := seen[num]; found {
			return true
		}
		seen[num] = struct{}{}
	}
	return false
}
