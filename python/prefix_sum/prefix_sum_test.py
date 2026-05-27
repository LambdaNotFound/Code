"""
Unit tests for product_except_self.py — Product of Array Except Self (LeetCode 238)
"""
import sys
import os

import pytest

sys.path.insert(0, os.path.dirname(__file__))
from prefix_sum import Solution


@pytest.mark.parametrize("nums, want", [
    # LeetCode examples
    ([1, 2, 3, 4],      [24, 12, 8, 6]),
    ([-1, 1, 0, -3, 3], [0, 0, 9, 0, 0]),

    # Two elements
    ([2, 3],            [3, 2]),

    # Contains one zero — only the zero position gets a non-zero product
    ([1, 2, 0, 4],      [0, 0, 8, 0]),

    # Contains two zeros — all products are zero
    ([0, 0, 3],         [0, 0, 0]),

    # All ones
    ([1, 1, 1, 1],      [1, 1, 1, 1]),

    # Negative numbers
    ([-1, -2, -3, -4],  [-24, -12, -8, -6]),

    # Mixed sign
    ([1, -1, 1, -1],    [1, -1, 1, -1]),

    # Large values (no overflow concern in Python)
    ([100, 200, 300],   [60000, 30000, 20000]),
])
def test_product_except_self(nums, want):
    assert Solution().productExceptSelf(nums) == want
    assert Solution().product_except_self(nums) == want
