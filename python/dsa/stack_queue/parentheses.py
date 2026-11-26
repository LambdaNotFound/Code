from collections import deque

"""
32. Longest Valid Parentheses
"""
class Solution:
    def longestValidParentheses(self, s: str) -> int:
        stack: deque[int] = deque()
        stack.append(-1)

        res = 0
        for i, ch in enumerate(s):
            if ch == '(':
                stack.append(i)
            else:
                stack.pop()

                if stack:
                    res = max(res, i - stack[-1])
                else:
                    stack.append(i)

        return res