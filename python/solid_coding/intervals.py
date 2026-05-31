from typing import List


class Solution:
    """
    435. Non-overlapping Intervals
    """

    def eraseOverlapIntervals(self, intervals: List[List[int]]) -> int:
        intervals.sort(key=lambda i: i[0])
        erase, prev = 0, 0
        for curr in range(1, len(intervals)):
            if intervals[prev][1] <= intervals[curr][0]:
                prev = curr
            else:  # overlapping intervals
                erase += 1
                if intervals[curr][1] < intervals[prev][1]:
                    prev = curr

        return erase
