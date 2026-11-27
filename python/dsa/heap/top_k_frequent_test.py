"""
Unit tests for top_k_frequent.py — LeetCode 347 (Top K Frequent Elements)
"""
import pytest

from dsa.heap.top_k_frequent import Solution


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
    result_set = set(Solution().topKFrequent(nums, k))
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
    result_set = set(Solution().topKFrequentClaude(nums, k))
    assert len(result_set) == k
    assert result_set <= want
