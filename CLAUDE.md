# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

Go commands run from `go/`, where the module's `go.mod` lives. Python and Rust have their own; see their sections below.

CI: [`.github/workflows/build.yml`](.github/workflows/build.yml) builds and tests all three projects on every push to `main`, as separate checks:
- **Go**, on the version `go.mod` declares: `go mod verify`, `build`, `vet`, and `test -race`
- **Rust**, on stable: `cargo build` and `cargo test`, both with `--all-targets --locked`
- **Python**, on 3.10 and 3.12: `compileall`, the full `pytest` run, and each test file run alone

Pushes by the two bot workflows don't trigger it, because they push with `GITHUB_TOKEN`. `actionlint` (`go install github.com/rhysd/actionlint/cmd/actionlint@latest`) lints workflow edits.

```bash
cd go

# Run all tests
go test ./...

# Run a single package
go test ./dsa/backtracking/

# Run a single test
go test ./dsa/backtracking/ -run Test_solveSudoku

# Run with verbose output
go test ./... -v

# Check coverage (summary per package)
go test ./... -cover

# Full coverage report
go test ./... -coverprofile=coverage.out && go tool cover -func=coverage.out
```

## Architecture

Source is grouped by language first (`go/`, `python/`, `rust/`), then by kind of work:

| Folder | Holds | Go | Python | Rust |
|---|---|---|---|---|
| `dsa/` | data-structure and algorithm problems, mostly LeetCode: one folder per LeetCode tag, with `python/dsa/` and `rust/dsa/` mirroring `go/dsa/` folder for folder and file for file | yes | yes | yes |
| `interview/` | company mock-interview problems, files prefixed by company | yes | yes | not yet |
| `design/` | low-level designs and design patterns | yes | not yet | not yet |
| `features/` | the language itself: concurrency, async runtimes, fundamentals | yes | not yet | yes |

A language gets one of these folders when it first has code of that kind: the first Python language demo starts a `features/` folder inside `python/`. `dsa/` is the exception: Python and Rust carry all 34 of `go/dsa/`'s tag folders and its `types/`, even folders that hold no solution yet, so a solution in any language goes where its Go counterpart is. The root holds only what belongs to no single language.

**`go/`** — one Go 1.23 module named `gocode`. Its `go.mod` is in `go/`, so import paths start below it (`gocode/dsa/types`).
- `go/dsa/` — one folder per leetcode.com topic tag, 34 in all: `array`, `backtracking`, `bfs`, `binary_search`, `binary_tree`, `bit_manipulation`, `bitmask`, `bst`, `counting`, `design`, `dfs`, `divide_and_conquer`, `dp`, `graph`, `greedy`, `hash_table`, `heap`, `linked_list`, `math`, `matrix`, `memoization`, `mono_queue`, `mono_stack`, `ordered_set`, `prefix_sum`, `shortest_path`, `simulation`, `sliding_window`, `stack`, `string`, `topological_sort`, `trie`, `two_pointers`, `union_find`. Folder names are the tag in snake_case, shortened to `bfs`, `dfs`, `bst`, `dp`, `heap`, `mono_stack`, and `mono_queue`; each package name is its folder name without underscores (`package binarytree`). See "Which tag folder a solution goes in" below. Three support packages sit beside the tag folders:
  - `go/dsa/types/` — shared LeetCode node definitions: `ListNode`, `TreeNode`, `Interval`, and `Node`, one struct shared by graph problems (`Neighbors`), doubly linked lists (`Next`, `Prev`), and random-pointer lists (`Random`)
  - `go/dsa/utils/` — test helpers that build and compare linked lists, trees, and graphs and copy matrices; imported only from `_test.go` files
  - `go/dsa/containers/` (package `containers`) — standalone reference implementations of data structures: heap, LRU cache, min-max stack, red-black treemap, queue, stack, hit counter. Only `go/dsa/mono_stack` imports it (as a dot import); solution packages otherwise write their own copy, so `LRUCache`, `Heap`, `Trie`, and `UnionFind` each have several definitions across packages, in keeping with the self-contained-solution rule below
