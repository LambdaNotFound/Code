package solid_coding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func deepCopyMatrix(m [][]int) [][]int {
	cp := make([][]int, len(m))
	for i := range m {
		cp[i] = append([]int(nil), m[i]...)
	}
	return cp
}

func Test_setZeroes(t *testing.T) {
	tests := []struct {
		name     string
		matrix   [][]int
		expected [][]int
	}{
		{
			"leetcode_1",
			[][]int{{1, 1, 1}, {1, 0, 1}, {1, 1, 1}},
			[][]int{{1, 0, 1}, {0, 0, 0}, {1, 0, 1}},
		},
		{
			"leetcode_2",
			[][]int{{0, 1, 2, 0}, {3, 4, 5, 2}, {1, 3, 1, 5}},
			[][]int{{0, 0, 0, 0}, {0, 4, 5, 0}, {0, 3, 1, 0}},
		},
		{
			"no_zeros",
			[][]int{{1, 2}, {3, 4}},
			[][]int{{1, 2}, {3, 4}},
		},
		{
			"all_zeros",
			[][]int{{0, 0}, {0, 0}},
			[][]int{{0, 0}, {0, 0}},
		},
		{
			"single_cell_zero",
			[][]int{{0}},
			[][]int{{0}},
		},
		{
			"corner_zero",
			[][]int{{0, 1, 1}, {1, 1, 1}, {1, 1, 1}},
			[][]int{{0, 0, 0}, {0, 1, 1}, {0, 1, 1}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m1 := deepCopyMatrix(tt.matrix)
			m2 := deepCopyMatrix(tt.matrix)
			setZeroes(m1)
			assert.Equal(t, tt.expected, m1, "setZeroes")
			setZeroesOptimal(m2)
			assert.Equal(t, tt.expected, m2, "setZeroesOptimal")
		})
	}
}

func deepCopyInts(nums []int) []int {
	cp := make([]int, len(nums))
	copy(cp, nums)
	return cp
}

func Test_spiralOrder(t *testing.T) {
	testCases := []struct {
		name     string
		matrix   [][]int
		expected []int
	}{
		{
			"case 1",
			[][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
			[]int{1, 2, 3, 6, 9, 8, 7, 4, 5},
		},
		{
			"case 2",
			[][]int{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}},
			[]int{1, 2, 3, 4, 8, 12, 11, 10, 9, 5, 6, 7},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := spiralOrder(tc.matrix)
			assert.Equal(t, tc.expected, result)
		})
	}
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
