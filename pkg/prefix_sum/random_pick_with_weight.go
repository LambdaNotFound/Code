package prefixsum

import (
	"math/rand"
	"sort"
)

/**
 * 528. Random Pick with Weight
 *
 * prefixSum + binarySearch
 */
type Solution struct {
	prefixSum []int
}

func Constructor(w []int) Solution {
	prefixSum := make([]int, len(w)+1)
	for i := 1; i < len(prefixSum); i++ {
		prefixSum[i] = prefixSum[i-1] + w[i-1]
	}

	return Solution{prefixSum: prefixSum}
}

func (s *Solution) PickIndex() int {
	r := rand.Intn(s.prefixSum[len(s.prefixSum)-1])
	return upperBound(s.prefixSum, r) - 1
}

func upperBound(array []int, target int) int {
	return sort.Search(len(array), func(i int) bool {
		return target < array[i] // lower bound, target <= array[i]
	})
}
