package prefixtree

/**
 * Unsafe Words
 */
type WordTrie struct {
	nodes      map[string]*WordTrie
	endOfWords bool
}

func ConstructorWordTrie() WordTrie {
	return WordTrie{make(map[string]*WordTrie), false}
}

func (t *WordTrie) Insert(words []string) {
	for _, c := range words {
		if _, found := t.nodes[c]; !found {
			node := ConstructorWordTrie()
			t.nodes[c] = &node
		}
		t = t.nodes[c]
	}
	t.endOfWords = true
}

func (t *WordTrie) Search(words []string) bool {
	for _, c := range words {
		if _, found := t.nodes[c]; !found {
			return false
		}
		t = t.nodes[c]
	}
	return t.endOfWords == true
}
