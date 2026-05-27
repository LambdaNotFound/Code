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
- `array/`, `dynamic_programming/`, `hashmap/`, `heap/`, `prefix_sum/`, `stack_queue/`, `interview/` — solutions grouped by technique
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

## Coverage auditing

To check coverage, diff solution methods against tested methods:

```bash
# methods defined in solution files
grep -rn "def " --include="*.py" . | grep -v test | grep -v __pycache__ | grep -v venv | grep -v common

# test functions defined
grep -rn "def test_" --include="*.py" . | grep -v __pycache__ | grep -v venv
```

Any solution file with no corresponding `_test.py` entry is a gap. Within a test file, check that every public method on `Solution` has at least one parametrize case.
