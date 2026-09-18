#!/usr/bin/env python3
"""Compute technical-analysis facts for one ticker from Alpha Vantage OHLCV data.

Stdlib only. Every number this prints is [COMPUTED] from the bars it was
given; the interpretation belongs to the caller (see ../references/indicators.md).

Inputs it accepts, in any of three shapes:
  * the raw CSV Alpha Vantage returns (timestamp,open,high,low,close,volume)
  * the file Claude Code writes for an oversized MCP result: {"result": "<csv>"}
  * the JSON Alpha Vantage returns with datatype=json ("Time Series (Daily)" or
    "Weekly Time Series")

Usage:
  ta.py --daily DAILY_FILE [--weekly WEEKLY_FILE] [--sma200 PRICE]
        [--symbol SYM] [--json] [--pivot-window N]
  ta.py --fetch SYM --out-dir DIR      # needs ALPHAVANTAGE_API_KEY in the env

--daily is the 100-bar "compact" daily series (the free tier's ceiling).
--weekly is the full weekly series (20+ years, free) and supplies the
52-week range, the 20/50/200-week averages, and the long trend.
--sma200 is the latest daily 200-SMA from the SMA endpoint; without it the
script substitutes a 40-week SMA and labels it approximate.
"""

import argparse
import csv
import io
import json
import math
import os
import sys
import urllib.parse
import urllib.request

# --------------------------------------------------------------------------- #
# Input
# --------------------------------------------------------------------------- #

FIELDS = ("date", "open", "high", "low", "close", "volume")


def _bar(date, o, h, l, c, adj, v):
    """One bar. When an adjusted close is present, scale open/high/low by the same
    factor so splits do not show up as cliffs; the "adjusted" flag records it."""
    o, h, l, c = float(o), float(h), float(l), float(c)
    adjusted = adj not in (None, "")
    if adjusted and c:
        k = float(adj) / c
        o, h, l, c = o * k, h * k, l * k, float(adj)
    return {"date": date, "open": o, "high": h, "low": l, "close": c, "volume": float(v), "adjusted": adjusted}


def _bars_from_csv(text):
    rows = list(csv.DictReader(io.StringIO(text.strip())))
    if not rows:
        raise ValueError("empty CSV")
    bars = []
    for r in rows:
        key = {k.strip().lower(): v for k, v in r.items() if k}
        date = key.get("timestamp") or key.get("date")
        if not date:
            raise ValueError("CSV has no timestamp/date column")
        bars.append(_bar(date.strip(), key["open"], key["high"], key["low"], key["close"],
                         key.get("adjusted close"), key["volume"]))
    return bars


def _bars_from_av_json(d):
    series_key = next((k for k in d if k.startswith(("Time Series", "Weekly", "Monthly"))), None)
    if series_key is None:
        raise ValueError("JSON has no Alpha Vantage time-series block")
    bars = []
    for date, row in d[series_key].items():
        col = {k.split(". ", 1)[-1]: v for k, v in row.items()}
        bars.append(_bar(date, col["open"], col["high"], col["low"], col["close"],
                         col.get("adjusted close"), col["volume"]))
    return bars


def parse_bars(text):
    """Return bars sorted oldest -> newest from any of the accepted shapes."""
    stripped = text.lstrip()
    if stripped.startswith("{"):
        d = json.loads(stripped)
        if "error" in d:
            raise ValueError(f"Alpha Vantage error: {d['error']}")
        if "result" in d and isinstance(d["result"], str):
            inner = d["result"].lstrip()
            bars = _bars_from_av_json(json.loads(inner)) if inner.startswith("{") else _bars_from_csv(inner)
        else:
            bars = _bars_from_av_json(d)
    else:
        bars = _bars_from_csv(text)
    bars.sort(key=lambda b: b["date"])
    return bars


def load_bars(path):
    with open(path, encoding="utf-8") as f:
        return parse_bars(f.read())


# --------------------------------------------------------------------------- #
# Indicators. Each returns a list aligned with its input; None where undefined.
# --------------------------------------------------------------------------- #

def sma(vals, n):
    out = [None] * len(vals)
    total = 0.0
    for i, v in enumerate(vals):
        total += v
        if i >= n:
            total -= vals[i - n]
        if i >= n - 1:
            out[i] = total / n
    return out


def ema(vals, n):
    out = [None] * len(vals)
    if len(vals) < n:
        return out
    k = 2.0 / (n + 1)
    prev = sum(vals[:n]) / n
    out[n - 1] = prev
    for i in range(n, len(vals)):
        prev = vals[i] * k + prev * (1 - k)
        out[i] = prev
    return out


def _ema_sparse(vals, n):
    """EMA over a list that may start with Nones (for MACD signal line)."""
    start = next((i for i, v in enumerate(vals) if v is not None), None)
    out = [None] * len(vals)
    if start is None:
        return out
    inner = ema(vals[start:], n)
    out[start:] = inner
    return out


