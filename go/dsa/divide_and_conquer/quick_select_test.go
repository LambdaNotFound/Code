package divide_and_conquer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// helper: copy before calling so the caller's slice is not mutated.
func quickSelect(nums []int, k int) int {
	cp := make([]int, len(nums))
	copy(cp, nums)
	return QuickSelect(cp, k)
}

func Test_QuickSelect(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		k        int
		expected int
	}{
		// -----------------------------------------------------------------------
		// Boundary ranks (k is the 0-indexed target in sorted order)
		// -----------------------------------------------------------------------
		{
			// k=0 → the smallest element.
			name:     "k0_smallest",
			nums:     []int{3, 2, 1, 5, 6, 4},
			k:        0,
			expected: 1,
		},
		{
			// k=n-1 → the largest element.
			name:     "k_last_largest",
			nums:     []int{3, 2, 1, 5, 6, 4},
			k:        5,
			expected: 6,
		},
		{
			// Middle rank.
			name:     "middle_rank",
			nums:     []int{3, 2, 1, 5, 6, 4},
			k:        2,
			expected: 3,
		},

		// -----------------------------------------------------------------------
		// Input orderings
		// -----------------------------------------------------------------------
		{
			// Already sorted ascending — pivot is always the largest remaining
			// element, forcing the recursion all the way left.
			name:     "already_sorted_ascending",
			nums:     []int{1, 2, 3, 4, 5},
			k:        1,
			expected: 2,
		},
		{
			// Reverse sorted — pivot is always the smallest remaining element.
			name:     "reverse_sorted",
			nums:     []int{5, 4, 3, 2, 1},
			k:        3,
			expected: 4,
		},

		// -----------------------------------------------------------------------
		// Duplicates and uniform arrays
		// -----------------------------------------------------------------------
		{
			// All values the same — result is always that value regardless of k.
			name:     "all_same_values",
			nums:     []int{7, 7, 7},
			k:        1,
			expected: 7,
		},
		{
			// Scattered duplicates; ensures partitioning handles equal elements.
			name:     "with_duplicates",
			nums:     []int{3, 1, 2, 3, 4, 3},
			k:        3,
			expected: 3,
		},

		// -----------------------------------------------------------------------
		// Negative numbers and single element
		// -----------------------------------------------------------------------
		{
			name:     "negative_numbers",
			nums:     []int{-5, 3, -1, 0, -2, 8},
			k:        1,
			expected: -2,
		},
		{
			name:     "single_element",
			nums:     []int{42},
			k:        0,
			expected: 42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, quickSelect(tt.nums, tt.k))
		})
	}
}
