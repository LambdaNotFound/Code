package slidingwindow

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_lengthOfLongestSubstring(t *testing.T) {
	asciiCases := []struct {
		name     string
		str      string
		expected int
	}{
		{"case 1", "abcabcbb", 3},
		{"case 2", "bbbbb", 1},
		{"case 3", "pwwkew", 3},
		{"case 4", "abc!@#abc", 6},
		{"empty", "", 0},
	}

	for _, tc := range asciiCases {
		t.Run(tc.name, func(t *testing.T) {
			// byte-based version — correct for ASCII input
			assert.Equal(t, tc.expected, lengthOfLongestSubstring(tc.str))
		})
	}
}

func Test_lengthOfLongestSubstringRune(t *testing.T) {
	testCases := []struct {
		name     string
		str      string
		expected int
	}{
		{"case 1", "abcabcbb", 3},
		{"case 2", "bbbbb", 1},
		{"case 3", "pwwkew", 3},
		{"case 4", "abc!@#abc", 6},
		{"unicode", "你好世界", 4}, // multi-byte chars — rune version handles correctly
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, lengthOfLongestSubstringRune(tc.str))
		})
	}
}

func Test_minWindow(t *testing.T) {
	testCases := []struct {
		name     string
		s        string
		t        string
		expected string
	}{
		{
			"case 1",
			"ADOBECODEBANC",
			"ABC",
			"BANC",
		},
		{
			"case 2",
			"a",
			"a",
			"a",
		},
		{
			"case 3",
			"a",
			"aa",
			"",
		},
		{
			"duplicates_in_t",
			"ABCDE",
			"AA",
			"",
		},
		{
			"t_longer",
			"AB",
			"ABC",
			"",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := minWindow(tc.s, tc.t)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func Test_minSubArrayLen(t *testing.T) {
	testCases := []struct {
		name     string
		target   int
		nums     []int
		expected int
	}{
		{
			"case 1",
			7,
			[]int{2, 3, 1, 2, 4, 3},
			2,
		},
		{
			"case 1",
			4,
			[]int{1, 4, 4},
			1,
		},
		{
			"case 1",
			11,
			[]int{1, 1, 1, 1, 1, 1, 1, 1},
			0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := minSubArrayLen(tc.target, tc.nums)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func Test_lengthOfLongestSubstringKDistinct(t *testing.T) {
	testCases := []struct {
		name     string
		s        string
		k        int
		expected int
	}{
		{name: "basic_eceba", s: "eceba", k: 2, expected: 3}, // "ece"
		{name: "basic_double_a", s: "aa", k: 1, expected: 2}, // "aa"
		{name: "single_char_k1", s: "a", k: 1, expected: 1},  // single char
		{name: "single_char_k2", s: "a", k: 2, expected: 1},  // still just 1 char

		// Edge cases
		{name: "empty_string", s: "", k: 2, expected: 0},              // empty string
		{name: "k_zero", s: "abc", k: 0, expected: 0},                 // k=0 means no substring
		{name: "k_greater_than_unique", s: "abc", k: 10, expected: 3}, // k > unique chars → whole string

		// Repeated patterns
		{name: "repeated_pattern", s: "abcadcacacaca", k: 3, expected: 11}, // "cadcacacaca"
		{name: "abaccc", s: "abaccc", k: 2, expected: 4},                   // "accc"

		// Larger repetition
		{name: "long_repetition", s: "aaaaabbbbbccccc", k: 2, expected: 10}, // "aaaaabbbbb"
		{name: "abcabcabc", s: "abcabcabc", k: 2, expected: 2},              // any "ab" or "bc"

		// Long alternating
		{name: "long_alternating", s: "abababababab", k: 2, expected: 12}, // whole string
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := lengthOfLongestSubstringKDistinct(tc.s, tc.k)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func Test_findAnagrams(t *testing.T) {
	testCases := []struct {
		name     string
		s        string
		p        string
		expected []int
	}{
		{
			"case 1",
			"cbaebabacd",
			"abc",
			[]int{0, 6},
		},
		{
			"case 2",
			"abab",
			"ab",
			[]int{0, 1, 2},
		},
		{
			"no_match",
			"af",
			"be",
			[]int{},
		},
		{
			"p_longer_than_s",
			"ab",
			"abc",
			[]int{},
		},
		{
			"whole_string_match",
			"abc",
			"bca",
			[]int{0},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := findAnagrams(tc.s, tc.p)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func Test_characterReplacement(t *testing.T) {
	testCases := []struct {
		name     string
		s        string
		k        int
		expected int
	}{
		{"case 1", "ABAB", 2, 4},
		{"case 2", "AABABBA", 1, 4},
		{"all_same", "AAAA", 0, 4},
		{"k_zero", "ABCD", 0, 1},
		{"single_char", "A", 1, 1},
		{"whole_string", "ABCD", 4, 4},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := characterReplacement(tc.s, tc.k)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func Test_containsNearbyDuplicate(t *testing.T) {
	testCases := []struct {
		name     string
		nums     []int
		k        int
		expected bool
	}{
		{"leetcode_example1", []int{1, 2, 3, 1}, 3, true},
		{"leetcode_example2", []int{1, 0, 1, 1}, 1, true},
		{"leetcode_example3", []int{1, 2, 3, 1, 2, 3}, 2, false},
		{"empty", []int{}, 1, false},
		{"single_element", []int{1}, 0, false},
		{"k_zero_no_duplicate", []int{1, 1}, 0, false},
		{"window_exact_k", []int{1, 2, 1}, 2, true},
		{"window_just_outside_k", []int{1, 2, 3, 1}, 2, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, containsNearbyDuplicate(tc.nums, tc.k))
		})
	}
}

func Test_findHighAccessEmployees(t *testing.T) {
	testCases := []struct {
		name        string
		accessTimes [][]string
		expected    []string
	}{
		{
			"leetcode_example",
			[][]string{{"a", "0549"}, {"b", "0457"}, {"a", "0532"}, {"a", "0621"}, {"b", "0540"}},
			[]string{"a"},
		},
		{
			"no_high_access",
			[][]string{{"a", "0800"}, {"a", "1400"}, {"a", "2000"}},
			[]string{},
		},
		{
			"boundary_exactly_one_hour",
			[][]string{{"a", "0100"}, {"a", "0159"}, {"a", "0200"}},
			[]string{},
		},
		{
			"boundary_within_one_hour",
			[][]string{{"a", "0100"}, {"a", "0159"}, {"a", "0159"}},
			[]string{"a"},
		},
		{
			"multiple_high_access",
			[][]string{{"a", "0100"}, {"a", "0110"}, {"a", "0120"}, {"b", "0200"}, {"b", "0210"}, {"b", "0220"}},
			[]string{"a", "b"},
		},
		{
			"fewer_than_3_accesses",
			[][]string{{"a", "0100"}, {"a", "0110"}},
			[]string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := findHighAccessEmployees(tc.accessTimes)
			assert.ElementsMatch(t, tc.expected, result)
		})
	}
}

func Test_longestSubstring(t *testing.T) {
	testCases := []struct {
		name     string
		s        string
		k        int
		expected int
	}{
		{"leetcode_example1", "aaabb", 3, 3},    // "aaa"
		{"leetcode_example2", "ababbc", 2, 5},   // "ababb"
		{"all_satisfy", "aaabbb", 1, 6},         // whole string
		{"none_satisfy", "abcdef", 2, 0},        // every char appears once
		{"empty_string", "", 1, 0},
		{"k_zero", "abc", 0, 3},                 // k=0: every char "satisfies"
		{"single_char_meets_k", "aaaa", 4, 4},
		{"single_char_short", "aaa", 4, 0},      // "a" only appears 3 times, k=4
		{"mixed", "ababacb", 3, 0},              // no char appears ≥ 3 in a valid window
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, longestSubstring(tc.s, tc.k))
		})
	}
}

func Test_findSubstring(t *testing.T) {
	testCases := []struct {
		name     string
		s        string
		words    []string
		expected []int
	}{
		{
			"leetcode_example1",
			"barfoothefoobarman",
			[]string{"foo", "bar"},
			[]int{0, 9},
		},
		{
			"leetcode_example2",
			"wordgoodgoodgoodbestword",
			[]string{"word", "good", "best", "word"},
			[]int{},
		},
		{
			"leetcode_example3",
			"barfoofoobarthefoobarman",
			[]string{"bar", "foo", "the"},
			[]int{6, 9, 12},
		},
		{
			"single_word",
			"aaa",
			[]string{"a"},
			[]int{0, 1, 2},
		},
		{
			"duplicate_words",
			"aaaa",
			[]string{"aa", "aa"},
			[]int{0},
		},
		{
			"no_match",
			"lingmindraboofooowingdrabo",
			[]string{"fooo", "barr", "wing", "ding", "wing"},
			[]int{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := findSubstring(tc.s, tc.words)
			assert.ElementsMatch(t, tc.expected, result)
		})
	}
}
