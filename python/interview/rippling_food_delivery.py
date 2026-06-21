"""
Design an in-memory service for a food-delivery platform to compute payouts and basic live analytics.
Tens of thousands of drivers, each submitting hundreds of deliveries per week.
Delivery details sent immediately after completion. Drivers have hourly rates that differ by driver.
A driver may take multiple deliveries simultaneously — each delivery is paid independently (overlaps earn the sum).
Use efficient data structures.
"""

from dataclasses import dataclass, field


@dataclass
class Delivery:
    driver_id: str
    start_time: float
    end_time: float


@dataclass
class Driver:
    driver_id: str
    hourly_rate: float
    deliveries: list[Delivery] = field(default_factory=list)


class Solution:
    def __init__(self):
        self.drivers: dict[str, Driver] = {}

    def add_driver(self, driver_id, hourly_rate):
        if driver_id not in self.drivers:
            self.drivers[driver_id] = Driver(driver_id, hourly_rate)

    def record_delivery(self, driver_id, start_time, end_time):
        if driver_id in self.drivers:
            self.drivers[driver_id].deliveries.append(
                Delivery(driver_id, start_time, end_time)
            )

    def get_total_cost(self) -> float:  # total cost across all drivers
        total_cost = 0
        for _, driver in self.drivers.items():
            for delivery in driver.deliveries:
                cost = (
                    driver.hourly_rate
                    * (delivery.end_time - delivery.start_time)
                    / 3600
                )
                total_cost += cost

        return total_cost

    def pay_up_to(pay_time):  # pay all wages before this timestamp
        pass

    def total_cost_unpaid() -> float:  # remaining unpaid amount
        pass
