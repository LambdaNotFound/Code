package graph

/**
 * Alien Dictionary
 *
 * A string a is lexicographically smaller than a string b if either of the following is true:
 *     a). The first letter where they differ is smaller in a than in b.
 *     b). a is a prefix of b and a.length < b.length.
 *
 * invalid lex order: words=["aaa", "aa", "a"]
 *
 * append res string when popping from queue
 */
func foreignDictionary(words []string) string {
	// Step 1: collect all unique characters
	charSet := make(map[byte]bool)
	for _, word := range words {
		for i := 0; i < len(word); i++ {
			charSet[word[i]] = true
		}
	}

	// Step 2: extract ordering rules by comparing adjacent words
	adjList := make(map[byte][]byte)
	for i := 0; i+1 < len(words); i++ {
		w1, w2 := words[i], words[i+1]
		minLen := min(len(w1), len(w2))
		isPrefixCase := len(w1) > len(w2)

		for j := 0; j < minLen; j++ {
			if w1[j] == w2[j] {
				// invalid: "apple" before "app" is impossible
				if j == minLen-1 && isPrefixCase {
					return ""
				}
				continue
			}
			// first differing char reveals ordering rule
			adjList[w1[j]] = append(adjList[w1[j]], w2[j])
			break
		}
	}

	// Step 3: compute indegree for each character
	indegree := make(map[byte]int)
	for _, neighbors := range adjList {
		for _, neighbor := range neighbors {
			indegree[neighbor]++
		}
	}

	// Step 4: topological sort via Kahn's algorithm
	// seed queue with all zero indegree characters
	queue := make([]byte, 0)
	for char := range charSet {
		if indegree[char] == 0 {
			queue = append(queue, char)
		}
	}

	result := make([]byte, 0, len(charSet))
	for len(queue) > 0 {
		char := queue[0]
		queue = queue[1:]
		result = append(result, char)

		for _, neighbor := range adjList[char] {
			indegree[neighbor]--
			if indegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	// cycle detected if not all characters are in result
	if len(result) != len(charSet) {
		return ""
	}

	return string(result)
}

func foreignDictionaryDFS(words []string) string {
	// Step 1: collect all unique characters
	charSet := make(map[byte]bool)
	for _, word := range words {
		for i := 0; i < len(word); i++ {
			charSet[word[i]] = true
		}
	}

	// Step 2: extract ordering rules by comparing adjacent words
	adjList := make(map[byte][]byte)
	for i := 0; i+1 < len(words); i++ {
		w1, w2 := words[i], words[i+1]
		minLen := min(len(w1), len(w2))
		for j := 0; j < minLen; j++ {
			if w1[j] == w2[j] {
				// invalid: "apple" before "app" is impossible
				if j == minLen-1 && len(w1) > len(w2) {
					return ""
				}
				continue
			}
			adjList[w1[j]] = append(adjList[w1[j]], w2[j])
			break
		}
	}

	// Step 3: DFS topological sort
	// three states:
	// 0 = unvisited
	// 1 = visiting (currently in DFS stack — cycle detection)
	// 2 = visited  (fully processed)
	state := make(map[byte]int)
	result := make([]byte, 0, len(charSet))

	var dfs func(char byte) bool
	dfs = func(char byte) bool {
		if state[char] == 1 {
			// currently in stack = cycle detected
			return false
		}
		if state[char] == 2 {
			// already fully processed — skip
			return true
		}

		// mark as visiting
		state[char] = 1

		for _, neighbor := range adjList[char] {
			if !dfs(neighbor) {
				return false
			}
		}

		// mark as fully visited — post order
		state[char] = 2
		// append AFTER all dependencies explored
		result = append(result, char)
		return true
	}

	for char := range charSet {
		if !dfs(char) {
			return ""
		}
	}

	// reverse post-order to get topological order
	for left, right := 0, len(result)-1; left < right; left, right = left+1, right-1 {
		result[left], result[right] = result[right], result[left]
	}

	return string(result)
}
