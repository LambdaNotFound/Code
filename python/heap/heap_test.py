"""
Unit tests for heap/ — LeetCode 23 (Merge k Sorted Lists) and 347 (Top K Frequent Elements)
"""
import sys
import os

import pytest

sys.path.insert(0, os.path.dirname(os.path.dirname(__file__)))
from common import build_list, list_to_vals
sys.path.insert(0, os.path.dirname(__file__))
from merge_k_lists import Solution as MergeKSolution
from top_k_frequent import Solution as TopKSolution


# ---------------------------------------------------------------------------
# 23. Merge k Sorted Lists
# ---------------------------------------------------------------------------

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
    result = MergeKSolution().mergeKLists(lists)
    assert list_to_vals(result) == want


# ---------------------------------------------------------------------------
# 347. Top K Frequent Elements
# ---------------------------------------------------------------------------

@pytest.mark.parametrize("nums, k, want", [
    # LeetCode examples
    ([1, 1, 1, 2, 2, 3],  2, {1, 2}),
    ([1],                  1, {1}),

    # All same frequency — any k elements are valid; use a set to check containment
    ([1, 2, 3],            2, {1, 2, 3}),   # any 2 of the 3

    # Tie-breaking — only the top-1 most frequent matters
    ([4, 4, 4, 2, 2, 1],  1, {4}),

    # Negative numbers
    ([-1, -1, 2, 2, 2],   1, {2}),

    # k equals total unique elements
    ([5, 5, 4, 4, 3, 3],  3, {3, 4, 5}),
])
def test_topKFrequent(nums, k, want):
    result_set = set(TopKSolution().topKFrequent(nums, k))
    assert len(result_set) == k
    # All returned elements must be in the valid answer set
    assert result_set <= want


@pytest.mark.parametrize("nums, k, want", [
    ([1, 1, 1, 2, 2, 3],  2, {1, 2}),
    ([1],                  1, {1}),
    ([4, 4, 4, 2, 2, 1],  1, {4}),
    ([-1, -1, 2, 2, 2],   1, {2}),
    ([5, 5, 4, 4, 3, 3],  3, {3, 4, 5}),
])
def test_topKFrequentClaude(nums, k, want):
    result_set = set(TopKSolution().topKFrequentClaude(nums, k))
    assert len(result_set) == k
    assert result_set <= want
