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
            else:  # overlapping
                intersection = [max(first[0], second[0]), min(first[1], second[1])]
                result.append(intersection)

                if first[1] < second[1]:
                    i += 1
                else:
                    j += 1
        return result
