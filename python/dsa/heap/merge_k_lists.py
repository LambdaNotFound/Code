import sys, os
sys.path.insert(0, os.path.dirname(os.path.dirname(__file__)))
from common import ListNode

from typing import List, Optional
import heapq

"""
23. Merge k Sorted Lists
"""
class Solution:
    def mergeKLists(self, lists: List[Optional[ListNode]]) -> Optional[ListNode]:
        minHeap: List[tuple[int, int, ListNode]] = []
        for i, node in enumerate(lists):
            if node:
                heapq.heappush(minHeap, (node.val, i, node))

        dummy: ListNode = ListNode()
        cur = dummy
        while len(minHeap) > 0:
            _, i, top = heapq.heappop(minHeap)
            cur.next = top
            cur = cur.next
            if top.next:
                heapq.heappush(minHeap, (top.next.val, i, top.next))

        return dummy.next
