package stack

/**
 * 239. Sliding Window Maximum
 *
 * Input: nums = [1,3,-1,-3,5,3,6,7], k = 3
 * Output: [3,3,5,5,6,7]
 * Explanation:
 * Window position                Max
 * ---------------               -----
 * [1  3  -1] -3  5  3  6  7       3
 * 1 [3  -1  -3] 5  3  6  7       3
 * 1  3 [-1  -3  5] 3  6  7       5
 * 1  3  -1 [-3  5  3] 6  7       5
 * 1  3  -1  -3 [5  3  6] 7       6
 * 1  3  -1  -3  5 [3  6  7]      7
 *
 */
func maxSlidingWindow(nums []int, k int) []int {
	deque := make([]int, 0)
	res := make([]int, 0)
	for l, r := 0, 0; r < len(nums); r++ {
		for len(deque) > 0 && deque[len(deque)-1] < nums[r] {
			deque = deque[:len(deque)-1]
		}
		deque = append(deque, nums[r]) // deque front tracks biggest num

		if r-l+1 == k {
			front := deque[0]
			res = append(res, front)

			if nums[l] == front {
				deque = deque[1:]
			}
			l++
		}
	}

	return res
}

func maxSlidingWindowClaude(nums []int, k int) []int {
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