def wilder_smooth(vals, n):
    """Wilder's smoothing: seed with the mean of the first n, then s = s + (v - s)/n."""
    out = [None] * len(vals)
    if len(vals) < n:
        return out
    prev = sum(vals[:n]) / n
    out[n - 1] = prev
    for i in range(n, len(vals)):
        prev = prev + (vals[i] - prev) / n
        out[i] = prev
    return out


def rsi(closes, n=14):
    out = [None] * len(closes)
    if len(closes) <= n:
        return out
    gains = [0.0] + [max(closes[i] - closes[i - 1], 0.0) for i in range(1, len(closes))]
    losses = [0.0] + [max(closes[i - 1] - closes[i], 0.0) for i in range(1, len(closes))]
    ag = wilder_smooth(gains[1:], n)
    al = wilder_smooth(losses[1:], n)
    for i in range(n, len(closes)):
        g, l = ag[i - 1], al[i - 1]
        if g is None:
            continue
        out[i] = 100.0 if l == 0 else 100.0 - 100.0 / (1.0 + g / l)
    return out


def macd(closes, fast=12, slow=26, signal=9):
    ef, es = ema(closes, fast), ema(closes, slow)
    line = [None if (a is None or b is None) else a - b for a, b in zip(ef, es)]
    sig = _ema_sparse(line, signal)
    hist = [None if (a is None or b is None) else a - b for a, b in zip(line, sig)]
    return line, sig, hist


def bbands(closes, n=20, k=2.0):
    mid = sma(closes, n)
    upper, lower = [None] * len(closes), [None] * len(closes)
    for i in range(n - 1, len(closes)):
        window = closes[i - n + 1:i + 1]
        m = mid[i]
        sd = math.sqrt(sum((x - m) ** 2 for x in window) / n)
        upper[i], lower[i] = m + k * sd, m - k * sd
    return upper, mid, lower


def true_range(bars):
    tr = [bars[0]["high"] - bars[0]["low"]]
    for i in range(1, len(bars)):
        h, l, pc = bars[i]["high"], bars[i]["low"], bars[i - 1]["close"]
        tr.append(max(h - l, abs(h - pc), abs(l - pc)))
    return tr


def atr(bars, n=14):
    return wilder_smooth(true_range(bars), n)


def adx(bars, n=14):
    """Return (adx, +DI, -DI), Wilder's method."""
    size = len(bars)
    none = [None] * size
    if size < 2 * n:
        return none, none, none
    tr = true_range(bars)
    pdm, mdm = [0.0], [0.0]
    for i in range(1, size):
        up = bars[i]["high"] - bars[i - 1]["high"]
        dn = bars[i - 1]["low"] - bars[i]["low"]
        pdm.append(up if up > dn and up > 0 else 0.0)
        mdm.append(dn if dn > up and dn > 0 else 0.0)
    # Wilder smoothing on the series starting at bar 1 (bar 0 has no move).
    str_ = wilder_smooth(tr[1:], n)
    spdm = wilder_smooth(pdm[1:], n)
    smdm = wilder_smooth(mdm[1:], n)
    pdi, mdi, dx = list(none), list(none), []
    for j in range(len(str_)):
        i = j + 1
        if str_[j] is None or str_[j] == 0:
            continue
        pdi[i] = 100.0 * spdm[j] / str_[j]
        mdi[i] = 100.0 * smdm[j] / str_[j]
        denom = pdi[i] + mdi[i]
        dx.append(0.0 if denom == 0 else 100.0 * abs(pdi[i] - mdi[i]) / denom)
    sadx = wilder_smooth(dx, n)
    out = list(none)
    first_dx_index = n  # dx[0] corresponds to bar index n (j = n-1, i = n)
    for j, v in enumerate(sadx):
        if v is not None:
            out[first_dx_index + j] = v
    return out, pdi, mdi


def stochastic(bars, k_n=14, k_smooth=3, d_n=3):
    """Slow stochastic: %K = SMA(k_smooth) of raw %K, %D = SMA(d_n) of %K."""
    raw = [None] * len(bars)
    for i in range(k_n - 1, len(bars)):
        window = bars[i - k_n + 1:i + 1]
        hh = max(b["high"] for b in window)
        ll = min(b["low"] for b in window)
        raw[i] = 50.0 if hh == ll else 100.0 * (bars[i]["close"] - ll) / (hh - ll)
    k = _sma_sparse(raw, k_smooth)
    d = _sma_sparse(k, d_n)
    return k, d


def _sma_sparse(vals, n):
    start = next((i for i, v in enumerate(vals) if v is not None), None)
    out = [None] * len(vals)
    if start is None:
        return out
    out[start:] = sma(vals[start:], n)
    return out


def obv(bars):
    out = [0.0]
    for i in range(1, len(bars)):
        c, pc, v = bars[i]["close"], bars[i - 1]["close"], bars[i]["volume"]
        out.append(out[-1] + (v if c > pc else -v if c < pc else 0.0))
    return out


# --------------------------------------------------------------------------- #
# Structure: pivots, levels, trend, divergence, gaps
# --------------------------------------------------------------------------- #

