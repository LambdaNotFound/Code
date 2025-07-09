package divide_and_conquer

import . "gocode/types"

/**
 * Merge Sort linked list (recursive structure)
 *
 * 148. Sort List
 */
func sortListMergeSort(head *ListNode) *ListNode {
    if head == nil || head.Next == nil {
        return head
    }
    slow, fast := head, head
    for fast.Next != nil && fast.Next.Next != nil {
        slow, fast = slow.Next, fast.Next.Next
    }

    firstHalfTail := slow
    slow = slow.Next
    firstHalfTail.Next = nil // divide the list into two

    firstHalf, secondHalf := sortListMergeSort(head), sortListMergeSort(slow)
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
