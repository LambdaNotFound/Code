package hashmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainsDuplicate(t *testing.T) {
    tests := []struct {
        name     string
        input    []int
        expected bool
    }{
        {
            name:     "Empty slice",
            input:    []int{},
            expected: false,
        },
        {
            name:     "Single element",
            input:    []int{1},
            expected: false,
        },
        {
            name:     "Two elements no duplicate",
            input:    []int{1, 2},
            expected: false,
        },
        {
            name:     "Two elements with duplicate",
            input:    []int{1, 1},
            expected: true,
        },
        {
            name:     "Multiple elements no duplicate",
            input:    []int{1, 2, 3, 4, 5},
            expected: false,
        },
        {
            name:     "Multiple elements with duplicate",
            input:    []int{1, 2, 3, 4, 2},
            expected: true,
        },
        {
            name:     "All duplicates",
            input:    []int{7, 7, 7, 7},
            expected: true,
        },
        {
            name:     "Negative numbers unique",
            input:    []int{-1, -2, -3},
            expected: false,
        },
        {
            name:     "Negative numbers with duplicate",
            input:    []int{-1, -1, 2},
            expected: true,
        },
        {
            name:     "Large variety",
            input:    []int{10, 20, 30, 40, 10},
            expected: true,
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            got := containsDuplicate(tc.input)
            assert.Equal(t, tc.expected, got)
        })
    }
}
