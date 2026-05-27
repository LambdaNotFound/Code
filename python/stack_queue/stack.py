from collections import deque

class Solution:
    """
    1209. Remove All Adjacent Duplicates in String II
    """
    def removeDuplicatesClaude(self, s: str, k: int) -> str:
        stack: deque[tuple[str, int]] = deque()
        for ch in s:
            if stack and stack[-1][0] == ch: # checks if stack is empty
                char, count = stack.pop()
                stack.append((char, count + 1))
            else:
                stack.append((ch, 1))

            if stack[-1][1] == k: # Time complexity: O(n)
                stack.pop()

        return ''.join(ch * count for ch, count in stack)

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
    
    """
    394. Decode String

    Input: s = "3[a]2[bc]"
    Output: "aaabcbc"

    Input: s = "3[a2[c]]"
    Output: "accaccacc"
    """
    def decodeString(self, s: str) -> str:
        stack: deque[str] = deque()

        for char in s:
            if char != ']':
                stack.append(char)
            else:
                letters, nums = '', ''
                while stack[-1] != '[':
                    letters = stack.pop() + letters
                stack.pop()
                while stack and '0' <= stack[-1] <= '9':
                    nums =  stack.pop() + nums
                stack.append(letters * int(nums))
        
        return ''.join(stack)