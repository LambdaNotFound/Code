from typing import List

"""
Sort by a computed value:

data.sort(key=lambda a: (a[0], -a[1], a[0] + a[1] + a[2]))
                           │      │         └─ tie on a[1] → sum all three, asc
                           │      └─ tie on a[0] → a[1] descending
                           └─ primary: a[0] ascending
"""


class Solution:
    """
    56. Merge Intervals
    """

    def merge(self, intervals: List[List[int]]) -> List[List[int]]:
        intervals.sort(key=lambda i: i[0])

        merged = [intervals[0]]
        for curr in range(1, len(intervals)):
            if merged[-1][1] >= intervals[curr][0]:  # overlapping intervals
                merged[-1] = [
                    min(merged[-1][0], intervals[curr][0]),
                    max(merged[-1][1], intervals[curr][1]),
                ]
            else:
                merged.append(intervals[curr])
        return merged

    """
    57. Insert Interval
    """

    def insert(
        self, intervals: List[List[int]], newInterval: List[int]
    ) -> List[List[int]]:
        before, after = [], []
        for curr in range(len(intervals)):
            if intervals[curr][1] < newInterval[0]:
                before.append(intervals[curr])
            elif newInterval[1] < intervals[curr][0]:
                after.append(intervals[curr])
            else:  # overlapping intervals
                newInterval = [
                    min(newInterval[0], intervals[curr][0]),
                    max(newInterval[1], intervals[curr][1]),
                ]

        return before + [newInterval] + after

    """
    435. Non-overlapping Intervals
    """

    def eraseOverlapIntervals(self, intervals: List[List[int]]) -> int:
        intervals.sort(key=lambda i: i[0])
        removed, prev = 0, 0
        for curr in range(1, len(intervals)):
            if intervals[prev][1] <= intervals[curr][0]:
                prev = curr
            else:  # overlapping intervals
                removed += 1
                if intervals[curr][1] < intervals[prev][1]:
                    prev = curr

        return removed

    """
    986. Interval List Intersections
    """

    def intervalIntersection(
        self, firstList: List[List[int]], secondList: List[List[int]]
    ) -> List[List[int]]:
        result: List[List[int]] = []
        i, j = 0, 0
        while i < len(firstList) and j < len(secondList):
            first, second = firstList[i], secondList[j]
            if first[1] < second[0]:
                i += 1
            elif first[0] > second[1]:
                j += 1
            else:  # overlapping intervals
                intersection = [max(first[0], second[0]), min(first[1], second[1])]
                result.append(intersection)

                if first[1] < second[1]:
                    i += 1
                else:
                    j += 1
        return result
