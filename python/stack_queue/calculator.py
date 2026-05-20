from collections import deque

"""
224. Basic Calculator

Input: s = "(1+(4+5+2)-3)+(6+8)"
Output: 23

"""
class Solution:
    def calculate(self, s: str) -> int:
        stack: deque[int] = deque()
