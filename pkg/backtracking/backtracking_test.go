package backtracking

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_permute(t *testing.T) {
    testCases := []struct {
        name     string
        nums     []int
        expected [][]int
    }{
        {
            "case 1",
            []int{1, 2, 3},
            [][]int{
                {1, 2, 3},
                {1, 3, 2},
                {2, 1, 3},
                {2, 3, 1},
                {3, 1, 2},
                {3, 2, 1},
            },
        },
        {
            "case 2",
            []int{0, 1},
            [][]int{
                {0, 1},
                {1, 0},
            },
        },
        {
            "case 3",
            []int{1},
            [][]int{
                {1},
            },
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := permute(tc.nums)
            assert.ElementsMatch(t, tc.expected, result)

            result = permuteWithSliceSpread(tc.nums)
            assert.ElementsMatch(t, tc.expected, result)
        })
    }
}

func Test_combine(t *testing.T) {
    testCases := []struct {
        name     string
        n        int
        k        int
        expected [][]int
    }{
        {
            "case 1",
            4,
            2,
            [][]int{
                {1, 2},
                {1, 3},
                {1, 4},
                {2, 3},
                {2, 4},
                {3, 4},
            },
        },
        {
            "case 2",
            1,
            1,
            [][]int{
                {1},
            },
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := combine(tc.n, tc.k)
            assert.ElementsMatch(t, tc.expected, result)
        })
    }
}

func Test_combinationSum(t *testing.T) {
    testCases := []struct {
        name       string
        candidates []int
        target     int
        expected   [][]int
    }{
        {
            "case 1",
            []int{2, 3, 6, 7},
            7,
            [][]int{
                {2, 2, 3},
                {7},
            },
        },
        {
            "case 2",
            []int{2, 3, 5},
            7,
            [][]int{
                {2, 2, 2, 2},
                {2, 3, 3},
                {3, 5},
            },
        },
        {
            "case 3",
            []int{2},
            1,
            [][]int{},
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := combinationSum(tc.candidates, tc.target)
            reflect.DeepEqual(tc.expected, result)
        })
    }
}

func Test_subsets(t *testing.T) {
    testCases := []struct {
        name     string
        nums     []int
        expected [][]int
    }{
        {
            "case 1",
            []int{1, 2, 3},
            [][]int{
                {1, 2, 3},
                {1, 2},
                {1, 3},
                {1},
                {2, 3},
                {2},
                {3},
                {},
            },
        },
        {
            "case 2",
            []int{0},
            [][]int{
                {0},
                {},
            },
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := subsets(tc.nums)
            assert.ElementsMatch(t, tc.expected, result)
        })
    }
}

func Test_letterCombinations(t *testing.T) {
    testCases := []struct {
        name     string
        digits   string
        expected []string
    }{
        {
            "case 1",
            "23",
            []string{"ad", "ae", "af", "bd", "be", "bf", "cd", "ce", "cf"},
        },
        {
            "case 2",
            "2",
            []string{"a", "b", "c"},
        },
        // Empty input triggers the early-return branch in both implementations.
        {
            "empty_digits",
            "",
            []string{},
        },
        // '7' maps to "pqrs" (4 letters) — exercises the 4-letter phoneMap entry.
        {
            "four_letter_key",
            "7",
            []string{"p", "q", "r", "s"},
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := letterCombinations(tc.digits)
            assert.Equal(t, tc.expected, result)

            result = letterCombinationsBacktrack(tc.digits)
            assert.Equal(t, tc.expected, result)
        })
    }
}

// Test_subsetsWithDup covers LeetCode 90 — subsets with duplicate elements.
//
// Branch coverage:
//   - duplicate-skip guard (i > start && nums[i] == nums[i-1])
//   - no duplicates: all elements distinct, guard never fires
//   - all-same: only one element at each recursion level survives
//   - single element: trivially two subsets — {} and {n}
func Test_subsetsWithDup(t *testing.T) {
    tests := []struct {
        name     string
        nums     []int
        expected [][]int
    }{
        {
            // LeetCode example: sort → [1,2,2]; second 2 skipped at top level.
            name:     "leetcode_example",
            nums:     []int{1, 2, 2},
            expected: [][]int{{}, {1}, {1, 2}, {1, 2, 2}, {2}, {2, 2}},
        },
        {
            // No duplicates — guard never fires; result equals plain subsets.
            name:     "no_duplicates",
            nums:     []int{1, 2, 3},
            expected: [][]int{{}, {1}, {1, 2}, {1, 2, 3}, {1, 3}, {2}, {2, 3}, {3}},
        },
        {
            // All same: [2,2,2] → only length-based subsets survive.
            name:     "all_same",
            nums:     []int{2, 2, 2},
            expected: [][]int{{}, {2}, {2, 2}, {2, 2, 2}},
        },
        {
            // Single element — no duplicates possible.
            name:     "single_element",
            nums:     []int{5},
            expected: [][]int{{}, {5}},
        },
        {
            // Two identical elements.
            name:     "two_same",
            nums:     []int{3, 3},
            expected: [][]int{{}, {3}, {3, 3}},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            assert.ElementsMatch(t, tt.expected, subsetsWithDup(tt.nums))
        })
    }
}

func Test_generateParenthesis(t *testing.T) {
    testCases := []struct {
        name     string
        n        int
        expected []string
    }{
        {
            "case 1",
            3,
            []string{
                "((()))",
                "(()())",
                "(())()",
                "()(())",
                "()()()",
            },
        },
        {
            "case 2",
            1,
            []string{
                "()",
            },
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := generateParenthesis(tc.n)
            assert.ElementsMatch(t, tc.expected, result)
        })
    }
}