def pivots(bars, k=3):
    """Fractal swing points: a high above the k bars either side, or a low below them."""
    out = []
    for i in range(k, len(bars) - k):
        h, l = bars[i]["high"], bars[i]["low"]
        left, right = bars[i - k:i], bars[i + 1:i + k + 1]
        if all(h > b["high"] for b in left) and all(h > b["high"] for b in right):
            out.append({"index": i, "date": bars[i]["date"], "kind": "H", "price": h})
        if all(l < b["low"] for b in left) and all(l < b["low"] for b in right):
            out.append({"index": i, "date": bars[i]["date"], "kind": "L", "price": l})
    return out


def trend_structure(pivs):
    """Classify the last two swing highs and lows. Returns (label, detail)."""
    highs = [p for p in pivs if p["kind"] == "H"][-2:]
    lows = [p for p in pivs if p["kind"] == "L"][-2:]
    if len(highs) < 2 or len(lows) < 2:
        return "undetermined", "fewer than two swing highs and two swing lows in the window"
    hh = highs[1]["price"] > highs[0]["price"]
    hl = lows[1]["price"] > lows[0]["price"]
    tag = ("HH" if hh else "LH") + "/" + ("HL" if hl else "LL")
    label = {"HH/HL": "uptrend", "LH/LL": "downtrend"}.get(tag, "mixed")
    detail = (f"{tag}: highs {highs[0]['price']:.2f} ({highs[0]['date']}) -> {highs[1]['price']:.2f} "
              f"({highs[1]['date']}); lows {lows[0]['price']:.2f} ({lows[0]['date']}) -> "
              f"{lows[1]['price']:.2f} ({lows[1]['date']})")
    return label, detail


def cluster_levels(pivs, price, tol=0.015):
    """Group pivot prices within tol of each other; return supports below and resistances above."""
    pts = sorted(pivs, key=lambda p: p["price"])
    clusters = []
    for p in pts:
        if clusters and abs(p["price"] - clusters[-1]["price"]) / clusters[-1]["price"] <= tol:
            c = clusters[-1]
            c["touches"] += 1
            c["price"] = (c["price"] * (c["touches"] - 1) + p["price"]) / c["touches"]
            c["last"] = max(c["last"], p["date"])
        else:
            clusters.append({"price": p["price"], "touches": 1, "last": p["date"]})
    supports = sorted((c for c in clusters if c["price"] < price), key=lambda c: price - c["price"])
    resistances = sorted((c for c in clusters if c["price"] > price), key=lambda c: c["price"] - price)
    return supports, resistances


def rsi_divergence(bars, pivs, rsi_vals, lookback=60):
    """Bearish: higher price high with lower RSI. Bullish: lower price low with higher RSI."""
    n = len(bars)
    recent = [p for p in pivs if p["index"] >= n - lookback and rsi_vals[p["index"]] is not None]
    highs = [p for p in recent if p["kind"] == "H"][-2:]
    lows = [p for p in recent if p["kind"] == "L"][-2:]
    found = []
    if len(highs) == 2 and highs[1]["price"] > highs[0]["price"] and rsi_vals[highs[1]["index"]] < rsi_vals[highs[0]["index"]]:
        found.append(f"bearish: price high {highs[0]['price']:.2f}->{highs[1]['price']:.2f} while RSI "
                     f"{rsi_vals[highs[0]['index']]:.1f}->{rsi_vals[highs[1]['index']]:.1f} ({highs[1]['date']})")
    if len(lows) == 2 and lows[1]["price"] < lows[0]["price"] and rsi_vals[lows[1]["index"]] > rsi_vals[lows[0]["index"]]:
        found.append(f"bullish: price low {lows[0]['price']:.2f}->{lows[1]['price']:.2f} while RSI "
                     f"{rsi_vals[lows[0]['index']]:.1f}->{rsi_vals[lows[1]['index']]:.1f} ({lows[1]['date']})")
    return found


def gaps(bars, lookback=20, min_pct=2.0):
    out = []
    start = max(1, len(bars) - lookback)
    for i in range(start, len(bars)):
        pc, o = bars[i - 1]["close"], bars[i]["open"]
        pct = 100.0 * (o - pc) / pc
        if abs(pct) < min_pct:
            continue
        after = bars[i:]
        if pct > 0:
            filled = any(b["low"] <= pc for b in after)
        else:
            filled = any(b["high"] >= pc for b in after)
        out.append({"date": bars[i]["date"], "pct": pct, "from": pc, "filled": filled})
    return out


def last_cross(fast, slow, dates):
    """Most recent bar where fast crossed slow. Returns (direction, date, bars_ago) or None."""
    for i in range(len(fast) - 1, 0, -1):
        if None in (fast[i], slow[i], fast[i - 1], slow[i - 1]):
            break
        above_now, above_prev = fast[i] > slow[i], fast[i - 1] > slow[i - 1]
        if above_now != above_prev:
            return ("bullish" if above_now else "bearish", dates[i], len(fast) - 1 - i)
    return None


