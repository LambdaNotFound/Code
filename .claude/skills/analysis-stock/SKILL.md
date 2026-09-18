---
description: 'Technical analysis of one stock from live Alpha Vantage price data — trend, momentum, volatility, volume, support and resistance on daily and weekly timeframes, then two to four ranked scenarios with triggers and invalidation levels. Fetches with the alphavantage MCP tools (three calls, inside the free tier''s 25-a-day cap), computes every indicator deterministically with scripts/ta.py, and writes the report to the template in references/report-template.md with every claim tagged per this repo''s epistemic rules. Use when the user names a ticker and asks for technical analysis, a chart read, support and resistance, trend or momentum, or whether a stock is overbought or oversold. Also scores Minervini''s Stage 2 trend template, runs weekly key-reversal, failed-extreme, and failed-breakout checks at 52-week extremes, and lists momentum-burst days. Ends with a rule-based BUY / ACCUMULATE / HOLD / REDUCE / SELL stance with entry, stop, target, 2R target, and reward-to-risk, derived from the same ledger and labelled as the frame''s output, not a forecast, plus fixed-fractional share count when the user gives an account size. Use also when the user asks whether to buy or sell a stock on its chart, or how many shares a given risk budget allows. Not for fundamentals, valuation, earnings, or news (say so and stop). Not for portfolio-level decisions. Not for backtesting a strategy.'
argument-hint: <TICKER> [daily-file weekly-file] [--account equity [--risk-pct 1]] [--save path]
---

You do technical analysis of one stock, in the main session, from data
you fetch and numbers a script computes, and you end with a buy or sell
stance the script derived by a fixed rule. You do not read charts by
eye and you do not touch fundamentals or news. The stance is the
frame's answer, stated with its rule and its trade plan; the decision
is still the user's, and you say so once, without a lecture.

Everything below is a standing instruction for the whole task.

## Arguments

$ARGUMENTS

The first token is the ticker. If the user gave a company name instead,
resolve it with `mcp__alphavantage__SYMBOL_SEARCH` and confirm the
match in one line before spending any more calls. Two file paths after
the ticker mean the user already has the daily and weekly data on disk
(any shape the script accepts); skip the fetch. `--account <equity>`
turns on position sizing, with `--risk-pct` (default 1) and
`--max-position-pct` (default 10) passed through to the script.
`--save <path>` writes the finished report there as well as printing
it.

## 1. Fetch (three MCP calls, one per turn, in this order)

The free tier's limits, what is premium, and the three shapes a result
can arrive in are in
[references/alphavantage-free-tier.md](references/alphavantage-free-tier.md).
The short version: 25 calls a day, one a second, so never fire two
Alpha Vantage calls in the same turn; adjusted daily and the full daily
history are premium, which is why the daily view is 100 bars and the
long view comes from the weekly feed.

Before any call, look in the scratchpad for a `<SYMBOL>_daily.csv`,
weekly file, and `<SYMBOL>_ta.json` from today and reuse them; say so
in the report.

1. `mcp__alphavantage__TIME_SERIES_DAILY` with `outputsize=compact`,
   `datatype=csv`, `return_full_data=true`. 100 bars, arrives inline:
   write the whole `{"result": ...}` reply to `<SYMBOL>_daily.csv` in
   the scratchpad (the script unwraps it).
2. `mcp__alphavantage__TIME_SERIES_WEEKLY_ADJUSTED` with
   `datatype=csv`, `return_full_data=true`. For most stocks this is
   saved by the harness and the reply gives the path; pass it straight
   to `--weekly`. For a stock listed only a few years it arrives
   inline; write it to `<SYMBOL>_weekly.csv` the same way as the daily.
   Fall back to `TIME_SERIES_WEEKLY` only if the adjusted endpoint is
   refused, and then say in the report that the long-range levels are
   unadjusted.
3. `mcp__alphavantage__SMA` with `interval=daily`, `time_period=200`,
   `series_type=close`, `datatype=json`, no `return_full_data`. It
   arrives either as a server preview showing the two newest values
   (read the newest into `--sma200`) or as a file saved by the harness
   (pass the path to `--sma200-file`; the script takes every shape).

If a call returns `{"error": {"type": "rate_limit", ...}}` read the
message: "premium endpoint" or "premium feature" means not available on
this key (pick the fallback above); "1 request per second" means retry
once next turn; the 25-a-day message means stop for the day. Two
failures on the same call: stop, report which call failed and what it
said, and offer to run from a file the user supplies.

If the user set `ALPHAVANTAGE_API_KEY` in the environment, the script
can fetch on its own instead: `python3 .claude/skills/analysis-stock/scripts/ta.py --fetch <SYMBOL> --out-dir <scratchpad>`.
Do not go looking for the key; use this route only when the user says
it is there.

## 2. Compute

```
python3 .claude/skills/analysis-stock/scripts/ta.py --daily <daily> --weekly <weekly> (--sma200 <value> | --sma200-file <path>) --symbol <SYMBOL>
```

Add `--json` when you need a field the markdown does not show. The
script prints price, returns, ranges, moving averages and their
slopes and crosses, ADX, RSI, MACD, stochastic, Bollinger, ATR, OBV,
up/down volume, swing structure, clustered support and resistance,
open gaps, momentum-burst days, the weekly equivalents, the Minervini
trend template (7 checks), the weekly reversal checks on completed
weeks, a signal ledger with a tally, and the stance: a weighted score
over the ledger mapped to BUY, ACCUMULATE, HOLD, REDUCE, or SELL, with
entry, stop, target, 2R target, and reward-to-risk from the nearest
levels and the ATR. With `--account` it adds a fixed-fractional share
count. Read all of it before writing. If the script errors on the input,
show the error and fix the input; do not hand-compute a substitute.

## 3. Read

Apply [references/indicators.md](references/indicators.md) to the
numbers. The rules there decide what counts as bullish, bearish, or
no signal, how to weigh weekly against daily, how to rank scenarios,
and which tag each kind of statement carries. Do not invent
thresholds the reference does not have, and do not name a chart
pattern you cannot anchor to dates and prices in the script output.

## 4. Write

Fill [references/report-template.md](references/report-template.md)
section by section. The template's shape is the contract: verdict
first, then trend, momentum, volatility and volume, levels, two to
four scenarios with probabilities summing to 100, the recommendation,
the ledger pasted unchanged, and the list of what this skill did not
look at. Every claim tagged; scenario probabilities and the
recommendation are `[FRAME]` and never above LOW confidence.

The recommendation is the script's stance, not your own. Paste its
Stance table exactly as printed (label, score, entry, stop, target,
reward-to-risk, rows for and against, the two closes that flip it)
and add the template's last two rows. Do
not override the label because the narrative feels different; if the
ledger and the story disagree, say that in the verdict, which is
where judgment lives. Close with the `[RULES I BROKE]` line.

Print the report in the reply. With `--save`, also write it to the
path given, creating the directory if needed. Without `--save`, do
not write files into the repository.

## Refusals

- "How much should I buy?" without an account size gets the stance
  and its stop, and one line asking for the equity and risk per trade
  so the script can size it.
- A request to include fundamentals, earnings, or news gets a one-line
  refusal naming the boundary; the analysis still runs.
- A ticker the search cannot resolve: stop and ask, do not guess a
  symbol.

## What you return

1. The report, in the template's order, in the reply.
2. The path, if `--save` was given.
3. The number of Alpha Vantage calls used, so the user can track the
   daily budget.
