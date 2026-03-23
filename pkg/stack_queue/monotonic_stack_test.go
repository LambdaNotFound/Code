package stack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_trap(t *testing.T) {
    testCases := []struct {
        name     string
        height   []int
        expected int
    }{
        {
            "case 1",
            []int{0, 1, 0, 2, 1, 0, 1, 3, 2, 1, 2, 1},
            6,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := trap(tc.height)
            assert.Equal(t, tc.expected, result)

            result = trap_slice(tc.height)
            assert.Equal(t, tc.expected, result)
        })
    }
}

func Test_largestRectangleArea(t *testing.T) {
    testCases := []struct {
        name     string
        height   []int
        expected int
    }{
        {
            "case 1",
            []int{2,1,5,6,2,3},
            10,
        },
        {
            "case 2",
            []int{2,4},
            4,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := largestRectangleArea(tc.height)
            assert.Equal(t, tc.expected, result)

            result = largestRectangleArea_slice(tc.height)
            assert.Equal(t, tc.expected, result)
        })
    }
}