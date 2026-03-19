package math

import "slices"

/**
 * 189. Rotate Array
 */
func rotateArraySlicesReverse(nums []int, k int) {
	index := len(nums) - k%len(nums)
	slices.Reverse(nums[:index])
	slices.Reverse(nums[index:])
	slices.Reverse(nums)
}

func rotateArray(nums []int, k int) {
	l := len(nums)
	if len(nums) != 0 {
		copy(nums, append(nums[l-k%l:], nums[0:l-k%l]...))
	}
}

/**
 * 48. Rotate Image
 */
func rotateImage(matrix [][]int) {
	slices.Reverse(matrix)

	for i := 0; i < len(matrix); i++ {
		for j := i; j < len(matrix[0]); j++ {
			matrix[i][j], matrix[j][i] = matrix[j][i], matrix[i][j]
		}
	}
}
