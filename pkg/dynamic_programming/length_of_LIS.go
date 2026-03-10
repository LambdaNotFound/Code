package dynamic_programming

/**
 * 300. Longest Increasing Subsequence
 */
func lengthOfLIS(nums []int) int {
	LIS := make([]int, len(nums))
	for i := range LIS {
		LIS[i] = 1
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
