package utils

import . "gocode/types"

func CreateLinkedList(arr []int) *ListNode {
    dummy := ListNode{}

    cur := &dummy
    for _, val := range arr {
        cur.Next = &ListNode{Val: val}
        cur = cur.Next
    }

    return dummy.Next
}

func VerifyLinkedLists(list1 *ListNode, list2 *ListNode) bool {
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