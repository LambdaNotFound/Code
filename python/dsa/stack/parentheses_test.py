"""
Unit tests for parentheses.py — LeetCode 32: Longest Valid Parentheses,
LeetCode 20: Valid Parentheses
"""

import pytest

from dsa.stack.parentheses import Solution


METHODS = [
    "longestValidParentheses",
    "longestValidParenthesesDP",
    "longestValidParenthesesTwoPointers",
]


@pytest.mark.parametrize("method", METHODS)
@pytest.mark.parametrize("s, want", [
    # same cases as Test_longestValidParentheses in go/dsa/stack/stack_extra_test.go
    ("(()",    2),
    (")()())", 4),
    ("",       0),
    ("(((",    0),
    (")))",    0),
    ("()()",   4),
    ("(())",   4),
    ("()(())", 6),
    ("()((",   2),
    ("()",     2),
    ("())()(", 2),

    # longest run in the middle, between resets
    ("(()))())(", 4),
    (")()())()()(", 4),
    ("()(()",  2),
    # unmatched '(' on the left: only the right-to-left pass finds it
    ("(()(((()", 2),
    # nested groups joined to earlier ones (DP case 2 with dp[j-1])
    (")(((((()())()()))()(()))(", 22),

    # single chars
    ("(",      0),
    (")",      0),
])
def test_longestValidParentheses(s, want, method):
    assert getattr(Solution(), method)(s) == want


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
