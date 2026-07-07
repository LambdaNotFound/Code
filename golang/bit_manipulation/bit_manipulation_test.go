package bit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getSum(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{name: "positive_both", a: 3, b: 5, expected: 8},
		{name: "positive_one", a: 1, b: 2, expected: 3},
		{name: "zero_b", a: 7, b: 0, expected: 7},
		{name: "zero_a", a: 0, b: 9, expected: 9},
		{name: "both_zero", a: 0, b: 0, expected: 0},
		{name: "negative_a", a: -3, b: 5, expected: 2},
		{name: "negative_b", a: 4, b: -1, expected: 3},
		{name: "both_negative", a: -2, b: -3, expected: -5},
		{name: "large_numbers", a: 1000, b: 2000, expected: 3000},
		{name: "carry_chain", a: 15, b: 1, expected: 16},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, getSum(tt.a, tt.b))
		})
	}
}

func Test_hammingWeight(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		expected int
	}{
		{name: "zero", n: 0, expected: 0},
		{name: "one", n: 1, expected: 1},
		{name: "two", n: 2, expected: 1},   // 10
		{name: "three", n: 3, expected: 2}, // 11
		{name: "seven", n: 7, expected: 3}, // 111
		{name: "eight", n: 8, expected: 1}, // 1000
		{name: "all_ones_8bit", n: 0xFF, expected: 8},
		{name: "all_ones_16bit", n: 0xFFFF, expected: 16},
		{name: "power_of_two", n: 1024, expected: 1},
		{name: "alternating_bits", n: 0b10101010, expected: 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, hammingWeight(tt.n))
		})
	}
}

func Test_countBits(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		expected []int
	}{
		{name: "n_0", n: 0, expected: []int{0}},
		{name: "n_1", n: 1, expected: []int{0, 1}},
		{name: "n_2", n: 2, expected: []int{0, 1, 1}},
		{name: "n_5", n: 5, expected: []int{0, 1, 1, 2, 1, 2}},
		{name: "n_8", n: 8, expected: []int{0, 1, 1, 2, 1, 2, 2, 3, 1}},
		{name: "n_13", n: 13, expected: []int{0, 1, 1, 2, 1, 2, 2, 3, 1, 2, 2, 3, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, countBits(tt.n))
		})
	}
}

func Test_missingNumber(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected int
	}{
		{name: "missing_2", nums: []int{3, 0, 1}, expected: 2},
		{name: "missing_2_v2", nums: []int{0, 1}, expected: 2},
		{name: "missing_8", nums: []int{9, 6, 4, 2, 3, 5, 7, 0, 1}, expected: 8},
		{name: "missing_0", nums: []int{1}, expected: 0},
		{name: "missing_1", nums: []int{0}, expected: 1},
		{name: "missing_last", nums: []int{0, 1, 2, 3, 4}, expected: 5},
		{name: "sequential_missing_middle", nums: []int{0, 1, 3, 4, 5}, expected: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, missingNumber(tt.nums))
		})
	}
}

func Test_reverseBits(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		expected int
	}{
		// 0...01 reversed → 10...0 = 2^31 = 2147483648
		{name: "one", n: 1, expected: 1 << 31},
		// 0 reversed → 0
		{name: "zero", n: 0, expected: 0},
		// 0b00000010110111011110000101001110 = 43261596
		// reversed: 0b01110010100001111011101101000000 = 964176192
		{name: "leetcode_example1", n: 43261596, expected: 964176192},
		// all 32 bits set → all 32 bits set (palindrome in bits)
		// But we use int (64-bit on Go), so 0xFFFFFFFF as input
		// reversing lower 32 bits of 0xFFFFFFFF → 0xFFFFFFFF = 4294967295
		{name: "all_ones_32bit", n: 0xFFFFFFFF, expected: 0xFFFFFFFF},
		// bit 31 set → bit 0 set = 1
		{name: "msb_set", n: 1 << 31, expected: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, reverseBits(tt.n))
		})
	}
}
