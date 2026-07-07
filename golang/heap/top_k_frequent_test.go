package heap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var topKFrequentCases = []struct {
	name     string
	nums     []int
	k        int
	expected []int
}{
	{name: "leetcode_example1", nums: []int{1, 1, 1, 2, 2, 3}, k: 2, expected: []int{1, 2}},
	{name: "leetcode_example2", nums: []int{1}, k: 1, expected: []int{1}},
	{name: "all_same", nums: []int{5, 5, 5, 5}, k: 1, expected: []int{5}},
	{name: "k_equals_len", nums: []int{1, 2, 3}, k: 3, expected: []int{1, 2, 3}},
	{name: "negatives", nums: []int{-1, -1, -2, -2, -2}, k: 1, expected: []int{-2}},
	{name: "two_tied_pick_both", nums: []int{1, 1, 2, 2}, k: 2, expected: []int{1, 2}},
	// Distinct frequencies for every element so the result is deterministic:
	// 4→freq3, 3→freq2, 2→freq1. No ties at the k-th boundary.
	{name: "large_k", nums: []int{4, 4, 4, 3, 3, 2}, k: 3, expected: []int{4, 3, 2}},
}

func Test_topKFrequent(t *testing.T) {
	for _, tt := range topKFrequentCases {
		t.Run(tt.name, func(t *testing.T) {
			assert.ElementsMatch(t, tt.expected, topKFrequent(tt.nums, tt.k))
		})
	}
}

func Test_topKFrequentClaude(t *testing.T) {
	for _, tt := range topKFrequentCases {
		t.Run(tt.name, func(t *testing.T) {
			assert.ElementsMatch(t, tt.expected, topKFrequentClaude(tt.nums, tt.k))
		})
	}
}

func Test_sanitize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "letters_only", input: "abc", expected: "abc"},
		{name: "mixed", input: "a!@#b!@c", expected: "abc"},
		{name: "leading_specials", input: "!@x#@yz", expected: "xyz"},
		{name: "trailing_specials", input: "def!@#", expected: "def"},
		{name: "all_specials", input: "!@#$%", expected: ""},
		{name: "empty", input: "", expected: ""},
		{name: "preserves_case", input: "Hello!World", expected: "HelloWorld"},
		{name: "interleaved_specials", input: "d#@!e!#@f!@", expected: "def"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, sanitize(tt.input))
		})
	}
}
