package math

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
