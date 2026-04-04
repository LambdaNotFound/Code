package stack

import (
	"strconv"
	"strings"
)

/**
 * 224. Basic Calculator
 *
 * Given a string s representing a valid expression, implement a basic calculator to evaluate it,
 * and return the result of the evaluation.
 *
 * s consists of digits, '+', '-', '(', ')', and ' '.
 *
 * Input: s = "1 + 1"
 * Output: 2
 * Input: s = " 2-1 + 2 "
 * Output: 3
 * Input: s = "(1+(4+5+2)-3)+(6+8)"
 * Output: 23
 *
 * stack stores stores result + sign before each '('
 */
func calculate(s string) int {
	stack := []int{} // stores result + sign before each '('
	result, num, sign := 0, 0, 1
	for i, ch := range s {
		switch {
		case ch >= '0' && ch <= '9':
			num = num*10 + int(ch-'0')

		case ch == '+' || ch == '-':
			result += sign * num // commit current number
			num = 0
			if ch == '+' {
				sign = 1
			} else {
				sign = -1
			}

		case ch == '(':
			// save current state, start fresh inside parens
			stack = append(stack, result, sign)
			result, sign = 0, 1

		case ch == ')':
			result += sign * num // commit last number in parens
			num = 0
			// restore outer context
			outerSign := stack[len(stack)-1]
			outerResult := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			result = outerResult + outerSign*result
		}
		if i == len(s)-1 {
			result += sign * num
		}
	}
	return result
}

/**
 * 227. Basic Calculator II
 *
 * Input: s = "3+2*2"
 * Output: 7
 */
func calculateII(s string) int {
	s = strings.Trim(s, " ")
	stack := []int{}
	num, op := 0, '+'
	for i, c := range s {
		if c >= '0' && c <= '9' {
			num = num*10 + int(c-'0')
		}
		if i == len(s)-1 || c == '+' || c == '-' || c == '*' || c == '/' {
			switch op {
			case '+':
				stack = append(stack, num)
			case '-':
				stack = append(stack, -num)
			case '*':
				top := stack[len(stack)-1]
				stack[len(stack)-1] = top * num
			case '/':
				top := stack[len(stack)-1]
				stack[len(stack)-1] = top / num
			}
			num = 0
			op = c
		}
	}

	res := 0
	for _, num := range stack {
		res += num
	}
	return res
}

/**
 * 772. Basic Calculator III
 *
 * Input:  "1+1"
 * Output: 2
 *
 * Input:  "6-4/2"
 * Output: 4
 *
 * Input:  "2*(5+5*2)/3+(6/2+8)"
 * Output: 21
 *
 * Input:  "(2+6*3+5-(3*14/7+2)*5)+3"
 * Output: -12
 */

/**
 * 394. Decode String
 *
 * Input: s = "3[a]2[bc]"
 * Output: "aaabcbc"
 *
 * Input: s = "3[a2[c]]"
 * Output: "accaccacc"
 *
 * Input: s = "2[abc]3[cd]ef"
 * Output: "abcabccdcdcdef"
 */
func decodeString(s string) string {
	// stack stores either string segments or digits
	stack := make([]string, 0)

	for _, c := range s {
		if c != ']' {
			stack = append(stack, string(c))
		} else {
			// pop until '[' to get the substring
			substr := ""
			for stack[len(stack)-1] != "[" {
				substr = stack[len(stack)-1] + substr
				stack = stack[:len(stack)-1]
			}
			stack = stack[:len(stack)-1] // pop '['

			// pop digits to get the repeat count
			k := ""
			for len(stack) > 0 && stack[len(stack)-1] >= "0" && stack[len(stack)-1] <= "9" {
				k = stack[len(stack)-1] + k
				stack = stack[:len(stack)-1]
			}
			num, _ := strconv.Atoi(k)

			// push expanded string back onto stack
			stack = append(stack, strings.Repeat(substr, num))
		}
	}

	return strings.Join(stack, "")
}
