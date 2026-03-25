package heap

import (
	"container/heap"
	. "gocode/types"
)

type ListNodeHeap struct {
	items []*ListNode
	less  func(*ListNode, *ListNode) bool
}

func (h *ListNodeHeap) Len() int               { return len(h.items) }
func (h *ListNodeHeap) Less(i int, j int) bool { return h.less(h.items[i], h.items[j]) }
func (h *ListNodeHeap) Swap(i int, j int)      { h.items[i], h.items[j] = h.items[j], h.items[i] }

func (h *ListNodeHeap) Push(item interface{}) {
	h.items = append(h.items, item.(*ListNode))
}
func (h *ListNodeHeap) Pop() interface{} {
	item := h.items[len(h.items)-1]
	h.items = h.items[:len(h.items)-1]
	return item
}

/**
 * 23. Merge k Sorted Lists
 */
func mergeKLists(lists []*ListNode) *ListNode {
	minHeap := &ListNodeHeap{
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
