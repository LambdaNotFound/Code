# Design — reconciling prose claims against agent and skill configuration

## Revision log

R0: initial draft. Claim taxonomy derived from a measured sweep of the
definition tree; extraction, subject binding, and reconciliation
specified; round-count claims taken in scope as symmetric agreement
classes rather than against an anointed source file; countable
inventories scoped out with a derivation.

R1: stage 3 rebuilt after executing the specified ladder over the whole
scan set — paragraph-unique and positional-anaphor rungs deleted as
measured false positives (R1-1, R1-2), replaced by role-token binding
with a cardinality-checked anaphor, with the binding decision for all 8
configuration claims tabulated. Existence checking (4b) deleted as
unforced (R1-4). Round counts switched to the brief's anointed-protocol
default with a classless disposition (R1-5, R1-6). Headline corpus
restated from the specified patterns rather than the prototype (R1-3).
Rename verification added to the plan and the hook suite (R1-7). F8 and
open question 4 dropped and citations re-verified at HEAD (R1-8).
Self-count coupling stated (R1-9).

## The question

The brief asks for one thing with four bindings: make a prose sentence
that states an agent's `model`, `effort`, or `maxTurns` fail a
validator run when it disagrees with that agent's frontmatter (goal 1),
without hardcoded names (goal 2), without a single false positive on
today's tree (goal 3), and with unparsed claims reported rather than
dropped (goal 4). This document decides the claim taxonomy — which
prose shapes are in scope and what each reconciles against — and
specifies the check that enforces it.

## Answer up front

**Three new checks in `validate-definitions.py` meet all four goals for
the 5 of 8 configuration claims whose subject binds; the other 3 are
reported, not checked.** The claim corpus measured by the patterns this
document specifies — not by the round-0 prototype, which was defective in
three ways (F1) — is 8 configuration claims and 12 round-budget claims
across the 33 files check 5 already scans, with no spurious hits
(`probe2.py`/`probe5.py`, run against HEAD `bfc9cec`; F1, F4, F5). That
smallness is what makes the design safe.

The guarantee is not uniform, and the caveats belong in the headline:

- A claim whose subject is a backticked agent name in the same sentence —
  the shape that actually went stale — is reconciled **exactly** against
  that agent's frontmatter. Drift fails the run. This covers
  `scope-problem/SKILL.md:18`, the brief's goal-1 target.
- A claim whose subject is a role noun ("the investigator and bar-raiser")
  or a plural anaphor ("Both run \`model: opus\`") is reconciled exactly
  **only when the noun identifies exactly one member of the file's own
  backticked roster, and only when an anaphor's count matches the number
  so resolved in the antecedent sentence**. Measured, that binds all
  three claims at `design-review-loop-agent-team-prompt.md:14-16`
  correctly. It is a heuristic and it declines rather than guesses.
- A claim whose subject does not bind is **reported as an advisory and
  not checked at all**. There is no existence fallback; it was cut for
  having no requirement behind it and a false-positive class of its own
  (R1-4). Three of today's 8 claims land here — `build-agent/SKILL.md:47`
  and `:48`, `design-review-loop-agent-team-prompt.md:19` — so the
  advisory baseline this design ships is **3, not 0**.
- Round counts have no frontmatter field. They are reconciled against the
  loop protocol that the claim's file resolves to, with the protocol's
  own statement anointed as the source (the brief's open-question-3
  default). Drift in the protocol alone still fails every other member,
  so the anointing costs no detection and buys a named culprit.

Cost is one extra pass over files the validator already reads; the added
wall time is unmeasured and named as a plan verification step.

## System map (observed)

Two validators, one hook, one test suite, one claim corpus.

- `agent-team-workspace/validate-definitions.py` — 287 lines, 13
  numbered check sections, 231 hard checks, exits non-zero on any
  failure (`validate-definitions.py:287`). Measured 150–159 ms over
  four runs including interpreter start; `231/231 hard checks passed, 0
  advisories` on HEAD. **observed**
- `agent-team-workspace/validate-skills.py` — 65 lines at HEAD `bfc9cec`,
  per-skill loop over the 8 skills, 0 failures. **observed**
- `.claude/hooks/validate-definitions.sh` — PostToolUse hook; filters
  on `.tool_input.file_path`, falling back to `.tool_input.command` for
  Bash, and runs the definition validator when the path or command
  names an agent, a skill, `CLAUDE.md`, or the validator
  (`validate-definitions.sh:35-52`). Fails open (`note()` at `:27-30`).
  **observed**
- `.claude/hooks/test-validate-definitions.sh` — 24 `run` cases against
  a throwaway copy under `/tmp`; helper contract is
  `run <desc> <expected-exit> <payload-json> [setup]`
  (`test-validate-definitions.sh:32-41`). **observed**
- The claim corpus is the file set check 5 already builds
  (`validate-definitions.py:106-109`): 13 agent files, 8 `SKILL.md`, 6
  reference files, 2 protocols, 3 agent-specs, `CLAUDE.md` — 33 files.
  **observed**
- Configuration lives in agent frontmatter: 13 agents, all carrying
  `model` and `maxTurns`, 4 carrying `effort`. Values in use today:
  `model` ∈ {`sonnet`, `inherit`, `opus`, `fable`}, `effort` ∈ {`max`,
  `xhigh`}, `maxTurns` ∈ {15, 20, 30, 40, 50, 60}. **observed**

## Findings

