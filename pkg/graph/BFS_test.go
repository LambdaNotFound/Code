package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
            result := orangesRotting_slice(tc.grid)
            assert.Equal(t, tc.expected, result)

            result = orangesRotting(tc.grid)
            assert.Equal(t, tc.expected, result)
        })
    }
}
