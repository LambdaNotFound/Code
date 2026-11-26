"""
Unit tests for binary_search.py — LC 33: Search in Rotated Sorted Array,
LC 81: Search in Rotated Sorted Array II, LC 278: First Bad Version
"""
import sys
import os

import pytest

sys.path.insert(0, os.path.dirname(__file__))
import binary_search as bs
from binary_search import Solution  # noqa: E402


@pytest.mark.parametrize("n, first_bad, want", [
    # LeetCode example
    (5, 4, 4),

    # First version is bad
    (1, 1, 1),

    # Last version is bad
    (10, 10, 10),

    # Second version is bad
    (5, 2, 2),

    # Large n
    (2**31 - 1, 2**31 - 1, 2**31 - 1),

    # Two versions, first is bad
    (2, 1, 1),

    # Two versions, second is bad
    (2, 2, 2),
])
def test_firstBadVersion(monkeypatch, n, first_bad, want):
    monkeypatch.setattr(bs, "isBadVersion", lambda v: v >= first_bad)
    assert Solution().firstBadVersion(n) == want


# ── 33. Search in Rotated Sorted Array ───────────────────────────────────────

@pytest.mark.parametrize("nums, target, want", [
    # LC example: rotated, target on the left segment
    ([4, 5, 6, 7, 0, 1, 2],    0,      4),

    # target not present
    ([4, 5, 6, 7, 0, 1, 2],    3,      -1),

    # single-element match
    ([1],                       1,      0),

    # single element, no match
    ([1],                       2,      -1),

    # no rotation (sorted): target at the beginning
    ([1, 2, 3, 4, 5],          1,      0),

    # no rotation: target at the end
    ([1, 2, 3, 4, 5],          5,      4),

    # rotation point at index 1
    ([2, 0, 1],                 0,      1),

    # target is the rotation pivot
    ([6, 7, 0, 1, 2, 3, 4, 5], 6,      0),

    # empty array
    ([],                        3,      -1),
])
def test_search(nums, target, want):
    assert Solution().search(nums, target) == want


# ── 81. Search in Rotated Sorted Array II ────────────────────────────────────

@pytest.mark.parametrize("nums, target, want", [
    # LC example: contains duplicates, target present
    ([2, 5, 6, 0, 0, 1, 2],    0,      True),

    # LC example: target absent
    ([2, 5, 6, 0, 0, 1, 2],    3,      False),

    # all duplicates, target matches
    ([1, 1, 1, 1, 1],           1,      True),

    # all duplicates, target absent
    ([1, 1, 1, 1, 1],           2,      False),

    # duplicates at both ends obscure rotation point
    ([3, 1, 2, 3, 3, 3, 3],    2,      True),

    # single element match
    ([1],                       1,      True),

    # single element no match
    ([1],                       0,      False),

    # no rotation, with duplicates
    ([1, 1, 2, 3, 4, 4],        4,      True),
    ([1, 1, 2, 3, 4, 4],        5,      False),

    # left half sorted (nums[mid] > nums[right]); target inside left half → right = mid-1
    ([3, 4, 5, 1, 1, 2],        4,      True),

    # left half sorted; target outside left half → left = mid+1
    ([3, 4, 5, 1, 1, 2],        1,      True),
])
def test_searchWithDuplicates(nums, target, want):
    assert Solution().searchWithDuplicates(nums, target) == want
