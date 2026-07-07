package unionfind

/**
 * 1061. Lexicographically Smallest Equivalent String
 *
 * Input: s1 = "parker", s2 = "morris", baseStr = "parser"
 * Output: "makkek"
 *
 */
func smallestEquivalentString(s1 string, s2 string, baseStr string) string {
	// parent[ch] = representative character of ch's equivalence class
	parent := make([]int, 26)
	for ch := range parent {
		parent[ch] = ch
	}

	// find returns the lexicographically smallest char in ch's class
	var find func(ch int) int
	find = func(ch int) int {
		if ch != parent[ch] {
			parent[ch] = find(parent[ch]) // path compression
		}
		return parent[ch]
	}

	// union merges two equivalence classes, smaller char becomes root
	union := func(ch1, ch2 int) {
		root1, root2 := find(ch1), find(ch2)
		if root1 == root2 {
			return
		}
		smaller := min(root1, root2)
		parent[root1] = smaller
		parent[root2] = smaller
	}

	// build equivalence classes from s1, s2 pairs
	for i := 0; i < len(s1); i++ {
		ch1, ch2 := int(s1[i]-'a'), int(s2[i]-'a')
		union(ch1, ch2)
	}

	// replace each char in baseStr with smallest equivalent char
	result := []byte(baseStr)
	for i, ch := range baseStr {
		result[i] = byte(find(int(ch-'a'))) + 'a'
	}
	return string(result)
}
