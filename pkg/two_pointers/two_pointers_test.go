package two_pointers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_maxArea(t *testing.T) {
    testCases := []struct {
        name     string
        height   []int
        expected int
    }{
        {
            "case 1",
            []int{1, 8, 6, 2, 5, 4, 8, 3, 7},
            49,
        },
        {
            "case 2",
            []int{1, 1},
            1,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := maxArea(tc.height)
            assert.Equal(t, tc.expected, result)
        })
    }
}

func Test_twoSum(t *testing.T) {
    testCases := []struct {
        name     string
        numbers  []int
        target   int
        expected []int
    }{
        {
            "case 1",
            []int{2, 7, 11, 15},
            9,
            []int{1, 2},
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := twoSum(tc.numbers, tc.target)
            assert.Equal(t, tc.expected, result)
        })
    }
}

func Test_threeSum(t *testing.T) {
    testCases := []struct {
        name     string
        numbers  []int
        expected [][]int
    }{
        {
            "case 1",
            []int{-1, 0, 1, 2, -1, -4},
            [][]int{
                {-1, -1, 2}, {-1, 0, 1},
            },
        },
        {
            "case 2",
            []int{0, 1, 1},
            [][]int{},
        },
        {
            "case 3",
            []int{0, 0, 0},
            [][]int{
                {0, 0, 0},
            },
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := threeSum(tc.numbers)
            assert.ElementsMatch(t, tc.expected, result)
        })
    }
}

func Test_sortColors(t *testing.T) {
    testCases := []struct {
        name     string
        numbers  []int
        expected []int
    }{
        {
            "case 1",
            []int{2, 0, 2, 1, 1, 0},
            []int{0, 0, 1, 1, 2, 2},
        },
        {
            "case 2",
            []int{2, 0, 1},
            []int{0, 1, 2},
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            sortColors(tc.numbers)
            assert.Equal(t, tc.expected, tc.numbers)
        })
    }
}
