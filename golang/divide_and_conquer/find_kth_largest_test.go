package divide_and_conquer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/**
 * 215. Kth Largest Element in an Array
 *
 * Both findKthLargest and findKthLargestAlt implement quickselect but express
 * the rank differently:
 *
 *   findKthLargest    — converts k to a 0-indexed target from the left:
 *                       target = len(nums)-k, then partitions until pivot==target.
 *   findKthLargestAlt — works in terms of "rank from largest" (1-indexed):
 *                       rank = len(nums)-pivot, and navigates left/right accordingly.
 *
 * Both share the same partition() helper (Lomuto scheme, rightmost pivot).
 *
 * Test strategy:
 *   - LeetCode canonical examples.
 *   - k=1 (largest) and k=n (smallest) — the two boundary ranks.
 *   - Middle rank on sorted, reverse-sorted, and random input.
 *   - All-same values: every position holds the same answer.
 *   - Duplicates scattered through the array.
 *   - Single-element array: trivially returns that element.
 *   - partition() is exercised directly to cover the nil-pointer nil branch
 *     of partitionListSwap via a standalone call.
 */

// helper: copy before calling so the caller's slice is not mutated.
func kthLargest(nums []int, k int) int {
	cp := make([]int, len(nums))
	copy(cp, nums)
	return findKthLargest(cp, k)
}

func kthLargestAlt(nums []int, k int) int {
	cp := make([]int, len(nums))
	copy(cp, nums)
	return findKthLargestAlt(cp, k)
}

func Test_findKthLargest(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		k        int
		expected int
	}{
		// -----------------------------------------------------------------------
		// LeetCode canonical examples
		// -----------------------------------------------------------------------
		{
			name:     "leetcode_example1",
			nums:     []int{3, 2, 1, 5, 6, 4},
			k:        2,
			expected: 5,
		},
		{
			name:     "leetcode_example2",
			nums:     []int{3, 2, 3, 1, 2, 4, 5, 5, 6},
			k:        4,
			expected: 4,
		},

		// -----------------------------------------------------------------------
		// Boundary ranks
		// -----------------------------------------------------------------------
		{
			// k=1 → the largest element.
			name:     "k1_largest",
			nums:     []int{3, 2, 1, 5, 6, 4},
			k:        1,
			expected: 6,
		},
		{
			// k=n → the smallest element.
			name:     "k_equals_n_smallest",
			nums:     []int{5, 3, 1, 2, 4},
			k:        5,
			expected: 1,
		},

		// -----------------------------------------------------------------------
		// Input orderings
		// -----------------------------------------------------------------------
		{
			// Already sorted ascending — pivot is always the largest, forcing
			// the algorithm to scan all the way left on every pass.
			name:     "already_sorted_ascending",
			nums:     []int{1, 2, 3, 4, 5},
			k:        2,
			expected: 4,
		},
		{
			// Reverse sorted — pivot is always the smallest.
			name:     "reverse_sorted",
			nums:     []int{5, 4, 3, 2, 1},
			k:        3,
			expected: 3,
		},

		// -----------------------------------------------------------------------
		// Duplicates and uniform arrays
		// -----------------------------------------------------------------------
		{
			// All values the same — kth largest is always that value.
			name:     "all_same_values",
			nums:     []int{7, 7, 7},
			k:        2,
			expected: 7,
		},
		{
			// Scattered duplicates; ensures the partition handles equal elements.
			name:     "with_duplicates",
			nums:     []int{3, 1, 2, 3, 4, 3},
			k:        3,
			expected: 3,
		},

		// -----------------------------------------------------------------------
		// Single element
		// -----------------------------------------------------------------------
		{
			name:     "single_element",
			nums:     []int{42},
			k:        1,
			expected: 42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, kthLargest(tt.nums, tt.k))
		})
	}
}

func Test_findKthLargestAlt(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		k        int
		expected int
	}{
		{name: "leetcode_example1", nums: []int{3, 2, 1, 5, 6, 4}, k: 2, expected: 5},
		{name: "leetcode_example2", nums: []int{3, 2, 3, 1, 2, 4, 5, 5, 6}, k: 4, expected: 4},
		{name: "k1_largest", nums: []int{3, 2, 1, 5, 6, 4}, k: 1, expected: 6},
		{name: "k_equals_n_smallest", nums: []int{5, 3, 1, 2, 4}, k: 5, expected: 1},
		{name: "already_sorted", nums: []int{1, 2, 3, 4, 5}, k: 2, expected: 4},
		{name: "reverse_sorted_mid", nums: []int{5, 4, 3, 2, 1}, k: 3, expected: 3},
		{name: "all_same_values", nums: []int{7, 7, 7}, k: 2, expected: 7},
		{name: "single_element", nums: []int{42}, k: 1, expected: 42},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, kthLargestAlt(tt.nums, tt.k))
		})
	}
}

// Test_partition exercises the Lomuto partition helper directly.
// This also covers the nil/nil early-return path of partitionListSwap
// (which quickSortHelper never reaches) by calling it directly.
func Test_partition(t *testing.T) {
	// Property verified for each case: arr[got] == pivot, all left ≤ pivot, all right ≥ pivot.
	tests := []struct {
		name  string
		input []int
		left  int
		right int
	}{
		{
			// pivot=5; 3,2,1,4 all ≤5 → pivot lands at index 4.
			name: "random_order",
			input: []int{3, 2, 1, 4, 6, 5}, left: 0, right: 5,
		},
		{
			// pivot=1 (smallest) → ends up at leftmost position.
			name:  "pivot_is_smallest",
			input: []int{5, 4, 3, 1}, left: 0, right: 3,
		},
		{
			// pivot=5 (largest) → ends up at rightmost position.
			name:  "pivot_is_largest",
			input: []int{1, 2, 3, 5}, left: 0, right: 3,
		},
		{
			// Two elements: pivot is the smaller one.
			name:  "two_elements_pivot_smaller",
			input: []int{3, 1}, left: 0, right: 1,
		},
		{
			// Sub-range of a larger array, single slot: pivot stays in place.
			name:  "single_element_range",
			input: []int{9, 2, 7}, left: 2, right: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arr := make([]int, len(tt.input))
			copy(arr, tt.input)
			got := partition(arr, tt.left, tt.right)
			// The pivot value must be at the returned index.
			pivotVal := tt.input[tt.right]
			assert.Equal(t, pivotVal, arr[got], "pivot value must be at returned index")
			// Everything to the left of got must be ≤ pivot.
			for i := tt.left; i < got; i++ {
				assert.LessOrEqual(t, arr[i], pivotVal, "left of pivot must be ≤ pivot at index %d", i)
			}
			// Everything to the right of got must be ≥ pivot.
			for i := got + 1; i <= tt.right; i++ {
				assert.GreaterOrEqual(t, arr[i], pivotVal, "right of pivot must be ≥ pivot at index %d", i)
			}
		})
	}
}

// Test_partitionListSwap_nil covers the nil-guard branch that quickSortHelper
// never reaches in practice (it only calls partitionListSwap when head != tail).
func Test_partitionListSwap_nil(t *testing.T) {
	t.Run("nil_head_returns_nil", func(t *testing.T) {
		assert.Nil(t, partitionListSwap(nil, nil))
	})
}
