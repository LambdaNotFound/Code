package hashmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_twoSum(t *testing.T) {
    tests := []struct {
        name     string
        nums     []int
        target   int
        expected []int
    }{
        {name: "basic_case", nums: []int{2, 7, 11, 15}, target: 9, expected: []int{0, 1}},
        {name: "another_case", nums: []int{3, 2, 4}, target: 6, expected: []int{1, 2}},
        {name: "repeated_numbers", nums: []int{3, 3}, target: 6, expected: []int{0, 1}},
        {name: "larger_array", nums: []int{1, 2, 3, 4, 5, 6}, target: 11, expected: []int{4, 5}},
        {name: "negatives", nums: []int{-1, -2, -3, -4, -5}, target: -8, expected: []int{2, 4}},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := twoSum(tt.nums, tt.target)

            assert.ElementsMatch(t, tt.expected, result)
        })
    }
}

func Test_isAnagram(t *testing.T) {
    tests := []struct {
        name     string
        s        string
        t        string
        expected bool
    }{
        {name: "simple_true", s: "anagram", t: "nagaram", expected: true},
        {name: "simple_false", s: "rat", t: "car", expected: false},
        {name: "different_lengths", s: "abc", t: "ab", expected: false},
        {name: "empty_strings", s: "", t: "", expected: true},
        {name: "unicode_false", s: "a", t: "á", expected: false},
        {name: "case_sensitive", s: "a", t: "A", expected: false},
        {name: "long_true", s: "abcdabcd", t: "dcbaabcd", expected: true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := isAnagram(tt.s, tt.t)
            assert.Equal(t, tt.expected, result)
        })
    }
}

func Test_canConstruct(t *testing.T) {
    tests := []struct {
        name     string
        s        string
        t        string
        expected bool
    }{
        {name: "simple_true", s: "anagram", t: "nagaram", expected: true},
        {name: "simple_false", s: "rat", t: "car", expected: false},
        {name: "different_lengths", s: "abc", t: "ab", expected: false},
        {name: "empty_strings", s: "", t: "", expected: true},
        {name: "unicode_false", s: "a", t: "á", expected: false},
        {name: "case_sensitive", s: "a", t: "A", expected: false},
        {name: "long_true", s: "abcdabcd", t: "dcbaabcd", expected: true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := canConstruct(tt.s, tt.t)
            assert.Equal(t, tt.expected, result)
        })
    }
}