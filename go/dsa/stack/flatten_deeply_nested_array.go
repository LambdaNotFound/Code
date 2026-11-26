package stack

import (
	"fmt"
	"strconv"
	"strings"
)

/**
 * 2625. Flatten Deeply Nested Array
 * Given a multi-dimensional array arr and a depth n, return a flattened version of that array.
 *
 * Approach: serialize to tokens ("[", "]", ints), then recursive descent —
 * when n > 0 and we see "[", recurse with n-1 and spread; otherwise keep as element.
 */

/**
 * Example 1:
 *
 * Input
 * arr = [1, 2, 3, [4, 5, 6], [7, 8, [9, 10, 11], 12], [13, 14, 15]]
 * n = 0
 * Output
 * [1, 2, 3, [4, 5, 6], [7, 8, [9, 10, 11], 12], [13, 14, 15]]
 *
 * Explanation
 * Passing a depth of n=0 will always result in the original array. This is because the smallest possible depth of a subarray (0) is not less than n=0. Thus, no subarray should be flattened.
 * Example 2:
 *
 * Input
 * arr = [1, 2, 3, [4, 5, 6], [7, 8, [9, 10, 11], 12], [13, 14, 15]]
 * n = 1
 * Output
 * [1, 2, 3, 4, 5, 6, 7, 8, [9, 10, 11], 12, 13, 14, 15]
 *
 * Explanation
 * The subarrays starting with 4, 7, and 13 are all flattened. This is because their depth of 0 is less than 1. However [9, 10, 11] remains unflattened because its depth is 1.
 * Example 3:
 *
 * Input
 * arr = [[1, 2, 3], [4, 5, 6], [7, 8, [9, 10, 11], 12], [13, 14, 15]]
 * n = 2
 * Output
 * [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15]
 *
 * Explanation
 * The maximum depth of any subarray is 1. Thus, all of them are flattened.
 */

func flat(arr []any, n int) []any { // nested array, type has to be []any
	s := fmt.Sprintf("%v", arr) // formats arr using Go's default format verb %v and returns it as a string.
	s = strings.ReplaceAll(s, "[", " [ ")
	s = strings.ReplaceAll(s, "]", " ] ")
	tokens := strings.Fields(s)

	pos := 0
	var parseArr func(n int) []any
	parseArr = func(n int) []any {
		pos++ // consume '['
		result := []any{}
		for tokens[pos] != "]" {
			if tokens[pos] == "[" {
				inner := parseArr(n - 1)
				if n > 0 {
					result = append(result, inner...) // expand
				} else {
					result = append(result, inner) // preserve
				}
			} else {
				num, _ := strconv.Atoi(tokens[pos])
				result = append(result, num)
				pos++
			}
		}
		pos++ // consume ']'
		return result
	}

	return parseArr(n)
}
