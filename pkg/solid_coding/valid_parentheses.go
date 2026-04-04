package solid_coding

/**
 * 20. Valid Parentheses
 *
 * Stack-based implementation (push expected closing bracket).
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
		default:
			if len(stack) == 0 || stack[len(stack)-1] != s[i] {
				return false
			}
			stack = stack[:len(stack)-1]
		}
	}
	return len(stack) == 0
}

/**
 * isValid_lookup uses a closing→opening map to match brackets.
 */
func isValid_lookup(s string) bool {
	match := map[byte]byte{')': '(', ']': '[', '}': '{'}
	stack := make([]byte, 0)
	for i := range s {
		c := s[i]
		if expected, isClosing := match[c]; isClosing {
			if len(stack) == 0 || stack[len(stack)-1] != expected {
				return false
			}
			stack = stack[:len(stack)-1]
		} else {
			stack = append(stack, c)
		}
	}
	return len(stack) == 0
}
