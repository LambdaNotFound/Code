package monostack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_dailyTemperatures(t *testing.T) {
	tests := []struct {
		name         string
		temperatures []int
		expected     []int
	}{
		{"leetcode_1", []int{73, 74, 75, 71, 69, 72, 76, 73}, []int{1, 1, 4, 2, 1, 1, 0, 0}},
		{"leetcode_2", []int{30, 40, 50, 60}, []int{1, 1, 1, 0}},
		{"leetcode_3", []int{30, 60, 90}, []int{1, 1, 0}},
		{"all_same", []int{70, 70, 70}, []int{0, 0, 0}},
		{"decreasing", []int{90, 80, 70}, []int{0, 0, 0}},
		{"single", []int{42}, []int{0}},
		{"valley", []int{50, 40, 60}, []int{2, 1, 0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input1 := append([]int(nil), tt.temperatures...)
			input2 := append([]int(nil), tt.temperatures...)
			assert.Equal(t, tt.expected, dailyTemperatures(input1), "dailyTemperatures")
			assert.Equal(t, tt.expected, dailyTemperaturesRightToLeft(input2), "dailyTemperaturesRightToLeft")
		})
	}
}
func Test_longestValidParentheses(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		expected int
	}{
		{"leetcode_1", "(()", 2},
		{"leetcode_2", ")()())", 4},
		{"empty", "", 0},
		{"all_open", "(((", 0},
		{"all_close", ")))", 0},
		{"full_match", "()()", 4},
		{"nested", "(())", 4},
		{"complex", "()(())", 6},
		{"trailing_open", "()((", 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, longestValidParentheses(tt.s), "stack")
			assert.Equal(t, tt.expected, longestValidParenthesesDP(tt.s), "dp")
			assert.Equal(t, tt.expected, longestValidParenthesesTwoPointers(tt.s), "two_pointers")
		})
	}
}

func Test_trap_extra(t *testing.T) {
	tests := []struct {
		name     string
		height   []int
		expected int
	}{
		{"all_zeros", []int{0, 0, 0}, 0},
		{"no_water", []int{1, 2, 3}, 0},
		{"flat", []int{3, 3, 3}, 0},
		{"single", []int{5}, 0},
		{"two_walls", []int{5, 0, 5}, 5},
		{"asymmetric", []int{4, 2, 0, 3, 2, 5}, 9},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input1 := append([]int(nil), tt.height...)
			input2 := append([]int(nil), tt.height...)
			assert.Equal(t, tt.expected, trap(input1), "trap")
			assert.Equal(t, tt.expected, trapSlice(input2), "trapSlice")
		})
	}
}
