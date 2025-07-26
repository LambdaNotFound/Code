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
