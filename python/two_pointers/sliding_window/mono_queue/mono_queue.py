"""
1438. Longest Continuous Subarray With Absolute Diff Less Than or Equal to Limit
"""
from collections import deque
class Solution:
    def longestSubarray(self, nums: List[int], limit: int) -> int:
        max_deque: deque[int] = deque()  # monotonic decreasing — front is max
        min_deque: deque[int] = deque()  # monotonic increasing — front is min
        result = 0

        left = 0
        for right, num in enumerate(nums):
            # maintain monotonic decreasing max deque
            while max_deque and max_deque[-1] < num:
                max_deque.pop()
            max_deque.append(num)

            # maintain monotonic increasing min deque
            while min_deque and min_deque[-1] > num:
                min_deque.pop()
            min_deque.append(num)

            # shrink window until constraint is satisfied
            while max_deque[0] - min_deque[0] > limit:
                evicted = nums[left]
                left += 1
                if evicted == max_deque[0]:
                    max_deque.popleft()
                if evicted == min_deque[0]:
                    min_deque.popleft()

            result = max(result, right - left + 1)

        return result