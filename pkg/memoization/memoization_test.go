package memoization

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClimbStairs(t *testing.T) {
    tests := []struct {
        name     string
        n        int
        expected int
    }{
        {
            name:     "n = 1",
            n:        1,
            expected: 1,
        },
        {
            name:     "n = 2",
            n:        2,
            expected: 2,
        },
        {
            name:     "n = 3",
            n:        3,
            expected: 3,
        },
        {
            name:     "n = 4",
            n:        4,
            expected: 5,
        },
        {
            name:     "n = 5",
            n:        5,
            expected: 8,
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            got := climbStairs1(tc.n)
            assert.Equal(t, tc.expected, got, "climbStairs(%d)", tc.n)
        
            got = climbStairs2(tc.n)
            assert.Equal(t, tc.expected, got, "climbStairs(%d)", tc.n)

            got = climbStairs3(tc.n)
            assert.Equal(t, tc.expected, got, "climbStairs(%d)", tc.n)

            got = climbStairs4(tc.n)
            assert.Equal(t, tc.expected, got, "climbStairs(%d)", tc.n)
        })
    }
}