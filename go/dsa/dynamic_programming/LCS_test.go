package dynamic_programming

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_longestCommonSubsequence(t *testing.T) {
	tests := []struct {
		name     string
		text1    string
		text2    string
		expected int
	}{
		{name: "example1", text1: "abcde", text2: "ace", expected: 3},
		{name: "example2", text1: "abc", text2: "abc", expected: 3},
		{name: "example3", text1: "abc", text2: "def", expected: 0},
		{name: "one_char_match", text1: "a", text2: "a", expected: 1},
		{name: "one_char_no_match", text1: "a", text2: "b", expected: 0},
		{name: "subsequence_interleaved", text1: "abcba", text2: "abcbcba", expected: 5},
		{name: "empty_like_single", text1: "z", text2: "z", expected: 1},
		{name: "long_lcs", text1: "AGGTAB", text2: "GXTXAYB", expected: 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, longestCommonSubsequence(tt.text1, tt.text2))
		})
	}
}

func Test_findLength(t *testing.T) {
	tests := []struct {
		name     string
		nums1    []int
		nums2    []int
		expected int
	}{
		{name: "example1", nums1: []int{1, 2, 3, 2, 1}, nums2: []int{3, 2, 1, 4, 7}, expected: 3},
		{name: "example2", nums1: []int{0, 0, 0, 0, 0}, nums2: []int{0, 0, 0, 0, 0}, expected: 5},
		{name: "no_common", nums1: []int{1, 2, 3}, nums2: []int{4, 5, 6}, expected: 0},
		{name: "single_match", nums1: []int{1, 2, 3}, nums2: []int{3, 4, 5}, expected: 1},
		{name: "identical_arrays", nums1: []int{1, 2, 3, 4}, nums2: []int{1, 2, 3, 4}, expected: 4},
		{name: "match_at_end", nums1: []int{1, 2, 3}, nums2: []int{4, 2, 3}, expected: 2},
		{name: "match_at_start", nums1: []int{1, 2, 3}, nums2: []int{1, 2, 5}, expected: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, findLength(tt.nums1, tt.nums2))
		})
	}
}

/**
 * 1048. Longest String Chain
 *
 * A word A is a predecessor of B if you can insert exactly one letter anywhere
 * in A (without reordering existing letters) to make B. The chain length is the
 * number of words in the longest chain w1 → w2 → … → wk where each word is a
 * predecessor of the next.
 *
 * Algorithm: sort by word length, then for each word try removing every single
 * character to find a shorter predecessor already seen in dp[].
 *
 * Test strategy:
 *   - LeetCode canonical examples (including the "no chain possible" case).
 *   - Degenerate inputs: single word, all same length, all length-1 words.
 *   - Prefix chains (trivially recognisable) and non-prefix chains (predecessor
 *     is found by removing a middle or leading character).
 *   - Input given in reverse-length order to verify the sort is load-bearing.
 *   - Multiple chains of different depths to confirm the max is taken.
 *   - Length gap: words separated by 2+ characters cannot form a chain even
 *     though they overlap in content.
 */
func Test_longestStrChain(t *testing.T) {
	tests := []struct {
		name     string
		words    []string
		expected int
	}{
		// -----------------------------------------------------------------------
		// LeetCode canonical examples
		// -----------------------------------------------------------------------
		{
			// chain: "a" → "ba" → "bda" → "bdca"
			name:     "leetcode_example1",
			words:    []string{"a", "b", "ba", "bca", "bda", "bdca"},
			expected: 4,
		},
		{
			// chain: "xb" → "xbc" → "cxbc" → "pcxbc" → "pcxbcf"
			name:     "leetcode_example2",
			words:    []string{"xbc", "pcxbcf", "xb", "cxbc", "pcxbc"},
			expected: 5,
		},
		{
			// "abcd" and "dbqca" differ by more than one insertion — no chain.
			name:     "leetcode_example3_no_chain",
			words:    []string{"abcd", "dbqca"},
			expected: 1,
		},

		// -----------------------------------------------------------------------
		// Degenerate / edge inputs
		// -----------------------------------------------------------------------
		{
			// A single word has a trivial chain of length 1.
			name:     "single_word",
			words:    []string{"a"},
			expected: 1,
		},
		{
			// Words of the same length can never be predecessors of each other.
			name:     "all_same_length_no_chain",
			words:    []string{"ab", "cd", "ef"},
			expected: 1,
		},
		{
			// Multiple length-1 words: no word has a shorter predecessor, so the
			// longest chain is 1.
			name:     "length_1_words_only",
			words:    []string{"a", "b", "c"},
			expected: 1,
		},
		{
			// Words exist at lengths 1, 3, 5 — the length gap of 2 means no word
			// can be a direct predecessor of another.
			name:     "length_gap_prevents_chain",
			words:    []string{"a", "abc", "abcde"},
			expected: 1,
		},
		{
			// No predecessor relationships across completely disjoint character sets.
			name:     "no_chain_across_disjoint_words",
			words:    []string{"abc", "def", "ghij"},
			expected: 1,
		},

		// -----------------------------------------------------------------------
		// Predecessor found by removing different positions
		// -----------------------------------------------------------------------
		{
			// Predecessor found by removing the first character each time:
			// "a" → "ba" → "cba" → "dcba"
			name:     "remove_leading_char",
			words:    []string{"a", "ba", "cba", "dcba"},
			expected: 4,
		},
		{
			// Predecessor found by removing the middle character:
			// "ab" → "aXb" (remove index 1 of "aXb" gives "ab")
			//        "aXb" → "aXYb" (remove index 2 of "aXYb" gives "aXb")
			name:     "remove_middle_char",
			words:    []string{"ab", "aXb", "aXYb"},
			expected: 3,
		},

		// -----------------------------------------------------------------------
		// Sorting is load-bearing
		// -----------------------------------------------------------------------
		{
			// Identical to leetcode_example1 but supplied in reverse length order.
			// The sort-by-length step must fire before dp is computed.
			name:     "input_in_reverse_length_order",
			words:    []string{"bdca", "bda", "ba", "a"},
			expected: 4,
		},

		// -----------------------------------------------------------------------
		// Linear prefix chain
		// -----------------------------------------------------------------------
		{
			// Simplest possible chain: each word extends the previous by one char.
			// "a" → "ab" → "abc" → "abcd"
			name:     "linear_chain_by_prefix",
			words:    []string{"a", "ab", "abc", "abcd"},
			expected: 4,
		},

		// -----------------------------------------------------------------------
		// Multiple chains — the longest must be returned
		// -----------------------------------------------------------------------
		{
			// Two chains share the root "a":
			//   "a" → "ab" → "abc"          (length 3)
			//   "a" → "ab" → "abc" → "abcd" (length 4)  ← winner
			// Also "b" → "bc" → "bcd" → "bcde" (length 4, tied).
			name:     "multiple_chains_pick_longest",
			words:    []string{"a", "ab", "abc", "b", "bc", "bcd", "bcde"},
			expected: 4,
		},
		{
			// The chain "a" → "ab" → "abc" → "abcd" reaches length 4.
			// "ac" is a dead-end side branch of depth 2.
			name:     "chain_through_branching_paths",
			words:    []string{"a", "ab", "ac", "abc", "abcd"},
			expected: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, longestStrChain(tt.words))
		})
	}
}
