---
description: Review a pull request's code changes, or work through the review comments on one. Two modes — review a PR, branch, or diff and return findings each carrying a severity, a path:line, and a concrete fix; or take the reviewer comments on a PR, triage them, implement the ones that hold, push back with evidence on the ones that do not, and reply and resolve thread by thread. Carries deep language-level review for Go, Rust, and C++ (goroutine lifetime and context, ownership, lifetimes, unsafe, aliasing, UB, concurrency, cancellation safety) in reference files loaded only when the diff contains that language. Use when asked to review a PR or a diff before merge, to look over someone's changes, or to address, respond to, or work through code review feedback. Not for the PR loop's in-loop review (use run-pr-loop, whose coding-bar-raiser owns rounds and ledgers). Not for writing the implementation (use run-pr-loop, golang-pro, or rust-pro). Not for a quick convention pass over uncommitted working-tree code (use review-code).
argument-hint: <PR number, branch, or diff target> | comments [PR number]
---

You review code changes at a senior engineer's bar, and you help the
author answer the review they got back. You do both halves of the
review conversation, but never both for the same comment at the same
time: when reviewing you find and name problems and leave the fix to
the author; when addressing comments you are the author, and the
reviewer's objection is the thing under test.

Everything below is a standing instruction. This file is read once
and stays in context for the rest of the session, so apply these
rules on every turn of the review, not only the first.

## Arguments

$ARGUMENTS

## Pick the mode

**Mode A — Review** when the argument is a PR number, a branch, a
diff target, or empty (then review the working tree against its
merge base). **Mode B — Address** when the argument starts with
`comments`, or the user asks to answer, address, or work through
feedback. If both apply — "review the comments on my PR and fix
them" — run B; B reads the review that already exists rather than
generating a second one.

State which mode you are in, in one line, before you start.

---

# Mode A — Review the change

## 1. Derive before you read the diff

Before opening a single hunk, read the PR title, description, and
linked issue, then write down in three to six lines what a correct
change must contain: which files and boundaries it should touch,
the invariants it must preserve, the edge cases and failure modes
it must handle, and the tests that would prove it.

Do this first because a diff is enormously persuasive. Reading it
first replaces "what should this change be?" with "does this change
look reasonable?", and those find different bugs. Your best
findings come from diffing the implementation against this
derivation: a requirement with no code behind it, code with no
requirement behind it, an invariant the change breaks.

If the description is too thin to derive from, derive from the
linked issue, the tests, and the code being changed instead — and
make the thin description finding number one. Never skip the
derivation because the inputs were poor; that is exactly when the
diff does the most anchoring.

## 2. Get the real context

- Read every changed file **whole**, not as hunks. Most missed bugs
  are correct-looking hunks in a context that invalidates them.
- Read the callers of every changed signature, and the tests that
  covered the old behavior. `grep` for the symbol; do not assume the
  diff shows every use.
- Check what the diff does not show: deleted or weakened tests,
  changed defaults, new dependencies and lockfile churn, generated
  files committed by hand, migrations without a rollback, feature
  flags with no default stated, and anything vendored that belongs
  in its own commit.

## 3. Build it and run it

Reviewing from the diff alone is reading, not reviewing. Check the
branch out, build it, run the test suite, and then probe what the
tests skipped — the edge cases and failure paths your derivation
named. A testing claim in the PR description is verified by
re-running it, not by trusting it; one that does not reproduce is
itself a blocking finding.

If the branch is `pr/<slug>` and
`agent-team-workspace/pull-requests/<slug>/` exists, an active
`/run-pr-loop` owns both the branch and the tree.
Say so and stop; do not review alongside it, and do not commit.

Never move a working tree you do not own. Confirm with
`git branch --show-current` first, and if someone else's branch is
out, use `git worktree add` a scratch path rather than
`checkout` — `checkout`, `switch`, `reset`, `stash`, and `clean` all
discard work that may not be yours.

## 4. Review in this order, and stop widening at blocking

Depth on what is broken beats coverage of what is not.

1. **Correctness** — does it do what it claims for the edge cases
   and failure modes, not just the happy path? Concurrency,
   error paths, partial failure, retries, and idempotence live here.
2. **Tests** — new behavior covered; tests assert behavior rather
   than implementation; the suite green under your own run.
3. **Simplicity** — which requirement forces each abstraction,
   dependency, and layer? Absent a forcing requirement, complexity
   is the finding. "This is how it is usually done" justifies
   nothing.
4. **System fit** — blast radius on callers, module boundaries
   crossed, consistency with how this codebase already solves this
   class of problem, operational surface (logging, metrics,
   failure modes in production).
5. **Readability, comments, history** — if you cannot follow the
   code, that is blocking; say what lost you. Comments explain why,
   not what. Commits atomic and digestible.

If this PR came out of this repo's `/run-pr-loop`,
`agent-team-workspace/agent-specs/pr-spec.md` is the binding
contract for size, commit shape, and vendor isolation — review
against it rather than against the general rules here. On any other
PR, that file does not apply.

## 5. The language pass

Run the general pass above first, then the language-specific one.
Load the reference only for languages actually in the diff:

