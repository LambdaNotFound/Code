package divide_and_conquer

/**
 * 215. Kth Largest Element in an Array
 *
 *
 * Time average O(n)
 * Time worst caseO(n²) — bad pivot choices (can use random pivot to avoid)
 * SpaceO(1) — in-place
 *
 */
func findKthLargest(nums []int, k int) int {
	// kth largest = (n-k)th smallest, convert to 0-indexed target
	target := len(nums) - k
	left, right := 0, len(nums)-1

	for left < right {
		pivot := partition(nums, left, right)
		if pivot == target {
			break
		} else if pivot < target {
			left = pivot + 1 // target is to the right
		} else {
			right = pivot - 1 // target is to the left
		}
	}
	return nums[target]
}

func findKthLargestAlt(nums []int, k int) int {
	left, right := 0, len(nums)-1

	for left < right {
		pivot := partition(nums, left, right)
		rank := len(nums) - pivot // pivot's rank from largest: 1-indexed

		if rank == k {
			return nums[pivot] // found exactly
		} else if rank < k {
			right = pivot - 1 // need something smaller, go left
		} else {
			left = pivot + 1 // need something larger, go right
		}
	}
	return nums[left]
}

func partition(nums []int, left, right int) int {
	pivot := nums[right] // choose rightmost as pivot
	store := left        // store index for elements <= pivot

	for i := left; i < right; i++ {
		if nums[i] <= pivot {
			nums[store], nums[i] = nums[i], nums[store]
			store++
		}
	}
	// place pivot in its final sorted position
	nums[store], nums[right] = nums[right], nums[store]
	return store
}
