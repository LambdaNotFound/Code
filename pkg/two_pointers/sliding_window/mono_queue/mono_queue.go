package monoqueue

import "math"

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
 * 1. deque x2
 *     maxDeque[0] tracks the max number in the window
 *     minDeque[0] tracks the min number in the window
 * 2. heap x2
 * 3. sorted set
 *
 * Time: O(n)
 * Space: O(n)
 */
func longestSubarray(nums []int, limit int) int {
	maxDeque, minDeque := []int{}, []int{}
	res := math.MinInt
	for left, right := 0, 0; right < len(nums); right++ {
		for len(maxDeque) > 0 && maxDeque[len(maxDeque)-1] < nums[right] {
			maxDeque = maxDeque[:len(maxDeque)-1]
		}
		maxDeque = append(maxDeque, nums[right])

		for len(minDeque) > 0 && minDeque[len(minDeque)-1] > nums[right] {
			minDeque = minDeque[:len(minDeque)-1]
		}
		minDeque = append(minDeque, nums[right])

		for maxDeque[0]-minDeque[0] > limit {
			num := nums[left]
			left++

			if num == maxDeque[0] {
				maxDeque = maxDeque[1:]
			}
			if num == minDeque[0] {
				minDeque = minDeque[1:]
			}
		}

		res = max(res, right-left+1)
	}

	return res
}

// Time: O(n^2)
func longestSubarrayBruteForce(nums []int, limit int) int {
	maxLen := 1
	for i := 0; i < len(nums); i++ {
		minVal, maxVal := nums[i], nums[i]
		for j := i + 1; j < len(nums); j++ {
			minVal = min(minVal, nums[j])
			maxVal = max(maxVal, nums[j])

			if maxVal-minVal <= limit {
				maxLen = max(maxLen, j-i+1)
			}
		}
	}
	return maxLen
}
