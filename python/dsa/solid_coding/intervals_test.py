import os
import sys

sys.path.insert(0, os.path.dirname(__file__))

import pytest
from intervals import Solution


@pytest.mark.parametrize(
    "intervals,expected",
    [
        ([[1, 3], [2, 6], [8, 10], [15, 18]], [[1, 6], [8, 10], [15, 18]]),
        ([[1, 4], [4, 5]], [[1, 5]]),
        ([[1, 2], [3, 4], [5, 6]], [[1, 2], [3, 4], [5, 6]]),
        ([[1, 5]], [[1, 5]]),
        ([[2, 6], [1, 3]], [[1, 6]]),
        ([[1, 10], [2, 5], [3, 7]], [[1, 10]]),
    ],
)
def test_merge(intervals, expected):
    assert Solution().merge(intervals) == expected


@pytest.mark.parametrize(
    "intervals,newInterval,expected",
    [
        ([[1, 3], [6, 9]], [2, 5], [[1, 5], [6, 9]]),
        (
            [[1, 2], [3, 5], [6, 7], [8, 10], [12, 16]],
            [4, 8],
            [[1, 2], [3, 10], [12, 16]],
        ),
        ([], [5, 7], [[5, 7]]),
        ([[1, 5]], [6, 8], [[1, 5], [6, 8]]),
        ([[3, 5]], [1, 2], [[1, 2], [3, 5]]),
        ([[1, 5]], [2, 3], [[1, 5]]),
        ([[1, 2], [3, 4]], [0, 10], [[0, 10]]),
    ],
)
def test_insert(intervals, newInterval, expected):
    assert Solution().insert(intervals, newInterval) == expected


@pytest.mark.parametrize(
    "intervals,expected",
    [
        ([[1, 2], [2, 3], [3, 4], [1, 3]], 1),
        ([[1, 2], [1, 2], [1, 2]], 2),
        ([[1, 2], [2, 3]], 0),
        ([[1, 100], [11, 22], [1, 11], [2, 12]], 2),
        ([[0, 1]], 0),
    ],
)
def test_eraseOverlapIntervals(intervals, expected):
    assert Solution().eraseOverlapIntervals(intervals) == expected


@pytest.mark.parametrize(
    "firstList,secondList,expected",
    [
        (
            [[0, 2], [5, 10], [13, 23], [24, 25]],
            [[1, 5], [8, 12], [15, 24], [25, 26]],
            [[1, 2], [5, 5], [8, 10], [15, 23], [24, 24], [25, 25]],
        ),
        ([[1, 3], [5, 9]], [], []),
        ([], [[4, 8], [10, 12]], []),
        ([[1, 7]], [[3, 10]], [[3, 7]]),
        ([[1, 4], [6, 8]], [[2, 3], [5, 7]], [[2, 3], [6, 7]]),
    ],
)
def test_intervalIntersection(firstList, secondList, expected):
    assert Solution().intervalIntersection(firstList, secondList) == expected
