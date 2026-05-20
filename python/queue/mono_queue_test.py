"""
Unit tests for mono_queue.py — longestSubarray (LeetCode 1438)
"""
import sys
import os

import pytest

sys.path.insert(0, os.path.dirname(__file__))
from queue.mono_queue import Solution  # noqa: E402


@pytest.mark.parametrize("nums, limit, want", [
    # LeetCode examples
    ([8, 2, 4, 7],              4,   2),
    ([10, 1, 2, 4, 7, 2],       5,   4),
    ([4, 2, 2, 2, 4, 4, 2, 2],  0,   3),
    # single element — always valid
    ([1],                        0,   1),
    # entire array fits within limit
    ([1, 2, 3, 4, 5],            4,   5),
    # no two adjacent elements fit (all diffs > limit)
    ([1, 5, 1, 5, 1],            3,   1),
    # identical elements, limit=0 → full array
    ([3, 3, 3, 3],               0,   4),
    # two elements within limit
    ([1, 100],                  99,   2),
    # two elements exceeding limit → 1
    ([1, 100],                  98,   1),
    # large window shrinks only at the boundary
    ([1, 3, 6, 10, 15],         5,    3),
])
def test_longest_subarray(nums, limit, want):
    assert Solution().longestSubarray(nums, limit) == want
