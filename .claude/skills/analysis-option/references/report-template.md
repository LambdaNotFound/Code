# Option analysis: <SIDE> <N>x <SYMBOL> <expiry> <strike> <call|put>

Quote taken <when, per the user>; underlying daily data through <date>;
rate <r> and dividend yield <q> assumed. Every number `[COMPUTED]` unless
marked. Black-Scholes on a driftless lognormal is a `[FRAME]`; nothing
below is a forecast or advice.

## Verdict

One paragraph. Is this a good risk/reward deal, in the model's terms,
and why: the two or three ledger rows that decide it, the one number
the trade depends on (usually the move needed versus the expected move,
or IV against HV), and what would change the label. Tag it `[INFERRED]`,
confidence LOW or lower. If the ledger's label and the story disagree,
say so here.

## The contract

| | |
|---|---|
| Position | |
| Days to expiry | |
| Spot / strike / moneyness | |
| Bid / ask / mid; premium used | |
| Spread and open interest | |
| Intrinsic / extrinsic | |

## Volatility and price

| | |
|---|---|
| Implied vol (source) | |
| Realized vol 20d / 60d | |
| IV / HV | |
| Model value at HV vs mid | |
| Delta / gamma / theta per day / vega | |
| Theta as % of mid per day | |

Two sentences: is the premium cheap or rich against what the stock has
been doing, and what the market seems to be pricing (`[COMMON]` reads).

## The move it needs

| | |
|---|---|
| Break-even at expiry | |
| Expected move to expiry, IV and HV | |
| Break-even in HV sigmas | |

## Risk and odds `[FRAME]`

| | |
|---|---|
| Max loss / max gain / capital tied up | |
| P&L at one sigma for, two sigma against | |
| Reward to risk | |
| P(profit) and EV at HV; at IV | |

One sentence on what P(profit) does and does not mean here.

## Scenarios at expiry

Paste the script's scenario table. Add one line on which scenario the
analysis-stock stance favours, if a TA JSON was given.

## Recommendation `[FRAME]`

| | |
|---|---|
| **Deal** | **<GOOD / FAIR / POOR>** |
| Score | |
| For | |
| Against | |
| Hard fails | |
| Flips on | the one or two changes that would move the label (a better fill, a different strike or expiry, IV coming in) |
| Confidence | LOW; fixed mechanical rule, no measured hit rate |
| Not seen | the thesis, the account, the rest of the portfolio, early assignment, and (unless given) the earnings date |

## Ledger

Paste the script's ledger table and tally unchanged.

## Not covered

Multi-leg alternatives, IV rank, term structure, margin, tax, early
exercise. Name the one or two that matter most for this position.

[RULES I BROKE]: which, where, why (or "none").
