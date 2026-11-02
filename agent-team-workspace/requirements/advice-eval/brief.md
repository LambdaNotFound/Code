# Brief: advice-eval: an eval framework for the advice written into the fi repo

Status: draft
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

**Why doing nothing loses (inferred):**

- About 70 predictions a week would go unscored, so how far to trust
  a call stays a guess.
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

Confirmed by the user in rounds 1–6. Each goal carries its check.

1. **Faithfulness:** every number a report states traces to a source:
   the data files, the script's JSON, or the positions file. Check: one
   command grades a report and lists every untraceable or mismatched
   number.
   - Tolerance (user, round 5): the source value, rounded to the digits
     the report prints, must equal the report's number. So 65.25
     matches 65.2500, and 1.98% matches 1.9812%.
   - A number with no source is reported as "untraceable", not as
     wrong.
2. **Rule compliance:** reports meet the fi repo's written rules:
   required sections, tags, a chart reason for every trade suggestion,
   and the account and tax rules. Check: the same command lists every
   rule a report breaks.
3. **Outcomes:** predictions and trade suggestions are scored after
   their check dates against what happened.
   - Predictions: hit or miss, a Brier score, and calibration over time.
   - Trades: the result over 1–4 weeks against three baselines (user,
     round 2): no trade, SPY over the same dates, and the rule's own
     entry, stop and target applied mechanically.
   - Taxes (user, round 2): show both. Every trade gets a pre-tax
     result; taxable-account trades also get an after-tax one, from the
     lot basis and the plan's tax rate.
   - Verdicts (user, round 4): hit rates and Brier scores carry 95%
     intervals. The scorecard calls a verdict ("better than a coin
     flip", "beats SPY") only when the interval excludes the baseline;
     until then it says "too few to tell".
   - Check: a scorecard report in fi.
   - Built first (user, round 4). About 70 predictions come due on
     11-06, 11-13 and 11-27.
   - Edge test (user, round 6): a live test with a power target. The
     scorecard reports the edge over each baseline with its 95%
     interval, and how many more scored calls a given edge would need
     before it could be detected. There is no historical replay, and no
     premium data is bought.
4. **Regression:** a change to the skills is caught if it moves the
   stance, the levels or the plan on frozen inputs. Check: automated
   tests that re-run the analysis on frozen data.
5. **Audit:** the checks in goals 1 and 2 run over every report already
   in fi. Check: an audit report listing every finding.
   - Method (user, round 3): Claude reads each old report once and
     writes its claims, trades and predictions to a JSON file under
     `fi/reports/evals/`. The grader runs on that file.
   - Old reports stay unedited.
6. **Machine-readable calls** (user, round 2): every new report ends
   with a fenced JSON block. It lists the report's numeric claims (each
   with its source), its trade suggestions and its predictions. Check:
   the grader reads the block, and a report without one fails.
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
   - They feed the "rule applied mechanically" baseline.
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
  (`.claude/skills/analysis-stock/references/alphavantage-free-tier.md`).
  Scoring outcomes needs bars from after each check date. The fi data
  rule caches every pull (`fi/CLAUDE.md:88`).
