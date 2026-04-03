package math

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func deepCopyInts(nums []int) []int {
	cp := make([]int, len(nums))
	copy(cp, nums)
	return cp
}

func deepCopyMatrix(m [][]int) [][]int {
	cp := make([][]int, len(m))
	for i, row := range m {
		cp[i] = make([]int, len(row))
		copy(cp[i], row)
	}
	return cp
}

func Test_rotateArray(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		k        int
		expected []int
	}{
		{name: "leetcode_example1", nums: []int{1, 2, 3, 4, 5, 6, 7}, k: 3, expected: []int{5, 6, 7, 1, 2, 3, 4}},
		{name: "leetcode_example2", nums: []int{-1, -100, 3, 99}, k: 2, expected: []int{3, 99, -1, -100}},
		{name: "k_zero", nums: []int{1, 2, 3}, k: 0, expected: []int{1, 2, 3}},
		{name: "k_equals_len", nums: []int{1, 2, 3}, k: 3, expected: []int{1, 2, 3}},
		{name: "k_greater_than_len", nums: []int{1, 2, 3}, k: 5, expected: []int{2, 3, 1}}, // 5%3=2 right-rotations
		{name: "single_element", nums: []int{42}, k: 7, expected: []int{42}},
		{name: "two_elements", nums: []int{1, 2}, k: 1, expected: []int{2, 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input1 := deepCopyInts(tt.nums)
			rotateArray(input1, tt.k)
			assert.Equal(t, tt.expected, input1, "rotateArray")

			input2 := deepCopyInts(tt.nums)
			rotateArraySlicesReverse(input2, tt.k)
			assert.Equal(t, tt.expected, input2, "rotateArraySlicesReverse")
		})
	}
}

func Test_rotateImage(t *testing.T) {
	tests := []struct {
		name     string
		matrix   [][]int
		expected [][]int
	}{
		{
			name:     "1x1",
			matrix:   [][]int{{1}},
			expected: [][]int{{1}},
		},
		{
			name:     "2x2",
			matrix:   [][]int{{1, 2}, {3, 4}},
			expected: [][]int{{3, 1}, {4, 2}},
		},
		{
			name: "leetcode_example1_3x3",
			matrix: [][]int{
				{1, 2, 3},
				{4, 5, 6},
				{7, 8, 9},
			},
			expected: [][]int{
				{7, 4, 1},
				{8, 5, 2},
				{9, 6, 3},
			},
		},
		{
			name: "leetcode_example2_4x4",
			matrix: [][]int{
				{5, 1, 9, 11},
				{2, 4, 8, 10},
				{13, 3, 6, 7},
				{15, 14, 12, 16},
			},
			expected: [][]int{
				{15, 13, 2, 5},
				{14, 3, 4, 1},
				{12, 6, 8, 9},
				{16, 7, 10, 11},
			},
		},
		{
			name:     "identity_3x3",
			matrix:   [][]int{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}},
			expected: [][]int{{0, 0, 1}, {0, 1, 0}, {1, 0, 0}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := deepCopyMatrix(tt.matrix)
			rotateImage(input)
			assert.Equal(t, tt.expected, input)
		})
	}
}
