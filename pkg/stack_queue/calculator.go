package stack

import (
	"strconv"
	"strings"
	"unicode"
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
 */
func calculate(s string) int {
	stack := make([]int, 0)
	result, num, sign := 0, 0, 1

	for i := 0; i < len(s); i++ {
		c := s[i]

		switch {
		case c >= '0' && c <= '9':
			// build multi-digit number
			num = num*10 + int(c-'0')

		case c == '+' || c == '-':
			// apply previous number with its sign
			result += sign * num
			num = 0
			if c == '+' {
				sign = 1
			} else {
				sign = -1
			}

		case c == '(':
			// save current state — result and sign before '('
			stack = append(stack, result)
			stack = append(stack, sign)
			// reset for expression inside parentheses
			result = 0
			sign = 1

		case c == ')':
			// apply last number inside parentheses
			result += sign * num
			num = 0

			// restore outer expression state
			outerSign := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			outerResult := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			// combine: outer + outerSign * inner
			result = outerResult + outerSign*result
		}
	}

	// apply last number
	result += sign * num
	return result
}

/**
 * 227. Basic Calculator II
 *
 * Input: s = "3+2*2"
 * Output: 7
 */
func calculate2(s string) int {
	s = strings.ReplaceAll(s, " ", "")
	st := []int{}
	op, num := '+', 0
	for i, c := range s {
		if unicode.IsDigit(c) {
			num = num*10 + (int(c) - '0')
		}
		if i == len(s)-1 || !unicode.IsDigit(c) {
			switch op {
			case '+':
				st = append(st, num)
			case '-':
				st = append(st, -num)
			case '*':
				operand := st[len(st)-1]
				st = st[:len(st)-1]
				st = append(st, operand*num)
			case '/':
				operand := st[len(st)-1]
				st = st[:len(st)-1]
				st = append(st, operand/num)
			}
			op = c
			num = 0
		}
	}

	res := 0
	for _, num := range st {
		res += num
	}
	return res
}

/*
 * expr    = term   (('+' | '-') term)*      ← lowest precedence
 * term    = factor (('*' | '/') factor)*    ← higher precedence
 * factor  = number | '(' expr ')'           ← highest precedence
 */
func calculateRecursiveDescentParser(s string) int {
	pos := 0

	var parseExpr func() int
	var parseTerm func() int
	var parseFactor func() int

	skipSpaces := func() {
		for pos < len(s) && s[pos] == ' ' {
			pos++
		}
	}

	parseFactor = func() int {
		skipSpaces()
		if s[pos] == '(' {
			pos++ // consume '('
			val := parseExpr()
			pos++ // consume ')'
			return val
		}
		// parse number
		num := 0
		for pos < len(s) && s[pos] >= '0' && s[pos] <= '9' {
			num = num*10 + int(s[pos]-'0')
			pos++
		}
		return num
	}

	parseTerm = func() int {
		val := parseFactor()
		skipSpaces()
		for pos < len(s) && (s[pos] == '*' || s[pos] == '/') {
			op := s[pos]
			pos++
			right := parseFactor()
			if op == '*' {
				val *= right
			} else {
				val /= right
			}
			skipSpaces()
		}
		return val
	}

	parseExpr = func() int {
		val := parseTerm()
		skipSpaces()
		for pos < len(s) && (s[pos] == '+' || s[pos] == '-') {
			op := s[pos]
			pos++
			right := parseTerm()
			if op == '+' {
				val += right
			} else {
				val -= right
			}
			skipSpaces()
		}
		return val
	}

	return parseExpr()
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
