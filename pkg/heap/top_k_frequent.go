package heap

import "container/heap"

type Pair struct {
	Val  int
	Freq int
}

type TopKHeap struct {
	items []Pair
	less  func(Pair, Pair) bool
}

func (t *TopKHeap) Less(i int, j int) bool { return t.less(t.items[i], t.items[j]) }
func (t *TopKHeap) Swap(i int, j int)      { t.items[i], t.items[j] = t.items[j], t.items[i] }
func (t *TopKHeap) Len() int               { return len(t.items) }

func (t *TopKHeap) Push(item interface{}) {
	t.items = append(t.items, item.(Pair))
}
func (t *TopKHeap) Pop() interface{} {
	item := t.items[len(t.items)-1]
	t.items = t.items[:len(t.items)-1]
	return item
}

/*
 * 347. Top K Frequent Elements
 */

func topKFrequent(nums []int, k int) []int {
	freqMap := make(map[int]int)
	for _, v := range nums {
		freqMap[v] += 1
	}

	maxHeap := &TopKHeap{
		less: func(i Pair, j Pair) bool {
			return i.Freq > j.Freq
		},
	}
	for k, v := range freqMap {
		heap.Push(maxHeap, Pair{
			Val:  k,
			Freq: v,
		})
	}

	res := make([]int, k)
	for i := 0; i < k; i++ {
		item := heap.Pop(maxHeap).(Pair)
		res[i] = item.Val
	}
	return res
}
