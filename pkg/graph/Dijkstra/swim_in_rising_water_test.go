package dijkstra

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_swimInWater(t *testing.T) {
	testCases := []struct {
		name     string
		grid     [][]int
		expected int
	}{
		{
			name:     "leetcode example 1 — 2x2",
			grid:     [][]int{{0, 2}, {1, 3}},
			expected: 3,
		},
		{
			name: "leetcode example 2 — 5x5",
			grid: [][]int{
				{0, 1, 2, 3, 4},
				{24, 23, 22, 21, 5},
				{12, 13, 14, 15, 16},
				{11, 17, 18, 19, 20},
				{10, 9, 8, 7, 6},
			},
			expected: 16,
		},
		{
			name:     "1x1 grid — already at destination",
			grid:     [][]int{{0}},
			expected: 0,
		},
		{
			name:     "1x1 grid — non-zero start",
			grid:     [][]int{{7}},
			expected: 7,
		},
		{
			name:     "straight path — bottleneck is max along path",
			grid:     [][]int{{0, 3}, {1, 2}},
			expected: 2, // path (0,0)→(1,0)→(1,1): max(0,1,2)=2
		},
		{
			name: "must go around high peak",
			grid: [][]int{
				{0, 9, 1},
				{8, 9, 2},
				{7, 6, 3},
			},
			expected: 8, // path (0,0)→(1,0)→(2,0)→(2,1)→(2,2): max(0,8,7,6,3)=8
		},
		{
			name: "diagonal obstacle forces high t",
			grid: [][]int{
				{0, 1, 2},
				{3, 8, 4},
				{7, 6, 5},
			},
			expected: 5, // path (0,0)→(0,1)→(0,2)→(1,2)→(2,2): max(0,1,2,4,5)=5
		},
		{
			name: "already sorted — answer equals bottom-right",
			grid: [][]int{
				{0, 1, 2},
				{3, 4, 5},
				{6, 7, 8},
			},
			expected: 8,
		},
		{
			name: "reverse sorted — greedy picks lowest elevation path",
			grid: [][]int{
				{8, 7, 6},
				{5, 4, 3},
				{2, 1, 0},
			},
			expected: 8, // must start at 8; no way to avoid it
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, swimInWater(deepCopy(tc.grid)))
		})
	}
}
