package apidesign

/**
 * 208. Implement Trie (Prefix Tree)
 */
type Trie struct {
	children map[rune]*Trie
	isEnd    bool
}

/** Initialize your data structure here. */
func ConstructorPrefixTree() Trie {
	return Trie{
		children: make(map[rune]*Trie),
	}
}

/** Inserts a word into the trie. */
func (this *Trie) Insert(word string) {
	cur := this
	for _, ch := range word {
		if _, exist := cur.children[ch]; !exist {
			cur.children[ch] = &Trie{
				children: make(map[rune]*Trie),
			}
		}
		cur = cur.children[ch]
	}
	cur.isEnd = true
}

/** Returns if the word is in the trie. */
func (this *Trie) Search(word string) bool {
	cur := this
	for _, ch := range word {
		if _, exist := cur.children[ch]; !exist {
			return false
		}
		cur = cur.children[ch]
	}
	return cur.isEnd
}

/** Returns if there is any word in the trie that starts with the given prefix. */
func (this *Trie) StartsWith(prefix string) bool {
	cur := this
	for _, ch := range prefix {
		if _, exist := cur.children[ch]; !exist {
			return false
		}
		cur = cur.children[ch]
	}
	return true
}
