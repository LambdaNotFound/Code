import os
import sys

import pytest

sys.path.insert(0, os.path.dirname(__file__))
from dynamic_programming.subsequence import Solution


@pytest.mark.parametrize(
    "nums, expected",
    [
        ([10, 9, 2, 5, 3, 7, 101, 18], 4),  # example from LC300
        ([0, 1, 0, 3, 2, 3], 4),  # interleaved values
        ([7, 7, 7, 7, 7], 1),  # all equal (strictly increasing)
        ([1], 1),  # single element
        ([1, 2, 3, 4, 5], 5),  # already sorted
        ([5, 4, 3, 2, 1], 1),  # strictly decreasing
    ],
)
def test_lengthOfLIS(nums, expected):
    assert Solution().lengthOfLIS(nums) == expected


@pytest.mark.parametrize(
    "nums, expected",
    [
        ([1, 3, 5, 4, 7], [[1, 3, 5, 7], [1, 3, 4, 7]]),
        ([7, 7, 7, 7, 7], [[7], [7], [7], [7], [7]]),
    ],
)
def test_computeLIS(nums, expected):
    assert Solution().computeLIS(nums) == expected
