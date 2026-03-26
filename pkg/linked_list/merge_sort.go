package linked_list

import (
	. "gocode/types"
)

/**
 * 148. Sort List
 */
func sortList(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}

	merge := func(a, b *ListNode) *ListNode {
		dummy := &ListNode{}
		cur := dummy
		for a != nil && b != nil {
			if a.Val < b.Val {
				cur.Next = a
				a = a.Next
			} else {
				cur.Next = b
				b = b.Next
			}
			cur = cur.Next
		}
		if a != nil {
			cur.Next = a
		} else {
			cur.Next = b
		}
		return dummy.Next
	}

	// fast starts at head.Next so slow stops ONE node before midpoint
	slow, fast := head, head.Next
	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
	}

	// slow IS the node before mid — no pre pointer needed
	mid := slow.Next
	slow.Next = nil

	l1 := sortList(head)
	l2 := sortList(mid)
	return merge(l1, l2)
}
