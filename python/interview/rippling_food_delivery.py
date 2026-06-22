"""
You are building a driver payment system for a food delivery company.
The accounting team needs to track how much money is owed to drivers and display this on a live dashboard.
"""

import bisect
from dataclasses import dataclass, field
from datetime import datetime
from math import ceil


@dataclass
class Delivery:
    driver_id: str
    start_time: datetime
    end_time: datetime
    cost: float = 0.0


@dataclass
class Driver:
    driver_id: str
    hourly_rate: float
    deliveries: list[Delivery] = field(default_factory=list)


class Solution:
    def __init__(self):
        self.drivers: dict[str, Driver] = {}
        self.total_cost = 0
        self.sorted_deliveries: list[Delivery] = []  # sorted by end_time
        self.paid_index: int = 0    # first unpaid delivery in sorted_deliveries
        self.paid_total: float = 0.0

    """
    Part 1:

    Implement the following functions:

    sign_up_driver(driver_id, hourly_rate)
    Registers a new driver with their unique ID
    Associates the driver with an hourly rate in USD
    record_delivery(driver_id, start_time, end_time)
    Records a completed delivery
    Start and end times have at least second-level precision
    Drivers are paid their full hourly rate for each delivery (regardless of duration)
    total_cost()
    Returns the total cost of all deliveries completed across all drivers

    Note: The choice of time format is an important design decision to discuss before implementation.
    """

    # T: O(1)  S: O(d) total across all drivers
    def add_driver(self, driver_id: str, hourly_rate: float):
        if driver_id not in self.drivers:
            self.drivers[driver_id] = Driver(driver_id, hourly_rate)

    # T: O(n) — bisect.insort is O(log n) search but O(n) list shift
    # S: O(1) per call; O(n) total across all deliveries in sorted_deliveries
    def record_delivery(self, driver_id: str, start_time: str, end_time: str):
        if driver_id in self.drivers:
            dt_start = datetime.fromisoformat(start_time)
            dt_end = datetime.fromisoformat(end_time)
            hours = ceil((dt_end - dt_start).total_seconds() / 3600)
            cost = self.drivers[driver_id].hourly_rate * hours
            delivery = Delivery(driver_id, dt_start, dt_end, cost)
            self.drivers[driver_id].deliveries.append(delivery)
            bisect.insort(self.sorted_deliveries, delivery, key=lambda d: d.end_time)
            self.total_cost += cost

    # T: O(1) — total_cost maintained incrementally
    def get_total_cost(self) -> float:
        return self.total_cost

    """
    Part 2: Pay-up-to and unpaid balance
    Given a specific timestamp, mark all wages earned before that point as paid. 
    Then total_cost_unpaid() returns only wages from deliveries after the last pay cutoff.
    """

    # T: O(k) amortized — each delivery is visited at most once across all pay_up_to calls
    def pay_up_to(self, pay_time: str):
        cutoff = datetime.fromisoformat(pay_time)
        while (self.paid_index < len(self.sorted_deliveries) and
               self.sorted_deliveries[self.paid_index].end_time <= cutoff):
            self.paid_total += self.sorted_deliveries[self.paid_index].cost
            self.paid_index += 1

    # T: O(1)  S: O(1)
    def total_cost_unpaid(self) -> float:
        return self.total_cost - self.paid_total

    """
    Part 3: Max simultaneous drivers (open-ended)
    same driver can have multiple concurrent deliveries
    """

    # T: O(n log n) — sort 2n events; O(n) sweep  S: O(n) for events + driver_active
    def max_simultaneous_drivers(self) -> int:
        events = []
        for delivery in self.sorted_deliveries:
            events.append((delivery.start_time, +1, delivery.driver_id))
            events.append((delivery.end_time, -1, delivery.driver_id))

        # at equal time, process ends (-1) before starts (+1)
        events.sort(key=lambda e: (e[0], e[1]))

        driver_active: dict[str, int] = {}  # driver_id → count of active deliveries
        active_drivers = 0
        max_drivers = 0

        for _, delta, driver_id in events:
            prev = driver_active.get(driver_id, 0)
            driver_active[driver_id] = prev + delta
            if prev == 0 and delta == +1:
                active_drivers += 1
            elif prev == 1 and delta == -1:
                active_drivers -= 1
            max_drivers = max(max_drivers, active_drivers)

        return max_drivers