def pct(a, b):
    return None if (a is None or b is None or b == 0) else 100.0 * (a - b) / b


def percentile_rank(vals, x):
    clean = [v for v in vals if v is not None]
    if not clean:
        return None
    return 100.0 * sum(1 for v in clean if v <= x) / len(clean)


def obv_trend(ob, n):
    """OBV can sit at or below zero, so compare it directly instead of by percent."""
    if len(ob) <= n:
        return None
    diff = ob[-1] - ob[-1 - n]
    return "rising" if diff > 0 else "falling" if diff < 0 else "flat"


def slope_sign(vals, n):
    """Sign of change over the last n bars: 'rising', 'falling', or 'flat' (< 0.1%)."""
    if len(vals) <= n or vals[-1] is None or vals[-1 - n] is None:
        return None
    ch = pct(vals[-1], vals[-1 - n])
    return "rising" if ch > 0.1 else "falling" if ch < -0.1 else "flat"


# --------------------------------------------------------------------------- #
# Analysis
# --------------------------------------------------------------------------- #

def analyse(daily, weekly=None, sma200_daily=None, pivot_window=3):
    closes = [b["close"] for b in daily]
    dates = [b["date"] for b in daily]
    n = len(daily)
    last = daily[-1]
    price = last["close"]

    s20, s50 = sma(closes, 20), sma(closes, 50)
    e12, e26 = ema(closes, 12), ema(closes, 26)
    r = rsi(closes)
    m_line, m_sig, m_hist = macd(closes)
    bb_u, bb_m, bb_l = bbands(closes)
    a = atr(daily)
    ax, pdi, mdi = adx(daily)
    st_k, st_d = stochastic(daily)
    ob = obv(daily)
    pivs = pivots(daily, pivot_window)

    sma200_source = "daily"
    if sma200_daily is None:
        s200 = sma(closes, 200)
        sma200_daily = s200[-1]
        if sma200_daily is None and weekly:
            wk = sma([b["close"] for b in weekly], 40)[-1]
            sma200_daily, sma200_source = wk, "approx-40-week"
    if sma200_daily is None:
        sma200_source = "unavailable"

    bw = [None if (u is None or m in (None, 0)) else (u - l) / m for u, m, l in zip(bb_u, bb_m, bb_l)]
    bw_now = bw[-1]
    bw_pct = percentile_rank(bw[-120:], bw_now) if bw_now is not None else None
    pct_b = None if bb_u[-1] is None or bb_u[-1] == bb_l[-1] else (price - bb_l[-1]) / (bb_u[-1] - bb_l[-1])

    vol20 = sma([b["volume"] for b in daily], 20)[-1]
    win = daily[-20:]
    up_vol = sum(b["volume"] for i, b in enumerate(win) if i and b["close"] > win[i - 1]["close"])
    dn_vol = sum(b["volume"] for i, b in enumerate(win) if i and b["close"] < win[i - 1]["close"])

    supports, resistances = cluster_levels(pivs, price)
    struct_label, struct_detail = trend_structure(pivs[-8:])
    cross_50_200 = last_cross(s50, sma(closes, 200), dates) if sma200_source == "daily" else None
    cross_20_50 = last_cross(s20, s50, dates)
    macd_cross = last_cross(m_line, m_sig, dates)

    hi20 = max(b["high"] for b in daily[-20:])
    lo20 = min(b["low"] for b in daily[-20:])

    out = {
        "as_of": last["date"],
        "bars": n,
        "price": price,
        "change_1d_pct": pct(price, daily[-2]["close"]) if n > 1 else None,
        "returns_pct": {
            "5d": pct(price, closes[-6]) if n > 5 else None,
            "20d": pct(price, closes[-21]) if n > 20 else None,
            "60d": pct(price, closes[-61]) if n > 60 else None,
        },
        "range_20d": {"high": hi20, "low": lo20, "pct_from_high": pct(price, hi20), "pct_from_low": pct(price, lo20)},
        "ma": {
            "sma20": s20[-1], "sma50": s50[-1], "sma200": sma200_daily, "sma200_source": sma200_source,
            "ema12": e12[-1], "ema26": e26[-1],
            "pct_vs_sma20": pct(price, s20[-1]), "pct_vs_sma50": pct(price, s50[-1]), "pct_vs_sma200": pct(price, sma200_daily),
            "sma20_slope": slope_sign(s20, 5), "sma50_slope": slope_sign(s50, 10),
            "cross_20_50": cross_20_50, "cross_50_200": cross_50_200,
        },
        "momentum": {
            "rsi14": r[-1], "rsi14_prev5": r[-6] if n > 5 else None,
            "macd": m_line[-1], "macd_signal": m_sig[-1], "macd_hist": m_hist[-1],
            "macd_hist_prev": m_hist[-2] if n > 1 else None, "macd_cross": macd_cross,
            "stoch_k": st_k[-1], "stoch_d": st_d[-1],
            "adx14": ax[-1], "plus_di": pdi[-1], "minus_di": mdi[-1],
            "rsi_divergence": rsi_divergence(daily, pivs, r),
        },
        "volatility": {
            "atr14": a[-1], "atr14_pct": pct(a[-1] + price, price) if a[-1] is not None else None,
            "bb_upper": bb_u[-1], "bb_mid": bb_m[-1], "bb_lower": bb_l[-1],
            "bb_pct_b": pct_b, "bb_bandwidth": bw_now, "bb_bandwidth_percentile_120": bw_pct,
        },
        "volume": {
            "last": last["volume"], "avg20": vol20, "ratio_vs_avg20": None if not vol20 else last["volume"] / vol20,
            "up_down_ratio_20": None if dn_vol == 0 else up_vol / dn_vol,
            "obv_trend_20": obv_trend(ob, 20),
            "price_trend_20": slope_sign(closes, 20),
        },
        "structure": {"label": struct_label, "detail": struct_detail, "pivot_window": pivot_window,
                      "recent_pivots": pivs[-6:]},
        "levels": {"supports": supports[:4], "resistances": resistances[:4]},
        "gaps_20d": gaps(daily),
        "weekly": analyse_weekly(weekly) if weekly else None,
    }
    out["ledger"] = ledger(out)
    return out


