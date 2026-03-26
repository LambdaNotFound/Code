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
	stack := make([]int, 0)

	for _, token := range tokens {
		val, err := strconv.Atoi(token)
		if err == nil {
			stack = append(stack, val)
			continue
		}

		// pop two operands — right before left
		right := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		left := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		switch token {
		case "+":
			stack = append(stack, left+right)
		case "-":
			stack = append(stack, left-right)
		case "*":
			stack = append(stack, left*right)
		case "/":
			stack = append(stack, left/right)
		}
	}

	return stack[0]
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
 * 844. Backspace String Compare
 */
func backspaceCompare(s string, t string) bool {
	process := func(str string) string {
		stack := make([]byte, 0)
		for i := range str {
			if str[i] == '#' {
				if len(stack) > 0 {
					stack = stack[:len(stack)-1]
				}
			} else {
				stack = append(stack, str[i])
			}
		}
		return string(stack)
	}

	return process(s) == process(t)
}

/**
 * 735. Asteroid Collision
 */
func asteroidCollision(asteroids []int) []int {
	st := []int{}
	for i := 0; i < len(asteroids); {
		push := true
		ast := asteroids[i]
		if len(st) > 0 {
			top := st[len(st)-1]
			if top > 0 && ast < 0 {
				push = false
				if top+ast == 0 {
					st = st[:len(st)-1]
					i++
				} else if top+ast > 0 {
					i++
				} else { // if st[len(st)-1] + ast < 0
					st = st[:len(st)-1]
				}
			}
		}
		if push {
			st = append(st, ast)
			i++
		}
	}
	return st
}

func asteroidCollisionCalude(asteroids []int) []int {
	st := []int{}
	for _, asteroid := range asteroids {
		alive := true
		for alive && asteroid < 0 && len(st) > 0 && st[len(st)-1] > 0 {
			top := st[len(st)-1]
			if top < -asteroid { // asteroid wins, pop and keep fighting
				st = st[:len(st)-1]
			} else if top == -asteroid { // mutual destruction
				st = st[:len(st)-1]
				alive = false
			} else { // top wins, asteroid dies
				alive = false
			}
		}
		if alive {
			st = append(st, asteroid)
		}
	}
	return st
}
