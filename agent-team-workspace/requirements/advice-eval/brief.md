# Brief: advice-eval: an eval framework for the advice written into the fi repo

Status: signed-off 2026-11-02
Slug: advice-eval
Started: 2026-11-02

> Privacy: the Code repo is public. GitHub listed LambdaNotFound/Code as
> `public` on 2026-11-02 (observed). This brief names no holdings,
> amounts or account details; those stay in the private fi repo.

## Problem

Restated in phase 1, then refined by the user's answers in rounds 1–7.

The user makes real trades on conclusions Claude writes into the private
fi repo, and some of those trades have tax consequences. The
conclusions come in several kinds:

- stock and option reports;
- weekly trade plans;
- portfolio and tax reviews;
- since 2026-11-02, short-term predictions with check dates.

Nothing checks those conclusions systematically:

- **Observed:** the computation scripts have unit tests:
  `fi/.claude/skills/analyse-stock/scripts/test_ta.py` (75 tests) and
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

**Why doing nothing loses (inferred):**

- About 70 predictions from this week's reports alone would go
  unscored, and each week adds more. How far to trust a call would stay
  a guess.
- Errors like the four above keep reaching decisions that cost tax,
  and each one is caught only if someone happens to re-read the report.
- The null option costs nothing up front. That's the only thing it has
  going for it.

The eval must support all four uses (user, round 1): catching errors
before acting, scoring the calls afterwards, protecting the analysis
code, and auditing the reports already written.

- **Observed:** the reports written on 2026-11-01 and 2026-11-02 use at
  least 8 different trade-table layouts and 2 different prediction
  layouts.
  - So a scorer that reads the Markdown would need a parser for each
    layout, and would break on the next one.
  - The file names are not cited here, because they would reveal
    holdings in a public repo.

## Goals

Confirmed by the user in rounds 1–7. Each goal carries its check.
Revised after the reader test (see "Reader test").

1. **Faithfulness:** every number a report states traces to a source.
   - Sources: the data files, the script's JSON, the positions file,
     any other fi file (accounts, profile, monthly costs, watchlist), a
     dated web page, or a derivation from other claims (its inputs and
     formula).
   - Tolerance (user, round 5): the printed value must sit within half
     a unit of its last printed digit of the source value. So 65.25
     matches 65.2500, and 1.98% matches 1.9812%. Either rounding mode
     passes.
   - Each claim records its unit, because the script's JSON stores some
     percentages as percent points and others as fractions
     (`fi/.claude/skills/analyse-stock/scripts/ta.py:1316`).
   - The printed value must also appear verbatim in the report's text.
     That catches an extraction that copied the source value instead
     of what the report says.
   - The grader also lists every number in the report's text that the
     block leaves out, as "unlisted".
   - A number with no source is reported as "untraceable", not as wrong.
   - Check: one command grades a report and lists every untraceable or
     mismatched number.
2. **Rule compliance:** reports meet the fi repo's written rules:
   required sections, tags, a chart reason for every trade suggestion,
   and the account and tax rules.
   - Each report is graded against the rules in force when it was
     committed, read from fi's git history. A rule added later is listed
     as "not yet required", not as a failure.
   - The required-sections rule applies to stock reports only.
   - Check: the same command lists every rule a report breaks.
