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
