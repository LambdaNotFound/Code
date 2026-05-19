from typing import List

"""
41. First Missing Positive

goal: [1, 2, 3, 4, 5]
       0  1  2  3  4
"""
class Solution:
    def firstMissingPositive(self, nums: List[int]) -> int:
        n = len(nums)
        
        # Phase 1: cyclic sort — place value v at index v-1 for v in [1, n]
        for i, _ in enumerate(nums):
            while 0 <= nums[i]-1 <= n-1 and nums[i] != nums[nums[i] - 1]:
                target = nums[i] - 1  # cache before swap — see notes below
                nums[i], nums[target] = nums[target], nums[i]
        
        # Phase 2: find first mismatch
        for i, _ in enumerate(nums):
            if nums[i] != i + 1:
                return i + 1
        
        return n + 1  # all [1..n] present → answer is n+1