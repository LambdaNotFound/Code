package prefixsum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_upperBound(t *testing.T) {
	tests := []struct {
		name     string
		array    []int
		target   int
		expected int
	}{
		// upperBound returns the first index i where target < array[i]
		{name: "first_bucket", array: []int{0, 1, 3, 6, 10}, target: 0, expected: 1},
		{name: "second_bucket", array: []int{0, 1, 3, 6, 10}, target: 1, expected: 2},
		{name: "third_bucket", array: []int{0, 1, 3, 6, 10}, target: 5, expected: 3},
		{name: "last_bucket", array: []int{0, 1, 3, 6, 10}, target: 9, expected: 4},
		{name: "two_elements", array: []int{0, 3}, target: 0, expected: 1},
		{name: "equal_weights", array: []int{0, 2, 4, 6}, target: 3, expected: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := upperBound(tt.array, tt.target)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test_lowerBound covers the lowerBound helper (first index i where target ≤ array[i]).
// This is the "lower bound" in the classic binary-search sense: the first element
// that is ≥ target. It differs from upperBound (first element strictly > target)
// by one position when target equals an element in the array.
func Test_lowerBound(t *testing.T) {
	tests := []struct {
		name     string
		array    []int
		target   int
		expected int
	}{
		// target exactly equals the first element → index 0.
		{name: "target_equals_first",   array: []int{0, 1, 3, 6, 10}, target: 0,  expected: 0},
		// target exactly equals an interior element → that element's index.
		{name: "target_equals_second",  array: []int{0, 1, 3, 6, 10}, target: 1,  expected: 1},
		{name: "target_equals_third",   array: []int{0, 1, 3, 6, 10}, target: 3,  expected: 2},
		// target falls between two elements → index of the first element > target.
		{name: "target_between_1_and_3", array: []int{0, 1, 3, 6, 10}, target: 2, expected: 2},
		{name: "target_between_3_and_6", array: []int{0, 1, 3, 6, 10}, target: 5, expected: 3},
		// target equals the last element.
		{name: "target_equals_last",    array: []int{0, 1, 3, 6, 10}, target: 10, expected: 4},
		// target beyond every element → returns len(array).
		{name: "target_beyond_last",    array: []int{0, 1, 3, 6, 10}, target: 11, expected: 5},
		// two-element array.
		{name: "two_elements_between",  array: []int{0, 3},            target: 1,  expected: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, lowerBound(tt.array, tt.target))
		})
	}
}

func Test_PickIndex(t *testing.T) {
	t.Run("single_weight_always_zero", func(t *testing.T) {
		s := Constructor([]int{10})
		for i := 0; i < 20; i++ {
			assert.Equal(t, 0, s.PickIndex())
		}
	})

	t.Run("returns_valid_index", func(t *testing.T) {
		weights := []int{1, 2, 3, 4}
		s := Constructor(weights)
		for i := 0; i < 100; i++ {
			idx := s.PickIndex()
			assert.True(t, idx >= 0 && idx < len(weights),
				"PickIndex %d out of range [0, %d)", idx, len(weights))
		}
	})

	t.Run("two_equal_weights_both_reachable", func(t *testing.T) {
		s := Constructor([]int{1, 1})
		seen := map[int]bool{}
		for i := 0; i < 200; i++ {
			seen[s.PickIndex()] = true
		}
		assert.True(t, seen[0], "index 0 should be reachable")
		assert.True(t, seen[1], "index 1 should be reachable")
	})
}
