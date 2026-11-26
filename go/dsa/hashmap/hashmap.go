package hashmap

import "strconv"

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

/*
 * Encode this code into a new string where for each digit d at index i in the original code,
 * “nd” will be appended to the output string, where n is the number of occurrences of d after and including index i.
 *
 * Input String: “13612344”
 * Expected Output: “2123161112132414”
 *
 */
func encode(s string) string {
	charCount := map[rune]int{}
	for _, ch := range s {
		charCount[ch]++
	}

	result := ""
	for _, ch := range s {
		count := strconv.Itoa(charCount[ch])
		result += count + string(ch)
		charCount[ch]--
	}

	return result
}
