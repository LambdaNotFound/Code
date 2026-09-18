#!/usr/bin/env python3
"""Risk/reward facts for one listed stock option. Stdlib only.

Everything printed is [COMPUTED] from the inputs; the reading belongs to the
caller (see ../references/options.md). The pricing model is Black-Scholes on a
driftless lognormal, which is a [FRAME]: it prices the contract consistently,
it does not know where the stock is going.

Usage:
  opt.py --symbol AAPL --expiry 2026-10-16 --strike 340 --type call [--side buy|sell]
         --daily DAILY_FILE
         (--quote bid=5.20,ask=5.40[,iv=0.31,oi=1200,volume=300,last=5.30] | --chain CHAIN_FILE)
         [--asof YYYY-MM-DD] [--rate 0.04] [--div-yield 0] [--ta-json TA_JSON]
         [--earnings YYYY-MM-DD] [--covered] [--contracts N] [--account EQUITY] [--json]

--daily is the underlying's daily OHLCV (raw CSV, the offloaded {"result": csv}
file, or Alpha Vantage JSON). It supplies the spot, realized volatility, and ATR.
--quote is the contract quote from your broker. --chain is an options-chain file
in Alpha Vantage's column layout (contractID, expiration, strike, type, bid, ask,
volume, open_interest, implied_volatility, ...); the contract is picked from it.
--ta-json is the output of the analysis-stock script with --json; it adds the
stance and the support/resistance levels to the scenario table.
"""

import argparse
import csv
import io
import json
import math
import sys
from datetime import date, datetime

# --------------------------------------------------------------------------- #
# Underlying bars (same three shapes the analysis-stock script accepts)
# --------------------------------------------------------------------------- #


def _rows_from_text(text):
    stripped = text.lstrip()
    if stripped.startswith("{"):
        d = json.loads(stripped)
        if "error" in d:
            raise ValueError(f"Alpha Vantage error: {d['error']}")
        if isinstance(d.get("result"), str):
            inner = d["result"].lstrip()
            if inner.startswith("{"):
                d = json.loads(inner)
            else:
                return list(csv.DictReader(io.StringIO(inner.strip())))
        key = next((k for k in d if k.startswith(("Time Series", "Weekly", "Monthly"))), None)
        if key is None:
            raise ValueError("JSON has no Alpha Vantage time-series block")
        rows = []
        for dt, row in d[key].items():
            col = {k.split(". ", 1)[-1]: v for k, v in row.items()}
            col["timestamp"] = dt
            rows.append(col)
        return rows
    return list(csv.DictReader(io.StringIO(text.strip())))


def load_bars(path):
    with open(path, encoding="utf-8") as fh:
        rows = _rows_from_text(fh.read())
    bars = []
    for r in rows:
        k = {kk.strip().lower(): v for kk, v in r.items() if kk}
        dt = (k.get("timestamp") or k.get("date") or "").strip()
        if not dt:
            continue
        bars.append({"date": dt, "high": float(k["high"]), "low": float(k["low"]),
                     "close": float(k.get("adjusted close") or k["close"])})
    bars.sort(key=lambda b: b["date"])
    if not bars:
        raise ValueError("no bars in daily file")
    return bars


def realized_vol(closes, n):
    """Annualised standard deviation of daily log returns over the last n bars."""
    if len(closes) <= n:
        return None
    rets = [math.log(closes[i] / closes[i - 1]) for i in range(len(closes) - n, len(closes))]
    m = sum(rets) / n
    var = sum((r - m) ** 2 for r in rets) / (n - 1)
    return math.sqrt(var) * math.sqrt(252)


def atr(bars, n=14):
    tr = [bars[0]["high"] - bars[0]["low"]]
    for i in range(1, len(bars)):
        h, l, pc = bars[i]["high"], bars[i]["low"], bars[i - 1]["close"]
        tr.append(max(h - l, abs(h - pc), abs(l - pc)))
    if len(tr) < n:
        return None
    s = sum(tr[:n]) / n
    for v in tr[n:]:
        s += (v - s) / n
    return s


# --------------------------------------------------------------------------- #
# Black-Scholes
# --------------------------------------------------------------------------- #


def norm_cdf(x):
    return 0.5 * (1.0 + math.erf(x / math.sqrt(2.0)))


def norm_pdf(x):
    return math.exp(-0.5 * x * x) / math.sqrt(2.0 * math.pi)


