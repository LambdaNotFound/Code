# Review ledger — prose-config-drift

## Round 1

**Independent derivation** (from `brief.md` + its 2026-08-31 amendment and the
tree, before reading `design.md`):

1. Invariants: (a) every in-scope prose claim reconciles against the file that
   *owns* the value, and disagreement calls `ck(False)` with file, line,
   claimed, actual; (b) subject binding reads the agent roster from the
   `.claude/agents/` glob at runtime — never a name or value list — so a swept
   rename rebinds and the old name has no subject; (c) zero new failures on
   today's tree, including the two correct claims at
   `.claude/skills/scope-problem/SKILL.md:18` and
   `agent-team-workspace/protocols/design-review-loop-agent-team-prompt.md:14-21`;
   (d) an unparseable candidate is emitted, never dropped — so the trigger set
   must be strictly wider than the parse set.
2. Hard constraints: the two existing validators only; PyYAML already imported,
   so "no new dependency" not "stdlib only" (amendment 2); the budget is the
   validator's own run, which I measured at 131/149/156/160/152 ms, mean 150
   (amendment 3); `ck` keeps returning its condition; new cases land in the
   24-case hook suite.
3. Minimum moving parts: an owner map (agent name → frontmatter dict, already
   built at `validate-definitions.py:122-125`); one scanner over check 5's
   existing `scan` list that finds a property noun, binds a subject in a bounded
   window, and extracts a numeral (digits + a small word table, brief Q1
   default); a three-way outcome — matched / mismatched (fail) / unparsed
   (warn). No new file, no cache, no allowlist.
4. Cost ceiling I measured before reading the design: 8 configuration claims and
   12 round-budget claims in the scan set, and the scan set is 33 files. The one
   genuinely hard case is the anaphoric subject at
   `design-review-loop-agent-team-prompt.md:15` ("Both run `model: opus`"),
   which the brief also requires to keep passing.

The proposal reaches broadly the same place: same location, same corpus, same
one-directional rule, and its trigger design is better than my sketch — I
reimplemented the six specified patterns and ran them over the scan set, getting
exactly 8 configuration and 12 round candidates with zero spurious hits, which
is the goal-3 evidence the document itself does not yet have. The diff that
produces objections is entirely in stage 3, subject binding: the ladder as
written does not do what the document says it does, and the document's
verification of it is wrong on the file.

### Objections

