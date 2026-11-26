"""
Unit tests for parentheses.py — LeetCode 32: Longest Valid Parentheses
"""
import sys
import os

import pytest

sys.path.insert(0, os.path.dirname(__file__))
from parentheses import Solution  # noqa: E402


@pytest.mark.parametrize("s, want", [
    # LeetCode examples
    ("(()",    2),
    (")()())", 4),
    ("",       0),

    # All valid
    ("()()",   4),
    ("(())",   4),

    # No valid pairs
    ("(((",    0),
    (")))",    0),

    # Valid suffix
    ("()(())", 6),

    # Disconnected valid segments — only longest counts
    ("()(())", 6),
    ("()(()",  2),

    # Single chars
    ("(",      0),
    (")",      0),
])
def test_longestValidParentheses(s, want):
    assert Solution().longestValidParentheses(s) == want
