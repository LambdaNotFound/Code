from collections import deque


class Solution:
    """
    32. Longest Valid Parentheses
     1. stack
     2. DP
     3. two pointers
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
    Case 1: s[i-1] == '('   → direct match "()"
            dp[i] = dp[i-2] + 2

    Case 2: s[i-1] == ')'   → previous char closes its own group "))"
        need to check if char BEFORE that group matches
        j = i - dp[i-1] - 1
        if s[j] == '(' :
            dp[i] = dp[i-1] + 2 + dp[j-1]

    translated from go/dsa/stack/parentheses.go
    """
    def longestValidParenthesesDP(self, s: str) -> int:
        n = len(s)
        dp = [0] * n
        res = 0

        for i in range(1, n):
            if s[i] == ')':
                if s[i - 1] == '(':
                    # case 1: direct "()" match
                    dp[i] = dp[i - 2] + 2 if i >= 2 else 2
                elif s[i - 1] == ')':
                    # case 2: ends with "))" — look before inner group
                    j = i - dp[i - 1] - 1
                    if j >= 0 and s[j] == '(':  # j >= 0: s[-1] is the last char
                        dp[i] = dp[i - 1] + 2
                        if j > 0:
                            dp[i] += dp[j - 1]
                res = max(res, dp[i])
        return res

    """
    Count '(' and ')'; a run is valid when the counts are equal. Left to right
    resets on an unmatched ')', but never sees the counts equal inside "(()",
    so a second pass right to left, resetting on an unmatched '(', finds it.
    """
    def longestValidParenthesesTwoPointers(self, s: str) -> int:
        res = 0

        # pass 1: left → right
        opens, closes = 0, 0
        for ch in s:
            if ch == '(':
                opens += 1
            else:
                closes += 1
            if opens == closes:
                res = max(res, 2 * closes)
            elif closes > opens:
                opens, closes = 0, 0  # unmatched ')' — reset

        # pass 2: right → left
        opens, closes = 0, 0
        for ch in reversed(s):
            if ch == '(':
                opens += 1
            else:
                closes += 1
            if opens == closes:
                res = max(res, 2 * opens)
            elif opens > closes:
                opens, closes = 0, 0  # unmatched '(' — reset

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
