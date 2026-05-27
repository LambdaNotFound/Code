package stack

import (
	"strconv"
	"strings"
)

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

/**
 * 1209. Remove All Adjacent Duplicates in String II
 *
 * a k duplicate removal consists of choosing k adjacent and equal letters from s
 *
 * We repeatedly make k duplicate removals on s until we no longer can.
 */
func removeDuplicates(s string, k int) string {
	stack := []rune{}
	for _, r := range s {
		stack = append(stack, r) // must push first, ("aaaa", 2) => ""
		if len(stack) >= k {
			i := len(stack) - k
			// check if last k runes from the stack are the same
			for i < len(stack) && stack[i] == stack[len(stack)-1] {
				i++
			}
			if i == len(stack) {
				stack = stack[:len(stack)-k]
			}
		}
	}
	return string(stack)
}

func removeDuplicatesClaude(s string, k int) string {
	type pair struct {
		val  byte
		freq int
	}

	stk := []pair{}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if len(stk) > 0 && stk[len(stk)-1].val == c {
			stk[len(stk)-1].freq++
		} else {
			stk = append(stk, pair{c, 1})
		}

		// pop immediately when freq reaches k
		if stk[len(stk)-1].freq == k {
			stk = stk[:len(stk)-1]
		}
	}

	var sb strings.Builder
	for _, p := range stk {
		sb.WriteString(strings.Repeat(string(p.val), p.freq))
	}
	return sb.String()
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
