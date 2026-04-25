package apidesign

/**
 * 745. Prefix and Suffix Search
 *
 * Design a special dictionary that searches the words in it by a prefix and a suffix.
 *
 * Implement the WordFilter class:
 *
 * WordFilter(string[] words) Initializes the object with the words in the dictionary.
 * f(string pref, string suff) Returns the index of the word in the dictionary,
 * which has the prefix pref and the suffix suff. If there is more than one valid index,
 * return the largest of them. If there is no such word in the dictionary, return -1.
 */

type TrieTreeNode struct {
	children [27]*TrieTreeNode
	maxIndex int
}

func newTrieNode() *TrieTreeNode {
	return &TrieTreeNode{maxIndex: -1}
}

func charToIndex(c byte) int {
	if c == '#' {
		return 26
	}
	return int(c - 'a')
}

type WordFilter struct {
	root *TrieTreeNode
}

func ConstructorSearch(words []string) WordFilter {
	wf := WordFilter{root: newTrieNode()}
	for idx, word := range words {
		for prefLen := 1; prefLen <= len(word); prefLen++ {
			for sufLen := 1; sufLen <= len(word); sufLen++ {
				key := word[:prefLen] + "#" + word[len(word)-sufLen:]
				wf.insert(key, idx)
			}
		}
	}
	return wf
}

func (wf *WordFilter) insert(key string, idx int) {
	cur := wf.root
	for i := 0; i < len(key); i++ {
		c := charToIndex(key[i])
		if cur.children[c] == nil {
			cur.children[c] = newTrieNode()
		}
		cur = cur.children[c]
	}
	// only store at terminal node — intermediate nodes must not be poisoned
	if idx > cur.maxIndex {
		cur.maxIndex = idx
	}
}

func (wf *WordFilter) F(pref string, suff string) int {
	key := pref + "#" + suff
	cur := wf.root
	for i := 0; i < len(key); i++ {
		c := charToIndex(key[i])
		if cur.children[c] == nil {
			return -1
		}
		cur = cur.children[c]
	}
	return cur.maxIndex
}
