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

func flat(arr []any, n int) []any {
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
