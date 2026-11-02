# Brief: advice-eval: an eval framework for the advice written into the fi repo

Status: draft
Slug: advice-eval
Started: 2026-11-02

> Privacy: the Code repo is public. GitHub listed LambdaNotFound/Code as
> `public` on 2026-11-02 (observed). This brief names no holdings,
> amounts or account details; those stay in the private fi repo.

## Problem

Draft. The user hasn't corrected it yet.

The user makes real trades on conclusions Claude writes into the private
fi repo, and some of those trades have tax consequences. The
conclusions come in several kinds:

- stock and option reports;
- weekly trade plans;
- portfolio and tax reviews;
- since 2026-11-02, short-term predictions with check dates.

Nothing checks those conclusions systematically:

- **Observed:** the computation scripts have unit tests:
  `.claude/skills/analysis-stock/scripts/test_ta.py` (75 tests) and
  `.claude/skills/analysis-option/scripts/test_opt.py` (23). Nothing
  checks the reports written from their output.
- **Observed:** this session, four kinds of error reached a written
  report or a piece of advice, and each was caught by hand:
  1. The script misread a quirk in the data feed: a week that spanned a
     split produced a fake all-time high. Fixed in commit 81c2125.
  2. A number in a report didn't match its source. A re-read caught it.
  3. An API reply was misread: an empty earnings calendar was reported
     as "no report in the next 3 months".
  4. A decision rule wasn't supported by the numbers: a sell trigger
     sat inside normal earnings-week noise, as the outlook computed
     later showed. Reversed the same day.
- **Observed:** prediction scoring is a manual to-do item in the fi
  repo. No outcome is measured today.
- **Inferred:** without checks, the user can't tell how much to trust a
  stance, a trade plan or a probability, and errors reach decisions
  that cost tax.

Gap: which decision the eval must support. Trust before acting?
Improving the skill? An audit after the fact? See open question 1.

## Goals

Placeholders taken from the request. None is confirmed, and each still
needs its check.

1. **Faithfulness** (assumed): every number a report states traces to a
   source, within a tolerance: the data files, the script's JSON, or
   the positions file. Check: to be decided.
2. **Rule compliance** (assumed): reports meet the fi repo's written
   rules: required sections, tags, a chart reason for every trade
   suggestion, and the account and tax rules. Check: to be decided.
3. **Outcomes** (assumed): predictions and trade suggestions are scored
   against what happened after their check dates, with calibration
   tracked over time. Check: to be decided.
4. **Regression** (inferred): a change to the skills is caught if it
   moves the stance or the plan on frozen inputs. Check: to be decided.

## Non-goals

Proposed, not confirmed:

- Judging fundamentals, valuation or news.
- Executing trades or placing orders.
- Proving the chart frame has an edge. A scoreboard can measure an
  edge; it can't promise one.

## Constraints and invariants

- **Observed:** the Code repo is public. The fi repo is private and
  holds positions and income (`fi/CLAUDE.md`: "The GitHub repo is
  private ... and must stay private").
  - So nothing private may land in the Code repo: no positions,
    amounts, trade plans, or eval results built from them.
- **Observed:** the fi repo limits where Claude writes. "Write only
  under `reports/` and `data/`, or in `scripts/` when I ask for a
  script" (`fi/CLAUDE.md:26`). Old reports are never edited; a new one
  is written instead.
- **Observed:** fi rules a checker could test:
  - one recommendation, with its reason, its main risk and a rules
    check (`fi/CLAUDE.md:48`);
  - a chart reason for every trade suggestion (`fi/CLAUDE.md:51`);
  - five required report sections (`fi/CLAUDE.md:119`);
  - taxes (`fi/CLAUDE.md:162`);
  - investment strategy (`fi/CLAUDE.md:199`);
  - "My rules" (`fi/profile.md:36`).
- **Observed:** the Alpha Vantage free tier allows 25 calls a day
  (`.claude/skills/analysis-stock/references/alphavantage-free-tier.md`).
  Scoring outcomes needs bars from after each check date. The fi data
  rule caches every pull (`fi/CLAUDE.md:88`).
- **Observed:** scripts are Python, standard library only (the fi
  scripts rule, and the skills' scripts).
- **Observed:** `ta.py` is deterministic. The outlook uses a fixed seed
  (`.claude/skills/analysis-stock/scripts/ta.py:984`), so frozen inputs
  give identical output.

## Decomposition

Not yet: that's phase 4.

## Open questions

The frame questions. Owner for all five: the user.

1. What should this change? Default: show which numbers, rules and
   calls in a report hold up, both before acting and after.
2. Who is it for? Default: the user alone, with Claude as the producer
   that gets corrected.
3. What does done look like? Default: a command that grades one report
   and lists every mismatch and rule break, plus a scoreboard of
   predictions and suggestions against outcomes.
4. What must or must not be used (deadline, tools, budget, where the
   code lives)? Default: standard-library Python, generic checkers in
   the public Code repo, private runs and results in fi.
5. Anything else before reading the code.

## Context

- `.claude/skills/analysis-stock/scripts/ta.py`:
  - `analyse()` at `:450`, `ledger()` at `:815`, `stance()` at `:915`;
  - `outlook()` at `:1166`, `render()` at `:1251`, `main()` at `:1514`.
  - The `--json` output is the machine-readable source a faithfulness
    check would compare a report against.
- `.claude/skills/analysis-stock/references/report-template.md`: the
  section contract. Short-term predictions start at `:78`.
- `.claude/skills/analysis-stock/references/indicators.md`: the tags
  and confidence caps a rule checker would enforce.
- `.claude/skills/analysis-option/scripts/opt.py`: the GOOD, FAIR and
  POOR option verdicts.
- The fi repo, private, at `/home/user/fi`:
  - `CLAUDE.md`: the rules;
  - `profile.md:36`: "My rules";
  - the redacted positions file under `accounts/`;
  - `data/`: cached Alpha Vantage CSVs;
  - `reports/stocks/`: 18 reports;
  - `reports/reviews/`: 13 files, including `to-do.md`.

## Questions asked

None yet.
