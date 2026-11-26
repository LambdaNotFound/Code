package array

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

// deepCopyInts is copied from matrix/matrix_test.go, which kept the original when 189 moved here.
func deepCopyInts(nums []int) []int {
	cp := make([]int, len(nums))
	copy(cp, nums)
	return cp
}