3. **Outcomes:** predictions and trade suggestions are scored after
   their check dates against what happened. Predictions answer two
   different questions, scored separately:
   - **Are the odds honest?** Ranges, level touches, closes beyond a
     level and order-fill odds come from a model with no view on
     direction. They carry no edge by construction, so they're scored
     for calibration: how often the event happened against the stated
     probability, by probability bucket.
     - "Miscalibrated" is called when the interval for (observed minus
       stated) excludes zero.
     - "Honest within 5 points" is called only when the whole interval
       sits inside plus or minus 5 points. That's an equivalence margin,
       not a failed test.
     - Otherwise the scorecard says "too few to tell".
   - **Is there skill?** Direction calls are scored against each call's
     no-edge probability: what the driftless model gives the same event.
     That's 50% only when the level is the as-of close. Two scores:
     - the Brier skill score against those no-edge probabilities;
     - direction accuracy (the predicted side against 50%), but only for
       calls whose level is the as-of close.
     - A skill verdict is called only when the interval excludes zero.
   - Otherwise, hit rates are never compared with a coin flip: a
     calibrated 68% range hits 68% of the time with no skill at all.
   - **Verdicts** (user, round 4): every score carries a 95% interval.
     No verdict is called before at least 8 distinct check dates have
     been scored.
   - **Correlation:** calls on the same date across tickers move
     together, and a ticker's 1-, 2- and 4-week calls share one price
     path. So intervals come from resampling whole as-of dates, with a
     fixed seed. Intervals that assume independent calls would be too
     narrow.
   - **Edge test** (user, round 6): a live test with a power target.
     The scorecard reports each edge with its interval, and how many
     more scored calls a given edge would need. There is no historical
     replay, and no premium data is bought.
   - **Trades:** each suggestion's result over 1, 2 and 4 weeks. The
     actual result, from the imported fills, is compared with no trade,
     with SPY over the same dates, and with the rule's own entry, stop
     and target simulated from the bars.
   - **Taxes** (user, round 2): every trade gets a pre-tax result.
     Taxable-account trades also get an after-tax one, using the rates
     in `fi/profile.md` and the tax rules in `fi/CLAUDE.md`, with basis
     and holding period from the lots in the trade-history import.
   - About 70 predictions from this week's reports alone come due on
     11-06, 11-13 and 11-27, and each week's reports add more. The
     scorecard is the first deliverable (user, round 4); piece 1 is its
     prerequisite.
   - Check: a scorecard report in fi.
4. **Regression:** a change to the skills is caught if it moves the
   stance, the levels, the plan or the outlook on frozen inputs. Check:
   automated tests that re-run the analysis on frozen data.
5. **Audit:** the checks in goals 1 and 2 run over every report already
   in fi: 30 reports, since `to-do.md` is not a report (`fi/CLAUDE.md:31`).
   - Method (user, round 3): Claude reads each old report once and
     writes its claims, trades and predictions to a JSON file under
     `fi/reports/evals/`. The grader runs on that file.
   - Old reports stay unedited.
   - Check: an audit report listing every finding.
6. **Machine-readable calls** (user, round 2): every new report carries
   a fenced JSON block, placed just before its closing `[RULES I BROKE]`
   line. It lists the report's numeric claims, trade suggestions and
   predictions, with the inputs the report was written from. Check: the
   grader reads the block, and a report without one fails.
7. **Judgment rules** (user, round 3): scripts check what's exact;
   a Claude grader scores the judgment rules.
   - Exact checks: numbers, sections, tags, and triggers inside the
     noise.
   - Judgment rules, for example "the chart reason cites a level with
     evidence" or "counterarguments come first", go to the grader. It
     scores them against a written rubric and quotes its evidence.
   - Only exact findings may trigger an automatic trade revision.
     Grader findings are flagged for the user.
   - Keeping the grader honest (user, round 4): the user labels about
     20 findings once. The grader's scores count only once they agree
     with those labels, and the check re-runs whenever the rubric
     changes.
8. **Actual fills** (user, round 3): the broker's trade-history CSVs
   are imported by a redacting script, like the positions import.
   - Fills give each trade's actual result. The mechanical baseline is
     simulated from bars, not from fills.
   - The import is expected to give each lot's basis and acquisition
     date, which the after-tax results need. This is assumed: no export
     has been seen yet.
   - This also closes the open "get the trade history" to-do item.
