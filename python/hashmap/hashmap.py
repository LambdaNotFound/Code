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
from collections import Counter
class Solution:
    def canConstruct(self, ransomNote: str, magazine: str) -> bool:
        counts = Counter(magazine)
        for ch in ransomNote:
            if counts[ch] > 0:
                counts[ch] -= 1
            else:
                return False
        return True

    def canConstruct2(self, ransomNote: str, magazine: str) -> bool:
        counts: dict[str, int] = {}
        for ch in magazine:
            counts[ch] = counts.get(ch, 0) + 1

        for ch in ransomNote:
            if counts.get(ch, 0) > 0:
                counts[ch] -= 1
            else:
                return False
        return True