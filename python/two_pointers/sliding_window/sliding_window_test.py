"""
Unit tests for sliding_window.py — minWindow (LeetCode 76)
"""
import sys
import os

import pytest

sys.path.insert(0, os.path.dirname(__file__))
from sliding_window import Solution


@pytest.mark.parametrize("s, t, want", [
    # canonical example from the problem
    ("ADOBECODEBANC",  "ABC",  "BANC"),
    # t is a single char present in s
    ("a",              "a",    "a"),
    # t equals s exactly
    ("abc",            "abc",  "abc"),
    # t not present in s at all
    ("aa",             "b",    ""),
    # duplicate chars in t — window must contain enough copies
    ("AABC",           "AA",   "AA"),          # two A's at the start — exact match
    ("ABAC",           "AA",   "ABA"),         # shortest window with two A's: positions 0-2
    ("ADOBECODEBANC",  "AA",   "ADOBECODEBA"), # A's at 0 and 10 — spans whole prefix
    ("xyz",            "aa",   ""),            # char present but not enough copies
    # surplus chars in s — verify surplus is skipped while shrinking
    ("BBBBBBAABC",     "ABC",  "ABC"),   # tightest window is the rightmost ABC
    # t longer than s
    ("ab",             "abc",  ""),
    # empty s
    ("",               "a",    ""),
    # empty t
    ("abc",            "",     ""),
    # t has chars not in s
    ("xyz",            "a",    ""),
    # multiple valid windows of same length — leftmost is returned
    ("BCABC",          "ABC",  "BCA"),   # BCA (0-2) and CAB/ABC (2-4) tie at 3; BCA found first
    # answer is at the start of s
    ("ABCDE",          "AB",   "AB"),
])
def test_min_window(s, t, want):
    assert Solution().minWindow(s, t) == want
