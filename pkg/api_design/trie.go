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

/**
 * 211. Design Add and Search Words Data Structure
 *
 * bool search(word) Returns true if there is any string in the data structure that matches word or false otherwise.
 * word may contain dots '.' where dots can be matched with any letter.
 *
 */
type WordDictionary struct {
	nodes     map[rune]*WordDictionary
	endOfWord bool
}

func ConstructorWordDictionary() *WordDictionary {
	return &WordDictionary{nodes: make(map[rune]*WordDictionary)}
}

func (this *WordDictionary) AddWord(word string) {
	for _, c := range word {
		if _, found := this.nodes[c]; !found {
			this.nodes[c] = ConstructorWordDictionary()
		}
		this = this.nodes[c]
	}
	this.endOfWord = true
}

func (this *WordDictionary) Search(word string) bool {
	for i, c := range word {
		if c == '.' {
			// wildcard — explore ALL children recursively
			for _, child := range this.nodes {
				if child.Search(word[i+1:]) {
					return true
				}
			}
			return false // no child led to a match
		}

		// regular character — follow single path
		if _, found := this.nodes[c]; !found {
			return false
		}
		this = this.nodes[c]
	}

	return this.endOfWord
}
