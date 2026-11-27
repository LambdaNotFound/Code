"""
Unit tests for merge_k_lists.py — LeetCode 23 (Merge k Sorted Lists)
"""
import pytest

from dsa.heap.merge_k_lists import Solution
from dsa.utils import build_list, list_to_vals


@pytest.mark.parametrize("lists_vals, want", [
    # LeetCode examples
    ([[1, 4, 5], [1, 3, 4], [2, 6]],  [1, 1, 2, 3, 4, 4, 5, 6]),
    ([],                               []),
    ([[]],                             []),

    # Single list
    ([[1, 2, 3]],                      [1, 2, 3]),

    # Two lists, one empty
    ([[1, 3], []],                     [1, 3]),

    # Already merged order
    ([[1], [2], [3]],                  [1, 2, 3]),

    # All same values
    ([[1, 1], [1], [1, 1]],           [1, 1, 1, 1, 1]),

    # Negative values
    ([[-3, -1], [-2, 0]],             [-3, -2, -1, 0]),
])
def test_mergeKLists(lists_vals, want):
    lists = [build_list(v) for v in lists_vals]
    result = Solution().mergeKLists(lists)
    assert list_to_vals(result) == want
