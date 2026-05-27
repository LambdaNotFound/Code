"""
Unit tests for perfect_squares.py — LeetCode 279: Perfect Squares
"""
import sys
import os

import pytest

sys.path.insert(0, os.path.dirname(__file__))
from perfect_squares import Solution  # noqa: E402


@pytest.mark.parametrize("n, want", [
    # LeetCode examples
    (12, 3),   # 4+4+4
    (13, 2),   # 4+9

    # Perfect squares themselves cost 1
    (1,  1),
    (4,  1),
    (9,  1),
    (16, 1),
    (25, 1),

    # Two squares
    (2,  2),   # 1+1
    (5,  2),   # 1+4
    (10, 2),   # 1+9

    # Three squares
    (3,  3),   # 1+1+1
    (6,  3),   # 1+1+4
    (11, 3),   # 1+1+9

    # Four squares (Lagrange: numbers of the form 4^a(8b+7))
    (7,  4),
    (15, 4),
])
def test_numSquares(n, want):
    assert Solution().numSquares(n) == want


def _is_perfect_square(x: int) -> bool:
    return int(x ** 0.5) ** 2 == x


@pytest.mark.parametrize("n, want_len", [
    # LeetCode examples
    (12, 3),   # e.g. [4, 4, 4]
    (13, 2),   # e.g. [4, 9]

    # Perfect squares themselves — one element equal to n
    (1,  1),
    (4,  1),
    (9,  1),
    (16, 1),
    (25, 1),

    # Two squares
    (2,  2),
    (5,  2),
    (10, 2),

    # Three squares
    (3,  3),
    (6,  3),
    (11, 3),

    # Four squares
    (7,  4),
    (15, 4),
])
def test_numSquaresWithList(n, want_len):
    result = Solution.numSquaresWithList(n)
    assert len(result) == want_len
    assert sum(result) == n
    assert all(_is_perfect_square(x) for x in result)
