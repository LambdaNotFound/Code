package dynamic_programming

/**
 * 300. Longest Increasing Subsequence
 *
 * LIS[i] length of longest increasing subsequence **ending at index i**
 *
 */
func lengthOfLIS(nums []int) int {
	LIS := make([]int, len(nums))
	for i := range LIS {
		LIS[i] = 1 // base case, LIS of single char
	}

	res := 1
	for i := 1; i < len(nums); i++ {
		for j := 0; j < i; j++ {
			if nums[j] < nums[i] {
				LIS[i] = max(LIS[i], LIS[j]+1)
				res = max(res, LIS[i])
			}
		}
	}

	return res
}

func lengthOfLISBinarySearch(nums []int) int {
	// tails[i] holds the smallest tail for increasing
	// subsequence of length i+1
	// tails is always maintained in sorted order
	tails := []int{}

	for _, num := range nums {
		// Binary search for leftmost position where tails[mid] >= num
		left, right := 0, len(tails)
		for left < right {
			mid := left + (right-left)/2
			if tails[mid] < num {
				left = mid + 1
			} else {
				right = mid
			}
		}

		// left is the insertion point
		if left == len(tails) {
			// num is larger than all tails → extend LIS by 1
			tails = append(tails, num)
		} else {
			// replace to maintain smallest possible tail
			tails[left] = num
		}
	}

	return len(tails)
}
