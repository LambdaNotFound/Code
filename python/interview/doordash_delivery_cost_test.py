import os
import sys
from datetime import datetime, timezone

import pytest

sys.path.insert(0, os.path.dirname(__file__))
from doordash_delivery_cost import Solution


def _epoch(iso: str) -> int:
    return int(datetime.fromisoformat(iso).replace(tzinfo=timezone.utc).timestamp())


@pytest.fixture
def solution():
    sol = Solution()
    sol.sign_up_driver("d1", 2000)  # $20.00/hr
    sol.sign_up_driver("d2", 3000)  # $30.00/hr
    return sol


class TestSingleDelivery:
    def test_one_hour_delivery(self, solution):
        solution.record_delivery("d1", _epoch("2026-03-16 12:00"), _epoch("2026-03-16 13:00"))
        assert solution.total_cost() == 2000

    def test_fractional_hour_delivery(self, solution):
        solution.record_delivery("d1", _epoch("2026-03-16 12:00"), _epoch("2026-03-16 12:30"))
        assert solution.total_cost() == 1000

    def test_zero_duration_delivery(self, solution):
        solution.record_delivery("d1", _epoch("2026-03-16 12:00"), _epoch("2026-03-16 12:00"))
        assert solution.total_cost() == 0

    def test_overnight_delivery_spans_midnight(self, solution):
        solution.record_delivery("d1", _epoch("2026-03-16 23:30"), _epoch("2026-03-17 00:30"))
        assert solution.total_cost() == 2000


class TestMultipleDeliveries:
    def test_same_driver_deliveries_are_summed(self, solution):
        solution.record_delivery("d1", _epoch("2026-03-16 09:00"), _epoch("2026-03-16 10:00"))
        solution.record_delivery("d1", _epoch("2026-03-16 11:00"), _epoch("2026-03-16 11:30"))
        assert solution.total_cost() == 3000

    def test_different_drivers_are_summed_together(self, solution):
        solution.record_delivery("d1", _epoch("2026-03-16 09:00"), _epoch("2026-03-16 10:00"))  # 2000
        solution.record_delivery("d2", _epoch("2026-03-16 09:00"), _epoch("2026-03-16 11:00"))  # 6000
        assert solution.total_cost() == 8000


class TestTotalCostBetween:
    def test_delivery_fully_inside_window_is_included(self, solution):
        solution.record_delivery("d1", _epoch("2026-03-16 09:00"), _epoch("2026-03-16 10:00"))
        cost = solution.total_cost_between(_epoch("2026-03-16 08:00"), _epoch("2026-03-16 12:00"))
        assert cost == 2000

    def test_delivery_fully_outside_window_is_excluded(self, solution):
        solution.record_delivery("d1", _epoch("2026-03-16 09:00"), _epoch("2026-03-16 10:00"))
        cost = solution.total_cost_between(_epoch("2026-03-16 11:00"), _epoch("2026-03-16 12:00"))
        assert cost == 0

    def test_delivery_partially_overlapping_window_is_excluded(self, solution):
        # starts before the window, ends inside it -> not fully contained
        solution.record_delivery("d1", _epoch("2026-03-16 07:30"), _epoch("2026-03-16 08:30"))
        cost = solution.total_cost_between(_epoch("2026-03-16 08:00"), _epoch("2026-03-16 12:00"))
        assert cost == 0

    def test_delivery_exactly_on_window_boundary_is_included(self, solution):
        solution.record_delivery("d1", _epoch("2026-03-16 08:00"), _epoch("2026-03-16 12:00"))
        cost = solution.total_cost_between(_epoch("2026-03-16 08:00"), _epoch("2026-03-16 12:00"))
        assert cost == 8000

    def test_only_deliveries_within_window_are_summed(self, solution):
        solution.record_delivery("d1", _epoch("2026-03-16 09:00"), _epoch("2026-03-16 10:00"))  # inside, 2000
        solution.record_delivery("d2", _epoch("2026-03-16 20:00"), _epoch("2026-03-16 21:00"))  # outside
        cost = solution.total_cost_between(_epoch("2026-03-16 08:00"), _epoch("2026-03-16 12:00"))
        assert cost == 2000