def bs_price(kind, s, k, t, r, sigma, q=0.0):
    """European price; kind is 'call' or 'put'. t in years."""
    if t <= 0 or sigma <= 0:
        return max(0.0, s - k) if kind == "call" else max(0.0, k - s)
    d1 = (math.log(s / k) + (r - q + 0.5 * sigma * sigma) * t) / (sigma * math.sqrt(t))
    d2 = d1 - sigma * math.sqrt(t)
    if kind == "call":
        return s * math.exp(-q * t) * norm_cdf(d1) - k * math.exp(-r * t) * norm_cdf(d2)
    return k * math.exp(-r * t) * norm_cdf(-d2) - s * math.exp(-q * t) * norm_cdf(-d1)


def bs_greeks(kind, s, k, t, r, sigma, q=0.0):
    """delta, gamma, theta (per calendar day), vega (per 1 vol point), rho (per 1%)."""
    if t <= 0 or sigma <= 0:
        itm = (s > k) if kind == "call" else (s < k)
        return {"delta": (1.0 if kind == "call" else -1.0) if itm else 0.0, "gamma": 0.0, "theta": 0.0, "vega": 0.0, "rho": 0.0}
    sq = math.sqrt(t)
    d1 = (math.log(s / k) + (r - q + 0.5 * sigma * sigma) * t) / (sigma * sq)
    d2 = d1 - sigma * sq
    eq, er = math.exp(-q * t), math.exp(-r * t)
    gamma = eq * norm_pdf(d1) / (s * sigma * sq)
    vega = s * eq * norm_pdf(d1) * sq / 100.0
    if kind == "call":
        delta = eq * norm_cdf(d1)
        theta = (-s * eq * norm_pdf(d1) * sigma / (2 * sq) - r * k * er * norm_cdf(d2) + q * s * eq * norm_cdf(d1)) / 365.0
        rho = k * t * er * norm_cdf(d2) / 100.0
    else:
        delta = -eq * norm_cdf(-d1)
        theta = (-s * eq * norm_pdf(d1) * sigma / (2 * sq) + r * k * er * norm_cdf(-d2) - q * s * eq * norm_cdf(-d1)) / 365.0
        rho = -k * t * er * norm_cdf(-d2) / 100.0
    return {"delta": delta, "gamma": gamma, "theta": theta, "vega": vega, "rho": rho}


def implied_vol(kind, price, s, k, t, r, q=0.0):
    """Bisection on sigma. Returns None when the price sits outside no-arbitrage bounds."""
    if t <= 0 or price <= 0:
        return None
    intrinsic = max(0.0, s * math.exp(-q * t) - k * math.exp(-r * t)) if kind == "call" else max(0.0, k * math.exp(-r * t) - s * math.exp(-q * t))
    if price < intrinsic - 1e-9:
        return None
    lo, hi = 1e-4, 5.0
    if bs_price(kind, s, k, t, r, hi, q) < price:
        return None
    for _ in range(100):
        mid = 0.5 * (lo + hi)
        if bs_price(kind, s, k, t, r, mid, q) > price:
            hi = mid
        else:
            lo = mid
    return 0.5 * (lo + hi)


# --------------------------------------------------------------------------- #
# Position maths
# --------------------------------------------------------------------------- #


def payoff_at_expiry(kind, side, k, premium, s_t):
    """Per-share P&L at expiry for one contract-share."""
    intrinsic = max(0.0, s_t - k) if kind == "call" else max(0.0, k - s_t)
    return (intrinsic - premium) if side == "buy" else (premium - intrinsic)


def breakeven(kind, k, premium):
    return k + premium if kind == "call" else k - premium


def lognormal_grid(s, sigma, t, n=4001, width=6.0):
    """(prices, probabilities) for a driftless lognormal over +/- width sigma."""
    if t <= 0 or sigma <= 0:
        return [s], [1.0]
    sd = sigma * math.sqrt(t)
    xs = [(-width + 2 * width * i / (n - 1)) for i in range(n)]
    probs = [norm_pdf(x) for x in xs]
    total = sum(probs)
    probs = [p / total for p in probs]
    prices = [s * math.exp(-0.5 * sd * sd + sd * x) for x in xs]
    return prices, probs


