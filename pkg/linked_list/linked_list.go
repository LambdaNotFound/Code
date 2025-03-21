package linked_list

import . "gocode/types"

/**
 * 24. Swap Nodes in Pairs
 */
func swapPairs(head *ListNode) *ListNode {
	dummy := ListNode{}
	dummy.Next = head

	pre, newTail := &dummy, dummy.Next
	for newTail != nil && newTail.Next != nil {
		newHead := newTail.Next
		newTail.Next = newHead.Next
		newHead.Next = newTail

		pre.Next = newHead
		pre = newTail
		newTail = pre.Next
	}

	return dummy.Next
}

/**
 * 21. Merge Two Sorted Lists
 * 23. Merge k Sorted Lists
 */
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
