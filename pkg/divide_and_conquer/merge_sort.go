package divide_and_conquer

import . "gocode/types"

/**
 * Merge Sort linked list (recursive structure)
 */

// MergeSort sorts an integer slice using merge sort
func MergeSort(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}

	mid := len(arr) / 2

	// Divide
	left := MergeSort(arr[:mid])
	right := MergeSort(arr[mid:])

	// Conquer (merge two sorted halves)
	return mergeSlices(left, right)
}

// merge combines two sorted slices
func mergeSlices(left, right []int) []int {
	result := make([]int, 0, len(left)+len(right))
	i, j := 0, 0

	for i < len(left) && j < len(right) {
		if left[i] < right[j] {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}

	// Append leftovers
	result = append(result, left[i:]...)
	result = append(result, right[j:]...)

	return result
}

/*
 * 148. Sort List
 */
func sortList(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}

	slow, fast := head, head.Next
	for fast.Next != nil && fast.Next.Next != nil {
		slow, fast = slow.Next, fast.Next.Next
	}
	mid := slow.Next // exact midpoint
	slow.Next = nil  // divide the list into two

	firstHalf, secondHalf := sortList(head), sortList(mid)
	return merge(firstHalf, secondHalf)
}

func merge(first *ListNode, second *ListNode) *ListNode {
	dummy := ListNode{}
	cur := &dummy

	for first != nil && second != nil {
		if first.Val < second.Val {
			cur.Next = first
			first = first.Next
		} else {
			cur.Next = second
			second = second.Next
		}
		cur = cur.Next
	}

	if first != nil {
		cur.Next = first
	} else if second != nil {
		cur.Next = second
	}
	return dummy.Next
}

/**
 * 109. Convert Sorted List to Binary Search Tree
 *
 */
func sortedListToBST(head *ListNode) *TreeNode {
	if head == nil {
		return nil
	}
	if head.Next == nil {
		return &TreeNode{Val: head.Val}
	}

	slow, fast := head, head.Next
	for fast.Next != nil && fast.Next.Next != nil {
		slow, fast = slow.Next, fast.Next.Next
	}
	mid := slow.Next // exact midpoint
	slow.Next = nil  // divide the list into two

	root := &TreeNode{Val: mid.Val}
	root.Left = sortedListToBST(head)
	root.Right = sortedListToBST(mid.Next)

	return root
}
