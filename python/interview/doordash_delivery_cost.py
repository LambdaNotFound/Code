import bisect
from dataclasses import dataclass
from datetime import datetime


@dataclass
class Driver:
    driver_id: str
    hourly_rate_cents: int


@dataclass
class Delivery:
    driver_id: str
    start_time: datetime
    end_time: datetime


# Design:
#   _drivers      dict[str, Driver]     — driver_id -> Driver, O(1) rate lookup
#   _deliveries   list[Delivery]        — all deliveries in insertion order
#                                          (not sorted by start/end time)
#
# Key decision: deliveries are kept as a flat list. Since
# the list isn't time-sorted, total_cost_between can't early-exit and must
# scan every delivery.
#
# Assumption: total_cost_between includes a delivery only if it falls fully
# within [start, end] (both start_time and end_time inside the window) —
# partially-overlapping deliveries are not prorated. Flagged as a design
# choice to confirm with the interviewer per the open-ended prompt.
class Solution:
    def __init__(self):  # T: O(1), S: O(1)
        self._drivers: dict[str, Driver] = {}
        self._deliveries: list[Delivery] = []

    def sign_up_driver(
        self, driver_id: str, hourly_rate: int
    ) -> None:  # T: O(1), S: O(1)
        self._drivers[driver_id] = Driver(driver_id, hourly_rate)

    def record_delivery(
        self, driver_id: str, start_time: str, end_time: str
    ) -> None:  # T: O(1), S: O(1)
        start = datetime.fromisoformat(start_time)
        end = datetime.fromisoformat(end_time)
        self._deliveries.append(Delivery(driver_id, start, end))

    def total_cost(self) -> int:  # T: O(n), S: O(1)
        return sum(self._pay(delivery) for delivery in self._deliveries)

    def total_cost_between(self, start: str, end: str) -> int:  # T: O(n), S: O(1)
        start_dt = datetime.fromisoformat(start)
        end_dt = datetime.fromisoformat(end)
        return sum(
            self._pay(delivery)
            for delivery in self._deliveries
            if delivery.start_time >= start_dt and delivery.end_time <= end_dt
        )

    def total_cost_between_sorted(
        self, start: str, end: str
    ) -> int:  # T: O(log n + k), S: O(1)
        # Precondition: _deliveries must be sorted by start_time, which holds
        # only if record_delivery is called with non-decreasing start_time
        # (a realistic assumption since delivery events are recorded as they
        # happen). If violated, bisect's results are undefined.
        #
        # Binary search narrows to the contiguous window [lo, hi) where
        # start_dt <= start_time <= end_dt; sorting by a single key can't
        # resolve the end_time <= end_dt half of containment, so that's
        # checked with a linear scan over just that window. k = hi - lo,
        # so this beats total_cost_between's O(n) when the window is narrow.
        start_dt = datetime.fromisoformat(start)
        end_dt = datetime.fromisoformat(end)
        lo = bisect.bisect_left(self._deliveries, start_dt, key=lambda d: d.start_time)
        hi = bisect.bisect_right(self._deliveries, end_dt, key=lambda d: d.start_time)
        total = 0
        for delivery in self._deliveries[lo:hi]:
            if delivery.end_time <= end_dt:
                total += self._pay(delivery)
        return total

    def _pay(self, delivery: Delivery) -> int:  # T: O(1), S: O(1)
        # Integer-only arithmetic throughout (duration in whole seconds,
        # rate in whole cents) so pay is exact — no float rounding error.
        # Fractional cents are rounded half-up (biasing the numerator by
        # half the divisor before floor-dividing), not truncated.
        rate_cents = self._drivers[delivery.driver_id].hourly_rate_cents
        duration = delivery.end_time - delivery.start_time
        duration_seconds = duration.days * 86400 + duration.seconds
        return (duration_seconds * rate_cents + 1800) // 3600  # ROUND_HALF_UP