**F1 — The claim corpus is 20 sentences, and a noun-anchored trigger
finds exactly them.** The six patterns specified below, run over check 5's
33-file scan list, fire on 8 configuration claims and 12 round-budget
claims and on nothing else. The 8: `build-agent/SKILL.md:47` and `:48`
(`` `model: fable` ``, `` `effort: max` ``), `scope-problem/SKILL.md:18`
(two claims — `effort=max` and `maxTurns=twenty` in one sentence),
`design-review-loop-agent-team-prompt.md:14`, `:15`, `:16`, `:19`. The 12
round claims are listed in F5. **observed**

The round-0 prototype reported 19 across 30 files, and all three
discrepancies were prototype defects, not corpus facts: it scanned line by
line and missed the claim wrapping
`pr-loop-agent-team-prompt.md:5-6` (F3); it added `+ 1` to the
frontmatter offset and reported `scope-problem/SKILL.md:19` for a line-18
claim; and it dropped a capture group so two round claims extracted
`None`. The scan list is 33 files — the validator itself prints
`33 files scanned` (`validate-definitions.py:119`). Every number in this
document is now from a run of the specified patterns. **observed**

**F2 — Naked number proximity is unusable; the property noun is what
does the work.** A sweep for numerals near configuration keywords
returned hundreds of hits — ordinary "one"/"two"/"three", and every
numbered list item in every agent's return contract
(`multi-agent-coordinator.md:164-173`, `rust-pro.md:118-129`,
`golang-pro.md:56-63`). The decisive case is
`multi-agent-coordinator.md:47`: "Concurrent subagents are capped (20
by default)" — a bound word, the number 20, inside an agent file whose
own `maxTurns` is 20. It is about a different platform limit entirely.
The noun-anchored trigger does not fire on it because "turns" is
absent. **observed**

**F3 — Claims wrap across source lines, so line-by-line extraction
under-detects.** `pr-loop-agent-team-prompt.md:5-6` reads
"`coding-bar-raiser` challenges for up to\n**4 rounds**". The probe
scanned line by line and missed it — it is the one round-budget claim
absent from the probe output despite being cited in the brief.
Paragraph normalization is required, not optional. **observed**

**F4 — Every proximity-based and position-based binding rule tested
produces a false positive on the exact lines the brief requires to
pass.** `design-review-loop-agent-team-prompt.md` lines 9-24 are one
paragraph carrying four configuration claims (`:14`, `:15`, `:16`,
`:19`) and exactly one backticked agent name, `architect-reviewer` at
`:23` — a name that appears only to say what to use *instead* of the
loop. Executed over the scan set: **observed**

| rule | subject it computes for `:15` | outcome |
|---|---|---|
| every name in the section | `{design-investigator, design-bar-raiser, ai-writing-auditor}` | fails: `ai-writing-auditor` is `model: sonnet` |
| unique name in the paragraph | `{architect-reviewer}` | fails: `model: inherit`, no `effort`, `maxTurns: 15` — 4 false positives, one per claim |
| last 2 names by position in the section | `{design-bar-raiser, ai-writing-auditor}` | fails: `ai-writing-auditor` again |
| nearest 2 names by distance | `{ai-writing-auditor, design-bar-raiser}` | fails: same |
| first 2 names by position | `{design-investigator, design-bar-raiser}` | passes, but only because this section names three agents in loop order; nothing derives it |

The correct subject is stated in prose the rules do not read: "The
investigator and bar-raiser are expert software engineers … and both run
at `effort: max`" (`:9-14`). What the rules can read is that
"investigator" and "bar-raiser" are kebab tokens naming exactly one
member each of this file's own backticked roster. That is the rule F8
measures. **observed**

**F5 — The two loop-round classes partition cleanly with no hardcoded
names.** Every file carrying a round-budget claim maps to exactly one
protocol, via backticked-roster membership (the parse check 11 already
does, `validate-definitions.py:229-234`) or via citing that protocol's
path. Measured: `design-bar-raiser.md` → design roster;
`design-investigator.md` → design roster and design path;
`coding-bar-raiser.md`, `coding-expert.md` → pr roster and pr path;
`run-design-loop/SKILL.md` → design path only; `run-pr-loop/SKILL.md`
→ pr path only. The two genuinely ambiguous files —
`ai-writing-auditor.md` (both rosters) and `CLAUDE.md` (both paths) —
carry no round-budget claim, so nothing is lost. **observed**

**F6 — The hook's sandbox stubs out the source trees, which kills any
check against `rust/` or `spaced_repetition/` counts.**
`test-validate-definitions.sh:21-25` copies only `.claude/`,
`agent-team-workspace/`, and `CLAUDE.md`, then creates empty stubs for
whatever the validator reports as dangling. A check counting
`rust/tokio_examples/*.rs` would see 0 in the sandbox against a claimed
12 and fail the suite's own baseline case. **inferred** from the
copy list at `:21` and the stub loop at `:22-25`; not executed.

**F7 — Design docs must stay out of the scan set, and this document
proves why.** The file you are reading quotes `` `maxTurns: 20` ``,
`` `model: opus` ``, "twenty turns", and "fifty turns". If
`agent-team-workspace/design-docs/**` entered the corpus, this design
would itself become a claim source and would fail the check the moment
any agent's cap changed. Check 5's scan list excludes it today
(`validate-definitions.py:106-109`); keep it excluded. **observed**

**F8 — Role-token binding resolves the claims positional rules cannot,
and mis-binds none of the 8.** A file's roster is the set of live agent
names it backticks anywhere; a kebab token of exactly one roster member
identifies that member. Executed over all 8 configuration claims:
**observed**

