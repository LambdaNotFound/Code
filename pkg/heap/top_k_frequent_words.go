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
	freqMap := make(map[string]int)
	for _, word := range words {
		freqMap[word]++
	}

	wordHeap := &WordHeap{
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

type Item struct {
	str  string
	freq int
}

type WordHeap struct {
	items []Item
	less  func(Item, Item) bool
}

func (w *WordHeap) Len() int           { return len(w.items) }
func (w *WordHeap) Less(i, j int) bool { return w.less(w.items[i], w.items[j]) }
func (w *WordHeap) Swap(i, j int) {
	w.items[i], w.items[j] = w.items[j], w.items[i]
}

func (w *WordHeap) Push(item interface{}) {
	w.items = append(w.items, item.(Item))
}

func (w *WordHeap) Pop() interface{} {
	item := w.items[len(w.items)-1]
	w.items = w.items[:len(w.items)-1]
	return item
}
