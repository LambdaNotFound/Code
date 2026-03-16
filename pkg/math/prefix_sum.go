package math

/**
 * 238. Product of Array Except Self
 *
 * [a, b, c, d]
 *
 * LEFT  [1      a     ab     abc]    prefix
 * RIGHT [bcd    cd    d       1 ]    suffix
 * ANS   [bcd    acd   abd    abc]
 */
func productExceptSelf(nums []int) []int {
	res := make([]int, len(nums))

	multiplier := 1
	for i := 0; i < len(nums); i++ {
		res[i] = multiplier
		multiplier *= nums[i]
	}

	multiplier = 1
	for i := len(res) - 1; i >= 0; i-- {
		res[i] *= multiplier
		multiplier *= nums[i]
	}

	return res
}

func productExceptSelfClaude(nums []int) []int {
	n := len(nums)
	res := make([]int, n)
	// forward pass: res[i] = product of all elements to the LEFT of i
	res[0] = 1
	for i := 1; i < n; i++ {
		res[i] = res[i-1] * nums[i-1]
	}
	// backward pass: multiply res[i] by product of all elements to the RIGHT of i
	suffix := 1
	for i := n - 2; i >= 0; i-- {
		suffix *= nums[i+1]
		res[i] *= suffix
	}
	return res
}
