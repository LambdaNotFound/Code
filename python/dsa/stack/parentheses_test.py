"""
Unit tests for parentheses.py — LeetCode 32: Longest Valid Parentheses
"""

import pytest

from dsa.stack.parentheses import Solution


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
