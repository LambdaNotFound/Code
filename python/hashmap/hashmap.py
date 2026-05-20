from typing import List

"""
1. Two Sum
"""
class Solution:
    def twoSum(self, nums: List[int], target: int) -> List[int]:
        complements: dict[int, int] = {}
        for i, num in enumerate(nums):
            complement = target - num
            if complement in complements:
                return [complements[complement], i]
            complements[num] = i
        return None

"""
383. Ransom Note
"""
from collections import defaultdict
class Solution:
    def canConstruct(self, ransomNote: str, magazine: str) -> bool:
        counts: defaultdict[int] = defaultdict(int)
        for ch in magazine:
            counts[ch] += 1

        for ch in ransomNote:
            if counts[ch] > 0:
                counts[ch] -= 1
            else:
                return False
        return True
