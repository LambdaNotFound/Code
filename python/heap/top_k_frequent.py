from collections import Counter
from typing import List
import heapq

"""
347. Top K Frequent Elements
"""
class Solution:
    def topKFrequent(self, nums: List[int], k: int) -> List[int]:
        count = Counter(nums)

        heap: List[tuple[int, int]] = []
        for num, freq in count.items():
            heapq.heappush(heap, (freq, num))
            if len(heap) > k:
                heapq.heappop(heap)   # evict least frequent

        return [num for _, num in heap]
    
    def topKFrequentClaude(self, nums: List[int], k: int) -> List[int]:
        count = Counter(nums)
        return heapq.nlargest(k, count.keys(), key=lambda x: count[x])