package interview

import (
	prefixtree "gocode/pkg/prefix_tree"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWordTrie_InsertSearch(t *testing.T) {
	testCases := []struct {
		name        string
		bannedWords []string
		text        string
		expected    string
	}{
		{
			name:        "basic trie behavior mirrors unsafe filter single word",
			bannedWords: []string{"hate", "spam"},
			text:        "this hate and spam here",
			expected:    "this **** and **** here",
		},
	}

	// This test now validates that the underlying trie behavior (via prefix_tree.Trie)
	// is consistent with what filterUnsafeWords expects, rather than relying on a
	// non-existent WordTrie type.
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			trie := prefixtree.ConstructorTrie()
			for _, word := range tc.bannedWords {
				trie.Insert(strings.ToLower(word))
			}

			words := strings.Fields(tc.text)
			for i, w := range words {
				clean := strings.ToLower(strings.Trim(w, ".,!?;:\"'"))
				if trie.Search(clean) {
					words[i] = "****"
				}
			}

			assert.Equal(t, tc.expected, strings.Join(words, " "))
		})
	}
}

func Test_filterUnsafeWords(t *testing.T) {
	testCases := []struct {
		name        string
		text        string
		bannedWords []string
		expected    string
	}{
		{
			name:        "single banned word replaced",
			text:        "this is hate today",
			bannedWords: []string{"hate"},
			expected:    "this is **** today",
		},
		{
			name:        "multiple banned words replaced",
			text:        "spam and hate here",
			bannedWords: []string{"spam", "hate"},
			expected:    "**** and **** here",
		},
		{
			name:        "case insensitive match",
			text:        "This is HATE and Spam",
			bannedWords: []string{"hate", "spam"},
			expected:    "This is **** and ****",
		},
		{
			name:        "word with punctuation stripped before lookup",
			text:        "this is hate, really",
			bannedWords: []string{"hate"},
			expected:    "this is **** really", // "hate," token replaced entirely, comma lost
		},
		{
			name:        "non-banned words unchanged",
			text:        "hello world",
			bannedWords: []string{"hate"},
			expected:    "hello world",
		},
		{
			name:        "empty banned list changes nothing",
			text:        "hello world",
			bannedWords: []string{},
			expected:    "hello world",
		},
		{
			name:        "empty text returns empty",
			text:        "",
			bannedWords: []string{"hate"},
			expected:    "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := filterUnsafeWords(tc.text, tc.bannedWords)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func Test_filterUnsafePhrases(t *testing.T) {
	testCases := []struct {
		name          string
		text          string
		bannedPhrases []string
		expected      string
	}{
		{
			name:          "single phrase in middle replaced",
			text:          "this hate speech is bad",
			bannedPhrases: []string{"hate speech"},
			expected:      "this **** **** is bad",
		},
		{
			name:          "phrase at start replaced",
			text:          "hate speech is wrong",
			bannedPhrases: []string{"hate speech"},
			expected:      "**** **** is wrong",
		},
		{
			name:          "phrase at end replaced",
			text:          "this is hate speech",
			bannedPhrases: []string{"hate speech"},
			expected:      "this is **** ****",
		},
		{
			name:          "multiple phrases replaced",
			text:          "this hate speech is very bad indeed today",
			bannedPhrases: []string{"hate speech", "very bad indeed"},
			expected:      "this **** **** is **** **** **** today",
		},
		{
			name:          "case insensitive phrase match",
			text:          "This HATE SPEECH is here",
			bannedPhrases: []string{"hate speech"},
			expected:      "This **** **** is here",
		},
		{
			name:          "no match leaves text unchanged",
			text:          "this is perfectly fine",
			bannedPhrases: []string{"hate speech"},
			expected:      "this is perfectly fine",
		},
		{
			name:          "single word phrase",
			text:          "this spam is bad",
			bannedPhrases: []string{"spam"},
			expected:      "this **** is bad",
		},
		{
			name:          "phrase with punctuation in text",
			text:          "this hate speech, really",
			bannedPhrases: []string{"hate speech"},
			expected:      "this **** **** really", // "speech," token replaced entirely, comma lost
		},
		{
			name:          "empty banned phrases changes nothing",
			text:          "hello world",
			bannedPhrases: []string{},
			expected:      "hello world",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := filterUnsafePhrases(tc.text, tc.bannedPhrases)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func Test_PhraseTrie(t *testing.T) {
	testCases := []struct {
		name    string
		inserts []string
		queries []struct {
			phrase   string
			contains bool
		}
		prefixQueries []struct {
			words    []string
			isPrefix bool
		}
	}{
		{
			name:    "exact phrase found",
			inserts: []string{"hate speech"},
			queries: []struct {
				phrase   string
				contains bool
			}{
				{"hate speech", true},
				{"hate", false},        // prefix only, not full phrase
				{"speech hate", false}, // wrong order
			},
			prefixQueries: []struct {
				words    []string
				isPrefix bool
			}{
				{[]string{"hate"}, true},           // valid prefix
				{[]string{"hate", "speech"}, true}, // full phrase is also a valid prefix
				{[]string{"speech"}, false},        // not a prefix
			},
		},
		{
			name:    "multiple phrases",
			inserts: []string{"hate speech", "very bad indeed"},
			queries: []struct {
				phrase   string
				contains bool
			}{
				{"hate speech", true},
				{"very bad indeed", true},
				{"very bad", false}, // prefix only
				{"hate", false},
			},
			prefixQueries: []struct {
				words    []string
				isPrefix bool
			}{
				{[]string{"very"}, true},
				{[]string{"very", "bad"}, true},
				{[]string{"bad"}, false},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			trie := NewPhraseTrie()
			for _, phrase := range tc.inserts {
				trie.Insert(phrase)
			}

			for _, q := range tc.queries {
				assert.Equal(t, q.contains, trie.Contains(q.phrase), "Contains(%q)", q.phrase)
			}

			for _, q := range tc.prefixQueries {
				assert.Equal(t, q.isPrefix, trie.ContainsPrefix(q.words), "ContainsPrefix(%v)", q.words)
			}
		})
	}
}
