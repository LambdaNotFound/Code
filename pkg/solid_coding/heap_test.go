package solid_coding

import (
	. "gocode/types"
	"gocode/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_mergeKLists(t *testing.T) {
    testCases := []struct {
        name     string
        input    []*ListNode
        expected *ListNode
    }{
        {
            "case 1",
            []*ListNode{
                utils.CreateLinkedList([]int{1, 4, 5}),
                utils.CreateLinkedList([]int{1, 3, 4}),
                utils.CreateLinkedList([]int{2, 6}),
            },
            utils.CreateLinkedList([]int{1, 1, 2, 3, 4, 4, 5, 6}),
        },
        {
            "case 2",
            []*ListNode{},
            utils.CreateLinkedList([]int{}),
        },
        {
            "case 3",
            []*ListNode{
                utils.CreateLinkedList([]int{}),
            },
            utils.CreateLinkedList([]int{}),
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := mergeKLists(tc.input)
            equal := utils.VerifyLinkedLists(tc.expected, result)
            assert.Equal(t, true, equal)
        })
    }
}
