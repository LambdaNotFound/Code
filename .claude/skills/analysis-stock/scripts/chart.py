#!/usr/bin/env python3
"""Render a multi-timeframe price chart for one ticker as a standalone SVG.

Stdlib only, no plotting library: candlesticks are rectangles and lines, moving
averages are polylines, levels are dashed rules. Everything drawn comes from the
same bars ta.py analyses, so the picture and the numbers cannot disagree.

Three stacked panels, newest bar on the right:
  * daily   - last N bars with SMA20/50/200, a volume strip, and gap markers
  * weekly  - last 104 bars with the 20/50-week averages
  * monthly - resampled off the weekly series, with the 12/24-month averages

Support and resistance clusters from ta.py are drawn on the daily panel and
labelled at the right edge.

Usage:
  chart.py --daily DAILY [--weekly WEEKLY] [--sma200 PRICE | --sma200-file F]
           --symbol SYM --out chart.svg [--daily-bars 100] [--theme light|dark]
"""

import argparse
import math
import os
import sys

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
import ta  # noqa: E402

THEMES = {
    "light": {"bg": "#ffffff", "panel": "#fbfbfc", "grid": "#e6e8eb", "axis": "#9aa2ad",
              "text": "#333a44", "muted": "#6b7480", "up": "#1a7f5a", "down": "#c0392b",
              "ma1": "#2563eb", "ma2": "#d97706", "ma3": "#7c3aed",
              "sup": "#1a7f5a", "res": "#c0392b", "vol": "#aab2bd"},
    "dark": {"bg": "#14171c", "panel": "#1a1e25", "grid": "#2a2f38", "axis": "#5a6472",
             "text": "#e3e7ec", "muted": "#98a2b0", "up": "#35c48c", "down": "#ef6d63",
             "ma1": "#6aa3ff", "ma2": "#f0b45a", "ma3": "#b490ff",
             "sup": "#35c48c", "res": "#ef6d63", "vol": "#39414d"},
}
W = 1160
PAD_L, PAD_R = 58, 86


def esc(s):
    return str(s).replace("&", "&amp;").replace("<", "&lt;").replace(">", "&gt;")


def nice_ticks(lo, hi, count=5):
    """Round price gridlines covering [lo, hi]."""
    if hi <= lo:
        return [lo]
    raw = (hi - lo) / count
    mag = 10 ** int(math.floor(math.log10(raw)))
    step = next(m * mag for m in (1, 2, 2.5, 5, 10) if m * mag >= raw)
    start = step * int(lo / step)
    out, v = [], start
    while v <= hi + step * 0.5:
        if v >= lo - step * 0.5:
            out.append(v)
        v += step
    return out


