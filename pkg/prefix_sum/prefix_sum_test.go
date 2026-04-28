package prefixsum

import (
	"math"
	"testing"

	. "gocode/types"

	"github.com/stretchr/testify/assert"
)

// buildTreePS constructs a binary tree from level-order values.
// math.MinInt is used as a sentinel for nil nodes.
func buildTreePS(vals []int) *TreeNode {
	if len(vals) == 0 || vals[0] == math.MinInt {
		return nil
	}
	root := &TreeNode{Val: vals[0]}
	queue := []*TreeNode{root}
	i := 1
	for len(queue) > 0 && i < len(vals) {
		node := queue[0]
		queue = queue[1:]
		if i < len(vals) && vals[i] != math.MinInt {
			node.Left = &TreeNode{Val: vals[i]}
			queue = append(queue, node.Left)
		}
		i++
		if i < len(vals) && vals[i] != math.MinInt {
			node.Right = &TreeNode{Val: vals[i]}
			queue = append(queue, node.Right)
		}
		i++
	}
	return root
}

func Test_productExceptSelf(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected []int
	}{
		{name: "four_elements", nums: []int{1, 2, 3, 4}, expected: []int{24, 12, 8, 6}},
		{name: "with_zero", nums: []int{-1, 1, 0, -3, 3}, expected: []int{0, 0, 9, 0, 0}},
		{name: "single_element", nums: []int{5}, expected: []int{1}},
		{name: "two_elements", nums: []int{2, 3}, expected: []int{3, 2}},
		{name: "two_zeros", nums: []int{0, 0, 1}, expected: []int{0, 0, 0}},
		{name: "negatives", nums: []int{-2, -3, -4}, expected: []int{12, 8, 6}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, productExceptSelf(tt.nums), "productExceptSelf")
		})
	}
}

func Test_findMaxLength(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected int
	}{
		{name: "leetcode_example1", nums: []int{0, 1}, expected: 2},
		{name: "leetcode_example2", nums: []int{0, 1, 0}, expected: 2},
		{name: "all_equal_length", nums: []int{0, 0, 0, 1, 1, 1}, expected: 6},
		{name: "partial_equal", nums: []int{0, 1, 1, 0, 1, 1, 1}, expected: 4},
		{name: "all_zeros", nums: []int{0, 0, 0}, expected: 0},
		{name: "all_ones", nums: []int{1, 1, 1}, expected: 0},
		{name: "single_zero", nums: []int{0}, expected: 0},
		{name: "empty", nums: []int{}, expected: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, findMaxLength(tt.nums), "findMaxLength")
			assert.Equal(t, tt.expected, findMaxLengthPrefixSum(tt.nums), "findMaxLengthPrefixSum")
		})
	}
}

func Test_subarraySum(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		k        int
		expected int
	}{
		{name: "leetcode_example1", nums: []int{1, 1, 1}, k: 2, expected: 2},
		{name: "leetcode_example2", nums: []int{1, 2, 3}, k: 3, expected: 2},
		{name: "negative_numbers", nums: []int{1, -1, 0}, k: 0, expected: 3},
		{name: "four_subarrays", nums: []int{3, 4, 7, 2, -3, 1, 4, 2}, k: 7, expected: 4},
		{name: "single_element_match", nums: []int{5}, k: 5, expected: 1},
		{name: "single_element_no_match", nums: []int{5}, k: 3, expected: 0},
		{name: "all_zeros_k0", nums: []int{0, 0, 0}, k: 0, expected: 6},
		{name: "no_match", nums: []int{1, 2, 3}, k: 10, expected: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, subarraySum(tt.nums, tt.k), "subarraySum")
			assert.Equal(t, tt.expected, subarraySumWithHashmap(tt.nums, tt.k), "subarraySumWithHashmap")
		})
	}
}

func Test_numSubarraysWithSum(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		goal     int
		expected int
	}{
		{name: "leetcode_example1", nums: []int{1, 0, 1, 0, 1}, goal: 2, expected: 4},
		{name: "leetcode_example2", nums: []int{0, 0, 0, 0, 0}, goal: 0, expected: 15},
		{name: "single_one", nums: []int{1, 0, 1}, goal: 1, expected: 4},
		{name: "goal_zero_with_ones", nums: []int{1, 1, 1}, goal: 0, expected: 0},
		{name: "goal_equals_all", nums: []int{1, 1, 1}, goal: 3, expected: 1},
		{name: "single_element_match", nums: []int{1}, goal: 1, expected: 1},
		{name: "single_element_no_match", nums: []int{0}, goal: 1, expected: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, numSubarraysWithSum(tt.nums, tt.goal), "hashmap")
			assert.Equal(t, tt.expected, numSubarraysWithSumPrefixSum(tt.nums, tt.goal), "prefix sum O(n²)")
		})
	}
}