def analyse_weekly(weekly):
    closes = [b["close"] for b in weekly]
    dates = [b["date"] for b in weekly]
    price = closes[-1]
    s10, s20, s50, s200 = sma(closes, 10), sma(closes, 20), sma(closes, 50), sma(closes, 200)
    r = rsi(closes)
    m_line, m_sig, m_hist = macd(closes)
    last52 = weekly[-52:]
    hi52, lo52 = max(b["high"] for b in last52), min(b["low"] for b in last52)
    hi52_date = next(b["date"] for b in last52 if b["high"] == hi52)
    lo52_date = next(b["date"] for b in last52 if b["low"] == lo52)
    pivs = pivots(weekly, 2)
    label, detail = trend_structure(pivs[-8:])
    sup, res = cluster_levels([p for p in pivs if p["index"] >= len(weekly) - 156], price, tol=0.02)
    adjusted = all(b["adjusted"] for b in weekly)
    # An unadjusted series still carries pre-split prices, so a high older than
    # the last split is fiction. Only report it when the feed was adjusted.
    all_time_high = max(b["high"] for b in weekly) if adjusted else None
    return {
        "adjusted": adjusted,
        "as_of": dates[-1], "bars": len(weekly),
        "returns_pct": {
            "13w": pct(price, closes[-14]) if len(closes) > 13 else None,
            "26w": pct(price, closes[-27]) if len(closes) > 26 else None,
            "52w": pct(price, closes[-53]) if len(closes) > 52 else None,
        },
        "range_52w": {"high": hi52, "high_date": hi52_date, "low": lo52, "low_date": lo52_date,
                      "pct_from_high": pct(price, hi52), "pct_from_low": pct(price, lo52)},
        "all_time_high": all_time_high, "pct_from_ath": pct(price, all_time_high),
        "ma": {"sma10": s10[-1], "sma20": s20[-1], "sma50": s50[-1], "sma200": s200[-1],
               "pct_vs_sma20": pct(price, s20[-1]), "pct_vs_sma50": pct(price, s50[-1]), "pct_vs_sma200": pct(price, s200[-1]),
               "sma20_slope": slope_sign(s20, 4), "sma50_slope": slope_sign(s50, 8),
               "cross_20_50": last_cross(s20, s50, dates)},
        "rsi14": r[-1], "macd_hist": m_hist[-1], "macd_cross": last_cross(m_line, m_sig, dates),
        "structure": {"label": label, "detail": detail},
        "levels": {"supports": sup[:3], "resistances": res[:3]},
    }


