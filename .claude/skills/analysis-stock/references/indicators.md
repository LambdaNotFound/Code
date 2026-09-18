# Indicator reading rules

How to turn the numbers `scripts/ta.py` prints into a read. Every threshold
here is the textbook default `[COMMON]`; none of them is a measured edge.
Technical analysis is a `[FRAME]`: internally coherent, not a claim about
what the stock will do. A conclusion drawn inside it stays inside it and
caps at LOW confidence when stated about the real world.

## Tagging the report

- A number from the script: `[COMPUTED]`.
- A textbook interpretation of that number (RSI 72 is "overbought"):
  `[COMMON]`.
- A read that combines several of them into a trend or a setup:
  `[INFERRED]`, with the rows it rests on named.
- A scenario probability: `[FRAME]`, and never above LOW confidence. The
  probabilities are a way of ranking scenarios against each other, not a
  forecast.
- Anything about earnings, news, fundamentals, or the macro backdrop is
  out of scope. If the user asks, say it is outside this skill rather
  than guessing.

## Trend

| Signal | Bullish | Bearish | Neutral |
|---|---|---|---|
| MA stack | price > 20 > 50 > 200, all rising | price < 20 < 50 < 200 | anything mixed |
| 50/200 cross | golden cross (50 up through 200) | death cross | none in window |
| ADX14 | >= 25 with +DI > -DI | >= 25 with -DI > +DI | < 25: no trend to read, whatever the DIs say |
| Swing structure | higher highs and higher lows | lower highs and lower lows | mixed |

An ADX under 20 with price chopping across the 20-day is a range, and
every trend-following read below is weaker inside a range. ADX above 40
is a trend already well underway; late, not early.

The script's 200-day comes from three sources, in this order of
trustworthiness: the SMA endpoint (`daily`), a 200-bar daily series
(`daily`), or a 40-week SMA of weekly closes (`approx-40-week`). Say
which one the report used.

## Momentum

| Signal | Read |
|---|---|
| RSI14 >= 70 | overbought; in a strong uptrend it can sit there for weeks, so it is a caution, not a sell |
| RSI14 <= 30 | oversold; same caveat in a downtrend |
| RSI14 40-60 while ADX < 20 | no momentum information |
| RSI divergence | the script checks only the last two swings; a divergence is a warning that the move is tiring, not a reversal signal on its own |
| MACD histogram | sign gives direction; the change from the previous bar says whether momentum is expanding or contracting. A shrinking positive histogram is a weakening uptrend, not a bearish signal yet |
| MACD cross | the script dates the last one; a cross more than ~15 bars old is stale |
| Stochastic %K >= 80 / <= 20 | overbought / oversold; faster and noisier than RSI, so weight it less |

## Volatility

| Signal | Read |
|---|---|
| ATR14 as % of price | the stock's normal daily range; use it to size a stop, and to judge whether a move was large for this stock rather than in absolute terms |
| Bollinger %B > 1 or < 0 | a close outside the bands; in a trend it is continuation, in a range it is mean-reversion. Decide which from ADX first |
| Bandwidth percentile <= 20 | squeeze: volatility compressed relative to its own recent history, which tends to precede a move of unknown direction |
| Bandwidth percentile >= 80 | expansion already happened |

## Volume

| Signal | Read |
|---|---|
| Last volume vs 20-day average >= 1.5x | the day mattered; note what price did on it |
| Up/down volume ratio (20 bars) >= 1.3 | accumulation |
| Up/down volume ratio <= 0.77 | distribution |
| OBV rising with price | confirms the advance |
| OBV falling while price rises | divergence: the advance is on thinning participation |

## Levels

The script clusters swing highs and lows within 1.5% (daily) or 2%
(weekly) of each other and counts touches. Reading them:

- More touches and a more recent last touch make a level more relevant.
  One touch is a candidate, not a level.
- A cluster of prior swing highs now below price is prior resistance
  acting as support (role reversal). Say so.
- Moving averages and the Bollinger mid-band are dynamic levels; report
  them alongside the static ones.
- The 20-day and 52-week highs and lows are levels whether or not they
  show up as swing clusters.
- An open gap (the script lists gaps of 2% or more in the last 20 bars)
  is a level: the gap's origin price is where a fill would end.

## Weekly versus daily

The weekly view decides the primary trend and the daily view times it.
When they disagree, say so, and let the weekly read carry the primary
scenario. A daily uptrend under a weekly downtrend is a counter-trend
rally until the weekly structure changes.

## Scenarios

Build two to four, mutually exclusive, probabilities summing to 100.
Each names its trigger (a level the price has to take out), its target
(the next cluster or MA in that direction), and the level that would
invalidate it. Rank them on the ledger: the side with more confirming
rows gets the higher number. Do not let the rank leak into confidence:
the whole exercise is `[FRAME]`.

## Stance (the buy/sell rule)

The script turns the ledger into one label by a fixed rule so that
two runs on the same data give the same answer:

- Trend rows weigh 2 (MA stack, price vs 200-day, ADX/DI, daily and
  weekly swing structure, weekly price vs 50/200-week). Momentum and
  volume rows weigh 1 (RSI, MACD histogram, stochastic, %B, RSI
  divergence, OBV agreement, up/down volume).
- Score = (bullish weight − bearish weight) / total weight, in [−1, 1].
- BUY at +0.50 and above, ACCUMULATE from +0.20, HOLD between −0.20
  and +0.20, REDUCE from −0.20 down, SELL at −0.50 and below.
- Plan: entry at the close; stop half an ATR beyond the nearest
  support (long) or resistance (short); target at the nearest level
  on the other side. A BUY whose reward-to-risk from the close is
  under 1.5 becomes an ACCUMULATE with the entry moved to the support,
  because the same trade is only worth taking from there.
- HOLD carries no plan, only the two closes that would change it.

What the rule is and is not: it is a transparent tally of textbook
reads, so it will be late at turns and wrong in ranges the way every
trend rule is. It is not a forecast, it has no measured hit rate, and
it knows nothing about earnings dates, news, or the user's account.
Report it as `[FRAME]` at LOW confidence, and name the rows behind it.
Never adjust the label by hand; if the story and the label disagree,
the verdict is where to say so.

## What the script does not check

Chart patterns (head and shoulders, cup and handle, flags, wedges),
candlestick patterns, Fibonacci retracements, Ichimoku, and relative
strength against an index. If you name one, it is your reading of the
numbers, tagged `[INFERRED]`, and it needs the levels it rests on in
the text. Do not name a pattern you cannot anchor to specific dates and
prices from the script output.
