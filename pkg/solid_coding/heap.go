package solid_coding

import (
	"container/heap"
	. "gocode/types"
)

// MinHeap struct
type MinHeap []*ListNode

// Implement heap.Interface methods
func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i].Val < h[j].Val } // Min Heap
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MinHeap) Push(x interface{}) {
    *h = append(*h, x.(*ListNode)) // Append dynamically
}

func (h *MinHeap) Pop() interface{} {
    old := *h
    n := len(old)
    x := old[n-1] // Remove last element
    *h = old[:n-1]
    return x
}

/**
 * 23. Merge k Sorted Lists
 */

func mergeKLists(lists []*ListNode) *ListNode {
    h := &MinHeap{}
    heap.Init(h) // Initialize heap

    // Dynamically insert elements
    for _, val := range lists {
        if val != nil {
            heap.Push(h, val)
        }
    }

    dummy := ListNode{}
    curr := &dummy
    for h.Len() > 0 {
        temp := heap.Pop(h).(*ListNode)
        curr.Next = temp
        curr = temp

        if temp != nil && temp.Next != nil {
            heap.Push(h, temp.Next)
        }
    }

    return dummy.Next
}