9. **Fundamentals and news** (user, round 6). Three checks:
   - **Cited facts:** when advice cites a P/E, an analyst target, an
     earnings date or a news item, it's checked against a dated source.
     Anything stale, unsourced or wrong is flagged.
   - **Ignored events:** advice is flagged when it ignores a known event
     inside its horizon: earnings, an ex-dividend date, a split, or an
     index change.
   - **Fundamental sense:** a Claude grader judges whether the advice
     holds up on valuation and growth. Its findings are flagged only,
     under goal 7's rules.
   - **Observed:** the free key returns the earnings calendar but not
     company fundamentals. `COMPANY_OVERVIEW` answered "premium
     endpoint" on 2026-11-02.
   - **Inferred:** most of these sources are therefore dated web pages,
     read in the session. Each must be cited with its date.

## Non-goals

Confirmed by the user in round 5:

- **Placing trades.** The eval never places or changes orders at the
  broker; the user executes every trade.
- **Automatic triggers.** No hook before commits, no schedule, no CI on
  push. The checks run on demand only.

Fundamentals and news, and the edge test, were proposed as non-goals
but brought into scope in round 5. Round 6 defined them: see goals 9
and 3.

Out of scope as a consequence of those answers:

- a historical replay of the stance rule;
- buying premium data.

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
  (`fi/.claude/skills/analyse-stock/references/alphavantage-free-tier.md`).
  Scoring outcomes needs bars from after each check date. The fi data
  rule caches every pull (`fi/CLAUDE.md:88`).
  - Caching began on 2026-11-02, with fi commit 45dcbee. That commit
    holds the 2026-10-30 close that every report so far was written
    from (observed).
- **Inferred:** the Alpha Vantage MCP tools exist only inside a Claude
  session, and `ta.py --fetch` needs a key the user hasn't set
  (`fi/CLAUDE.md:110-114`).
  - So the eval scripts work offline, from `data/`.
  - Before a run, Claude fetches any missing bars under the cache rule
    and records the calls it spent in the run report.
- **Observed:** fi's scripts are Python 3, "standard library first"
  (`fi/CLAUDE.md:61`). The skills' scripts use the standard library
  only.
- **Observed:** `ta.py` is deterministic. The outlook uses a fixed seed
  (`fi/.claude/skills/analyse-stock/scripts/ta.py:984`), so frozen inputs
  give identical output.
- **User, round 1:** the results go to the user, and Claude fixes what
  fails without asking each time.
  - So each finding must be specific enough to act on: the file, the
    claim, the expected value, the value found, and the source.
  - A fix to a report is a new correction report, because old reports
    are never edited (`fi/CLAUDE.md:28`).
  - A fix to the analysis skills goes to the Code repo with tests and
    no private data. A fix to the eval itself goes to fi.
- **User, round 2:** Claude may do four things without asking:
  - write correction reports;
  - fix the code, with tests;
  - revise trade advice;
  - spend Alpha Vantage calls on scoring.
- **User, round 4:** automatic revisions may move prices and triggers.
  A revision that adds a taxable sale, or flips a buy to a sell, is
  proposed to the user instead of made.
- **Inferred:** a revision is never silent.
  - Claude writes it as a new dated report that names the old advice,
    the new advice and the finding behind it.
  - Every run ends by listing its revisions for the user.
  - Claude never places orders. The user still executes every trade.
- **Inferred:** spending calls stays inside the free tier's 25 a day
  and the fi data cache rule. When the budget runs out, scoring is
  marked pending; it never estimates the missing bars.
- **Imposed, round 3:** all the eval code lives in the private fi repo,
  under `fi/scripts/`. The user chose this over a public engine in the
  Code repo.
  - Consequences:
    - golden regression data can be real and private;
    - the eval runs the Code repo's `ta.py` by path;
    - a change to `ta.py` is checked only when the eval is run on
      demand.
  - **Observed:** fi's `.gitignore` skips many folder names at any
    depth, such as `lib/`, `build/` and `env/`, so each new path must
    pass `git check-ignore -v` (`fi/CLAUDE.md`, scripts rule).
