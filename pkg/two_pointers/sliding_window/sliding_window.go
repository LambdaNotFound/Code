package slidingwindow

import (
	"math"
	"sort"
	"strconv"
)

/**
 * Sliding Window
 *
 * 1). Fixed-size window
 *
 * 2). Variable-size window (expand + shrink)
 *
 * 3). Anagram / frequency counter windows
 */

/**
 * 3. Longest Substring Without Repeating Characters
 *
 * Sliding window w/ hashmap tracking if a char appeared before
 *     Two Pointers: i move if char seen before, otherwise j move
 *
 * Time: O(n)
 * Space: O(1) — bounded by alphabet size, not input length
 */
func lengthOfLongestSubstring(s string) int {
	freq := make(map[byte]int)
	maxLen := 0
	for left, right := 0, 0; right < len(s); right++ {
		char := s[right]
		freq[char]++
		for freq[char] > 1 {
			freq[s[left]]--
			left++
		}

		maxLen = max(maxLen, right-left+1)
	}
	return maxLen
}

// lengthOfLongestSubstringRune handles multi-byte Unicode characters correctly
// by operating on runes instead of bytes.
func lengthOfLongestSubstringRune(s string) int {
	freq := make(map[rune]int)
	runes := []rune(s)
	maxLen := 0
	for left, right := 0, 0; right < len(runes); right++ {
		char := runes[right]
		freq[char]++
		for freq[char] > 1 {
			freq[runes[left]]--
			left++
		}
		maxLen = max(maxLen, right-left+1)
	}
	return maxLen
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

	res, minLen, needed := "", math.MaxInt, len(t)
	for left, right := 0, 0; right < len(s); right++ {
		freqMap[s[right]]--
		if freqMap[s[right]] >= 0 {
			needed--
		}
		for needed == 0 {
			// shrink past non-required chars first
			for freqMap[s[left]] < 0 {
				freqMap[s[left]]++
				left++
			}
			// now s[left] is a required char — record window
			if right-left+1 < minLen {
				minLen = right - left + 1
				res = s[left : left+minLen]
			}
			// invalidate window by removing s[left]
			freqMap[s[left]]++
			needed++
			left++
		}
	}
	return res
}

/**
 * 424. Longest Repeating Character Replacement
 *
 * You are given a string s and an integer k. You can choose any character of the string and
 * change it to any other uppercase English character. You can perform this operation at most k times.
 */
func characterReplacement(s string, k int) int {
	freq := [26]int{}
	maxFreq, maxLen := 0, 0
	left := 0
	for right := 0; right < len(s); right++ {
		freq[s[right]-'A']++
		maxFreq = max(maxFreq, freq[s[right]-'A'])

		// window is invalid: shrink by exactly one step
		if right-left+1 > k+maxFreq {
			freq[s[left]-'A']--
			left++
		}

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
	freqMap := make(map[byte]int)
	for i := range p {
		freqMap[p[i]]++
	}

	res := []int{}
	for left, right := 0, 0; right < len(s); right++ {
		freqMap[s[right]]--
		// shrink until window is valid again
		for freqMap[s[right]] < 0 {
			freqMap[s[left]]++
			left++
		}
		if right-left+1 == len(p) {
			res = append(res, left)
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
	minLen, sum := math.MaxInt32, 0
	for left, right := 0, 0; right < len(nums); right++ {
		sum += nums[right]

		for sum >= target {
			minLen = min(minLen, right-left+1)

			sum -= nums[left]
			left += 1
		}
	}

	if minLen == math.MaxInt32 {
		return 0
	}
	return minLen
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
		timeVal, _ := strconv.Atoi(entry[1])
		timesByEmployee[name] = append(timesByEmployee[name], timeVal)
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

/**
 * 219. Contains Duplicate II
 *
 * Complexity: O(n)
 */
func containsNearbyDuplicate(nums []int, k int) bool {
	lastSeenAt := make(map[int]int)
	for i, num := range nums {
		if lastIdx, found := lastSeenAt[num]; found && i-lastIdx <= k {
			return true
		}
		lastSeenAt[num] = i
	}
	return false
}
