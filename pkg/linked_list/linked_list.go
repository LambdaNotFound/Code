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
	pre := (*ListNode)(nil)
	for cur := head; cur != nil; {
		next := cur.Next
		cur.Next = pre
		pre = cur
		cur = next
	}
	return pre
}

func reverseList_recursive(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}

	newHead := reverseList_recursive(head.Next)
	head.Next.Next = head
	head.Next = nil

	return newHead
}

/**
 * 24. Swap Nodes in Pairs
 *
 * //    node1    node2    node3
 * //    nTail    nHead
 * //    nHead.Next        nTail.Next
 * //             pre.Next
 * //    pre      nTail
 *
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
 * 25. Reverse Nodes in k-Group
 */

/**
 * 19. Remove Nth Node From End of List
 */
func removeNthFromEnd(head *ListNode, n int) *ListNode {
	fast, slow := head, head
	for i := 0; i < n; i++ {
		fast = fast.Next
	}

	dummy := &ListNode{}
	dummy.Next = head
	pre := dummy
	for fast != nil {
		fast = fast.Next
		pre = slow
		slow = slow.Next
	}
	pre.Next = slow.Next

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

/**
 * 61. Rotate List
 */
func rotateRight(head *ListNode, k int) *ListNode {
	if head == nil {
		return head
	}

	// single pass: compute length AND capture tail
	length := 1
	tail := head
	for tail.Next != nil {
		tail = tail.Next
		length++
	}

	offset := length - k%length
	if offset == length {
		return head
	}

	// find split point
	cur := head
	for i := 0; i < offset-1; i++ {
		cur = cur.Next
	}

	// rewire
	newHead := cur.Next
	cur.Next = nil
	tail.Next = head
	return newHead
}

/**
 * 2. Add Two Numbers
 */
func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
	dummy := &ListNode{}
	cur := dummy
	carry := 0

	val := func(node *ListNode) int {
		if node == nil {
			return 0
		}
		return node.Val
	}

	for l1 != nil || l2 != nil || carry != 0 {
		sum := val(l1) + val(l2) + carry
		carry = sum / 10
		cur.Next = &ListNode{Val: sum % 10}
		cur = cur.Next
		if l1 != nil {
			l1 = l1.Next
		}
		if l2 != nil {
			l2 = l2.Next
		}
	}

	return dummy.Next
}

/**
 * 328. Odd Even Linked List
 */
func oddEvenList(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}

	dummyOdd, dummyEven := &ListNode{}, &ListNode{}
	preOdd, preEven := dummyOdd, dummyEven
	cur := head
	for cur != nil && cur.Next != nil {
		next := cur.Next
		cur.Next = next.Next
		next.Next = nil

		preOdd.Next = cur
		preOdd = preOdd.Next

		cur = cur.Next

		preEven.Next = next
		preEven = preEven.Next
	}

	if cur != nil {
		cur.Next = dummyEven.Next
	} else {
		preOdd.Next = dummyEven.Next
	}

	return dummyOdd.Next
}

func oddEvenListCalude(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}

	odd, even := head, head.Next
	evenHead := even // anchor even list's head

	for even != nil && even.Next != nil {
		odd.Next = even.Next
		odd = odd.Next
		even.Next = odd.Next
		even = even.Next
	}

	odd.Next = evenHead // splice even list onto odd list
	return head
}
