from collections import defaultdict
from typing import List


class Solution:
    """
    1. Two Sum
    """

    def twoSum(self, nums: List[int], target: int) -> List[int]:
        complements: dict[int] = {}
        for i, num in enumerate(nums):
            complement = target - num
            if complement in complements:
                return [complements[complement], i]
            complements[num] = i
        return None

    """
    383. Ransom Note

    Time: O(m + n)
    Space: O(k), effectively O(1) if ASCII

    follow up: If magazine is massive but ransomNote is short?
    
    Optimization: scan ransomNote first, then stream magazine with early exit
    Build a "needs" map from the small ransomNote, then stop scanning magazine the moment all needs are satisfied:
    """

    def canConstruct(self, ransomNote: str, magazine: str) -> bool:
        needs: defaultdict = defaultdict(int)
        for ch in magazine:
            needs[ch] += 1

        for ch in ransomNote:
            if needs[ch] > 0:
                needs[ch] -= 1
            else:
                return False
        return True
