"""
Unit tests for sensor_coverage.py — LeetCode 475: Heaters
"""

import pytest

from dsa.binary_search.sensor_coverage import Solution


@pytest.mark.parametrize("houses, heaters, want", [
    # LC example 1: one heater in the middle
    ([1, 2, 3],         [2],        1),

    # LC example 2: heaters at both ends
    ([1, 2, 3, 4],      [1, 4],     1),

    # LC example 3: only one heater, far from the last house
    ([1, 5],            [2],        3),

    # house sits exactly on a heater
    ([1],               [1],        0),

    # all houses cluster around a single heater
    ([1, 2, 3],         [2],        1),

    # heater to the right of all houses: farthest house (1) is 9 away
    ([1, 2, 3],         [10],       9),

    # heater to the left of all houses: farthest house (7) is 6 away
    ([5, 6, 7],         [1],        6),

    # multiple heaters: each house is at most 1 away from heater 2 or 6
    ([1, 3, 5, 7],      [2, 6],     1),

    # heaters unsorted in input: must still produce correct result
    ([1, 2, 3, 4],      [4, 1],     1),

    # single house, multiple heaters
    ([5],               [1, 3, 7],  2),
])
def test_findRadius(houses, heaters, want):
    assert Solution.findRadius(houses, heaters) == want
