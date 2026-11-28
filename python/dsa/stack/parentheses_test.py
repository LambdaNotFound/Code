"""
Unit tests for parentheses.py — LeetCode 32: Longest Valid Parentheses,
LeetCode 20: Valid Parentheses
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


@pytest.mark.parametrize("s, want", [
    # same cases as Test_isValid in go/dsa/stack/stack_test.go
    ("",       True),
    ("()",     True),
    ("()[]{}", True),
    ("(]",     False),
    ("([{}])", True),
    ("(((",    False),
    (")))",    False),
    ("([)]",   False),

    # closer before any opener, extra closer, extra opener
    ("){",     False),
    ("())",    False),
    ("(()",    False),
])
def test_isValid(s, want):
    assert Solution().isValid(s) == want
