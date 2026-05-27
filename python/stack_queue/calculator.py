from collections import deque

class Solution:
    """
    Basic Calculator Template
    """
    def calculate_template(self, s: str) -> int:
        s = s.replace(" ", "")
        pos = [0]

        def parse() -> int:
            stack = []
            current_number, pending_sign = 0, '+'

            while pos[0] < len(s):
                ch = s[pos[0]]
                pos[0] += 1

                if '0' <= ch <= '9':
                    current_number = current_number * 10 + int(ch)

                if ch == '(':
                    current_number = parse()

                if pos[0] == len(s) or ch in '+-*/)': # commit on operator, end, or closing paren
                    if pending_sign == '+':
                        stack.append(current_number)
                    elif pending_sign == '-':
                        stack.append(-current_number)
                    elif pending_sign == '*':
                        stack[-1] *= current_number
                    elif pending_sign == '/':
                        stack[-1] = int(stack[-1] / current_number)

                    if ch == ')':
                        break

                    current_number = 0
                    pending_sign = ch

            return sum(stack)

        return parse()

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
                if char == '(':
                    current_num = parse()
                if pos[0] == len(s) or char == '+' or char == '-' or char == ')':
                    stack.append(pending_sign * current_num)
                    if char == ')':
                        break
                    current_num = 0
                    pending_sign = 1 if char == '+' else -1
            return sum(stack)

        return parse()
            
    def calculate_iterative(self, s: str) -> int:
        s = s.replace(" ", "")
        stack: deque[int] = deque() # stores result + sign before each '('
        result, current_num, pending_sign = 0, 0, 1 # result -> total accumulator
        for char in s:
            if '0' <= char <= '9':
                current_num = current_num * 10 + int(char)
            elif char == '+' or char == '-':
                result += pending_sign * current_num
                current_num = 0
                pending_sign = 1 if char == '+' else -1
            elif char == '(':
                stack.append(result)
                stack.append(pending_sign)
                result, pending_sign = 0, 1
            elif char == ')':
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
        s = s.replace(' ', '')
        stack: deque[int] = deque()

        current_number, pending_sign = 0, '+'
        for i, char in enumerate(s):
            if '0' <= char <= '9':
                current_number = current_number * 10 + int(char)
            if i == len(s)-1 or char in '+-*/':
                if pending_sign == '+':
                    stack.append(current_number)
                elif pending_sign == '-':
                    stack.append(-current_number)
                elif pending_sign == '*':
                    stack.append(stack.pop() * current_number)
                elif pending_sign == '/':
                    stack.append(int(stack.pop() / current_number))
                current_number, pending_sign = 0, char

        return sum(stack)

    """
    772. Basic Calculator III
    Input: s = "(2+6*3+5-(3*14/7+2)*5)+3"
    Output: -12
    """
    def calculate3(self, s: str) -> int:
        s = s.replace(" ", "")
        pos = [0]

        def parse() -> int:
            stack: deque[int] = deque()
            current_number, pending_sign = 0, '+'

            while pos[0] < len(s):
                char = s[pos[0]]
                pos[0] += 1

                if '0' <= char <= '9':
                    current_number = current_number * 10 + int(char)
                if char == '(':
                    current_number = parse()
                if pos[0] == len(s) or char in '+-*/)':
                    if pending_sign == '+':
                        stack.append(current_number)
                    elif pending_sign == '-':
                        stack.append(-current_number)
                    elif pending_sign == '*':
                        stack.append(stack.pop() * current_number)
                    elif pending_sign == '/':
                        stack.append(int(stack.pop() / current_number))

                    if char == ')':
                        break
                    current_number, pending_sign = 0, char

            return sum(stack)

        return parse()