- **Observed:** a raw broker export holds account numbers and must
  never be committed (`fi/CLAUDE.md:45`).
  - The to-do list notes that `.gitignore` catches only
    the positions export, not the trade-history export.
  - So the history importer must redact before anything is staged.
    The ignore pattern needs the user's go-ahead, because `.gitignore`
    is the user's file.
- **Observed:** both build loops assume they run in this repo and
  commit their state to it.
  - The PR loop creates branch `pr/<slug>` and commits `pr.md` and
    `review.md` under `agent-team-workspace/pull-requests/<slug>/` as
    `[loop]` commits, pushed every round
    (`agent-team-workspace/protocols/pr-loop-agent-team-prompt.md:30-51`,
    `:64`, `:113`).
  - The design loop commits and pushes `design.md` and `review.md`
    under `agent-team-workspace/design-docs/<slug>/`
    (`agent-team-workspace/protocols/design-review-loop-agent-team-prompt.md:46-49`,
    `:136-139`).
  - **Inferred:** in a public repo, design notes and review rounds about
    the private reports would be published. The final commit that
    removes the loop state doesn't erase it from pushed history.
  - **Inferred:** the code itself belongs in fi (round 3), which the
    loops don't target.

## Decomposition

**Route, for every piece (imposed, user, round 7):** built directly in
the main session, under `fi/scripts/`, with tests. Each piece is
committed to fi under its rules (`YYYY-MM-DD: summary` subjects, pushed
right away) and reviewed by the user.

- No agent loop is used, so no piece brief is copied to a loop
  directory. This brief is the single source.
- Each piece's paragraph below is written to be built from without the
  scoping conversation.

Order: dependency first, then the deadline, then risk.

| # | Piece (slug) | Goals | Route | Depends on |
|---|---|---|---|---|
| 1 | `advice-eval-calls-schema` | 6, 3 | direct in fi | none |
| 2 | `advice-eval-scorecard` | 3 | direct in fi | 1 |
| 3 | `advice-eval-grader` | 1, 2, 7 | direct in fi | 1 |
| 4 | `advice-eval-regression` | 4 | direct in fi | none |
| 5 | `advice-eval-audit` | 5 | direct in fi | 1, 3 |
| 6 | `advice-eval-fundamentals` | 9 | direct in fi | 3 |
| 7 | `advice-eval-fills` | 8 | direct in fi | the user's history download |
| 8 | `advice-eval-trades` | 3 | direct in fi | 1, 2, 7 |

### 1. `advice-eval-calls-schema`

Goal: a machine-readable record of every report's calls.

- **The block:** a fenced JSON block placed just before a report's
  closing `[RULES I BROKE]` line. The exact field names are piece 1's
  output, approved by the user before piece 2 starts. At a minimum it
  records:
  - **the report:** its path, its date, the close its data ends at (the
    as-of date), and its inputs: the fi commit holding the `data/` files
    used, the `ta.py` commit, and the flags passed.
    - The inputs are required for new reports.
    - A calls file extracted from an old report marks itself as legacy.
      It may give "unknown" for any input it can't recover, with the
      reason;
  - **each claim:** an id, the printed text exactly as it appears, the
    value, the unit of the printed value (price, percent points,
    fraction, count or date), its source's unit when that differs, its
    tag, and its source. A source is one of:
    - a file, with a locator (a field, a JSON path, a row key or a
      line);
    - the script's JSON, with a JSON path;
    - a web page, with its URL and date;
    - a derivation, with the ids of its inputs and the formula;
  - **each trade:** an id, the ticker, the account type (taxable, Roth,
    HSA, 401k), side, quantity, order type, prices, trigger, fallback,
    horizon, and the stance it rests on (label and score). Also either
    its chart reason, or "no technical reason" with the real reason;
  - **each prediction:** an id, the ticker, its kind, the level or range
    and the direction (above or below), and the probability.
    - Also the no-edge probability: what the driftless model gives the
      same event. That's 50% for a close above the as-of close, and the
      stated probability itself for a model-derived odd.
    - Also the as-of close, the check date, and which table it came
      from.