def ledger(f):
    """One row per indicator: (name, value, read). read in {bullish, bearish, neutral}."""
    rows = []
    ma, mo, vo, vl, wk = f["ma"], f["momentum"], f["volatility"], f["volume"], f["weekly"]
    p = f["price"]

    def add(name, value, read, rule):
        rows.append({"indicator": name, "value": value, "read": read, "rule": rule})

    if ma["sma20"] and ma["sma50"]:
        if ma["sma200"]:
            aligned_up = p > ma["sma20"] > ma["sma50"] > ma["sma200"]
            aligned_dn = p < ma["sma20"] < ma["sma50"] < ma["sma200"]
            add("MA stack (20/50/200)", f"{ma['sma20']:.2f}/{ma['sma50']:.2f}/{ma['sma200']:.2f}",
                "bullish" if aligned_up else "bearish" if aligned_dn else "neutral",
                "bullish only if price > 20 > 50 > 200; bearish only if fully inverted")
        else:
            add("MA stack (20/50)", f"{ma['sma20']:.2f}/{ma['sma50']:.2f}",
                "bullish" if p > ma["sma20"] > ma["sma50"] else "bearish" if p < ma["sma20"] < ma["sma50"] else "neutral",
                "no 200-day available")
    if ma["sma200"]:
        add("Price vs 200-day", f"{ma['pct_vs_sma200']:+.1f}%", "bullish" if p > ma["sma200"] else "bearish", "above/below")
    if mo["adx14"] is not None:
        strong = mo["adx14"] >= 25
        direction = "bullish" if mo["plus_di"] > mo["minus_di"] else "bearish"
        add("ADX14 / DI", f"{mo['adx14']:.1f} (+DI {mo['plus_di']:.1f} / -DI {mo['minus_di']:.1f})",
            direction if strong else "neutral", "ADX >= 25 trending, direction from DI; below 25 no trend read")
    if mo["rsi14"] is not None:
        v = mo["rsi14"]
        add("RSI14", f"{v:.1f}", "bearish" if v >= 70 else "bullish" if v <= 30 else "neutral",
            ">= 70 overbought (bearish), <= 30 oversold (bullish); in-between is no signal")
    if mo["macd_hist"] is not None:
        h, hp = mo["macd_hist"], mo["macd_hist_prev"]
        read = "bullish" if h > 0 else "bearish"
        add("MACD histogram", f"{h:+.3f} (prev {hp:+.3f})" if hp is not None else f"{h:+.3f}", read,
            "sign of histogram; watch whether it is expanding or contracting")
    if mo["stoch_k"] is not None:
        k = mo["stoch_k"]
        add("Stochastic %K", f"{k:.1f}", "bearish" if k >= 80 else "bullish" if k <= 20 else "neutral",
            ">= 80 overbought, <= 20 oversold")
    if vo["bb_pct_b"] is not None:
        b = vo["bb_pct_b"]
        add("Bollinger %B", f"{b:.2f}", "bearish" if b > 1 else "bullish" if b < 0 else "neutral",
            "> 1 closed above upper band, < 0 below lower band")
    if vl["obv_trend_20"] and vl["price_trend_20"]:
        agree = vl["obv_trend_20"] == vl["price_trend_20"]
        add("OBV vs price (20 bars)", f"OBV {vl['obv_trend_20']}, price {vl['price_trend_20']}",
            ("bullish" if vl["price_trend_20"] == "rising" else "bearish") if agree and vl["price_trend_20"] != "flat" else "neutral",
            "confirms when both move the same way; disagreement is a divergence, not a signal")
    if vl["up_down_ratio_20"] is not None:
        u = vl["up_down_ratio_20"]
        add("Up/down volume (20 bars)", f"{u:.2f}", "bullish" if u >= 1.3 else "bearish" if u <= 0.77 else "neutral",
            ">= 1.3 accumulation, <= 0.77 distribution")
    add("Swing structure (daily)", f["structure"]["label"],
        {"uptrend": "bullish", "downtrend": "bearish"}.get(f["structure"]["label"], "neutral"), "HH/HL vs LH/LL of last swings")
    for d in mo["rsi_divergence"]:
        add("RSI divergence", d, "bearish" if d.startswith("bearish") else "bullish", "two most recent swing points")
    if wk:
        add("Weekly structure", wk["structure"]["label"],
            {"uptrend": "bullish", "downtrend": "bearish"}.get(wk["structure"]["label"], "neutral"), "HH/HL vs LH/LL of weekly swings")
        if wk["ma"]["sma50"] and wk["ma"]["sma200"]:
            add("Weekly price vs 50/200-week", f"{wk['ma']['pct_vs_sma50']:+.1f}% / {wk['ma']['pct_vs_sma200']:+.1f}%",
                "bullish" if p > wk["ma"]["sma50"] > wk["ma"]["sma200"] else "bearish" if p < wk["ma"]["sma50"] < wk["ma"]["sma200"] else "neutral",
                "bullish only if price > 50w > 200w")
    tally = {k: sum(1 for r in rows if r["read"] == k) for k in ("bullish", "bearish", "neutral")}
    return {"rows": rows, "tally": tally}


# --------------------------------------------------------------------------- #
# Output
# --------------------------------------------------------------------------- #

def _f(v, nd=2, suffix=""):
    if v is None:
        return "n/a"
    return f"{v:,.{nd}f}{suffix}"


def _p(v):
    return "n/a" if v is None else f"{v:+.1f}%"


def _cross(c):
    if not c:
        return "none in window"
    d, date, ago = c
    return f"{d} on {date} ({ago} bars ago)"


