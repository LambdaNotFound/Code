---
description: Run the PR loop — coding-expert turns a requirements/design brief into code on a branch, code-bar-raiser reviews and challenges it for up to 4 rounds, ai-writing-auditor audits the PR description, ending in a high-quality PR ready for human review. Use when the user asks to run the PR loop, wants a design or requirements implemented as a reviewed PR, or wants to resume an interrupted PR loop by slug. Not for quick edits (use golang-pro/rust-pro) or producing the design itself (use design-loop).
argument-hint: <brief or design-doc path, or slug to resume> [slug: <slug>] [open PR: yes/no]
---

You are the lead of the PR loop. The protocol is
`docs/pr-loop-agent-team-prompt.md` — read it now and follow it
exactly. This skill is only the entry point; it restates nothing.
The agents are `coding-expert` and `code-bar-raiser`
(`.claude/agents/`), plus `ai-writing-auditor` for the PR
description; the output contract is `docs/pr-spec.md`.

## Arguments

$ARGUMENTS

## Before round 0

1. Determine the slug: use the one given, else derive a short
   kebab-case slug from the change.
2. If `docs/pr-loop/<slug>/` or branch `pr/<slug>` exists, this is
   a **resume**: run the protocol's resume derivation on the files
   and branch, and continue from the state it yields. Files outrank
   these arguments on any disagreement.
3. Fresh start: write `docs/pr-loop/<slug>/brief.md` verbatim from
   the arguments — requirements or the named design doc's content,
   target area, constraints — and create branch `pr/<slug>` from
   the default branch, before any agent runs. A design-loop
   `design.md` may be copied in as the brief. If the arguments name
   a goal but no requirements to freeze, stop and ask; never invent
   the brief.
4. Record whether the user authorized opening the PR (`open PR:
   yes`). Without that, the loop ends with everything ready and no
   PR opened — never open one unasked.

## While the loop runs

- Follow the protocol to the letter: relay only slugs, round
  numbers, and verdicts; never touch `pr.md`, `review.md`, or the
  branch's code; amend `brief.md` only by appended dated
  `## Amendment` sections on the user's instruction.
- Commit `docs/pr-loop/<slug>/` as a `[loop]` commit and push the
  branch after every round — the agents never commit loop state,
  and an uncommitted ledger does not survive the session.
- Keep `pr/<slug>` checked out for the whole loop; the agents share
  your working tree and may not switch it.
- Report at the end, or on escalate, block, split-proposal, or
  stall: verdict, branch, PR link or ready state, residual risks,
  deferred follow-ups, and the human-side etiquette reminders from
  `docs/pr-spec.md`. Not a play-by-play.
