# Design-Review Loop — Agent Team Prompt

Goal: produce a high-quality engineering design (or implementation
plan) with three agents — `research-investigator` authors,
`design-bar-raiser` challenges for up to **5 rounds** until approval
or escalation, and `ai-writing-auditor` gives the approved document
a final editorial pass.

All three are defined in `.claude/agents/`. The investigator and
bar-raiser are expert software engineers at both design altitudes —
high-level (architecture, boundaries, data models, consistency,
scaling, failure domains) and low-level (interfaces, data
structures, algorithms, concurrency, error semantics) — and both run at
`effort: max`, so a full loop is deliberately expensive: use it for
designs that matter. Both run `model: opus`, capped at
`maxTurns: 20`. Budget in tens of minutes per round — a measured round
0 plus one review round on a small greenfield design ran about 27
minutes and 190k tokens at `opus`. That measurement predates the
20-turn cap, so a round can now stop earlier than it did: if a design
or a review arrives with sections named unfinished, that is the cap
rather than the agent giving up, and the fix is to raise `maxTurns`,
not to re-prompt. For a cheap one-shot
opinion on a human-authored design, use `architect-reviewer`
instead.

## Invocation

The `/design-loop` skill (`.claude/skills/design-loop/`) is the
entry point that wraps this protocol; invoking it or writing the
prompt below are equivalent.

> Run the design-review loop on: <brief — the problem, requirements,
> and any design doc or code paths to start from>. Topic slug:
> `<slug>`. Deliverable: `RFC` (optional).

When the deliverable is an RFC, `agent-team-workspace/agent-specs/rfc-spec.md` is the shared
output contract: the investigator shapes `design.md` to it and the
bar-raiser reviews against it — including its high-level-only cap,
which overrides the loop's usual demand for low-level depth in the
document (the research still goes as deep as the proof requires).

## State: single writer per file

| Path | Sole writer | Role |
|---|---|---|
| `agent-team-workspace/design-docs/<slug>/brief.md` | lead (round 0 + appended amendments) | requirements |
| `agent-team-workspace/design-docs/<slug>/design.md` | research-investigator | the design |
| `agent-team-workspace/design-docs/<slug>/review.md` | design-bar-raiser | append-only round ledger |
| `agent-team-workspace/design-docs/<slug>/design.rewritten.md` | ai-writing-auditor | editorial intermediate |

The brief, the design, and the ledger are the entire loop state.
The investigator and bar-raiser are stateless between rounds and
re-read the files on every invocation, so fresh invocations and
resumed ones behave identically. The lead writes `brief.md` before
round 0 — the problem, the requirements, the deliverable (RFC or
design doc), and any constraints, verbatim from the user — and
after that touches it only to append a dated `## Amendment` section,
only when the user changes the requirements; standing text is never
rewritten or deleted. An amendment is revised material: it licenses
new blocking objections without goalpost-drift penalty, and both
agents read the requirements as brief plus amendments, the latest
winning on conflict. Both agents take the
requirements from `brief.md`, never from the invoking prompt or from
the design's restatement of them.

## Protocol (lead session)

1. **Round 0 — author.** Write `brief.md` first, verbatim — if it
   already conforms to `agent-team-workspace/agent-specs/brief-spec.md` (`Status: signed-off`,
   e.g. from `/scoping`), keep it exactly as given; do not
   re-author a scoped brief. Then invoke `research-investigator`
   with the slug. It researches from first principles, designs at
   both altitudes, and writes `design.md` (with `## Revision log`
   and `## Objection responses` sections).
2. **Round N (1..5) — challenge.** Invoke `design-bar-raiser` with
   the slug and round number. It derives the requirements
   independently from `brief.md`, challenges the high-level design and verifies the
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
   `ai-writing-auditor` on `agent-team-workspace/design-docs/<slug>/design.md`. It
   writes the cleaned prose to `design.rewritten.md` (it never
   overwrites its source) and returns a claim-inventory report.
   - Claim inventory clean → invoke `research-investigator` for the
     adoption pass: it diffs the rewrite against the design for
     technical meaning, replaces `design.md`'s prose with the
     audited prose, and logs one `editorial:` entry. If it returns
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
- An RFC brief changes the document, not the protocol: rounds,
  verdicts, objection IDs, and the editorial pass all run the same.
- When an agent reports a brief correction — the brief asserts
  something the codebase contradicts — put it to the user with the
  evidence and offer the `## Amendment` that fixes it. A factual
  correction is not a requirements change and does not reopen the
  design, but leaving it unamended makes every later agent
  re-derive it, and makes the frozen brief a source of falsehood.
  Never amend on your own; the brief is the user's.
- The lead owns every commit: after each round it stages and
  commits `agent-team-workspace/design-docs/<slug>/`, and pushes where the session's
  git conventions allow. The agents write files and never commit —
  an uncommitted round does not survive the session.
- If a round produces no visible change (no revision log entry, no
  new review section), the loop is stuck: stop and tell the user
  which agent stalled and on what.

## Resuming an interrupted loop

The slug directory is the checkpoint; no conversation context is
needed. When the files and a resume prompt disagree — on the round
number, the verdict, anything — the files win. To resume, read all
files in the slug directory and take the first matching state:

1. No `brief.md` → the loop never started. Get the brief from the
   user; do not reconstruct it from memory.
2. No `design.md`, or an empty revision log → round 0 pending:
   invoke the investigator.
3. Last `review.md` verdict is `revise` or `reject-approach`, and
   the revision log has no `R<N>:` entry for that round → the
   response is
   pending: invoke the investigator (step 3).
4. Same verdict, but the `R<N>:` revision log entry exists →
   invoke the bar-raiser as the next round (step 2).
5. Verdict `approve`/`approve-with-risks`, no `design.rewritten.md`,
   no `editorial:` entry in the revision log → editorial pass
   pending:
   invoke the auditor (step 4).
6. `design.rewritten.md` present, no `editorial:` entry → adoption
   pending: invoke the investigator's adoption pass.
7. `editorial:` entry present → the loop is closed; delete a leftover
   `design.rewritten.md` and report the deliverable.
8. Verdict `escalate` → closed pending the human: put the ledger's
   escalation paragraph to the user.

A `review.md` round section with no `Verdict:` line is an
incomplete round: re-invoke the bar-raiser for that same round
number and say the section is unfinished.

## Deliverable

The final `design.md` — approved and editorially audited; when an
RFC was requested, it is the RFC, shaped by `agent-team-workspace/agent-specs/rfc-spec.md` — plus the full `review.md` ledger as the audit trail, and,
on escalate, the one paragraph the human must decide.
