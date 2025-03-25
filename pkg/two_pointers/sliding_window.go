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

        length := i - 1 - left + 1
        if length > res {
            res = length
        }

        table[s[i]] = i
    }

    return res
}
