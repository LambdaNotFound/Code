---
description: 'Risk/reward analysis of one listed stock option contract: is this call or put, bought or sold, a good deal at the quoted price. Takes the contract quote from the user (bid, ask, and optionally IV, open interest, volume; the free Alpha Vantage tier has no options chain) or from a chain file, fetches the underlying''s daily bars through the alphavantage MCP server, and computes with scripts/opt.py: implied vol solved from the mid, realized vol, Black-Scholes greeks and model value, theta burn, break-even in sigmas, expected move, payoff scenarios, P(profit) and expected value, reward-to-risk, a ledger, and a GOOD / FAIR / POOR verdict with the rows behind it. Can take the analysis-stock JSON to score alignment with the chart stance. Use when the user names an option (ticker, expiry, strike, call or put) and asks whether to buy or sell it, whether the premium is cheap or rich, what the break-even or odds are, or whether the risk/reward is worth it. Not for the stock''s own chart read (use analysis-stock). Not for multi-leg strategies, portfolio hedging, or margin and tax questions. Not for fundamentals or news.'
argument-hint: <TICKER> <YYYY-MM-DD> <STRIKE> <call|put> [buy|sell] --quote bid=,ask=[,iv=,oi=,volume=] [--earnings date] [--with-stance] [--account equity]
---

You judge one option contract on its risk and reward, in the main
session, from a quote the user gives you and numbers a script computes.
You do not price multi-leg structures, you do not touch fundamentals or
news, and you do not invent a quote. The verdict is the model's answer,
stated with its rule and the rows behind it; the decision is the
user's, and you say so once, without a lecture.

Everything below is a standing instruction for the whole task.

## Arguments

$ARGUMENTS

Order: ticker, expiry (YYYY-MM-DD), strike, `call` or `put`, then an
optional `buy` (default) or `sell`. Flags: `--quote bid=,ask=[,iv=,oi=,volume=,last=]`
(IV as a decimal, 0.32 not 32), `--chain <file>` instead of `--quote`,
`--earnings <date>`, `--covered` for a short call backed by stock,
`--contracts N`, `--account <equity>`, `--rate` and `--div-yield`
(decimals; defaults 0.04 and 0), `--with-stance` to run the
analysis-stock script first and score chart alignment, `--save <path>`.

If the user named the option in prose ("the October 340 call on
Apple") resolve it to those tokens and confirm in one line. If they
gave no quote and no chain file, stop and ask for bid and ask, and when
they took them. Do not estimate a quote from the model; the point of
the skill is to judge the price the market is showing.

## 1. Fetch (one MCP call; two more with --with-stance)

The options-chain endpoints (`HISTORICAL_OPTIONS`, `REALTIME_OPTIONS`)
are premium on the free key; do not call them unless the user says
their key has them. The underlying's bars are free.

1. `mcp__alphavantage__TIME_SERIES_DAILY` with `outputsize=compact`,
   `datatype=csv`, `return_full_data=true`. Save an inline result to
   the scratchpad as `<SYMBOL>_daily.csv`; use the saved path if the
   harness wrote one.
2. With `--with-stance`: also `TIME_SERIES_WEEKLY_ADJUSTED` and the
   daily 200-`SMA`, exactly as [../analysis-stock/SKILL.md](../analysis-stock/SKILL.md)
   describes, then run that skill's script with `--json` and save the
   output as `<SYMBOL>_ta.json`. Skip this when the user already ran
   analysis-stock this session and the JSON is in the scratchpad.

Rate-limit handling is the same as analysis-stock: "premium endpoint"
means pick the fallback; anything else, wait a few seconds and retry
once. The earnings calendar endpoint has not returned usable data on
this key; ask the user for the date instead when it matters (expiry
more than two weeks out).

## 2. Compute

```
python3 .claude/skills/analysis-option/scripts/opt.py --symbol <SYM> --expiry <date> --strike <K> --type <call|put> --side <buy|sell> \
  --daily <daily file> --quote bid=<b>,ask=<a>[,iv=,oi=,volume=] [--earnings <date>] [--ta-json <SYM>_ta.json] [--covered] [--contracts N] [--account E]
```

The script prints the quote and liquidity, implied vol (quoted or
solved from the mid) against 20- and 60-day realized vol, the model
value at realized vol, greeks, theta as a share of premium per day,
break-even and the move it needs in sigmas, expected move at IV and at
HV, max loss and gain, P&L at one sigma for and two sigma against,
reward-to-risk, P(profit) and expected value under a driftless
lognormal at HV and at IV, a scenario table at expiry (with the
analysis-stock levels when a TA JSON was given), a ledger, and the
verdict. `--json` gives every field. Read all of it before writing. If
it errors, show the error and fix the input; never hand-compute a
substitute. If implied vol comes back unavailable, the quote is stale
or crossed relative to the spot: say so and ask for a fresh one.

## 3. Read

Apply [references/options.md](references/options.md): what each line
means, how buyers and sellers read the same number differently, what
earnings inside the window does to every volatility line, the position
table (a naked short call is a hard fail), and the tag each kind of
statement carries. Do not invent thresholds the reference does not
have.

## 4. Write

Fill [references/report-template.md](references/report-template.md)
section by section: verdict first, the contract, volatility and price,
the move it needs, risk and odds, scenarios, the recommendation table
with the script's verdict pasted as is, the ledger unchanged, and what
was not covered. Every claim tagged; P(profit), EV, and the verdict are
`[FRAME]` and never above LOW confidence. Say when the quote was taken.
Close with the `[RULES I BROKE]` line.

The recommendation is the script's verdict, not your own. Do not
override the label because the narrative feels different; if they
disagree, say that in the verdict paragraph. Print the report in the
reply. With `--save`, also write it to the path given. Without it, do
not write files into the repository.

## Refusals

- A spread, straddle, collar, or any second leg: say the skill prices
  one contract, offer to run each leg separately, and stop.
- "Will it go up?" gets the analysis-stock stance if a TA JSON exists,
  otherwise the sentence that this skill prices the contract and the
  chart read lives in analysis-stock.
- A quote with no time attached: run, and say in the first line of the
  report that the quote's time is unknown.

## What you return

1. The report, in the template's order, in the reply.
2. The path, if `--save` was given.
3. The number of Alpha Vantage calls used.
