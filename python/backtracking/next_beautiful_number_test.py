"""
Unit tests for next_beautiful_number.py — Next Greater Numerically Balanced Number (LeetCode 2048)
"""
import sys
import os

import pytest

sys.path.insert(0, os.path.dirname(__file__))
from next_beautiful_number import Solution

METHODS = ["nextBeautifulNumber", "next_beautiful_number"]


@pytest.mark.parametrize("method", METHODS)
@pytest.mark.parametrize("n, want", [
    # LeetCode examples
    (1,    22),    # next balanced > 1 is 22 (two 2s); 1 itself is balanced but not strictly greater
    (1000, 1333),

    # Just below a balanced number
    (1332, 1333),

    # Exact balanced number — must return strictly next
    (1333, 3133),  # 2222 is invalid (digit 2 appears 4 times, not 2); next is 3133

    # Boundary at the start
    (0,    1),

    # n=3000: 3111 is invalid (digit 1 appears 3 times, not 1); next is 3133
    (3000, 3133),

    # Last 4-digit balanced number; 4444 (digit 4 appears 4 times) is next
    (3333, 4444),

    # 13333 is invalid (digit 3 appears 4 times, not 3); next is 14444
    (13000, 14444),
])
def test_next_beautiful_number(method, n, want):
    assert getattr(Solution(), method)(n) == want
