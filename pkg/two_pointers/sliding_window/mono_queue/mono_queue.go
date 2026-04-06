package monoqueue

/**
 * Increasing Monotonic Queue: It only keeps elements in increasing order,
 *    and any element that is smaller than the current minimum is removed.
 * Decreasing Monotonic Queue: It only keeps elements in decreasing order,
 *    and any element that is larger than the current maximum is removed.
 *
 *
 * We use a MonoQueue as the data structure to query both the minimum and maximum values in O(1) time complexity.
 */

/**
 * 239. Sliding Window Maximum
 *
 * Input: nums = [1,3,-1,-3,5,3,6,7], k = 3
 * Output: [3,3,5,5,6,7]
 *                                             Decreasing Monotonic Queue
 *  Window position                Max         Window
 *  ---------------               -----       --------
 * [1  3  -1] -3  5  3  6  7       3          [1, 2]
 *  1 [3  -1  -3] 5  3  6  7       3          [1, 2, 3]
 *  1  3 [-1  -3  5] 3  6  7       5          [4]
 *  1  3  -1 [-3  5  3] 6  7       5          [4, 5]
 *  1  3  -1  -3 [5  3  6] 7       6          [6]
 *  1  3  -1  -3  5 [3  6  7]      7          [7]
 */
func maxSlidingWindow(nums []int, k int) []int {
	result := make([]int, 0, len(nums)-k+1)
	deque := make([]int, 0) // stores indices, decreasing order of values

	for i := 0; i < len(nums); i++ {
		// remove indices outside current window from front
		for len(deque) > 0 && deque[0] < i-k+1 {
			deque = deque[1:]
		}

		// maintain decreasing order — pop smaller values from back
		for len(deque) > 0 && nums[deque[len(deque)-1]] < nums[i] {
			deque = deque[:len(deque)-1]
		}

		deque = append(deque, i)

		// window fully formed — front of deque is max
		if i >= k-1 {
			result = append(result, nums[deque[0]])
		}
	}

	return result
}

/**
 * 1438. Longest Continuous Subarray With Absolute Diff Less Than or Equal to Limit
 *
 */
func longestSubarray(nums []int, limit int) int {
	res := 1
	for i := 0; i < len(nums); i++ {
		minVal, maxVal := nums[i], nums[i]
		for j := i + 1; j < len(nums); j++ {
			minVal = min(minVal, nums[j])
			maxVal = max(maxVal, nums[j])

			if maxVal-minVal <= limit {
				res = max(res, j-i+1)
			}
		}
	}
	return res
}
