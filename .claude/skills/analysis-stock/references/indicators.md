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

## Trend template (Minervini Stage 2) `[COMMON]`

Seven pass/fail checks: price above the 150- and 200-day; 150 above 200;
200-day rising over the last month; 50 above 150 above 200; price above
the 50-day; price at least 30% above the 52-week low; price within 25%
of the 52-week high. The eighth (relative-strength rank of 70 or more)
needs a universe and is not scored. Six or more of seven is a Stage 2
uptrend and reads bullish; two or fewer reads bearish. Two inputs are
approximated on the free tier and the script says so: the 150-day is the
30-week SMA and the 200-day slope is the 40-week SMA now against four
weeks ago.

## Weekly reversal checks `[COMMON]`

Run on completed weeks only (the last weekly bar is the week in
progress). Extremes are the prior 52 weeks; the window is the last 8
completed weeks.

| Check | Long-side trap (bearish) | Short-side mirror (bullish) |
|---|---|---|
| Key reversal | new 52-week high intraweek, close below the prior week's low | new 52-week low, close above the prior week's high |
| Failed extreme | traded above the prior 52-week high, closed back below it | traded below the prior low, closed back above it |
| Failed breakout | a weekly close above the prior high, then a close back below within 3 weeks (dated on the failure week) | mirror on lows |
| Continuation veto | a later weekly close at a new 52-week closing high negates the bearish signals | a new closing low negates the bullish ones |

A signal that stands reads against the trend it interrupted. It is a
warning about the crowd's position at an extreme, not a reversal by
itself; it enters the ledger at weight 1.

## Burst days `[COMMON]`

A close up 4% or more on volume above the prior day, or a daily range
wider than each of the prior three ranges when the prior day was not
already extended. Listed with the volume ratio and where the close sat
in the day's range (above 0.7 is strong). They are facts about
participation, not a ledger row: a burst on 2x volume that closed near
its high near a level is the evidence a breakout scenario needs.

## Position size (only with an account size)

Fixed fractional: risk dollars = account × risk percent (default 1);
shares = risk dollars ÷ (entry − stop), capped so the position is at
most 10% of the account by default. The tighter constraint wins and the
script names it. Two percent risk per trade is the ceiling anyone should
pass; the default stays at one.

## Stance (the buy/sell rule)

The script turns the ledger into one label by a fixed rule so that
two runs on the same data give the same answer:

- Trend rows weigh 2 (MA stack, ADX/DI, daily and weekly swing
  structure, weekly price vs 50/200-week, trend template). Momentum
  and volume rows weigh 1 (RSI, MACD histogram, stochastic, %B, RSI
  divergence, OBV agreement, up/down volume, weekly reversal checks),
  and so does price vs 200-day, because the MA stack and the trend
  template already count it; at weight 2 a clean Stage 2 stock carried
  four trend votes for one fact.
- Score = (bullish weight − bearish weight) / total weight, in [−1, 1].
- BUY at +0.50 and above, ACCUMULATE from +0.20, HOLD between −0.20
  and +0.20, REDUCE from −0.20 down, SELL at −0.50 and below.
- Plan: entry at the close; stop half an ATR beyond the nearest
  support (long) or resistance (short); target at the nearest level
  on the other side. A BUY whose reward-to-risk from the close is
  under 1.5 becomes an ACCUMULATE with the entry moved to the support,
  because the same trade is only worth taking from there. The short
  side mirrors it: a REDUCE or SELL whose reward-to-risk from the close
  is under 1.5 moves the exit up to the resistance, which for a holder
  means trim into strength, not at the close.
- The plan also prints a 2R target (entry plus twice the risk) beside
  the level target; when the level target is under 2R, the level is
  the honest one.
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