- **Prediction kinds:**
  - close above a level on a date;
  - close inside a range on a date;
  - trades at a level by a date;
  - a daily close beyond a level by a date;
  - a Friday (weekly) close beyond a level by a date.
- **Validator:** a script that checks a block against the schema.
- **Rule:** a rule in `fi/CLAUDE.md` making the block required, edited
  at the user's request.
- **First extraction:** the two prediction reports dated 2026-11-02
  under `fi/reports/stocks/` (one for a single ticker, one for 13
  tickers).
  - Every probability with a check date is extracted, from the numbered
    predictions and from the odds tables, each tagged with its table.
  - The calls files go under `fi/reports/evals/`, so piece 2 can score
    them on 11-06.
- **Done when:** both files validate; a deliberately broken block fails
  with a message naming the field; and the user approves the schema.

### 2. `advice-eval-scorecard`

Goal: prediction outcomes, with honest uncertainty.

- **Inputs:** the calls files, plus daily bars covering each check date,
  read from `data/`.
  - The script runs offline.
  - Before a run, Claude fetches any missing bars under the cache rule.
    If it can't, the prediction is marked pending.
- **Windows:**
  - date kinds ("close above", "close inside a range") use the close on
    the check date;
  - path kinds use the sessions from the first one after the as-of close
    through the check date:
    - "trades at" uses every session's high and low;
    - "daily close beyond" uses every session's close;
    - "Friday close beyond" uses only each week's last-session close;
  - a check date that isn't a trading day falls back to the last session
    before it.
  - The printed date governs, not a count of sessions. `ta.py` counts
    weekdays and ignores holidays
    (`fi/.claude/skills/analyse-stock/scripts/ta.py:990`).
- **Splits:** the daily bars are unadjusted. When the weekly adjusted
  series shows a split inside a window, the level is scaled by the split
  ratio before scoring, and the report says so.
- **Known bias:** the outlook simulated 20 weekdays for the 2026-11-27
  calls, but only 19 sessions trade because of Thanksgiving. The
  scorecard notes this and doesn't correct for it.
- **Output:**
  - a calibration table of observed frequency against stated
    probability, by bucket, for the model-derived odds;
  - for direction calls:
    - the Brier score, and the Brier skill score against each call's
      no-edge probability;
    - direction accuracy, where the level is the as-of close;
  - 95% intervals from resampling whole as-of dates, with a fixed seed;
  - the verdict rules from goal 3;
  - how many more independent calls would confirm the skill the calls
    claim, computed from their own probabilities (open question 2);
  - the calls Claude spent fetching bars, passed in as an option.
- **Report:** a dated report under `fi/reports/evals/`.
- **Done when:** a test fixture with known outcomes scores to
  hand-computed values, and the 11-06 predictions are scored that day.

### 3. `advice-eval-grader`

Goal: one command that grades one report.

- **Exact checks:**
  - every claim matches its source at the printed precision, and a
    claim with no source is reported as "untraceable";
  - the required sections are present (`fi/CLAUDE.md:119`);
  - every claim's tag, read from the block, is one of `[KNOWN]`,
    `[COMPUTED]`, `[INFERRED]`, `[COMMON]`, `[FRAME]` or `[GUESS]` (the
    user's tag set), and the same tag also appears in the report text on
    the claim's line;
  - every trade has a chart reason or "no technical reason"
    (`fi/CLAUDE.md:51`);
  - every recommendation has its reason, main risk and rules check
    (`fi/CLAUDE.md:48`);
  - any trade price (limit, trigger, fallback) or flip level with
    better-than-even odds of a daily close beyond it within 10
    sessions is flagged in the report, as the rule requires
    (`fi/CLAUDE.md:147-149`).
    - The grader computes those odds itself with `ta.py`'s simulation,
      for any price, not only `ta.py`'s own levels.
    - It fails a report only when such a level goes unflagged.
  - the account rules hold (`fi/profile.md:36`), for example covered
    calls only and the leveraged-ETF cap.
- **Claude grader:** a written rubric for the judgment rules, applied
  with quoted evidence.
  - Its scores count only after agreeing with the user's labels:
    matching the user's pass or fail on at least 16 of 20 (report,
    rule) items.
  - The rubric lives at `fi/scripts/eval/rubric.md`, and the labels at
    `fi/reports/evals/labels.json` (defaults).
