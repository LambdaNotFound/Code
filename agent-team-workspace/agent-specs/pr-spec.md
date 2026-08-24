# PR Output Contract

When the PR loop runs, this file is the output contract:
`coding-expert` builds to it and `code-bar-raiser` reviews against
it. It distills the team's PR etiquette (background: Google's code
review guidelines) into the parts agents can enforce; the
human-process parts are listed at the end for the lead to hand back
to the user.

## Scope

- PRs are small and easily digestible by the reviewing audience —
  one concern per PR. When a brief implies more than one digestible
  PR, the loop stops and proposes a split before implementing; it
  never ships an unreviewably large diff because the brief asked
  for a lot.
- Vendor and generated-file changes ship in a separate atomic PR,
  or at minimum in their own clearly marked commit that reviewers
  can skip.

## Commit history

- Digestible, atomic commits that tell the story of the change in
  review order: each commit builds, each message says why.
- During review rounds, pushed history is never rewritten —
  objections are answered with new commits so the reviewer sees
  exactly what changed since their last pass. The human may squash
  on merge.

## PR description (`pr.md`)

- **Title** — imperative, specific, small enough to be true.
- **What and why** — the problem and the change, tied to the brief;
  a reader who was not in the loop understands both.
- **How tested** — the actual commands run and their results. Never
  a test that was not run.
- **Risks and rollback** — what could break, how to undo.
- **Review focus** — where the reviewer's attention buys the most.
- Prose is audited (ai-writing-auditor) before the PR opens; code
  and comments are not — comment quality is the reviewer's.

## Tests

- New behavior is covered by tests that assert behavior, not
  implementation. The full suite is green before every review
  round — run by the author, re-run by the reviewer.

## Review standards (per round)

- Order: correctness, tests, simplicity (no more complex than
  requirements force), system fit, readability and comments.
- Severity: `blocking` / `should-fix` / `nit`; every objection
  names what would resolve it.
- If the reviewer cannot follow the code, that is a blocking
  objection and the author clarifies before review continues —
  code the reviewer cannot read, other engineers cannot read.
- At least one substantive comment every round, approvals
  included; a blanket approval teaches nothing.
- The reviewer builds and runs the branch the lead left checked
  out; review from the diff alone does not count. Every agent
  shares the lead's working tree, so nobody but the lead changes
  which branch is out.

## Human-side etiquette (out of scope for agents)

The lead reminds the user of these in its final report; agents do
not perform them: posting the PR to the team's channel and cc'ing
reviewers, bumping for another look, requesting a short meeting
when discussion stalls, and the one-business-day review-response
norm.
