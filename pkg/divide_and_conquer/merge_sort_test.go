package divide_and_conquer

import (
	"gocode/utils"

	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_sortList(t *testing.T) {
    testCases := []struct {
        name     string
        list     []int
        expected []int
    }{
        {
            "case 1",
            []int{4, 2, 1, 3},
            []int{1, 2, 3, 4},
        },
        {
            "case 2",
            []int{-1, 5, 3, 4, 0},
            []int{-1, 0, 3, 4, 5},
        },
        {
            "case 3",
            []int{},
            []int{},
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            list := utils.CreateLinkedList(tc.list)
            result := sortListMergeSort(list)
            expected := utils.CreateLinkedList(tc.expected)
            isEqual := utils.VerifyLinkedLists(result, expected)
            assert.Equal(t, true, isEqual)
        })
    }
}

func Test_MergeSort(t *testing.T) {
    tests := []struct {
        name     string
        input    []int
        expected []int
    }{
        {
            name:     "Empty slice",
            input:    []int{},
            expected: []int{},
        },
        {
            name:     "Single element",
            input:    []int{42},
            expected: []int{42},
        },
        {
            name:     "Already sorted",
            input:    []int{1, 2, 3, 4, 5},
            expected: []int{1, 2, 3, 4, 5},
        },
        {
            name:     "Reverse order",
            input:    []int{5, 4, 3, 2, 1},
            expected: []int{1, 2, 3, 4, 5},
        },
        {
            name:     "With duplicates",
            input:    []int{4, 2, 4, 1, 3, 2},
            expected: []int{1, 2, 2, 3, 4, 4},
        },
        {
            name:     "All same elements",
            input:    []int{7, 7, 7, 7},
            expected: []int{7, 7, 7, 7},
        },
        {
            name:     "Contains negative numbers",
            input:    []int{-3, -1, -7, 2, 0},
            expected: []int{-7, -3, -1, 0, 2},
        },
        {
            name:     "Random order",
            input:    []int{10, 1, 5, 8, 3, 2, 9, 4, 7, 6},
            expected: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            got := MergeSort(tc.input)
            assert.Equal(t, tc.expected, got)
        })
    }
}