def distribution_stats(kind, side, k, premium, s, sigma, t):
    prices, probs = lognormal_grid(s, sigma, t)
    pnl = [payoff_at_expiry(kind, side, k, premium, p) for p in prices]
    p_profit = sum(pr for pr, v in zip(probs, pnl) if v > 0)
    ev = sum(pr * v for pr, v in zip(probs, pnl))
    gains = [(pr, v) for pr, v in zip(probs, pnl) if v > 0]
    losses = [(pr, v) for pr, v in zip(probs, pnl) if v <= 0]
    avg_gain = (sum(pr * v for pr, v in gains) / sum(pr for pr, _ in gains)) if gains else 0.0
    avg_loss = (sum(pr * v for pr, v in losses) / sum(pr for pr, _ in losses)) if losses else 0.0
    return {"p_profit": p_profit, "ev": ev, "avg_gain": avg_gain, "avg_loss": avg_loss}


def scenario_table(kind, side, k, premium, s, sigma, t, levels=None):
    """P&L at expiry at +/- 1 and 2 sigma, the strike, break-even, and any TA levels."""
    sd = sigma * math.sqrt(t) if t > 0 else 0.0
    pts = [("-2 sigma", s * math.exp(-2 * sd)), ("-1 sigma", s * math.exp(-sd)), ("unchanged", s),
           ("+1 sigma", s * math.exp(sd)), ("+2 sigma", s * math.exp(2 * sd)),
           ("strike", k), ("break-even", breakeven(kind, k, premium))]
    for name, price in (levels or []):
        pts.append((name, price))
    pts.sort(key=lambda x: x[1])
    return [{"scenario": n, "price": p, "move_pct": 100.0 * (p / s - 1.0), "pnl": payoff_at_expiry(kind, side, k, premium, p)} for n, p in pts]


# --------------------------------------------------------------------------- #
# Chain / quote input
# --------------------------------------------------------------------------- #


def parse_quote(text):
    out = {}
    for part in text.split(","):
        if "=" not in part:
            continue
        key, val = part.split("=", 1)
        out[key.strip().lower()] = float(val)
    if "bid" not in out or "ask" not in out:
        raise ValueError("--quote needs at least bid= and ask=")
    return out


def load_chain(path):
    with open(path, encoding="utf-8") as fh:
        rows = _rows_from_text(fh.read())
    return [{k.strip().lower(): v for k, v in r.items() if k} for r in rows]


def pick_contract(chain, expiry, strike, kind):
    for r in chain:
        try:
            if r.get("expiration") == expiry and abs(float(r.get("strike", "nan")) - strike) < 1e-6 and r.get("type", "").lower() == kind:
                q = {"bid": float(r["bid"]), "ask": float(r["ask"])}
                for src, dst in (("implied_volatility", "iv"), ("open_interest", "oi"), ("volume", "volume"), ("last", "last"),
                                 ("delta", "delta"), ("theta", "theta"), ("gamma", "gamma"), ("vega", "vega")):
                    if r.get(src) not in (None, ""):
                        q[dst] = float(r[src])
                q["contract_id"] = r.get("contractid") or r.get("contract_id")
                return q
        except ValueError:
            continue
    raise ValueError(f"contract {expiry} {strike} {kind} not in chain")


def chain_context(chain, expiry, kind, strike):
    """Same-expiry neighbours: IV by strike, to show skew around the contract."""
    rows = []
    for r in chain:
        try:
            if r.get("expiration") == expiry and r.get("type", "").lower() == kind and r.get("implied_volatility") not in (None, ""):
                rows.append((float(r["strike"]), float(r["implied_volatility"]), float(r.get("open_interest") or 0)))
        except ValueError:
            continue
    rows.sort()
    idx = next((i for i, (k, _, _) in enumerate(rows) if abs(k - strike) < 1e-6), None)
    if idx is None:
        return []
    return [{"strike": k, "iv": iv, "oi": oi} for k, iv, oi in rows[max(0, idx - 3): idx + 4]]


# --------------------------------------------------------------------------- #
# Analysis
# --------------------------------------------------------------------------- #


