package heap

import (
	"container/heap"
	"strings"
)

type Pair struct {
	Val  int
	Freq int
}

/*
 * 347. Top K Frequent Elements
 *
 * follow up: what if there is a steam of words?
 */

// maxHeap: O(n + m log m) time, O(m) space. Better when k is close to m — you're popping most of the heap anyway,
// so the extra log m per operation doesn't hurt much compared to the constant overhead of capping a min-heap.
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

// minHeap: O(n + m log k) time, O(k) space. Better when k << m — keeps only k elements in the heap at a time,
// so each push/evict is log k instead of log m, and memory stays bounded by k regardless of how many unique elements there are.
func topKFrequentClaude(nums []int, k int) []int {
	freqMap := make(map[int]int)
	for _, v := range nums {
		freqMap[v]++
	}

	minHeap := &Heap[Pair]{
		less: func(i, j Pair) bool {
			return i.Freq < j.Freq
		},
	}
	for val, freq := range freqMap {
		heap.Push(minHeap, Pair{Val: val, Freq: freq})
		if minHeap.Len() > k {
			heap.Pop(minHeap)
		}
	}

	res := make([]int, k)
	for i := range k {
		res[i] = heap.Pop(minHeap).(Pair).Val
	}
	return res
}

/*
 * Given a space-delimited string sentence and an integer k, output a list of the k most frequent words in descending order.
 *
 * Example: `"    a!@#b!@c    d#@!e!#@f!@    !@x#@yz    d!@#ef$#@!   x#@y!@z$   !@d$#e!f g$@h!i", 2` ⇒ `["def", xyz"]`
 */
func sanitize(word string) string {
	var b strings.Builder
	for _, ch := range word {
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') {
			b.WriteRune(ch)
		}
	}
	return b.String()
}

/*
var (
    nonAlpha   = regexp.MustCompile(`[^a-zA-Z\s]`)
    multiSpace = regexp.MustCompile(`\s+`)
)

func sanitizeInput(s string) string {
    s = nonAlpha.ReplaceAllString(s, "")     // strip special chars, keep letters and spaces
    s = multiSpace.ReplaceAllString(s, " ")  // collapse runs of whitespace to one space
    return strings.TrimSpace(s)              // drop leading/trailing space before Split
}
*/
