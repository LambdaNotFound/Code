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
        target        int
        nums        []int
        expected int
    }{
        {
            "case 1",
            7,
            []int{2,3,1,2,4,3},
            2,
        },
        {
            "case 1",
            4,
            []int{1,4,4},
            1,
        },
        {
            "case 1",
            11,
            []int{1,1,1,1,1,1,1,1},
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