class Panel:
    """One price pane: maps bar index to x and price to y, and emits SVG."""

    def __init__(self, bars, top, height, theme, title, x0=PAD_L, x1=W - PAD_R):
        self.bars, self.top, self.h, self.c, self.title = bars, top, height, theme, title
        self.x0, self.x1 = x0, x1
        lows = [b["low"] for b in bars]
        highs = [b["high"] for b in bars]
        self.lo, self.hi = min(lows), max(highs)
        pad = (self.hi - self.lo) * 0.06 or 1.0
        self.lo -= pad
        self.hi += pad
        self.n = len(bars)
        self.bw = (self.x1 - self.x0) / max(self.n, 1)

    def x(self, i):
        return self.x0 + (i + 0.5) * self.bw

    def y(self, p):
        return self.top + (self.hi - p) / (self.hi - self.lo) * self.h

    def frame(self, out):
        c = self.c
        out.append(f'<rect x="{self.x0}" y="{self.top}" width="{self.x1-self.x0:.1f}" height="{self.h}" '
                   f'fill="{c["panel"]}" stroke="{c["grid"]}" stroke-width="1"/>')
        out.append(f'<text x="{self.x0}" y="{self.top-8}" font-size="13" font-weight="600" '
                   f'fill="{c["text"]}" font-family="system-ui,-apple-system,Segoe UI,sans-serif">{esc(self.title)}</text>')
        for t in nice_ticks(self.lo, self.hi):
            yy = self.y(t)
            if not (self.top <= yy <= self.top + self.h):
                continue
            out.append(f'<line x1="{self.x0}" y1="{yy:.1f}" x2="{self.x1}" y2="{yy:.1f}" '
                       f'stroke="{c["grid"]}" stroke-width="1"/>')
            out.append(f'<text x="{self.x0-6}" y="{yy+4:.1f}" font-size="11" text-anchor="end" '
                       f'fill="{c["axis"]}" font-family="system-ui,sans-serif">{t:,.0f}</text>')

    def candles(self, out):
        c = self.c
        bw = max(self.bw * 0.62, 1.0)
        for i, b in enumerate(self.bars):
            up = b["close"] >= b["open"]
            col = c["up"] if up else c["down"]
            cx = self.x(i)
            out.append(f'<line x1="{cx:.1f}" y1="{self.y(b["high"]):.1f}" x2="{cx:.1f}" '
                       f'y2="{self.y(b["low"]):.1f}" stroke="{col}" stroke-width="1"/>')
            top = self.y(max(b["open"], b["close"]))
            bot = self.y(min(b["open"], b["close"]))
            out.append(f'<rect x="{cx-bw/2:.1f}" y="{top:.1f}" width="{bw:.1f}" '
                       f'height="{max(bot-top,1):.1f}" fill="{col}"/>')

    def line(self, out, values, colour, width=1.6, dash=None):
        pts = [f"{self.x(i):.1f},{self.y(v):.1f}" for i, v in enumerate(values) if v is not None]
        if len(pts) < 2:
            return
        d = f' stroke-dasharray="{dash}"' if dash else ""
        out.append(f'<polyline points="{" ".join(pts)}" fill="none" stroke="{colour}" '
                   f'stroke-width="{width}"{d} stroke-linejoin="round"/>')

    def level(self, out, price, colour, label):
        if not (self.lo < price < self.hi):
            return
        yy = self.y(price)
        out.append(f'<line x1="{self.x0}" y1="{yy:.1f}" x2="{self.x1}" y2="{yy:.1f}" stroke="{colour}" '
                   f'stroke-width="1" stroke-dasharray="5 4" opacity="0.85"/>')
        out.append(f'<text x="{self.x1+5}" y="{yy+4:.1f}" font-size="10.5" fill="{colour}" '
                   f'font-family="system-ui,sans-serif">{esc(label)}</text>')

    def dates(self, out, every=None):
        c = self.c
        every = every or max(1, self.n // 9)
        for i in range(0, self.n, every):
            out.append(f'<text x="{self.x(i):.1f}" y="{self.top+self.h+14:.1f}" font-size="10" '
                       f'text-anchor="middle" fill="{c["axis"]}" font-family="system-ui,sans-serif">'
                       f'{esc(self.bars[i]["date"])}</text>')

    def last_price_tag(self, out):
        c = self.c
        p = self.bars[-1]["close"]
        yy = self.y(p)
        out.append(f'<rect x="{self.x1+1}" y="{yy-9:.1f}" width="56" height="18" rx="3" fill="{c["text"]}"/>')
        out.append(f'<text x="{self.x1+29}" y="{yy+4:.1f}" font-size="11" font-weight="600" text-anchor="middle" '
                   f'fill="{c["bg"]}" font-family="system-ui,sans-serif">{p:,.2f}</text>')


def legend(out, x, y, theme, items):
    for i, (label, colour) in enumerate(items):
        xx = x + i * 118
        out.append(f'<line x1="{xx}" y1="{y}" x2="{xx+20}" y2="{y}" stroke="{colour}" stroke-width="2.4"/>')
        out.append(f'<text x="{xx+25}" y="{y+4}" font-size="11" fill="{theme["muted"]}" '
                   f'font-family="system-ui,sans-serif">{esc(label)}</text>')


def build(symbol, daily, weekly, facts, theme_name="light", daily_bars=100):
    c = THEMES[theme_name]
    d = daily[-daily_bars:]
    out = []
    wk = weekly[-104:] if weekly else None
    mo = ta.to_monthly(weekly)[-72:] if weekly else None
    vol_h = 70
    h_daily, h_weekly, h_monthly = 300, 210, 190
    y = 72
    # each block is panel + 14 (date row) + 22 (legend row) + 38 (gap to the next title)
    BLOCK = 74
    total = (y + h_daily + BLOCK + vol_h + 44
             + (h_weekly + BLOCK if wk else 0)
             + (h_monthly + BLOCK if mo else 0) + 24)

    out.append(f'<svg xmlns="http://www.w3.org/2000/svg" width="{W}" height="{total}" '
               f'viewBox="0 0 {W} {total}" font-family="system-ui,-apple-system,Segoe UI,sans-serif">')
    out.append(f'<rect width="{W}" height="{total}" fill="{c["bg"]}"/>')
    last = daily[-1]
    chg = 100.0 * (last["close"] - daily[-2]["close"]) / daily[-2]["close"] if len(daily) > 1 else 0.0
    out.append(f'<text x="{PAD_L}" y="32" font-size="19" font-weight="700" fill="{c["text"]}">{esc(symbol)}'
               f'<tspan dx="14" font-size="15" font-weight="400" fill="{c["up"] if chg>=0 else c["down"]}">'
               f'{last["close"]:,.2f}</tspan>'
               f'<tspan dx="8" font-size="15" font-weight="400" fill="{c["up"] if chg>=0 else c["down"]}">'
               f'{chg:+.2f}%</tspan></text>')
    out.append(f'<text x="{W-PAD_R}" y="32" font-size="11.5" text-anchor="end" fill="{c["muted"]}">'
               f'as of {esc(last["date"])} \u00b7 close-only, no intraday</text>')

    # ---- daily price ----
    closes = [b["close"] for b in daily]
    s20, s50, s200 = ta.sma(closes, 20), ta.sma(closes, 50), ta.sma(closes, 200)
    off = len(daily) - len(d)
    p = Panel(d, y, h_daily, c, f"Daily \u00b7 last {len(d)} bars")
    p.frame(out)
    p.line(out, s20[off:], c["ma1"])
    p.line(out, s50[off:], c["ma2"])
    ma200 = (facts or {}).get("ma", {}).get("sma200")
    src = (facts or {}).get("ma", {}).get("sma200_source", "")
    if any(v is not None for v in s200[off:]):
        p.line(out, s200[off:], c["ma3"])
    elif ma200:
        p.level(out, ma200, c["ma3"], f"200d {ma200:,.0f}")
    p.candles(out)
    lv = (facts or {}).get("levels", {})
    for cl in lv.get("supports", [])[:3]:
        p.level(out, cl["price"], c["sup"], f'S {cl["price"]:,.2f}')
    for cl in lv.get("resistances", [])[:3]:
        p.level(out, cl["price"], c["res"], f'R {cl["price"]:,.2f}')
    for g in (facts or {}).get("gaps_20d", []):
        idx = next((i for i, b in enumerate(d) if b["date"] == g["date"]), None)
        if idx is None:
            continue
        col = c["up"] if g["pct"] > 0 else c["down"]
        out.append(f'<circle cx="{p.x(idx):.1f}" cy="{p.top+11:.1f}" r="3.6" fill="none" stroke="{col}" stroke-width="1.6"/>')
        if not g["filled"]:
            out.append(f'<text x="{p.x(idx):.1f}" y="{p.top+28:.1f}" font-size="9" text-anchor="middle" '
                       f'fill="{col}">gap</text>')
    p.last_price_tag(out)
    p.dates(out)
    legend(out, PAD_L, y + h_daily + 36, c,
           [("SMA20", c["ma1"]), ("SMA50", c["ma2"]), (f"SMA200 ({src or 'n/a'})", c["ma3"]),
            ("support", c["sup"]), ("resistance", c["res"])])

    # ---- volume strip ----
    vy = y + h_daily + BLOCK
    vmax = max(b["volume"] for b in d) or 1
    out.append(f'<rect x="{PAD_L}" y="{vy}" width="{W-PAD_R-PAD_L}" height="{vol_h}" fill="{c["panel"]}" '
               f'stroke="{c["grid"]}" stroke-width="1"/>')
    out.append(f'<text x="{PAD_L+6}" y="{vy+14}" font-size="11" fill="{c["muted"]}">Volume (20-bar average line)</text>')
    avg = ta.sma([b["volume"] for b in d], 20)
    bw = max(p.bw * 0.62, 1.0)
    for i, b in enumerate(d):
        hh = b["volume"] / vmax * (vol_h - 6)
        up = b["close"] >= b["open"]
        out.append(f'<rect x="{p.x(i)-bw/2:.1f}" y="{vy+vol_h-hh:.1f}" width="{bw:.1f}" height="{hh:.1f}" '
                   f'fill="{c["up"] if up else c["down"]}" opacity="0.55"/>')
    pts = [f"{p.x(i):.1f},{vy+vol_h-(v/vmax*(vol_h-6)):.1f}" for i, v in enumerate(avg) if v is not None]
    if len(pts) > 1:
        out.append(f'<polyline points="{" ".join(pts)}" fill="none" stroke="{c["muted"]}" stroke-width="1.3"/>')

    cursor = vy + vol_h + 44

    # ---- weekly ----
    if wk:
        wc = [b["close"] for b in weekly]
        w20, w50 = ta.sma(wc, 20), ta.sma(wc, 50)
        woff = len(weekly) - len(wk)
        pw = Panel(wk, cursor, h_weekly, c, f"Weekly \u00b7 last {len(wk)} bars")
        pw.frame(out)
        pw.line(out, w20[woff:], c["ma1"])
        pw.line(out, w50[woff:], c["ma2"])
        pw.candles(out)
        pw.last_price_tag(out)
        pw.dates(out)
        legend(out, PAD_L, cursor + h_weekly + 36, c, [("20-week", c["ma1"]), ("50-week", c["ma2"])])
        cursor += h_weekly + BLOCK

    # ---- monthly ----
    if mo:
        mc = [b["close"] for b in ta.to_monthly(weekly)]
        m12, m24 = ta.sma(mc, 12), ta.sma(mc, 24)
        moff = len(mc) - len(mo)
        pm = Panel(mo, cursor, h_monthly, c, f"Monthly \u00b7 last {len(mo)} bars, resampled from weekly")
        pm.frame(out)
        pm.line(out, m12[moff:], c["ma1"])
        pm.line(out, m24[moff:], c["ma2"])
        pm.candles(out)
        pm.last_price_tag(out)
        pm.dates(out, every=max(1, len(mo) // 8))
        legend(out, PAD_L, cursor + h_monthly + 36, c, [("12-month", c["ma1"]), ("24-month", c["ma2"])])

    out.append('</svg>')
    return "\n".join(out)


def main(argv=None):
    ap = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    ap.add_argument("--daily", required=True)
    ap.add_argument("--weekly")
    ap.add_argument("--sma200", type=float)
    ap.add_argument("--sma200-file")
    ap.add_argument("--symbol", default="?")
    ap.add_argument("--out", required=True)
    ap.add_argument("--daily-bars", type=int, default=100)
    ap.add_argument("--theme", choices=sorted(THEMES), default="light")
    args = ap.parse_args(argv)

    daily = ta.load_bars(args.daily)
    weekly = ta.load_bars(args.weekly) if args.weekly else None
    sma200 = args.sma200
    if args.sma200_file and sma200 is None:
        sma200 = ta.latest_sma_from_file(args.sma200_file)
    facts = ta.analyse(daily, weekly, sma200)
    svg = build(args.symbol.upper(), daily, weekly, facts, args.theme, args.daily_bars)
    with open(args.out, "w", encoding="utf-8") as fh:
        fh.write(svg)
    print(args.out)


if __name__ == "__main__":
    main()
