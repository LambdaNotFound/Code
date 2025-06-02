package solid_coding

import (
	"math"
	"strconv"
	"strings"
	"unicode"
)

/**
 *    for i, rune := range string {
 *
 *    for i := 0; i < len(string); i++ {
 *        byte := string[i]
 *
 */

/**
 * 8. String to Integer (atoi)
 */
func myAtoi(s string) int {
    s = strings.TrimSpace(s)
    multiplier := 1
    if len(s) == 0 {
        return 0
    } else if s[0] == '-' {
        multiplier = -1
        s = s[1:]
    } else if s[0] == '+' {
        s = s[1:]
    }

    res := 0
    for _, r := range s {
        if !unicode.IsDigit(r) {
            break
        }
        curr, _ := strconv.Atoi(string(r))

        if multiplier == 1 && (res*10 > math.MaxInt32-curr) {
            return math.MaxInt32
        }
        if multiplier == -1 && (-res*10 < math.MinInt32+curr) {
            return math.MinInt32
        }

        res = res*10 + curr
    }
    return multiplier * res
}

/**
 * 20. Valid Parentheses
 */
func isValid(s string) bool {
    stack := []byte{}
    for i := 0; i < len(s); i++ {
        if s[i] == '(' || s[i] == '[' || s[i] == '{' {
            stack = append(stack, s[i])
        } else {
            l := len(stack)
            if l == 0 {
                return false
            } else {
                var expected byte
                switch s[i] {
                case ')':
                    expected = '('
                case ']':
                    expected = '['
                case '}':
                    expected = '{'
                }
                if expected != stack[l-1] {
                    return false
                }
            }
            stack = stack[:l-1]
        }

    }

    return len(stack) == 0
}

func isValid_lookup(s string) bool {
    stack := []rune{} // Stack for opening brackets
    hash := map[rune]rune{')': '(', ']': '[', '}': '{'}

    for _, char := range s {
        if match, found := hash[char]; found {
            // Check if stack is non-empty and matches
            if len(stack) > 0 && stack[len(stack)-1] == match {
                stack = stack[:len(stack)-1] // Pop
            } else {
                return false // Invalid
            }
        } else {
            stack = append(stack, char) // Push opening bracket
        }
    }
    return len(stack) == 0 // Valid if stack is empty
}

/**
 * 67. Add Binary
 */
func addBinary(a string, b string) string {
    finalstr := ""
    v1, v2, rem := 0, 0, 0

    for l1, l2 := len(a)-1, len(b)-1; l1 >= 0 || l2 >= 0 || rem != 0; {
        if l1 >= 0 {
            v1, _ = strconv.Atoi(string(a[l1]))
        }
        if l2 >= 0 {
            v2, _ = strconv.Atoi(string(b[l2]))
        }

        sum := v1 + v2 + rem

        // according to sum append appropriate character in finalstr
        switch sum {
        case 3:
            finalstr = "1" + finalstr
            rem = 1
        case 2:
            finalstr = "0" + finalstr
            rem = 1
        case 1:
            finalstr = "1" + finalstr
            rem = 0
        case 0:
            finalstr = "0" + finalstr
            rem = 0
        }

        v1, v2 = 0, 0
        l1 -= 1
        l2 -= 1
    }

    return finalstr
}
