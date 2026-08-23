# Brief Output Contract

A brief is the frozen input the loops build on — `/scoping` produces
one with the user, `design-loop` and `pr-loop` consume it. This file
defines what a finalized brief contains; a brief that meets it and
carries the user's sign-off satisfies the loops' "verbatim from the
user" requirement, because the user co-authored and approved every
line.

## Header

- `Status: draft` while scoping is underway; `Status: signed-off
  <date>` once the user approves. Loops consume signed-off briefs
  only. After sign-off, changes arrive solely as the loops' dated
  `## Amendment` sections.

## Required sections

**Problem.** Problem-shaped, not solution-shaped: what hurts or is
missing, for whom, and what happens if nothing is done. A requested
solution ("add a cache") appears here only as context; the
requirement underneath it ("p95 under X", "cost below Y") is what
binds.

**Goals.** Numbered and testable — each one states how you would
check it is done. "Improve performance" is not a goal; "the
scheduler run completes without manual retries" is.

**Non-goals.** What this work deliberately does not do. An empty
non-goals section means the scope was not challenged.

**Constraints and invariants.** What must stay true: technical
(from the codebase, cited `path:line`), and imposed (compatibility,
deadlines, tools). A user-insisted solution choice lands here,
marked as imposed rather than derived.

**Decomposition.** The problem broken into independently
deliverable pieces, ordered by dependency and risk. One line each:
`piece | goal it serves | route (design-loop / pr-loop / no code) |
depends on`. A piece too big for its route gets split before
sign-off, not discovered mid-loop.

**Open questions.** Each with an owner and a default that applies
if unanswered. A question with no owner and no default blocks
sign-off.

**Context.** Pointers into the codebase (`path:line`) and to any
research that grounds the above — enough that a loop agent starts
reading in the right place.

## Quality bar

- Every goal traces to the problem; every decomposition piece
  traces to a goal. Anything that traces to nothing gets cut.
- The null option was considered: the brief can say why doing
  nothing loses.
- Nothing in the brief contradicts the codebase as read; where the
  user's belief and the code disagreed, the brief records what the
  code actually does.
