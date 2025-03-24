package linked_list

import (
	"gocode/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
