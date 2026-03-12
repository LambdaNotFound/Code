package apidesign

/**
 * 208. Implement Trie (Prefix Tree)
 */
type Trie struct {
	nodes     map[rune]*Trie
	endOfWord bool // endOfWord string
}

func ConstructorTrie() Trie {
	return Trie{make(map[rune]*Trie), false}
}

func (t *Trie) Insert(word string) {
	for _, c := range word {
		if _, found := t.nodes[c]; !found {
			node := ConstructorTrie()
			t.nodes[c] = &node
		}
		t = t.nodes[c]
	}
	t.endOfWord = true
}

func (t *Trie) Search(word string) bool {
	for _, c := range word {
		if _, found := t.nodes[c]; !found {
			return false
		}
		t = t.nodes[c]
	}
	return t.endOfWord == true
}

func (t *Trie) StartsWith(prefix string) bool {
	for _, c := range prefix {
		if _, found := t.nodes[c]; !found {
			return false
		}
		t = t.nodes[c]
	}
	return true
}
