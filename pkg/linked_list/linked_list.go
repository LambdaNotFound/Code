package linked_list

import (
	. "gocode/types"
)

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

/**
 * 206. Reverse Linked List
 * 1). iterative impl
 * 2). recursive impl
 */
func reverseList(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}

	tail := head.Next
	head.Next = nil // reset pointer
	newHead := reverseList(tail)
	tail.Next = head

	return newHead
}

func reverseList_iterative(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}

	var prev, curr *ListNode = nil, head // 3 pointers
	for curr != nil {
		next := curr.Next
		curr.Next = prev
		prev = curr
		curr = next
	}
	return prev
}

func reverseList_recursive2(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}

	newHead := reverseList_recursive2(head.Next)
	head.Next.Next = head
	head.Next = nil

	return newHead
}

/**
 * 24. Swap Nodes in Pairs
 * 25. Reverse Nodes in k-Group
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
 * 141. Linked List Cycle
 *
 * if there is a loop and the faster pointer is approaching the
 * slow pointer, there can only be 2 cases:
 *     1. the faster pointer is 1 step behind the slow pointer
 *     2. the faster pointer is 2 step behind the slow pointer
 *
 * in case 1, fater pointer will meet slow pointer in next step
 * in case 2, it will reduces to case 1 in next step
 */
func hasCycle(head *ListNode) bool {
	fast, slow := head, head
	for fast != nil && fast.Next != nil {
		fast = fast.Next.Next
		slow = slow.Next
		if fast == slow {
			return true
		}
	}
	return false
}

/**
 * 876. Middle of the Linked List
 */
func middleNode(head *ListNode) *ListNode {
	fast, slow := head, head
	for fast != nil && fast.Next != nil {
		fast = fast.Next.Next
		slow = slow.Next
	}

	return slow
}

/**
 * 138. Copy List with Random Pointer
 *
 * 1. Hash map. O(n) time and O(n) space.
 * 2. In-place weaving. No extra space needed.
 *
 */
func copyRandomList(head *Node) *Node {
	for cur := head; cur != nil; { // Weave copied nodes into the original list
		copy := &Node{Val: cur.Val}
		next := cur.Next

		copy.Next = next
		cur.Next = copy

		cur = next
	}

	for cur := head; cur != nil; { // Set random pointers for copies before unweaving
		copy := cur.Next
		if cur.Random != nil {
			copy.Random = cur.Random.Next
		}
		cur = copy.Next
	}

	dummy := &Node{}
	copyTail := dummy
	for cur := head; cur != nil; {
		copy := cur.Next
		nextOriginal := copy.Next

		copyTail.Next = copy // Unweave the lists: separate copy
		copyTail = copyTail.Next

		cur.Next = nextOriginal // Unweave the lists: separate original
		cur = nextOriginal
	}

	return dummy.Next
}
