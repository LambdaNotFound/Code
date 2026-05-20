from typing import List
from collections import defaultdict

"""
316. Remove Duplicate Letters
"""
class Solution:
    def removeDuplicateLetters(self, s: str) -> str:
        lastIndex: dict = {}
        for i, ch in enumerate(s):
            lastIndex[ch] = i

        stack: List = []
        inStack: defaultdict = defaultdict(bool)
        for i, ch in enumerate(s):
            if inStack[ch]:
                continue
            while stack and stack[-1] > ch and lastIndex[stack[-1]] > i:
                inStack[stack[-1]] = False
                stack.pop()
            inStack[ch] = True
            stack.append(ch)

        return ''.join(stack)