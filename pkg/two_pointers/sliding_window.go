package two_pointers

import (
	"math"
	"sort"
	"strconv"
)

/**
 * 1. Fixed-size window
 *
 * 2. Variable-size window (expand + shrink)
 *
 * 3. Anagram / frequency counter windows
 *
 */

/**
 * 3. Longest Substring Without Repeating Characters
 *
 * Sliding window w/ hashmap tracking if a char appeared before
 *     Two Pointers: i move if char seen before, otherwise j move
 *
 * Time: O(n)
 * Space: O(1)
 */
func lengthOfLongestSubstring(s string) int {
	res, freqMap := 0, make(map[byte]int)
	for left, right := 0, 0; right < len(s); right++ {
		freqMap[s[right]]++
		for freqMap[s[right]] > 1 {
			freqMap[s[left]]--
			left++
		}
		res = max(res, right-left+1)
	}
	return res
}

/*
 * optimize to store index in a map
 */
func lengthOfLongestSubstringTrackIndex(s string) int {
	res := 0
	hashmap := make(map[byte]int)
	for left, right := 0, 0; right < len(s); right++ {
		_, exist := hashmap[s[right]]
		if exist {
			preIndex := hashmap[s[right]]
			left = max(left, preIndex+1)
		}
		hashmap[s[right]] = right
		res = max(res, right-left+1)
	}
	return res
}

func lengthOfLongestSubstring_rune(s string) int {
	lastSeen := make(map[rune]int) // stores last seen index of each rune
	start := 0                     // start index of current window
	maxLen := 0

	for i, r := range []rune(s) { // iterate over runes, not bytes
		if prevIndex, found := lastSeen[r]; found && prevIndex >= start {
			// Move start to right of previous duplicate
			start = prevIndex + 1
		}
		lastSeen[r] = i

		// Current window length = i - start + 1
		if i-start+1 > maxLen {
			maxLen = i - start + 1
		}
	}
	return maxLen
}

func lengthOfLongestSubstring_alt(s string) int {
	res := 0
	hashmap := make(map[byte]bool)
	for left, right := 0, 0; right < len(s); {
		_, exist := hashmap[s[right]]
		if !exist {
			hashmap[s[right]] = true
			res = max(res, right-left+1)

			right += 1
		} else {
			for ; left < right; left++ {
				delete(hashmap, s[left])
				if s[left] == s[right] {
					break
				}
			}
			left += 1
		}
	}
	return res
}

/**
 * 76. Minimum Window Substring
 *
 * return the minimum window substring of s such that
 * every character in t (including duplicates) is included in the window
 */
func minWindow(s string, t string) string {
	freqMap := make(map[byte]int)
	for i := range t {
		freqMap[t[i]]++
	}

	res, size, cnt := "", math.MaxInt, len(t)
	for left, right := 0, 0; right < len(s); right++ {
		freqMap[s[right]]--
		if freqMap[s[right]] >= 0 {
			cnt--
			for cnt == 0 {
				if right-left+1 < size {
					size = right - left + 1
					res = s[left : left+size]
				}

				if freqMap[s[left]] == 0 {
					cnt++
				}
				freqMap[s[left]]++
				left++
			}
		}
	}

	return res
}

func minWindowClaude(s string, t string) string {
	freqMap := make(map[byte]int)
	for i := range t {
		freqMap[t[i]]++
	}

	res := ""
	minSize := math.MaxInt
	cnt := len(t) // number of chars still needed

	for left, right := 0, 0; right < len(s); right++ {
		// expand window: add s[right]
		freqMap[s[right]]--
		if freqMap[s[right]] >= 0 {
			cnt-- // one more required char satisfied
		}

		// shrink window from left while valid
		for cnt == 0 {
			if right-left+1 < minSize {
				minSize = right - left + 1
				res = s[left : left+minSize]
			}

			// remove s[left] from window
			if freqMap[s[left]] == 0 {
				cnt++ // losing a required char
			}
			freqMap[s[left]]++
			left++
		}
	}

	return res
}

/**
 * 209. Minimum Size Subarray Sum
 *
 * return the minimal length of a subarray whose sum is greater than or equal to target
 */
func minSubArrayLen(target int, nums []int) int {
	res := math.MaxInt32

	sum := 0
	for left, right := 0, 0; right < len(nums); right++ {
		sum += nums[right]

		for sum >= target {
			res = min(res, right-left+1)

			sum -= nums[left]
			left += 1
		}
	}

	if res == math.MaxInt32 {
		return 0
	}
	return res
}

/**
 * 340. Longest Substring with At Most K Distinct Characters
 */
func lengthOfLongestSubstringKDistinct(s string, k int) int {
	if k == 0 || len(s) == 0 {
		return 0
	}

	left, maxLen := 0, 0
	freq := make(map[byte]int)

	for right := 0; right < len(s); right++ {
		// expand window
		freq[s[right]]++

		// shrink window if too many distinct chars
		for len(freq) > k { // len(map[byte]int)
			freq[s[left]]--
			if freq[s[left]] == 0 {
				delete(freq, s[left])
			}
			left++
		}

		// update max length
		maxLen = max(maxLen, right-left+1)
	}

	return maxLen
}

/**
 * 438. Find All Anagrams in a String
 *
 * Given two strings s and p, return an array of all the start
 * indices of p's anagrams in s
 *
 */
func findAnagrams(s string, p string) []int {
	hashmap := make(map[byte]int)
	for i := range p {
		hashmap[p[i]]++
	}

	res := make([]int, 0)
	for left, right := 0, 0; right < len(s); {
		cnt, exist := hashmap[s[right]]
		if !exist {
			for left < right { // move left side
				hashmap[s[left]] += 1
				left += 1
			}
			right += 1
			left = right
		} else {
			if cnt > 0 {
				hashmap[s[right]] -= 1
				if right-left+1 == len(p) {
					res = append(res, left)
				}
				right += 1
			} else {
				hashmap[s[left]] += 1
				left += 1
			}
		}
	}

	return res
}

/**
 * 2933. High-Access Employees
 *
 * Input: access_times = [["a","0549"],["b","0457"],["a","0532"],["a","0621"],["b","0540"]]
 * Output: ["a"]
 * Explanation: "a" has three access times in the one-hour period of [05:32, 06:31] which are 05:32, 05:49, and 06:21.
 * But "b" does not have more than two access times at all.
 * So the answer is ["a"].
 *
 */
func findHighAccessEmployees(accessTimes [][]string) []string {
	// group and parse access times by employee name
	timesByEmployee := make(map[string][]int)
	for _, entry := range accessTimes {
		name := entry[0]
		t, _ := strconv.Atoi(entry[1])
		timesByEmployee[name] = append(timesByEmployee[name], t)
	}

	highAccess := make([]string, 0)
	for name, times := range timesByEmployee {
		if len(times) < 3 {
			continue
		}
		sort.Ints(times)

		// sliding window of size 3 — check if 3 accesses fit within 1 hour
		// HHMM format: difference < 100 means within same hour window
		for i := 0; i+2 < len(times); i++ {
			if times[i+2]-times[i] < 100 {
				highAccess = append(highAccess, name)
				break
			}
		}
	}

	return highAccess
}
