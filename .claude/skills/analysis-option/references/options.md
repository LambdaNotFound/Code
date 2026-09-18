# Reading an option quote

How to turn what `scripts/opt.py` prints into a judgement. Every rule
here is a textbook default `[COMMON]`, not a measured edge. The pricing
model is Black-Scholes on a driftless lognormal: a `[FRAME]` that prices
the contract consistently and says nothing about direction. Any claim
about what the trade will make is capped at LOW confidence.

## Tagging

- A number from the script: `[COMPUTED]`. Say which inputs it rests on
  when they were assumptions (the rate, the dividend yield, the quote's
  timestamp).
- A textbook read of a number (spread over 10% is illiquid): `[COMMON]`.
- The verdict and any probability of profit: `[FRAME]`, LOW confidence.
  P(profit) at HV assumes the stock's recent volatility persists and
  that it has no drift; both are wrong in ways that matter.
- The quote itself is `[KNOWN]` only as of the time the user took it.
  Options quotes go stale in minutes; say when it was taken.

## What the script computes and what to make of it

| Line | Read |
|---|---|
| Premium used | a buyer pays the ask, a seller receives the bid. Mid is the reference for fairness, not what you get |
| Spread as % of mid | under 5% is fine, 5 to 10 costs a round trip's worth of edge, over 10 means the market for this strike is thin. Over 25 is a hard fail |
| Open interest | under 100 contracts and you may not be able to leave; 500 or more is a real market |
| IV source | "quoted" came from the user or the chain; "solved from mid" is the script's own inversion. If the two disagree by more than 2 vol points, the quote and the spot are from different moments |
| IV / HV | implied over realized (60-day, or 20-day when there are fewer bars). Buyers want it under 0.9, sellers over 1.2. Around 1 the option is priced at what the stock has actually been doing |
| Model value at HV | what Black-Scholes says the contract is worth if the future looks like the last 60 days. Mid above it means the market expects more movement than the past shows, or an event |
| Theta as % of mid per day | how fast the premium decays today. A buyer paying 2% a day needs the move to come fast |
| Break-even in HV sigmas | the move the buyer needs, in units of the stock's own expected move. Inside half a sigma is ordinary; beyond one sigma needs a catalyst |
| Expected move (IV vs HV) | the market's one-sigma range to expiry and the realized one. When IV's is much wider, the market is pricing an event |
| P(profit) at HV | the chance the position ends above zero at expiry under the driftless lognormal. It is not the chance of making money before expiry, which is higher for a buyer who can exit on a move |
| EV at HV | expected P&L under the same model. It is near zero by construction for a fairly priced option; the sign comes from IV/HV and the spread. Over +5% of premium is an edge in the model's terms, under -10% is paying up |
| Reward to risk | buyer: P&L at a one-sigma move in your favour divided by max loss. Seller: premium divided by the loss at a two-sigma move against. Neither is a forecast |
| Scenarios | P&L at expiry at plus and minus one and two sigma, the strike, the break-even, and the analysis-stock levels when a TA JSON was given |

## Position types

| Position | Max loss | Max gain | What the script assumes |
|---|---|---|---|
| buy call | premium | unbounded | |
| buy put | premium | strike minus premium | |
| sell put | strike minus premium | premium | cash-secured: capital tied up is the strike times 100 |
| sell call, covered | in the stock, not the option | premium plus strike minus cost basis | `--covered`; the script scores the option leg only |
| sell call, naked | unbounded | premium | hard fail: the verdict is POOR whatever the other rows say |

Assignment: a short American-style put or call can be assigned before
expiry, most likely when it is in the money near an ex-dividend date or
with little extrinsic value left. The script does not model early
exercise; say so when the position is short and in the money.

## Earnings and events

An earnings date inside the window changes the meaning of every
volatility line. For a seller it is gap risk the premium may not cover;
the ledger marks it bad. For a buyer it is why IV is high, and the
implied vol will collapse after the print, so the stock has to move
more than the expected move for the position to profit. The script only
knows the date if the user gives it. Ask for it when the expiry is more
than two weeks out and the user did not say.

## Verdict rule

Every ledger row reads good, bad, or neutral. Score = (good − bad) /
rows. GOOD at +0.40 and above, POOR at −0.20 and below, FAIR between.
Hard fails override to POOR: spread over 25% of mid, P(profit) under
10%, no solvable implied vol (a stale or crossed quote), or a naked
short call. The rule is a transparent tally of textbook checks; it has
no measured hit rate and it does not know the user's thesis, the
account, or the rest of the portfolio. Report it as `[FRAME]` at LOW
confidence and name the rows behind it. Never adjust the label by hand;
if the story and the label disagree, the verdict paragraph is where to
say so.

## Alignment with the chart

With a `--ta-json` from the analysis-stock script, the ledger scores
whether the position's direction matches the chart stance: long call or
short put is bullish, long put or short call is bearish. A HOLD stance
is neutral. This row is one vote among many; the chart stance is itself
a `[FRAME]`.

## What the script does not check

Multi-leg structures (spreads, straddles, collars), early exercise,
dividend capture, margin requirements, tax treatment, IV rank or
percentile (needs an IV history the free data does not supply), and
the term structure across expiries. Say which of these matter for the
position and that the skill did not look at them.
