"""
Unit tests for monotonic_stack.py — removeDuplicateLetters (LeetCode 316)

Goal: return the lexicographically smallest subsequence that contains every
unique letter exactly once.
"""

import pytest

from dsa.mono_stack.monotonic_stack import Solution


@pytest.mark.parametrize("s, want", [
    # LeetCode examples
    ("bcabc",    "abc"),
    ("cbacdcbc", "acdb"),

    # Single character — trivial
    ("a",        "a"),

    # All same — deduplicated to one
    ("aaaa",     "a"),

    # Already the lexicographic minimum, no duplicates
    ("abc",      "abc"),

    # Each character appears once — order is forced (can't improve)
    ("dcba",     "dcba"),

    # Later occurrences allow earlier ones to be dropped
    ("abacb",    "abc"),

    # 'b' appears only before 'a', so 'b' must precede 'a' in result
    ("bbcaac",   "bac"),

    # Two-char deduplicate
    ("ba",       "ba"),
    ("ab",       "ab"),
    ("aabb",     "ab"),

    # Longer cascade: 'c' can be dropped because it reappears later
    ("cbacbca",  "abc"),
])
def test_remove_duplicate_letters(s, want):
    assert Solution().removeDuplicateLetters(s) == want
