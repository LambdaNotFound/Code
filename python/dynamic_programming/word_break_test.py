"""
Unit tests for word_break.py — LeetCode 140: Word Break II
"""
import sys
import os

import pytest

sys.path.insert(0, os.path.dirname(__file__))
from word_break import Solution  # noqa: E402


@pytest.mark.parametrize("s, wordDict, want", [
    # LeetCode examples
    (
        "catsanddog",
        ["cat", "cats", "and", "sand", "dog"],
        ["cats and dog", "cat sand dog"],
    ),
    (
        "pineapplepenapple",
        ["apple", "pen", "applepen", "pine", "pineapple"],
        ["pine apple pen apple", "pineapple pen apple", "pine applepen apple"],
    ),
    (
        "catsandog",
        ["cats", "dog", "sand", "and", "cat"],
        [],
    ),

    # Single word that is itself in the dict
    ("apple", ["apple"], ["apple"]),

    # Word reused multiple times
    ("aaa", ["a", "aa"], ["a a a", "a aa", "aa a"]),

    # No valid segmentation
    ("abc", ["a", "b"], []),

    # Empty-ish: single character not in dict
    ("z", ["a"], []),

    # Single character in dict
    ("a", ["a"], ["a"]),
])
def test_wordBreak(s, wordDict, want):
    result = Solution().wordBreak(s, wordDict)
    assert sorted(result) == sorted(want)
