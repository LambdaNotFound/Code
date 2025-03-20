package linked_list

import . "gocode/types"

func mergeTwoLists(list1 *ListNode, list2 *ListNode) *ListNode {
	dummy := ListNode{}
	cur := &dummy

	for list1 != nil && list2 != nil {
		if list1.Val < list2.Val {
			cur.Next = list1
			list1 = list1.Next
		} else {
			cur.Next = list2
			list2 = list2.Next
		}
		cur = cur.Next
	}

	if list1 != nil {
		cur.Next = list1
	} else if list2 != nil {
		cur.Next = list2
	}

	return dummy.Next
}