def analyse(symbol, expiry, strike, kind, side, bars, quote, asof=None, rate=0.04, div_yield=0.0,
            ta=None, earnings=None, covered=False, contracts=1, account=None):
    closes = [b["close"] for b in bars]
    s = closes[-1]
    asof_d = datetime.strptime(asof or bars[-1]["date"], "%Y-%m-%d").date()
    exp_d = datetime.strptime(expiry, "%Y-%m-%d").date()
    dte = (exp_d - asof_d).days
    if dte <= 0:
        raise ValueError("expiry is not after the as-of date")
    t = dte / 365.0

    bid, ask = quote["bid"], quote["ask"]
    mid = 0.5 * (bid + ask)
    premium = ask if side == "buy" else bid  # what you actually pay or receive
    spread = ask - bid
    spread_pct = 100.0 * spread / mid if mid else None

    iv_quoted = quote.get("iv")
    iv_solved = implied_vol(kind, mid, s, strike, t, rate, div_yield)
    iv = iv_quoted if iv_quoted else iv_solved
    hv20, hv60 = realized_vol(closes, 20), realized_vol(closes, 60)
    hv = hv60 or hv20
    greeks = bs_greeks(kind, s, strike, t, rate, iv, div_yield) if iv else None
    fair_at_hv = bs_price(kind, s, strike, t, rate, hv, div_yield) if hv else None
    intrinsic = max(0.0, s - strike) if kind == "call" else max(0.0, strike - s)
    extrinsic = mid - intrinsic
    moneyness = 100.0 * (s / strike - 1.0) if kind == "call" else 100.0 * (strike / s - 1.0)

    be = breakeven(kind, strike, premium)
    move_to_be_pct = 100.0 * (be / s - 1.0)
    exp_move_iv = s * iv * math.sqrt(t) if iv else None
    exp_move_hv = s * hv * math.sqrt(t) if hv else None
    be_in_sigmas = abs(be - s) / exp_move_hv if exp_move_hv else None

    stats_hv = distribution_stats(kind, side, strike, premium, s, hv, t) if hv else None
    stats_iv = distribution_stats(kind, side, strike, premium, s, iv, t) if iv else None

    per_contract = 100 * contracts
    if side == "buy":
        max_loss = premium * per_contract
        max_gain = None if kind == "call" else (strike - premium) * per_contract
        capital = max_loss
    else:
        max_gain = premium * per_contract
        if kind == "put":
            max_loss = (strike - premium) * per_contract
            capital = strike * per_contract  # cash-secured
        else:
            max_loss = None if not covered else None  # unbounded naked; covered call loss lives in the stock
            capital = None
    sd_hv = hv * math.sqrt(t) if hv else None
    adverse_2s = s * math.exp(2 * sd_hv) if (kind == "call") else s * math.exp(-2 * sd_hv) if sd_hv else None
    if side == "buy":
        adverse_2s = s * math.exp(-2 * sd_hv) if kind == "call" else s * math.exp(2 * sd_hv)
    favorable_1s = s * math.exp(sd_hv) if kind == "call" else s * math.exp(-sd_hv)
    if side == "sell":
        favorable_1s = s  # a short wants nothing to happen
    pnl_fav_1s = payoff_at_expiry(kind, side, strike, premium, favorable_1s) * per_contract
    pnl_adv_2s = payoff_at_expiry(kind, side, strike, premium, adverse_2s) * per_contract
    if side == "buy":
        reward_risk = (pnl_fav_1s / max_loss) if max_loss else None
    else:
        reward_risk = (max_gain / -pnl_adv_2s) if pnl_adv_2s < 0 else None

    theta_day = greeks["theta"] if greeks else None
    theta_pct = (100.0 * abs(theta_day) / mid) if (theta_day is not None and mid) else None

    levels = []
    stance = None
    if ta:
        stance = (ta.get("stance") or {}).get("label")
        for c in (ta.get("levels") or {}).get("supports", [])[:2]:
            levels.append((f"support {c['price']:.2f}", c["price"]))
        for c in (ta.get("levels") or {}).get("resistances", [])[:2]:
            levels.append((f"resistance {c['price']:.2f}", c["price"]))
    earnings_inside = None
    if earnings:
        e = datetime.strptime(earnings, "%Y-%m-%d").date()
        earnings_inside = asof_d < e <= exp_d

    out = {
        "symbol": symbol, "contract": f"{symbol} {expiry} {strike:g} {kind}", "side": side, "covered": covered,
        "contracts": contracts, "as_of": asof_d.isoformat(), "dte": dte,
        "underlying": {"spot": s, "hv20": hv20, "hv60": hv60, "atr14": atr(bars), "atr14_pct": (100.0 * atr(bars) / s) if atr(bars) else None},
        "quote": {"bid": bid, "ask": ask, "mid": mid, "last": quote.get("last"), "premium_used": premium,
                  "spread": spread, "spread_pct_of_mid": spread_pct, "oi": quote.get("oi"), "volume": quote.get("volume"),
                  "contract_id": quote.get("contract_id")},
        "pricing": {"iv": iv, "iv_source": "quoted" if iv_quoted else ("solved from mid" if iv_solved else "unavailable"),
                    "iv_solved_from_mid": iv_solved, "hv_used": hv, "iv_over_hv": (iv / hv) if (iv and hv) else None,
                    "fair_value_at_hv": fair_at_hv, "mid_vs_fair_pct": (100.0 * (mid / fair_at_hv - 1.0)) if fair_at_hv else None,
                    "intrinsic": intrinsic, "extrinsic": extrinsic, "moneyness_pct": moneyness,
                    "rate": rate, "div_yield": div_yield, "greeks": greeks,
                    "theta_per_day": theta_day, "theta_pct_of_mid_per_day": theta_pct},
        "move": {"break_even": be, "move_to_break_even_pct": move_to_be_pct,
                 "expected_move_iv": exp_move_iv, "expected_move_iv_pct": (100.0 * exp_move_iv / s) if exp_move_iv else None,
                 "expected_move_hv": exp_move_hv, "expected_move_hv_pct": (100.0 * exp_move_hv / s) if exp_move_hv else None,
                 "break_even_in_hv_sigmas": be_in_sigmas},
        "risk": {"max_loss": max_loss, "max_gain": max_gain, "capital_at_risk": capital,
                 "pnl_favorable_1s": pnl_fav_1s, "favorable_1s_price": favorable_1s,
                 "pnl_adverse_2s": pnl_adv_2s, "adverse_2s_price": adverse_2s, "reward_to_risk": reward_risk},
        "odds_hv": stats_hv, "odds_iv": stats_iv,
        "scenarios": scenario_table(kind, side, strike, premium, s, hv or iv or 0.0, t, levels),
        "ta_stance": stance, "earnings": earnings, "earnings_inside": earnings_inside,
        "position_size": None,
    }
    for row in out["scenarios"]:
        row["pnl"] *= per_contract
    if account and (capital or max_loss):
        risk = max_loss if max_loss is not None else None
        out["position_size"] = {"account": account, "risk_pct_of_account": (100.0 * risk / account) if risk else None,
                                "capital_pct_of_account": (100.0 * capital / account) if capital else None}
    out["ledger"] = ledger(out)
    out["verdict"] = verdict(out)
    return out


