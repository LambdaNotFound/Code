package apidesign

/**
 * 303. Range Sum Query - Immutable
 */
type NumArray struct {
	prefixSum []int
}

func ConstructorNumArray(nums []int) NumArray {
	prefixSum := make([]int, len(nums)+1)
	for i, num := range nums {
		prefixSum[i+1] = prefixSum[i] + num
	}
	return NumArray{prefixSum: prefixSum}
}

func (na *NumArray) SumRange(left int, right int) int {
	return na.prefixSum[right+1] - na.prefixSum[left]
}