| claim | rung | subject | claimed vs actual |
|---|---|---|---|
| `scope-problem/SKILL.md:18` `effort=max` | 1 | `design-investigator` | max = max |
| `scope-problem/SKILL.md:18` `maxTurns=twenty` | 1 | `design-investigator` | 20 = 20 |
| `design-review-loop-…:14` `effort=max` | 2 | `design-investigator`, `design-bar-raiser` | max = max, max |
| `design-review-loop-…:15` `model=opus` | 3 | `design-investigator`, `design-bar-raiser` | opus = opus, opus |
| `design-review-loop-…:16` `maxTurns=20` | 3 | `design-investigator`, `design-bar-raiser` | 20 = 20, 20 |
| `design-review-loop-…:19` `maxTurns=20` | — | none | advisory |
| `build-agent/SKILL.md:47` `model=fable` | — | none | advisory |
| `build-agent/SKILL.md:48` `effort=max` | — | none | advisory |

Rung 2 fires on `:14` because that sentence contains "investigator" and
"bar-raiser", each a token of exactly one roster member
(roster: `design-investigator`, `design-bar-raiser`,
`ai-writing-auditor`, `architect-reviewer`). Rung 3 fires on `:15-16`
because "Both" ⇒ 2 and the antecedent sentence resolves exactly 2. `:19`
("That measurement predates the 20-turn cap") resolves nothing and
declines. `build-agent/SKILL.md` backticks no live agent name at all, so
its roster is empty and both its claims decline. Five binds, all correct
against frontmatter; zero mis-binds; three declines. **observed**

**F9 — Both validators import PyYAML, which is not the standard
library.** `validate-definitions.py:15` and `validate-skills.py:1` both
`import yaml`; resolved here to
`/usr/lib/python3/dist-packages/yaml/__init__.py` under Python
3.11.15. The brief's constraint "Python 3, standard library only,
matching both existing validators" is half right — Python 3 with no
`pip install` step, but not stdlib-only. This design adds no
dependency, so the constraint binds either way. **observed**

## High-level design

### Where the check lives, and why

`validate-definitions.py`, as two new numbered sections (14 and 15).

Derived, not assumed: the claim corpus spans agents, skills, protocols,
specs, and `CLAUDE.md`, and a single claim relates a sentence in one
file to frontmatter in another. `validate-skills.py` iterates one skill
at a time (`validate-skills.py:16`) and has no place to hang a claim
made in a protocol about an agent. The brief's own constraint assigns
cross-file checks to `validate-definitions.py`.

The hook needs no change. It already runs the definition validator on
edits to agents, skills, `CLAUDE.md`, and the validator itself
(`validate-definitions.sh:38`, `:48`) — which is precisely the set of
files that can introduce or invalidate a claim. The ~13 ms
early-exit path for unrelated calls is untouched, because this design
adds nothing to the shell script.

### The pipeline

Four stages, each with one responsibility.

```
files ──▶ [1 normalize] ──▶ paragraphs ──▶ [2 trigger] ──▶ candidates
                                                              │
                              ┌───────────────────────────────┤
                              ▼                               ▼
                        [3 bind subject]                 (no value)
                              │                               │
             ┌────────────────┴────────────────┐              ▼
             ▼                                 ▼          unparsed
     name / role noun / anaphor           declines        (advisory)
             │                                 │
             ▼                                 ▼
       [4 reconcile]                       advisory
             │                          "no subject bound"
             ▼
   claimed == frontmatter?
             │
             └──▶ fail if not
```

**Stage 1 — normalize.** Split each file's body into paragraphs on
blank lines, join each paragraph's lines with single spaces, and record
the file line each paragraph starts on. Frontmatter is excluded from
the body (`fm_body` already returns the split,
`validate-definitions.py:35-41`) and the `description` field is scanned
separately as a pseudo-paragraph, since round claims live there
(`design-bar-raiser.md` description, `run-pr-loop/SKILL.md`
description). Required by F3.

**Stage 2 — trigger.** Two families, both anchored on a property noun.
Family A is the literal frontmatter echo — a backticked ``key: value``
pair whose key is a real frontmatter field. Family B is the English
restatement — a bound word followed by a value followed by the property
noun. Required by F2: without the noun anchor, `multi-agent-coordinator.md:47`
is a false positive.

**Stage 3 — bind subject.** A decision ladder, most specific first,
with an explicit "decline to bind" terminal. Required by F4: every
proximity- and position-based version fails lines the brief requires to
pass, so the ladder must resolve a subject from the words in the
sentence rather than from what is nearby.

**Stage 4 — reconcile.** Exact comparison against the bound agent's
frontmatter. A claim that did not bind is not compared at all; it is
reported (stage 5). Round claims take their own path: comparison against
the loop protocol the claim's file resolves to.

### Data ownership

| datum | owner | read by |
|---|---|---|
| an agent's `model`/`effort`/`maxTurns` | that agent's frontmatter | stage 4 |
| the set of live agent names | the `.claude/agents/` glob at runtime | stage 3 |
| a file's roster (which agents it talks about) | that file's own backticked text | stage 3, rungs 2-3 |
| a loop's round budget | that loop's protocol file | round check |
| the scan set | check 5's `scan` list | stage 1 |

No new file, no new state, no cache. The check is a pure function of
the tree.

### Failure domains and degradation

- **A claim the extractor cannot parse** degrades to an advisory line,
  not a failure. Derivation: this validator runs inside a PostToolUse
  hook that shows exit-2 output to Claude; a hard failure on an odd
  sentence would block work unrelated to the sentence. This is the
  brief's own default for open question 2, and it is also what goal 4
  asks for — reported, not skipped.
- **A subject that cannot be bound** degrades from checked to reported.
  The claim appears in the advisory list with the reason binding failed
  and is never compared. This loses detection, never soundness: the
  design would rather miss a stale sentence than fail a true one, which
  is what goal 3 makes the hard goal and goal 1 does not.
