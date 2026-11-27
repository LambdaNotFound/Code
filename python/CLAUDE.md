# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Setup

Needs Python 3.10+ (annotations such as `str | None` fail on import before that) and pytest. This project uses a local virtualenv; see [README.md](README.md) to create it. Activate it before running anything:

```bash
source venv/bin/activate
```

## Commands

```bash
# Run all tests
pytest

# Run a single package
pytest dsa/stack/

# Run a single test file
pytest dsa/stack/stack_test.py

# Run a single test by name
pytest dsa/stack/stack_test.py -k test_decodeString

# Verbose output
pytest -v
```

## Architecture

Two folders, matching the repo-wide layout (see the root `CLAUDE.md`): `dsa/` holds data-structure and algorithm solutions, and `interview/` holds company mock-interview problems. Run every command from `python/`, the folder this file is in.

**`dsa/` mirrors `go/dsa/`.** It has `go/dsa/`'s 34 folders, one per leetcode.com tag (`array`, `backtracking`, … `union_find`), plus `types/` and `utils/`. `go/dsa/containers/` has no Python counterpart yet. A Python solution goes in the same folder and the same file as the Go solution to that problem, with the file's `.go` swapped for `.py`: Two Sum is `go/dsa/hash_table/hashmap.go` in Go and `dsa/hash_table/hashmap.py` here. When Go has no solution to the problem, the root `CLAUDE.md`'s "Which tag folder a solution goes in" rule picks the folder. Every folder is a package with an `__init__.py`, even while it holds no solution yet.
- `dsa/types/` — LeetCode's node classes, `ListNode` (`list.py`) and `TreeNode` (`tree.py`), the counterpart of `go/dsa/types/`. Import them as `from dsa.types import ListNode`.
- `dsa/utils/` — test helpers, the counterpart of `go/dsa/utils/`: `build_list` and `list_to_vals` (`list_utils.py`), and `build_tree` (`tree_utils.py`). Import them as `from dsa.utils import build_tree`.
- `interview/` — company mock-interview problems; see "Interview problems" below

**Import conventions:**

`dsa/` imports package-style, from `python/`:
```python
from dsa.types import ListNode
from dsa.stack.stack import Solution
```
This works when a test file runs on its own (`pytest dsa/stack/stack_test.py`) because `dsa/` and each folder in it have an `__init__.py`: pytest then puts the first folder above them with no `__init__.py`, `python/`, on `sys.path`. Never insert a `dsa/` folder into `sys.path`. With `dsa/` itself on the path, the `types/` and `string/` packages would shadow the standard library modules of the same names, and so would `array/` and `math/` on Python builds that ship those two as files (Homebrew's and Apple's do; uv's builds compile them in). `python -m pytest` run from inside `dsa/` does exactly that, because `-m` puts the current folder on the path, and it stops at startup with `No module named 'dsa'`; plain `pytest` from any folder is fine. Package imports are also what let two folders hold files of the same name, as Go's do (`dp/word_break.py` and `trie/word_break.py`).

`interview/` has no `__init__.py`, so its tests still insert their own directory and import the solution by bare name:
```python
sys.path.insert(0, os.path.dirname(__file__))
from rippling_task_manager import TaskManager
```
Do not import `interview/` package-style (`from interview.foo import X`): nothing puts `python/` on `sys.path` for it, so such a test fails when run alone. Every test file must pass when run by itself: `pytest path/to/x_test.py`.

## Testing conventions

- Table-driven tests using `@pytest.mark.parametrize`
- When a function has multiple valid outputs (e.g. set of top-k elements, any optimal decomposition), validate properties — length, sum, containment — rather than exact equality
- When two implementations of the same function exist, parametrize over both method names via a `METHODS` list and `getattr(Solution(), method)`
- `dsa/utils/` helpers (`build_list`, `list_to_vals`, `build_tree`) should be used in tests rather than duplicating linked-list/tree construction logic inline; add new helpers to `dsa/utils/` when they would apply across packages
- Static methods (no `self`) are called as `Solution.method(args)`, not `Solution().method(args)`
- Test file naming: one `<module>_test.py` beside each solution file, as Go has one `_test.go` per file (e.g. `dsa/heap/merge_k_lists_test.py` and `dsa/heap/top_k_frequent_test.py`)

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
grep -rn "def " --include="*.py" . | grep -v test | grep -v __pycache__ | grep -v venv | grep -v 'dsa/types/\|dsa/utils/'

# test functions defined
grep -rn "def test_" --include="*.py" . | grep -v __pycache__ | grep -v venv
```

Any solution file with no corresponding `_test.py` entry is a gap. Within a test file, check that every public method on `Solution` has at least one parametrize case.