- **Rules in force:** each report is graded against the rules in force
  when it was committed (goal 2).
  - The fi rule says "five sections" over four bullets, because
    support and resistance count as two (`fi/CLAUDE.md:119`).
  - Piece 3 proposes clearer wording to the user, and doesn't edit the
    rule without asking.
- **Findings:** each one names the file, the claim, the expected value,
  the value found, and the source.
- **Done when:**
  - seeded broken reports (a changed number, a missing section, a
    missing chart reason, an unflagged trigger inside the noise) each
    produce
    exactly their finding;
  - a clean report produces none;
  - the grader meets the agreement bar on the labels.

### 4. `advice-eval-regression`

Goal: catch changes to the analysis code that move its outputs.

- **How:** run `ta.py`, from the Code repo by path, on `fi/data/`
  pinned at a git ref, and diff the JSON against expected files stored
  in fi.
  - The diff covers the stance, the score, the levels, the plan and the
    outlook.
  - Expected changes are accepted by re-pinning, with a commit message
    saying why.
  - The first pin is the fi commit whose `data/` holds the 2026-10-30
    close for all 14 cached tickers.
  - Expected outputs go under `fi/scripts/eval/golden/`, once
    `git check-ignore` passes.
- **Done when:** the run passes today, and fails after a deliberate
  one-line change to the stance bands.

### 5. `advice-eval-audit`

Goal: grade every existing report.

- **How:** Claude reads each of the 30 reports once and writes its calls
  file under `fi/reports/evals/`. Then the grader runs on all of them,
  against the rules in force when each was committed, and the result is
  an audit report.
- **Fixes:** made under the auto-fix rules:
  - a correction report for each wrong number or misread source;
  - automatic revisions of prices and triggers, from exact findings
    only;
  - new taxable sales and buy-to-sell flips are proposed, not made.
- **Done when:** every report has a validated calls file and a graded
  result, and every finding is closed, proposed or listed.

### 6. `advice-eval-fundamentals`

Goal: the fundamentals and news checks in goal 9.

- **Cited facts:** checked against dated sources, read in the session.
  The source and its date are recorded.
- **Ignored events:** an event inside a call's horizon (earnings, an
  ex-dividend date, a split, an index change) that the report doesn't
  mention is flagged.
- **Fundamental sense:** graded by the Claude grader and flagged only.
- **Done when:** a seeded report citing a stale P/E and ignoring an
  earnings date inside its horizon produces both findings.

### 7. `advice-eval-fills`

Goal: actual fills, without exposing account numbers.

- **How:** a redacting importer for the broker's trade-history CSVs,
  like the existing positions importer under `fi/scripts/`. The raw
  file is never staged.
- The approved `.gitignore` pattern for the raw history files is
  added.
- **Done when:** a sample export imports with no account number in the
  output, and `git check-ignore` matches the raw file.

### 8. `advice-eval-trades`

Goal: trade results against the three baselines.

- **Baselines:** each trade's actual result over 1, 2 and 4 weeks, from
  the imported fills. It's compared with no trade, with SPY over the
  same dates, and with the rule's own entry, stop and target simulated
  from the bars.
- **Taxes:** pre-tax for every trade. Taxable-account trades also get an
  after-tax result from three sources:
  - the rates in `fi/profile.md`;
  - the tax rules in `fi/CLAUDE.md`;
  - each lot's basis and holding period, from the import.
  - When a holding period is unknown, both the short-term and the
    long-term results are shown, labeled.
  - When a lot is older than the imported history, the positions file's
    average basis is used, labeled.
