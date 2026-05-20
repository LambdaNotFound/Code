"""
Unit tests for array.py — firstMissingPositive (LeetCode 41)
"""
import importlib.util
import os

import pytest

# array.py conflicts with stdlib; load via spec
_spec = importlib.util.spec_from_file_location(
    "array_sol", os.path.join(os.path.dirname(__file__), "array.py")
)
_mod = importlib.util.module_from_spec(_spec)
_spec.loader.exec_module(_mod)
Solution = _mod.Solution


@pytest.mark.parametrize("nums, want", [
    # smallest positive entirely absent
    ([1, 2, 0],          3),
    # 1 is missing, negative present
    ([3, 4, -1, 1],      2),
    # all values exceed n → answer is 1
    ([7, 8, 9, 11, 12],  1),
    # single missing beyond the array
    ([1],                2),
    # complete [1..n] → n+1
    ([1, 2, 3],          4),
    # two elements in wrong order
    ([2, 1],             3),
    # duplicates, 2 missing
    ([1, 1],             2),
    # all negative → 1
    ([-1, -2],           1),
    # zero only → 1
    ([0],                1),
    # 1 absent, only 2 present
    ([2],                1),
    # longer shuffled range
    ([5, 3, 2, 1, 4],   6),
    # large gap after 1
    ([1, 100],           2),
])
def test_first_missing_positive(nums, want):
    assert Solution().firstMissingPositive(nums) == want
