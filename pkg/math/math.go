package math

/**
 * 189. Rotate Array
 */
func rotate(nums []int, k int) {
	l := len(nums)
	if len(nums) != 0 {
		copy(nums, append(nums[l-k%l:], nums[0:l-k%l]...))
	}
}
