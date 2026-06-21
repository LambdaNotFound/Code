import sys, os
sys.path.insert(0, os.path.dirname(__file__))
from rippling_food_delivery import Solution

import pytest


@pytest.fixture
def sol():
    return Solution()


def test_basic_flow(sol):
    sol.add_driver("driver_1", 20.0)
    sol.add_driver("driver_2", 25.0)

    sol.record_delivery("driver_1", "2024-01-15T10:00:00", "2024-01-15T10:30:00")
    sol.record_delivery("driver_1", "2024-01-15T11:00:00", "2024-01-15T11:45:00")
    sol.record_delivery("driver_2", "2024-01-15T09:00:00", "2024-01-15T09:15:00")

    # driver_1: 2 deliveries × $20.00 = $40.00
    # driver_2: 1 delivery  × $25.00 = $25.00
    assert sol.get_total_cost() == 65.0


def test_multiple_drivers(sol):
    sol.add_driver("driver_a", 15.0)
    sol.add_driver("driver_b", 30.0)
    sol.add_driver("driver_c", 22.5)

    sol.record_delivery("driver_a", "2024-01-15T08:00:00", "2024-01-15T08:30:00")
    sol.record_delivery("driver_b", "2024-01-15T09:00:00", "2024-01-15T09:20:00")
    sol.record_delivery("driver_c", "2024-01-15T10:00:00", "2024-01-15T10:45:00")
    sol.record_delivery("driver_a", "2024-01-15T11:00:00", "2024-01-15T11:15:00")

    # driver_a: 2 deliveries × $15.00 = $30.00
    # driver_b: 1 delivery  × $30.00 = $30.00
    # driver_c: 1 delivery  × $22.50 = $22.50
    assert sol.get_total_cost() == 82.5
