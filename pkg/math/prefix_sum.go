package math

/**
 * 238. Product of Array Except Self
 *
 * [a, b, c, d]
 *
 * LEFT  [1     a    ab    abc]
 * RIGHT [bcd  cd     d     1 ]
 * ANS  [bcd   acd   abd   abc]
 *
 */
func productExceptSelf(nums []int) []int {
    l := len(nums)
    left := make([]int, l)
    right := make([]int, l)
    ans := make([]int, l)
    left[0] = 1
    right[l-1] = 1
    for i := 1; i < l; i++ {
        j := l - i - 1
        left[i] = nums[i-1] * left[i-1]
        right[j] = nums[j+1] * right[j+1]
    }
    for i := 0; i < l; i++ {
        ans[i] = left[i] * right[i]
    }
    return ans
}

func productExceptSelfWithMultiplier(nums []int) []int {
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
