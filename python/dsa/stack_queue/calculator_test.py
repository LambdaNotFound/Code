"""
Unit tests for calculator.py — Basic Calculator I/II/III (LeetCode 224/227/772)
"""
import sys
import os

import pytest

sys.path.insert(0, os.path.dirname(__file__))
from calculator import Solution  # noqa: E402


@pytest.mark.parametrize("s, want", [
    # LeetCode examples
    ("1 + 1",               2),
    (" 2-1 + 2 ",           3),
    ("(1+(4+5+2)-3)+(6+8)", 23),

    # Single number
    ("0",                   0),
    ("42",                  42),

    # Unary-style negation via subtraction from zero
    ("0-5",                 -5),

    # Nested parentheses
    ("((3))",               3),
    ("(1+(2+(3+(4+5))))",   15),

    # Subtraction only
    ("10-3-2",              5),

    # Mixed add/subtract, no parens
    ("1+2+3+4+5",           15),
    ("10-1-2-3",            4),

    # Parentheses change sign context
    ("3-(2+1)",             0),
    ("5-(3-1)",             3),
    ("(5-(3-1))+3",         6),

    # Leading/trailing spaces
    (" 100 ",               100),
])
def test_calculate(s, want):
    assert Solution().calculate(s) == want
    assert Solution().calculate_iterative(s) == want


@pytest.mark.parametrize("s, want", [
    # LeetCode examples
    ("3+2*2",          7),
    (" 3/2 ",          1),
    (" 3+5 / 2 ",      5),

    # Single number
    ("42",             42),

    # Multiplication only
    ("2*3*4",          24),

    # Division truncates toward zero
    ("7/2",            3),
    ("14-3/2",         13),

    # Mixed precedence
    ("1+2*3-4/2",      5),
    ("2*3+4",          10),
    ("100*2+3",        203),

    # Left-to-right for equal precedence
    ("10-3-2",         5),
    ("6/2/3",          1),
])
def test_calculate2(s, want):
    assert Solution().calculate2(s) == want


@pytest.mark.parametrize("s, want", [
    # LeetCode examples
    ("1+2",                              3),
    ("6-4/2",                            4),
    ("2*(5+5*2)/3+(6/2+8)",             21),
    ("(2+6*3+5-(3*14/7+2)*5)+3",       -12),

    # Single number
    ("0",                                0),

    # Nested parentheses with mixed ops
    ("(2+3)*4",                         20),
    ("10/(2+3)",                         2),
    ("2*(3+(4*5))",                     46),
    ("((2+3))*((4-1))",                 15),

    # Parens that cancel to subtraction
    ("5-(3-1)",                          3),
    ("3*(1+2)-4/(1+1)",                  7),

    # No parens — same as calc II
    ("3+2*2",                            7),
    (" 3/2 ",                            1),
])
def test_calculate3(s, want):
    assert Solution().calculate3(s) == want
    assert Solution().calculate_template(s) == want