class TestTotalCostBetweenSorted:
    # These mirror TestTotalCostBetween's cases, calling total_cost_between_sorted
    # instead. All deliveries below are recorded in non-decreasing start_time
    # order, satisfying the method's sortedness precondition.

    def test_delivery_fully_inside_window_is_included(self, solution):
        solution.record_delivery("d1", _epoch("2026-03-16 09:00"), _epoch("2026-03-16 10:00"))
        cost = solution.total_cost_between_sorted(_epoch("2026-03-16 08:00"), _epoch("2026-03-16 12:00"))
        assert cost == 2000

    def test_delivery_fully_outside_window_is_excluded(self, solution):
        solution.record_delivery("d1", _epoch("2026-03-16 09:00"), _epoch("2026-03-16 10:00"))
        cost = solution.total_cost_between_sorted(_epoch("2026-03-16 11:00"), _epoch("2026-03-16 12:00"))
        assert cost == 0

    def test_delivery_partially_overlapping_window_is_excluded(self, solution):
        # starts before the window, ends inside it -> not fully contained
        solution.record_delivery("d1", _epoch("2026-03-16 07:30"), _epoch("2026-03-16 08:30"))
        cost = solution.total_cost_between_sorted(_epoch("2026-03-16 08:00"), _epoch("2026-03-16 12:00"))
        assert cost == 0

    def test_delivery_exactly_on_window_boundary_is_included(self, solution):
        solution.record_delivery("d1", _epoch("2026-03-16 08:00"), _epoch("2026-03-16 12:00"))
        cost = solution.total_cost_between_sorted(_epoch("2026-03-16 08:00"), _epoch("2026-03-16 12:00"))
        assert cost == 8000

    def test_only_deliveries_within_window_are_summed(self, solution):
        solution.record_delivery("d1", _epoch("2026-03-16 09:00"), _epoch("2026-03-16 10:00"))  # inside, 2000
        solution.record_delivery("d2", _epoch("2026-03-16 20:00"), _epoch("2026-03-16 21:00"))  # outside
        cost = solution.total_cost_between_sorted(_epoch("2026-03-16 08:00"), _epoch("2026-03-16 12:00"))
        assert cost == 2000

    def test_agrees_with_unsorted_variant_on_sorted_input(self, solution):
        solution.record_delivery("d1", _epoch("2026-03-16 07:00"), _epoch("2026-03-16 07:45"))
        solution.record_delivery("d2", _epoch("2026-03-16 09:00"), _epoch("2026-03-16 11:00"))
        solution.record_delivery("d1", _epoch("2026-03-16 12:00"), _epoch("2026-03-16 13:15"))
        window = (_epoch("2026-03-16 08:00"), _epoch("2026-03-16 12:30"))
        assert solution.total_cost_between_sorted(*window) == solution.total_cost_between(*window)


class TestEdgeCases:
    def test_no_deliveries_returns_zero(self, solution):
        assert solution.total_cost() == 0

    def test_unrecorded_driver_id_raises_key_error(self, solution):
        solution.record_delivery("unknown", _epoch("2026-03-16 09:00"), _epoch("2026-03-16 10:00"))
        with pytest.raises(KeyError):
            solution.total_cost()


class TestRounding:
    def test_fractional_cent_below_half_rounds_down(self, solution):
        # 44 seconds at 2000 cents/hr = 24.444... cents -> rounds down to 24
        solution.record_delivery("d1", _epoch("2026-03-16 09:00:00"), _epoch("2026-03-16 09:00:44"))
        assert solution.total_cost() == 24

    def test_fractional_cent_above_half_rounds_up(self, solution):
        # 61 seconds at 2000 cents/hr = 33.888... cents -> rounds up to 34
        solution.record_delivery("d1", _epoch("2026-03-16 09:00:00"), _epoch("2026-03-16 09:01:01"))
        assert solution.total_cost() == 34

    def test_exact_half_cent_rounds_up(self, solution):
        solution.sign_up_driver("d3", 2001)  # $20.01/hr
        # 30 minutes at 2001 cents/hr = 1000.5 cents exactly -> rounds up to 1001
        solution.record_delivery("d3", _epoch("2026-03-16 09:00:00"), _epoch("2026-03-16 09:30:00"))
        assert solution.total_cost() == 1001
