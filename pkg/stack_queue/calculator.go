package stack

import (
	"regexp"
	"strings"
)

// Basic Calculator template
func calculateTemplate(s string) int {
	s = strings.ReplaceAll(s, " ", "")
	pos := 0

	var dfs func() int
	dfs = func() int {
		stack := []int{}
		currentNumber, pendingSign := 0, '+'

		for pos < len(s) {
			char := s[pos]
			pos++

			if char >= '0' && char <= '9' {
				currentNumber = currentNumber*10 + int(char-'0')
			}

			if char == '(' {
				currentNumber = dfs()
			}

			// commit on operator, end, or closing paren
			if pos == len(s) || char == '+' || char == '-' || char == '*' || char == '/' || char == ')' {
				switch pendingSign {
				case '+':
					stack = append(stack, currentNumber)
				case '-':
					stack = append(stack, -currentNumber)
				case '*':
					stack[len(stack)-1] *= currentNumber
				case '/':
					stack[len(stack)-1] /= currentNumber
				}

				if char == ')' {
					break
				}
				currentNumber = 0
				pendingSign = rune(char) // ← only update for real operators
			}
		}

		result := 0
		for _, num := range stack {
			result += num
		}
		return result
	}

	return dfs()
}

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
	whitespace := regexp.MustCompile(`[^0-9+\-()]`)
	s = whitespace.ReplaceAllString(s, "")

	pos := 0
	var parse func() int
	parse = func() int {
		stack := []int{}
		currentNumber, pendingSign := 0, 1
		for pos < len(s) {
			char := s[pos]
			pos++

			if char >= '0' && char <= '9' {
				currentNumber = currentNumber*10 + int(char-'0')
			}
			if char == '(' {
				currentNumber = parse()
			}
			if pos == len(s) || char == '+' || char == '-' || char == ')' {
				stack = append(stack, pendingSign*currentNumber)
				if char == ')' {
					break
				}
				currentNumber = 0
				switch char {
				case '+':
					pendingSign = 1
				case '-':
					pendingSign = -1
				}
			}
		}

		total := 0
		for _, val := range stack {
			total += val
		}
		return total
	}

	return parse()
}

func calculateIterative(s string) int {
	reg := regexp.MustCompile(`[^0-9+\-()]`)
	s = reg.ReplaceAllString(s, "")

	stack := []int{} // stores result + sign before each '('
	result, currentNumber, pendingSign := 0, 0, 1
	for _, char := range s {
		if char >= '0' && char <= '9' {
			currentNumber = currentNumber*10 + int(char-'0')
		} else if char == '+' || char == '-' {
			result += pendingSign * currentNumber // commit current number
			currentNumber = 0
			if char == '+' {
				pendingSign = 1
			} else {
				pendingSign = -1
			}
		} else if char == '(' {
			stack = append(stack, result) // save current number
			stack = append(stack, pendingSign)
			result, pendingSign = 0, 1 // start fresh inside parens
		} else if char == ')' {
			result += pendingSign * currentNumber // commit last number in parens
			currentNumber = 0

			outerSign := stack[len(stack)-1] // restore outer context
			stack = stack[:len(stack)-1]
			outerResult := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			result = outerResult + outerSign*result
		}

	}
	result += pendingSign * currentNumber

	return result
}

/**
 * 227. Basic Calculator II
 *
 * Input: s = "3+2*2"
 * Output: 7
 */
func calculate2(s string) int {
	reg := regexp.MustCompile(`[^0-9+\-*/]`)
	s = reg.ReplaceAllString(s, "")

	stack := []int{}
	currentNumber, pendingSign := 0, '+'
	for i, char := range s {
		if char >= '0' && char <= '9' {
			currentNumber = currentNumber*10 + int(char-'0')
		}
		if i == len(s)-1 || char == '+' || char == '-' || char == '*' || char == '/' {
			switch pendingSign {
			case '+':
				stack = append(stack, currentNumber)
			case '-':
				stack = append(stack, -currentNumber)
			case '*':
				top := stack[len(stack)-1]
				stack[len(stack)-1] = top * currentNumber
			case '/':
				top := stack[len(stack)-1]
				stack[len(stack)-1] = top / currentNumber
			}
			currentNumber = 0
			pendingSign = char
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
func calculate2Variant(s string) int {
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
func calculate3(s string) int {
	reg := regexp.MustCompile(`[^0-9+\-*/%()]`)
	s = reg.ReplaceAllString(s, "")

	pos := 0

	var parse func() int
	parse = func() int {
		stack := []int{}
		currentNumber, pendingSign := 0, byte('+')

		for pos < len(s) {
			ch := s[pos]
			pos++

			if ch >= '0' && ch <= '9' {
				currentNumber = currentNumber*10 + int(ch-'0')
			}

			if ch == '(' {
				currentNumber = parse() // recurse into parens
			}

			isOperator := ch == '+' || ch == '-' || ch == '*' || ch == '/' || ch == ')'
			if pos == len(s) || isOperator {
				switch pendingSign {
				case '+':
					stack = append(stack, currentNumber)
				case '-':
					stack = append(stack, -currentNumber)
				case '*':
					stack[len(stack)-1] *= currentNumber
				case '/':
					stack[len(stack)-1] /= currentNumber
				}

				if ch == ')' {
					break // commit and return to caller
				}
				currentNumber, pendingSign = 0, ch
			}

		}

		result := 0
		for _, num := range stack {
			result += num
		}
		return result
	}

	return parse()
}
