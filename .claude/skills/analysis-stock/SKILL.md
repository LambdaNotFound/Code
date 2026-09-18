---
description: 'Technical analysis of one stock from live Alpha Vantage price data — trend, momentum, volatility, volume, support and resistance on daily and weekly timeframes, then two to four ranked scenarios with triggers and invalidation levels. Fetches with the alphavantage MCP tools (three calls, inside the free tier''s 25-a-day cap), computes every indicator deterministically with scripts/ta.py, and writes the report to the template in references/report-template.md with every claim tagged per this repo''s epistemic rules. Use when the user names a ticker and asks for technical analysis, a chart read, support and resistance, trend or momentum, or whether a stock is overbought or oversold. Ends with a rule-based BUY / ACCUMULATE / HOLD / REDUCE / SELL stance with entry, stop, target, and reward-to-risk, derived from the same ledger and labelled as the frame''s output, not a forecast. Use also when the user asks whether to buy or sell a stock on its chart. Not for fundamentals, valuation, earnings, or news (say so and stop). Not for position sizing or portfolio decisions. Not for backtesting a strategy.'
argument-hint: <TICKER> [daily-file weekly-file] [--save path]
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
(any shape the script accepts); skip the fetch. `--save <path>` writes
the finished report there as well as printing it.

## 1. Fetch (three MCP calls, in this order)

The free Alpha Vantage tier allows 25 requests a day and one a second.
Every call below costs one; do not add exploratory ones. Adjusted daily
data and the full daily history are premium, which is why the daily
view is 100 bars and the long view comes from the weekly feed.

1. `mcp__alphavantage__TIME_SERIES_DAILY` with `outputsize=compact`,
   `datatype=csv`, `return_full_data=true`. 100 bars. If the result
   arrives inline, write it to a file in the scratchpad directory
   named `<SYMBOL>_daily.csv`. If it arrives as a saved file, use
   that path.
2. `mcp__alphavantage__TIME_SERIES_WEEKLY_ADJUSTED` with
   `datatype=csv`, `return_full_data=true`. Twenty-plus years; it
   always arrives as a saved file (the harness writes oversized
   results to a tool-results folder under the user's Claude projects
   directory and prints the path). Use that path directly. Fall back to
   `TIME_SERIES_WEEKLY` only if the adjusted endpoint is refused, and
   then say in the report that the long-range levels are unadjusted.
3. `mcp__alphavantage__SMA` with `interval=daily`, `time_period=200`,
   `series_type=close`, `datatype=json`. The preview shows the latest
   value; that number is `--sma200`. Skip the `return_full_data` flag
   here; the preview is all you need.

If a call returns `{"error": {"type": "rate_limit", ...}}` read the
message: "premium endpoint" means the endpoint is not available on
this key (pick the fallback above); anything else means wait a few
seconds and retry once. Two failures on the same call: stop, report
which call failed and what it said, and offer to run from a file the
user supplies.

If the user set `ALPHAVANTAGE_API_KEY` in the environment, the script
can fetch on its own instead: `python3 .claude/skills/analysis-stock/scripts/ta.py --fetch <SYMBOL> --out-dir <scratchpad>`.
Do not go looking for the key; use this route only when the user says
it is there.

## 2. Compute

```
python3 .claude/skills/analysis-stock/scripts/ta.py --daily <daily> --weekly <weekly> --sma200 <value> --symbol <SYMBOL>
```

Add `--json` when you need a field the markdown does not show. The
script prints price, returns, ranges, moving averages and their
slopes and crosses, ADX, RSI, MACD, stochastic, Bollinger, ATR, OBV,
up/down volume, swing structure, clustered support and resistance,
open gaps, the weekly equivalents, a signal ledger with a tally, and
the stance: a weighted score over the ledger mapped to BUY,
ACCUMULATE, HOLD, REDUCE, or SELL, with entry, stop, target, and
reward-to-risk from the nearest levels and the ATR. Read all of it
before writing. If the script errors on the input,
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

- "How much should I buy?" gets the stance and its stop, and the
  sentence that sizing depends on the account, which this skill does
  not see.
- A request to include fundamentals, earnings, or news gets a one-line
  refusal naming the boundary; the analysis still runs.
- A ticker the search cannot resolve: stop and ask, do not guess a
  symbol.

## What you return

1. The report, in the template's order, in the reply.
2. The path, if `--save` was given.
3. The number of Alpha Vantage calls used, so the user can track the
   daily budget.
