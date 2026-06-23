import sys, os
sys.path.insert(0, os.path.dirname(__file__))
from rippling_task_manager import (
    NotDoneFilter, UnassignedFilter, DeduplicationFilter, CompositeFilter,
    DueDateComparator, HighPriorityComparator, CreatedAtComparator, CompositeComparator,
    TaskManager,
)

import pytest


# ── Helpers ───────────────────────────────────────────────────────────────────

def make_task(id, description="task", is_done=False, due_date=None,
              parent_id=None, assignee=None, is_high_priority=False, created_at=0):
    return {
        "id": id,
        "description": description,
        "is_done": is_done,
        "due_date": due_date,
        "parent_id": parent_id,
        "assignee": assignee,
        "is_high_priority": is_high_priority,
        "created_at": created_at,
    }


def ids(tasks: list) -> list:
    return [t["id"] for t in tasks]


def printed_lines(capsys) -> list[str]:
    return capsys.readouterr().out.strip().splitlines()


# ── TestFilters ───────────────────────────────────────────────────────────────

class TestNotDoneFilter:
    def test_excludes_done_tasks(self):
        tasks = [make_task(1, is_done=True), make_task(2, is_done=False)]
        assert ids(NotDoneFilter().filter(tasks)) == [2]

    def test_all_done_returns_empty(self):
        tasks = [make_task(1, is_done=True), make_task(2, is_done=True)]
        assert NotDoneFilter().filter(tasks) == []

    def test_none_done_returns_all(self):
        tasks = [make_task(1), make_task(2)]
        assert ids(NotDoneFilter().filter(tasks)) == [1, 2]


class TestUnassignedFilter:
    def test_excludes_assigned_tasks(self):
        tasks = [make_task(1, assignee="alice"), make_task(2)]
        assert ids(UnassignedFilter().filter(tasks)) == [2]

    def test_all_assigned_returns_empty(self):
        tasks = [make_task(1, assignee="alice"), make_task(2, assignee="bob")]
        assert UnassignedFilter().filter(tasks) == []


class TestDeduplicationFilter:
    def test_keeps_first_seen_duplicate(self):
        tasks = [
            make_task(1, description="Fix bug", due_date="2024-01-01"),
            make_task(2, description="Fix bug", due_date="2024-01-01"),
        ]
        result = DeduplicationFilter().filter(tasks)
        assert ids(result) == [1]

    def test_same_description_different_due_date_kept(self):
        tasks = [
            make_task(1, description="Fix bug", due_date="2024-01-01"),
            make_task(2, description="Fix bug", due_date="2024-01-02"),
        ]
        assert ids(DeduplicationFilter().filter(tasks)) == [1, 2]

    def test_same_description_one_none_due_date_kept(self):
        tasks = [
            make_task(1, description="Fix bug", due_date=None),
            make_task(2, description="Fix bug", due_date=None),
        ]
        assert ids(DeduplicationFilter().filter(tasks)) == [1]

    def test_no_duplicates_returns_all(self):
        tasks = [make_task(1, description="A"), make_task(2, description="B")]
        assert ids(DeduplicationFilter().filter(tasks)) == [1, 2]


class TestCompositeFilter:
    def test_chains_all_filters(self):
        tasks = [
            make_task(1, is_done=True),           # removed by NotDoneFilter
            make_task(2, assignee="alice"),         # removed by UnassignedFilter
            make_task(3, description="X", due_date="2024-01-01"),
            make_task(4, description="X", due_date="2024-01-01"),  # dupe of 3
        ]
        result = CompositeFilter([
            NotDoneFilter(), UnassignedFilter(), DeduplicationFilter()
        ]).filter(tasks)
        assert ids(result) == [3]


# ── TestComparators ───────────────────────────────────────────────────────────

class TestDueDateComparator:
    cmp = DueDateComparator()

    def test_earlier_due_date_comes_first(self):
        a = make_task(1, due_date="2024-01-01")
        b = make_task(2, due_date="2024-06-01")
        assert self.cmp.compare(a, b) < 0

    def test_task_with_due_date_before_task_without(self):
        a = make_task(1, due_date="2030-01-01")
        b = make_task(2, due_date=None)
        assert self.cmp.compare(a, b) < 0

    def test_task_without_due_date_after_task_with(self):
        a = make_task(1, due_date=None)
        b = make_task(2, due_date="2024-01-01")
        assert self.cmp.compare(a, b) > 0

    def test_same_due_date_returns_zero(self):
        a = make_task(1, due_date="2024-01-01")
        b = make_task(2, due_date="2024-01-01")
        assert self.cmp.compare(a, b) == 0

    def test_both_no_due_date_returns_zero(self):
        a = make_task(1, due_date=None)
        b = make_task(2, due_date=None)
        assert self.cmp.compare(a, b) == 0


