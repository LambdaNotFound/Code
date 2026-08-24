---
description: Create new subagents and design multi-agent team loops using this repo's hardened process — role definition by negative space, verified frontmatter, minimal tool grants, file-based loop state with single writers, checkpoint/resume, convergence discipline, and a first-principles hardening review before shipping. Use when the user wants to create an agent, design a pair or team of agents, build an agent loop or workflow, or harden existing ones. Not for running the design-review loop itself (use design-loop).
argument-hint: <what the agent or team should do> [name(s)] [loop?]
---

You are building agents and agent-team loops the way this repo's
design-review loop was built. The exemplars are the canon — read the
relevant one before writing anything, and copy its shape, not its
content:

- Single agent: `docs/agents-hardened/research-investigator.md`
- Adversarial counterpart: `docs/agents-hardened/design-bar-raiser.md`
- Loop protocol: `docs/design-review-loop-agent-team-prompt.md`
- Shared output contract: `docs/rfc-spec.md`
- Entry-point skill: `.claude/skills/design-loop/SKILL.md`

## Step 0 — Survey before writing

Check `.claude/agents/`, `docs/agents-hardened/`, and the built-in
agent types. If an existing agent covers the role, extend or harden
it; do not create a near-duplicate. If the request is a team, first
decide whether a team is warranted: a loop costs one context per
agent per round and adds failure modes a single session does not
have. Recommend one agent when the work is sequential, touches the
same files, or fits one context. The cheapest correct design is
often one agent — say so when it is true.

## Step 1 — Define the role by its negative space

One agent, one job. The description is the router — agents are
selected on it alone — so it carries both halves: when to use this
agent (including proactive triggers), and explicit boundaries in
the form "Not for X (use Y)" naming the sibling that owns X. The
body opens with an identity paragraph in the same shape: "You do X.
You do not Y; Y belongs to Z."

## Step 2 — Frontmatter, verified not recalled

Fields to decide: `name`, `description`, `tools`, `model`,
`effort`, `maxTurns`, `memory`, `permissionMode`, `hooks`. Rules:

- Verify field names, allowed values, and model capabilities
  against the live Claude Code docs or the claude-api skill. Never
  configure from memory. If the user asks for a capability you
  cannot find documented, say plainly that it does not exist and
  map their intent to the nearest real knob (this repo's precedent:
  a request for an "endgame" ability became `model: fable` +
  `effort: max` after verification).
- Tools: the minimum the role needs. The tool list cannot scope
  paths, so a read-only role that must write its own reports gets
  Write plus a body rule naming its only writable paths.
- `memory`: grant only with a scope rule in the body — process
  lessons and codebase geography, never opinions, verdicts, or
  topic content — and state that files outrank memory on conflict.
- Match capability to stakes: highest model and effort for
  judgment-heavy daily-use roles; cheaper settings for mechanical
  ones.

## Step 3 — The body every hardened agent carries

In the house voice — terse, second person, every rule earns its
line:

1. Identity and refusals (what it is, what it will not do).
2. Method, ordered: how it works, in the sequence it should work.
3. Evidence and honesty rules: claims carry citations
   (`path:line`, command output, URL + date); statements labeled
   observed / inferred / assumed; no numbers nobody measured.
4. "What you return": only the final message survives the agent's
   context — everything the caller needs goes there, as an exact
   numbered contract, no preamble. Anything that must survive
   longer than the caller's context goes in a file, not the
   message.

## Step 4 — Team loops, when warranted

- **Adversarial pairing with hard role boundaries.** The author
  never reviews; the reviewer never authors; disagreement between
  them is signal to surface, not smooth over. Give the reviewer an
  independent-derivation step (derive from the frozen inputs
  before reading the proposal) and give the author an
  anti-capitulation rule (never accept an objection it can refute
  with evidence).
- **State is files, one writer each.** Append-only ledgers for
  reviewers; write-once input briefs with an explicit amendment
  channel (append-only, dated, user-authorized, counted as revised
  material). Everything a resume needs lives in the files:
  verdicts persisted in the ledger, log entries keyed (`R<N>:`,
  `editorial:`), because final messages die with the caller.
