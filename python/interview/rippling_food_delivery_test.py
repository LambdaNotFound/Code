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


def test_pay_up_to(sol):
    sol.add_driver("driver_1", 20.0)
    sol.add_driver("driver_2", 25.0)

    # end_times: 09:15, 10:30, 11:45
    sol.record_delivery("driver_1", "2024-01-15T10:00:00", "2024-01-15T10:30:00")  # $20
    sol.record_delivery("driver_1", "2024-01-15T11:00:00", "2024-01-15T11:45:00")  # $20
    sol.record_delivery("driver_2", "2024-01-15T09:00:00", "2024-01-15T09:15:00")  # $25

    # before any payment, all $65 is unpaid
    assert sol.total_cost_unpaid() == 65.0

    # pay up to 10:00 — only driver_2's delivery (ends 09:15) is paid
    sol.pay_up_to("2024-01-15T10:00:00")
    assert sol.total_cost_unpaid() == 40.0  # driver_1's two deliveries remain

    # pay up to 11:00 — driver_1's first delivery (ends 10:30) also paid
    sol.pay_up_to("2024-01-15T11:00:00")
    assert sol.total_cost_unpaid() == 20.0  # only driver_1's second delivery remains

    # pay up to end — nothing unpaid
    sol.pay_up_to("2024-01-15T12:00:00")
    assert sol.total_cost_unpaid() == 0.0


def test_max_simultaneous_drivers(sol):
    sol.add_driver("driver_1", 20.0)
    sol.add_driver("driver_2", 25.0)
    sol.add_driver("driver_3", 15.0)

    # 10:00-11:00: driver_1 only → 1
    # 10:30-11:30: driver_2 overlaps driver_1 → 2
    # 10:45-11:15: driver_3 overlaps both → 3 (peak)
    # driver_1 also has a second concurrent delivery 10:15-10:50 → still counts as 1 driver
    sol.record_delivery("driver_1", "2024-01-15T10:00:00", "2024-01-15T11:00:00")
    sol.record_delivery("driver_1", "2024-01-15T10:15:00", "2024-01-15T10:50:00")  # concurrent, same driver
    sol.record_delivery("driver_2", "2024-01-15T10:30:00", "2024-01-15T11:30:00")
    sol.record_delivery("driver_3", "2024-01-15T10:45:00", "2024-01-15T11:15:00")

    assert sol.max_simultaneous_drivers() == 3