# --------------------------------------------------------------------------- #
# Ledger and verdict
# --------------------------------------------------------------------------- #


def ledger(f):
    rows = []
    side, q, p, m, r, o = f["side"], f["quote"], f["pricing"], f["move"], f["risk"], f["odds_hv"]

    def add(name, value, read, rule):
        rows.append({"check": name, "value": value, "read": read, "rule": rule})

    if q["spread_pct_of_mid"] is not None:
        sp = q["spread_pct_of_mid"]
        add("Bid-ask spread", f"{sp:.1f}% of mid", "good" if sp <= 5 else "neutral" if sp <= 10 else "bad", "<= 5% good, > 10% bad")
    if q["oi"] is not None:
        add("Open interest", f"{q['oi']:,.0f}", "good" if q["oi"] >= 500 else "neutral" if q["oi"] >= 100 else "bad", ">= 500 good, < 100 bad")
    if p["iv_over_hv"] is not None:
        ratio = p["iv_over_hv"]
        if side == "buy":
            read = "good" if ratio <= 0.9 else "bad" if ratio >= 1.2 else "neutral"
        else:
            read = "good" if ratio >= 1.2 else "bad" if ratio <= 0.9 else "neutral"
        add("IV vs realized vol", f"IV {100*p['iv']:.1f}% / HV {100*p['hv_used']:.1f}% = {ratio:.2f}", read,
            "buyer wants cheap vol (<= 0.9), seller wants rich vol (>= 1.2)")
    dte = f["dte"]
    if side == "buy":
        read = "bad" if dte < 14 else "good" if 30 <= dte <= 120 else "neutral"
        rule = "buyer: < 14 days is theta's territory; 30-120 good"
    else:
        read = "bad" if (dte < 7 or dte > 90) else "good" if 20 <= dte <= 45 else "neutral"
        rule = "seller: 20-45 days good; < 7 gamma risk, > 90 slow decay"
    add("Days to expiry", str(dte), read, rule)
    if p["theta_pct_of_mid_per_day"] is not None:
        tp = p["theta_pct_of_mid_per_day"]
        read = ("bad" if tp > 2 else "good" if tp < 1 else "neutral") if side == "buy" else ("good" if tp > 1.5 else "bad" if tp < 0.5 else "neutral")
        add("Theta burn", f"{tp:.2f}% of mid per day", read, "buyer: > 2%/day bad; seller: > 1.5%/day good")
    if m["break_even_in_hv_sigmas"] is not None:
        bz = m["break_even_in_hv_sigmas"]
        if side == "buy":
            read = "good" if bz <= 0.5 else "bad" if bz >= 1.0 else "neutral"
            rule = "buyer: break-even inside 0.5 sigma good, beyond 1 sigma bad"
        else:
            read = "good" if bz >= 1.0 else "bad" if bz <= 0.5 else "neutral"
            rule = "seller: break-even beyond 1 sigma good, inside 0.5 sigma bad"
        add("Break-even distance", f"{m['move_to_break_even_pct']:+.1f}% = {bz:.2f} HV sigmas over {dte}d", read, rule)
    if o:
        pp = 100.0 * o["p_profit"]
        if side == "buy":
            read = "good" if pp >= 40 else "bad" if pp < 25 else "neutral"
            rule = "buyer: >= 40% good, < 25% bad"
        else:
            read = "good" if pp >= 70 else "bad" if pp < 55 else "neutral"
            rule = "seller: >= 70% good, < 55% bad"
        add("P(profit) at HV", f"{pp:.0f}%", read, rule)
        base = q["premium_used"] if q["premium_used"] else None
        if base:
            evp = 100.0 * o["ev"] / base
            add("Expected value at HV", f"{evp:+.0f}% of premium", "good" if evp > 5 else "bad" if evp < -10 else "neutral",
                "> +5% of premium good, < -10% bad; a driftless model, so near zero is normal")
    if r["reward_to_risk"] is not None:
        rr = r["reward_to_risk"]
        if side == "buy":
            read = "good" if rr >= 1.0 else "bad" if rr < 0.5 else "neutral"
            rule = "buyer: gain at +1 sigma / max loss; >= 1 good, < 0.5 bad"
        else:
            read = "good" if rr >= 0.5 else "bad" if rr < 0.25 else "neutral"
            rule = "seller: premium / loss at 2 sigma against; >= 0.5 good, < 0.25 bad"
        add("Reward to risk", f"{rr:.2f}", read, rule)
    if f["ta_stance"]:
        bullish_pos = (f["side"] == "buy") == (f["contract"].endswith("call"))
        st = f["ta_stance"]
        bull = st in ("BUY", "ACCUMULATE")
        bear = st in ("SELL", "REDUCE")
        read = "good" if (bullish_pos and bull) or (not bullish_pos and bear) else "bad" if (bullish_pos and bear) or (not bullish_pos and bull) else "neutral"
        add("Alignment with analysis-stock stance", st, read, "position direction vs the chart stance")
    if f["earnings_inside"] is not None:
        if f["earnings_inside"]:
            add("Earnings inside the window", f["earnings"], "bad" if side == "sell" else "neutral",
                "seller: gap risk; buyer: IV crush after the print, so the move must beat it")
        else:
            add("Earnings inside the window", "none", "neutral", "no scheduled print before expiry")
    tally = {k: sum(1 for x in rows if x["read"] == k) for k in ("good", "bad", "neutral")}
    return {"rows": rows, "tally": tally}


