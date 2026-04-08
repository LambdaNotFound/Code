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

func Test_romanToInt(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		expected int
	}{
		{"III", "III", 3},
		{"LVIII", "LVIII", 58},
		{"MCMXCIV", "MCMXCIV", 1994},
		{"IV", "IV", 4},
		{"IX", "IX", 9},
		{"XL", "XL", 40},
		{"XC", "XC", 90},
		{"CD", "CD", 400},
		{"CM", "CM", 900},
		{"single_M", "M", 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, romanToInt(tt.s))
		})
	}
}

func Test_firstMissingPositive(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected int
	}{
		{"leetcode_1", []int{1, 2, 0}, 3},
		{"leetcode_2", []int{3, 4, -1, 1}, 2},
		{"leetcode_3", []int{7, 8, 9, 11, 12}, 1},
		{"single_one", []int{1}, 2},
		{"single_two", []int{2}, 1},
		{"consecutive", []int{1, 2, 3, 4, 5}, 6},
		{"all_negative", []int{-1, -2, -3}, 1},
		{"with_duplicates", []int{1, 1, 2, 2}, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nums := append([]int(nil), tt.nums...)
			assert.Equal(t, tt.expected, firstMissingPositive(nums))
		})
	}
}

func Test_removeDuplicatesFromSortedArray(t *testing.T) {
	tests := []struct {
		name         string
		nums         []int
		expectedLen  int
		expectedNums []int // first expectedLen elements
	}{
		{
			name:         "no_duplicates",
			nums:         []int{1, 2, 3, 4},
			expectedLen:  4,
			expectedNums: []int{1, 2, 3, 4},
		},
		{
			name:         "all_same",
			nums:         []int{2, 2, 2, 2},
			expectedLen:  2,
			expectedNums: []int{2, 2},
		},
		{
			name:         "each_appears_twice",
			nums:         []int{1, 1, 2, 2, 3, 3},
			expectedLen:  6,
			expectedNums: []int{1, 1, 2, 2, 3, 3},
		},
		{
			name:         "three_duplicates_trimmed",
			nums:         []int{1, 1, 1, 2, 2, 3},
			expectedLen:  5,
			expectedNums: []int{1, 1, 2, 2, 3},
		},
		{
			name:         "single_element",
			nums:         []int{7},
			expectedLen:  1,
			expectedNums: []int{7},
		},
		{
			name:         "two_same",
			nums:         []int{5, 5},
			expectedLen:  2,
			expectedNums: []int{5, 5},
		},
		{
			name:         "three_same_trimmed_to_two",
			nums:         []int{5, 5, 5},
			expectedLen:  2,
			expectedNums: []int{5, 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nums := append([]int(nil), tt.nums...)
			got := removeDuplicatesFromSortedArray(nums)
			assert.Equal(t, tt.expectedLen, got)
			assert.Equal(t, tt.expectedNums, nums[:got])
		})
	}
}
