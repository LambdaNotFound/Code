package two_pointers

/**
 * 3. Longest Substring Without Repeating Characters
 *
 * Sliding window w/ hashmap tracking if a char appeared before
 * [i -> if char seen before, j ->]
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

/**
 * optimize to store index in a map
 */
func lengthOfLongestSubstring_optimized(s string) int {
    if len(s) == 0 {
        return 0
    }

    table := make(map[byte]int)
    res, left := 1, -1
    for i := 0; i < len(s); i++ {
        if _, ok := table[s[i]]; ok {
            prevIndex := table[s[i]]
            if left < prevIndex {
                left = prevIndex
            }
        }

        length := (i - 1) - left + 1
        if length > res {
            res = length
        }

        table[s[i]] = i
    }

    return res
}

/**
 * 438. Find All Anagrams in a String
 */
func findAnagrams(s string, p string) []int {
    var res []int
    charToCnt := make(map[rune]int)
    for _, ch := range p {
        charToCnt[ch] += 1
    }
    i := 0
    j := 0
    for j < len(s) {
        ch := rune(s[j])
        charCnt, ok := charToCnt[ch]
        if !ok {
            for i < j {
                ch := rune(s[i])
                charToCnt[ch] += 1
                i += 1
            }
            i += 1
            j += 1
        } else if charCnt == 0 {
            ch = rune(s[i])
            charToCnt[ch] += 1
            i += 1
        } else {
            charToCnt[ch] -= 1
            if charToCnt[ch] == 0 && (j-i+1) == len(p) {
                res = append(res, i)
            }
            j += 1
        }
    }
    return res
}
