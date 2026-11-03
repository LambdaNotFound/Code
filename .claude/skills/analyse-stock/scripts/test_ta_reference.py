"""Differential tests: every number ta.py reports, re-computed by code written apart from it.

Each reference below follows the textbook definition named in its docstring
and shares no code with ta.py. The tests run both on five seeded synthetic
markets (market.py: an uptrend, a downtrend, a range, a volatile walk and a
gappy one) and compare field by field. Floats match at 1e-9 relative unless a
test says otherwise.

Run from this directory: python3 -m unittest test_ta_reference -v
"""
import math
import os
import statistics
import sys
import unittest

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
import ta  # noqa: E402
from market import REGIMES, market  # noqa: E402

TOL = 1e-9
# One market per regime, plus seeds picked because they reach the rarer branches:
# a BUY, a SELL, a weekly key reversal, a vetoed reversal, an unfilled gap, an exit
# moved up to resistance (and seed 11's uptrend: an RSI divergence and an entry
# moved down to support).
CASES = [(seed, regime) for seed, regime in enumerate(REGIMES, start=11)] + [
    (23, "range"), (17, "downtrend"), (24, "uptrend"), (13, "uptrend"), (20, "volatile"), (11, "downtrend")]


# --------------------------------------------------------------------------- #
# Reference implementations
# --------------------------------------------------------------------------- #

def r_sma_series(v, n):
    """Simple moving average, None until n values exist."""
    return [None] * (n - 1) + [statistics.fmean(v[i - n + 1:i + 1]) for i in range(n - 1, len(v))]


def r_ema_series(v, n):
    """EMA with alpha 2/(n+1), seeded with the mean of the first n values."""
    if len(v) < n:
        return [None] * len(v)
    a, out = 2.0 / (n + 1), [None] * (n - 1)
    e = statistics.fmean(v[:n])
    out.append(e)
    for x in v[n:]:
        e = a * x + (1 - a) * e
        out.append(e)
    return out


def r_wilder_series(v, n):
    """Wilder's running average: seed with the mean of the first n, then (prev*(n-1) + x)/n."""
    if len(v) < n:
        return [None] * len(v)
    out = [None] * (n - 1)
    s = statistics.fmean(v[:n])
    out.append(s)
    for x in v[n:]:
        s = (s * (n - 1) + x) / n
        out.append(s)
    return out


def r_rsi_at(closes, end, n=14):
    """Wilder RSI on closes[:end+1]."""
    ch = [b - a for a, b in zip(closes[:end + 1], closes[1:end + 1])]
    ag = r_wilder_series([max(x, 0.0) for x in ch], n)[-1]
    al = r_wilder_series([max(-x, 0.0) for x in ch], n)[-1]
    if ag == 0 and al == 0:
        return 50.0  # no movement: no signal
    return 100.0 if al == 0 else 100.0 - 100.0 / (1.0 + ag / al)


def r_macd(closes):
    """MACD(12, 26, 9): EMA12 - EMA26; signal = EMA9 of the line from its first value."""
    e12, e26 = r_ema_series(closes, 12), r_ema_series(closes, 26)
    line = [a - b for a, b in zip(e12[25:], e26[25:])]
    sig = r_ema_series(line, 9)
    hist = [None if s is None else l_ - s for l_, s in zip(line, sig)]
    return line, sig, hist


def r_true_range(bars):
    """True range; the first bar has no prior close, so its range is high - low."""
    out = [bars[0]["high"] - bars[0]["low"]]
    for p, b in zip(bars, bars[1:]):
        out.append(max(b["high"] - b["low"], abs(b["high"] - p["close"]), abs(b["low"] - p["close"])))
    return out


def r_adx(bars, n=14):
    """Wilder's +DI, -DI and ADX: smoothed directional movement over smoothed true range."""
    tr, pdm, mdm = [], [], []
    for p, b in zip(bars, bars[1:]):
        up, dn = b["high"] - p["high"], p["low"] - b["low"]
        tr.append(max(b["high"] - b["low"], abs(b["high"] - p["close"]), abs(b["low"] - p["close"])))
        pdm.append(up if up > dn and up > 0 else 0.0)
        mdm.append(dn if dn > up and dn > 0 else 0.0)
    s_tr, s_p, s_m = r_wilder_series(tr, n), r_wilder_series(pdm, n), r_wilder_series(mdm, n)
    pdi = [None if t is None else 100 * a / t for a, t in zip(s_p, s_tr)]
    mdi = [None if t is None else 100 * a / t for a, t in zip(s_m, s_tr)]
    dx = [100 * abs(a - b) / (a + b) if a + b else 0.0 for a, b in zip(pdi[n - 1:], mdi[n - 1:])]
    return r_wilder_series(dx, n)[-1], pdi[-1], mdi[-1]


