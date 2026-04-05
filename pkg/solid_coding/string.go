package solid_coding

import (
	"math"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

/**
 *    for i, rune := range string {
 *
 *    for i := 0; i < len(string); i++ {
 *        byte := string[i]
 *
 */

/**
 * 8. String to Integer (atoi)
 *
 * Implement the myAtoi(string s) function, which converts a string to a 32-bit signed integer.
 */
func myAtoi(s string) int {
	s = strings.TrimSpace(s)
	multiplier := 1
	if len(s) == 0 {
		return 0
	} else if s[0] == '-' {
		multiplier = -1
		s = s[1:]
	} else if s[0] == '+' {
		s = s[1:]
	}

	res := 0
	for i, r := range s {
		if unicode.IsDigit(r) {
			num := s[i] - '0'
			res = res*10 + int(num)

			if multiplier*res > math.MaxInt32 {
				return math.MaxInt32
			}
			if multiplier*res < math.MinInt32 {
				return math.MinInt32
			}
		} else {
			break
		}
	}
	return multiplier * res
}

/**
 * 67. Add Binary
 */
func addBinary(a string, b string) string {
	finalstr := ""
	v1, v2, rem := 0, 0, 0

	for l1, l2 := len(a)-1, len(b)-1; l1 >= 0 || l2 >= 0 || rem != 0; {
		if l1 >= 0 {
			v1, _ = strconv.Atoi(string(a[l1]))
		}
		if l2 >= 0 {
			v2, _ = strconv.Atoi(string(b[l2]))
		}

		sum := v1 + v2 + rem

		// according to sum append appropriate character in finalstr
		switch sum {
		case 3:
			finalstr = "1" + finalstr
			rem = 1
		case 2:
			finalstr = "0" + finalstr
			rem = 1
		case 1:
			finalstr = "1" + finalstr
			rem = 0
		case 0:
			finalstr = "0" + finalstr
			rem = 0
		}

		v1, v2 = 0, 0
		l1 -= 1
		l2 -= 1
	}

	return finalstr
}

/**
 * 409. Longest Palindrome
 *
 * return the length of the longest palindrome that can be built with those letters
 */
func longestPalindromeLength(s string) int {
	res := 0
	charMap := make(map[rune]int)
	for _, val := range s {
		charMap[val] += 1
		count := charMap[val]
		if count == 2 {
			res += 2
			delete(charMap, val)
		}
	}

	if len(charMap) != 0 {
		res += 1
	}
	return res
}

/*
 * 49. Group Anagrams
 */
func groupAnagrams(strs []string) [][]string {
	var sortRunes func(string) string
	sortRunes = func(s string) string {
		runes := []rune(s)
		sort.Slice(runes, func(i, j int) bool {
			return runes[i] < runes[j]
		})
		return string(runes)
	}

	anagramsMap := make(map[string][]string)
	for _, str := range strs {
		key := sortRunes(str)
		anagramsMap[key] = append(anagramsMap[key], str)
	}

	res := make([][]string, 0)
	for _, anagrams := range anagramsMap {
		res = append(res, anagrams)
	}
	return res
}

/*
 * 125. Valid Palindrome
 */
func isPalindrome(s string) bool {
	var cleaned strings.Builder
	for _, char := range strings.ToLower(s) {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			cleaned.WriteRune(char)
		}
	}
	cleanedStr := cleaned.String()

	for i, j := 0, len(cleanedStr)-1; i < j; i, j = i+1, j-1 {
		if cleanedStr[i] != cleanedStr[j] {
			return false
		}
	}
	return true
}

/**
 * 680. Valid Palindrome II
 *
 * Given a string s, return true if the s can be palindrome after
 * deleting at most one character from it.
 *
 * substr[l+1, r] or substr[l, r-1]
 */
func validPalindrome(s string) bool {
	l, r := 0, len(s)-1

	for l < r {
		if s[l] != s[r] {
			return isPalindromeHelper(s, l+1, r) || isPalindromeHelper(s, l, r-1)
		}
		l++
		r--
	}
	return true
}

