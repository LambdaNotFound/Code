package stack

import (
	"regexp"
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
 * stack stores [..., result,  sign] before each '('
 *    when process ')', res = outerResult + outerSign * result
 */
func calculate(s string) int {
	reg := regexp.MustCompile(`[^0-9+\-()]`)
	s = reg.ReplaceAllString(s, "")

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
	reg := regexp.MustCompile(`[^0-9+\-*/]`)
	s = reg.ReplaceAllString(s, "")

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
 * 227. Basic Calculator II variant
 *
 * Addition, Subtraction, Multiplication, Division
 *
 * Input: s = "3add2mul2"
 * Output: 7
 */
func calculateIIVariant(s string) int {
	s = strings.ReplaceAll(s, " ", "")

	stack := []int{}
	num, preOp, curOp := 0, "add", ""
	for i, ch := range s {
		if ch >= '0' && ch <= '9' {
			num = num*10 + int(ch-'0')
		} else {
			curOp += string(ch)
		}

		if i == len(s)-1 || curOp == "add" || curOp == "sub" || curOp == "mul" || curOp == "div" {
			switch {
			case preOp == "add":
				stack = append(stack, num)
			case preOp == "sub":
				stack = append(stack, -num)
			case preOp == "mul":
				stack[len(stack)-1] = stack[len(stack)-1] * num
			case preOp == "div":
				stack[len(stack)-1] = stack[len(stack)-1] / num
			}

			num = 0
			preOp = curOp
			curOp = ""
		}
	}

	res := 0
	for i := range stack {
		res = res + stack[i]
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
func calculateIII(s string) int {
	stack := []int{}
	num, op := 0, '+'
	for i := 0; i < len(s); i++ {
		c := s[i]

		if c == '(' {
			count, next := 1, i+1
			for next < len(s) {
				if s[next] == '(' {
					count++
				}
				if s[next] == ')' {
					count--
				}
				if count == 0 {
					break
				}
				next++
			}
			num = calculateIII(s[i+1 : next])
			i = next
		}

		if c >= '0' && c <= '9' {
			num = num*10 + int(c-'0')
		}
		if c == '+' || c == '-' || c == '*' || c == '/' || i == len(s)-1 {
			switch op {
			case '+':
				stack = append(stack, num)
			case '-':
				stack = append(stack, -num)
			case '*':
				stack[len(stack)-1] *= num
			case '/':
				stack[len(stack)-1] /= num
			}
			num = 0
			op = rune(c)
		}
	}

	res := 0
	for _, v := range stack {
		res += v
	}
	return res
}

func calculateClaude(s string) int {
	reg := regexp.MustCompile(`[^0-9+\-*/%()]`)
	s = reg.ReplaceAllString(s, "")

	pos := 0

	var parse func() int
	parse = func() int {
		stack := []int{}
		num, op := 0, byte('+')

		for pos < len(s) {
			ch := s[pos]
			pos++

			if ch >= '0' && ch <= '9' {
				num = num*10 + int(ch-'0')
			}

			if ch == '(' {
				num = parse() // recurse into parens
			}

			if (ch != ' ' && ch < '0') || pos == len(s) {
				switch op {
				case '+':
					stack = append(stack, num)
				case '-':
					stack = append(stack, -num)
				case '*':
					stack[len(stack)-1] *= num
				case '/':
					stack[len(stack)-1] /= num
				}
				num, op = 0, ch
			}

			if ch == ')' {
				break // return to caller
			}
		}

		result := 0
		for _, v := range stack {
			result += v
		}
		return result
	}

	return parse()
}

// Basic Calculator template
func calculateT(s string) int {
	pos := 0

	var dfs func() int
	dfs = func() int {
		stk := []int{}
		num, op := 0, '+'

		for pos < len(s) {
			ch := s[pos]
			pos++

			if ch >= '0' && ch <= '9' {
				num = num*10 + int(ch-'0')
			} else if ch == '(' { // ← else if: mutually exclusive
				num = dfs()
			}

			// commit on operator, end, or closing paren
			if pos == len(s) || ch == '+' || ch == '-' || ch == '*' || ch == '/' || ch == ')' {
				switch op {
				case '+':
					stk = append(stk, num)
				case '-':
					stk = append(stk, -num)
				case '*':
					stk[len(stk)-1] *= num
				case '/':
					stk[len(stk)-1] /= num
				}
				num = 0
				if ch == '+' || ch == '-' || ch == '*' || ch == '/' {
					op = rune(ch) // ← only update for real operators
				}
				if ch == ')' {
					break
				}
			}
		}

		res := 0
		for _, v := range stk {
			res += v
		}
		return res
	}

	return dfs()
}

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
