package divide_and_conquer

/**
 * QuickSelect
 *
 * Time average O(n) — each partition eliminates roughly half the remaining
 *                      range, and only the side containing k is recursed into
 * Time worst case O(n²) — bad pivot choices (e.g. already-sorted input with a
 *                          last-element pivot); can use random pivot to avoid
 * Space average O(log n) — recursion stack depth; no extra data is copied
 *                           since nums is partitioned in place
 * Space worst case O(n) — recursion stack depth when partitions are maximally
 *                          unbalanced
 *
 */
func QuickSelect(nums []int, k int) int {
	var quickSelect func(nums []int, lo, hi, k int) int

	// qs returns the k-th smallest element (0-indexed) in nums[lo..hi].
	quickSelect = func(nums []int, lo, hi, k int) int {
		if lo == hi {
			return nums[lo]
		}
		// partition: last element as pivot, returns its final sorted index
		pivot := nums[hi]
		i := lo
		for j := lo; j < hi; j++ {
			if nums[j] <= pivot {
				nums[i], nums[j] = nums[j], nums[i]
				i++
			}
		}
		nums[i], nums[hi] = nums[hi], nums[i]
		p := i

		if k == p {
			return nums[p]
		} else if k < p {
			return quickSelect(nums, lo, p-1, k)
		}
		return quickSelect(nums, p+1, hi, k)
	}

	return quickSelect(nums, 0, len(nums)-1, k)
}
