package heap

import (
	"container/heap"
	"slices"
)

/**
 * 692. Top K Frequent Words
 *
 * Return the answer sorted by the frequency from highest to lowest.
 * Sort the words with the same frequency by their lexicographical order.
 */
func topKFrequentWords(words []string, k int) []string {
	type Item struct {
		str  string
		freq int
	}

	freqMap := make(map[string]int)
	for _, word := range words {
		freqMap[word]++
	}

	wordHeap := &Heap[Item]{
		less: func(a Item, b Item) bool {
			if a.freq == b.freq {
				return a.str > b.str
			}
			return a.freq < b.freq
		},
	}
	for word, v := range freqMap {
		item := Item{
			str:  word,
			freq: v,
		}

		heap.Push(wordHeap, item)
		if wordHeap.Len() > k {
			heap.Pop(wordHeap)
		}
	}

	res := []string{}
	for wordHeap.Len() > 0 {
		item := heap.Pop(wordHeap).(Item)
		res = append(res, item.str)
	}
	slices.Reverse(res)

	return res
}
