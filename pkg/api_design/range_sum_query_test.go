package apidesign

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NumArray(t *testing.T) {
	testCases := []struct {
		name     string
		nums     []int
		queries  []struct{ left, right, expected int }
	}{
		{
			name: "leetcode_example",
			nums: []int{-2, 0, 3, -5, 2, -1},
			queries: []struct{ left, right, expected int }{
				{0, 2, 1},  // -2+0+3 = 1
				{2, 5, -1}, // 3-5+2-1 = -1
				{0, 5, -3}, // -2+0+3-5+2-1 = -3
			},
		},
		{
			name: "single_element",
			nums: []int{7},
			queries: []struct{ left, right, expected int }{
				{0, 0, 7},
			},
		},
		{
			name: "all_zeros",
			nums: []int{0, 0, 0},
			queries: []struct{ left, right, expected int }{
				{0, 2, 0},
				{1, 1, 0},
			},
		},
		{
			name: "full_range",
			nums: []int{1, 2, 3, 4, 5},
			queries: []struct{ left, right, expected int }{
				{0, 4, 15},
				{0, 0, 1},
				{4, 4, 5},
				{1, 3, 9},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			na := ConstructorNumArray(tc.nums)
			for _, q := range tc.queries {
				assert.Equal(t, q.expected, na.SumRange(q.left, q.right))
			}
		})
	}
}
