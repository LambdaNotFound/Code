"""
Unit tests for binary_search.py — LeetCode 278: First Bad Version
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
