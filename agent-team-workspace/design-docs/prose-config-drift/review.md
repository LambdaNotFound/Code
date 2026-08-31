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