def r_stochastic(bars, k_n=14, smooth=3):
    """Slow stochastic: raw %K over k_n bars, %K = 3-bar mean of raw, %D = 3-bar mean of %K."""
    raw = []
    for i in range(k_n - 1, len(bars)):
        w = bars[i - k_n + 1:i + 1]
        hh, ll = max(b["high"] for b in w), min(b["low"] for b in w)
        raw.append(50.0 if hh == ll else 100 * (bars[i]["close"] - ll) / (hh - ll))
    k = [statistics.fmean(raw[i - smooth + 1:i + 1]) for i in range(smooth - 1, len(raw))]
    d = [statistics.fmean(k[i - smooth + 1:i + 1]) for i in range(smooth - 1, len(k))]
    return k[-1], d[-1]


def r_slope(series, n):
    """'rising' / 'falling' over n bars when the change exceeds 0.1%, else 'flat'."""
    a, b = series[-1 - n], series[-1]
    if a is None or b is None:
        return None
    ch = 100 * (b - a) / a
    return "rising" if ch > 0.1 else "falling" if ch < -0.1 else "flat"


def r_last_cross(fast, slow, dates):
    """The latest bar where fast moved to the other side of slow."""
    for i in range(len(fast) - 1, 0, -1):
        pair = (fast[i], slow[i], fast[i - 1], slow[i - 1])
        if any(x is None for x in pair):
            return None
        if (fast[i] > slow[i]) != (fast[i - 1] > slow[i - 1]):
            return ("bullish" if fast[i] > slow[i] else "bearish", dates[i], len(fast) - 1 - i)
    return None


def r_pivots(bars, k):
    """Fractal swings: a high above the k bars on each side; a low below them."""
    out = []
    for i in range(k, len(bars) - k):
        side = bars[i - k:i] + bars[i + 1:i + k + 1]
        if all(bars[i]["high"] > b["high"] for b in side):
            out.append((i, "H", bars[i]["high"]))
        if all(bars[i]["low"] < b["low"] for b in side):
            out.append((i, "L", bars[i]["low"]))
    return out


def r_structure(pivs):
    """Uptrend on a higher high and a higher low; downtrend on lower of both; else mixed."""
    hs = [p[2] for p in pivs if p[1] == "H"][-2:]
    ls = [p[2] for p in pivs if p[1] == "L"][-2:]
    if len(hs) < 2 or len(ls) < 2:
        return "undetermined"
    return {(True, True): "uptrend", (False, False): "downtrend"}.get((hs[1] > hs[0], ls[1] > ls[0]), "mixed")


def r_gaps(bars, lookback=20, min_pct=2.0):
    """Opens at least min_pct from the prior close in the last lookback bars, and whether
    a later bar traded back to that close."""
    out = []
    for i in range(max(1, len(bars) - lookback), len(bars)):
        pc = bars[i - 1]["close"]
        g = 100 * (bars[i]["open"] - pc) / pc
        if abs(g) >= min_pct:
            later = bars[i:]
            filled = any(b["low"] <= pc for b in later) if g > 0 else any(b["high"] >= pc for b in later)
            out.append((bars[i]["date"], filled))
    return out


def r_bursts(bars, lookback=20, min_pct=4.0):
    """Close up or down min_pct on rising volume, or a range above each of the prior three
    when the prior day moved less than min_pct."""
    out = []
    for i in range(max(3, len(bars) - lookback), len(bars)):
        b, p, pp = bars[i], bars[i - 1], bars[i - 2]
        chg, prev_chg = 100 * (b["close"] / p["close"] - 1), 100 * (p["close"] / pp["close"] - 1)
        tags = []
        if chg >= min_pct and b["volume"] > p["volume"]:
            tags.append("4pct_breakout")
        if chg <= -min_pct and b["volume"] > p["volume"]:
            tags.append("4pct_breakdown")
        if b["high"] - b["low"] > max(x["high"] - x["low"] for x in bars[i - 3:i]) and abs(prev_chg) < min_pct:
            tags.append("range_expansion")
        if tags:
            out.append((b["date"], tags))
    return out


