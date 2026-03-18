package stack

import (
	"strconv"
	"strings"
)

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

/**
 * 150. Evaluate Reverse Polish Notation
 */
func evalRPN(tokens []string) int {
	if len(tokens) == 1 {
		num, _ := strconv.Atoi(tokens[0])
		return num
	}

	stack := []int{}
	res := 0
	for _, token := range tokens {
		if num, ok := strconv.Atoi(token); ok == nil {
			stack = append(stack, num)
		} else {
			l := len(stack)
			a, b := stack[l-2], stack[l-1]
			stack = stack[:l-2]

			switch token {
			case "+":
				res = a + b
			case "-":
				res = a - b
			case "*":
				res = a * b
			case "/":
				res = a / b
			}

			stack = append(stack, res)
		}
	}
	return res
}

/**
 * 224. Basic Calculator
 *
 * Given a string s representing a valid expression, implement a basic calculator to evaluate it,
 * and return the result of the evaluation.
 *
 * s consists of digits, '+', '-', '(', ')', and ' '.
 */
func calculate(s string) int {
	if len(s) == 0 {
		return 0
	}

	result, sign, num := 0, 1, 0

	var stk []int
	stk = append(stk, sign)

	for i := range s {
		if s[i] >= '0' && s[i] <= '9' {
			num = num*10 + int(s[i]-'0')
		} else if s[i] == '+' || s[i] == '-' {
			result += sign * num
			sign = stk[len(stk)-1]
			if s[i] != '+' {
				sign *= -1
			}
			num = 0
		} else if s[i] == '(' {
			stk = append(stk, sign)
		} else if s[i] == ')' {
			stk = stk[:len(stk)-1]
		}
	}

	result += sign * num
	return result
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
