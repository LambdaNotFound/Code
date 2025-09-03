package two_pointers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_lengthOfLongestSubstring(t *testing.T) {
    testCases := []struct {
        name     string
        str      string
        expected int
    }{
        {
            "case 1",
            "abcabcbb",
            3,
        },
        {
            "case 2",
            "bbbbb",
            1,
        },
        {
            "case 3",
            "pwwkew",
            3,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := lengthOfLongestSubstring(tc.str)
            assert.Equal(t, tc.expected, result)

            result = lengthOfLongestSubstring_optimized(tc.str)
            assert.Equal(t, tc.expected, result)
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
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := findAnagrams(tc.s, tc.p)
            assert.Equal(t, tc.expected, result)
        })
    }
}
