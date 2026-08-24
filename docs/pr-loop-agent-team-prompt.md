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
as brief plus amendments, latest winning). The lead never writes
`pr.md`, `review.md`, or any code — its only work on the branch is
creating it, committing loop state, and pushing.

The lead owns every loop-state commit: after each round it stages
`docs/pr-loop/<slug>/`, commits it with a `[loop]` prefix so
reviewers can skip those commits, and pushes. Agents never commit
loop state — `coding-expert` commits code only, `code-bar-raiser`
commits nothing at all. The branch's final commit removes
`docs/pr-loop/<slug>/` so the PR's net diff carries only the change
itself.

Every agent in this loop shares the lead's working tree. The lead
creates `pr/<slug>` and leaves it checked out for the whole loop,
and is the only one that may change which branch is out; the agents
confirm the branch and stop rather than switching it.

## Protocol (lead session)

1. **Round 0 — implement.** Write `brief.md` — if it already
   conforms to `docs/brief-spec.md` (`Status: signed-off`, e.g.
   from `/scoping`), keep it exactly as given; do not re-author a
   scoped brief. Create branch `pr/<slug>` from the default branch
   and leave it checked out, then invoke `coding-expert` with the
   slug. It implements in digestible commits, gets the suite green,
   and writes `pr.md`. If it returns `blocked` or `split proposed`,
   stop and put that to the user.
2. **Round N (1..4) — review.** Invoke `code-bar-raiser` with the
   slug and round number. It derives independently from the brief,
   builds and runs the branch you left checked out, appends
   `## Round N` with a `Verdict:` line to `review.md`:
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
   open one otherwise. Use whatever GitHub access this session
   actually has: the `gh` CLI is absent in some environments, and
   the GitHub MCP tools are the fallback. If neither is available,
   stop with the branch pushed and hand the user the compare URL to
   open it themselves. Report to the user: verdict, branch, PR link
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
- After every round, stage and commit `docs/pr-loop/<slug>/` as a
  `[loop]` commit and push the branch — the pushed branch carrying
  both the code and the state directory is the only durable
  checkpoint, so a round whose ledger was never committed is a
  round that did not survive.
- If a round produces no visible change (no new commits or
  dispositions, no new ledger section), the loop is stuck: stop and
  tell the user which agent stalled and on what.

## Resuming an interrupted loop

The slug directory and the branch are the checkpoint; files outrank
any resume prompt. Read them and take the first matching state:

0. Branch `pr/<slug>` exists but `docs/pr-loop/<slug>/` does not →
   the loop already closed, because step 5 removes that directory.
   Never restart it: check whether the PR is open, finish step 5 if
   it is not, and report. A branch that already carries the work
   never needs a fresh brief.
1. No `brief.md` and no branch → never started: get the brief from
   the user.
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
