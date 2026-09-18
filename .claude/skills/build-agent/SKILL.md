---
description: Create new subagents, repo skills, and multi-agent team loops using this repo's hardened process — role definition by negative space, verified frontmatter and verified runtime, minimal tool grants, mechanical work in scripts rather than prose, file-based loop state with single writers, checkpoint/resume, convergence discipline, a first-principles hardening review, and one live run before shipping. Use when the user wants to create an agent or a skill, design a pair or team of agents, build an agent loop or workflow, or harden existing ones. Not for running the design-review loop itself (use run-design-loop).
argument-hint: <what the agent or team should do> [name(s)] [loop?]
---

You are building agents and agent-team loops the way this repo's
design-review loop was built. The exemplars are the canon — read the
relevant one before writing anything, and copy its shape, not its
content:

- Single agent: `.claude/agents/design-investigator.md`
- Adversarial counterpart: `.claude/agents/design-bar-raiser.md`
- Loop protocol: `agent-team-workspace/protocols/design-review-loop-agent-team-prompt.md`
- Shared output contract: `agent-team-workspace/agent-specs/rfc-spec.md`
- Entry-point skill: `.claude/skills/run-design-loop/SKILL.md`
- Skill backed by a script: `.claude/skills/draw-diagram/SKILL.md` —
  the model keeps the judgement,
  `.claude/skills/draw-diagram/scripts/build_excalidraw.py` owns
  every number

## Step 0 — Survey before writing

Check `.claude/agents/` and the built-in
agent types. If an existing agent covers the role, extend or harden
it; do not create a near-duplicate. If the request is a team, first
decide whether a team is warranted: a loop costs one context per
agent per round and adds failure modes a single session does not
have. Recommend one agent when the work is sequential, touches the
same files, or fits one context. The cheapest correct design is
often one agent — say so when it is true.

When the user points at an external reference — a skill or agent from
another repo — review it before copying anything. Say what to keep
(its method: what it decides, in what order) and what to drop (its
mechanism and its dependencies), and say which is which. A
reference's structural choices are not requirements. Precedent: the
Excalidraw reference behind `draw-diagram` forbade a generator
script and hand-typed every coordinate; this repo's version keeps
its design method and moves every number into a compiler, because
hand arithmetic is where this repo's own agents slip.

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
- The same rule covers anything the agent or its scripts talk to: a
  file format, a protocol, a CLI. Read the implementation, not your
  memory of it — `draw-diagram`'s font ids, roundness types, and
  load-time behaviour came from grepping the Excalidraw package
  pulled from npm — and label what you could not find as inferred,
  with what it affects if wrong.

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

## Step 3b — Prose or program

Before writing an instruction, ask whether a script could do the
work. Arithmetic, geometry, id bookkeeping, file formats, and any
checklist item that can be computed go in `scripts/`, and the
SKILL.md tells the model to run it. The model keeps the judgement
(what to draw, what to say, what the reader must see) and gives up
the bookkeeping (where, how wide, which id binds to which). The
reference behind `draw-diagram` needed a three-phase workflow to get
a few hundred hand-typed coordinates past the output limit, and a
27-item checklist to catch what came out wrong; a compiler made most
of that checklist mechanical and the workflow one command.

Rules for a script:

- Stdlib only unless the repo already depends on the package, so
  the first run needs nothing installed.
- Probe the machine before designing around it — which binaries,
  which network hosts — and fail soft with a message naming what is
  missing. The reference fetched Excalidraw from a CDN on every
  render; that host is blocked by this session's egress policy, so
  its mandatory validation step could not run here at all.
- Measure environment quirks, never hardcode them. Chromium's new
  headless mode returned a PNG shorter than the canvas; the fix
  reads the PNG header and reshoots with the difference, not a
  magic 86 px that the next Chromium changes.
- Same input, same bytes. Seed anything random from stable ids so
  a rebuild diffs cleanly.
- Tests beside it (`test_<name>.py`, stdlib `unittest`), including
  one that executes every example in the skill's references. That
  test found a defect in
  `.claude/skills/draw-diagram/references/visual-patterns.md`
  minutes after it was written: a label overlapping the box it
  pointed at, in a doc that claimed every snippet built clean.

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
- **The brief is evidence, not truth.** Give the author a first step
  that verifies the brief's factual claims against the codebase and
  returns corrections; the lead appends them as a dated amendment.
  In this repo's first live design-loop run the lead's brief carried
  four false statements in a page of text — a document quoted for
  numbers it did not contain, a "standard library only" claim beside
  an `import yaml`, a timing figure for the wrong code path, and a
  goal whose check already existed — all caught by the investigator,
  none by the lead who wrote it.

