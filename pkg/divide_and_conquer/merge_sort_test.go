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