package linked_list

import (
	. "gocode/types"
)

/**
 * 143. Reorder List
 */
func reorderList(head *ListNode) {
	if head == nil || head.Next == nil {
		return
	}

	// Step 1: Find middle via fast/slow pointers
	slow, fast := head, head
	for fast.Next != nil && fast.Next.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
	}
	// slow is now at the middle — cut the list here
	secondHalf := slow.Next
	slow.Next = nil

	// Step 2: Reverse second half
	var prev *ListNode
	for curr := secondHalf; curr != nil; {
		next := curr.Next
		curr.Next = prev
		prev = curr
		curr = next
	}
	// prev is now the head of the reversed second half

	// Step 3: Merge two halves — l1 from front, l2 from back
	l1, l2 := head, prev
	for l2 != nil {
		l1Next, l2Next := l1.Next, l2.Next
		l1.Next = l2
		l2.Next = l1Next
		l1, l2 = l1Next, l2Next
	}
}

/**
 * 234. Palindrome Linked List
 */
func isPalindrome(head *ListNode) bool {
	// Step 1: find middle via fast/slow pointers
	slow, fast := head, head
	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
	}

	// Step 2: reverse second half
	var prev *ListNode
	for cur := slow; cur != nil; {
		next := cur.Next
		cur.Next = prev
		prev = cur
		cur = next
	}

	// Step 3: compare first and reversed second half
	left, right := head, prev
	for right != nil {
		if left.Val != right.Val {
			return false
		}
		left = left.Next
		right = right.Next
	}

	return true
}
