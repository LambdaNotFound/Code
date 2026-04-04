package stack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_backspaceCompare(t *testing.T) {
	tests := []struct {
		name     string
		s, t_str string
		expected bool
	}{
		{"leetcode_1", "ab#c", "ad#c", true},
		{"leetcode_2", "ab##", "c#d#", true},
		{"leetcode_3", "a#c", "b", false},
		{"both_empty", "###", "##", true},
		{"no_backspace", "abc", "abc", true},
		{"different", "abc", "def", false},
		{"backspace_at_start", "#a", "a", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, backspaceCompare(tt.s, tt.t_str))
		})
	}
}

func Test_asteroidCollision(t *testing.T) {
	tests := []struct {
		name      string
		asteroids []int
		expected  []int
	}{
		{"leetcode_1", []int{5, 10, -5}, []int{5, 10}},
		{"leetcode_2", []int{8, -8}, []int{}},
		{"leetcode_3", []int{10, 2, -5}, []int{10}},
		{"no_collision", []int{1, 2, 3}, []int{1, 2, 3}},
		{"all_left", []int{-1, -2, -3}, []int{-1, -2, -3}},
		{"mutual_destroy", []int{-2, -1, 1, 2}, []int{-2, -1, 1, 2}},
		{"chain", []int{1, -2, 2, -1}, []int{-2, 2}},
		{"single", []int{5}, []int{5}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input1 := append([]int(nil), tt.asteroids...)
			input2 := append([]int(nil), tt.asteroids...)
			assert.Equal(t, tt.expected, asteroidCollision(input1))
			assert.Equal(t, tt.expected, asteroidCollisionCalude(input2))
		})
	}
}

func Test_removeDuplicateLetters(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		expected string
	}{
		{"leetcode_1", "bcabc", "abc"},
		{"leetcode_2", "cbacdcbc", "acdb"},
		{"single_char", "a", "a"},
		{"all_same", "aaaa", "a"},
		{"already_lex", "abc", "abc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, removeDuplicateLetters(tt.s))
		})
	}
}

func Test_removeDuplicatesK(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		k        int
		expected string
	}{
		{"leetcode_1", "abcd", 2, "abcd"},
		{"leetcode_2", "deeedbbcccbdaa", 3, "aa"},
		{"leetcode_3", "pbbcggttciiippooaais", 2, "ps"},
		{"no_removal", "abc", 3, "abc"},
		{"all_removed", "aaa", 3, ""},
		{"nested", "aaabbbccc", 3, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, removeDuplicates(tt.s, tt.k), "removeDuplicates")
			assert.Equal(t, tt.expected, removeDuplicatesClaude(tt.s, tt.k), "removeDuplicatesClaude")
		})
	}
}

func Test_decodeString(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		expected string
	}{
		{"leetcode_1", "3[a]2[bc]", "aaabcbc"},
		{"leetcode_2", "3[a2[c]]", "accaccacc"},
		{"leetcode_3", "2[abc]3[cd]ef", "abcabccdcdcdef"},
		{"single_repeat", "1[a]", "a"},
		{"nested_repeat", "2[3[a]]", "aaaaaa"},
		{"no_brackets", "abc", "abc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, decodeString(tt.s))
		})
	}
}
