# PR Loop — Agent Team Prompt

Goal: turn an engineering design or requirements brief into a
high-quality Pull Request for human review, with three agents —
`coding-expert` implements, `code-bar-raiser` challenges for up to
**4 rounds**, and `ai-writing-auditor` audits the PR description.
The output contract is `docs/pr-spec.md`; the human review the PR
receives afterward is the backstop, not the excuse.

Both engineering agents think from first principles and are defined
in `.claude/agents/` (hardened masters in `docs/agents-hardened/`).
A full loop is deliberately expensive; use it for changes that
deserve a real review, and use `golang-pro`/`rust-pro` directly for
quick edits.

## Invocation

The `/pr-loop` skill (`.claude/skills/pr-loop/`) is the entry point
that wraps this protocol; invoking it or writing the prompt below
are equivalent.

> Run the PR loop on: <brief — the requirements or design doc
> (a design-loop `design.md` works as-is), target area of the
> codebase, constraints>. Slug: `<slug>`. Open PR: yes/no.

## State: single writer per file

| Path | Sole writer | Role |
|---|---|---|
| `docs/pr-loop/<slug>/brief.md` | lead (round 0 + appended amendments) | requirements |
| `docs/pr-loop/<slug>/pr.md` | coding-expert | PR title/description, keyed revision log, objection responses |
| `docs/pr-loop/<slug>/review.md` | code-bar-raiser | append-only round ledger with verdicts |
| `docs/pr-loop/<slug>/pr.rewritten.md` | ai-writing-auditor | editorial intermediate |
| branch `pr/<slug>` | coding-expert | the code, in digestible commits |

The files plus the branch are the entire loop state; both
engineering agents are stateless and re-read them on every
invocation. The lead writes `brief.md` before round 0 (verbatim
from the user; amendments only by appended dated `## Amendment`
sections on the user's instruction, standing text never rewritten —
an amendment is revised material and both agents read requirements
as brief plus amendments, latest winning). The lead never touches
the other files or the branch.

Loop-state commits on the branch are prefixed `[loop]` so reviewers
can skip them; the branch's final commit removes `docs/pr-loop/<slug>/`
so the PR's net diff carries only the change itself.

## Protocol (lead session)

1. **Round 0 — implement.** Write `brief.md`, create branch
   `pr/<slug>` from the default branch, then invoke `coding-expert`
   with the slug. It implements in digestible commits, gets the
   suite green, and writes `pr.md`. If it returns `blocked` or
   `split proposed`, stop and put that to the user.
2. **Round N (1..4) — review.** Invoke `code-bar-raiser` with the
   slug and round number. It derives independently from the brief,
   checks out and runs the branch, appends `## Round N` with a
   `Verdict:` line to `review.md`:
   - `approve` / `approve-with-risks` → step 4.
   - `revise` → step 3.
   - `escalate` → **stop**; put the ledger's escalation paragraph
     to the user.
3. **Respond and revise.** Invoke `coding-expert` with the slug and
   round number. It answers every objection by ID — accepted (with
   the resolving commit), rebutted (with evidence), deferred (with
   the follow-up logged) — in new commits, never rewriting pushed
   history, and logs `R<N>:` in `pr.md`. Back to step 2 as round
   N+1.
4. **Editorial pass — after approval only.** Invoke
   `ai-writing-auditor` on `docs/pr-loop/<slug>/pr.md`; it writes
   `pr.rewritten.md` and returns a claim-inventory report. Clean →
   `coding-expert` adoption pass (meaning-diff, adopt, `editorial:`
   entry; corrections reported route the drifted sections back to
   the bar-raiser). Failed rewrite → keep `pr.md` as approved and
   tell the user. The auditor never touches code or comments —
   comment quality was the bar-raiser's job.
5. **Close.** Final commit removes `docs/pr-loop/<slug>/` from the
   branch; push. If the invocation said `Open PR: yes`, open the PR
   (draft where supported) with `pr.md`'s title and body — never
   open one otherwise. Report to the user: verdict, branch, PR link
   or ready-to-open state, residual risks, deferred follow-ups, and
   the human-side etiquette steps from `docs/pr-spec.md` (post to
   the team channel, cc reviewers).
6. **Round 4 is the floor for a decision, not a target.**

## Loop rules

- Relay only slugs, round numbers, and verdicts; substance travels
  through the files and the branch — never paraphrased prompts.
- No agent plays another's role: the author never writes
  `review.md`; the reviewer never writes code, `pr.md`, or
  `brief.md`; the auditor touches `pr.md` prose only. Disagreement
  between author and reviewer is signal — surface it, never smooth
  it.
- Push the branch after every round; the pushed branch plus the
  state directory is the only durable checkpoint.
- If a round produces no visible change (no new commits or
  dispositions, no new ledger section), the loop is stuck: stop and
  tell the user which agent stalled and on what.

## Resuming an interrupted loop

The slug directory and the branch are the checkpoint; files outrank
any resume prompt. Read them and take the first matching state:

1. No `brief.md` → never started: get the brief from the user.
2. No branch `pr/<slug>`, or no `R0:` entry in `pr.md` → round 0
   pending: invoke coding-expert.
3. Last `review.md` verdict `revise`, no `R<N>:` entry for that
   round in `pr.md` → response pending: invoke coding-expert.
4. Same verdict, `R<N>:` entry exists → invoke code-bar-raiser as
   round N+1.
5. Verdict `approve`/`approve-with-risks`, no `pr.rewritten.md`,
   no `editorial:` entry → editorial pass pending.
6. `pr.rewritten.md` present, no `editorial:` entry → adoption
   pending.
7. `editorial:` entry present → close per step 5.
8. Verdict `escalate` → closed pending the human.

A round section without a `Verdict:` line is an incomplete round:
re-invoke the bar-raiser for that round number.

## Deliverable

A branch whose net diff is the change alone, a PR description that
survived the claim inventory, a review ledger as the audit trail —
and a PR the human reviewer can trust was already held to the bar.