## Step 5 — Deploy

- **One copy, at `.claude/agents/<name>.md`.** That is the only
  path Claude Code loads, and git already holds the history, the
  authorship, and the rollback. This repo previously kept a second
  "hardened master" mirror; 7 of its 20 files had silently
  diverged before anyone noticed, and the one file it did not
  duplicate was the hook nothing could resolve. If you need a
  staging area for hardening you are not ready to deploy, use a
  branch — a branch cannot silently diverge, because that is what
  branches are for. The same holds for skills at
  `.claude/skills/<name>/SKILL.md`.
- **Build order.** Write `SKILL.md` before anything else in a new
  skill directory. The PostToolUse hook validates after every write
  under `.claude/`, and a directory without its `SKILL.md` fails
  that validation on every subsequent write until it exists. Use
  the Write tool for large files there: on failure the hook echoes
  a Bash command's full text back, heredoc included.
- Register: the entry-point skill for loops, and CLAUDE.md's
  skills/agents lists so fresh sessions can discover it. Names,
  never counts: CLAUDE.md said the validator ran "231 checks" and
  was wrong the moment the next skill landed. A number a program
  computes does not belong in prose; the program prints it.
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
   the return contract needs a slot to say so. The mechanism is
   write-early: the agent writes its file first and refines it in
   place, so a cutoff leaves a partial file rather than nothing.
   Round 1 of this repo's first live design-loop run hit the cap
   with its work only in the unsent final message and lost the
   round; with write-early in the prompt, every later truncation
   kept its work. Put the rule in the agent definition, not the
   invoking prompt — as of this writing it lives only in the lead's
   prompts. And a turn is a model response, not a tool call: 24
   tool calls fit under a 20-turn cap once and 21 did not, so size
   the cap with margin rather than by counting calls.
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

**Pass C — run it once, for real.**

14. **A live run on the smallest real input.** Passes A and B and
    the validator were all green before this repo's first live
    design-loop run. The run found six defects none of them could
    see: a turn cap changed in frontmatter and not in the skill
    prose that described it; a check that verified 0 of 8 skills
    while printing `0 failures`; four false statements in the
    lead's brief; a round claim hidden by literal spaces where text
    wraps; English role nouns as a rename surface no check tracked;
    and an arithmetic slip that survived two review rounds
    (`agent-team-workspace/design-docs/prose-config-drift/design.md`
    line 648 splits 10 claims as 5 and 4). Static review reads what
    a definition says; only a run shows what it does.
15. **Look at the artifact, not the exit code.** `draw-diagram`'s
    first PNG was written, non-empty, and clipped at the bottom;
    the tool exited 0. Only opening the image showed it. Whatever
    the component produces — a file, a ledger, a render — open it.

All three passes are judgement. What is mechanical, run instead:
`python3 agent-team-workspace/validate-definitions.py` checks
frontmatter parses, agent names are unique (duplicates load by
filesystem read order, not precedence), hook targets resolve and are
executable, no repo path dangles, every "(use X)" boundary names
something real, no skill shadows a bundled one, resume derivations
are contiguous, no reference file is orphaned, and CLAUDE.md's
roster matches the tree. `python3 agent-team-workspace/validate-skills.py`
is the per-skill pass. Both exit non-zero on failure. Run them
before you ship and after every fix.

Three lessons from maintaining those checks, each a real defect:

- **A check tolerates the states earlier checks report.** The
  validator raised a traceback on a skill directory without its
  `SKILL.md`, and the traceback took every later check with it:
  check 4 had already recorded the failure; checks 5 and 12 opened
  the missing file anyway. A check that raises on a state already
  flagged hides everything after it.
- **Skipped is not passed.** `validate-skills.py` carried a
  git-based line-count check whose precondition stopped holding in
  the commit that introduced it; it skipped 6 of 8 skills, compared
  the other 2 to themselves, and printed `0 failures` every run
  until it was deleted. A check that cannot run says so in the
  summary line.
- **Judge by structure, not prior knowledge.** The recurring defect
  in this repo's checks is a hardcoded allowlist, a rename map, or
  a regex that matches a name; each goes stale silently. Derive the
  roster from the filesystem and frontmatter at run time. One such
  list remains: each validator keeps its own set of bundled skills
  a description may route to, so a new `(use X)` toward a bundled
  skill edits both, and the day they disagree one of them is wrong.

Fix the majors before shipping — then re-run Pass A on the files
you just edited. Fixes introduce defects at a high rate: in this
repo's second review round, two of three majors were damage from
the first round's fixes. Report anything left open, ranked, with
the concrete fix for each — the user decides, and an issue reported
honestly beats one padded over.
