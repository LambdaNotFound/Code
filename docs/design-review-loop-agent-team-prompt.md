# Design-Review Loop — Agent Team Prompt

Goal: produce a high-quality engineering design (or implementation
plan) with three agents — `research-investigator` authors,
`design-bar-raiser` challenges for up to **5 rounds** until approval
or escalation, and `ai-writing-auditor` gives the approved document
a final editorial pass.

All three are defined in `.claude/agents/` (hardened masters for the
first two in `docs/agents-hardened/`). The investigator and
bar-raiser are expert software engineers at both design altitudes —
high-level (architecture, boundaries, data models, consistency,
scaling, failure domains) and low-level (interfaces, data
structures, algorithms, concurrency, error semantics) — and both run
`model: fable` at `effort: max`, so a full loop is deliberately
expensive: use it for designs that matter. For a cheap one-shot
opinion on a human-authored design, use `architect-reviewer`
instead.

## Invocation

> Run the design-review loop on: <brief — the problem, requirements,
> and any design doc or code paths to start from>. Topic slug:
> `<slug>`.

## State: single writer per file

| Path | Sole writer | Role |
|---|---|---|
| `docs/research/<slug>/design.md` | research-investigator | the design |
| `docs/research/<slug>/review.md` | design-bar-raiser | append-only round ledger |
| `docs/research/<slug>/design.rewritten.md` | ai-writing-auditor | editorial intermediate |

The first two files are the entire loop state. The investigator and
bar-raiser are stateless between rounds and re-read both files on
every invocation, so fresh invocations and resumed ones behave
identically. The lead never edits any of the three.

## Protocol (lead session)

1. **Round 0 — author.** Invoke `research-investigator` with the
   brief and the slug. It researches from first principles, designs
   at both altitudes, and writes `design.md` (with `## Revision log`
   and `## Objection responses` sections).
2. **Round N (1..5) — challenge.** Invoke `design-bar-raiser` with
   the slug and round number. It derives the requirements
   independently, challenges the high-level design and verifies the
   low-level design, spot-checks citations against the codebase,
   appends `## Round N` to `review.md`, and returns a verdict:
   - `approve` / `approve-with-risks` → step 4, the editorial pass.
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
4. **Editorial pass — after approval only.** Invoke
   `ai-writing-auditor` on `docs/research/<slug>/design.md`. It
   writes the cleaned prose to `design.rewritten.md` (it never
   overwrites its source) and returns a claim-inventory report.
   - Claim inventory clean → invoke `research-investigator` for the
     adoption pass: it diffs the rewrite against the design for
     technical meaning, replaces `design.md`'s prose with the
     audited prose, and logs one editorial entry. If it returns
     `Adopted`, the loop is done. If it returns `Adopted with
     corrections`, put the named corrections to `design-bar-raiser`
     for a check of only the drifted sections before closing.
   - Report says failed rewrite (a lost claim, a changed section) →
     keep `design.md` exactly as approved, discard the rewrite, and
     tell the user the audit failed and why.
   Either way, delete `design.rewritten.md` once the loop closes;
   it is an intermediate, not a deliverable.
5. **Round 5 is the floor for a decision, not a target.** The
   bar-raiser must land approve / approve-with-risks / escalate by
   round 5; converging in 2 rounds because round 1 was thorough is
   the success case.

## Loop rules

- Relay only slugs, round numbers, and verdicts between agents; the
  substance travels through the files. Do not paraphrase objections
  or responses into the prompts — paraphrase drifts.
- Prefer resuming each agent by name across rounds where the
  platform supports it (context carries over); fresh invocations are
  equally correct since the files hold all state.
- No agent plays another's role. The investigator never writes
  `review.md`; the bar-raiser never writes `design.md`; the auditor
  touches prose only — technical verdicts and objections never route
  through it, and an editorial pass never reopens the review.
- The editorial pass runs once, at the end. Do not run the auditor
  between rounds: mid-loop prose churn invalidates the bar-raiser's
  citation spot-checks for no quality gain.
- If a round produces no visible change (no revision log entry, no
  new review section), the loop is stuck: stop and tell the user
  which agent stalled and on what.

## Deliverable

The final `design.md` — approved at both altitudes and editorially
audited — plus the full `review.md` ledger as the audit trail, and,
on escalate, the one paragraph the human must decide.