def verdict(f):
    rows = f["ledger"]["rows"]
    t = f["ledger"]["tally"]
    scored = t["good"] + t["bad"] + t["neutral"]
    score = (t["good"] - t["bad"]) / scored if scored else 0.0
    label = "GOOD" if score >= 0.4 else "POOR" if score <= -0.2 else "FAIR"
    hard = []
    sp = f["quote"]["spread_pct_of_mid"]
    if sp is not None and sp > 25:
        hard.append(f"spread is {sp:.0f}% of mid")
    if f["odds_hv"] and f["odds_hv"]["p_profit"] < 0.10:
        hard.append(f"P(profit) {100*f['odds_hv']['p_profit']:.0f}%")
    if f["pricing"]["iv"] is None:
        hard.append("no implied vol: the mid is outside model bounds or the quote is stale")
    if f["side"] == "sell" and not f["contract"].endswith("put") and not f["covered"]:
        hard.append("naked short call: loss is unbounded")
    if hard:
        label = "POOR"
    return {"label": label, "score": score, "hard_fails": hard,
            "for": [x["check"] for x in rows if x["read"] == "good"], "against": [x["check"] for x in rows if x["read"] == "bad"]}


# --------------------------------------------------------------------------- #
# Output
# --------------------------------------------------------------------------- #


