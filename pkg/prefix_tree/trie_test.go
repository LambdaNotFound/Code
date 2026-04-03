package prefixtree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Trie(t *testing.T) {
	t.Run("leetcode_example", func(t *testing.T) {
		trie := ConstructorTrie()
		trie.Insert("apple")
		assert.True(t, trie.Search("apple"))
		assert.False(t, trie.Search("app"))
		assert.True(t, trie.StartsWith("app"))
		trie.Insert("app")
		assert.True(t, trie.Search("app"))
	})

	t.Run("empty_trie_search", func(t *testing.T) {
		trie := ConstructorTrie()
		assert.False(t, trie.Search("anything"))
		assert.False(t, trie.StartsWith("a"))
	})

	t.Run("single_char_word", func(t *testing.T) {
		trie := ConstructorTrie()
		trie.Insert("a")
		assert.True(t, trie.Search("a"))
		assert.False(t, trie.Search("ab"))
		assert.True(t, trie.StartsWith("a"))
		assert.False(t, trie.StartsWith("b"))
	})

	t.Run("prefix_not_word", func(t *testing.T) {
		trie := ConstructorTrie()
		trie.Insert("hello")
		assert.False(t, trie.Search("hell"))
		assert.True(t, trie.StartsWith("hell"))
		assert.True(t, trie.StartsWith("hello"))
		assert.False(t, trie.StartsWith("hellop"))
	})

	t.Run("multiple_words_common_prefix", func(t *testing.T) {
		trie := ConstructorTrie()
		trie.Insert("car")
		trie.Insert("card")
		trie.Insert("care")
		trie.Insert("careful")
		assert.True(t, trie.Search("car"))
		assert.True(t, trie.Search("card"))
		assert.True(t, trie.Search("care"))
		assert.True(t, trie.Search("careful"))
		assert.False(t, trie.Search("ca"))
		assert.False(t, trie.Search("cares"))
		assert.True(t, trie.StartsWith("car"))
		assert.True(t, trie.StartsWith("care"))
		assert.False(t, trie.StartsWith("cards"))
	})

	t.Run("no_false_prefix_match", func(t *testing.T) {
		trie := ConstructorTrie()
		trie.Insert("abc")
		assert.False(t, trie.Search("ab"))
		assert.False(t, trie.Search("a"))
		assert.False(t, trie.Search("abcd"))
		assert.True(t, trie.StartsWith("a"))
		assert.True(t, trie.StartsWith("ab"))
		assert.True(t, trie.StartsWith("abc"))
		assert.False(t, trie.StartsWith("abcd"))
	})

	t.Run("duplicate_insert", func(t *testing.T) {
		trie := ConstructorTrie()
		trie.Insert("word")
		trie.Insert("word")
		assert.True(t, trie.Search("word"))
	})
}

func Test_WordDictionary(t *testing.T) {
	t.Run("leetcode_example", func(t *testing.T) {
		wd := ConstructorWordDictionary()
		wd.AddWord("bad")
		wd.AddWord("dad")
		wd.AddWord("mad")
		assert.False(t, wd.Search("pad"))
		assert.True(t, wd.Search("bad"))
		assert.True(t, wd.Search(".ad"))
		assert.True(t, wd.Search("b.."))
	})

	t.Run("exact_match", func(t *testing.T) {
		wd := ConstructorWordDictionary()
		wd.AddWord("word")
		assert.True(t, wd.Search("word"))
		assert.False(t, wd.Search("ward"))
		assert.False(t, wd.Search("wor"))
		assert.False(t, wd.Search("words"))
	})

	t.Run("single_wildcard", func(t *testing.T) {
		wd := ConstructorWordDictionary()
		wd.AddWord("bad")
		wd.AddWord("dad")
		wd.AddWord("mad")
		assert.True(t, wd.Search(".ad"))  // matches bad, dad, mad
		assert.True(t, wd.Search("b.."))  // matches bad
		assert.False(t, wd.Search("p..")) // no word starting with p
		assert.False(t, wd.Search("ba"))  // too short
	})

	t.Run("all_wildcards", func(t *testing.T) {
		wd := ConstructorWordDictionary()
		wd.AddWord("cat")
		assert.True(t, wd.Search("..."))
		assert.False(t, wd.Search(".."))
		assert.False(t, wd.Search("...."))
	})

	t.Run("wildcard_at_end", func(t *testing.T) {
		wd := ConstructorWordDictionary()
		wd.AddWord("abc")
		wd.AddWord("abd")
		assert.True(t, wd.Search("ab."))
		assert.False(t, wd.Search("ab"))
		assert.False(t, wd.Search("ab.."))
	})

	t.Run("empty_dictionary", func(t *testing.T) {
		wd := ConstructorWordDictionary()
		assert.False(t, wd.Search("a"))
		assert.False(t, wd.Search("."))
	})

	t.Run("multiple_wildcards", func(t *testing.T) {
		wd := ConstructorWordDictionary()
		wd.AddWord("hello")
		assert.True(t, wd.Search("h...."))
		assert.True(t, wd.Search(".e..."))
		assert.True(t, wd.Search("....o"))
		assert.True(t, wd.Search("....."))
		assert.False(t, wd.Search("h..."))  // too short
		assert.False(t, wd.Search("h.....")) // too long
	})
}