- **Edge test:** the same intervals, verdict rule and power target as
  piece 2.
- **Done when:** a fixture with known fills and bars gives hand-
  computed results.

## Open questions

Rounds 1–8 are answered (see the log). At sign-off the user accepted
the defaults for questions 1, 2 and 5. Questions 3 and 4 are Claude's
to settle while building, under their defaults.

1. **What agreement with the user's labels makes the grader count?**
   Answered at sign-off: at least 16 of 20 (80%), re-checked when the
   rubric changes.
2. **Which skill test drives the power target?** Answered at sign-off:
   report both. At
   80% power and 95% confidence:
   - **Direction accuracy against 50%:** for a true 55% hit rate, about
     620–780 independent calls (616 one-sided, 782 two-sided).
   - **Brier skill score:** for calls stated at 55% that are exactly
     right, about 2,450 independent calls. A 55% forecast beats 50% in
     squared error only when the true rate is above 52.5%.
   - Same-date calls aren't independent, so either one means well over
     a year of weekly reports.
   - Default: report both, each computed from the calls' own stated
     probabilities.
3. **What if a source is unavailable when a claim is checked?** Owner:
   Claude, in piece 3. Default: mark the claim "unverified" and list it,
   but don't count it as a failure.
4. **Where do outputs go, and what shape is the tool?** Owner: Claude.
   Default:
   - one CLI, `fi/scripts/advice_eval.py`, with subcommands: validate,
     score, trades, grade, regress and audit;
   - Python 3, as fi's scripts already use;
   - tested with `unittest` in `fi/scripts/test_advice_eval.py`;
   - one dated report per run under `fi/reports/evals/`, listing its
     findings, its revisions, what's pending and the calls it spent;
   - every new path passes `git check-ignore` first.
5. **What calibration margin, and how many check dates before a
   verdict?** Answered at sign-off: plus or minus 5 points, and at least
   8 distinct check dates.

Settled with less questioning, because the blast radius is small:

- **Scale (inferred):** one user, about 14 tickers, about 70
  predictions from each week's reports so far, and 30 old reports, all
  run on demand. Standard-library Python is enough.
- **Blast radius:** the new rule in `fi/CLAUDE.md` changes every future
  report. `ta.py` is unchanged. The eval reads old reports but never
  edits them.
- **Deadline and budget:** the only deadline is piece 2 by 11-06. The
  budget is the free tier's calls, with no paid services (round 6).
  The Claude grader runs inside the session.

## Reader test (phase 5)

Run on 2026-11-02. A design-investigator agent read the brief with no
conversation context.

- **What it got right:**
  - the shape of pieces 1 and 2: a standard-library schema and
    validator, and an offline scorer with Brier scores, intervals and a
    dated report;
  - the privacy split;
  - the build route.
- **What it missed, now fixed in the brief:**
  - Goal 3 scored range and touch odds against a coin flip, which a
    calibrated model beats with no skill. Calibration and skill are now
    scored separately.
  - The intervals assumed independent calls. They now resample whole
    check dates.
  - The schema had no field list, no weekly-close kind, no as-of date
    or window start, no tag or unit fields, and no derived or non-data
    sources.
  - Holidays: the printed check date now governs, not a count of
    sessions.
  - The noise check would have failed reports that follow the rule. It
    now fails only unflagged levels, and covers any order price.
  - Old reports are graded against the rules in force when they were
    committed.
  - Several things were undefined and are now stated: fills versus the
    mechanical baseline, the tax-rate source, holding periods, the CLI
    shape, and where the rubric, labels and golden files live.
  - Small corrections:
    - 30 reports, not 31;
    - "standard library first";
    - caching began on 2026-11-02;
    - about 70 predictions due across three dates, not 70 a week;
    - the block sits just before `[RULES I BROKE]`.
