# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Run all tests
go test ./...

# Run a single package
go test ./pkg/backtracking/

# Run a single test
go test ./pkg/backtracking/ -run Test_solveSudoku

# Run with verbose output
go test ./... -v
```

## Architecture

This is a Go (1.23) repository of LeetCode and algorithmic problem solutions, organized by technique.

**Top-level packages:**
- `pkg/` — LeetCode solutions grouped by algorithm category: `backtracking`, `binary_search`, `divide_and_conquer`, `dynamic_programming`, `graph`, `greedy`, `tree`, `tree_map`, `two_pointers`, `heap`, `hashmap`, `prefix_sum`, `prefix_tree`, `linked_list`, `stack_queue`, `recursion`, `memoization`, `bit_manipulation`, `math`, `design`, `api_design`, `oo_design`, `solid_coding`, `interview`
- `containers/` — custom data structure implementations used across problems: heap, LRU cache, min-max stack, red-black treemap, queue, stack, hit counter
- `concurrency/` — Go concurrency patterns: channels, fan-in/out, lock-free stack/queue/counter, mutex patterns, select patterns
- `types/` — shared LeetCode node definitions (`ListNode`, `TreeNode`, `Node` for graphs)
- `utils/` — test helpers for constructing and comparing linked lists, trees, and graphs
- `fixtures/` — static test data files (e.g. input text for API design problems)
- `notes/` — system design and architecture reference notes (e.g. sharding, replication patterns)

**Import conventions:**
- Packages that use `types` use a dot import: `. "gocode/types"` so `ListNode`, `TreeNode`, etc. are available unqualified
- The module is named `gocode` (see `go.mod`)

**Testing:**
- All tests use `github.com/stretchr/testify/assert`
- Test functions are co-located with implementation files in the same package (no separate `_test` packages)
- `utils/` provides helpers like `CreateLinkedList`, `VerifyLinkedLists`, and `GraphsEqual` for test setup

## Python sub-project

`python/` is a parallel Python LeetCode repository with its own virtualenv and `CLAUDE.md`. It is completely independent of the Go module — do not mix imports or test commands.

- See [python/CLAUDE.md](python/CLAUDE.md) for setup, commands, and conventions specific to that sub-project.
- Python packages mirror the Go structure: `array`, `backtracking`, `binary_search`, `dynamic_programming`, `graph`, `hashmap`, `heap`, `interview`, `prefix_sum`, `solid_coding`, `stack_queue`, `tree`, `two_pointers`; shared types live in `common/`.

## LeetCode conventions

- Strictly follow the provided function signature — do not change return types or parameters
- Assume input constraints are met; do not add extra validation unless the problem specifies it
- Prefer runtime optimization over memory when there is no explicit space constraint
- Add helper functions only when genuinely necessary; keep solutions self-contained
- Think through edge cases and T/S complexity before writing code
