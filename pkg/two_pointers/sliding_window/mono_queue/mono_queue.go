package monoqueue

/**
 * Increasing Monotonic Queue: It only keeps elements in increasing order,
 *    and any element that is smaller than the current minimum is removed.
 * => tracking the MAX within window
 *
 * Decreasing Monotonic Queue: It only keeps elements in decreasing order,
 *    and any element that is larger than the current maximum is removed.
 * => tracking the Min within window
 *
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
 *
 * Time: O(n)
 * Space: O(n)
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
 * 1. deque x2, 2. heap x2, 3. sorted set
 * Time: O(n)
 * Space: O(n)
 */
func longestSubarray(nums []int, limit int) int {
	type pair struct {
		index, val int
	}

	res := 1
	maxDeque := []pair{} // decreasing → front is max
	minDeque := []pair{} // increasing → front is min
	l := 0

	for r := 0; r < len(nums); r++ {
		// maintain max deque: pop smaller from back
		for len(maxDeque) > 0 && maxDeque[len(maxDeque)-1].val < nums[r] {
			maxDeque = maxDeque[:len(maxDeque)-1]
		}
		// maintain min deque: pop larger from back
		for len(minDeque) > 0 && minDeque[len(minDeque)-1].val > nums[r] {
			minDeque = minDeque[:len(minDeque)-1]
		}
		maxDeque = append(maxDeque, pair{r, nums[r]})
		minDeque = append(minDeque, pair{r, nums[r]})

		// shrink window from left until valid
		for maxDeque[0].val-minDeque[0].val > limit {
			l++
			if maxDeque[0].index < l {
				maxDeque = maxDeque[1:]
			}
			if minDeque[0].index < l {
				minDeque = minDeque[1:]
			}
		}

		res = max(res, r-l+1)
	}

	return res
}

// Time: O(n^2)
func longestSubarrayBruteForce(nums []int, limit int) int {
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
