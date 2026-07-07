package math

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_twoEggDrop(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		expected int
	}{
		{"n_1", 1, 1},
		{"n_2", 2, 2},
		{"n_3", 3, 2},
		{"n_4", 4, 3},
		{"n_6", 6, 3},
		{"n_7", 7, 4},
		{"n_10", 10, 4},
		{"n_14", 14, 5},
		{"n_100", 100, 14},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, twoEggDrop(tt.n))
		})
	}
}

func Test_countPoints(t *testing.T) {
	tests := []struct {
		name     string
		points   [][]int
		queries  [][]int
		expected []int
	}{
		{
			name:     "leetcode_example1",
			points:   [][]int{{1, 3}, {3, 3}, {5, 3}, {2, 2}},
			queries:  [][]int{{2, 3, 1}, {4, 3, 1}, {1, 1, 2}},
			expected: []int{3, 2, 2},
		},
		{
			name:     "leetcode_example2",
			points:   [][]int{{1, 1}, {2, 2}, {3, 3}, {4, 4}, {5, 5}},
			queries:  [][]int{{1, 2, 2}, {2, 2, 2}, {4, 3, 2}, {4, 3, 3}},
			expected: []int{2, 3, 2, 4},
		},
		{
			name:     "point_on_boundary",
			points:   [][]int{{1, 0}},
			queries:  [][]int{{0, 0, 1}},
			expected: []int{1},
		},
		{
			name:     "point_outside_circle",
			points:   [][]int{{5, 5}},
			queries:  [][]int{{0, 0, 3}},
			expected: []int{0},
		},
		{
			name:     "single_point_single_query",
			points:   [][]int{{0, 0}},
			queries:  [][]int{{0, 0, 0}},
			expected: []int{1},
		},
		{
			name:     "no_points_in_any_circle",
			points:   [][]int{{10, 10}, {20, 20}},
			queries:  [][]int{{0, 0, 5}},
			expected: []int{0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, countPoints(tt.points, tt.queries))
		})
	}
}
