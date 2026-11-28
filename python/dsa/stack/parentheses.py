from collections import deque


class Solution:
    """
    32. Longest Valid Parentheses
    """
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

    """
    20. Valid Parentheses

    push the closer each opener expects; a closer must match the top
    translated from go/dsa/stack/parentheses.go
    """
    def isValid(self, s: str) -> bool:
        closer = {'(': ')', '[': ']', '{': '}'}
        stack: list[str] = []
        for ch in s:
            if ch in closer:
                stack.append(closer[ch])
            elif not stack or stack.pop() != ch:  # closing bracket
                return False
        return not stack
