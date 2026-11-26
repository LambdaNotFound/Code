package binarysearch

/**
 * Find all elements in a sorted array appearing more than n/3 times.
 *
 * Follow-up: solve in better than O(n) time.
 *
 * At most 2 elements can appear more than n/3 times (pigeonhole).
 */
func majorityElements(arr []int) []int {
	n := len(arr)
	threshold := n / 3

	// Extract at most 2 candidates from anchor positions.
	// Any element with freq > n/3 must cover index n/3 or 2n/3.
	candidates := []int{arr[n/3], arr[2*n/3]}

	// Deduplicate candidates — both anchors may land on the same value.
	if candidates[0] == candidates[1] {
		candidates = candidates[:1]
	}

	// lowerBound returns the first index where arr[i] >= val.
	lowerBound := func(arr []int, val int) int {
		left, right := 0, len(arr)
		for left < right {
			mid := left + (right-left)/2
			if arr[mid] < val {
				left = mid + 1
			} else {
				right = mid
			}
		}
		return left
	}

	upperBound := func(arr []int, val int) int {
		left, right := 0, len(arr)
		for left < right {
			mid := left + (right-left)/2
			if arr[mid] <= val {
				left = mid + 1
			} else {
				right = mid
			}
		}
		return left
	}

	count := func(arr []int, val int) int {
		return upperBound(arr, val) - lowerBound(arr, val)
	}

	var result []int
	for _, candidate := range candidates {
		if count(arr, candidate) > threshold {
			result = append(result, candidate)
		}
	}
	return result
}
