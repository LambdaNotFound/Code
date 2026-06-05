package heap

import (
	"container/heap"
	. "gocode/types"
)

/**
 * 23. Merge k Sorted Lists
 */
func mergeKLists(lists []*ListNode) *ListNode {
	minHeap := &Heap[*ListNode]{
		less: func(i *ListNode, j *ListNode) bool {
			return i.Val < j.Val
		},
	}

	for _, list := range lists {
		if list != nil {
			heap.Push(minHeap, list)
		}
	}

	dummy := &ListNode{}
	cur := dummy
	for minHeap.Len() > 0 {
		tmp := heap.Pop(minHeap).(*ListNode)
		cur.Next = tmp
		cur = tmp

		if tmp.Next != nil { // push the next item in sorted linked list
			heap.Push(minHeap, tmp.Next)
		}
	}

	return dummy.Next
}
