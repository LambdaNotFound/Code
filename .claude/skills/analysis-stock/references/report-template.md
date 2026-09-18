# Technical analysis: <SYMBOL>

Data: daily through <date> (<n> bars), weekly through <date>; Alpha Vantage,
<adjusted | unadjusted> weekly. Every number `[COMPUTED]` unless marked.
Technical analysis is a `[FRAME]`; nothing below is a forecast or advice.

## Verdict

One paragraph. Primary trend (weekly), the daily posture inside it, the
one thing that would change the read, and the ledger tally. Tag it
`[INFERRED]`, confidence LOW or lower.

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
timeframe.

## Momentum

RSI14, MACD line/signal/histogram and its direction, stochastic, any
divergence the script found. One read per indicator, each tagged
`[COMMON]` for the textbook threshold.

## Volatility and volume

ATR as % of price, Bollinger position and bandwidth percentile, volume
vs 20-day average, up/down volume ratio, OBV agreement. State whether
volume confirms the price action or not.

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

**<LABEL>** (score <s>, net <n> of <t>). Confidence LOW.

- Entry / stop / target / reward:risk, as the script printed them.
- For: the ledger rows that pulled the score up. Against: the rows
  that pulled it down.
- Flips to <label> on a close above <resistance 1>; to <label> on a
  close below <support 1>.
- One sentence: the rule is fixed and mechanical, and it does not see
  earnings, news, or the account; the decision and the size are the
  reader's.

## Signal ledger

Paste the script's ledger table, tally, and stance block unchanged.

## Not covered

Fundamentals, earnings dates, news, sector and index relative strength,
options positioning. Say which of these the user should check before
acting on anything above.

[RULES I BROKE]: which, where, why (or "none").