- **A malformed frontmatter block** is already caught upstream by
  checks 1 and 3b, which run first; stage 4 reads the same parsed
  dicts and inherits their guarantees.
- **The whole section raising** would take checks after it with it. It
  must not index into dicts without `.get`, and must not assume any
  agent has any field — 9 of 13 agents have no `effort` key.

### Operational surface

No deploy, no migration, no rollback. Adding checks changes the printed
total (`231/231`), which is stated in `CLAUDE.md:37` and must be
updated in the same commit — the hook forces a validator run on
`CLAUDE.md` edits (`validate-definitions.sh:38`), so the new number is
displayed to whoever makes the change. Observability is the existing
two-list summary: hard failures under `FAIL`, unparsed claims under
`warn` (`validate-definitions.py:283-286`).

## Low-level design

### Numeral parsing (open question 1: both digits and words)

The live defect was the word "fifty" (`git log -S'fifty turns'`,
commit `a9471ca` introduced it, `ea59d23` fixed it). Digits-only would
have missed the only real instance this repo has produced.

```python
_UNITS = {w: i for i, w in enumerate(
    "zero one two three four five six seven eight nine ten eleven twelve "
    "thirteen fourteen fifteen sixteen seventeen eighteen nineteen".split())}
_TENS  = {w: (i + 2) * 10 for i, w in enumerate(
    "twenty thirty forty fifty sixty seventy eighty ninety".split())}

def _numeral(tok: str) -> int | None:
    """'50' | 'fifty' | 'fifty-five' -> int; anything else -> None."""
```

O(1) per token, 28 table entries covering 0–99. This is a hardcoded
list, and the constraint against hardcoded lists does not reach it: the
banned kind names repo entities that go stale on rename (an agent
allowlist, a five-agent regex). English numerals name nothing in this
repo and do not change when an agent is renamed. State the distinction
in the code comment, because it is the first thing a reader will
challenge.

### Trigger patterns

```python
_FIELDS = ('model', 'effort', 'maxTurns')
_NUM    = r'\d{1,4}|' + '|'.join(_TENS) + '|' + <tens-unit compounds> + '|' + '|'.join(_UNITS)
_BOUND  = r'(?:up to|capped at|cap(?:ped)? of|at most|a maximum of|maximum of|limit of|for)'

ECHO    = re.compile(rf'`({"|".join(_FIELDS)})\s*:\s*([A-Za-z0-9_-]+)`')
TURNS   = re.compile(rf'{_BOUND}\s+({_NUM})[- ]turns?\b', re.I)
TURNS2  = re.compile(rf'\b({_NUM})[- ]turn (?:cap|budget|limit)\b', re.I)
EFFORT  = re.compile(r'\b(?:at|runs? at|of)\s+(max|maximum|xhigh|high|medium|low)\s+effort\b', re.I)
UNPARSED_TURNS  = re.compile(rf'{_BOUND}\s+((?:[a-z-]+\s+){{1,3}})turns?\b', re.I)
UNPARSED_EFFORT = re.compile(r'\b(?:at|runs? at|of)\s+((?:[a-z-]+\s+){1,3})effort\b', re.I)
```

`TURNS2` exists because `design-review-loop-agent-team-prompt.md:19`
says "the 20-turn cap" — noun-first, no bound word. The `UNPARSED_*`
patterns run only on spans no strong pattern consumed, and only fire
when `_numeral` returns `None` for every token in the captured span.
That ordering is what stops "Your turns are capped" — present in all 13
agent bodies — from producing 13 advisories: it has no bound word
*before* the noun, so neither pattern matches.

`effort` levels are read from the union of values actually present in
agent frontmatter plus the platform's documented set. Reading them from
the tree alone would mean a level no agent currently uses is
unrecognized; hardcoding them alone is the banned pattern. Use the
union and let an unrecognized level fall to `UNPARSED_EFFORT`.

### Claim record

```python
Claim = collections.namedtuple('Claim', 'path line field claimed subjects raw')
# path      str   file the claim appears in
# line      int   1-based line in the whole file, frontmatter included
# field     str   'model' | 'effort' | 'maxTurns' | 'rounds'
# claimed   str   normalized value: str(int) for maxTurns/rounds, lowercase token otherwise
# subjects  tuple agent names, empty when unbound
# raw       str   the trimmed sentence, for the failure message
```

The line number is the paragraph's start line plus the newline count in
the paragraph text before the match offset. Getting this wrong by one
is the easiest bug here: for a file with frontmatter, the body's first
line is `frontmatter_block.count('\n') + 1`, and the prototype's
`+ 1` on top of that reported `scope-problem/SKILL.md:19` for a claim
that lives on line 18. Protocols have no frontmatter and were
unaffected, which is exactly how such a bug survives a spot check.

### Subject binding

```python
def _bind(para: str, offset: int, names: frozenset[str], self_agent: str | None
          ) -> tuple[tuple[str, ...], str]:
    """-> (subjects, reason). Empty subjects means unbound; reason names why."""
```

The ladder, in order:

1. **Same sentence.** Split the paragraph on sentence boundaries; if
   the claim's sentence contains exactly one backticked token in
   `names`, bind to it. Covers `scope-problem/SKILL.md:18`.
2. **Same paragraph, unique.** If the paragraph contains exactly one
   name from `names`, bind to it.
3. **Self-reference.** If the file is an agent definition and the
   sentence's subject is second person (`\b(you|your)\b` before the
   property noun), bind to that agent. No current claim needs this; it
   is here because agent bodies are the natural place for one to appear.
