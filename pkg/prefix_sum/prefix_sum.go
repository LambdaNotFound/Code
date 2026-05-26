package prefixsum

import (
	. "gocode/types"
)

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
	n := len(nums)
	prefix := make([]int, n)
	suffix := make([]int, n)
	prefix[0] = 1
	for i := 1; i < n; i++ {
		prefix[i] = prefix[i-1] * nums[i-1]
	}
	suffix[n-1] = 1
	for i := n - 2; i >= 0; i-- {
		suffix[i] = suffix[i+1] * nums[i+1]
	}
	result := make([]int, n)
	for i := range result {
		result[i] = prefix[i] * suffix[i]
	}
	return result
}

/**
 * 525. Contiguous Array
 */
func findMaxLength(nums []int) int {
	// firstSeen[prefixSum] = earliest index where this prefix sum occurred
	firstSeen := map[int]int{0: -1} // base case: sum=0 seen before index 0

	prefixSum := 0
	maxLen := 0

	for i, num := range nums {
		// transform: 0 → -1, 1 → +1
		if num == 0 {
			prefixSum--
		} else {
			prefixSum++
		}

		if firstIdx, seen := firstSeen[prefixSum]; seen {
			// same prefix sum seen before → subarray [firstIdx+1..i] has sum 0
			maxLen = max(maxLen, i-firstIdx)
		} else {
			// first time seeing this prefix sum — record earliest index
			firstSeen[prefixSum] = i
		}
	}

	return maxLen
}

/**
 * 560. Subarray Sum Equals K
 */
func subarraySum(nums []int, k int) int {
	n := len(nums)
	prefixSum := make([]int, n+1)
	for i, num := range nums {
		prefixSum[i+1] = prefixSum[i] + num
	}

	count := 0
	for i := 0; i < n; i++ {
		for j := i + 1; j <= n; j++ { // check all pairs (i, j) — O(n²)
			if prefixSum[j]-prefixSum[i] == k {
				count++
			}
		}
	}

	return count
}

func subarraySumWithHashmap(nums []int, k int) int {
	// prefixSumFreq[sum] = freq of this prefix sum has been seen
	prefixSumFreq := map[int]int{0: 1} // base case: empty prefix has sum 0
	prefixSum := 0
	count := 0

	// Instead of checking all pairs, use a hashmap to look up
	// the complementc in O(1):
	// prefixSum[j] - prefixSum[i] = k
	// => prefixSum[i] = prefixSum[j] - k
	for _, num := range nums {
		prefixSum += num

		// how many previous prefixSums equal prefixSum-k?
		// each one forms a valid subarray ending at current index
		count += prefixSumFreq[prefixSum-k]

		// record current prefix sum
		prefixSumFreq[prefixSum]++
	}

	return count
}

/**
 * 930. Binary Subarrays With Sum
 *
 * count subarrays ending at j with sum == goal
 * [0, ..., i]            prefixSum[i]
 * [0, ..., ..., j]       prefixSum[j]
 *            delta       subarray's sum
 * prefixSum[i] + prefixSum[j] == goal
 *                ^ should be added to map AFTER lookup
 *                countMap[prefixSum]++
 */
func numSubarraysWithSum(nums []int, goal int) int {
	countMap, prefixSum := map[int]int{0: 1}, 0
	count := 0
	for _, num := range nums {
		prefixSum += num

		count += countMap[prefixSum-goal]
		countMap[prefixSum]++
	}
	return count
}

func numSubarraysWithSumPrefixSum(nums []int, goal int) int {
	n := len(nums)
	prefixSum := make([]int, n+1)
	for i := 0; i < n; i++ {
		prefixSum[i+1] = prefixSum[i] + nums[i]
	}

	count := 0
	for start := 0; start <= n; start++ {
		for end := start + 1; end <= n; end++ {
			if prefixSum[end]-prefixSum[start] == goal {
				count++
			}
		}
	}
	return count
}

/**
 * 523. Continuous Subarray Sum
 *
 * its length is at least two, and the sum of the elements of the subarray is a multiple of k.
 *
 * Math: (a - b) % k == 0  ↔  a % k == b % k
 */
func checkSubarraySum(nums []int, k int) bool {
	prefixSum := 0
	remainderToIndex := map[int]int{0: -1} // remainder -> earliest index where this remainder was seen

	for i, num := range nums {
		prefixSum += num
		remainder := prefixSum % k
		if firstIdx, found := remainderToIndex[remainder]; found {
			if i-firstIdx >= 2 { // subarray nums[firstIdx+1 .. i] has length i - firstIdx
				return true
			}
		} else {
			remainderToIndex[remainder] = i
		}
	}
	return false
}

func checkSubarraySumNaive(nums []int, k int) bool {
	n := len(nums)
	prefixSum := make([]int, n+1)
	for i := 1; i <= n; i++ {
		prefixSum[i] = prefixSum[i-1] + nums[i-1]
	}

	for i := 0; i <= n; i++ {
		for j := i + 2; j <= n; j++ {
			if (prefixSum[j]-prefixSum[i])%k == 0 {
				return true
			}
		}
	}
	return false
}