R1-1 | blocking | Ladder rung 2 ("same paragraph, exactly one name → bind", design.md:369-370) binds all three `ECHO` claims in `design-review-loop-agent-team-prompt.md`'s lines 9-24 paragraph to `architect-reviewer` — the only agent name in that 16-line paragraph, at `:23` — producing three hard failures on today's tree (`effort: max` vs `None`, `model: opus` vs `inherit`, `maxTurns: 20` vs `15`). Rung 2 precedes rung 4, so the anaphor rung that F4 says rescues this line never runs. I executed the specified ladder to confirm. This is a false positive on the exact line the brief's Context names as must-keep-passing, so goal 3 fails and F4's conclusion is unreachable. | A paragraph-scoped bind needs a bound — proximity, sentence distance, or a cardinality guard like rung 4's — and the fix must be evidenced by running the whole ladder over the 33-file scan set and showing the binding decision for all 8 configuration claims, not narrated for one.
R1-2 | blocking | Rung 4's stated rule and its stated result contradict each other. The rule is "collect distinct names in the enclosing section before the claim, take the last *k* by position" (design.md:377-381); the names before `:15` in the section opened by the `#` heading at `:1` are `design-investigator` (`:4`), `design-bar-raiser` (`:5`), `ai-writing-auditor` (`:6`), so the last 2 by position are `design-bar-raiser` and `ai-writing-auditor` — not the pair F4 claims. Binding `model: opus` against `ai-writing-auditor` (`model: sonnet`, `maxTurns: 30`) fails, which is the very outcome F4 says "recency plus cardinality excludes". The design kills the bind-to-all alternative (design.md:567-569) on the strength of this line, so the kill rests on an unsound replacement. | Restate the rung so its output on this file is the pair the design claims, and show the computed subject set for `:15` from the rule as written, not asserted. If no positional rule yields `{design-investigator, design-bar-raiser}` without reading the prose cue at `:9-10`, rung 4 should be dropped (design.md:395-396 says it is separable) and the site accepted as unbound.
R1-3 | should-fix | The headline's safety argument — "the measured claim corpus is 8 configuration claims and 11 round-budget claims across 30 scanned files, and a trigger ... fires on exactly those 19 and on nothing else (probe output)" (design.md:26-30) — attributes a measurement to a prototype the document elsewhere reports as defective in three ways: it missed `pr-loop-agent-team-prompt.md:5-6` (F3), it reported `scope-problem/SKILL.md:19` for a line-18 claim (design.md:352-354), and it dropped a capture group so two round claims extracted `None` (design.md:465-467). The corpus is 20, not 19; the scan set is 33 files, not 30 (`validate-definitions.py` check 5 prints "33 files scanned", and the design's own system map and complexity section both say 33). | Restate the headline as the corpus measured by the *specified* patterns and cite that run. For what it is worth, I ran the six specified patterns over the 33-file scan set and got 8 configuration and 12 round candidates with zero spurious hits and only the two intended `UNPARSED_*` overlaps at `scope-problem/SKILL.md:18` — the substance holds; the attribution does not.
R1-4 | should-fix | Stage 4b, the existence check (design.md:409-414), has no forcing requirement: goal 1 asks that a claim disagreeing with *that agent's* frontmatter fail, and 4b cannot detect that; no goal asks for "a value no agent currently has". Its only live subjects are `build-agent/SKILL.md:47-48`, which is a historical anecdote ("a request for an 'endgame' ability became `model: fable` + `effort: max` after verification") — a sentence that stays true after `fable` is retired but that 4b would then fail. So the component adds a false-positive class against a hard goal while satisfying none. | Name the requirement that forces 4b, or kill it in favour of the simpler rule the design already has to hand: an unbound claim reports as advisory and nothing else. If 4b stays, state why a true sentence about a past decision failing the run is acceptable.
R1-5 | should-fix | `_loop_class` returns `None` for 0 or 2+ matches (design.md:443-447), but the design never says what happens to a round claim whose class is `None`, and the four advisory reasons at design.md:496-499 have no entry for it — so such a claim is silently skipped, which is what goal 4 forbids. It is not hypothetical: I computed the mapping and `CLAUDE.md` matches both protocols while `review-pr/SKILL.md` and `scope-problem/SKILL.md` match neither. Relatedly, symmetric agreement is vacuous for a class of size 1, which the design does not address. | Specify the disposition of a classless round claim (advisory with a reason) and the behaviour of a single-member class.
R1-6 | should-fix | The kill reason given for the brief's open-question-3 default — anointing the protocol as the round-count source "requires a hardcoded path-to-class mapping, which is the failure mode the constraints name" (design.md:433-436) — is refuted by the design's own `_loop_class` (design.md:443-452), which computes that mapping from rosters and path citations with nothing hardcoded and already maps "protocol → itself". The anointed variant is one line on top of the design's own function. The second reason, "it buys nothing", is false for a single-member class and for culprit attribution in the failure message. | Either give a kill reason the design's own mechanism does not refute, or take the brief's default. The decision is the design's to make; the stated grounds are not available.
R1-7 | should-fix | Goal 2 as amended ("rename an agent, sweep every reference, and confirm the new reconciliation check still reconciles that agent's claims under its new name and reports no missing subject for the old one") has no verification anywhere: no plan step exercises it and none of step 7's five hook cases is a rename case. The mechanism is sound and I accept it — stage 3 reads names from the `.claude/agents/` glob at runtime (design.md:246) — but the goal whose whole point is this repo's three hardcoded-list defects is the one goal with no test. | Add the amended check as a plan step and as a sixth hook case in `.claude/hooks/test-validate-definitions.sh`.
R1-8 | should-fix | Three citations no longer resolve after commit `1525e8a`, which the amendment announces: the system map calls `validate-skills.py` "76 lines" (it is 65), F8 cites `validate-skills.py:62`'s `IndexError` (the check was deleted), and F8's support "two of them self-mapped (`validate-skills.py:5`)" is stale (`OLD2NEW` lost its identity entries and now has 6). Open question 4 rests on F8. | Re-verify against HEAD and drop or restate F8 and open question 4.
R1-9 | nit | Step 8 treats `CLAUDE.md:37`'s "231" as a one-time edit, but the new sections call `ck` once per claim (`validate-definitions.py:19-21` increments per call), so the printed total becomes a function of prose content and goes stale whenever a claim is added or removed — which makes optional step 9's self-count check fail on ordinary prose edits. The coupling already exists through checks 7, 11 and 11b, so this enlarges it rather than creating it. | State the coupling where step 9 is proposed, or have sections 14 and 15 contribute a count independent of corpus size.

