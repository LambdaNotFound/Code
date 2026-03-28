package heap

import (
	"container/heap"
	"math"
)

/**
 * 632. Smallest Range Covering Elements from K Lists
 */
func smallestRange(nums [][]int) []int {
	minHeap := &ItemHeap{
		less: func(a, b NumItem) bool {
			return a.num < b.num
		},
	}
	upper := math.MinInt
	for i := range nums {
		val := nums[i][0]
		upper = max(upper, val)

		heap.Push(minHeap, NumItem{
			num:      val,
			idx:      0,
			arrayIdx: i,
		})
	}

	res := []int{minHeap.items[0].num, upper}
	for minHeap.Len() > 0 {
		top := heap.Pop(minHeap).(NumItem)

		if upper-top.num < res[1]-res[0] {
			res = []int{top.num, upper}
		}

		if top.idx == len(nums[top.arrayIdx])-1 {
			break
		}

		nextNum := nums[top.arrayIdx][top.idx+1]
		upper = max(upper, nextNum)
		heap.Push(minHeap, NumItem{
			num:      nextNum,
			idx:      top.idx + 1,
			arrayIdx: top.arrayIdx,
		})
	}

	return res
}

type NumItem struct {
	num      int
	idx      int
	arrayIdx int
}

type ItemHeap struct {
	items []NumItem
	less  func(NumItem, NumItem) bool
}

func (h *ItemHeap) Len() int           { return len(h.items) }
func (h *ItemHeap) Swap(i, j int)      { h.items[i], h.items[j] = h.items[j], h.items[i] }
func (h *ItemHeap) Less(i, j int) bool { return h.less(h.items[i], h.items[j]) }

func (h *ItemHeap) Push(item interface{}) {
	h.items = append(h.items, item.(NumItem))
}

func (h *ItemHeap) Pop() interface{} {
	item := h.items[len(h.items)-1]
	h.items = h.items[:len(h.items)-1]
	return item
}
