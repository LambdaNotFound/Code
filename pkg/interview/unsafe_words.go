package interview

import "strings"

/**
 * Unsafe Words
 *
 * single string/word => phrase multiple-word
 *
 * banned: ["hate speech", "very bad indeed"]
 * text:   "this hate speech is very bad indeed today"
 */
func filterUnsafeWords(text string, bannedWords []string) string {
	trie := ConstructorTrie()
	for _, word := range bannedWords {
		trie.Insert(strings.ToLower(word))
	}

	words := strings.Fields(text) // split on whitespace
	for i, word := range words {
		// strip punctuation for clean lookup
		clean := strings.ToLower(strings.Trim(word, ".,!?;:\"'"))
		if trie.Search(clean) {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}

func filterUnsafePhrases(text string, bannedPhrases []string) string {
	trie := NewPhraseTrie()
	for _, phrase := range bannedPhrases {
		trie.Insert(phrase)
	}

	// tokenize: strip punctuation per word, preserve original for output
	rawWords := strings.Fields(text)
	cleanWords := make([]string, len(rawWords))
	for i, w := range rawWords {
		cleanWords[i] = strings.ToLower(strings.Trim(w, ".,!?;:\"'"))
	}

	masked := make([]bool, len(cleanWords))

	// for each start position, try to match a phrase in trie
	for i := 0; i < len(cleanWords); i++ {
		node := trie.root
		for j := i; j < len(cleanWords); j++ {
			next, found := node.children[cleanWords[j]]
			if !found {
				break // ContainsPrefix would return false here
			}
			node = next
			if node.isEnd {
				for k := i; k <= j; k++ {
					masked[k] = true
				}
				// continue for longer matches
			}
		}
	}

	for i, mask := range masked {
		if mask {
			rawWords[i] = "****"
		}
	}
	return strings.Join(rawWords, " ")
}

type PhraseTrieNode struct {
	children map[string]*PhraseTrieNode
	isEnd    bool
}

type PhraseTrie struct {
	root *PhraseTrieNode
}

func NewPhraseTrie() *PhraseTrie {
	return &PhraseTrie{root: &PhraseTrieNode{
		children: make(map[string]*PhraseTrieNode),
	}}
}

// insert phrase as sequence of word tokens
func (t *PhraseTrie) Insert(phrase string) {
	node := t.root
	for _, word := range strings.Fields(phrase) {
		word = strings.ToLower(word)
		if node.children[word] == nil {
			node.children[word] = &PhraseTrieNode{
				children: make(map[string]*PhraseTrieNode),
			}
		}
		node = node.children[word]
	}
	node.isEnd = true
}

// Contains checks if exact phrase exists in trie
func (t *PhraseTrie) Contains(phrase string) bool {
	node := t.root
	for _, word := range strings.Fields(phrase) {
		word = strings.ToLower(word)
		next, found := node.children[word]
		if !found {
			return false
		}
		node = next
	}
	return node.isEnd
}

// ContainsPrefix checks if any banned phrase STARTS with this prefix
// useful for early termination during scanning
func (t *PhraseTrie) ContainsPrefix(words []string) bool {
	node := t.root
	for _, word := range words {
		word = strings.ToLower(word)
		next, found := node.children[word]
		if !found {
			return false
		}
		node = next
	}
	return true // valid prefix exists (may or may not be isEnd)
}