### Spot-checks

- `validate-definitions.py:19-22` (`ck` returns its condition) | held
- `validate-definitions.py:106-109` scan list = 33 files | held (validator prints "33 files scanned"); design.md:27's "30" | failed
- `validate-definitions.py:229-234` is check 11's protocol-token loop | held
- F2 `multi-agent-coordinator.md:47` "Concurrent subagents are capped (20 by default)" | held
- F4 `design-review-loop-agent-team-prompt.md:4-6` names three agents, claim at `:15` | held on the names, failed on the binding conclusion (R1-2)
- F3 `pr-loop-agent-team-prompt.md:5-6` wraps across lines | held
- F6 `test-validate-definitions.sh:21-25` copies only `.claude/`, `agent-team-workspace/`, `CLAUDE.md` and stubs danglers | held
- `test-validate-definitions.sh:32-41` run helper contract | held (`:33-41`)
- `git log -S'fifty turns'` → introduced `a9471ca`, fixed `ea59d23` | held
- `231/231 hard checks passed, 0 advisories`, exit 0 | held
- Timing 150–159 ms | held (I measured 131/149/156/160/152, mean 150)
- `validate-skills.py` "76 lines" | failed (65)
- F8 `validate-skills.py:62` `IndexError` | failed (deleted in `1525e8a`)
- F1's 8 configuration claims and their sites | held (reimplemented and re-ran the specified patterns; also 12, not 11, round claims)

Sample widened after the F4 and `validate-skills.py` citation failures: I
reimplemented all six specified trigger patterns plus ladder rungs 1–2 and ran
them over the full 33-file scan set rather than spot-reading them.

### Brief corrections

None beyond the four already recorded in the amendment, all of which I
independently confirmed (PyYAML at `/usr/lib/python3/dist-packages/yaml/`; the
validator's own run at ~150 ms; `CLAUDE.md` stating neither "13 agents" nor "8
skills"; `validate-definitions.py:229-258` already failing an unswept rename).
One small thing for the record: `.claude/hooks/validate-definitions.sh:45` still
says "worth one 174ms run", which the same measurement makes stale.

Verdict: revise

## Round 3

**Independent derivation (from brief + amendment + codebase, before reading design.md).**
Invariants: (a) every in-scope prose claim resolves to a subject drawn from the
filesystem/frontmatter at runtime — never a name list (brief constraint, three prior
failures); (b) mismatch ⇒ non-zero exit naming file, line, claimed, actual (goal 1);
(c) detection must be strictly wider than parsing, because a single regex that fails to
match is silent and goal 4 forbids silence — so trigger-set minus parsed-set is the
unparsed report, and open Q2's "warn" default keeps that report from breaking goal 3;
(d) the corpus today must produce zero new *failures* (goal 3), which the
`.../design-review-loop-agent-team-prompt.md:14-16` "both run..." anaphor stresses, since
its subject is not in its own sentence.
Constraints: Python 3, no new dependency (amendment 2, not stdlib-only); code lands in the
two existing validators; `ck` must keep returning its condition; new cases in the 24-case
hook suite; runtime inside the 146 ms full-run budget (amendment 3), not the 16 ms early exit.
Minimum parts: reuse the existing scan corpus + `agents` frontmatter map; one subject
resolver (nearest in-scope backticked agent name, with an explicit rule for multi-subject
anaphora); per-attribute trigger regex + strict parse; digit-and-word number normalizer
(Q1 default: both, the live defect was "fifty"); round-count claims reconciled against the
protocol file as sole source (Q3 default); test cases. No new module, no sidecar registry,
no annotation syntax in the prose, no NLP.

**Diff against the proposal.** Same location, same corpus, same one-directional rule,
same three-way outcome. The proposal carries two parts my derivation does not: role-token
binding (rung 2) and the cardinality-checked anaphor (rung 3). Both are forced by the
brief's must-keep-passing site at `design-review-loop-agent-team-prompt.md:14-16`, whose
subject is not in its own sentence, and the design kills my simpler "nearest name"
resolver with a measurement I reproduced in round 1 (F4). Different route, same place —
no objection from the diff. What I did instead was re-measure the round-2 revision's
central evidence, since every load-bearing claim in it is a measurement claim.

**Re-measurement (independent reimplementation of the six specified patterns, the
`ROUND` pattern, `_loop_class`, and the full binding ladder over check 5's 33-file scan
list).** All of the following reproduced exactly, against HEAD `4df7121`:

- 8 configuration candidates, 0 spurious — the sites F1 lists, with F1's line numbers.
- 13 round-budget candidates, 0 spurious — the sites and values F5 tabulates, including
  `design-review-loop-agent-team-prompt.md:108`, whose numeral sits on `:109` across a
  line wrap. My round-1 count of 12 was the one that was wrong; R1-3 closes with the
  author correcting the reviewer.
- F8's binding table, row for row: rung 1 binds `scope-problem/SKILL.md:18` (both
  claims), rung 2 binds `:14`, rung 3 binds `:15` and `:16`, three decline. 5 checked,
  0 mis-binds, 3 advisories — the headline's stated baseline.