4. **Cardinality-checked anaphor.** If the sentence opens with a plural
   anaphor — `both`/`either`/`neither` ⇒ 2, `all three` ⇒ 3 — collect
   distinct names from `names` mentioned in the enclosing section
   (since the last `#`-heading line) *before* the claim, take the last
   *k* by position, and bind only if exactly *k* are available. For
   `design-review-loop-agent-team-prompt.md:15`, "Both" ⇒ 2, and the
   two nearest preceding names are `design-bar-raiser` (`:5`) and
   `design-investigator` (`:4`) — the correct pair. The section names a
   third agent, `ai-writing-auditor` (`:6`), whose `model: sonnet`
   would make a bind-to-all rule fail (F4); recency plus cardinality
   excludes it.
5. **Decline.** Return empty subjects with a reason string.

Rung 4 is the only heuristic in the design, and it is the one a
reviewer should attack. It is included because without it goal 1 is
only partly met: `design-review-loop-agent-team-prompt.md:15-16` is a
sentence stating an agent's `model` and `maxTurns`, and under
existence-checking alone its drift would not fail. It is bounded by
cardinality so that it declines rather than guesses whenever the count
does not match, and its one live instance is verified above. If a
reviewer judges the risk unacceptable, deleting rung 4 costs exactly
one site's precision and nothing else — the ladder degrades to rung 5.

### Reconciliation

```python
for c in claims:
    if c.subjects:                                          # 4a exact
        for s in c.subjects:
            actual = AGENT_FM[s].get(c.field)
            ck(actual is not None and _norm(actual) == c.claimed,
               "prose contradicts agent frontmatter",
               f"{c.path}:{c.line} claims {c.field}={c.claimed} for '{s}', "
               f"frontmatter says {actual!r} | {c.raw}")
    else:                                                   # 4b existence
        pool = {_norm(fm[c.field]) for fm in AGENT_FM.values() if c.field in fm}
        ck(c.claimed in pool,
           "prose states a config value no agent has",
           f"{c.path}:{c.line} claims {c.field}={c.claimed}; "
           f"live values are {sorted(pool)} | {c.raw}")
```

The failure message carries file, line, claimed value, and actual
value — goal 1's stated check verbatim. `ck` is used, not `warn`, and
its return value is ignored here because nothing gates on it; the
return still matters and must not be removed (`validate-definitions.py:19-22`).

4b's weakness is real and must be stated where a reader will see it:
with `maxTurns: 20` on four agents today, an unbound claim of "20"
passes even if the agent the sentence meant moved to 30. 4b catches
only a value that has left the tree entirely — for example, if
`model: fable` were retired, `build-agent/SKILL.md:47` would fail. That
is a genuine class (it is how a renamed model tier would surface) and a
narrow one.

### Round-count agreement (open question 3: in scope, symmetric)

Round budgets have no frontmatter field, so there is no source of
truth. The brief's default anoints the protocol file as the source.
Rejected: it requires a hardcoded path-to-class mapping, which is the
failure mode the constraints name, and it buys nothing — under either
rule, editing the protocol from 5 to 6 alone makes the class
inconsistent and fails.

The symmetric rule: partition round claims into classes, require every
claim in a class to state the same number.

```python
def _loop_class(path: str, rosters: dict[str, set[str]],
                cites: dict[str, set[str]]) -> str | None:
    """Protocol -> itself. Other file -> the single protocol whose backticked
    roster names it or whose path it cites. 0 or 2+ matches -> None."""
```

`rosters` reuses the backticked-kebab parse check 11 already performs
(`validate-definitions.py:231`); `cites` is a substring test for each
protocol path. Zero hardcoded names, zero hardcoded values, and it
survives a rename because both inputs are computed from the tree.
Verified on today's tree (F5): six claim-bearing files each map to
exactly one class, and the two ambiguous files carry no claim.

Trigger, measured against the round-noise this repo contains:

```python
ROUND = re.compile(rf'(?:up to|by|of)\s+\**({_NUM})\**\s*\**rounds?\b'
                   rf'|\b({_NUM})\s+rounds? is the budget'
                   rf'|\bN of (\d+)\b'
                   rf'|\bby round\s+\**({_NUM})', re.I)
```

The second alternative carries a capture group the prototype omitted,
which is why `coding-bar-raiser.md:127` and `design-bar-raiser.md:160`
extracted `None` in the probe run. Verified non-firing on the ordinal
uses in the same files: "converging in two rounds because round 1 was
thorough" (`design-bar-raiser.md:160-161`), "Round 1 casts the widest
net" (`:164`), "From round 2 on" (`:165`), "`R2-3` = round 2, objection
3" (`design-investigator.md:240`), "at round 0" (`:224`), "rounds 1 and
2" (`:253`), and `scope-problem/SKILL.md:240-242`, whose "Round 1" and
"Round 2 and beyond" are *interview* rounds — a different concept in a
file that maps to no loop class.

Known miss: "five rounds of drift is not rigor"
(`design-bar-raiser.md:161`) states the budget value with no budget
marker and is not extracted. Under-detection, not a false positive.

Failure message names both the class and every member claim, because a
disagreement has no single culprit:

```
round budget disagrees within the design-review loop:
  design-bar-raiser.md:160 says 5
  run-design-loop/SKILL.md:description says 6
```

### Unparsed reporting (open question 2: advisory)

```python
warn(False, "config claim not parsed",
     f"{path}:{line} {reason} | {sentence}")
```

