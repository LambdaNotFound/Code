# Technical analysis: <SYMBOL>

Data: daily through <date> (<n> bars), weekly through <date>; Alpha Vantage,
<adjusted | unadjusted> weekly. Every number `[COMPUTED]` unless marked.
Technical analysis is a `[FRAME]`; nothing below is a forecast or advice.

## Verdict

One paragraph. Primary trend (weekly), the daily posture inside it, the
one thing that would change the read, and the ledger tally. Tag it
`[INFERRED]`, confidence LOW or lower.

## Chart

The rendered SVG, sent as a file. One line naming what is on it and
the one thing worth looking at.

## Price and trend

| | Daily | Weekly |
|---|---|---|
| Close | | |
| vs 20 / 50 / 200 MA | | |
| MA stack | | |
| Swing structure | | |
| ADX (+DI / -DI) | daily only | |
| Last MA cross | | |

Two or three sentences: is this a trend or a range, and on which
timeframe. Then the trend template line: N/7 passed, which failed, and
the approximations the script named.

## Momentum

RSI14, MACD line/signal/histogram and its direction, stochastic, any
divergence the script found. One read per indicator, each tagged
`[COMMON]` for the textbook threshold.

## Volatility and volume

ATR as % of price, Bollinger position and bandwidth percentile, volume
vs 20-day average, up/down volume ratio, OBV agreement, burst days in
the last 20 bars. State whether volume confirms the price action or not.

## Monthly (multi-year)

The 12- and 24-month averages and price against them, the returns, the
full-history range, monthly swing structure, and the monthly levels.
State the two caveats the script prints: partial newest month, and
weeks assigned to the month they closed in.

## Expected move

Paste the script's table. One sentence that this is realized
volatility, not what the options market is pricing, and that
analysis-option with a live quote gives the comparison.

## Weekly reversal checks

The script's read and each triggered check with its week and level, or
"none in the last 8 completed weeks". If a check was vetoed, say what
vetoed it.

## Levels

| Level | Price | Kind | Evidence |
|---|---|---|---|
| Resistance 2 | | swing cluster / MA / 52w high / gap | touches, last date |
| Resistance 1 | | | |
| **Current** | | | |
| Support 1 | | | |
| Support 2 | | | |

Nearest first on each side. Name the source of every row.

## Scenarios `[FRAME]`

### 1. <name> — <p>%
Trigger, target, invalidation, and which ledger rows support it.

### 2. <name> — <p>%
Same shape.

### 3. <name> — <p>%
Same shape. Probabilities sum to 100. Confidence on any of them: LOW.

## Recommendation `[FRAME]`

Paste the script's Stance table as is, then fill the last two rows.

| | |
|---|---|
| **Stance** | **<LABEL>** |
| Score | <s> (net <n> of <t>) |
| Entry | |
| Stop | |
| Target | |
| Reward:risk | |
| 2R target | |
| For | ledger rows that pulled the score up |
| Against | ledger rows that pulled it down |
| Flips up on | close above <resistance 1> |
| Flips down on | close below <support 1> |
| Confidence | LOW; fixed mechanical rule, no measured hit rate |
| Not seen | earnings date, news, the account; size is the reader's |

With an account size given, paste the script's Position size table
under the recommendation.

## Signal ledger

Paste the script's ledger table, tally, and stance block unchanged.

## Not covered

Fundamentals, earnings dates, news, sector and index relative strength,
options positioning. Say which of these the user should check before
acting on anything above.

[RULES I BROKE]: which, where, why (or "none").
