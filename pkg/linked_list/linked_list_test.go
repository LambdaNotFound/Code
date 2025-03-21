package linked_list

import (
	. "gocode/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createLinkedList(arr []int) *ListNode {
	dummy := ListNode{}

	cur := &dummy
	for _, val := range arr {
		cur.Next = &ListNode{Val: val}
		cur = cur.Next
	}

	return dummy.Next
}

func verifyLinkedLists(list1 *ListNode, list2 *ListNode) bool {
	for list1 != nil && list2 != nil {
		if list1.Val != list2.Val {
			return false
		}
		list1 = list1.Next
		list2 = list2.Next
	}

	if list1 != nil || list2 != nil {
		return false
	}

	return true
}

func Test_mergeTwoLists(t *testing.T) {
	list1 := createLinkedList([]int{1, 3, 5, 7})
	list2 := createLinkedList([]int{2, 4, 6, 8})

	want := createLinkedList([]int{1, 2, 3, 4, 5, 6, 7, 8})
	got := mergeTwoLists(list1, list2)
	equal := verifyLinkedLists(want, got)

	assert.Equal(t, true, equal)
}