- **Observed:** scripts are Python, standard library only (the fi
  scripts rule, and the skills' scripts).
- **Observed:** `ta.py` is deterministic. The outlook uses a fixed seed
  (`.claude/skills/analysis-stock/scripts/ta.py:984`), so frozen inputs
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

- **The block:** a fenced JSON block that ends each new report. It has
  three lists:
  - **claims:** each with the value as printed, its source (a file and
    a field, or a URL and a date), and its kind (price, percent, count,
    date);
  - **trades:** each with the account type (taxable, Roth, HSA, 401k),
    side, quantity, order type and price, trigger, fallback, the stance
    it rests on, and its chart reason or "no technical reason";
  - **predictions:** each with its kind (close above, close inside a
    range, trades at a level, daily close beyond a level), level or
    range, probability and check date.
- **Validator:** a script that checks a block against the schema.
- **Rule:** a rule in `fi/CLAUDE.md` making the block required, which
  is edited at the user's request.
- **First extraction:** the two prediction reports dated 2026-11-02
  under `fi/reports/stocks/` (one for a single ticker, one for 13
  tickers) are extracted into calls files under `fi/reports/evals/`, so piece 2 can score them
  on 11-06.
- **Done when:** both files validate; a deliberately broken block fails
  with a message naming the field; and the user approves the schema.

### 2. `advice-eval-scorecard`

Goal: prediction outcomes, with honest uncertainty.

- **Inputs:** the calls files, plus daily bars covering each check date.
  The bars come from `data/` when it's current; otherwise the scorecard
  spends calls under the cache rule, or marks the prediction pending.
- **Hit or miss, by kind:**
  - "close above" and "close inside a range" use the close on the
    check date;
  - "trades at a level" uses each session's high and low up to the
    check date;
  - "close beyond a level" uses the daily closes up to the check date.
- **Output:** the Brier score, hit rates with Wilson 95% intervals, and
  the edge over a coin flip and over the no-edge model's probabilities.
  - A verdict is printed only when the interval excludes the baseline.
  - It also states how many more scored calls it would take to detect
    the edge.
- **Report:** a dated report under `fi/reports/evals/`.
- **Done when:** a test fixture with known outcomes scores to hand-
  computed values, and the 11-06 predictions are scored that day.

### 3. `advice-eval-grader`

Goal: one command that grades one report.

- **Exact checks:**
  - every claim matches its source at the printed precision, and a
    claim with no source is reported as "untraceable";
  - the required sections are present (`fi/CLAUDE.md:119`);
  - every claim carries a tag;
  - every trade has a chart reason or "no technical reason"
    (`fi/CLAUDE.md:51`);
  - every recommendation has its reason, main risk and rules check
    (`fi/CLAUDE.md:48`);
  - no plan level or flip level has better-than-even odds, in `ta.py`'s
    outlook, of a daily close beyond it within 10 sessions;
  - the account rules hold (`fi/profile.md:36`), for example covered
    calls only and the leveraged-ETF cap.
- **Claude grader:** a written rubric for the judgment rules, applied
  with quoted evidence.
  - Its scores count only after agreeing with the user's 20 labels.
- **Findings:** each one names the file, the claim, the expected value,
  the value found, and the source.
- **Done when:**
  - seeded broken reports (a changed number, a missing section, a
    missing chart reason, a trigger inside the noise) each produce
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
- **Done when:** the run passes today, and fails after a deliberate
  one-line change to the stance bands.

### 5. `advice-eval-audit`

Goal: grade every existing report.

- **How:** Claude reads each of the 31 reports once and writes its calls
  file under `fi/reports/evals/`. Then the grader runs on all of them,
  and the result is an audit report.
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

- **Baselines:** each trade suggestion's result over 1, 2 and 4 weeks,
  compared with no trade, with SPY over the same dates, and with the
  rule's own entry, stop and target applied mechanically.
- **Taxes:** pre-tax for every trade. Taxable-account trades also get
  an after-tax result, from the lot basis and the plan's tax rate.
- **Edge test:** the same intervals, verdict rule and power target as
  piece 2.
- **Done when:** a fixture with known fills and bars gives hand-
  computed results.

## Open questions

Rounds 1–7 are answered (see the log). Each question below has an owner and a default that applies if it isn't
answered.

1. **What agreement with the user's labels makes the grader count?**
   Owner: the user. Default: at least 16 of 20 (80%), re-checked when
   the rubric changes.
2. **What edge should the power target aim at?** Owner: the user.
   Default: a 5-point gain in hit rate over the no-edge model's
   probabilities, at 80% power and a 95% confidence level.
3. **How are check dates counted around market holidays?** Owner:
   Claude, in piece 2. Default: by trading sessions, using the bars
   that actually exist. A prediction whose date falls on a holiday is
   scored on the last session before it.
4. **What if a source is unavailable when a claim is checked?** Owner:
   Claude, in piece 3. Default: mark the claim "unverified" and list
   it, but don't count it as a failure.
5. **Where do outputs go?** Owner: Claude. Default:
   `fi/reports/evals/`, one dated report per run. Each lists its
   findings, its revisions, what's pending, and the calls it spent.

Settled with less questioning, because the blast radius is small:

- **Scale (inferred):** one user, about 14 tickers, about 70
  predictions a week, and 31 old reports, all run on demand. Standard
  library Python is enough.
- **Blast radius:** the new rule in `fi/CLAUDE.md` changes every future
  report. `ta.py` is unchanged. The eval reads old reports but never
  edits them.
- **Deadline and budget:** the only deadline is piece 2 by 11-06. The
  budget is the free tier's calls, with no paid services (round 6).
  The Claude grader runs inside the session.

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
