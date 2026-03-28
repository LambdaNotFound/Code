package stack

import (
	"strconv"
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
