package hashmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
