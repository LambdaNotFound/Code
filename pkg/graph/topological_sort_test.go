package graph

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
        {
            "case 1",
            2,
            [][]int{
                {1, 0},
            },
            true,
        },
        {
            "case 2",
            2,
            [][]int{
                {1, 0},
                {0, 1},
            },
            false,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := canFinish(tc.numCourses, tc.prerequisites)
            assert.Equal(t, tc.expected, result)
        })
    }
}

func Test_scheduleCourse(t *testing.T) {
    tests := []struct {
        name     string
        courses  [][]int
        expected int
    }{
        {
            name:     "Example 1 (LeetCode)",
            courses:  [][]int{{100, 200}, {200, 1300}, {1000, 1250}, {2000, 3200}},
            expected: 3,
        },
        {
            name:     "All fit perfectly",
            courses:  [][]int{{1, 2}, {2, 3}, {3, 4}},
            expected: 3,
        },
        {
            name:     "None fit (all deadlines too short)",
            courses:  [][]int{{3, 2}, {4, 3}},
            expected: 0,
        },
        {
            name:     "Some fit (pick optimal subset)",
            courses:  [][]int{{5, 5}, {4, 6}, {2, 6}},
            expected: 2, // Can take {4,6} and {2,6} or {5,5} and {2,6}
        },
        {
            name:     "Tight deadlines require optimal replacement",
            courses:  [][]int{{3, 2}, {4, 3}, {2, 4}, {1, 100}},
            expected: 3,
        },
        {
            name:     "Empty input",
            courses:  [][]int{},
            expected: 0,
        },
        {
            name:     "Single course fits",
            courses:  [][]int{{5, 10}},
            expected: 1,
        },
        {
            name:     "Single course doesn't fit",
            courses:  [][]int{{10, 5}},
            expected: 0,
        },
        {
            name:     "Courses with same deadlines but different durations",
            courses:  [][]int{{1, 5}, {2, 5}, {3, 5}, {4, 5}},
            expected: 2, // Best pick shortest ones: {1,5} + {2,5}
        },
        {
            name:     "Large course late in schedule",
            courses:  [][]int{{2, 2}, {3, 4}, {10, 100}},
            expected: 3, // Enough total time for all
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            got := scheduleCourseHeap(tc.courses)
            assert.Equal(t, tc.expected, got)

            got = scheduleCourseHeap(tc.courses)
            assert.Equal(t, tc.expected, got)
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