class TestHighPriorityComparator:
    cmp = HighPriorityComparator()

    def test_high_priority_comes_first(self):
        a = make_task(1, is_high_priority=True)
        b = make_task(2, is_high_priority=False)
        assert self.cmp.compare(a, b) < 0

    def test_low_priority_comes_after(self):
        a = make_task(1, is_high_priority=False)
        b = make_task(2, is_high_priority=True)
        assert self.cmp.compare(a, b) > 0

    def test_same_priority_returns_zero(self):
        a = make_task(1, is_high_priority=True)
        b = make_task(2, is_high_priority=True)
        assert self.cmp.compare(a, b) == 0


class TestCreatedAtComparator:
    cmp = CreatedAtComparator()

    def test_older_task_comes_first(self):
        a = make_task(1, created_at=100)
        b = make_task(2, created_at=200)
        assert self.cmp.compare(a, b) < 0

    def test_newer_task_comes_after(self):
        a = make_task(1, created_at=200)
        b = make_task(2, created_at=100)
        assert self.cmp.compare(a, b) > 0

    def test_same_created_at_returns_zero(self):
        a = make_task(1, created_at=100)
        b = make_task(2, created_at=100)
        assert self.cmp.compare(a, b) == 0


# ── TestTaskManager ───────────────────────────────────────────────────────────

class TestTaskManager:
    def test_deduplication_keeps_first_seen(self, capsys):
        tasks = [
            make_task(1, description="Fix bug", due_date="2024-01-01"),
            make_task(2, description="Fix bug", due_date="2024-01-01"),
        ]
        TaskManager(tasks).print_prioritized_tasks()
        lines = printed_lines(capsys)
        assert len(lines) == 1
        assert "Task ID: 1" in lines[0]

    def test_mixed_due_date_ordering(self, capsys):
        # task with far-future due_date should come before high-priority no-date task
        tasks = [
            make_task(1, description="No date task", due_date=None, is_high_priority=True),
            make_task(2, description="Far future task", due_date="2099-01-01"),
        ]
        TaskManager(tasks).print_prioritized_tasks()
        lines = printed_lines(capsys)
        assert "Task ID: 2" in lines[0]
        assert "Task ID: 1" in lines[1]

    def test_priority_tiebreak_same_priority_oldest_wins(self, capsys):
        tasks = [
            make_task(1, description="Newer", due_date=None, is_high_priority=False, created_at=200),
            make_task(2, description="Older", due_date=None, is_high_priority=False, created_at=100),
        ]
        TaskManager(tasks).print_prioritized_tasks()
        lines = printed_lines(capsys)
        assert "Task ID: 2" in lines[0]  # oldest created_at first
        assert "Task ID: 1" in lines[1]

    def test_same_due_date_high_priority_comes_first(self, capsys):
        tasks = [
            make_task(1, description="Low prio",  due_date="2024-06-01", is_high_priority=False),
            make_task(2, description="High prio", due_date="2024-06-01", is_high_priority=True),
        ]
        TaskManager(tasks).print_prioritized_tasks()
        lines = printed_lines(capsys)
        assert "Task ID: 2" in lines[0]
        assert "Task ID: 1" in lines[1]

    def test_parent_filtered_out_omits_suffix(self, capsys):
        tasks = [
            make_task(1, description="Parent", is_done=True),
            make_task(2, description="Child", parent_id=1),
        ]
        TaskManager(tasks).print_prioritized_tasks()
        lines = printed_lines(capsys)
        assert len(lines) == 1
        assert "parent:" not in lines[0]

    def test_parent_present_shows_suffix(self, capsys):
        tasks = [
            make_task(1, description="Parent task", due_date="2024-01-01"),
            make_task(2, description="Child task",  due_date="2024-01-02", parent_id=1),
        ]
        TaskManager(tasks).print_prioritized_tasks()
        lines = printed_lines(capsys)
        child_line = next(l for l in lines if "Task ID: 2" in l)
        assert ", parent: Parent task" in child_line

    def test_done_and_assigned_tasks_excluded(self, capsys):
        tasks = [
            make_task(1, description="Done",     is_done=True),
            make_task(2, description="Assigned", assignee="alice"),
            make_task(3, description="Active"),
        ]
        TaskManager(tasks).print_prioritized_tasks()
        lines = printed_lines(capsys)
        assert len(lines) == 1
        assert "Task ID: 3" in lines[0]

    def test_output_format(self, capsys):
        tasks = [make_task(42, description="Write tests")]
        TaskManager(tasks).print_prioritized_tasks()
        lines = printed_lines(capsys)
        assert lines[0] == "Task ID: 42, Description: Write tests"

    def test_empty_input(self, capsys):
        TaskManager([]).print_prioritized_tasks()
        assert printed_lines(capsys) == []
