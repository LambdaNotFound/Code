from abc import ABC, abstractmethod
from functools import cmp_to_key


# ── Filtering ─────────────────────────────────────────────────────────────────

class TaskFilter(ABC):
    @abstractmethod
    def filter(self, tasks: list) -> list:
        pass


class CompositeFilter(TaskFilter):
    def __init__(self, filters: list[TaskFilter]):
        self._filters = filters

    def filter(self, tasks: list) -> list:
        for f in self._filters:
            tasks = f.filter(tasks)
        return tasks


class NotDoneFilter(TaskFilter):
    def filter(self, tasks: list) -> list:
        return [t for t in tasks if not t["is_done"]]


class UnassignedFilter(TaskFilter):
    def filter(self, tasks: list) -> list:
        return [t for t in tasks if t["assignee"] is None]


class DeduplicationFilter(TaskFilter):
    def filter(self, tasks: list) -> list:
        seen: set = set()
        result = []
        for t in tasks:
            key = (t["description"], t["due_date"])
            if key not in seen:
                seen.add(key)
                result.append(t)
        return result


# ── Prioritization ────────────────────────────────────────────────────────────

class TaskComparator(ABC):
    @abstractmethod
    def compare(self, a: dict, b: dict) -> int:
        pass  # negative = a first, 0 = tie/delegate, positive = b first


class CompositeComparator(TaskComparator):
    def __init__(self, comparators: list[TaskComparator]):
        self._chain = comparators

    def compare(self, a: dict, b: dict) -> int:
        for comp in self._chain:
            result = comp.compare(a, b)
            if result != 0:
                return result
        return 0


class DueDateComparator(TaskComparator):
    def compare(self, a: dict, b: dict) -> int:
        a_date, b_date = a["due_date"], b["due_date"]
        if a_date and b_date:
            if a_date < b_date:
                return -1
            if a_date > b_date:
                return 1
            return 0  # same due date — delegate to next rule
        if a_date:
            return -1  # a has due date, b doesn't — a first
        if b_date:
            return 1   # b has due date, a doesn't — b first
        return 0  # both have no due date — delegate


class HighPriorityComparator(TaskComparator):
    def compare(self, a: dict, b: dict) -> int:
        if a["is_high_priority"] == b["is_high_priority"]:
            return 0
        return -1 if a["is_high_priority"] else 1


class CreatedAtComparator(TaskComparator):
    def compare(self, a: dict, b: dict) -> int:
        return a["created_at"] - b["created_at"]


# ── TaskManager ───────────────────────────────────────────────────────────────

class TaskManager:
    def __init__(self, tasks: list):
        self._tasks = tasks
        self._filter = CompositeFilter([
            NotDoneFilter(),
            UnassignedFilter(),
            DeduplicationFilter(),
        ])
        self._comparator = CompositeComparator([
            DueDateComparator(),
            HighPriorityComparator(),
            CreatedAtComparator(),
        ])

    def print_prioritized_tasks(self):
        filtered = self._filter.filter(self._tasks)
        task_map = {t["id"]: t for t in filtered}
        sorted_tasks = sorted(filtered, key=cmp_to_key(self._comparator.compare))
        for task in sorted_tasks:
            line = f"Task ID: {task['id']}, Description: {task['description']}"
            parent_id = task.get("parent_id")
            if parent_id is not None and parent_id in task_map:
                line += f", parent: {task_map[parent_id]['description']}"
            print(line)