`reason` is one of `no numeral in bound phrase`, `unknown effort
level`, `anaphor cardinality k, n candidates`, `no subject resolved and
field absent from every agent`. The last is the only one that also
fails; the rest are advisory only.

Today's tree produces zero of these: all 8 configuration claims either
bind (rungs 1 and 4) or pass existence, and the 11 round claims all
resolve. Preserving the current `0 advisories` baseline is deliberate —
an advisory list that is routinely non-empty is one nobody reads.

### Complexity and cost

Let *F* = 33 files, *L* ≈ 5,000 total lines, *A* = 13 agents.
Normalization is O(*L*); triggering is O(*L*) with a fixed number of
compiled patterns; binding is O(paragraph length × *A*) worst case,
bounded by the section scan in rung 4; reconciliation is O(claims × *A*)
= O(19 × 13). Space is O(*L*) for the paragraph list, dropped per file.

Files are read once into memory and shared with check 5, which
currently re-opens each file (`validate-definitions.py:112`); folding
them into one read is a small net saving that partly offsets the new
pass. The added wall time is **unmeasured**. The measurement that
settles it — five runs of the validator before and after, comparing
medians against the 150–159 ms baseline — is step 7 of the plan.

## What is deliberately out of scope, and why

**Countable inventories in `CLAUDE.md`.** The brief's Problem section
names them; no goal binds them. Restricting them to what the Non-goals
allow leaves almost nothing:

- "165 problems", "24 Rust examples", "12 + 12 example files" reconcile
  against `spaced_repetition/` and `rust/` — excluded by the Non-goal
  on Go, Python, and Rust source, and independently unimplementable in
  the hook's sandbox, which stubs those trees empty (F6).
- "13 agents" and "8 skills" **are not stated in `CLAUDE.md`**. A grep
  for both digit and word forms adjacent to "agent"/"skill" returns
  nothing. There is no claim to check. See Brief corrections.
- "231 checks" (`CLAUDE.md:37`) counts the validator's own runtime
  counter and "24 cases" (`:73`) counts `^run ` lines in the hook
  suite. Both are countable, and both are counts of the checking
  machinery rather than of agent or skill definitions. Checking "231"
  is also self-referential: the check increments the counter it
  validates, so its correctness depends on being evaluated last. A
  non-self-referential form exists — compare after all checks have run,
  appending to `fails` without incrementing `checks` — and the plan
  carries it as an optional final step so a reviewer can take it or
  leave it. The argument against is that the true number is printed on
  every run, to the same person the hook has just forced a run on.

**`agent-team-workspace/design-docs/**` and `research/**`.** Frozen
artifacts, and including them would make this document a claim source
(F7).

**Prose that describes configuration without stating a value.** "Your
turns are capped" asserts nothing checkable and is left alone. This is
the one-directional rule the brief's Non-goals set: a claim that exists
must be true; no claim need exist.

## Alternatives rejected

- **Scan every numeral near a configuration keyword.** Simplest thing
  that could work. Killed by F2: hundreds of hits, dominated by
  ordinary English and numbered list items, and one genuine trap at
  `multi-agent-coordinator.md:47`. Goal 3 is a hard goal.
- **Generate the prose from frontmatter (templating or an include
  directive).** Eliminates drift by construction rather than detecting
  it. Killed by the file format: these are markdown files loaded
  verbatim as agent and skill instructions, with no build step and no
  include mechanism; introducing one means every reader and every
  loader now needs the expander.
- **Bind anaphors to every agent named in the enclosing section.**
  Killed by F4 — it fails
  `design-review-loop-agent-team-prompt.md:15` today.
- **Anoint the protocol file as the round-count source.** The brief's
  default. Killed by the hardcoded-mapping constraint; symmetric
  agreement detects the same drift with nothing hardcoded.
- **Put the check in `validate-skills.py`.** Killed by structure: it
  iterates one skill at a time and has no place for a protocol's claim
  about an agent.
- **Fail on unparsed claims.** Killed by the hook: exit 2 on an odd
  sentence blocks unrelated work, and this is a lint.

## Plan

Ordered by risk. Step 1 is the spike: if the trigger set is noisier in
practice than the prototype measured, the taxonomy is wrong and steps
2–8 do not survive.

**1. Spike — extraction only, no reconciliation.**
Files: `agent-team-workspace/validate-definitions.py`.
Add stages 1 and 2 plus a temporary dump of every candidate with its
file, line, field, and value. Depends on nothing.
Verify: the dump lists exactly the 8 configuration claims and 12 round
claims (11 from the probe plus `pr-loop-agent-team-prompt.md:5-6`,
which paragraph joining must now catch, F3), and nothing else. Line
numbers match: `scope-problem/SKILL.md:18`, `build-agent/SKILL.md:47`
and `:48`, `design-review-loop-agent-team-prompt.md:14`, `:15`, `:16`,
`:19`. Any extra candidate stops the plan.

**2. Subject binding, rungs 1–3 and 5.**
Files: `agent-team-workspace/validate-definitions.py`.
Depends on 1.
Verify: `scope-problem/SKILL.md:18` binds to `design-investigator`;
the four protocol claims and the two `build-agent` claims report
unbound with a reason.

**3. Reconciliation 4a and 4b, wired as section 14.**
Files: `agent-team-workspace/validate-definitions.py`.
Depends on 2.
Verify: exit 0 on HEAD with the advisory list still empty; then edit
`scope-problem/SKILL.md:18` to "fifty turns" in a scratch copy and
confirm exit 1 with a message naming the file, line 18, claimed 50, and
actual 20. Revert.