NICE = (0.01, 0.02, 0.05, 0.1, 0.2, 0.25, 0.5, 1, 2, 2.5, 5, 10, 20, 25, 50, 100, 200, 250, 500)


def r_cost(bars, price):
    """Volume at price: each bar's volume spread evenly over its low-high range, in bins of
    the largest round step at or under 1% of price. Exact overlaps, no sub-sampling."""
    bars = [b for b in bars if b["volume"] > 0]
    step = max([s for s in NICE if s <= price / 100] or [NICE[0]])
    base = math.floor(min(b["low"] for b in bars) / step) * step
    bins, total, profit, wsum = {}, 0.0, 0.0, 0.0
    for b in bars:
        lo, hi, v = b["low"], b["high"], b["volume"]
        total += v
        wsum += v * (lo + hi) / 2
        if hi <= lo:
            j = int((lo - base) // step)
            bins[j] = bins.get(j, 0.0) + v
            profit += v if lo <= price else 0.0
            continue
        profit += v * min(1.0, max(0.0, (price - lo) / (hi - lo)))
        for j in range(int((lo - base) // step), int((hi - base) // step) + 1):
            a, z = max(lo, base + j * step), min(hi, base + (j + 1) * step)
            if z > a:
                bins[j] = bins.get(j, 0.0) + v * (z - a) / (hi - lo)
    shares = {round(base + j * step, 6): bins[j] / total for j in sorted(bins)}
    # exact quantile of the piecewise-uniform distribution
    pieces = sorted((b["low"], b["high"], b["volume"]) for b in bars)

    def cdf(x):
        return sum(v * (1.0 if x >= hi else 0.0 if x <= lo else (x - lo) / (hi - lo)) if hi > lo else (v if x >= lo else 0.0)
                   for lo, hi, v in pieces) / total

    def q(p):
        a, z = min(b["low"] for b in bars), max(b["high"] for b in bars)
        for _ in range(80):
            m = (a + z) / 2
            a, z = (m, z) if cdf(m) < p else (a, m)
        return (a + z) / 2

    return {"bin_size": step, "shares": shares, "profit_ratio": profit / total, "average_cost": wsum / total,
            "q": {p: q(p) for p in (0.05, 0.15, 0.85, 0.95)}}


def r_weekly_reversals(weekly, ext=52, recent=8, fail=3):
    """Key reversal, failed extreme and failed breakout at a 52-week extreme, on completed
    weeks (the last bar is the week in progress), with the continuation veto."""
    bars = weekly[:-1]
    n = len(bars)
    if n < ext + recent + 1:
        return None
    H, L, C = [b["high"] for b in bars], [b["low"] for b in bars], [b["close"] for b in bars]
    bear, bull = [], []
    for i in range(n - recent, n):
        ph, pl = max(H[i - ext:i]), min(L[i - ext:i])
        if H[i] > ph:
            if C[i] < L[i - 1]:
                bear.append(("key_reversal", bars[i]["date"]))
            elif C[i] < ph:
                bear.append(("failed_extreme", bars[i]["date"]))
        if L[i] < pl:
            if C[i] > H[i - 1]:
                bull.append(("key_reversal", bars[i]["date"]))
            elif C[i] > pl:
                bull.append(("failed_extreme", bars[i]["date"]))
    for b in range(n - recent - fail, n - 1):
        ph, pl = max(H[b - ext:b]), min(L[b - ext:b])
        for out, broke, back in ((bear, C[b] > ph, lambda f: C[f] < ph), (bull, C[b] < pl, lambda f: C[f] > pl)):
            if broke:
                for f in range(b + 1, min(b + 1 + fail, n)):
                    if back(f) and f >= n - recent:
                        out.append(("failed_breakout", bars[f]["date"]))
                        break

    def vetoed(sig, up):
        if not sig:
            return False
        idx = [b["date"] for b in bars].index(max(s[1] for s in sig))
        return any((C[i] > max(C[i - ext:i])) if up else (C[i] < min(C[i - ext:i])) for i in range(idx + 1, n))

    bv, uv = vetoed(bear, True), vetoed(bull, False)
    live_bear, live_bull = bool(bear) and not bv, bool(bull) and not uv
    read = "bearish" if live_bear and not live_bull else "bullish" if live_bull and not live_bear else "neutral"
    return {"bearish": sorted(bear), "bullish": sorted(bull), "read": read, "bearish_vetoed": bv, "bullish_vetoed": uv}


# --------------------------------------------------------------------------- #
# Tests
# --------------------------------------------------------------------------- #

class Base(unittest.TestCase):
    def close(self, got, want, what, tol=TOL):
        self.assertIsNotNone(got, what)
        self.assertTrue(abs(got - want) <= tol * max(1.0, abs(want)), f"{what}: ta.py {got}, reference {want}")


class DailyIndicators(Base):
    def test_every_daily_number(self):
        for seed, regime in CASES:
            with self.subTest(regime=regime):
                daily, weekly = market(seed, regime)
                f = ta.analyse(daily, weekly)
                c = [b["close"] for b in daily]
                p = c[-1]
                dates = [b["date"] for b in daily]
                self.close(f["price"], p, "price")
                self.close(f["change_1d_pct"], 100 * (p / c[-2] - 1), "change_1d_pct")
                for k, n in (("5d", 5), ("20d", 20), ("60d", 60)):
                    self.close(f["returns_pct"][k], 100 * (p / c[-1 - n] - 1), f"returns {k}")
                hi20, lo20 = max(b["high"] for b in daily[-20:]), min(b["low"] for b in daily[-20:])
                self.close(f["range_20d"]["high"], hi20, "range high")
                self.close(f["range_20d"]["low"], lo20, "range low")
                self.close(f["range_20d"]["pct_from_high"], 100 * (p / hi20 - 1), "pct from high")
                s20, s50, s200 = r_sma_series(c, 20), r_sma_series(c, 50), r_sma_series(c, 200)
                ma = f["ma"]
                self.close(ma["sma20"], s20[-1], "sma20")
                self.close(ma["sma50"], s50[-1], "sma50")
                self.close(ma["sma200"], s200[-1], "sma200 from 260 daily bars")
                self.assertEqual(ma["sma200_source"], "daily")
                self.close(ma["ema12"], r_ema_series(c, 12)[-1], "ema12")
                self.close(ma["ema26"], r_ema_series(c, 26)[-1], "ema26")
                self.close(ma["pct_vs_sma50"], 100 * (p / s50[-1] - 1), "pct vs sma50")
                self.assertEqual(ma["sma20_slope"], r_slope(s20, 5))
                self.assertEqual(ma["sma50_slope"], r_slope(s50, 10))
                self.assertEqual(ma["cross_20_50"], r_last_cross(s20, s50, dates))
                self.assertEqual(ma["cross_50_200"], r_last_cross(s50, s200, dates))
                mo = f["momentum"]
                self.close(mo["rsi14"], r_rsi_at(c, len(c) - 1), "rsi14")
                self.close(mo["rsi14_prev5"], r_rsi_at(c, len(c) - 6), "rsi14 five bars back")
                line, sig, hist = r_macd(c)
                self.close(mo["macd"], line[-1], "macd")
                self.close(mo["macd_signal"], sig[-1], "macd signal")
                self.close(mo["macd_hist"], hist[-1], "macd hist")
                self.close(mo["macd_hist_prev"], hist[-2], "macd hist prev")
                k, d = r_stochastic(daily)
                self.close(mo["stoch_k"], k, "stoch %K")
                self.close(mo["stoch_d"], d, "stoch %D")
                adx, pdi, mdi = r_adx(daily)
                self.close(mo["adx14"], adx, "adx14")
                self.close(mo["plus_di"], pdi, "+DI")
                self.close(mo["minus_di"], mdi, "-DI")
                vo = f["volatility"]
                atr = r_wilder_series(r_true_range(daily), 14)[-1]
                self.close(vo["atr14"], atr, "atr14")
                self.close(vo["atr14_pct"], 100 * atr / p, "atr14 %")
                mid, sd = statistics.fmean(c[-20:]), statistics.pstdev(c[-20:])
                up, lo = mid + 2 * sd, mid - 2 * sd
                self.close(vo["bb_mid"], mid, "bb mid")
                self.close(vo["bb_upper"], up, "bb upper")
                self.close(vo["bb_lower"], lo, "bb lower")
                self.close(vo["bb_pct_b"], (p - lo) / (up - lo), "%B")
                self.close(vo["bb_bandwidth"], (up - lo) / mid, "bandwidth")
                bw = [(4 * statistics.pstdev(c[i - 19:i + 1])) / statistics.fmean(c[i - 19:i + 1]) for i in range(19, len(c))]
                self.close(vo["bb_bandwidth_percentile_120"], 100 * sum(x <= bw[-1] for x in bw[-120:]) / len(bw[-120:]),
                           "bandwidth percentile", tol=1e-6)
                vols = [b["volume"] for b in daily]
                self.close(f["volume"]["avg20"], statistics.fmean(vols[-20:]), "avg volume")
                self.close(f["volume"]["ratio_vs_avg20"], vols[-1] / statistics.fmean(vols[-20:]), "volume ratio")
                w = daily[-20:]
                upv = sum(b["volume"] for a, b in zip(w, w[1:]) if b["close"] > a["close"])
                dnv = sum(b["volume"] for a, b in zip(w, w[1:]) if b["close"] < a["close"])
                self.close(f["volume"]["up_down_ratio_20"], upv / dnv, "up/down volume")
                obv = [0.0]
                for a, b in zip(daily, daily[1:]):
                    obv.append(obv[-1] + (b["volume"] if b["close"] > a["close"] else -b["volume"] if b["close"] < a["close"] else 0))
                diff = obv[-1] - obv[-21]
                self.assertEqual(f["volume"]["obv_trend_20"], "rising" if diff > 0 else "falling" if diff < 0 else "flat")
                self.assertEqual(f["volume"]["price_trend_20"], r_slope(c, 20))

    def test_structure_gaps_and_bursts(self):
        for seed, regime in CASES:
            with self.subTest(regime=regime):
                daily, weekly = market(seed, regime)
                f = ta.analyse(daily, weekly)
                piv = r_pivots(daily, 3)
                self.assertEqual(f["structure"]["label"], r_structure(piv[-8:]))
                self.assertEqual([(p["index"], p["kind"], p["price"]) for p in f["structure"]["recent_pivots"]], piv[-6:])
                self.assertEqual([(g["date"], g["filled"]) for g in f["gaps_20d"]], r_gaps(daily))
                self.assertEqual([(b["date"], b["tags"]) for b in f["burst_days_20"]], r_bursts(daily))


class WeeklyNumbers(Base):
    def test_every_weekly_number(self):
        for seed, regime in CASES:
            with self.subTest(regime=regime):
                daily, weekly = market(seed, regime)
                wk = ta.analyse(daily, weekly)["weekly"]
                c = [b["close"] for b in weekly]
                p = c[-1]
                for k, n in (("13w", 13), ("26w", 26), ("52w", 52)):
                    self.close(wk["returns_pct"][k], 100 * (p / c[-1 - n] - 1), f"weekly return {k}")
                last52 = weekly[-52:]
                hi = max(last52, key=lambda b: b["high"])
                lo = min(last52, key=lambda b: b["low"])
                self.assertEqual((wk["range_52w"]["high"], wk["range_52w"]["high_date"]), (hi["high"], hi["date"]))
                self.assertEqual((wk["range_52w"]["low"], wk["range_52w"]["low_date"]), (lo["low"], lo["date"]))
                self.close(wk["all_time_high"], max(b["high"] for b in weekly), "all-time high")
                for n in (10, 20, 50, 200):
                    self.close(wk["ma"][f"sma{n}"], statistics.fmean(c[-n:]), f"weekly sma{n}")
                s20, s50 = r_sma_series(c, 20), r_sma_series(c, 50)
                self.assertEqual(wk["ma"]["sma20_slope"], r_slope(s20, 4))
                self.assertEqual(wk["ma"]["sma50_slope"], r_slope(s50, 8))
                self.close(wk["rsi14"], r_rsi_at(c, len(c) - 1), "weekly rsi14")
                self.close(wk["macd_hist"], r_macd(c)[2][-1], "weekly macd hist")
                self.assertEqual(wk["structure"]["label"], r_structure(r_pivots(weekly, 2)[-8:]))

    def test_weekly_reversals_match_the_rules(self):
        for seed, regime in CASES:
            with self.subTest(regime=regime):
                daily, weekly = market(seed, regime)
                got = ta.analyse(daily, weekly)["weekly_reversals"]
                want = r_weekly_reversals(weekly)
                self.assertEqual(sorted((x["check"], x["week"]) for x in got["bearish"]), want["bearish"])
                self.assertEqual(sorted((x["check"], x["week"]) for x in got["bullish"]), want["bullish"])
                self.assertEqual((got["read"], got["bearish_vetoed"], got["bullish_vetoed"]),
                                 (want["read"], want["bearish_vetoed"], want["bullish_vetoed"]))
                self.assertEqual(got["completed_week"], weekly[-2]["date"])  # the last bar is in progress


class WeeklyManySeries(Base):
    """Rules whose output is a label or a list change rarely between nearby settings,
    so these run on 60 weekly histories (12 seeds x 5 regimes) rather than a handful."""

    SERIES = [market(seed, regime, daily_n=60)[1] for seed in range(11, 23) for regime in REGIMES]

    def test_weekly_labels_and_windows(self):
        for i, weekly in enumerate(self.SERIES):
            with self.subTest(series=i):
                wk = ta.analyse_weekly(weekly)
                c = [b["close"] for b in weekly]
                s20, s50 = r_sma_series(c, 20), r_sma_series(c, 50)
                self.assertEqual(wk["ma"]["sma20_slope"], r_slope(s20, 4))
                self.assertEqual(wk["ma"]["sma50_slope"], r_slope(s50, 8))
                self.assertEqual(wk["ma"]["cross_20_50"], r_last_cross(s20, s50, [b["date"] for b in weekly]))
                self.assertEqual(wk["structure"]["label"], r_structure(r_pivots(weekly, 2)[-8:]))
                got, want = ta.weekly_reversals(weekly), r_weekly_reversals(weekly)
                self.assertEqual((got["extreme_weeks"], got["window_weeks"]), (52, 8))
                self.assertEqual(sorted((x["check"], x["week"]) for x in got["bearish"]), want["bearish"])
                self.assertEqual(sorted((x["check"], x["week"]) for x in got["bullish"]), want["bullish"])
                self.assertEqual(got["read"], want["read"])


class TrendTemplateAndCost(Base):
    def test_trend_template_checks(self):
        for seed, regime in CASES:
            with self.subTest(regime=regime):
                daily, weekly = market(seed, regime)
                f = ta.analyse(daily, weekly)
                c, wc = [b["close"] for b in daily], [b["close"] for b in weekly]
                p, s50, s200 = c[-1], statistics.fmean(c[-50:]), statistics.fmean(c[-200:])
                s150 = statistics.fmean(wc[-30:])                     # the 30-week stand-in
                s200_then = statistics.fmean(wc[-44:-4])              # the 40-week SMA four weeks back
                lo52 = min(b["low"] for b in weekly[-52:])
                hi52 = max(b["high"] for b in weekly[-52:])
                want = [p > s150 and p > s200, s150 > s200, s200 > s200_then, s50 > s150 > s200, p > s50,
                        p >= 1.3 * lo52, p >= 0.75 * hi52]
                tt = f["trend_template"]
                self.assertEqual([x["pass"] for x in tt["checks"]], want)
                self.assertEqual(tt["passed"], sum(want))

    def test_cost_distribution_against_exact_overlaps(self):
        for seed, regime in CASES:
            with self.subTest(regime=regime):
                daily, weekly = market(seed, regime)
                f = ta.analyse(daily, weekly)
                for key, bars in (("daily", daily), ("weekly_52w", weekly[-52:])):
                    got, want = f["cost_distribution"][key], r_cost(bars, daily[-1]["close"])
                    self.assertEqual(got["bin_size"], want["bin_size"])
                    shares = {round(b["low"], 6): b["share"] for b in got["bins"]}
                    self.assertEqual(sorted(shares), sorted(want["shares"]), f"{key}: bin edges")
                    for edge, s in want["shares"].items():
                        self.close(shares[edge], s, f"{key} share at {edge}", tol=1e-9)
                    self.close(got["profit_ratio"], want["profit_ratio"], f"{key} profit ratio")
                    self.close(got["average_cost"], want["average_cost"], f"{key} average cost")
                    self.assertEqual(got["peak"]["share"], max(b["share"] for b in got["bins"]))
                    fine = got["bin_size"] / 10  # ta.py interpolates inside tenth-of-a-bin steps
                    for name, lo_q, hi_q in (("range_70", 0.15, 0.85), ("range_90", 0.05, 0.95)):
                        self.assertLessEqual(abs(got[name]["low"] - want["q"][lo_q]), fine + 1e-9, f"{key} {name} low")
                        self.assertLessEqual(abs(got[name]["high"] - want["q"][hi_q]), fine + 1e-9, f"{key} {name} high")


class LedgerAndStance(Base):
    def read_of(self, f, name):
        return next((r["read"] for r in f["ledger"]["rows"] if r["indicator"] == name), None)

    def test_each_ledger_row_follows_its_rule(self):
        for seed, regime in CASES:
            with self.subTest(regime=regime):
                daily, weekly = market(seed, regime)
                f = ta.analyse(daily, weekly)
                p, ma, mo, vo, vl, wk = f["price"], f["ma"], f["momentum"], f["volatility"], f["volume"], f["weekly"]
                want = {
                    "MA stack (20/50/200)": ("bullish" if p > ma["sma20"] > ma["sma50"] > ma["sma200"]
                                             else "bearish" if p < ma["sma20"] < ma["sma50"] < ma["sma200"] else "neutral"),
                    "Price vs 200-day": "bullish" if p > ma["sma200"] else "bearish",
                    "ADX14 / DI": ("neutral" if mo["adx14"] < 25 else "bullish" if mo["plus_di"] > mo["minus_di"] else "bearish"),
                    "RSI14": "bearish" if mo["rsi14"] >= 70 else "bullish" if mo["rsi14"] <= 30 else "neutral",
                    "MACD histogram": "bullish" if mo["macd_hist"] > 0 else "bearish" if mo["macd_hist"] < 0 else "neutral",
                    "Stochastic %K": "bearish" if mo["stoch_k"] >= 80 else "bullish" if mo["stoch_k"] <= 20 else "neutral",
                    "Bollinger %B": "bearish" if vo["bb_pct_b"] > 1 else "bullish" if vo["bb_pct_b"] < 0 else "neutral",
                    "Up/down volume (20 bars)": ("bullish" if vl["up_down_ratio_20"] >= 1.3
                                                 else "bearish" if vl["up_down_ratio_20"] <= 0.77 else "neutral"),
                    "Swing structure (daily)": {"uptrend": "bullish", "downtrend": "bearish"}.get(f["structure"]["label"], "neutral"),
                    "Weekly structure": {"uptrend": "bullish", "downtrend": "bearish"}.get(wk["structure"]["label"], "neutral"),
                    "Weekly price vs 50/200-week": ("bullish" if p > wk["ma"]["sma50"] > wk["ma"]["sma200"]
                                                    else "bearish" if p < wk["ma"]["sma50"] < wk["ma"]["sma200"] else "neutral"),
                    "Trend template (Minervini)": ("bullish" if f["trend_template"]["passed"] >= 6
                                                   else "bearish" if f["trend_template"]["passed"] <= 2 else "neutral"),
                    "Weekly reversal checks (8w)": f["weekly_reversals"]["read"],
                }
                agree = vl["obv_trend_20"] == vl["price_trend_20"]
                want["OBV vs price (20 bars)"] = (("bullish" if vl["price_trend_20"] == "rising" else "bearish")
                                                  if agree and vl["price_trend_20"] != "flat" else "neutral")
                for name, read in want.items():
                    self.assertEqual(self.read_of(f, name), read, name)
                tally = {k: sum(r["read"] == k for r in f["ledger"]["rows"]) for k in ("bullish", "bearish", "neutral")}
                self.assertEqual(f["ledger"]["tally"], tally)

    def test_score_label_levels_and_plan(self):
        for seed, regime in CASES:
            with self.subTest(regime=regime):
                daily, weekly = market(seed, regime)
                f = ta.analyse(daily, weekly, account=100_000)
                st, rows, p, a = f["stance"], f["ledger"]["rows"], f["price"], f["volatility"]["atr14"]
                w = [ta.WEIGHTS.get(r["indicator"], 0) for r in rows]
                net = sum(x * {"bullish": 1, "bearish": -1}.get(r["read"], 0) for x, r in zip(w, rows))
                self.assertEqual((st["net"], st["total"]), (net, sum(w)))
                self.close(st["score"], net / sum(w), "score")
                band = next(name for floor, name in ta.BANDS if st["score"] >= floor)
                demoted = "moves down to support" in st["plan"]["note"]
                self.assertEqual(st["label"], "ACCUMULATE" if demoted else band)
                if demoted:
                    self.assertIn(band, ("BUY", "ACCUMULATE"))
                sup = [x["price"] for x in f["levels"]["supports"]]
                res = [x["price"] for x in f["levels"]["resistances"]]
                s1 = next((x for x in sup if p - x >= 0.5 * a), sup[0] if sup else None)
                r1 = next((x for x in res if x - p >= 0.5 * a), res[0] if res else None)
                if r1 is None and f["weekly"]["range_52w"]["high"] > p:
                    r1 = f["weekly"]["range_52w"]["high"]
                if s1 is None and f["weekly"]["range_52w"]["low"] < p:
                    s1 = f["weekly"]["range_52w"]["low"]
                self.assertEqual((st["support_1"], st["resistance_1"]), (s1, r1))
                pl = st["plan"]
                moved_up = "moves up to resistance" in pl["note"]
                if st["label"] in ("BUY", "ACCUMULATE"):
                    entry = s1 if demoted else p                       # the rule: entry at the price, or at support
                    stop = (s1 - 0.5 * a) if s1 is not None else p - 2 * a
                    self.assertEqual((pl["entry"], pl["target"]), (entry, r1))
                    self.close(pl["stop"], stop, "stop half an ATR under support")
                elif st["label"] in ("REDUCE", "SELL"):
                    entry = r1 if moved_up else p
                    stop = (r1 + 0.5 * a) if r1 is not None else p + 2 * a
                    self.assertEqual((pl["entry"], pl["target"]), (entry, s1))
                    self.close(pl["stop"], stop, "stop half an ATR over resistance")
                if st["label"] in ("BUY", "ACCUMULATE"):
                    self.assertLess(pl["stop"], pl["entry"])
                    if pl["target"] is not None:
                        self.assertGreater(pl["target"], pl["entry"])
                        self.close(pl["reward_risk"], (pl["target"] - pl["entry"]) / (pl["entry"] - pl["stop"]), "R:R")
                elif st["label"] in ("REDUCE", "SELL"):
                    self.assertGreater(pl["stop"], pl["entry"])
                    if pl["target"] is not None:
                        self.assertLess(pl["target"], pl["entry"])
                        self.close(pl["reward_risk"], (pl["entry"] - pl["target"]) / (pl["stop"] - pl["entry"]), "R:R")
                else:
                    self.assertEqual((pl["entry"], pl["stop"], pl["target"]), (None, None, None))
                if pl["entry"] is not None:
                    self.close(pl["target_2r"], pl["entry"] + 2 * (pl["entry"] - pl["stop"]), "2R target")
                    ps = f["position_size"]
                    risk = 100_000 * 0.01
                    by_risk, by_cap = int(risk // abs(pl["entry"] - pl["stop"])), int(10_000 // pl["entry"])
                    self.assertEqual(ps["shares"], max(0, min(by_risk, by_cap)))
                    self.close(ps["actual_risk"], ps["shares"] * abs(pl["entry"] - pl["stop"]), "actual risk")
                else:
                    self.assertIsNone(f["position_size"])
                pulls = [r["indicator"] for r in rows if r["read"] == "bullish" and ta.WEIGHTS.get(r["indicator"])]
                drags = [r["indicator"] for r in rows if r["read"] == "bearish" and ta.WEIGHTS.get(r["indicator"])]
                self.assertEqual((st["for"], st["against"]), (pulls, drags))


if __name__ == "__main__":
    unittest.main()
