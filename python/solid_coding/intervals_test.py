import sys, os
sys.path.insert(0, os.path.dirname(__file__))

import pytest
from intervals import Solution


@pytest.mark.parametrize("intervals,expected", [
    ([[1, 2], [2, 3], [3, 4], [1, 3]], 1),
    ([[1, 2], [1, 2], [1, 2]], 2),
    ([[1, 2], [2, 3]], 0),
    ([[1, 100], [11, 22], [1, 11], [2, 12]], 2),
    ([[0, 1]], 0),
])
def test_eraseOverlapIntervals(intervals, expected):
    assert Solution().eraseOverlapIntervals(intervals) == expected
