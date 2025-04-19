package solid_coding

import "strconv"

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