/**
 * 523. Continuous Subarray Sum
 *
 * Returns true when there is a subarray of length ≥ 2 whose sum is a multiple
 * of k. Uses a remainder map: equal remainders at indices i and j (j-i ≥ 2)
 * mean the subarray between them has sum divisible by k.
 *
 * Test strategy:
 *   - LeetCode canonical examples (true / true / false).
 *   - Minimum-length subarray (exactly 2 elements).
 *   - Single element — never valid regardless of value.
 *   - Two zeros — sum is 0, which is a multiple of any k.
 *   - length_1_gap_blocked: remainder repeats but the two occurrences are only
 *     one index apart (i-firstIdx == 1 < 2), so the function must return false.
 */
func Test_checkSubarraySum(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		k        int
		expected bool
	}{
		{
			// [2,4] sums to 6 = 1×6.
			name: "leetcode_example1",
			nums: []int{23, 2, 4, 6, 7}, k: 6, expected: true,
		},
		{
			// Entire array sums to 42 = 7×6.
			name: "leetcode_example2",
			nums: []int{23, 2, 6, 4, 7}, k: 6, expected: true,
		},
		{
			// No subarray of length ≥ 2 sums to a multiple of 13.
			name: "leetcode_example3_false",
			nums: []int{23, 2, 6, 4, 7}, k: 13, expected: false,
		},
		{
			// [0,0] → sum=0, multiple of any k; length=2 ≥ 2.
			name: "two_zeros_multiple_of_any_k",
			nums: []int{0, 0}, k: 7, expected: true,
		},
		{
			// Exactly two elements summing to a multiple: 3+3=6=2×3.
			name: "exact_two_elements_multiple",
			nums: []int{3, 3}, k: 3, expected: true,
		},
		{
			// Single element can never form a length-2 subarray.
			name: "single_element_never_valid",
			nums: []int{6}, k: 6, expected: false,
		},
		{
			// 1+2=3, not a multiple of 5.
			name: "two_elements_not_multiple",
			nums: []int{1, 2}, k: 5, expected: false,
		},
		{
			// Remainder 1 appears at index 1 and index 2 (gap = 1 < 2).
			// The guard i-firstIdx >= 2 must block the early return.
			name: "length_1_gap_blocked",
			nums: []int{5, 1, 5}, k: 5, expected: false,
		},
		{
			// [1,2] sums to 3 = 1×3; length=2 ≥ 2.
			name: "sum_equals_k",
			nums: []int{1, 2, 3}, k: 3, expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, checkSubarraySum(tt.nums, tt.k))
		})
	}
}

func Test_pathSum(t *testing.T) {
	ni := math.MinInt // nil sentinel

	tests := []struct {
		name      string
		vals      []int
		targetSum int
		expected  int
	}{
		{
			name:      "nil_root",
			vals:      []int{},
			targetSum: 5,
			expected:  0,
		},
		{
			name:      "single_node_match",
			vals:      []int{5},
			targetSum: 5,
			expected:  1,
		},
		{
			name:      "single_node_no_match",
			vals:      []int{1},
			targetSum: 0,
			expected:  0,
		},
		{
			// Tree:   1
			//        / \
			//       2   3
			// Paths summing to 3: [1,2] and [3]
			name:      "simple_tree_two_paths",
			vals:      []int{1, 2, 3},
			targetSum: 3,
			expected:  2,
		},
		{
			// LeetCode example: root=[10,5,-3,3,2,null,11,3,-2,null,1], targetSum=8
			// Paths: 5→3, 5→2→1, -3→11
			name:      "leetcode_example",
			vals:      []int{10, 5, -3, 3, 2, ni, 11, 3, -2, ni, 1},
			targetSum: 8,
			expected:  3,
		},
		{
			name:      "no_matching_paths",
			vals:      []int{1, 2, 3},
			targetSum: 100,
			expected:  0,
		},
		{
			// Tree: 0
			//      / \
			//     1  -1
			// Paths summing to 0: [0,-1... wait no.
			// Paths: [0]=0 ✓, [1]=1, [-1]=-1, [0,1]=1, [0,-1]=-1
			// Only [0] sums to 0.
			name:      "root_matches_zero",
			vals:      []int{0, 1, -1},
			targetSum: 0,
			expected:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := buildTreePS(tt.vals)
			result := pathSum(root, tt.targetSum)
			assert.Equal(t, tt.expected, result)
		})
	}
}
