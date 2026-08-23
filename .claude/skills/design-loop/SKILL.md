---
description: Run the design-review loop — research-investigator authors a design or RFC, design-bar-raiser challenges it for up to 5 rounds, ai-writing-auditor gives the approved doc a final editorial pass. Use when the user asks to run the design-review loop, wants a bar-raised design/RFC produced by the agent loop, or wants to resume an interrupted loop by slug.
argument-hint: <brief, or slug to resume> [slug: <slug>] [Deliverable: RFC]
---

You are the lead of the design-review loop. The protocol is
`docs/design-review-loop-agent-team-prompt.md` — read it now and
follow it exactly. This skill is only the entry point; it restates
nothing, so the protocol file stays the single source of truth. The
agents are `research-investigator` and `design-bar-raiser`
(`.claude/agents/`), plus `ai-writing-auditor` for the editorial
pass.

## Arguments

$ARGUMENTS

## Before round 0

1. Determine the slug: use the one given, else derive a short
   kebab-case slug from the topic.
2. If `docs/research/<slug>/` already exists, this is a **resume**:
   run the protocol's "Resuming an interrupted loop" derivation on
   the files and continue from the state it yields. The files
   outrank these arguments on any disagreement, round numbers
   included.
3. Fresh start: write `docs/research/<slug>/brief.md` verbatim from
   the arguments — problem, requirements, deliverable, constraints —
   before any agent runs. If the arguments name a topic but carry no
   requirements to freeze, stop and ask the user for them; never
   invent the brief.
4. If the deliverable is an RFC, `docs/rfc-spec.md` is the output
   contract (the protocol covers how).

## While the loop runs

- Follow the protocol's steps and loop rules to the letter: relay
  only slugs, round numbers, and verdicts between agents; never
  edit `design.md` or `review.md`; change `brief.md` only by
  appending a dated `## Amendment` section, and only on the user's
  instruction.
- Commit `docs/research/<slug>/` after each round unless the user's
  git conventions say otherwise — the directory is the loop's only
  durable checkpoint.
- Report to the user at the end, or on escalate or stall: verdict,
  file paths, residual risks. Not a play-by-play of rounds.
