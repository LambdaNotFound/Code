package apidesign

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/**
 * 745. Prefix and Suffix Search
 *
 * Test strategy:
 *   - LeetCode canonical example (single word "apple").
 *   - Exact full-word match for both prefix and suffix.
 *   - No-match cases: prefix exists but suffix doesn't, and vice versa.
 *   - Multiple words where the same (pref, suff) pair matches more than one word
 *     — the function must return the largest index.
 *   - Two distinct words with non-overlapping pref/suff combinations.
 *   - Single-character word: degenerate case where pref == suff == the word itself.
 *   - charToIndex coverage: '#' → 26, lowercase letter → 0–25.
 */
func Test_WordFilter(t *testing.T) {
	t.Run("leetcode_example_single_word", func(t *testing.T) {
		wf := ConstructorSearch([]string{"apple"})
		assert.Equal(t, 0, wf.F("a", "e"))        // prefix "a", suffix "e"
		assert.Equal(t, 0, wf.F("apple", "apple")) // exact full match
		assert.Equal(t, 0, wf.F("ap", "le"))       // interior prefix and suffix
		assert.Equal(t, -1, wf.F("b", "e"))        // prefix not in dictionary
		assert.Equal(t, -1, wf.F("a", "b"))        // suffix not in dictionary
	})

	t.Run("return_largest_index_on_tie", func(t *testing.T) {
		// Both words share prefix "a" and suffix "b"; largest index must win.
		wf := ConstructorSearch([]string{"ab", "ab"})
		assert.Equal(t, 1, wf.F("a", "b"))
		assert.Equal(t, 1, wf.F("ab", "ab"))
	})

	t.Run("two_distinct_words", func(t *testing.T) {
		// "abc" at index 0, "xyz" at index 1 — no pref/suff overlap between them.
		wf := ConstructorSearch([]string{"abc", "xyz"})
		assert.Equal(t, 0, wf.F("a", "c"))  // only "abc" matches
		assert.Equal(t, 1, wf.F("x", "z"))  // only "xyz" matches
		assert.Equal(t, -1, wf.F("a", "z")) // cross-word: no single word satisfies both
	})

	t.Run("longer_suffix_supersedes_earlier_index", func(t *testing.T) {
		// "caab" (0) and "ccbca" (1): F("c","a") — "ccbca" ends with 'a', "caab" ends with 'b'.
		// Only index 1 qualifies; expect 1.
		wf := ConstructorSearch([]string{"caab", "ccbca"})
		assert.Equal(t, 1, wf.F("c", "a"))
		assert.Equal(t, 0, wf.F("c", "b")) // only "caab" ends with 'b'
		assert.Equal(t, -1, wf.F("z", "a"))
	})

	t.Run("single_char_word", func(t *testing.T) {
		wf := ConstructorSearch([]string{"a"})
		assert.Equal(t, 0, wf.F("a", "a"))  // pref == suff == word
		assert.Equal(t, -1, wf.F("b", "a")) // wrong prefix
		assert.Equal(t, -1, wf.F("a", "b")) // wrong suffix
	})

	t.Run("multiple_words_prefer_latest", func(t *testing.T) {
		// "ba" at 0, "ab" at 1, "ab" at 2. F("a","b") should return 2 (latest "ab").
		wf := ConstructorSearch([]string{"ba", "ab", "ab"})
		assert.Equal(t, 2, wf.F("a", "b"))
		assert.Equal(t, 0, wf.F("b", "a")) // only "ba" matches
	})
}
