package stack

import "strings"

/**
 * 32. Longest Valid Parentheses
 *    track the invalid indices
 */
func longestValidParentheses(s string) int {
	stack := []int{-1} // sentinel: boundary before string starts
	res := 0

	for i, ch := range s {
		if ch == '(' {
			stack = append(stack, i) // push index of '('
		} else {
			stack = stack[:len(stack)-1] // pop matching '('

			if len(stack) == 0 {
				stack = append(stack, i) // new boundary: unmatched ')'
			} else {
				res = max(res, i-stack[len(stack)-1]) // valid length
			}
		}
	}
	return res
}

/*
 * Case 1: s[i-1] == '('   → direct match "()"
 *	       dp[i] = dp[i-2] + 2
 *
 * Case 2: s[i-1] == ')'   → previous char closes its own group "))"
 *	need to check if char BEFORE that group matches
 *	j = i - dp[i-1] - 1
 *	if s[j] == '(' :
 *	    dp[i] = dp[i-1] + 2 + dp[j-1]
 */
func longestValidParenthesesDP(s string) int {
	n := len(s)
	dp := make([]int, n)
	res := 0

	for i := 1; i < n; i++ {
		if s[i] == ')' {
			if s[i-1] == '(' {
				// case 1: direct "()" match
				if i >= 2 {
					dp[i] = dp[i-2] + 2
				} else {
					dp[i] = 2
				}
			} else if dp[i-1] > 0 {
				// case 2: ends with "))" — look before inner group
				j := i - dp[i-1] - 1
				if j >= 0 && s[j] == '(' {
					dp[i] = dp[i-1] + 2
					if j > 0 {
						dp[i] += dp[j-1]
					}
				}
			}
			res = max(res, dp[i])
		}
	}
	return res
}

func longestValidParenthesesTwoPointers(s string) int {
	res := 0

	// pass 1: left → right
	open, close := 0, 0
	for i := 0; i < len(s); i++ {
		if s[i] == '(' {
			open++
		} else {
			close++
		}
		if open == close {
			res = max(res, 2*close)
		} else if close > open {
			open, close = 0, 0 // unmatched ')' — reset
		}
	}

	// pass 2: right → left
	open, close = 0, 0
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '(' {
			open++
		} else {
			close++
		}
		if open == close {
			res = max(res, 2*open)
		} else if open > close {
			open, close = 0, 0 // unmatched '(' — reset
		}
	}

	return res
}

/**
 * 1249. Minimum Remove to Make Valid Parentheses
 */
func minRemoveToMakeValid(s string) string {
	var stack []int
	chars := []rune(s)

	// first pass: track indices of '(' and mark ')' for removal
	for i, char := range chars {
		if char == '(' {
			stack = append(stack, i)
		} else if char == ')' {
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			} else {
				chars[i] = '*' // mark ')' for removal
			}
		}
	}

	// mark remaining '(' for removal
	for _, index := range stack {
		chars[index] = '*' // mark '(' for removal
	}

	result := make([]rune, 0, len(chars))
	for _, char := range chars {
		if char != '*' {
			result = append(result, char)
		}
	}

	return string(result)
}

func minRemoveToMakeValidClaude(s string) string {
	// Step 1: find all unmatched bracket indices
	stack := make([]int, 0) // indices of unmatched brackets
	for i := 0; i < len(s); i++ {
		if s[i] == '(' {
			stack = append(stack, i) // candidate for removal
		} else if s[i] == ')' {
			if len(stack) > 0 && s[stack[len(stack)-1]] == '(' {
				stack = stack[:len(stack)-1] // matched — pop
			} else {
				stack = append(stack, i) // unmatched ')' — mark for removal
			}
		}
	}

	// Step 2: build result skipping unmatched indices
	// convert stack to set for O(1) lookup
	remove := make(map[int]bool, len(stack))
	for _, idx := range stack {
		remove[idx] = true
	}

	var sb strings.Builder
	for i := 0; i < len(s); i++ {
		if !remove[i] {
			sb.WriteByte(s[i])
		}
	}

	return sb.String()
}

/**
 * 921. Minimum Add to Make Parentheses Valid
 *
 * "()))((" => 4
 */
func minAddToMakeValid(s string) int {
	unmatchedOpen := 0  // '(' seen with no matching ')' yet
	unmatchedClose := 0 // ')' seen with no matching '(' available

	for _, ch := range s {
		if ch == '(' {
			unmatchedOpen++
		} else if ch == ')' {
			if unmatchedOpen > 0 {
				unmatchedOpen-- // matched with a prior '('
			} else {
				unmatchedClose++ // no '(' available to match
			}
		}
	}
	// each unmatched '(' needs a ')' added, each unmatched ')' needs a '(' added
	return unmatchedOpen + unmatchedClose
}

/**
 * 20. Valid Parentheses
 */
func isValid(s string) bool {
	stack := make([]byte, 0)
	for i := range s {
		switch s[i] {
		case '(':
			stack = append(stack, ')')
		case '[':
			stack = append(stack, ']')
		case '{':
			stack = append(stack, '}')
		default: // closing bracket
			if len(stack) == 0 || stack[len(stack)-1] != s[i] {
				return false
			}
			stack = stack[:len(stack)-1]
		}
	}
	return len(stack) == 0
}

func isValidClaude(s string) bool {
	// map closing bracket to its expected opening bracket
	match := map[byte]byte{
		')': '(',
		']': '[',
		'}': '{',
	}

	stack := make([]byte, 0)
	for i := range s {
		c := s[i]
		if expected, isClosing := match[c]; isClosing {
			// closing bracket — check top of stack matches
			if len(stack) == 0 || stack[len(stack)-1] != expected {
				return false
			}
			stack = stack[:len(stack)-1]
		} else {
			// opening bracket — push to stack
			stack = append(stack, c)
		}
	}

	return len(stack) == 0
}
