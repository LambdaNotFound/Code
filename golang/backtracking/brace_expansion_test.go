package backtracking

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_expand(t *testing.T) {
	testCases := []struct {
		name     string
		s        string
		expected []string
	}{
		{
			name:     "leetcode_example_1",
			s:        "{a,b}c{d,e}f",
			expected: []string{"acdf", "acef", "bcdf", "bcef"},
		},
		{
			name:     "leetcode_example_2",
			s:        "abcd",
			expected: []string{"abcd"},
		},
		{
			name:     "single_char",
			s:        "a",
			expected: []string{"a"},
		},
		{
			name:     "single_brace_group",
			s:        "{a,b,c}",
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "brace_group_sorted",
			s:        "{c,a,b}",
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "adjacent_brace_groups",
			s:        "{a,b}{c,d}",
			expected: []string{"ac", "ad", "bc", "bd"},
		},
		{
			name:     "three_brace_groups",
			s:        "{a,b,c}{d,e}{f,g}",
			expected: []string{"adf", "adg", "aef", "aeg", "bdf", "bdg", "bef", "beg", "cdf", "cdg", "cef", "ceg"},
		},
		{
			name:     "literal_between_braces",
			s:        "{a,b,c}d{e,f}",
			expected: []string{"ade", "adf", "bde", "bdf", "cde", "cdf"},
		},
		{
			name:     "all_literals",
			s:        "xyz",
			expected: []string{"xyz"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := expand(tc.s)
			assert.Equal(t, tc.expected, result)
		})
	}
}
