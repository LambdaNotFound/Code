package slidingwindow

/**
 * 30. Substring with Concatenation of All Words
 *
 * Input: s = "barfoothefoobarman", words = ["foo","bar"]
 *
 * Output: [0,9]
 *
 * window always = wordCount words wide
 * right adds:  s[i : i+wordLen]
 * left drops:  s[i-windowLen : i-windowLen+wordLen]  (only once window is full)
 *
 */
func findSubstring(s string, words []string) []int {
	res := []int{}
	wordCount := make(map[string]int)
	for _, w := range words {
		wordCount[w]++
	}

	wordLen, windowLen := len(words[0]), len(words[0])*len(words)

	for offset := 0; offset < wordLen; offset++ {
		count := len(words)
		window := make(map[string]int)
		for k, v := range wordCount {
			window[k] = v // copy wordCount for this offset's window
		}

		for i := offset; i+wordLen <= len(s); i += wordLen {
			// add incoming word on the right
			current := s[i : i+wordLen]
			if window[current] > 0 {
				count--
			}
			window[current]--

			// remove outgoing word on the left (once window is full)
			start := i - windowLen // left edge of the window
			if start >= offset {
				previous := s[start : start+wordLen]
				window[previous]++
				if window[previous] > 0 {
					count++
				}
			}

			if count == 0 {
				res = append(res, start+wordLen)
			}
		}
	}

	return res
}
