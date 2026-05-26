from collections import deque

class Solution:
    """
    224. Basic Calculator

    Input: s = "(1+(4+5+2)-3)+(6+8)"
    Output: 23
    """
    def calculate(self, s: str) -> int:
        s = s.replace(" ", "")
        pos = [0]

        def parse() -> int:
            stack: deque[int] = deque() # stores sign * number for current call stack
            current_num, pending_sign = 0, 1

            while pos[0] < len(s):
                char = s[pos[0]]
                pos[0] += 1

                if '0' <= char <= '9':
                    current_num = current_num * 10 + int(char)
                elif char == '+' or char == '-':
                    stack.append(pending_sign * current_num)
                    current_num = 0
                    pending_sign = 1 if char == '+' else -1
                elif char == '(':
                    sub_result = parse()
                    stack.append(pending_sign * sub_result)
                elif char == ')':
                    break

            stack.append(pending_sign * current_num)
            return sum(stack)

        return parse()
            
    def calculate_iterative(self, s: str) -> int:
        s = s.replace(" ", "")
        stack: deque[int] = deque() # stores result + sign before each '('
        result, current_num, pending_sign = 0, 0, 1 # result -> total accumulator

        for char in s:
            if '0' <= char <= '9':
                current_num = current_num * 10 + int(char)

            if char == '+' or char == '-':
                result += pending_sign * current_num
                current_num = 0
                pending_sign = 1 if char == '+' else -1

            if char == '(':
                stack.append(result)
                stack.append(pending_sign)
                result, pending_sign = 0, 1

            if char == ')':
                result += pending_sign * current_num
                current_num = 0

                outer_sign = stack.pop()
                outer_result = stack.pop()
                result = outer_result + outer_sign * result
        result += pending_sign * current_num
        
        return result
    """
    227. Basic Calculator II

    Input: s = "3+2*2"
    Output: 7
    """
    def calculate2(self, s: str) -> int:
        pass

    """
    772. Basic Calculator III
    Input: s = "(2+6*3+5-(3*14/7+2)*5)+3"
    Output: -12
    """
    def calculate3(self, s: str) -> int:
        pass