from collections import deque

""""
1209. Remove All Adjacent Duplicates in String II
"""
class Solution:
    def removeDuplicates(self, s: str, k: int) -> str:
        stack: deque = deque()
        for _, ch in enumerate(s):
            stack.append(ch)
            if len(stack) >= k:
                i = len(stack) - k
                while i < len(stack) and stack[i] == stack[-1]: # Time: O(n·k) overall
                    i += 1
                if i == len(stack):
                    for _ in range(k):
                       stack.pop()
        
        return ''.join(stack)
        
    def removeDuplicatesClaude(self, s: str, k: int) -> str:
        stack: deque[tuple[str, int]] = deque()
        for ch in s:
            if stack and stack[-1][0] == ch: # checks if stack is empty
                char, count = stack.pop()
                stack.append((char, count + 1))
            else:
                stack.append((ch, 1))

            if stack[-1][1] == k: # Time: O(n)
                stack.pop()
        return ''.join(ch * count for ch, count in stack)