def render(f, symbol):
    L = []
    ma, mo, vo, vl, wk = f["ma"], f["momentum"], f["volatility"], f["volume"], f["weekly"]
    L.append(f"# TA facts: {symbol} (daily as of {f['as_of']}, {f['bars']} bars"
             + (f"; weekly as of {wk['as_of']}, {wk['bars']} bars" if wk else "") + ")")
    L.append("")
    L.append("## Price")
    L.append(f"- Close {_f(f['price'])} ({_p(f['change_1d_pct'])} on the day)")
    r = f["returns_pct"]
    L.append(f"- Returns: 5d {_p(r['5d'])}, 20d {_p(r['20d'])}, 60d {_p(r['60d'])}"
             + (f", 13w {_p(wk['returns_pct']['13w'])}, 26w {_p(wk['returns_pct']['26w'])}, 52w {_p(wk['returns_pct']['52w'])}" if wk else ""))
    rg = f["range_20d"]
    L.append(f"- 20-day range {_f(rg['low'])} - {_f(rg['high'])} (from high {_p(rg['pct_from_high'])}, from low {_p(rg['pct_from_low'])})")
    if wk:
        w52 = wk["range_52w"]
        L.append(f"- 52-week range {_f(w52['low'])} ({w52['low_date']}) - {_f(w52['high'])} ({w52['high_date']}); "
                 f"from high {_p(w52['pct_from_high'])}, from low {_p(w52['pct_from_low'])}; from all-time high "
                 + (_p(wk['pct_from_ath']) if wk['adjusted'] else "n/a (unadjusted feed; use TIME_SERIES_WEEKLY_ADJUSTED)"))
    L.append("")
    L.append("## Trend (daily)")
    L.append("| MA | Value | Price vs MA | Slope |")
    L.append("|---|---|---|---|")
    L.append(f"| SMA20 | {_f(ma['sma20'])} | {_p(ma['pct_vs_sma20'])} | {ma['sma20_slope'] or 'n/a'} |")
    L.append(f"| SMA50 | {_f(ma['sma50'])} | {_p(ma['pct_vs_sma50'])} | {ma['sma50_slope'] or 'n/a'} |")
    L.append(f"| SMA200 ({ma['sma200_source']}) | {_f(ma['sma200'])} | {_p(ma['pct_vs_sma200'])} | n/a |")
    L.append(f"- 20/50 cross: {_cross(ma['cross_20_50'])}; 50/200 cross: {_cross(ma['cross_50_200']) if ma['sma200_source'] == 'daily' else 'needs 200+ daily bars'}")
    L.append(f"- ADX14 {_f(mo['adx14'], 1)} (+DI {_f(mo['plus_di'], 1)}, -DI {_f(mo['minus_di'], 1)})")
    L.append(f"- Swing structure: {f['structure']['label']} ({f['structure']['detail']})")
    L.append("")
    L.append("## Momentum")
    L.append(f"- RSI14 {_f(mo['rsi14'], 1)} (5 bars ago {_f(mo['rsi14_prev5'], 1)})")
    L.append(f"- MACD {_f(mo['macd'], 3)} / signal {_f(mo['macd_signal'], 3)} / hist {_f(mo['macd_hist'], 3)} (prev {_f(mo['macd_hist_prev'], 3)}); last cross: {_cross(mo['macd_cross'])}")
    L.append(f"- Stochastic %K {_f(mo['stoch_k'], 1)} / %D {_f(mo['stoch_d'], 1)}")
    L.append("- RSI divergence: " + ("; ".join(mo["rsi_divergence"]) if mo["rsi_divergence"] else "none between the last two swings"))
    L.append("")
    L.append("## Volatility")
    L.append(f"- ATR14 {_f(vo['atr14'])} ({_f(vo['atr14_pct'], 2, '%')} of price)")
    L.append(f"- Bollinger(20,2): upper {_f(vo['bb_upper'])} / mid {_f(vo['bb_mid'])} / lower {_f(vo['bb_lower'])}; %B {_f(vo['bb_pct_b'])}; "
             f"bandwidth {_f(vo['bb_bandwidth'], 3)} (percentile of last 120 bars: {_f(vo['bb_bandwidth_percentile_120'], 0)})")
    L.append("")
    L.append("## Volume")
    L.append(f"- Last {_f(vl['last'], 0)} vs 20-day avg {_f(vl['avg20'], 0)} (x{_f(vl['ratio_vs_avg20'])})")
    L.append(f"- Up/down volume ratio (20 bars) {_f(vl['up_down_ratio_20'])}; OBV {vl['obv_trend_20'] or 'n/a'} while price {vl['price_trend_20'] or 'n/a'}")
    L.append("")
    L.append("## Levels (daily swing clusters, nearest first)")
    for kind in ("supports", "resistances"):
        items = f["levels"][kind]
        L.append(f"- {kind.capitalize()}: " + (", ".join(f"{_f(c['price'])} ({c['touches']} touch{'es' if c['touches'] > 1 else ''}, last {c['last']})" for c in items) if items else "none in window"))
    if wk:
        for kind in ("supports", "resistances"):
            items = wk["levels"][kind]
            L.append(f"- Weekly {kind} (3y): " + (", ".join(f"{_f(c['price'])} ({c['touches']}, last {c['last']})" for c in items) if items else "none"))
    L.append(f"- Gaps (20d, >=2%): " + ("; ".join(f"{g['date']} {g['pct']:+.1f}% from {_f(g['from'])} ({'filled' if g['filled'] else 'open'})" for g in f["gaps_20d"]) if f["gaps_20d"] else "none"))
    L.append("")
    if wk:
        L.append("## Weekly")
        wm = wk["ma"]
        L.append("| MA | Value | Price vs MA | Slope |")
        L.append("|---|---|---|---|")
        L.append(f"| 10-week | {_f(wm['sma10'])} | {_p(pct(f['price'], wm['sma10']))} | n/a |")
        L.append(f"| 20-week | {_f(wm['sma20'])} | {_p(wm['pct_vs_sma20'])} | {wm['sma20_slope'] or 'n/a'} |")
        L.append(f"| 50-week | {_f(wm['sma50'])} | {_p(wm['pct_vs_sma50'])} | {wm['sma50_slope'] or 'n/a'} |")
        L.append(f"| 200-week | {_f(wm['sma200'])} | {_p(wm['pct_vs_sma200'])} | n/a |")
        L.append(f"- 20/50-week cross: {_cross(wm['cross_20_50'])}; weekly RSI14 {_f(wk['rsi14'], 1)}; weekly MACD hist {_f(wk['macd_hist'], 3)}; weekly MACD cross: {_cross(wk['macd_cross'])}")
        L.append(f"- Weekly swing structure: {wk['structure']['label']} ({wk['structure']['detail']})")
        L.append("")
    L.append("## Signal ledger")
    L.append("| Indicator | Value | Read | Rule |")
    L.append("|---|---|---|---|")
    for row in f["ledger"]["rows"]:
        L.append(f"| {row['indicator']} | {row['value']} | {row['read']} | {row['rule']} |")
    t = f["ledger"]["tally"]
    L.append(f"- Tally: {t['bullish']} bullish / {t['bearish']} bearish / {t['neutral']} neutral. A tally counts rows; it does not weight them.")
    return "\n".join(L)


