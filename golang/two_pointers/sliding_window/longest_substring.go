package slidingwindow

/**
 * 395. Longest Substring with At Least K Repeating Characters
 *
 * Time: O(26 × n) = O(n)
 * Space: O(1) — fixed size [26]int array
 */
func longestSubstring(s string, k int) int {
	if k == 0 {
		return len(s) // every char trivially appears >= 0 times
	}
	maxLen := 0
	// fix the number of unique characters in the window
	for uniqueTarget := 1; uniqueTarget <= 26; uniqueTarget++ {
		freq := [26]int{}
		unique := 0     // distinct chars in window
		satisfying := 0 // chars with freq >= k
		left := 0

		for right := 0; right < len(s); right++ {
			// expand: add s[right]
			rightIdx := s[right] - 'a'
			if freq[rightIdx] == 0 {
				unique++
			}
			freq[rightIdx]++
			if freq[rightIdx] == k {
				satisfying++
			}

			// shrink: too many unique chars in window
			for unique > uniqueTarget {
				leftIdx := s[left] - 'a'
				if freq[leftIdx] == k {
					satisfying--
				}
				freq[leftIdx]--
				if freq[leftIdx] == 0 {
					unique--
				}
				left++
			}

			// record: window has exactly uniqueTarget chars, all satisfying
			if unique == uniqueTarget && satisfying == uniqueTarget {
				maxLen = max(maxLen, right-left+1)
			}
		}
	}
	return maxLen
}

/**
 * 340. Longest Substring with At Most K Distinct Characters
 *
 * Input: s = "eceba", k = 2
 * Output: 3
 * Explanation: T is "ece" which its length is 3.
 *
 * Input: s = "aa", k = 1
 * Output: 2
 * Explanation: T is "aa" which its length is 2.
 *
 * Time: O(n) — both pointers move forward only
 * Space: O(k) — map holds at most k+1 entries before shrinking
 */
func lengthOfLongestSubstringKDistinct(s string, k int) int {
	freq := make(map[byte]int)
	maxLen := 0
	for left, right := 0, 0; right < len(s); right++ {
		freq[s[right]]++
		for len(freq) > k {
			freq[s[left]]--
			if freq[s[left]] == 0 {
				delete(freq, s[left])
			}
			left++
		}
		maxLen = max(maxLen, right-left+1)
	}
	return maxLen
}