- All 13 round claims classify; 3 protocol sources, 10 compared, 0 disagreements.
- 0 unparsed from stages 1-2, with the consumed-span suppression rule.
- Goal 1's stated check: on a scratch copy with "twenty turns" reverted to "fifty
  turns", the ladder yields exactly one mismatch — `scope-problem/SKILL.md:18`,
  `maxTurns` claimed 50, frontmatter 20, subject `design-investigator`.

### Objections

R3-1 | should-fix | Plan step 5 defines the swept rename as "the file, the `name:` field, and every reference the validator's checks 11 and 11b would otherwise flag", then asserts the outcome "rung 2 resolves 'author' against the renamed roster, since the prose noun is swept with everything else … the advisory count stays 3". Checks 11 and 11b match only backticked kebab tokens (`validate-definitions.py:231`, `:250`), so they never flag the bare English role noun in "The investigator and bar-raiser are expert software engineers…" (`design-review-loop-agent-team-prompt.md:9-10`). Measured on scratch copies both ways: renaming `design-investigator` → `design-author` across every token occurrence leaves `:15` and `:16` declining with `anaphor wants 2, antecedent resolves 1` and binds `:14` to `design-bar-raiser` alone — advisory count 5, not 3. Only when the untracked prose noun is *also* rewritten do all three rebind and reconcile. The mechanism is sound and the loss is reported, not silent, so amended goal 2 is met in substance; the plan's stated expected outcome is wrong on its own definition of the sweep. | Restate step 5's expected result for a token-only sweep (5 advisories, `:14` bound to one subject), and name English role nouns as a rename surface no check in this repo tracks — or widen the definition of "sweep" to include them and say what detects an unswept one.
R3-2 | should-fix | Step 2's verification gate — "exit 0 on HEAD with exactly 3 advisories, all reason `no agent name or role noun in sentence`, at `build-agent/SKILL.md:47`, `:48` and `design-review-loop-agent-team-prompt.md:19`" — contradicts the design's own body, which says "Rung 1 alone satisfies goal 1's stated check and leaves 6 of the 8 claims unchecked, including the whole `design-review-loop-agent-team-prompt.md` paragraph". Step 2 ships rung 1 plus decline only, and none of `:14`, `:15`, `:16` contains a backticked live agent name in its sentence — that is precisely why they need rungs 2 and 3 — so all three decline at step 2. The correct gate is 6 advisories at step 2, dropping to 3 at step 3. An implementer who trusts the stated gate will read a correct step-2 build as broken, and the cheapest way out is to pull rungs 2-3 forward, collapsing the separability the design's fallback argument rests on. | Restate step 2's gate as 6 advisories with the three extra sites named, and step 3's as the drop from 6 to 3.
R3-3 | should-fix | "**Every multi-word literal in every pattern uses `\s+`, never a literal space.**" is contradicted by three of the patterns printed in the same document: `_BOUND` (`up to`, `capped at`, `at most`, `a maximum of`, `maximum of`, `limit of`), `TURNS2` (`turn (?:cap|budget|limit)`), and `EFFORT` (`runs? at`). Only `ROUND` and `ECHO` comply. The prose two paragraphs later says the correction "applies to `up to` in `ROUND` and to every phrase in `_BOUND`", but the block a reader copies is the uncorrected one, and this is the exact defect class that hid `design-review-loop-agent-team-prompt.md:108` from two prior measurement rounds. | Print the patterns with `\s+` already applied, or label the block as pre-correction and state which alternatives are corrected where.
R3-4 | nit | `Claim.line` is documented as "`int` 1-based line in the whole file, frontmatter included", but stage 1 scans `description` as a pseudo-paragraph and the design's own round-count failure message prints `.claude/skills/run-design-loop/SKILL.md:description`. 4 of the 13 round claims and 0 of the 8 configuration claims live in a description, so the field is `int | str` in practice. | Type the field for both cases and state what a description-borne claim reports as its locus.

