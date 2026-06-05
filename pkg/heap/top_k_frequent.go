package heap

import "container/heap"

type Pair struct {
	Val  int
	Freq int
}

/*
 * 347. Top K Frequent Elements
 */

func topKFrequent(nums []int, k int) []int {
	freqMap := make(map[int]int)
	for _, v := range nums {
		freqMap[v] += 1
	}

	maxHeap := &Heap[Pair]{
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