- **Stateless agents.** Each agent's file says: you may be invoked
  fresh at any round; the files are the state; read them in a
  fixed order (order matters wherever independence requires
  deriving before reading). Files outrank the invoking prompt and
  memory on any disagreement.
- **A lead protocol doc as single source of truth.** The lead
  relays only ids, round numbers, and verdicts — substance travels
  through files, never paraphrased prompts. Include a resume
  section: an ordered state derivation from the files alone.
- **Convergence discipline.** A round budget with mandatory landing
  states (approve / approve-with-risks / escalate); new blockers
  after round 1 must cite revised material or admit "missed and
  critical"; approval requires closure on evidence, never fatigue;
  escalation to the human is a designed exit, not a failure.
- **Wiring an existing agent in** (an auditor, a formatter): read
  its definition file first and design around its actual contract —
  output paths, guards, refusals. Never assume.
- **Entry point.** A thin wrapper skill that reads the protocol
  doc (restating nothing), handles naming and resume detection,
  freezes the input brief before round 0, and commits the state
  directory per round.

## Step 5 — Deploy

- Hardened master in `docs/agents-hardened/<name>.md`, identical
  live copy in `.claude/agents/<name>.md`; re-sync both on every
  edit.
- Register: the entry-point skill for loops, and CLAUDE.md's
  skills/agents lists so fresh sessions can discover it.
- Commit and push per the session's git conventions.

## Step 6 — Harden before shipping

Two passes, because they find different defects. Every item below
was a real defect found and fixed in this repo's own loops, and the
order matters: three review rounds there each turned up fresh
majors, and the class shifted every round — inside single agents,
then between an agent and its contract, then between loops. A
component-only review will not find the last two kinds.

**Pass A — each component alone.**

1. **Entry point** — can a fresh session with no context discover
   and start it?
2. **Contradictions** — can every pair of instructions in the file
   hold at once? (Reading order vs independence is the classic.)
3. **State channels** — does memory or anything else bypass "the
   files are the state"? Scope it or drop it.
4. **Amendment path** — can immutable inputs legally change when
   the user changes their mind?
5. **Resume** — is the next action derivable from files alone?
   Verdicts persisted, log entries keyed? Check the end states too:
   if closing deletes state, a resume must not read the absence as
   "never started" and redo finished work.
6. **Append safety** — can a whole-file write silently clobber
   append-only history? An agent told to append needs `Edit`;
   without it there is no append mechanism, only discipline.
7. **Turn budget** — what happens at `maxTurns`? A hard cutoff
   mid-work must degrade to saved work and a partial report, and
   the return contract needs a slot to say so.
8. **Correlated blind spots** — same model on both sides of an
   adversarial pair is procedural, not epistemic, independence;
   acknowledge it and keep the human escape hatch. Check the
   direction too: a reviewer weaker than the author it challenges
   undoes the premise.
9. **Copy drift** — are dual copies actually identical?

**Pass B — the system between components.**

10. **Restatements** — after changing any rule, grep every file
    that restates it. A rule fixed in the agent and left standing
    in its contract is worse than never fixing it, because the two
    now disagree and the agent reads both.
11. **Contract versus implementation** — read each contract file
    against the agents that consume it. A contract asserting what
    no agent enforces is decoration; one forbidding what they all
    allow is a trap.
12. **Seams** — for every handoff, does the producer's output
    actually satisfy the consumer's input contract, and is it
    addressed so two runs cannot collide? Decomposition is the
    usual failure: N pieces sharing one slug overwrite each other.
    Check that the best input is not penalized by a check written
    for the worst.
13. **Shared resources** — subagents share the lead's filesystem
    and working tree. Anything that switches branches, deletes, or
    resets is acting on everyone; name one owner and forbid the
    rest.

Fix the majors before shipping — then re-run Pass A on the files
you just edited. Fixes introduce defects at a high rate: in this
repo's second review round, two of three majors were damage from
the first round's fixes. Report anything left open, ranked, with
the concrete fix for each — the user decides, and an issue reported
honestly beats one padded over.