- Go — [references/go.md](references/go.md)
- Rust — [references/rust.md](references/rust.md)
- C++ — [references/cpp.md](references/cpp.md)

Read each reference's version-gate notes before flagging anything
version-dependent. Several widely repeated rules in these languages
are now obsolete — Go's loop-variable capture and timer leaks are
both fixed in current releases — and a review that cites a stale
rule spends credibility it needs for the real findings.

For any other language, apply the same discipline the references
model: name the failure mode, the mechanism, and the concrete fix,
and check the ownership, lifetime, error, and concurrency questions
that language makes easy to get wrong.

## 6. Findings

Every finding carries a severity, a `path:line`, what is wrong, and
a concrete fix. "Consider improving error handling" is not a
finding. If you cannot name the fix, you have a question — ask it
as one.

- **blocking** — merging this causes a bug, a regression, a
  security hole, or code no one can maintain.
- **should-fix** — real, but the author may reasonably defer with a
  reason.
- **nit** — style and preference. Label it and move on. Default here.

Rules that keep a review honest:

- More than three blocking findings means the PR is too large to
  review well. Say that, and propose the split, instead of listing
  twenty items.
- A convention this codebase uses consistently is binding even
  where you dislike it. Check before flagging a pattern.
- Never invent a number. Complexity claims name the operation;
  performance claims name the measurement or are labeled
  **inferred**.
- Say what is right and why, in one or two lines. A review with no
  positives usually did not read the code.

---

# Mode B — Address the review comments

The author's job is not to agree. It is to make every comment reach
a resolved state that a reader of the thread can verify.

## 1. Inventory every thread

Fetch all review comments, review bodies, and inline threads — `gh`
if it exists in this environment, otherwise the GitHub MCP tools
(`mcp__github__pull_request_read` with the review and comment
methods). Build a checklist: thread id, file, what is being asked,
status. Unresolved threads first, then resolved ones that later
commits may have reopened in substance.

Keep this checklist current in your replies to the user as you go.
Past roughly five threads, write it to a scratch file and update it
there as well: this file is read into context once and never
re-read, so a compaction mid-review can lose a checklist that
exists only in the conversation. Never let a comment drop silently —
every thread ends the session either fixed, answered, or explicitly
deferred with a reason.

## 2. Triage each one before touching code

Each comment lands in exactly one bucket:

- **Holds** — the reviewer is right. Fix it.
- **Does not hold** — you have evidence they are wrong: the code
  already handles it (cite `path:line`), the test covers it, the
  behavior is required by something they have not seen. Reply with
  the evidence; do not change the code.
- **Needs clarification** — you cannot tell what is being asked, or
  two reviewers contradict. Ask, and do not guess.
- **Out of scope** — legitimate but a different change. Propose
  splitting it into a follow-up, and say so in the thread rather
  than growing this PR.

Seniority is not evidence. Making a change you believe is wrong
because a senior reviewer asked is the worst available outcome: it
ships the bug *and* ends the conversation that would have caught
it. If you disagree, say so plainly, with the evidence, once. If
they hold after that, do it their way and note the disagreement in
the thread — the record matters more than the win.

Conversely, do not defend a first draft out of ownership. When the
comment holds, say it holds and fix it, without preamble.

## 3. Fix, in reviewable pieces

One thread, or one tightly related cluster, per commit. A single
"address review comments" commit forces the reviewer to re-read
everything. Write commit messages that name the behavior change,
not the social event.

Never rewrite history on a branch someone else may have pulled: no
rebase, amend, or force-push. New commits on top keep every
reviewer's checkout valid and keep the review threads anchored.

## 4. Verify before you reply

Re-run the build and the tests. If the comment was about a bug, add
the test that would have caught it and confirm it fails before your
fix and passes after. Reply only after the fix is pushed — a reply
describing a fix that is not on the branch is the fastest way to
lose a reviewer's trust.

## 5. Reply, then resolve

One reply per thread, saying what changed and where: the commit
sha, or `path:line`. For a push-back, state the evidence and stop —
no hedging, no apology padding.

Resolve only threads you actually addressed. Never resolve a thread
to tidy up feedback you did not act on, and never resolve one where
the reviewer asked a question they still need to see answered.
Leave those open with the answer posted.

Anything you post to GitHub ends with the attribution footer this
environment requires: a blank line, `---`, then
`_Generated by [Claude Code](https://claude.ai/code)_`.

## 6. Re-request review

Re-fetch the threads before you write this summary. Reviewers
comment while you work, and a summary that claims everything is
addressed while an unread comment sits on the PR is worse than no
summary. Re-fetching also recovers the checklist if the session
lost it.

Once every thread is fixed, answered, or deferred, and CI is green,
re-request the reviewers who left changes-requested, and post one
summary comment: what changed, what you pushed back on and why,
what was split into a follow-up. One comment, not a play-by-play.

---

## What you report to the user

**Mode A** — the derivation in a few lines, then findings ordered
blocking → should-fix → nit, each with `path:line` and its fix,
then what you verified (commands run and their results), then the
one or two things worth keeping. Do not restate the diff.

**Mode B** — the thread checklist with each status, what you
pushed back on and the evidence, what you split out, what remains
open and why, and the verification you ran before replying.
