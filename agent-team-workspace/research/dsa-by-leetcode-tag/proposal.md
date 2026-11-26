# Reorganizing `go/dsa` by LeetCode tag

Status: done on 11-26-2026. [`mapping.tsv`](mapping.tsv) records where every solution went and why: its old file, new file, LeetCode tags, and the reason for its folder. For how to place a new solution, see "Which tag folder a solution goes in" in the root `CLAUDE.md`.

## Where the tags come from

The tags for all 262 problems come from [doocs/leetcode](https://github.com/doocs/leetcode) at commit `776cc2e` (11-26-2026), one front-matter block per problem. Its scripts read tags from leetcode.cn. That site adds tags of its own (Knapsack, Dijkstra, Tree DP, Parentheses, Lowest Common Ancestor, and others) to the leetcode.com set. Filtered to the 71 leetcode.com tags, the doocs tags match `spaced_repetition/problems.json` exactly on all 157 problems the two files share. `mapping.tsv` keeps the leetcode.cn-only tags in a separate column.

## The rule

Each solution goes to the folder of one leetcode.com tag: the one naming the algorithm the solution uses.

1. **Translate the old folder.** `go/dsa` was already grouped by algorithm, so the old folder usually named it. Translate that to LeetCode's tag (`prefix_tree` became `trie`, `hashmap` became `hash_table`, `graph/dijkstra` became `shortest_path`). For 262 of 282 solutions, LeetCode tags the problem with that same tag, and that tag is the folder.
2. **Read the code for the other 20.** In each case I read the code and picked the tag naming what it does; see "Verdicts" below. Where LeetCode does not tag a problem with the algorithm the code uses, the algorithm still decides. This happens four times: 39, 437, 778 and 1631, and 2933.
3. **A file with several solutions to one problem goes where its first solution goes.** For example, `number_of_islands.go` holds DFS, BFS, and union-find versions; the DFS version comes first, so the file went to `dfs`.
4. **Tree problems go to `binary_tree` or `bst`, even when solved by BFS or DFS.**
5. **`types/`, `utils/`, and `containers/` stay as they are.** They hold support code, not problems.

## Decisions taken

These were the user's four choices, made on 11-26-2026:

1. **Tree problems solved by BFS or DFS go to `binary_tree`.** This moved 102, 103, 199, 863, 958, 104, and 543.
2. **Groupings leetcode.com has no tag for are folded into the plain tag.** Knapsack went to `dp`, multi-source and bidirectional BFS to `bfs`, and Dijkstra and Bellman-Ford to `shortest_path`.
3. **Every tag gets its own folder, even with one file.** More solutions will be added to these later.
4. **Folder names are short.** `breadth-first-search` became `bfs`, `depth-first-search` became `dfs`, `binary-search-tree` became `bst`, `dynamic-programming` became `dp`, and "Heap (Priority Queue)" became `heap`. The repo's earlier `mono_stack` and `mono_queue` were kept. Every other folder is its tag in snake_case. Each package name is its folder name without underscores (`package binarytree`).

## Layout

| Folder | Problems | Files | Came from |
|---|---:|---:|---|
| `dp` | 29 | 19 | `dynamic_programming`, its `knapsack/`, two solutions from `backtracking`; plus rod cutting |
| `binary_tree` | 23 | 8 | `tree`, the tree problems in `graph/bfs` and `graph/dfs`, one from `api_design` |
| `design` | 19 | 18 | `api_design`, 1429 from `linked_list`; plus the transactional KV store |
| `backtracking` | 16 | 6 | `backtracking`, the second 465 from `api_design` |
| `linked_list` | 16 | 4 | `linked_list` |
| `binary_search` | 13 | 5 | `binary_search` |
| `stack` | 13 | 4 | `stack_queue` |
| `string` | 13 | 2 | `solid_coding/string` |
| `two_pointers` | 13 | 4 | `two_pointers`, 31 from `backtracking`, 1570 from `math` |
| `array` | 11 | 4 | `array` (including the intervals file), `solid_coding/ranges.go`, 189 |
| `heap` | 11 | 11 | `heap` |
| `sliding_window` | 11 | 3 | `two_pointers/sliding_window` |
| `bfs` | 9 | 10 | `graph/bfs`, `graph/multi_source`, `graph/bidirectional_bfs`, 721 and 815 from `graph`, the second 1293 |
| `hash_table` | 8 | 5 | `hashmap`, 36 from `math` |
| `prefix_sum` | 7 | 3 | `prefix_sum` |
| `shortest_path` | 7 | 4 | `graph/dijkstra` (with its `generic_heap.go`), `graph/bellman_ford` |
| `bst` | 6 | 3 | the BST problems in `tree`; plus the BST implementation, `bst_core.go` |
| `dfs` | 6 | 4 | `graph/dfs` (733, 827), `number_of_islands.go` (200, 463), `recursion` (339, 341) |
| `greedy` | 6 | 2 | `greedy`, 630 from `graph/topological_sort` |
| `bit_manipulation` | 5 | 1 | `bit_manipulation` |
| `topological_sort` | 5 | 4 | `graph/topological_sort` except 630 |
| `trie` | 5 | 4 | `prefix_tree` |
| `union_find` | 5 | 4 | `graph/union_find`, `graph_valid_tree.go` |
| `matrix` | 4 | 2 | `solid_coding/matrix.go`, 311 from `math` |
| `mono_stack` | 4 | 1 | `stack_queue/mono_stack` |
| `divide_and_conquer` | 3 | 4 | `divide_and_conquer`; plus QuickSelect |
| `math` | 3 | 2 | `math` (1828, 1884), 50 from `binary_search` (now `pow.go`) |
| `mono_queue` | 2 | 1 | `two_pointers/sliding_window/mono_queue` |
| `counting` | 1 | 1 | 169 (Boyer-Moore voting) from `math` |
| `graph` | 1 | 1 | `graph/graph.go` (133, plus BFS and DFS templates) |
| `memoization` | 1 | 1 | `memoization` |
| `ordered_set` | 1 | 1 | `tree_map` |
| `simulation` | 1 | 1 | 157 (`read4.go`) from `api_design` |
| `bitmask` | 0 | 1 | TSP bitmask DP (not a LeetCode problem) |
| `types`, `utils`, `containers` | | | unchanged |

`api_design/pinterest_violation.go` moved to `go/interview/`. It is a company problem, not a LeetCode one.

## Verdicts from reading code

These are the solutions whose old folder is not one of LeetCode's tags for the problem, plus the multi-solution files in `graph/`.

| Problem | From | To | Why |
|---|---|---|---|
| 31 Next Permutation | `backtracking` | `two_pointers` | find the pivot, swap, reverse the suffix |
| 36 Valid Sudoku | `math` | `hash_table` (as `valid_sudoku.go`) | a seen-set per row, column, and box |
| 39 Combination Sum (DP copy) | `dynamic_programming/knapsack` | `dp` | DP formulation; LeetCode tags only backtracking |
| 50 Pow(x, n) | `binary_search` | `math` (as `pow.go`) | binary exponentiation |
| 105 Construct Binary Tree (second copy) | `api_design` | `binary_tree` (as `build_tree_alt.go`) | recursive build using `slices.Index` |
| 169 Majority Element | `math` | `counting` | Boyer-Moore voting (a leetcode.cn tag) |
| 200 and 463 (islands) | `graph` | `dfs` | `numIslandsDFS` comes first; the BFS and union-find versions go with it |
| 311 Sparse Matrix Multiplication | `math` | `matrix` | compresses each row to its non-zeros |
| 339 Nested List Weight Sum | `recursion` | `dfs` | recursive DFS |
| 341 Flatten Nested List Iterator | `recursion` | `dfs` | DFS flatten in the constructor |
| 377 Combination Sum IV (memo copy) | `backtracking` | `dp` | memoized DFS is top-down DP |
| 437 Path Sum III | `prefix_sum` | `prefix_sum` (stays) | prefix sums on root-to-node paths; LeetCode tags only tree and DFS |
| 465 Optimal Account Balancing (second copy) | `api_design` | `backtracking` | exhaustive DFS over balances |
| 630 Course Schedule III | `graph/topological_sort` | `greedy` | greedy with a max-heap; not a topological sort |
| 721 Accounts Merge | `graph` | `bfs` | `accountsMerge` is BFS; the union-find version goes with it |
| 778 Swim in Rising Water | `graph/dijkstra` | `shortest_path` | Dijkstra; leetcode.com tags heap, BFS, DFS, and union-find, while leetcode.cn tags Dijkstra |
| 1293 Shortest Path with Obstacles (DP copy) | `dynamic_programming` | `bfs` | BFS over (row, column, eliminations) states; there is no DP in it |
| 1429 First Unique Number | `linked_list` | `design` | a class built on a linked hash map |
| 1570 Dot Product of Sparse Vectors | `math` | `two_pointers` | two-pointer merge of the sparse entries |
| 1631 Path With Minimum Effort | `graph/dijkstra` | `shortest_path` | Dijkstra; same tag gap as 778 |
| 2625 Flatten Deeply Nested Array | `stack_queue` | `stack` | a JavaScript-only problem with no leetcode.com tags |
| 2933 High-Access Employees | `two_pointers/sliding_window` | `sliding_window` (stays) | sorts, then slides a window of 3; LeetCode tags sorting |

## What the move changed besides file locations

- **Split files.** 8 source files and 9 test files were split, with the minority parts in new files:
  - `binary_tree/tree_bfs.go`, `binary_tree/tree_dfs.go`, `binary_tree/ancestors.go`, `binary_tree/build_tree_alt.go`
  - `dp/combination_sum_iv_memoization.go`, `two_pointers/next_permutation.go`, `two_pointers/sparse_vector_dot_product.go`
  - `greedy/course_schedule_iii.go`, `array/rotate_array.go`
  - a matching `_test.go` for each
- **Six name collisions,** each between two different implementations. Nothing was deleted; the second copy was renamed:
  - `findCheapestPrice` became `findCheapestPriceBellmanFord`
  - `combinationSum4` became `combinationSum4Memoization`
  - `maxProfit` (rod cutting) became `maxProfitRodCutting`
  - the second copies of `minTransfers`, `shortestPath` (1293), and `buildTree` got an `Alt` suffix
  - their tests were renamed to match
- **Helpers that several folders need were copied, not imported.** Each copy carries a comment naming where the original is:
  - the inorder, preorder, and postorder traversals at the end of `bst/bst_core.go`
  - `isBalanced` in `bst/BST_test.go`
  - `node` in `binary_tree/ancestors_test.go`
  - `treeEqual` in `binary_tree/build_tree_alt_test.go`
  - `deepCopyInts` in `array/rotate_array_test.go`
- **One helper was moved instead.** The BFS tests' level-order tree builder moved to `binary_tree/tree_bfs_test.go` as `buildTreeLevelOrder`, because in `binary_tree` the name `buildTree` belongs to the 105 solution.
- **One call was fixed.** `bst/BST.go`'s `lowestCommonAncestorIterative` handled p > q by calling the general binary-tree LCA (236) instead of itself. The result was the same but took O(n) instead of O(h). It now calls itself with p and q swapped.
- **One import path changed outside `go/dsa`:** `go/interview/unsafe_words*.go` now imports `gocode/dsa/trie`. It keeps its `prefixtree` alias, because the code has local variables named `trie`.
- **`math/rotate_test.go` is an empty stub,** containing only a package clause, from the 07-06-2026 copy out of `pkg/`. It was kept, in `math/`, and is safe to delete.

## How the move was checked

- `go build`, `go vet`, and `go test -race ./...` pass.
- All 3,592 passing test cases are the same ones that passed before the move. The comparison went name by name, with the five renamed tests mapped back to their old names.
- All 313 README links resolve. 306 land on their problem's header line, and the other 7 are deliberate whole-file links.
- No file is newly out of `gofmt` format. Seven of the new files had a stray blank line where an import or a problem was cut out; `gofmt -w` was run on exactly those seven.