**4. Rung 4, the cardinality-checked anaphor.**
Files: `agent-team-workspace/validate-definitions.py`.
Depends on 3. Separated from step 2 so it can be dropped without
unpicking the ladder.
Verify: `design-review-loop-agent-team-prompt.md:15-16` binds to
`{design-investigator, design-bar-raiser}` and passes; a scratch edit
of `design-bar-raiser.md`'s `maxTurns` to 30 makes it fail naming both
agents.

**5. Round-count classes, wired as section 15.**
Files: `agent-team-workspace/validate-definitions.py`.
Depends on 1 (not on 2–4; it uses its own subject rule).
Verify: exit 0 on HEAD; a scratch edit of `run-pr-loop/SKILL.md`'s
description from 4 to 5 rounds fails, listing every pr-class member and
its number; the design class is unaffected.

**6. Unparsed advisory tier.**
Files: `agent-team-workspace/validate-definitions.py`.
Depends on 3.
Verify: HEAD still prints `0 advisories`; a scratch edit to "up to a
couple dozen turns" prints one advisory naming the file, line, and
`no numeral in bound phrase`, and exit stays 0.

**7. Hook test cases and the cost measurement.**
Files: `.claude/hooks/test-validate-definitions.sh`.
Depends on 3, 5, 6. Add five `run` cases under `=== UPDATE ===`, using
the existing `run <desc> <exit> <payload> [setup]` contract
(`test-validate-definitions.sh:32-41`): the "fifty turns" regression
(exit 2), a `model` drift in a bound claim (exit 2), a round-count
drift in one class member (exit 2), an unparseable phrasing (exit 0,
goal 4), and a benign prose edit near a claim (exit 0).
Verify: the suite reports 29 passed, 0 failed. Separately, time five
validator runs and compare the median against the 150–159 ms baseline;
record the delta in the commit message.

**8. Update `CLAUDE.md`.**
Files: `CLAUDE.md`.
Depends on 7. Change the `231` at `CLAUDE.md:37` to the new total, the
`24 cases` at `:73` to 29, and extend the validator's one-line
description to name the prose-reconciliation checks.
Verify: `python3 agent-team-workspace/validate-definitions.py` exits 0
and its printed total equals the number now written at `:37`.

**9. Optional — check the two self-counts.**
Files: `agent-team-workspace/validate-definitions.py`, `CLAUDE.md`.
Depends on 8. Compare `CLAUDE.md`'s stated check count against the
final `checks` value and its stated case count against
`grep -c '^run '`, appending to `fails` without incrementing `checks`.
Verify: exit 0; then decrement the number in a scratch copy and confirm
exit 1. Drop this step if a reviewer judges the added edit-coupling not
worth it.

### What would invalidate this plan

- **Step 1's dump is noisier than 20 candidates.** The taxonomy is then
  wrong, not the implementation. Verified by measurement in step 1.
- **The measured wall-time delta is large enough to matter.**
  **Unverified** — no measurement taken. Mitigation is sharing check
  5's file reads; if that is not enough, the pass can be gated behind a
  cheap substring pre-filter.
- **Rung 4's recency binding is judged too clever.** The design
  survives without it at the cost of one site (step 4 is separable).
- **A future claim shape appears that none of the six patterns match.**
  It is silently unchecked, not falsely failed. This is the
  one-directional design the Non-goals chose; the advisory tier only
  catches near-misses of the existing patterns.
- **Goal 2's stated verification measures the wrong check.** Renaming
  an agent without sweeping references already fails checks 11 and 11b
  (`validate-definitions.py:229-258`). The rename test must sweep
  references, or it reports pre-existing failures as this check's.

## Open questions

1. **Should the countable inventories be in scope at all?** Owner:
   lead. This design scopes them out (derivation above), which leaves
   the brief's Problem statement partly unaddressed. Settled by the
   lead saying whether "231 checks" and "24 cases" drifting is a cost
   worth an edit-coupling; step 9 implements it if so.
2. **What is the added wall time?** Owner: implementer. Settled by
   step 7's measurement — five runs before and after, median compared
   to the 150–159 ms baseline measured here.
3. **Does the `effort` level vocabulary need a source beyond the
   tree?** Owner: implementer. Only `max` and `xhigh` are in use, so a
   claim naming a valid-but-unused level falls to the advisory tier.
   Settled by checking the platform's documented `effort` values
   against what `build-agent/SKILL.md:42` says about verifying field
   values — not read for this design.
4. **Is `validate-skills.py:62`'s `IndexError` (F8) worth fixing here?**
   Owner: lead. Out of this brief's scope; it will crash the skills
   validator the first time a ninth skill is added.

## Brief corrections

1. **"The eight countable claims in `CLAUDE.md` (231 checks, 24 hook
   cases, 13 agents, 8 skills, ...)"** — `CLAUDE.md` does not state "13
   agents" or "8 skills" in any digit or word form. The countables it
   does state are 231 (`:37`), 24 cases (`:73`), 165 problems (`:82`),
   24 / 12 / 12 examples (`:95`), and 23 GoF patterns (`:65`), which
   the brief does not list. Ground truth confirmed: 13 agent files, 8
   skill directories, 231/231 checks, 24 `^run ` cases, 165 entries in
   `problems.json`, 12 files in each Rust example directory — every
   stated number is correct, but two of the eight are not stated.
2. **"Python 3, standard library only, matching both existing
   validators"** — both validators import PyYAML
   (`validate-definitions.py:15`, `validate-skills.py:1`), which is not
   the standard library. The operative constraint is "no new
   dependency", which this design meets.
