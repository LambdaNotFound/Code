package two_pointers

/**
 * 3. Longest Substring Without Repeating Characters
 *
 * a). Sliding window w/ hashmap tracking if a char appeared before
 *
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
