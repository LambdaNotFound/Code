"""
Unit tests for hashmap.py — Two Sum (1) and Ransom Note (383)
"""
import sys
import os

import pytest

sys.path.insert(0, os.path.dirname(__file__))
from hashmap import Solution  # noqa: E402


# ── 1. Two Sum ────────────────────────────────────────────────────────────────

@pytest.mark.parametrize("nums, target, want", [
    # answer is the first two elements
    ([2, 7, 11, 15],                9,          [0, 1]),
    # answer sits in the middle of the array
    ([3, 2, 4],                     6,          [1, 2]),
    # same value used twice — indices must differ
    ([3, 3],                        6,          [0, 1]),
    # negative numbers
    ([-1, -2, -3, -4, -5],          -8,         [2, 4]),
    # complement appears late; 7 + 2 = 9
    ([1, 5, 3, 7, 2],               9,          [3, 4]),
    # minimal two-element array
    ([0, 4],                        4,          [0, 1]),
    # zero target; -3 + 3 = 0
    ([-3, 0, 3],                    0,          [0, 2]),
    # large values; 1_000_000 + 1 = 1_000_001
    ([1_000_000, 999_999, 1],       1_000_001,  [0, 2]),
    # no valid pair → None
    ([1, 2, 3],                     10,         None),
])
def test_two_sum(nums, target, want):
    assert Solution().twoSum(nums, target) == want


# ── 383. Ransom Note ──────────────────────────────────────────────────────────

@pytest.mark.parametrize("ransom_note, magazine, want", [
    # letter missing from magazine entirely
    ("a",       "b",    False),
    # magazine has exactly the letters needed
    ("abc",     "abc",  True),
    # magazine has extra letters — still constructible
    ("ab",      "abcd", True),
    # only one 'a' in magazine, need two
    ("aa",      "ab",   False),
    # two 'a's present — sufficient
    ("aa",      "aab",  True),
    # empty ransom note is always constructible
    ("",        "abc",  True),
    # both empty — trivially constructible
    ("",        "",     True),
    # non-empty ransom note, empty magazine
    ("a",       "",     False),
    # all same characters, exact count
    ("aaa",     "aaa",  True),
    # letter not present in magazine at all
    ("z",       "abc",  False),
    # ransom note longer than magazine
    ("abcde",   "abc",  False),
    # single-character match inside a longer magazine
    ("x",       "xyz",  True),
    # case-sensitive: 'A' ≠ 'a'
    ("A",       "a",    False),
    # each magazine letter usable at most once
    ("aabb",    "aabb", True),
    ("aabb",    "ab",   False),
])
def test_can_construct(ransom_note, magazine, want):
    assert Solution().canConstruct(ransom_note, magazine) == want
