package two_pointers

import "math"

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
    if len(s) == 0 {
        return 0
    }

    table := make(map[byte]bool)
    res := 1
    for i, j := 0, 0; i < len(s) && j < len(s); {
        if !table[s[j]] {
            table[s[j]] = true
            length := j - i + 1
            if length > res {
                res = length
            }
            j++
        } else {
            table[s[i]] = false // flip
            i++
        }
    }

    return res
}

/*
 * optimize to store index in a map
 */
func lengthOfLongestSubstring_optimized(s string) int {
    if len(s) == 0 {
        return 0
    }

    table := make(map[byte]int)
    res, left := 1, -1 // left + 1 is the beginning of non-repeating str
    for i := 0; i < len(s); i++ {
        if _, ok := table[s[i]]; ok {
            prevIndex := table[s[i]]
            left = max(left, prevIndex) // (left, i]
        }

        length := i - (left + 1) + 1
        res = max(res, length)

        table[s[i]] = i
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
    if len(s) < len(t) {
        return ""
    }

    freqMap := [128]int{}
    count := len(t)
    start, minStart, minLen := 0, 0, math.MaxInt32

    // Initialize the frequency map with characters from t
    for _, c := range t {
        freqMap[c]++
    }

    // Start the sliding window
    for end := 0; end < len(s); end++ {
        if freqMap[s[end]] > 0 {
            count--
        }
        freqMap[s[end]]--

        // Try to minimize the window
        for count == 0 {
            if end-start+1 < minLen {
                minStart = start
                minLen = end - start + 1
            }

            freqMap[s[start]]++
            if freqMap[s[start]] > 0 {
                count++
            }
            start++
        }
    }

    if minLen == math.MaxInt32 {
        return ""
    }
    return s[minStart : minStart+minLen]
}

/**
 * 438. Find All Anagrams in a String
 *
 * Given two strings s and p, return an array of all the start
 * indices of p's anagrams in s
 *
 */
func findAnagrams(s string, p string) []int {
    var res []int
    charToCnt := make(map[byte]int)
    for i, _ := range p {
        charToCnt[p[i]] += 1 // -1 if match a char, +1 to recover
    }

    i := 0
    j := 0
    for j < len(s) {
        charCnt, ok := charToCnt[s[j]]
        if !ok { // move left, right as char @ j not in anagram
            for i < j {
                charToCnt[s[i]] += 1
                i += 1
            }
            i += 1
            j += 1
        } else if charCnt == 0 { // move left, as no more char can be used @ j
            charToCnt[s[i]] += 1
            i += 1
        } else {
            charToCnt[s[j]] -= 1
            if charToCnt[s[j]] == 0 && (j-i+1) == len(p) {
                res = append(res, i)
            }
            j += 1
        }
    }
    return res
}
