"""
Unit tests for array.py — can_arrange (team photo) and firstMissingPositive (LeetCode 41)
"""
import importlib.util
import os

import pytest

# array.py conflicts with stdlib; load via spec
_spec = importlib.util.spec_from_file_location(
    "array_sol", os.path.join(os.path.dirname(__file__), "array.py")
)
_mod = importlib.util.module_from_spec(_spec)
_spec.loader.exec_module(_mod)
Solution = _mod.Solution


@pytest.mark.parametrize("h1, h2, want", [
    # heights1 is shorter — it goes in front
    ([1, 3],    [2, 4, 5],  True),   # 1<2, 3<4 — matchable
    ([5, 6],    [1, 2, 3],  False),  # front taller than every back player
    ([3],       [1, 4, 2],  True),   # single front player, one back covers it
    # heights2 is shorter — it goes in front
    ([2, 4, 5], [1, 3],     True),   # symmetric of first case
    ([1, 2, 3], [5, 6],     False),  # front [5,6] > back [1,2,3]
    # equal size — either team may lead
    ([1, 2],    [3, 4],     True),   # heights1 as front works
    ([3, 4],    [1, 2],     True),   # heights2 as front works
    ([1, 3],    [2, 3],     False),  # neither order allows all strict matches
    ([1, 2],    [1, 2],     False),  # equal heights — strict inequality fails
    ([2],       [1],        True),   # single element — heights2 leads as front
    ([1],       [2],        True),   # single element — heights1 leads as front
    # duplicates within a team
    ([1, 1],    [2, 3, 4],  True),   # two 1s both covered by distinct back players
    ([3, 3],    [1, 2, 4],  False),  # only one back player (4) > 3, but two front players need coverage
])
def test_can_arrange(h1, h2, want):
    assert Solution().can_arrange(h1, h2) == want


@pytest.mark.parametrize("nums, want", [
    # smallest positive entirely absent
    ([1, 2, 0],          3),
    # 1 is missing, negative present
    ([3, 4, -1, 1],      2),
    # all values exceed n → answer is 1
    ([7, 8, 9, 11, 12],  1),
    # single missing beyond the array
    ([1],                2),
    # complete [1..n] → n+1
    ([1, 2, 3],          4),
    # two elements in wrong order
    ([2, 1],             3),
    # duplicates, 2 missing
    ([1, 1],             2),
    # all negative → 1
    ([-1, -2],           1),
    # zero only → 1
    ([0],                1),
    # 1 absent, only 2 present
    ([2],                1),
    # longer shuffled range
    ([5, 3, 2, 1, 4],   6),
    # large gap after 1
    ([1, 100],           2),
])
def test_first_missing_positive(nums, want):
    assert Solution().firstMissingPositive(nums) == want


@pytest.mark.parametrize("nums, target, want", [
    ([2, 7, 11, 15],           9,  [0, 1]),
    ([3, 2, 4],                6,  [1, 2]),
    ([3, 3],                   6,  [0, 1]),
    ([-1, -2, -3, -4, -5],     -8, [2, 4]),
    ([0, 4, 3, 0],              0, [0, 3]),
    ([1, 2, 3, 4, 5],           9, [3, 4]),
])
def test_two_sum(nums, target, want):
    assert Solution().twoSum(nums, target) == want