### Spot-checks

- `validate-definitions.py:112` — check 5 re-opens each file per line | held
- `validate-definitions.py:231` — check 11's backticked-kebab parse, reused by `_loop_class` | held
- `validate-definitions.py:283-286` — two-list FAIL/warn summary | held
- `validate-definitions.py:287` — `sys.exit(1 if fails else 0)`, file is 287 lines | held
- `validate-definitions.sh:45` "worth one 174ms run" (brief correction 3's stale-figure note) | held
- `validate-definitions.sh:56-58` — output discarded on exit 0 | held
- `test-validate-definitions.sh:21-25` — copies only `.claude/`, `agent-team-workspace/`, `CLAUDE.md`, then stubs danglers (F6's basis) | held
- `grep -c '^run '` = 24 | held
- `CLAUDE.md:37` states "231 structural…checks"; `:73` states "24 cases" | held
- `build-agent/SKILL.md:42` — "Verify field names, allowed values, and model capabilities" (open question 3's referent) | held
- `design-review-loop-agent-team-prompt.md:108-109` — "by\nround 5" wraps; brief cites `:109`, design cites `:108` and says so | held
- `design-review-loop-agent-team-prompt.md:19` — sentence names no agent and no role token | held
- `build-agent/SKILL.md` backticks no live agent name; roster empty | held
- F1's 8 configuration sites and F5's 13 round sites, reimplemented and re-run | held
- F8's full binding table, all 8 rows | held
- Headline's "advisory baseline is 3, not 0" | held (measured 3)
- `231/231 hard checks passed, 0 advisories`, exit 0 at HEAD `4df7121` | held
- Validator runtime | held (I measured 123/129/129/132/183 ms, median 129, against the design's ~150 and the amendment's 146 — same order, no claim rests on the difference)

No citation failed this round; the sample was not widened.

### Closed

R1-1 resolved · R1-2 resolved · R1-3 resolved (and the author corrected my count: 13
round claims, not 12) · R1-4 resolved · R1-5 resolved · R1-6 resolved · R1-7 resolved
(step 5 and hook case 6 exist; their stated outcome is R3-1) · R1-8 resolved · R1-9
resolved. No regressions.

### Brief corrections

None new. The four in the amendment all reconfirm, and the design records them.

### Residual risks accepted with this approval

1. **Role-noun coverage decays on rename.** Trigger: any agent rename. Symptom: the
   advisory count rises above the enumerated 3, with reason `anaphor wants k, antecedent
   resolves n`. Reported, not silent; costs detection, never soundness.
2. **Plan gates state numbers the implementation will contradict** (R3-1, R3-2).
   Trigger: running steps 2 and 5 as written. Fix is textual and must land before the
   PR loop starts.
3. **The printed patterns under-detect wrapped phrases** (R3-3). Trigger: any future
   claim phrase that wraps a source line. Fix is textual.
4. **Countable inventories stay unchecked**, which leaves the brief's Problem statement
   partly unaddressed. Trigger: the lead answering open question 1 yes; step 8 exists
   for it, with its edit-coupling cost stated.

Verdict: approve-with-risks