# --------------------------------------------------------------------------- #
# Optional direct fetch (no MCP): needs ALPHAVANTAGE_API_KEY
# --------------------------------------------------------------------------- #

AV_URL = "https://www.alphavantage.co/query"


def fetch(symbol, out_dir):
    key = os.environ.get("ALPHAVANTAGE_API_KEY")
    if not key:
        sys.exit("ALPHAVANTAGE_API_KEY is not set; use the MCP tools instead")
    os.makedirs(out_dir, exist_ok=True)
    paths = {}
    for name, params in (
        ("daily", {"function": "TIME_SERIES_DAILY", "outputsize": "compact"}),
        ("weekly", {"function": "TIME_SERIES_WEEKLY_ADJUSTED"}),
        ("sma200", {"function": "SMA", "interval": "daily", "time_period": 200, "series_type": "close"}),
    ):
        q = dict(params, symbol=symbol, apikey=key, datatype="csv")
        with urllib.request.urlopen(AV_URL + "?" + urllib.parse.urlencode(q), timeout=30) as resp:
            body = resp.read().decode("utf-8")
        if body.lstrip().startswith("{"):
            sys.exit(f"{name}: {body.strip()[:300]}")
        p = os.path.join(out_dir, f"{symbol.upper()}_{name}.csv")
        with open(p, "w", encoding="utf-8") as fh:
            fh.write(body)
        paths[name] = p
    return paths


def latest_sma_from_csv(path):
    """The SMA endpoint's CSV is 'time,SMA' newest first; return the newest value."""
    with open(path, encoding="utf-8") as fh:
        rows = list(csv.DictReader(fh))
    return float(rows[0]["SMA"]) if rows else None


def main(argv=None):
    ap = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    ap.add_argument("--daily", help="daily OHLCV file (CSV, {result: csv}, or AV JSON)")
    ap.add_argument("--weekly", help="weekly OHLCV file (same shapes)")
    ap.add_argument("--sma200", type=float, help="latest daily 200-SMA from the SMA endpoint")
    ap.add_argument("--symbol", default="?", help="ticker, for the heading only")
    ap.add_argument("--json", action="store_true", help="emit the full JSON instead of markdown")
    ap.add_argument("--pivot-window", type=int, default=3, help="bars either side for a swing point (default 3)")
    ap.add_argument("--fetch", metavar="SYM", help="fetch daily/weekly/SMA200 with ALPHAVANTAGE_API_KEY, then analyse")
    ap.add_argument("--out-dir", default=".", help="where --fetch writes its CSV files")
    args = ap.parse_args(argv)

    if args.fetch:
        paths = fetch(args.fetch, args.out_dir)
        args.daily, args.weekly, args.symbol = paths["daily"], paths["weekly"], args.fetch.upper()
        args.sma200 = latest_sma_from_csv(paths["sma200"])
    if not args.daily:
        ap.error("--daily is required (or --fetch)")

    daily = load_bars(args.daily)
    weekly = load_bars(args.weekly) if args.weekly else None
    facts = analyse(daily, weekly, args.sma200, args.pivot_window)
    if args.json:
        print(json.dumps(facts, indent=2, default=str))
    else:
        print(render(facts, args.symbol.upper()))


if __name__ == "__main__":
    main()
