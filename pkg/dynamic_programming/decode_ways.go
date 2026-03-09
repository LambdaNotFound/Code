package dynamic_programming

import "strconv"

/**
 * 91. Decode Ways
 */
func numDecodings(s string) int {
	table := make([]int, len(s)+1)
	table[0] = 1
	if s[0] != '0' {
		table[1] = 1
	}
	for i := 2; i <= len(s); i++ {
		if num, _ := strconv.Atoi(s[i-2 : i]); num >= 10 && num <= 26 {
			table[i] += table[i-2]
		}
		if num, _ := strconv.Atoi(s[i-1 : i]); num < 10 && num > 0 {
			table[i] += table[i-1]
		}
	}

	return table[len(s)]
}
