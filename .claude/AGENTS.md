# .claude/

Claude Code configuration for this repo: subagent definitions, skills, hooks, and settings. Root [CLAUDE.md](../CLAUDE.md) has the full agent/skill catalogue by role — this file is about how the pieces in this folder actually fit together, for an agent editing them rather than just using them.

## Layout

- `CLAUDE.md` — project instructions scoped to this folder itself (currently just the epistemic-rules block for this session).
- `settings.json` — the one file that wires everything below into Claude Code's runtime (see below).
- `agents/*.md` — one file per subagent, plus `agents/hooks/kubectl-guard.sh` (a guard script for `kubernetes-specialist`, not a Claude Code lifecycle hook — don't confuse it with `.claude/hooks/`).
- `skills/<name>/SKILL.md` — one folder per skill. A skill can carry a `references/*.md` subfolder for material loaded only when relevant (e.g. `review-pr/references/{go,rust,cpp}.md` loaded only when the diff contains that language; `write-rust/references/*.md` loaded for borrow-checker/async/Go-mapping guidance). Skill folder names are verb-first (`run-pr-loop`, `review-pr`, `scope-problem`, `write-rust`, ...) — this is a recent rename from noun-first names (`pr-loop`, `pr-review`, `scoping`, `rust-expert`); expect stale references to the old names in anything written before that rename.
- `hooks/*.sh` — Claude Code lifecycle hooks (currently just `PostToolUse`) plus their test suites, distinct from `agents/hooks/`.

## Agent definition conventions

Each `agents/*.md` file is YAML frontmatter + a system prompt. Frontmatter fields actually in use here, beyond the required `name`/`description`: `tools` (explicit allowlist — keep it minimal per agent, not a blanket grant), `model` (overrides the default, including non-Sonnet choices like `model: fable`), `effort` (e.g. `max` for bar-raiser agents), `maxTurns`, and `memory: project` (used by `code-bar-raiser` to persist state across PR-loop rounds). A subagent's `description` is load-bearing: it's what the dispatching agent reads to decide when to use it, so the "Not for X (use Y)" negative-space clauses at the end of each description are what keep near-duplicate agents (e.g. `code-reviewer` vs `code-bar-raiser` vs `leetcode-reviewer`) from being reached for interchangeably.

## `settings.json`

- `permissions.allow`/`deny` — broadly permissive (`Edit(*)`, `Write(*)`, `Bash(*)`) with a short explicit denylist (`rm -rf`, `sudo`, `curl`, `wget`). `skipDangerousModePermissionPrompt: true` goes with this — the guardrails here are the denylist, not interactive confirmation.
- `hooks.PostToolUse` wires `hooks/validate-definitions.sh` to fire after `Edit|Write|MultiEdit|Bash`. The hook itself does the path filtering (matching on the tool payload, not the matcher) so a single entry catches both direct edits and Bash-driven deletes/moves — see the hook's own header comment before changing the matcher.
- `effortLevel: xhigh` and `env.CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS` are project-wide defaults for this session type, unrelated to any one agent or skill.

## Hooks

`hooks/validate-definitions.sh` is a lint, not a safety guard: `PostToolUse` fires after the edit already happened, so it can only report a broken agent/skill definition back to Claude (exit 2 + stderr), never block it. It fails open (a missing interpreter reports itself and gets out of the way) rather than blocking all edits on a broken toolchain. `hooks/test-validate-definitions.sh` is its test suite (24 cases per root CLAUDE.md) — run it after changing either the hook or `agent-team-workspace/validate-definitions.py`, since the hook is a thin wrapper around that validator.
