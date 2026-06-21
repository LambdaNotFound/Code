"""
You are building a driver payment system for a food delivery company.
The accounting team needs to track how much money is owed to drivers and display this on a live dashboard.
"""

from dataclasses import dataclass, field
from datetime import datetime
from math import ceil


@dataclass
class Delivery:
    driver_id: str
    start_time: datetime
    end_time: datetime


@dataclass
class Driver:
    driver_id: str
    hourly_rate: float
    deliveries: list[Delivery] = field(default_factory=list)


class Solution:
    def __init__(self):
        self.drivers: dict[str, Driver] = {}
        self.total_cost = 0

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

    def add_driver(self, driver_id: str, hourly_rate: float):
        if driver_id not in self.drivers:
            self.drivers[driver_id] = Driver(driver_id, hourly_rate)

    def record_delivery(self, driver_id: str, start_time: str, end_time: str):
        if driver_id in self.drivers:
            dt_start = datetime.fromisoformat(start_time)
            dt_end = datetime.fromisoformat(end_time)
            self.drivers[driver_id].deliveries.append(
                Delivery(driver_id, dt_start, dt_end)
            )
            hours = ceil((dt_end - dt_start).total_seconds() / 3600)
            self.total_cost += self.drivers[driver_id].hourly_rate * hours

    def get_total_cost(self) -> float:  # total cost across all drivers
        return self.total_cost

    def pay_up_to(pay_time: str):  # pay all wages before this timestamp
        pass

    def total_cost_unpaid() -> float:  # remaining unpaid amount
        pass
