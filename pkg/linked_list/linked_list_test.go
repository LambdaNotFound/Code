package linked_list

import (
	. "gocode/types"
	"gocode/utils"

	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_reverseList(t *testing.T) {
    testCases := []struct {
        name     string
        input  []int
        expected *ListNode
    }{
        {
            "case 1",
            []int{1,2,3,4,5},
            utils.CreateLinkedList([]int{5,4,3,2,1}),
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            head := utils.CreateLinkedList(tc.input)
            result := reverseList(head)
            equal := utils.VerifyLinkedLists(tc.expected, result)
            assert.Equal(t, true, equal)

            head = utils.CreateLinkedList(tc.input)
            result = reverseList_recursive(head)
            equal = utils.VerifyLinkedLists(tc.expected, result)
            assert.Equal(t, true, equal)

            head = utils.CreateLinkedList(tc.input)
            result = reverseList_recursive2(head)
            equal = utils.VerifyLinkedLists(tc.expected, result)
            assert.Equal(t, true, equal)
       })
    }
}

func Test_mergeTwoLists(t *testing.T) {
    list1 := utils.CreateLinkedList([]int{1, 3, 5, 7})
    list2 := utils.CreateLinkedList([]int{2, 4, 6, 8})

    want := utils.CreateLinkedList([]int{1, 2, 3, 4, 5, 6, 7, 8})
    got := mergeTwoLists(list1, list2)
    equal := utils.VerifyLinkedLists(want, got)

    assert.Equal(t, true, equal)
}

func Test_swapPairs(t *testing.T) {
    list1 := utils.CreateLinkedList([]int{1, 2, 3, 4})

    want := utils.CreateLinkedList([]int{2, 1, 4, 3})
    got := swapPairs(list1)
    equal := utils.VerifyLinkedLists(want, got)

    assert.Equal(t, true, equal)
}
