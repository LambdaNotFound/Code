package dynamic_programming

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_canJump(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected bool
	}{
		{name: "example1", nums: []int{2, 3, 1, 1, 4}, expected: true},
		{name: "example2", nums: []int{3, 2, 1, 0, 4}, expected: false},
		{name: "single_element", nums: []int{0}, expected: true},
		{name: "two_reachable", nums: []int{1, 0}, expected: true},
		{name: "two_unreachable", nums: []int{0, 1}, expected: false},
		{name: "large_first_jump", nums: []int{5, 0, 0, 0, 0}, expected: true},
		{name: "all_ones", nums: []int{1, 1, 1, 1, 1}, expected: true},
		{name: "zero_stuck_in_middle", nums: []int{1, 0, 1, 0}, expected: false},
		{name: "just_barely_reachable", nums: []int{2, 0, 0}, expected: true},
		{name: "all_zeros_single", nums: []int{0}, expected: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, canJump(tt.nums))
		})
	}
}
