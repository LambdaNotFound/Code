package math

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_canJump covers LeetCode 55 — Jump Game (greedy reachability).
//
// Branch coverage:
//   - early return false: i > reachablePosition (blocked midway)
//   - reachablePosition update: i+nums[i] > reachablePosition
//   - return true: reached end without being blocked
func Test_canJump(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected bool
	}{
		// LeetCode canonical examples.
		{name: "leetcode_example1_true", nums: []int{2, 3, 1, 1, 4}, expected: true},
		{name: "leetcode_example2_false", nums: []int{3, 2, 1, 0, 4}, expected: false},
		// Single element — already at the last index.
		{name: "single_element", nums: []int{0}, expected: true},
		// Zero at the start with two elements — can't move.
		{name: "zero_first_two_elements", nums: []int{0, 1}, expected: false},
		// Large first jump clears the whole array.
		{name: "large_first_jump", nums: []int{5, 0, 0, 0, 0}, expected: true},
		// Blocked in the middle by a zero.
		{name: "blocked_by_zero_middle", nums: []int{1, 0, 2}, expected: false},
		// Every step is exactly 1 — just barely reaches the end.
		{name: "all_ones", nums: []int{1, 1, 1, 1}, expected: true},
		// Reachable position updated multiple times before blocked.
		{name: "reachable_update_then_pass", nums: []int{2, 0, 1, 0, 1}, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, canJump(tt.nums))
		})
	}
}

func Test_canCompleteCircuit(t *testing.T) {
	testCases := []struct {
		name     string
		gas      []int
		cost     []int
		expected int
	}{
		{"case 1", []int{1, 2, 3, 4, 5}, []int{3, 4, 5, 1, 2}, 3},
		{"case 2", []int{2, 3, 4}, []int{3, 4, 3}, -1},
		{"single_station", []int{5}, []int{4}, 0},
		{"all_equal", []int{2, 2, 2}, []int{2, 2, 2}, 0},
		{"impossible", []int{1, 1, 1}, []int{2, 2, 2}, -1},
		{"start_at_zero", []int{3, 1, 1}, []int{1, 2, 2}, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := canCompleteCircuit(tc.gas, tc.cost)
			assert.Equal(t, tc.expected, result)

			result = canCompleteCircuitBruteForce(tc.gas, tc.cost)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func Test_jump(t *testing.T) {
	testCases := []struct {
		name     string
		nums     []int
		expected int
	}{
		{"case 1", []int{2, 3, 1, 1, 4}, 2},
		{"case 2", []int{2, 3, 0, 1, 4}, 2},
		{"single_element", []int{0}, 0},
		{"already_at_end", []int{5}, 0},
		{"linear_steps", []int{1, 1, 1, 1}, 3},
		{"three_jumps", []int{1, 2, 1, 1, 1}, 3},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := jump(tc.nums)
			assert.Equal(t, tc.expected, result)

			result = jumpDP(tc.nums)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func Test_leastInterval(t *testing.T) {
	testCases := []struct {
		name     string
		tasks    []byte
		n        int
		expected int
	}{
		{"case 1", []byte{'A', 'A', 'A', 'B', 'B', 'B'}, 2, 8},
		{"case 2", []byte{'A', 'C', 'A', 'B', 'D', 'B'}, 1, 6},
		{"case 3", []byte{'A', 'A', 'A', 'B', 'B', 'B'}, 3, 10},
		{"single_task", []byte{'A'}, 5, 1},
		{"no_idle_needed", []byte{'A', 'B', 'C', 'A', 'B', 'C'}, 2, 6},
		{"many_max_freq_tasks", []byte{'A', 'A', 'A', 'B', 'B', 'B', 'C', 'C', 'C'}, 2, 9},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := leastInterval(tc.tasks, tc.n)
			assert.Equal(t, tc.expected, result)
		})
	}
}
