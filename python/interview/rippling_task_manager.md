### Overview
Design a task management system, it takes a list of tasks and out put the tasks in order based on set of rules. The set of rules can be extended. 

### Data Model
Task:
task = {
    "id": int, # unique identifier
    "description": str,
    "is_done": bool, # true if done
    "due_date": str | None, # ISO format: "YYYY-MM-DD" or None/null
    "parent_id": int | None, # unique identifier of parent task, if any
    "assignee": str | None, # user identifier or None/null if unassigned
    "is_high_priority": bool, # True if task is marked high priority
    "created_at": int # Unix epoch time in seconds
}

Rule:
1. Filtering Rules:
Only include tasks where the task is not done, not assigned, and there isn't already a task with the same description and due date (sometimes duplicates are generated). Keep the first seen

2. Prioritization Rules:
Tasks must be printed in a priority order following these rules:

Order tasks with below priority:
- 1. Tasks with a due_date come before tasks without one; among tasks with a due_date, sort by earliest first.
- 2. Among tasks without a due_date or with same due_date, prefer high priority over low priority.
- 3. Tiebreak: prefer oldest (smallest created_at).

### API
print_prioritized_tasks() should print tasks in priority order in the following format:
    "Task ID: {id}, Description: {description}"
If the parent is filtered out or has no parent, omit the suffix. If there is a parent task, add the additional suffix ", parent: {parent_task_description}" to the output. 

### High Level Design
1. Follow best OOP practices. 

Filter Rule uses strategy pattern, with a filter() behavior. 

class TaskFilter(ABC):
    def filter(self, tasks: list) -> list: ...

class CompositeFilter:
    def __init__(self, filters: list[TaskFilter]):
        self._filters = filters
    
    def filter(self, tasks: list) -> list:
        for f in self._filters:
            tasks = f.filter(tasks)
        return tasks

Prioritization Rules use composite comparator pattern, with a compare() behavior. e.g.

class TaskComparator(ABC):
    @abstractmethod
    def compare(self, a: dict, b: dict) -> int:
        pass  # negative = a first, 0 = tie/delegate, positive = b first

class CompositeComparator(TaskComparator):
    def __init__(self, comparators: list[TaskComparator]):
        self._chain = comparators

    def compare(self, a, b) -> int:
        for comp in self._chain:
            result = comp.compare(a, b)
            if result != 0:
                return result
        return 0

comparator = CompositeComparator([
    DueDateComparator(),
    HighPriorityComparator(),
    CreatedAtComparator(),
])

2. use a hashmap to store the id => task mapping, this hashmap will be used in the parent task look-up during the print_prioritized_tasks() process. The hashmap is init after the filtering

### Test Plan
Deduplication: two tasks with same description + due_date → keep first
Mixed due_date ordering: task with far-future due_date vs. high-priority no-date task
Priority tiebreak: two no-date tasks, same priority → oldest created_at wins
Parent in filtered-out list → no suffix
Parent present in the filtered task list → suffix shows parent description.
Same due_date, different priority → high priority task comes first