- **Second round** (same reader, revised sections only):
  - It confirmed most fixes: 6 of 10 questions and 12 of 17 defects.
  - It found seven new problems, all fixed:
    - The "honest" verdict was backwards. It now needs an equivalence
      margin.
    - The Friday-close kind read every session; it now reads only each
      week's last.
    - Old reports can't supply their inputs. They're now marked legacy,
      with "unknown" allowed.
    - The power target used the wrong metric. Both metrics are now
      reported.
    - The no-edge reference is 50% only at the as-of close; it's now
      defined for any level.
    - Two stale counts were corrected.
    - The missing trades subcommand was added.
  - Also fixed:
    - numbers left out of the block are listed;
    - splits inside a window are handled;
    - the 19-versus-20-session bias is noted;
    - tags must appear in the text;
    - calls spent are passed in;
    - lots older than the history use the average basis.
- **Still open:** open questions 2 and 5.

## Context

- `fi/.claude/skills/analyse-stock/scripts/ta.py`:
  - `analyse()` at `:450`, `ledger()` at `:815`, `stance()` at `:915`;
  - `outlook()` at `:1166`, `render()` at `:1251`, `main()` at `:1514`.
  - The `--json` output is the machine-readable source a faithfulness
    check would compare a report against.
- `fi/.claude/skills/analyse-stock/references/report-template.md`: the
  section contract. Short-term predictions start at `:78`.
- `fi/.claude/skills/analyse-stock/references/indicators.md`: the tags
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

Round 1, 2026-11-02, interactive:

1. **What should it do?** All four: catch errors before acting, score
   the calls afterwards, protect the analysis code, and audit existing
   reports.
2. **Who reads and acts on the results?** The user, plus Claude fixes
   what fails without asking each time.
3. **What does done mean?** All four: a report-grading command, a
   prediction scorecard, trade results against a baseline, and
   regression tests on the same data.
4. **When does it run?** On demand only.

Round 2, 2026-11-02, interactive:

1. **Machine-readable calls?** A JSON block inside each report.
2. **Baselines for trade results?** All three: no trade, SPY over the
   same window, and the rule applied mechanically.
3. **Taxes in trade results?** Show both pre-tax and after-tax.
4. **Fixes Claude may make without asking?** All four: correction
   reports, code fixes with tests, revised trade advice, and spending
   data calls.

Round 3, 2026-11-02, interactive:

1. **Where does the code live?** Everything in fi.
2. **How are judgment rules graded?** Scripts plus a Claude grader.
3. **Where do actual fills come from?** A broker history import.
4. **How are old reports audited?** A one-time extraction.

Round 4, 2026-11-02, interactive:

1. **When is a verdict called?** Only when the 95% range clears the
   baseline.
2. **Limits on automatic revisions?** New taxable sales and buy-to-sell
   flips are flagged, not made.
3. **First piece?** The prediction scorecard.
4. **How is the grader kept honest?** It's checked against about 20 of
   the user's labels.

Round 5, 2026-11-02, interactive:

1. **What stays out of scope?** Only placing trades and automatic
   triggers. Fundamentals and news, and proving an edge, are now in
   scope.
2. **Number tolerance?** Match at the printed precision.
3. **Anything else?** Nothing else.

Round 6, 2026-11-02, interactive:

1. **What to do with fundamentals and news?** All three: check cited
   facts, flag ignored events, and grade fundamental sense.
2. **How to test for an edge?** A live test with a power target.

Round 7, 2026-11-02, interactive:

1. **How do the pieces get built?** Directly in fi, in the main
   session, with no agent loops.
2. **May the trade-history pattern go into fi's `.gitignore`?** Yes.

Round 8, sign-off, 2026-11-02, interactive:

1. **Sign off?** Yes.
2. **Skill test?** Report both.
3. **Thresholds?** Plus or minus 5 points, and 8 check dates.
4. **Grader bar?** 16 of 20.
