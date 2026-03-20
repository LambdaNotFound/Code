package hashmap

/**
 * 1. Two Sum
 *
 * one-pass, to avoid using same number twice
 */
func twoSum(nums []int, target int) []int {
	complements := make(map[int]int)
	for i, num := range nums {
		if index, found := complements[target-num]; found {
			return []int{index, i}
		} else {
			complements[num] = i
		}
	}
	return nil
}

/**
 * 383. Ransom Note
 */
func canConstruct(ransomNote string, magazine string) bool {
	hashmap := make(map[rune]int)
	for _, val := range magazine {
		hashmap[val] += 1
	}

	for _, val := range ransomNote {
		if hashmap[val] > 0 {
			hashmap[val] -= 1
		} else {
			return false
		}
	}
	return true
}

/**
 * 242. Valid Anagram
 */
func isAnagram(s string, t string) bool {
	chars := make(map[rune]int)
	for _, v := range s {
		chars[v]++
	}

	for _, v := range t {
		if number, exists := chars[v]; !exists || number == 0 {
			return false
		}
		chars[v]--
	}

	return len(s) == len(t)
}

/**
 * 525. Contiguous Array
 */
func findMaxLength(nums []int) int {
	// balance: +1 for 0, -1 for 1
	// equal 0s and 1s → balance returns to same value
	// store first occurrence of each balance value
	balanceToIndex := map[int]int{0: -1}

	res, balance := 0, 0
	for i, num := range nums {
		if num == 0 {
			balance++
		} else {
			balance--
		}

		if idx, found := balanceToIndex[balance]; found {
			res = max(res, i-idx)
		} else {
			balanceToIndex[balance] = i
		}
	}

	return res
}