- `go/interview/` — company mock-interview problems (`affirm_*`, `codesignal_*`, `otterai_*`, `starlink_*`, and others). Its one data file, `go/interview/testdata/file.txt`, is read by the demo in `affirm_spreadsheet.go` as `testdata/file.txt`, relative to the package directory `go test` runs in; Go tooling skips any directory named `testdata/`
- `go/design/lld/` (package `lld`) — low-level designs: booking, elevator, Excel sum formula (LC 631), parking lot, rate limiter, meeting-room scheduling
- `go/design/patterns/` (package `patterns`) — GoF design patterns plus three company OO-design problems; see "Design patterns" below
- `go/features/concurrency/` (package `concurrency`) — Go concurrency patterns: channels, fan-in/out, lock-free stack/queue/counter, mutex patterns, select patterns, a leaky-bucket rate limiter, and a waiter. Most patterns live as `Example*` functions in plain `.go` files (not `_test.go`) for narrative documentation — `go test` does **not** execute these (Example funcs only run when declared in a `_test.go` file); actual assertions live in the sibling `_test.go` files (e.g. `channels_test.go`), which don't yet cover every pattern
- `go/features/basics/` (package `basics`) — Go fundamentals as runnable tests: named returns and slice and map behavior, `slices.Clone`/`maps.Clone`, bytes versus runes, `iota` enums, panic/recover in place of try/catch, `time` layouts, integer-cents money, and the generic `RemoveIf` slice filter

**`python/`** and **`rust/`** — see the Python and Rust sections below.