func isPalindromeHelper(s string, l, r int) bool {
	for l < r { // covers both cases: odd/even length
		if s[l] != s[r] {
			return false
		}
		l++
		r--
	}
	return true
}

// delete at most k chars
func validPalindromeAtMostK(s string, k int) bool {
	memo := make(map[[3]int]bool)

	var dfs func(left, right, remaining int) bool
	dfs = func(left, right, remaining int) bool {
		// base case: pointers met or crossed
		if left >= right {
			return true
		}

		key := [3]int{left, right, remaining}
		if val, found := memo[key]; found {
			return val
		}

		var res bool
		if s[left] == s[right] {
			res = dfs(left+1, right-1, remaining)
		} else if remaining == 0 {
			res = false
		} else {
			// try deleting left char OR right char
			res = dfs(left+1, right, remaining-1) ||
				dfs(left, right-1, remaining-1)
		}

		memo[key] = res
		return res
	}

	return dfs(0, len(s)-1, k)
}

/**
 * 409. Longest Palindrome
 *
 *  Given a string s which consists of lowercase or uppercase letters,
 * return the length of the longest palindrome that can be built with those letters.
 */
func longestPalindrome(s string) int {
	freqMap := make(map[byte]int)
	for i := range s {
		freqMap[s[i]]++
	}

	res, hasOdd := 0, false
	for _, v := range freqMap {
		res += v - (v % 2) // take even part: v if even, v-1 if odd
		if v%2 == 1 {
			hasOdd = true // at least one odd freq char exists
		}
	}

	// place one odd char in the center
	if hasOdd {
		res++
	}

	return res
}

/**
 * 14. Longest Common Prefix
 */
func longestCommonPrefix(strs []string) string {
	res := ""
	for i := 0; i < len(strs[0]); i++ {
		char := strs[0][i]
		for j := 1; j < len(strs); j++ {
			if i >= len(strs[j]) || strs[j][i] != char {
				return res // ← return immediately on mismatch
			}
		}
		res += string(char)
	}
	return res
}

/**
 * 179. Largest Number
 *
 * sort strings
 */
func largestNumber(nums []int) string {
	strs := make([]string, len(nums))
	for i, num := range nums {
		strs[i] = strconv.Itoa(num)
	}

	sort.Slice(strs, func(i, j int) bool {
		return strs[i]+strs[j] > strs[j]+strs[i]
	})

	// edge case: all zeros
	if strs[0] == "0" {
		return "0"
	}

	return strings.Join(strs, "")
}

/**
 * 1047. Remove All Adjacent Duplicates In String
 */
func removeDuplicates(s string) string {
	stack := []rune{}
	for _, c := range s {
		if len(stack) > 0 && stack[len(stack)-1] == c {
			stack = stack[:len(stack)-1]
		} else {
			stack = append(stack, c)
		}
	}
	return string(stack)
}

func removeKDuplicates(s string, k int) string {
	type frame struct {
		ch    byte
		count int
	}

	stack := []frame{}

	for i := 0; i < len(s); i++ {
		ch := s[i]
		if len(stack) > 0 && stack[len(stack)-1].ch == ch {
			stack[len(stack)-1].count++
			if stack[len(stack)-1].count == k {
				stack = stack[:len(stack)-1] // pop when count reaches k
			}
		} else {
			stack = append(stack, frame{ch, 1})
		}
	}

	// reconstruct string from stack
	var sb strings.Builder
	for _, f := range stack {
		for i := 0; i < f.count; i++ {
			sb.WriteByte(f.ch)
		}
	}
	return sb.String()
}

/**
 * 13. Roman to Integer
 */
func romanToInt(s string) int {
	toInt := map[byte]int{
		'I': 1,
		'V': 5,
		'X': 10,
		'L': 50,
		'C': 100,
		'D': 500,
		'M': 1000,
	}

	res := 0
	for i := range s {
		if i+1 < len(s) && toInt[s[i]] < toInt[s[i+1]] {
			res -= toInt[s[i]]
		} else {
			res += toInt[s[i]]
		}
	}

	return res
}
