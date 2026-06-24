# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Setup

This project uses a local virtualenv. Activate it before running anything:

```bash
source venv/bin/activate
```

## Commands

```bash
# Run all tests
pytest

# Run a single package
pytest stack_queue/

# Run a single test file
pytest stack_queue/stack_test.py

# Run a single test by name
pytest stack_queue/stack_test.py -k test_decodeString

# Verbose output
pytest -v
```

## Architecture

Python LeetCode solutions organized by algorithm category. Each category is a directory containing one or more solution files and a co-located `_test.py` file.

**Packages:**
- `array/`, `backtracking/`, `binary_search/`, `dynamic_programming/`, `graph/`, `hashmap/`, `heap/`, `interview/`, `prefix_sum/`, `solid_coding/`, `stack_queue/`, `tree/`, `two_pointers/` — solutions grouped by technique
- `common/` — shared node definitions (`ListNode`, `TreeNode`) and test helpers (`build_list`, `list_to_vals`, `build_tree`); imported with `sys.path.insert` since there is no `setup.py`

**Import conventions:**

Solution files that need shared types insert the repo root into `sys.path`:
```python
import sys, os
sys.path.insert(0, os.path.dirname(os.path.dirname(__file__)))
from common import ListNode
```

Test files insert their own directory to import the solution module:
```python
sys.path.insert(0, os.path.dirname(__file__))
from stack import Solution
```

**Name collision with stdlib — `array/`:**
`array.py` conflicts with the stdlib `array` module. Its test loads it via `importlib.util.spec_from_file_location` instead of a plain import. Follow the same pattern if any future file name collides with stdlib (e.g. `queue`, `calendar`, `math`).

## Testing conventions

- Table-driven tests using `@pytest.mark.parametrize`
- When a function has multiple valid outputs (e.g. set of top-k elements, any optimal decomposition), validate properties — length, sum, containment — rather than exact equality
- When two implementations of the same function exist, parametrize over both method names via a `METHODS` list and `getattr(Solution(), method)`
- `common/` helpers (`build_list`, `list_to_vals`, `build_tree`) should be used in tests rather than duplicating linked-list/tree construction logic inline; add new helpers to `common/` when they would apply across packages
- Static methods (no `self`) are called as `Solution.method(args)`, not `Solution().method(args)`
- Test file naming: `<module>_test.py` co-located with the solution file; when a directory has multiple solution files, a single `<dir>_test.py` covers them all (e.g. `heap/heap_test.py` tests both `merge_k_lists.py` and `top_k_frequent.py`)

## Interview problems (`interview/`)

**Module-level design comment** — every interview solution opens with a block comment (not a docstring) describing the data model, key invariants, and complexity rationale. Format:

```python
# Design:
#   _field   type   — purpose and any invariant
#
# Key decision: why this data structure was chosen over alternatives.
```

**Method complexity annotations** — inline comment on each method signature:

```python
def process(self, document: str) -> int:  # T: O(n), S: O(1)
```

**Private attributes** — all instance state uses `_` prefix (`_map`, `_scores`, `_next_id`).

**Common patterns used across problems:**

- *Strategy pattern*: `Handler` ABC with `process()` (document processor, logger). Concrete handlers compose into a manager class that holds `list[Handler]`.
- *Bounded recent-N list*: for "last N unique items," maintain a list capped at N. On insert of X: remove X if present (O(N) = O(1) since N is fixed), append X, pop front if over N. O(1) insert and retrieval.
- *Fixed-size min-heap for top-k*: push `(score, item)`, pop when `len > k`. Uses `_RevStr` wrapper to invert string comparison when lexically-smallest ties must be kept (plain `(count, word)` would pop the lexically smallest on ties, keeping the largest — wrong).
- *Incremental score maintenance*: first vote ±1; direction change ±2. Avoids recomputing from history.
- *`heapq.nlargest(k, items, key=...)`*: O(n log k) top-k; prefer over sorting when k ≪ n.
- *Python method overloading*: use a default `None` parameter and branch on it (`get(k, snap_id=None)`).

**Test structure for interview problems** — use `@pytest.fixture` for shared setup and group tests into classes by feature area (e.g. `TestBasicMapOps`, `TestTakeSnapshot`). Not everything needs to be parametrized; explicit named tests are fine when the case has a distinct semantic.

## Coverage auditing

To check coverage, diff solution methods against tested methods:

```bash
# methods defined in solution files
grep -rn "def " --include="*.py" . | grep -v test | grep -v __pycache__ | grep -v venv | grep -v common

# test functions defined
grep -rn "def test_" --include="*.py" . | grep -v __pycache__ | grep -v venv
```

Any solution file with no corresponding `_test.py` entry is a gap. Within a test file, check that every public method on `Solution` has at least one parametrize case.