def _f(v, nd=2, suffix=""):
    return "n/a" if v is None else f"{v:,.{nd}f}{suffix}"


def _pct(v, nd=1):
    return "n/a" if v is None else f"{v:+.{nd}f}%"


def render(f):
    q, p, m, r, u = f["quote"], f["pricing"], f["move"], f["risk"], f["underlying"]
    L = [f"# Option facts: {f['side'].upper()} {f['contracts']}x {f['contract']}" + (" (covered)" if f["covered"] else ""),
         f"As of {f['as_of']}; {f['dte']} days to expiry; spot {_f(u['spot'])}; rate {100*p['rate']:.2f}% and dividend yield {100*p['div_yield']:.2f}% assumed.", ""]
    L += ["## Quote and liquidity", "| | |", "|---|---|",
          f"| Bid / ask / mid | {_f(q['bid'])} / {_f(q['ask'])} / {_f(q['mid'])} |",
          f"| Premium used ({'pay the ask' if f['side'] == 'buy' else 'receive the bid'}) | {_f(q['premium_used'])} per share = {_f(q['premium_used'] * 100 * f['contracts'], 0)} total |",
          f"| Spread | {_f(q['spread'])} = {_f(q['spread_pct_of_mid'], 1, '%')} of mid |",
          f"| Open interest / volume | {_f(q['oi'], 0)} / {_f(q['volume'], 0)} |",
          f"| Intrinsic / extrinsic | {_f(p['intrinsic'])} / {_f(p['extrinsic'])} |",
          f"| Moneyness | {_pct(p['moneyness_pct'])} ({'ITM' if p['moneyness_pct'] > 0 else 'OTM' if p['moneyness_pct'] < 0 else 'ATM'}) |", ""]
    g = p["greeks"] or {}
    L += ["## Volatility and pricing", "| | |", "|---|---|",
          f"| Implied vol ({p['iv_source']}) | {_f(100*p['iv'] if p['iv'] else None, 1, '%')} |",
          f"| Realized vol 20d / 60d | {_f(100*u['hv20'] if u['hv20'] else None, 1, '%')} / {_f(100*u['hv60'] if u['hv60'] else None, 1, '%')} |",
          f"| IV / HV | {_f(p['iv_over_hv'])} |",
          f"| Model value at HV | {_f(p['fair_value_at_hv'])} (mid is {_pct(p['mid_vs_fair_pct'])} vs it) |",
          f"| Delta / gamma / theta per day / vega per vol point | {_f(g.get('delta'), 3)} / {_f(g.get('gamma'), 4)} / {_f(g.get('theta'), 3)} / {_f(g.get('vega'), 3)} |",
          f"| Theta as % of mid per day | {_f(p['theta_pct_of_mid_per_day'], 2, '%')} |",
          f"| ATR14 | {_f(u['atr14'])} ({_f(u['atr14_pct'], 2, '%')} of spot) |", ""]
    L += ["## Move required", "| | |", "|---|---|",
          f"| Break-even at expiry | {_f(m['break_even'])} ({_pct(m['move_to_break_even_pct'])} from spot) |",
          f"| Expected move to expiry, IV | {_f(m['expected_move_iv'])} ({_pct(m['expected_move_iv_pct'])}) |",
          f"| Expected move to expiry, HV | {_f(m['expected_move_hv'])} ({_pct(m['expected_move_hv_pct'])}) |",
          f"| Break-even in HV sigmas | {_f(m['break_even_in_hv_sigmas'])} |", ""]
    L += ["## Risk and odds (per position, at expiry)", "| | |", "|---|---|",
          f"| Max loss | {_f(r['max_loss'], 0) if r['max_loss'] is not None else 'unbounded'} |",
          f"| Max gain | {_f(r['max_gain'], 0) if r['max_gain'] is not None else 'unbounded'} |",
          f"| Capital tied up | {_f(r['capital_at_risk'], 0) if r['capital_at_risk'] is not None else 'margin-dependent'} |",
          f"| P&L at +1 HV sigma in your favour ({_f(r['favorable_1s_price'])}) | {_f(r['pnl_favorable_1s'], 0)} |",
          f"| P&L at 2 HV sigma against ({_f(r['adverse_2s_price'])}) | {_f(r['pnl_adverse_2s'], 0)} |",
          f"| Reward to risk | {_f(r['reward_to_risk'])} |"]
    for label, o in (("HV", f["odds_hv"]), ("IV", f["odds_iv"])):
        if o:
            L.append(f"| P(profit) / EV / avg gain / avg loss at {label} | {100*o['p_profit']:.0f}% / {_f(o['ev']*100*f['contracts'], 0)} / {_f(o['avg_gain']*100*f['contracts'], 0)} / {_f(o['avg_loss']*100*f['contracts'], 0)} |")
    if f["position_size"]:
        ps = f["position_size"]
        L.append(f"| Of a {_f(ps['account'], 0)} account | risk {_f(ps['risk_pct_of_account'], 2, '%')}, capital {_f(ps['capital_pct_of_account'], 2, '%')} |")
    L.append("")
    L += ["## Scenarios at expiry", "| Scenario | Price | Move | P&L |", "|---|---|---|---|"]
    for row in f["scenarios"]:
        L.append(f"| {row['scenario']} | {_f(row['price'])} | {_pct(row['move_pct'])} | {_f(row['pnl'], 0)} |")
    L.append("")
    L += ["## Ledger", "| Check | Value | Read | Rule |", "|---|---|---|---|"]
    for x in f["ledger"]["rows"]:
        L.append(f"| {x['check']} | {x['value']} | {x['read']} | {x['rule']} |")
    t = f["ledger"]["tally"]
    L.append(f"- Tally: {t['good']} good / {t['bad']} bad / {t['neutral']} neutral.")
    L.append("")
    v = f["verdict"]
    L += ["## Verdict (rule-based)", "| | |", "|---|---|",
          f"| **Deal** | **{v['label']}** |",
          f"| Score | {v['score']:+.2f} (good minus bad over all rows; GOOD >= +0.40, POOR <= -0.20) |",
          f"| For | {', '.join(v['for']) or 'none'} |",
          f"| Against | {', '.join(v['against']) or 'none'} |",
          f"| Hard fails | {'; '.join(v['hard_fails']) or 'none'} |"]
    return "\n".join(L)


