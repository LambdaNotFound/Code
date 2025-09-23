package math

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_canCompleteCircuit(t *testing.T) {
    testCases := []struct {
        name     string
        gas      []int
        cost     []int
        expected int
    }{
        {
            "case 1",
            []int{1, 2, 3, 4, 5},
            []int{3, 4, 5, 1, 2},
            3,
        },
        {
            "case 2",
            []int{2, 3, 4},
            []int{3, 4, 3},
            -1,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := canCompleteCircuit(tc.gas, tc.cost)
            assert.Equal(t, tc.expected, result)
        })
    }
}

func Test_jump(t *testing.T) {
    testCases := []struct {
        name     string
        nums      []int
        expected int
    }{
        {
            "case 1",
            []int{2,3,1,1,4},
            2,
        },
        {
            "case 1",
            []int{2,3,0,1,4},
            2,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := jump(tc.nums)
            assert.Equal(t, tc.expected, result)
        })
    }
}

func Test_leastInterval(t *testing.T) {
    testCases := []struct {
        name     string
        tasks    []byte
        n        int
        expected int
    }{
        {
            "case 1",
            []byte{'A', 'A', 'A', 'B', 'B', 'B'},
            2,
            8,
        },
        {
            "case 2",
            []byte{'A', 'C', 'A', 'B', 'D', 'B'},
            1,
            6,
        },
        {
            "case 3",
            []byte{'A', 'A', 'A', 'B', 'B', 'B'},
            3,
            10,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := leastInterval(tc.tasks, tc.n)
            assert.Equal(t, tc.expected, result)
        })
    }
}

func TestProductExceptSelf(t *testing.T) {
    tests := []struct {
        name     string
        nums     []int
        expected []int
    }{
        {
            name:     "basic example",
            nums:     []int{1, 2, 3, 4},
            expected: []int{24, 12, 8, 6},
        },
        {
            name:     "with zero",
            nums:     []int{-1, 1, 0, -3, 3},
            expected: []int{0, 0, 9, 0, 0},
        },
        {
            name:     "single element",
            nums:     []int{5},
            expected: []int{1}, // usually problem guarantees len >= 2, but we handle it
        },
        {
            name:     "two elements",
            nums:     []int{2, 3},
            expected: []int{3, 2},
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            got := productExceptSelf(tc.nums)
            assert.Equal(t, tc.expected, got)

            got = productExceptSelfWithMultiplier(tc.nums)
            assert.Equal(t, tc.expected, got)
        })
    }
}