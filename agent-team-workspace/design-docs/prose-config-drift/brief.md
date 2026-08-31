# Brief — reconciling prose claims against agent and skill configuration

Status: draft
Origin: lead, 2026-08-31. Written as an ad-hoc test run of the design
loop at its current `maxTurns: 20` cap. The requirements below come
from the lead, not from a user scoping session, and the user has not
signed off on them. Treat them as a real problem stated by a
colleague rather than as a frozen contract.

## Problem

Definition files carry configuration in frontmatter — `model`,
`effort`, `maxTurns` on the 13 agents; round counts and agent rosters
in the two loop protocols; countable inventories in `CLAUDE.md`. The
same facts are also written out in prose, in other files, for a
reader who will never open the frontmatter.

Nothing reconciles the two. When configuration changes, the prose is
updated by whoever remembers.

This is not hypothetical. `.claude/skills/scope-problem/SKILL.md:18`
told its reader the design investigator "runs at max effort for up to
fifty turns, so it is the most expensive way to guess wrong" for two
commits after the agent was capped at 20. The line survived
`validate-definitions.py` (231 checks), `validate-skills.py`, and a
24-case hook suite, and was found only because a human read the
sentence while doing something else.

The cost is specific: these files are read by agents to decide how to
act. A skill that overstates an agent's turn budget by 2.5x is
advising against a dispatch on false grounds. Left alone, every
number in prose decays into a claim nobody can trust, and the cheapest
response — deleting the numbers — makes the documents worse, because
the numbers are what make the guidance actionable.

## Goals

1. A prose sentence that states an agent's `model`, `effort`, or
   `maxTurns` and disagrees with that agent's frontmatter fails a
   validator run. Check: reintroduce "fifty turns" at
   `scope-problem/SKILL.md:18` and see a non-zero exit naming the
   file, the line, the claimed value, and the actual value.
2. The check survives a rename. Check: rename an agent and re-run;
   the check still reconciles the renamed agent's claims and does not
   report the old name as missing.
3. False positives are zero on the tree as it stands today. Check:
   run against `main` at `8832723` and get no new failures. The eight
   countable claims in `CLAUDE.md` (231 checks, 24 hook cases, 13
   agents, 8 skills, 165 problems, 24 Rust examples, 12 + 12
   example files) are all currently correct and must stay passing.
4. A claim the checker cannot parse is reported as unparsed, not
   silently skipped. Check: a deliberately odd phrasing appears in
   the output as unhandled rather than passing.

## Non-goals

- Rewriting the prose to remove numbers. The numbers stay; the check
  makes them trustworthy.
- Reconciling prose against anything other than agent and skill
  definition files and the loop protocols. Go, Python, and Rust
  source are out of scope.
- Natural-language understanding. A bounded set of claim shapes is
  the target; open-ended English is not.
- Enforcing that every configuration value be mentioned in prose. The
  check is one-directional: a claim that exists must be true.
- Changing any agent's actual configuration.

## Constraints and invariants

- Checks live in `agent-team-workspace/validate-definitions.py`
  (agents, CLAUDE.md, cross-file) or
  `agent-team-workspace/validate-skills.py` (per-skill). Both exit
  non-zero on failure. `validate-definitions.py:19` shows the check
  primitive: `ck(cond, label, detail)`, which must return its
  condition — several checks gate on it, and a version that dropped
  the return silently killed 51 of them.
- `.claude/hooks/validate-definitions.sh` re-runs the definition
  validator after any `Edit`, `Write`, `MultiEdit`, or `Bash` naming
  an agent, a skill, `CLAUDE.md`, or the validator. It fails open and
  costs ~13ms on unrelated calls; whatever this design adds runs
  inside that budget or justifies exceeding it.
- Python 3, standard library only, matching both existing validators.
- A check that judges by a hardcoded list of names or values is the
  failure mode this repo has hit three times: an allowlist vouching
  for a deleted agent, a filter that could only validate already-known
  targets, a regex naming five agents. Judge by structure.
- New checks need coverage in `.claude/hooks/test-validate-definitions.sh`,
  which runs its 24 cases against a throwaway copy under `/tmp`.

## Decomposition

| piece | goal it serves | route | depends on |
|---|---|---|---|
| Claim taxonomy: which prose shapes are in scope, and what each reconciles against | 1, 4 | run-design-loop | — |
| Extraction and reconciliation check in the validators | 1, 2, 3 | run-pr-loop | taxonomy |
| Test cases in the hook suite, including the "fifty turns" regression | 1, 3 | run-pr-loop | check |

## Open questions

1. Do number-words ("twenty") count, or only digits? Owner: design.
   Default if unanswered: both, since the live defect was a
   number-word.
2. Does an unparsed claim fail the run or warn? Owner: design.
   Default: warn, so an odd sentence cannot block an unrelated commit.
3. Are round counts ("up to 5 rounds", stated in three files each)
   in scope, given the value lives in prose rather than in any
   frontmatter field? Owner: design. Default: in scope, reconciled
   against the protocol file as the single source.

## Context

- The live defect and its fix: `.claude/skills/scope-problem/SKILL.md:18`.
- A correct claim of the same shape, which must keep passing:
  `agent-team-workspace/protocols/design-review-loop-agent-team-prompt.md:15-16`
  ("Both run `model: opus`, capped at `maxTurns: 20`").
- Round-count claims restated across files:
  `design-review-loop-agent-team-prompt.md:5` and `:109`,
  `.claude/agents/design-bar-raiser.md:161`,
  `.claude/agents/design-investigator.md:218`,
  `agent-team-workspace/protocols/pr-loop-agent-team-prompt.md:6`.
- Countable claims in `CLAUDE.md`: lines 73, 82, 95.
- Existing check styles to match: `validate-definitions.py:19` (`ck`),
  `:25` (`_load`), `:104` (`ROOTS`), `:135` (`BUILTIN_AGENTS`),
  `:241` (`CM`, the CLAUDE.md text used by cross-file checks).
- The 13 agents' frontmatter is the source of truth for goal 1:
  `.claude/agents/*.md`.
