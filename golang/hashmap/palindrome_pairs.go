package hashmap

import "slices"

/**
 * 336. Palindrome Pairs
 *
 * Case1: If s1 is a blank string,
 *    then for any string that is palindrome s2, s1+s2 and s2+s1 are palindrome. <= "" + a, a + ""
 *
 * Case 2: If s2 is the reversing string of s1,
 *    then s1+s2 and s2+s1 are palindrome. <= abc + cba, cba + abc
 *
 * Case 3: If s1[0:cut] is palindrome and there exists s2 is the reversing string of s1[cut+1:],
 *    then s2+s1 is palindrome. <= fed + abadef*
 *    [0: cut] [cut+1: ]
 *     palin    reverse
 *
 * Case 4: If s1[cut+1:] is palindrome and there exists s2 is the reversing string of s1[0:cut],
 *    then s1+s2 is palindrome. <= abcded + cba
 *    [0: cut] [cut+1: ]
 *     reverse  palin
 */
func palindromePairs(words []string) [][]int {
	idxMap := make(map[string]int)
	for i, word := range words {
		idxMap[reverse(word)] = i
	}

	res := [][]int{}
	for i, word := range words {
		if isPalindrome(word) {
			if idx, found := idxMap[""]; found && i != idx {
				res = append(res, []int{i, idx})
				res = append(res, []int{idx, i})
			}
		}

		if idx, found := idxMap[word]; found && idx != i {
			res = append(res, []int{i, idx})
		}

		for j := 1; j < len(word); j++ {
			prefix, suffix := word[:j], word[j:]
			if isPalindrome(prefix) {
				if idx, found := idxMap[suffix]; found && i != idx {
					res = append(res, []int{idx, i})
				}
			}

			if isPalindrome(suffix) {
				if idx, found := idxMap[prefix]; found && i != idx {
					res = append(res, []int{i, idx})
				}
			}
		}
	}

	return res
}

func reverse(s string) string {
	bytes := []byte(s)
	slices.Reverse(bytes)
	return string(bytes)
}

func isPalindrome(s string) bool {
	for l, r := 0, len(s)-1; l < r; l, r = l+1, r-1 {
		if s[l] != s[r] {
			return false
		}
	}
	return true
}
