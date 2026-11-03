"""Seeded synthetic markets for the analyse-stock tests.

`market(seed, ...)` draws a long daily random walk with overnight gaps, intraday
ranges and lognormal volume, then returns the last `daily_n` daily bars and the
whole history aggregated into weeks (open of the first session, highest high,
lowest low, last close, summed volume, dated on the week's last session). The
weekly bars are marked adjusted, as the weekly-adjusted endpoint's are. Prices
are rounded to 4 decimals and volumes to whole shares, like Alpha Vantage's.

The same seed always gives the same bars, so tests can pin expected outputs.
"""
import datetime
import math
import random

REGIMES = {          # name: (drift per session, daily sigma, chance of a 4% gap per session)
    "uptrend": (0.0012, 0.012, 0.0),
    "downtrend": (-0.0012, 0.015, 0.0),
    "range": (0.0, 0.010, 0.0),
    "volatile": (0.0, 0.032, 0.0),
    "gappy": (0.0004, 0.016, 0.06),
}


def weekdays_ending(last, n):
    out, d = [], datetime.date.fromisoformat(last)
    while len(out) < n:
        if d.weekday() < 5:
            out.append(d.isoformat())
        d -= datetime.timedelta(days=1)
    return list(reversed(out))


def market(seed, regime="range", daily_n=260, history=1500, last="2026-10-30", start=100.0):
    drift, vol, gap_p = REGIMES[regime]
    rng = random.Random(seed)
    bars, c = [], start
    for d in weekdays_ending(last, history):
        o = c * math.exp(rng.gauss(0.0, vol * 0.3))
        if gap_p and rng.random() < gap_p:
            o = c * math.exp(rng.choice((-1, 1)) * 0.04)
        cl = o * math.exp(rng.gauss(drift, vol))
        hi = max(o, cl) * (1 + abs(rng.gauss(0.0, vol * 0.5)))
        lo = min(o, cl) * (1 - abs(rng.gauss(0.0, vol * 0.5)))
        v = int(1_000_000 * math.exp(rng.gauss(0.0, 0.4)))
        bars.append({"date": d, "open": round(o, 4), "high": round(hi, 4), "low": round(lo, 4),
                     "close": round(cl, 4), "volume": float(v), "adjusted": False})
        c = cl
    weekly, cur = [], None
    for b in bars:
        wk = datetime.date.fromisoformat(b["date"]).isocalendar()[:2]
        if cur is None or cur[0] != wk:
            cur = (wk, dict(b, adjusted=True))
            weekly.append(cur[1])
        else:
            w = cur[1]
            w.update(date=b["date"], high=max(w["high"], b["high"]), low=min(w["low"], b["low"]),
                     close=b["close"], volume=w["volume"] + b["volume"])
    return [dict(b) for b in bars[-daily_n:]], weekly


def to_csv(bars, weekly=False):
    """Alpha Vantage CSV, newest row first."""
    if weekly:
        head = "timestamp,open,high,low,close,adjusted close,volume,dividend amount"
        rows = [f"{b['date']},{b['open']},{b['high']},{b['low']},{b['close']},{b['close']},{int(b['volume'])},0.0000"
                for b in reversed(bars)]
    else:
        head = "timestamp,open,high,low,close,volume"
        rows = [f"{b['date']},{b['open']},{b['high']},{b['low']},{b['close']},{int(b['volume'])}" for b in reversed(bars)]
    return head + "\n" + "\n".join(rows) + "\n"


def to_av_json(bars):
    """Alpha Vantage's TIME_SERIES_DAILY JSON shape."""
    return {"Meta Data": {"2. Symbol": "X"},
            "Time Series (Daily)": {b["date"]: {"1. open": str(b["open"]), "2. high": str(b["high"]),
                                                "3. low": str(b["low"]), "4. close": str(b["close"]),
                                                "5. volume": str(int(b["volume"]))} for b in bars}}
