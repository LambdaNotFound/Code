package topologicalsort

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_canFinish(t *testing.T) {
	testCases := []struct {
		name          string
		numCourses    int
		prerequisites [][]int
		expected      bool
	}{
		{"case 1", 2, [][]int{{1, 0}}, true},
		{"case 2", 2, [][]int{{1, 0}, {0, 1}}, false},
		{name: "no_prereqs", numCourses: 3, prerequisites: [][]int{}, expected: true},
		{name: "single_course", numCourses: 1, prerequisites: [][]int{}, expected: true},
		{name: "linear_chain", numCourses: 4, prerequisites: [][]int{{1, 0}, {2, 1}, {3, 2}}, expected: true},
		{name: "two_separate_cycles", numCourses: 4, prerequisites: [][]int{{0, 1}, {1, 0}, {2, 3}, {3, 2}}, expected: false},
		{name: "self_loop", numCourses: 2, prerequisites: [][]int{{0, 0}}, expected: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := canFinish(tc.numCourses, tc.prerequisites)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func Test_findOrder(t *testing.T) {
	testCases := []struct {
		name          string
		numCourses    int
		prerequisites [][]int
		expected      []int
	}{
		{
			"case 1",
			2,
			[][]int{
				{1, 0},
			},
			[]int{0, 1},
		},
		{
			"case 2",
			4,
			[][]int{
				{1, 0},
				{2, 0},
				{3, 1},
				{3, 2},
			},
			[]int{0, 2, 1, 3},
		},
		{
			"case 3",
			1,
			[][]int{},
			[]int{0},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := findOrder(tc.numCourses, tc.prerequisites)
			assert.ElementsMatch(t, tc.expected, result)
		})
	}
}

func Test_findMinHeightTrees(t *testing.T) {
	testCases := []struct {
		name     string
		n        int
		edges    [][]int
		expected []int
	}{
		{
			"case 1",
			4,
			[][]int{
				{1, 0},
				{1, 2},
				{1, 3},
			},
			[]int{1},
		},
		{
			"case 2",
			6,
			[][]int{
				{3, 0},
				{3, 1},
				{3, 2},
				{3, 4},
				{5, 4},
			},
			[]int{3, 4},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := findMinHeightTrees(tc.n, tc.edges)
			assert.ElementsMatch(t, tc.expected, result)
		})
	}
}

func Test_longestIncreasingPath(t *testing.T) {
	testCases := []struct {
		name     string
		matrix   [][]int
		expected int
	}{
		{
			"case 1",
			[][]int{
				{9, 9, 4},
				{6, 6, 8},
				{2, 1, 1},
			},
			4,
		},
		{
			"case 2",
			[][]int{
				{3, 4, 5},
				{3, 2, 6},
				{2, 2, 1},
			},
			4,
		},
		{
			"case 3",
			[][]int{
				{1},
			},
			1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := longestIncreasingPath(tc.matrix)
			assert.Equal(t, tc.expected, result)

			result = longestIncreasingPathMemoization(tc.matrix)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func Test_findMinHeightTreesTwoPassBFS(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		edges    [][]int
		expected []int
	}{
		{
			name:     "single_node",
			n:        1,
			edges:    [][]int{},
			expected: []int{0},
		},
		{
			name:     "two_nodes",
			n:        2,
			edges:    [][]int{{0, 1}},
			expected: []int{0, 1},
		},
		{
			name:     "star_graph_one_center",
			n:        4,
			edges:    [][]int{{1, 0}, {1, 2}, {1, 3}},
			expected: []int{1},
		},
		{
			name:     "path_graph_two_centers",
			n:        6,
			edges:    [][]int{{3, 0}, {3, 1}, {3, 2}, {3, 4}, {5, 4}},
			expected: []int{3, 4},
		},
		{
			name:     "chain_of_five_one_center",
			n:        5,
			edges:    [][]int{{0, 1}, {1, 2}, {2, 3}, {3, 4}},
			expected: []int{2},
		},
		{
			name:     "chain_of_four_two_centers",
			n:        4,
			edges:    [][]int{{0, 1}, {1, 2}, {2, 3}},
			expected: []int{1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findMinHeightTreesTwoPassBFS(tt.n, tt.edges)
			assert.ElementsMatch(t, tt.expected, got)
		})
	}
}
