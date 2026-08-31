# .claude/

Claude Code configuration for this repo: subagent definitions, skills, hooks, and settings. Root [CLAUDE.md](../CLAUDE.md) has the full agent/skill catalogue by role — this file is about how the pieces in this folder actually fit together, for an agent editing them rather than just using them.

## Layout

- `CLAUDE.md` — project instructions scoped to this folder itself (currently just the epistemic-rules block for this session).
- `settings.json` — the one file that wires everything below into Claude Code's runtime (see below).
- `agents/*.md` — one file per subagent, plus `agents/hooks/kubectl-guard.sh` — a real Claude Code `PreToolUse` hook, just declared in `kubernetes-specialist.md`'s own frontmatter rather than in `settings.json` (see "Two kinds of hook" below).
- `skills/<name>/SKILL.md` — one folder per skill. A skill can carry a `references/*.md` subfolder for material loaded only when relevant (e.g. `review-pr/references/{go,rust,cpp}.md` loaded only when the diff contains that language; `write-rust/references/*.md` loaded for borrow-checker/async/Go-mapping guidance). Skill folder names are verb-first (`run-pr-loop`, `review-pr`, `scope-problem`, `write-rust`, ...) — this is a recent rename from noun-first names (`pr-loop`, `pr-review`, `scoping`, `rust-expert`); expect stale references to the old names in anything written before that rename.
- `hooks/*.sh` — the global hooks wired in `settings.json` (currently just `PostToolUse`) plus their test suites. Distinct from `agents/hooks/`, which holds hooks scoped to one specific agent.

## Agent definition conventions

Each `agents/*.md` file is YAML frontmatter + a system prompt. Frontmatter fields actually in use here, beyond the required `name`/`description`: `tools` (explicit allowlist — keep it minimal per agent, not a blanket grant), `model` (overrides the default, including non-Sonnet choices like `model: fable`), `effort` (e.g. `max` for bar-raiser agents), `maxTurns`, and `memory: project` (used by `code-bar-raiser` to persist state across PR-loop rounds). A subagent's `description` is load-bearing: it's what the dispatching agent reads to decide when to use it, so the "Not for X (use Y)" negative-space clauses at the end of each description are what keep near-duplicate agents (e.g. `code-reviewer` vs `code-bar-raiser` vs `leetcode-reviewer`) from being reached for interchangeably.

## `settings.json`

- `permissions.allow`/`deny` — broadly permissive (`Edit(*)`, `Write(*)`, `Bash(*)`) with a short explicit denylist (`rm -rf`, `sudo`, `curl`, `wget`). `skipDangerousModePermissionPrompt: true` goes with this — the guardrails here are the denylist, not interactive confirmation.
- `hooks.PostToolUse` wires `hooks/validate-definitions.sh` to fire after `Edit|Write|MultiEdit|Bash`. The hook itself does the path filtering (matching on the tool payload, not the matcher) so a single entry catches both direct edits and Bash-driven deletes/moves — see the hook's own header comment before changing the matcher.
- `effortLevel: xhigh` and `env.CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS` are project-wide defaults for this session type, unrelated to any one agent or skill.

## Two kinds of hook

Both are real Claude Code hooks; they differ in scope and in whether they can actually block anything.

- **Global, report-only:** `hooks/validate-definitions.sh`, wired in `settings.json` under `PostToolUse` for every session. `PostToolUse` fires after the edit already happened, so it can only report a broken agent/skill definition back to Claude (exit 2 + stderr) — it never blocks. It fails open (a missing interpreter reports itself and gets out of the way) rather than blocking all edits on a broken toolchain. `hooks/test-validate-definitions.sh` is its test suite (24 cases per root CLAUDE.md) — run it after changing either the hook or `agent-team-workspace/validate-definitions.py`, since the hook is a thin wrapper around that validator.
- **Per-agent, enforcing:** `agents/hooks/kubectl-guard.sh`, wired in `kubernetes-specialist.md`'s own frontmatter under `PreToolUse`, active only while that subagent runs. `PreToolUse` fires before the command executes, so it can actually `deny` — it blocks mutating `kubectl`/`helm`/`oc` commands and `Secret` value reads before Bash ever runs them, and fails closed if `jq` is missing. If you need a new hook that must genuinely stop an action rather than just flag it after the fact, this is the pattern to copy — a global `PostToolUse` hook structurally cannot do that.
