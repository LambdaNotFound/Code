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

# Check coverage (summary per package)
go test ./... -cover

# Full coverage report
go test ./... -coverprofile=coverage.out && go tool cover -func=coverage.out
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

**`pkg/graph` sub-packages:**
- `BFS/` — BFS on graphs and grids
- `Bellman_Ford/` — Bellman-Ford SSSP; handles negative weights and detects negative cycles; standard adjacency-list signature `graph [][][2]int`
- `DFS/` — DFS traversal patterns
- `Dijkstra/` — Dijkstra with a generic `Heap[T]` (see `generic_heap.go`); adjacency-list signature `graph [][][2]int, src int`; also covers LC 787, 743, 1631, 1514, 499, 505
- `bidirectional_BFS/` — meet-in-the-middle BFS
- `multi_source/` — multi-source BFS (e.g. rotting oranges, 01-matrix)
- `topological_sort/` — Kahn's algorithm and DFS-based topo sort
- `union_find/` — Union-Find with path compression

**Import conventions:**
- Packages that use `types` use a dot import: `. "gocode/types"` so `ListNode`, `TreeNode`, etc. are available unqualified
- The module is named `gocode` (see `go.mod`)

**Testing:**
- All tests use `github.com/stretchr/testify/assert`
- Test functions are co-located with implementation files in the same package (no separate `_test` packages)
- `utils/` provides helpers like `CreateLinkedList`, `VerifyLinkedLists`, and `GraphsEqual` for test setup

**Naming conventions:**
- Functions suffixed with `Claude` (e.g. `asteroidCollisionClaude`) are alternative implementations of the same problem — a different algorithm or data structure approach, not a replacement
- Functions suffixed with `Naive` are brute-force or O(n²) baselines kept for comparison

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

**Block comments:** when commenting out multiple lines, use a `/** ... */` block with every line prefixed by `* ` (GoLand-style), not `//` per line. Example:
```go
/**
 * s = "abba"
 * wordDict = ["ab", "a", "abb"]
 * Trace:
 * i=1: "a" matches, freq["a"]-- → 0, dp[1]=true
 */
```
