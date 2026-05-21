from collections import deque

"""
224. Basic Calculator

Input: s = "(1+(4+5+2)-3)+(6+8)"
Output: 23
"""
class Solution:
    def calculate(self, s: str) -> int:
        s = s.replace(" ", "")
        self.pos = 0

        def parse() -> int:
            stack: deque[int] = deque()
            num, sign = 0, 1

            while self.pos < len(s):
                ch = s[self.pos]
                self.pos += 1

                if '0' <= ch <= '9':
                    num = num * 10 + int(ch)
                elif ch == '+' or ch == '-':
                    stack.append(sign * num)
                    num = 0
                    sign = 1 if ch == '+' else -1
                elif ch == '(':
                    sub = parse()
                    stack.append(sign * sub)
                elif ch == ')':
                    break

            stack.append(sign * num)
            return sum(stack)

        return parse()
            
