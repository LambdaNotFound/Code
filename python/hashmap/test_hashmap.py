"""
Unit tests for hashmap.py

The source file defines two classes both named `Solution` — a LeetCode convention.
The second definition shadows the first at module level, so we exec each class in
its own namespace to test both independently without touching the source.
"""

import os
from typing import List, Optional
from collections import Counter

import pytest

# ── Load both Solution classes from the source file ───────────────────────────

_src_path = os.path.join(os.path.dirname(__file__), "hashmap.py")
with open(_src_path) as _f:
    _lines = _f.read().splitlines()

_class_starts = [i for i, line in enumerate(_lines) if line.strip().startswith("class Solution")]

# First Solution — twoSum; needs List from typing
_ns1: dict = {"List": List}
exec("\n".join(_lines[_class_starts[0] : _class_starts[1]]), _ns1)
TwoSum = _ns1["Solution"]

# Second Solution — canConstruct / canConstruct2; needs Counter
_ns2: dict = {"Counter": Counter}
exec("\n".join(_lines[_class_starts[1] :]), _ns2)
RansomNote = _ns2["Solution"]


# ── 1. Two Sum ────────────────────────────────────────────────────────────────

@pytest.mark.parametrize("nums, target, want", [
    # name                          nums                    target  want
    # answer is the first two elements
    ([2, 7, 11, 15],                9,                      [0, 1]),
    # answer sits in the middle of the array
    ([3, 2, 4],                     6,                      [1, 2]),
    # same value used twice — indices must differ
    ([3, 3],                        6,                      [0, 1]),
    # negative numbers
    ([-1, -2, -3, -4, -5],          -8,                     [2, 4]),
    # complement appears late; 7 + 2 = 9
    ([1, 5, 3, 7, 2],               9,                      [3, 4]),
    # minimal two-element array
    ([0, 4],                        4,                      [0, 1]),
    # zero target; -3 + 3 = 0
    ([-3, 0, 3],                    0,                      [0, 2]),
    # large values; 1_000_000 + 1 = 1_000_001
    ([1_000_000, 999_999, 1],       1_000_001,              [0, 2]),
    # no valid pair → None
    ([1, 2, 3],                     10,                     None),
])
def test_two_sum(nums, target, want):
    assert TwoSum().twoSum(nums, target) == want


# ── 383. Ransom Note ──────────────────────────────────────────────────────────

# Both canConstruct (Counter-based) and canConstruct2 (manual dict) must behave
# identically. We parametrize over the method name so every case runs twice.

@pytest.mark.parametrize("method", ["canConstruct", "canConstruct2"])
@pytest.mark.parametrize("ransom_note, magazine, want", [
    # name                          ransom_note     magazine        want
    # letter missing from magazine entirely
    ("a",                           "b",            False),
    # magazine has exactly the letters needed
    ("abc",                         "abc",          True),
    # magazine has extra letters — still constructible
    ("ab",                          "abcd",         True),
    # only one 'a' in magazine, need two
    ("aa",                          "ab",           False),
    # two 'a's present — sufficient
    ("aa",                          "aab",          True),
    # empty ransom note is always constructible
    ("",                            "abc",          True),
    # both empty — trivially constructible
    ("",                            "",             True),
    # non-empty ransom note, empty magazine
    ("a",                           "",             False),
    # all same characters, exact count
    ("aaa",                         "aaa",          True),
    # letter not present in magazine at all
    ("z",                           "abc",          False),
    # ransom note longer than magazine
    ("abcde",                       "abc",          False),
    # single-character match inside a longer magazine
    ("x",                           "xyz",          True),
    # case-sensitive: 'A' ≠ 'a'
    ("A",                           "a",            False),
    # each magazine letter usable at most once
    ("aabb",                        "aabb",         True),
    ("aabb",                        "ab",           False),
])
def test_can_construct(method, ransom_note, magazine, want):
    fn = getattr(RansomNote(), method)
    assert fn(ransom_note, magazine) == want
