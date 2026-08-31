# RFC Output Contract

When the design-review loop's deliverable is an RFC, this file is
the output contract: `design-investigator` shapes
`agent-team-workspace/design-docs/<slug>/design.md` to it, and `design-bar-raiser`
reviews against it. It also stands alone as the house definition of
a complete RFC.

## Required sections

**Proposed solution.** Every RFC clearly outlines the proposed
solution: after this section alone, a reader knows what is being
proposed and what changes if it ships.

**Pros and cons.** A thorough discussion of the tradeoffs of the
chosen solution, and of every steel-manned alternative.

- Alternatives are steel-manned: each presented in the strongest
  form its advocate would recognize — real advantages first — then
  killed with a stated reason tied to a requirement or constraint.
- The chosen solution carries real cons. Every design pays
  something; a tradeoff section where the winner has no genuine
  cost is advocacy, not analysis.

## Technical discussion — high level only

The technical discussion of the proposed solution is limited to
high level details. The areas below are the menu; an RFC includes
the subset that applies, and only that subset — not all RFCs
include all areas.

- **APIs** — proposed proto service definitions: RPCs with the
  supporting proto request and response messages (or the project's
  equivalent interface definitions).
- **Data model** — proposed core data models, particularly those
  shared between components.
- **Storage** — proposed storage solutions, each with the
  justification for choosing that storage solution.
- **Business logic** — a high level description of the business
  logic required.
- **Example** — an application of the proposed solution to the
  problem the RFC is addressing.

## Scope rules

- High level is a ceiling, not a floor to pad toward. Proto and API
  definitions, data models, and storage choices are as deep as an
  RFC goes; function-level signatures beyond the API definitions,
  algorithm internals, and concurrency mechanics belong in a
  follow-up design doc, not the RFC.
- Omitting an area is fine where it does not apply; omitting an
  area that is load-bearing for the problem is a review objection.
  The example is the cheapest proof the proposal is understood, and
  is rarely the right area to omit.