def main(argv=None):
    ap = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    ap.add_argument("--symbol", required=True)
    ap.add_argument("--expiry", required=True, help="YYYY-MM-DD")
    ap.add_argument("--strike", type=float, required=True)
    ap.add_argument("--type", choices=["call", "put"], required=True)
    ap.add_argument("--side", choices=["buy", "sell"], default="buy")
    ap.add_argument("--daily", required=True, help="underlying daily OHLCV file")
    ap.add_argument("--quote", help="bid=,ask=[,iv=,oi=,volume=,last=]  (iv as a decimal, 0.32)")
    ap.add_argument("--chain", help="options chain file in Alpha Vantage layout")
    ap.add_argument("--asof", help="valuation date YYYY-MM-DD (default: last daily bar)")
    ap.add_argument("--rate", type=float, default=0.04, help="risk-free rate, decimal (default 0.04)")
    ap.add_argument("--div-yield", type=float, default=0.0, help="continuous dividend yield, decimal")
    ap.add_argument("--ta-json", help="analysis-stock --json output for stance and levels")
    ap.add_argument("--earnings", help="next earnings date YYYY-MM-DD, if known")
    ap.add_argument("--covered", action="store_true", help="short call is covered by stock")
    ap.add_argument("--contracts", type=int, default=1)
    ap.add_argument("--account", type=float, help="account equity, to express risk as a share of it")
    ap.add_argument("--json", action="store_true")
    args = ap.parse_args(argv)

    if not args.quote and not args.chain:
        ap.error("one of --quote or --chain is required")
    bars = load_bars(args.daily)
    if args.chain:
        chain = load_chain(args.chain)
        quote = pick_contract(chain, args.expiry, args.strike, args.type)
        context = chain_context(chain, args.expiry, args.type, args.strike)
    else:
        quote, context = parse_quote(args.quote), []
    ta = None
    if args.ta_json:
        with open(args.ta_json, encoding="utf-8") as fh:
            ta = json.load(fh)
    facts = analyse(args.symbol.upper(), args.expiry, args.strike, args.type, args.side, bars, quote, args.asof,
                    args.rate, args.div_yield, ta, args.earnings, args.covered, args.contracts, args.account)
    facts["chain_context"] = context
    if args.json:
        print(json.dumps(facts, indent=2, default=str))
    else:
        print(render(facts))
        if context:
            print("\n## Same-expiry neighbours (IV by strike)")
            print("| Strike | IV | OI |")
            print("|---|---|---|")
            for c in context:
                print(f"| {c['strike']:g} | {100*c['iv']:.1f}% | {c['oi']:,.0f} |")


if __name__ == "__main__":
    main()
