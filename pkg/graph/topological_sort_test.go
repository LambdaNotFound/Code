package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
