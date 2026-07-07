package hashmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_longestConsecutive(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected int
	}{
		{name: "leetcode_example1", nums: []int{100, 4, 200, 1, 3, 2}, expected: 4},
		{name: "leetcode_example2", nums: []int{0, 3, 7, 2, 5, 8, 4, 6, 0, 1}, expected: 9},
		{name: "empty", nums: []int{}, expected: 0},
		{name: "single", nums: []int{5}, expected: 1},
		{name: "all_same", nums: []int{3, 3, 3}, expected: 1},
		{name: "no_consecutive", nums: []int{10, 20, 30}, expected: 1},
		{name: "negatives", nums: []int{-3, -2, -1, 0, 1}, expected: 5},
		{name: "two_sequences", nums: []int{1, 2, 3, 10, 11, 12, 13}, expected: 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, longestConsecutive(tt.nums))
		})
	}
}

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
