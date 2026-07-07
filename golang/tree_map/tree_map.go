package treemap

import "sort"

/**
 * Implement a minimal sorted set with sort.Search (fastest to write)
 * If you only need ceiling/floor queries, sort.Search on a sorted slice is your friend:
 *
 * // Ceiling: first index where keys[i] >= target
 * i := sort.Search(len(keys), func(i int) bool { return keys[i] >= target })
 *
 * // Floor: last index where keys[i] <= target
 * i := sort.Search(len(keys), func(i int) bool { return keys[i] > target }) - 1
 */
type TreeMap struct {
	keys []int // <- sorted
	m    map[int]int
}

func (t *TreeMap) Put(k, v int) {
	if _, ok := t.m[k]; !ok {
		i := sort.SearchInts(t.keys, k)
		t.keys = append(t.keys, 0)
		copy(t.keys[i+1:], t.keys[i:]) // shift
		t.keys[i] = k
	}
	t.m[k] = v
}

func (t *TreeMap) Floor(k int) (int, bool) {
	i := sort.SearchInts(t.keys, k+1) - 1
	if i < 0 {
		return 0, false
	}
	return t.keys[i], true
}

func (t *TreeMap) Ceil(k int) (int, bool) {
	i := sort.SearchInts(t.keys, k)
	if i >= len(t.keys) {
		return 0, false
	}
	return t.keys[i], true
}

/**
 * 1438. Longest Continuous Subarray With Absolute Diff Less Than or Equal to Limit
 *
 * window = sorted map { value: frequency }
 *
 * for each right pointer:
 *    add nums[right] to window
 *
 *    while window.max - window.min > limit:
 *        remove nums[left] from window
 *        left++
 *
 *     ans = max(ans, right - left + 1)
 */
func longestSubarray(nums []int, limit int) int {
	// TreeMap: sorted keys + count map
	keys := []int{}
	freq := map[int]int{}

	insert := func(v int) {
		freq[v]++
		if freq[v] == 1 { // new key
			i := sort.SearchInts(keys, v)
			keys = append(keys, 0)
			copy(keys[i+1:], keys[i:])
			keys[i] = v
		}
	}

	delete := func(v int) {
		freq[v]--
		if freq[v] == 0 {
			delete(freq, v)
			i := sort.SearchInts(keys, v)
			keys = append(keys[:i], keys[i+1:]...)
		}
	}

	left, ans := 0, 0
	for right := 0; right < len(nums); right++ {
		insert(nums[right])

		for keys[len(keys)-1]-keys[0] > limit {
			delete(nums[left])
			left++
		}

		if right-left+1 > ans {
			ans = right - left + 1
		}
	}
	return ans
}
