package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func deepCopy2D(src [][]int) [][]int {
    dst := make([][]int, len(src))
    for i := range src {
        dst[i] = make([]int, len(src[i]))
        copy(dst[i], src[i])
    }
    return dst
}

func Test_orangesRotting(t *testing.T) {
    testCases := []struct {
        name     string
        grid     [][]int
        expected int
    }{
        {
            "case 1",
            [][]int{{2, 1, 1}, {1, 1, 0}, {0, 1, 1}},
            4,
        },
        {
            "case 2",
            [][]int{{2, 1, 1}, {0, 1, 1}, {1, 0, 1}},
            -1,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            grid := deepCopy2D(tc.grid)
            result := orangesRotting(grid)
            assert.Equal(t, tc.expected, result)

            grid = deepCopy2D(tc.grid)
            result = orangesRotting_slice(grid)
            assert.Equal(t, tc.expected, result)
        })
    }
}
