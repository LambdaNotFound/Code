# Design-Review Loop — Agent Team Prompt

Goal: produce a high-quality engineering design (or implementation
plan) by alternating two agents — `research-investigator` authors,
`design-bar-raiser` challenges — for up to **5 rounds**, until the
bar-raiser approves or escalates.

Both agents are defined in `.claude/agents/` (hardened masters in
`docs/agents-hardened/`). Both run `model: fable` at `effort: max`,
so a full loop is deliberately expensive: use it for designs that
matter. For a cheap one-shot opinion on a human-authored design, use
`architect-reviewer` instead.

## Invocation

> Run the design-review loop on: <brief — the problem, requirements,
> and any design doc or code paths to start from>. Topic slug:
> `<slug>`.

## State: two files, one writer each

| Path | Sole writer |
|---|---|
| `docs/research/<slug>/design.md` | research-investigator |
| `docs/research/<slug>/review.md` | design-bar-raiser |

The files are the entire loop state. Both agents are stateless
between rounds and re-read both files on every invocation, so fresh
invocations and resumed ones behave identically. The lead never
edits either file.

## Protocol (lead session)

1. **Round 0 — author.** Invoke `research-investigator` with the
   brief and the slug. It researches from first principles and
   writes `design.md` (with `## Revision log` and
   `## Objection responses` sections).
2. **Round N (1..5) — challenge.** Invoke `design-bar-raiser` with
   the slug and round number. It derives the requirements
   independently, spot-checks the design's citations against the
   codebase, appends `## Round N` to `review.md`, and returns a
   verdict:
   - `approve` / `approve-with-risks` → **stop**; report the verdict,
     residual risks, and both file paths to the user.
   - `revise` → step 3.
   - `reject-approach` → step 3, but the investigator restarts the
     design core rather than patching it. Costs a round.
   - `escalate` → **stop**; put the escalation paragraph to the
     user. Do not keep looping around an irreconcilable core.
3. **Respond and revise.** Invoke `research-investigator` with the
   slug and round number. It answers every objection by ID —
   accepted (with the revision), rebutted (with evidence), or
   deferred (with why) — revises `design.md` in place, and logs the
   round in the revision log. Then back to step 2 as round N+1.
4. **Round 5 is the floor for a decision, not a target.** The
   bar-raiser must land approve / approve-with-risks / escalate by
   round 5; converging in 2 rounds because round 1 was thorough is
   the success case.

## Loop rules

- Relay only slugs, round numbers, and verdicts between agents; the
  substance travels through the two files. Do not paraphrase
  objections or responses into the prompts — paraphrase drifts.
- Prefer resuming each agent by name across rounds where the
  platform supports it (context carries over); fresh invocations are
  equally correct since the files hold all state.
- Neither agent may be asked to play the other's role. The
  investigator never writes `review.md`; the bar-raiser never writes
  `design.md`. Disagreement between them is signal — surface it,
  never smooth it over.
- If a round produces no visible change (no revision log entry, no
  new review section), the loop is stuck: stop and tell the user
  which agent stalled and on what.

## Deliverable

The final `design.md` (approved or approved-with-risks), the full
`review.md` ledger as the audit trail, and — on escalate — the one
paragraph the human must decide.
