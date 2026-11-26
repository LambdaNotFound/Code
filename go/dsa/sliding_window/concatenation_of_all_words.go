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
	starts := []int{}
	wordCount := make(map[string]int)
	for _, word := range words {
		wordCount[word]++
	}

	wordLen, windowLen := len(words[0]), len(words[0])*len(words)

	for offset := 0; offset < wordLen; offset++ {
		needed := len(words)
		windowFreq := make(map[string]int)
		for word, cnt := range wordCount {
			windowFreq[word] = cnt // copy wordCount for this offset's window
		}

		for i := offset; i+wordLen <= len(s); i += wordLen {
			// add incoming word on the right
			rightWord := s[i : i+wordLen]
			if windowFreq[rightWord] > 0 {
				needed--
			}
			windowFreq[rightWord]--

			// remove outgoing word on the left (once window is full)
			dropPos := i - windowLen // one wordLen before the window's actual start
			if dropPos >= offset {
				leftWord := s[dropPos : dropPos+wordLen]
				windowFreq[leftWord]++
				if windowFreq[leftWord] > 0 {
					needed++
				}
			}

			if needed == 0 {
				starts = append(starts, dropPos+wordLen)
			}
		}
	}

	return starts
}
