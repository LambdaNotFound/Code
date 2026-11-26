package prefixtree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_shortestSubstrings(t *testing.T) {
	t.Run("leetcode_example_1", func(t *testing.T) {
		// "ab" < "ca" lexicographically, both length 2 and unique to "cab"
		// "ad" is a substring of "bad", so "ad" has no unique substring
		// "ba" is unique to "bad" (length 2, shorter than "bad")
		// "c" appears in "cab", so "c" has no unique substring
		arr := []string{"cab", "ad", "bad", "c"}
		expected := []string{"ab", "", "ba", ""}
		assert.Equal(t, expected, shortestSubstrings(arr))
	})

	t.Run("leetcode_example_2", func(t *testing.T) {
		// every substring of "abc" and "bcd" also appears in "abcd"
		// "abcd" itself is the only substring unique to arr[2]
		arr := []string{"abc", "bcd", "abcd"}
		expected := []string{"", "", "abcd"}
		assert.Equal(t, expected, shortestSubstrings(arr))
	})

	t.Run("single_string", func(t *testing.T) {
		// Every substring is unique since there's only one string; shortest is first char
		arr := []string{"abc"}
		result := shortestSubstrings(arr)
		assert.Equal(t, 1, len(result[0]), "shortest unique substring of a lone string should be length 1")
		assert.Equal(t, "a", result[0])
	})

	t.Run("identical_strings", func(t *testing.T) {
		// Both strings are the same — no substring is unique to either
		arr := []string{"abc", "abc"}
		expected := []string{"", ""}
		assert.Equal(t, expected, shortestSubstrings(arr))
	})

	t.Run("one_char_strings_distinct", func(t *testing.T) {
		arr := []string{"a", "b", "c"}
		expected := []string{"a", "b", "c"}
		assert.Equal(t, expected, shortestSubstrings(arr))
	})

	t.Run("one_char_strings_duplicate", func(t *testing.T) {
		arr := []string{"a", "a"}
		expected := []string{"", ""}
		assert.Equal(t, expected, shortestSubstrings(arr))
	})

	t.Run("substring_overlap", func(t *testing.T) {
		// every substring of "ab" also appears in "abc", so "ab" has no unique substring
		// "c" (length 1) is unique to "abc" — shorter than "bc" or "abc"
		arr := []string{"ab", "abc"}
		expected := []string{"", "c"}
		assert.Equal(t, expected, shortestSubstrings(arr))
	})

	t.Run("repeated_chars_in_string", func(t *testing.T) {
		// "aaa" — all its substrings are "a", "aa", "aaa".
		// "bbb" has "b","bb","bbb"; none overlap with "aaa".
		arr := []string{"aaa", "bbb"}
		expected := []string{"a", "b"}
		assert.Equal(t, expected, shortestSubstrings(arr))
	})

	t.Run("lexicographic_tiebreak", func(t *testing.T) {
		// "ba" and "ab" both have length-2 unique substrings; lexicographically smaller wins
		// arr[0]="ba": substrings "b","a","ba" — "b" not in "ac","ca"; "a" appears in "ac","ca"
		// arr[1]="ac": "a" in "ba","ca"; "c" in "ca"; "ac" unique
		// arr[2]="ca": "c" in "ac"; "a" in "ac","ba"; "ca" unique
		arr := []string{"ba", "ac", "ca"}
		result := shortestSubstrings(arr)
		// "b" is only in "ba"
		assert.Equal(t, "b", result[0])
		// "ac" is the shortest unique for arr[1]
		assert.Equal(t, "ac", result[1])
		// "ca" is the shortest unique for arr[2]
		assert.Equal(t, "ca", result[2])
	})
}
