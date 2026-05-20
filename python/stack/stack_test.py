"""
Unit tests for stack.py — removeDuplicates / removeDuplicatesClaude (LeetCode 1209)

Both methods must produce identical output; we parametrize over both.
"""
import sys
import os

import pytest

sys.path.insert(0, os.path.dirname(__file__))
from stack import Solution  # noqa: E402


METHODS = ["removeDuplicates", "removeDuplicatesClaude"]

@pytest.mark.parametrize("method", METHODS)
@pytest.mark.parametrize("s, k, want", [
    # cascade: eee removed → bbb merged → ddd removed → "aa"
    ("deeedbbcccbdaa",      3,  "aa"),
    # interleaved pairs all collapse
    ("pbbcggttciiippooaais", 2,  "ps"),
    # no duplicates of size k — unchanged
    ("abcd",                2,  "abcd"),
    # all pairs removed
    ("aabbcc",              2,  ""),
    # k=1 removes every character
    ("abc",                 1,  ""),
    # single character, no removal
    ("a",                   2,  "a"),
    # cascading triples
    ("aaabcccbbddd",        3,  ""),
    # partial match: two b's remain after removing ccc
    ("abbccc",              3,  "abb"),
    # only one type of char, exact multiple
    ("aaaa",               2,  ""),
    # only one type of char, not a multiple
    ("aaaaa",              2,  "a"),
])
def test_remove_duplicates(method, s, k, want):
    fn = getattr(Solution(), method)
    assert fn(s, k) == want
