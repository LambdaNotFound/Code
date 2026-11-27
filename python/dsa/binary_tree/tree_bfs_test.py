"""
Unit tests for tree_bfs.py — LeetCode 199: Binary Tree Right Side View
"""

import pytest

from dsa.binary_tree.tree_bfs import Solution
from dsa.utils import build_tree


@pytest.mark.parametrize("values, want", [
    # LeetCode examples
    ([1, 2, 3, None, 5, None, 4], [1, 3, 4]),
    ([1, None, 3],                [1, 3]),
    ([],                          []),

    # Single node
    ([1],                         [1]),

    # Left-skewed — only left children; right side sees each node
    ([1, 2, None, 3],             [1, 2, 3]),

    # Right-skewed
    ([1, None, 2, None, 3],       [1, 2, 3]),

    # Complete tree — rightmost node at each level
    ([1, 2, 3, 4, 5, 6, 7],      [1, 3, 7]),

    # Left subtree deeper than right — deepest left node visible
    ([1, 2, 3, 4],                [1, 3, 4]),
])
def test_rightSideView(values, want):
    root = build_tree(values)
    assert Solution().rightSideView(root) == want