3. **Constraint "`.claude/hooks/validate-definitions.sh` ... costs ~13
   ms on unrelated calls; whatever this design adds runs inside that
   budget"** — the ~13 ms is the hook's early-exit path for unrelated
   calls, which this design does not touch. The budget it must fit is
   the validator's own run, measured here at 150–159 ms; the hook's own
   comment cites 174 ms (`validate-definitions.sh:45`).
4. **Goal 2's check, "rename an agent and re-run; the check ... does
   not report the old name as missing"** — an unswept rename already
   fails checks 11 and 11b today, so the check as written will see
   failures that are not this design's. The rename must be swept for
   the test to measure what it intends.

## Objection responses

**R1-1 | accepted.** I executed the specified ladder over the 33-file scan
set and reproduced the failure exactly: rung 2 ("same paragraph, exactly one
name") binds all four claims in
`design-review-loop-agent-team-prompt.md`'s lines 9-24 paragraph to
`architect-reviewer`, the only agent name in it (`:23`), against whose
frontmatter (`model: inherit`, no `effort`, `maxTurns: 15`) every one of them
fails. Rung 2 is deleted, not bounded: measured over all 8 configuration
claims it binds nothing correctly and produces four false positives, so no
requirement holds it up. The binding decision for all 8 claims under the
replacement ladder is tabulated in F8.

**R1-2 | accepted, with a different remedy than the one suggested.** The
positional rule and its stated result do contradict each other: computed on
the file, "the last 2 by position" is
`{design-bar-raiser, ai-writing-auditor}`, and binding `model: opus` against
`ai-writing-auditor` (`model: sonnet`) fails. I then tested whether any
positional rule yields `{design-investigator, design-bar-raiser}` without
reading prose: last-*k* and nearest-*k* both give the wrong pair, and first-*k*
gives the right one only by coincidence on a three-name section. So the
positional rung is dropped. It is replaced by role-token binding (new ladder
rung 2) plus a cardinality-checked anaphor over the antecedent sentence (rung
3), which is structural rather than positional — it reads the file's own
backticked roster at runtime and resolves the kebab-case tokens that identify
exactly one of its members. Measured: it binds 5 of the 8 configuration claims,
all correctly, and declines on the other 3.

**R1-3 | accepted.** The headline now states the corpus measured by the
specified patterns — 8 configuration claims and 12 round-budget claims across
33 files — and cites the run that produced it. The prototype's three defects
are why the old number was wrong; F1 now reports the specified patterns' output
and keeps the prototype only as the thing whose defects the specification
fixes.

**R1-4 | accepted.** No goal forces 4b, and its only live subjects are
`build-agent/SKILL.md:47-48`, which I read: a parenthetical about a past
decision ("this repo's precedent: a request for an \"endgame\" ability became
`model: fable` + `effort: max` after verification"). That sentence stays true
after `fable` is retired and 4b would then fail it. 4b is deleted; an unbound
claim now reports as advisory and nothing else. The consequence — the advisory
baseline is 3, not 0 — is stated in the headline and in Unparsed reporting.

**R1-5 | accepted.** A round claim whose class is `None` is now an advisory
with reason `no loop class`, listed among the four advisory reasons. Taking the
brief's default under R1-6 makes the class-of-one problem disappear rather than
requiring a rule: a claim is compared against its protocol's value, not against
its siblings, so a class with one non-protocol member is still checked, and a
protocol carrying two different round numbers is a hard failure against itself.

**R1-6 | accepted.** The kill reason was wrong: `_loop_class` computes the
path-to-class mapping from rosters and path citations with nothing hardcoded,
so the anointed variant costs one line on top of it. The design takes the
brief's default. Two things it buys that symmetric agreement does not: culprit
attribution in the failure message, and a live check for a class whose only
other member is the protocol.

**R1-7 | accepted.** The amended goal-2 check is now plan step 5 (a swept
rename executed against a scratch copy) and hook case 6 in step 8.

**R1-8 | accepted.** Re-verified at HEAD `bfc9cec`: `validate-skills.py` is 65
lines, `OLD2NEW` has 6 entries and no identity entries, and the check F8 cited
is gone. F8 and open question 4 are deleted; the system map's line count is
corrected. F8's slot now carries the role-token binding measurement, which is
the evidence this round needed.

**R1-9 | accepted.** The coupling is stated where step 9 is proposed and again
in the operational surface: the new sections call `ck` once per checked claim,
so `CLAUDE.md:37`'s total becomes a function of prose content and step 9 would
fail on any prose edit that adds or removes a claim. Step 9 stays optional and
now carries that cost explicitly.

## Sources

All evidence is from the repository at commit `92a0545` (HEAD), with
the tree state matching `8832723` for every file cited except
`brief.md`. Commands run:

- `python3 agent-team-workspace/validate-definitions.py` — 231/231, 0
  advisories, exit 0.
- `python3 agent-team-workspace/validate-skills.py` — 8 skills, 0
  failures, exit 0.
- Timing loop over four validator runs — 150, 154, 156, 159 ms.
- `git log -S'fifty turns' -- .claude/skills/scope-problem/SKILL.md` —
  introduced `a9471ca`, fixed `ea59d23`.
- Prototype extractor run against the scan set — 8 configuration
  triggers, 11 round triggers, 0 spurious.
- `ls .claude/agents/*.md | wc -l` → 13; `ls -d .claude/skills/*/ | wc
  -l` → 8; `ls rust/tokio_examples/*.rs | wc -l` → 12;
  `ls rust/basics_examples/*.rs | wc -l` → 12;
  `len(json.load(open('spaced_repetition/problems.json')))` → 165;
  `grep -c '^run ' .claude/hooks/test-validate-definitions.sh` → 24.

No external sources were consulted; every claim in this document is
grounded in the repository.