**Root folders that belong to no language:**
- `spaced_repetition/` — the daily practice scheduler; see its section below
- `system_design/` — written system design solutions/notes (e.g. URL shortener, autopay, hotel booking, sharding, feed generation), not part of the Go module
- `agent-team-workspace/` — the PR loop's documents and the definition validators: `pull-requests/<topic>/` (run-pr-loop state), `agent-specs/` (`pr-spec`, the PR loop's output contract), `protocols/` (the PR loop protocol), `design-docs/prose-config-drift/` (the one run of this repo's former design-review loop, kept as a record: its frozen `brief.md`, the `design.md` deliverable, and the `review.md` round ledger), `validate-definitions.py` (structural and semantic checks over every agent and skill definition; it prints the count it ran) and `validate-skills.py` (per-skill regression pass: frontmatter, description cap, reference links, routing targets, repo paths, and stale names; both exit non-zero on failure), and `research/` (standalone investigations and the older ad-hoc agent-team prompts under `legacy-prompts/`). Each document artifact gets its own folder
- `notes/project_rules/` — interview-prep "project rules" prompts (behavioral, coding/DSA, system design mock-interview instructions) used to configure Claude Code sessions for interview practice; unrelated to the Go source
- `scripts/` — `sync-claude.sh` copies this repo's `.claude/` folder into `~/.claude/` so its agents, skills, and settings apply in every project on the machine (idempotent; `--dry-run`, `--exclude`, `--no-prune`, `--strict`); `test-sync-claude.sh` is its test suite
- `diagrams/` — two diagrams, `design-review-loop/` and `url-shortener/`, each holding a `.spec.json` source plus the compiled `.excalidraw` and its `.png`/`.svg` renders. They were built with this repo's former `draw-diagram`, a JSON-spec compiler retired on 11-26-2026 (git history has it), so the `.spec.json` files no longer build here. New diagrams come from the system-design repo's `draw-diagram`
- Editor config at the root: `.cursorrules` (Cursor's copy of the LeetCode rules below) and `.vscode/settings.json` (format-on-save; Go uses 4-space indentation)

**Which tag folder a solution goes in.** This is how `go/dsa` was organized on 11-26-2026; [agent-team-workspace/research/dsa-by-leetcode-tag/](agent-team-workspace/research/dsa-by-leetcode-tag/proposal.md) has the analysis, and its `mapping.tsv` records where every solution went and why. Apply the same rule to a new solution:
- It goes in the folder of the one leetcode.com tag that names the algorithm the code uses, normally one of the tags LeetCode lists on the problem's page. When LeetCode does not list the tag for the algorithm the code uses, the algorithm still decides: 437 Path Sum III sits in `prefix_sum`, and 778 and 1631 sit in `shortest_path`.
- Tree problems go in `binary_tree` or `bst`, even when solved by BFS or DFS. `bfs` and `dfs` hold graph and grid problems.
- Groupings without a leetcode.com tag go in the plain tag folder: knapsack in `dp`, multi-source and bidirectional BFS in `bfs`, Dijkstra and Bellman-Ford in `shortest_path` (Dijkstra's generic `Heap[T]` is `shortest_path/generic_heap.go`).
- A file holding several solutions to one problem goes where its first solution goes; for example, `dfs/number_of_islands.go` also holds the BFS and union-find versions.
- Every tag gets its own folder, even one with a single file. Adding a solution never moves an existing one.
- Solutions that are not LeetCode problems go with their algorithm: QuickSelect in `divide_and_conquer`, TSP in `bitmask`, rod cutting in `dp`, the BST implementation (`bst/bst_core.go`) in `bst`.
- The rule covers all three languages. A Python or Rust solution goes in the folder and the file its Go counterpart uses, with `.py` or `.rs` for `.go`, even when its own approach differs (Python's 268 sorts, but sits in `bit_manipulation/` beside Go's). When Go has the problem in two files, one per approach (23 is in both `heap/merge_k_lists.go` and `linked_list/linked_list.go`), it goes beside the approach it uses. When Go has no solution, the rules above pick the folder.
- File names are lowercase snake_case and never start with a digit, so one name works as a Go file, a Python module and a Rust module: `grid_paths_2d.go`, not `2D_grid_paths.go`. Eleven Go files were renamed to this on 11-26-2026. This volume is case-insensitive, so a rename that only changes case takes two `git mv` steps through a temporary name.

**Import conventions:**
- Packages that use `types` use a dot import: `. "gocode/dsa/types"` so `ListNode`, `TreeNode`, etc. are available unqualified
- The module is named `gocode` (see `go/go.mod`)

**Testing:**
- All tests use `github.com/stretchr/testify/assert`
- Test functions are co-located with implementation files in the same package (no separate `_test` packages)
- `go/dsa/utils/` provides helpers like `CreateLinkedList`, `VerifyLinkedLists`, and `GraphsEqual` for test setup
- `go/dsa/types` has no tests (definitions only), and neither does `go/dsa/shortest_path/bellman_ford.go`
- A helper that several folders need is copied into each one, with a comment naming where the original lives (for example, the traversals at the end of `bst/bst_core.go`). Packages never import one another's solutions.
- Exception to the split-file norm: `go/features/basics/`'s demos (all but `remove_if.go`) and `go/design/patterns/`'s GoF-pattern files (`abstract_factory_test.go`, `adapter_test.go`, `command_test.go`, `composite_test.go`, `decorator_test.go`, `facade_test.go`, `factory_method_test.go`, `iterator_test.go`, `mediator_test.go`, `observer_test.go`) combine implementation + `Test_` functions in one `_test.go` file with no separate non-test source — each file is a self-contained, runnable demonstration of a single concept

**Naming conventions:**
- Functions suffixed with `Claude` (e.g. `asteroidCollisionClaude`) are alternative implementations of the same problem — a different algorithm or data structure approach, not a replacement
- Functions suffixed with `Naive` are brute-force or O(n²) baselines kept for comparison
- When two independent solutions to one problem share a folder, the second is suffixed with its approach (`findCheapestPriceBellmanFord`, `combinationSum4Memoization`), or with `Alt` when both use the same approach (`minTransfersAlt`, `shortestPathAlt`, `buildTreeAlt`)

**Design patterns (`go/design/patterns/`):** see [go/design/patterns/README.md](go/design/patterns/README.md) for the full catalogue of which of the 23 classic GoF patterns (Creational/Structural/Behavioral) are implemented here vs. not yet. `affirm_card_game.go`, `affirm_filter_match.go`, `brex_stervo_game.go` (and their test pairs) in the same package are unrelated, company-specific mock-interview OO-design problems, not named GoF patterns.

## Commit messages

Every commit subject is `MM-DD-YYYY: <short description>`, e.g. `11-24-2026: Fix off-by-one in binary search`, at most 72 characters including the date. Use the local date (`date +%m-%d-%Y`). The body and trailers such as `Co-Authored-By` go below a blank line and are not checked.

- Enforced by `.githooks/commit-msg`, which git runs only after `git config core.hooksPath .githooks` (local config, so once per clone). It rejects an impossible date (`02-30-2026`) and a subject past 72 characters, and exempts git's own `Merge branch …`/`Revert "…"` subjects and `fixup!`/`squash!`/`amend!` markers. `git commit --no-verify` skips it for one commit.
- Its tests are `bash .githooks/test-commit-msg.sh`: unit cases on message files, plus end-to-end commits in a throwaway repo. Run them after changing the hook.
- The two GitHub Actions workflows commit in the same format, dated with the Pacific date the run worked under: `MM-DD-YYYY: Record daily LeetCode plan` or, while paused, `MM-DD-YYYY: Keep daily LeetCode plan paused`, plus `MM-DD-YYYY: Log missed LeetCode solves` when its missed-log check logged something (`daily-leetcode.yml`), and `MM-DD-YYYY: Log LeetCode solves from issue #N` (`log-solve.yml`). The hook is local and never runs on GitHub's runners, so nothing checks those subjects; keep them in the format by hand when editing either workflow. All earlier commits were reworded to this format on 11-24-2026, so a hash recorded before then (for example in `agent-team-workspace/design-docs/prose-config-drift/`) no longer resolves.

## Root README (`README.md`)

A manually-maintained index table (Problem Tag | Problem Name | LeetCode Link | Go Source | Python Source) of every solved problem, linked to the exact file and line. It is not auto-generated by any script or workflow — when adding a new solution, add a corresponding row (grouped under its existing tag section, or a new one if the tag doesn't exist yet). The Daily Practice page reads this file from `main` for its "My solution" links: in each row, the cell linking to `/problems/<slug>` must come before the source cells, and source links stay relative to the repo root (`parseSolutions` in the Pages repo's `practice/src/lib/practice.ts`).

- The `#L` anchors move whenever a solution file is moved, split, or gains or loses lines above a solution, so fix them in the same change. To find an anchor's new line, look up the old line's text in `git show HEAD:<old path>` and find that text in the new file.
- There is no Rust Source column yet; `rust/dsa/hash_table/hashmap.rs` (Two Sum) is not linked. Adding a column changes the cells `parseSolutions` reads, so check the page against it first.

## Claude Code configuration

- `.claude/CLAUDE.md` — the epistemic rules (claim tags, confidence bands, anti-sycophancy, the `[RULES I BROKE]` footer) loaded alongside this file
- `system_design/`, `spaced_repetition/`, and `notes/project_rules/` each carry an `AGENTS.md` describing that directory's files and editing conventions; read it before changing anything there
- `.claude/settings.json` wires one `PostToolUse` hook: after any `Edit`, `Write`, `MultiEdit`, or `Bash`, `.claude/hooks/validate-definitions.sh` re-runs the definition validator if the edited path — or, for `Bash`, the command text, which is how deletes and moves arrive — names an agent, a skill, the PR loop's protocol or contract under `agent-team-workspace/`, `CLAUDE.md`, or the validator itself, and shows the failures to Claude on exit 2. Without PyYAML it skips the check and says so. An unrelated `Bash` call costs ~13ms. `PostToolUse` cannot block, so this catches a broken definition immediately after it is written rather than whenever someone remembers to run the script. Its own tests are `.claude/hooks/test-validate-definitions.sh` — 27 cases over add, remove, update, and a validator that cannot run, run against a throwaway copy of the definitions under `/tmp` so the real tree is never mutated. Run them after changing either the hook or the validator.

- Scoping, design, design review, building a design as a crate, diagrams, concept notes, recall drills, and slide decks belong to the system-design repo (local clone `/Volumes/Untitled/Workspace/system-design`), not to this one. Its `scripts/sync-skills` installs those skills user-level, so `/scope-problem` and `/draw-diagram` load in GoCode sessions too, but they write system-design's layout (a `design-doc/<slug>/` folder, a Designs table in the root README, a `src/` Cargo workspace), so run them there. This repo's own scoping skill, design-review loop, diagram compiler, their three agents, and the brief and RFC contracts were retired on 11-26-2026; git history has them. The PR loop stays, and takes that repo's signed-off `requirement.md` or exported `rfc.md` as its brief
- `.claude/agents/leetcode-reviewer.md` — subagent that reviews Go LeetCode solutions for correctness, edge cases, complexity, and idiom; invoke proactively after writing/modifying a solution
- The rest of `.claude/agents/`, by role: `coding-expert` and `coding-bar-raiser` (run-pr-loop pair), `ai-writing-auditor` (final editorial pass on the PR loop's description, and on any prose on request), `architect-reviewer` (one-shot review of a design someone else wrote), `code-reviewer` (one-shot diff review outside a loop), `golang-pro` and `rust-pro` (write and modify Go and Rust; `rust-pro` is the only Rust agent), `kubernetes-specialist` (read-only against live clusters, guarded by `.claude/agents/hooks/kubectl-guard.sh`), and `multi-agent-coordinator` (plan how several agents sequence work; produces a plan, does not run it)
- `.claude/skills/` — `review-code` (Go/Python convention pass over working-tree code, with a naming-only mode that absorbed the former review-naming skill), `summarize-changes`, `run-pr-loop` (entry point for the PR loop: `coding-expert` implements, `coding-bar-raiser` reviews, protocol in [agent-team-workspace/protocols/pr-loop-agent-team-prompt.md](agent-team-workspace/protocols/pr-loop-agent-team-prompt.md), contract in [agent-team-workspace/agent-specs/pr-spec.md](agent-team-workspace/agent-specs/pr-spec.md)), `review-pr` (review a PR's code changes, or work through the review comments on one; carries Go, Rust, and C++ depth in `references/` files loaded only when the diff contains that language), `write-rust` (write Rust interactively — ownership and type design first, then compiler-driven implementation; `references/` carry a borrow-checker error-to-design-fix triage, a Go-to-Rust idiom mapping, and async shape selection), `check-usage` (report the organization's Claude API token usage over a lookback window in hours or days, by model or by key, straight from the Usage and Cost Admin API; needs an Admin API key; `scripts/test_usage.py` is its unittest suite), and `build-agent` (process for creating new subagents, skills, and agent-team loops). This repo's stock and option analysis skills were removed on 11-26-2026 (git history has them); that work lives in the fi repo (local clone `/Volumes/Untitled/Workspace/fi`), whose own analyse-stock and analyse-option skills use the Alpha Vantage MCP server

## Spaced-repetition practice (`spaced_repetition/`)

Daily LeetCode practice scheduler (Python, stdlib only; not part of the Go module). `problems.json` holds the 165-problem set (Grind 75/169 + Blind 75, bit-manipulation excluded); `state.json` tracks SM-2 review state. From 3 PM PT, a GitHub Actions workflow (`daily-leetcode.yml`) opens a new dated "Daily Practice YYYY-MM-DD" issue (named "Daily LeetCode …" before 11-25-2026) with the day's plan (earlier issues are left open), and `log-solve.yml` logs solves from the user's comments on that issue (commenting still works on closed issues) or from a "Practice log YYYY-MM-DD" issue the user opens, which the Daily Practice page (https://lambdanotfound.github.io/practice/, its own repository) prefills; after a week without a log plans pause until the user logs a solve or comments `resume`; both commit `state.json` under the repo owner's identity, so `git pull` before touching state and `git push` after logging locally. When the user reports solving a problem (e.g. "solved 200 good"), run `python3 spaced_repetition/sr.py log <id> <again|hard|good|easy>` (synonyms in `GRADE_ALIASES` also accepted; the grade is attempt performance, not LeetCode difficulty). See [spaced_repetition/README.md](spaced_repetition/README.md).

- Only `status`, `stats` and `due` are read-only; every other command writes `state.json`. `today` is the one that surprises: it records the plan it picks under `served`, and piping its output away does not stop that. Try commands on a copy of `spaced_repetition/` in the scratchpad, or restore with `git checkout -- spaced_repetition/state.json` before anything is committed.
- The Daily Practice page lives in its own repository, `LambdaNotFound/lambdanotfound.github.io` (local clone `/Volumes/Untitled/Workspace/lambdanotfound.github.io`, page in `practice/`). It has no backend: every visit reads `state.json`, `problems.json` and the root `README.md` from `main` on raw.githubusercontent.com. Renaming a field in either JSON file, or changing the shape of a plan in `served`, breaks the page, so change its types in `practice/src/lib/practice.ts` in the same piece of work. Its log button opens a "Practice log YYYY-MM-DD" issue whose body is plain `<id> <grade>` lines, so keep `parse_comment` accepting that. Log entries gained a `prev` field on 11-25-2026 (the card before the log, for `undo`); the page's `LogEntry` type lists it. With a token the viewer saves, the page's "Log now" creates that issue through the GitHub API as the token's owner, which is what `log-solve.yml`'s owner check needs. The repository answers to both `LambdaNotFound/GoCode` (the page's name for it) and `LambdaNotFound/Code` (this clone's remote).
- Comments log in a strict format: one problem per line or per comma, `20 good` or `#54 hard (needed a hint)`, plus `undo 54` and `resume`. `parse_comment` reports anything else that names a problem instead of guessing, because a wrong grade silently moves a review date. A second log of a problem on the same day replaces the first.
- Both workflows write `state.json` the same way, and editing either one must keep all of it: fetch the newest `main`, apply the change, commit, push, and on a refused push start again; a retry that finds its own earlier commit on `main` stops there, because that push landed. `log-solve.yml` has no concurrency group on purpose: a group holds one waiting run and cancels it when the next arrives, which silently lost logs. `state["events"]` records each applied comment, so retries and re-runs never log twice, and the daily run's missed-log check (`sr.py sweep`) logs comments from the last three days that have no record.
- Neither workflow has automated tests, and `test_sr.py` can't reach the bugs between them. Check a workflow change by running each step's `run:` script with `bash -e` against a shallow clone of a throwaway bare repo used as `origin`, with `gh` stubbed, including a refused push and a push that lands but reports failure.

## Python sub-project

`python/` holds the Python solutions, with its own virtualenv and `CLAUDE.md`. It is completely independent of the Go module — do not mix imports or test commands. Run `pytest` from `python/`.

- See [python/CLAUDE.md](python/CLAUDE.md) for setup, commands, and conventions specific to that sub-project.
- `python/dsa/` mirrors `go/dsa/`: the same 34 tag folders plus `types/` and `utils/`, each a package, and each Python solution in the file its Go counterpart uses, with `.py` for `.go` (`python/dsa/hash_table/hashmap.py` beside `go/dsa/hash_table/hashmap.go`). LeetCode's node classes live in `python/dsa/types/` and the test helpers in `python/dsa/utils/`. Imports are package-style (`from dsa.types import TreeNode`); [python/CLAUDE.md](python/CLAUDE.md) says why, and why nothing may put a `dsa/` folder on `sys.path`.
- `python/interview/` holds company mock-interview problems, several with a `.md` problem statement beside the solution.

## Rust sub-project

`rust/` is a Cargo package (`tokio_hello_world`, edition 2024) for practicing Rust: LeetCode solutions and async Rust with Tokio. It is independent of the Go module, with its own `Cargo.toml`/`Cargo.lock` (committed, per Cargo's guidance for executables). `cargo run` prints a greeting via a task spawned with `tokio::spawn`; `cargo test` runs the async tests in `src/tests.rs`, which use `#[tokio::test(start_paused = true)]` (the `test-util` dev-dependency feature) to fast-forward simulated time instead of actually sleeping, and the unit tests under `dsa/`.

**`rust/dsa/`** is the package's library target (`[lib] name = "dsa"`, `path = "dsa/lib.rs"`) and mirrors `go/dsa/`. It has the same 34 tag folders plus `types/`, each a module whose `mod.rs` declares its files, and `lib.rs` declares all 35 modules.
- A solution goes in the folder and file its Go counterpart uses, with `.rs` for `.go`: Two Sum is `rust/dsa/hash_table/hashmap.rs`, beside `go/dsa/hash_table/hashmap.go`. Adding a file means adding a `pub mod <file>;` line to that folder's `mod.rs`, since Cargo compiles only modules declared from `lib.rs` down.
- Each file defines its own `pub struct Solution;` with the solution as an associated function in `impl Solution`, the way LeetCode's Rust template does. Its tests sit in the same file, in a `#[cfg(test)] mod tests`.
- `dsa/types/` holds LeetCode's own `ListNode` and `TreeNode` definitions (`ListNode::new` is made `pub`), re-exported as `crate::types::ListNode`. `go/dsa/utils/` and `go/dsa/containers/` get Rust counterparts when a solution first needs them.
- `lib.rs` allows `clippy::module_inception`, because files named after their folder, like `heap/heap.rs`, are what the mirror produces.
- CI runs neither rustfmt nor clippy, so check those by hand: `cargo clippy --all-targets --locked -- -D warnings`, and `rustfmt --check --edition 2024 dsa/lib.rs`, which follows the `mod` lines into every module. Don't use plain `cargo fmt --check`: it also reports 15 files under `features/` whose formatting predates `dsa/`.

Beyond `src/`, the crate holds 24 runnable teaching examples under `features/`: `features/tokio/` (12 async patterns — cancellation, channels, cooperative scheduling, error handling, leaky bucket, runtime flavors, `select!`, shared state, sleep/timeout, spawn/join, `spawn_blocking`, structured concurrency) and `features/basics/` (12 language fundamentals). Two conventions hold across all of them: every file opens with a `//! Pattern:` line followed by a `//! Go comparison:` paragraph mapping the Rust shape onto a Go engineer's instinct, and — because neither directory is named `examples/`, so Cargo will not auto-discover them — **each example needs an explicit `[[example]]` block in `Cargo.toml`** naming its path.

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
