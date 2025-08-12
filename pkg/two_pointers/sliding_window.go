package two_pointers

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
