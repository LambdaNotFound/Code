"""
Unit tests for path_sum.py — Path Sum III (LeetCode 437)
"""
import sys
import os

import pytest

sys.path.insert(0, os.path.dirname(os.path.dirname(__file__)))
sys.path.insert(0, os.path.dirname(__file__))
from common import build_tree
from path_sum import Solution


@pytest.mark.parametrize("values, targetSum, want", [
    # LeetCode examples
    ([10, 5, -3, 3, 2, None, 11, 3, -2, None, 1], 8,  3),
    ([5, 4, 8, 11, None, 13, 4, 7, 2, None, None, 5, 1], 22, 3),

    # Empty tree
    ([], 0, 0),

    # Single node equals target
    ([5], 5, 1),

    # Single node does not equal target
    ([1], 5, 0),

    # Root-only, target is zero
    ([0], 0, 1),

    # Two paths: [1,2] and [3]
    ([1, 2, 3], 3,   2),

    # Only [1]: [-2,3] sums to 1 but is not a downward path from root
    ([1, -2, 3], 1,  1),

    # Target not reachable
    ([1, 2, 3], 100, 0),

    # Negative target: [-1,-2] and [-3]
    ([-1, -2, -3], -3, 2),

    # All zeros: single nodes (3) + two-node paths root->left, root->right (2)
    ([0, 0, 0], 0, 5),
])
def test_path_sum(values, targetSum, want):
    root = build_tree(values)
    assert Solution().pathSum(root, targetSum) == want
