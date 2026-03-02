package heap

import (
	"container/heap"
	. "gocode/types"
)

/**
 * 23. Merge k Sorted Lists
 */
func mergeKLists(lists []*ListNode) *ListNode {
	h := &ListNodeMinHeap{}
	heap.Init(h)
	for _, val := range lists {
		if val != nil {
			heap.Push(h, val)
		}
	}

	dummy := ListNode{}
	cur := &dummy
	for h.Len() > 0 {
		tmp := heap.Pop(h).(*ListNode)
		cur.Next = tmp
		cur = tmp

		if tmp.Next != nil {
			heap.Push(h, tmp.Next)
		}
	}

	return dummy.Next
}
