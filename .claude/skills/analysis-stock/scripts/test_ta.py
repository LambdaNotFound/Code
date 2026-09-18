"""Tests for ta.py. Run from this directory: python3 -m unittest test_ta -v"""

import json
import os
import sys
import unittest

sys.path.insert(0, os.path.dirname(__file__))
import ta  # noqa: E402


def bars_from_closes(closes, spread=1.0, volume=100.0):
    return [{"date": f"2026-01-{i + 1:02d}", "open": c, "high": c + spread, "low": c - spread,
             "close": c, "volume": volume, "adjusted": False} for i, c in enumerate(closes)]


class ParseTests(unittest.TestCase):
    CSV = "timestamp,open,high,low,close,volume\r\n2026-01-03,3,4,2,3.5,30\r\n2026-01-02,2,3,1,2.5,20\r\n2026-01-01,1,2,0,1.5,10\r\n"

    def test_raw_csv_sorted_ascending(self):
        bars = ta.parse_bars(self.CSV)
        self.assertEqual([b["date"] for b in bars], ["2026-01-01", "2026-01-02", "2026-01-03"])
        self.assertEqual(bars[-1]["close"], 3.5)
        self.assertFalse(bars[0]["adjusted"])

    def test_mcp_offload_wrapper(self):
        bars = ta.parse_bars(json.dumps({"result": self.CSV}))
        self.assertEqual(len(bars), 3)

    def test_alpha_vantage_json(self):
        d = {"Meta Data": {}, "Time Series (Daily)": {
            "2026-01-02": {"1. open": "2", "2. high": "3", "3. low": "1", "4. close": "2.5", "5. volume": "20"},
            "2026-01-01": {"1. open": "1", "2. high": "2", "3. low": "0", "4. close": "1.5", "5. volume": "10"}}}
        bars = ta.parse_bars(json.dumps(d))
        self.assertEqual(bars[0]["date"], "2026-01-01")
        self.assertEqual(bars[1]["volume"], 20.0)

    def test_adjusted_close_scales_ohl(self):
        csv = "timestamp,open,high,low,close,adjusted close,volume,dividend amount\r\n2020-08-21,400,440,380,400,100,1,0\r\n"
        b = ta.parse_bars(csv)[0]
        self.assertTrue(b["adjusted"])
        self.assertEqual((b["open"], b["high"], b["low"], b["close"]), (100.0, 110.0, 95.0, 100.0))

    def test_error_payload_raises(self):
        with self.assertRaises(ValueError):
            ta.parse_bars(json.dumps({"error": {"type": "rate_limit"}}))


class IndicatorTests(unittest.TestCase):
    def test_sma(self):
        self.assertEqual(ta.sma([1, 2, 3, 4, 5], 3), [None, None, 2.0, 3.0, 4.0])

    def test_ema_of_constant_is_constant(self):
        out = ta.ema([5.0] * 30, 10)
        self.assertIsNone(out[8])
        self.assertTrue(all(abs(v - 5.0) < 1e-12 for v in out[9:]))

    def test_ema_first_step(self):
        out = ta.ema([1, 1, 1, 4], 3)  # seed 1.0, k = 0.5 -> 4*0.5 + 1*0.5
        self.assertAlmostEqual(out[3], 2.5)

    def test_rsi_all_gains_is_100_all_losses_is_0(self):
        up = list(range(1, 40))
        self.assertAlmostEqual(ta.rsi(up)[-1], 100.0)
        self.assertAlmostEqual(ta.rsi(up[::-1])[-1], 0.0)

    def test_rsi_alternating_stays_near_50(self):
        # +1, -1, +1 ... Wilder smoothing never lets the two averages equalise
        # exactly (the last change always tips one), but they stay close.
        closes = [10 + (i % 2) for i in range(40)]
        self.assertTrue(45.0 < ta.rsi(closes)[-1] < 55.0)

    def test_rsi_seed_matches_formula(self):
        # First RSI value is the plain average of the first 14 gains and losses.
        closes = [44.34, 44.09, 44.15, 43.61, 44.33, 44.83, 45.10, 45.42, 45.84, 46.08, 45.89,
                  46.03, 45.61, 46.28, 46.28]
        changes = [b - a for a, b in zip(closes, closes[1:])]
        ag = sum(c for c in changes if c > 0) / 14
        al = sum(-c for c in changes if c < 0) / 14
        self.assertAlmostEqual(ta.rsi(closes)[-1], 100 - 100 / (1 + ag / al))
        self.assertIsNone(ta.rsi(closes)[-2])

    def test_rsi_second_value_uses_wilder_smoothing(self):
        closes = [10.0] * 15 + [12.0, 11.0]  # flat seed window, then +2, -1
        r = ta.rsi(closes)
        # after +2: ag = 2/14, al = 0 -> 100. after -1: ag = (2/14)*(13/14), al = 1/14
        ag = (2 / 14) * (13 / 14)
        al = 1 / 14
        self.assertAlmostEqual(r[-2], 100.0)
        self.assertAlmostEqual(r[-1], 100 - 100 / (1 + ag / al))

    def test_macd_of_constant_is_zero(self):
        line, sig, hist = ta.macd([7.0] * 60)
        self.assertAlmostEqual(line[-1], 0.0)
        self.assertAlmostEqual(sig[-1], 0.0)
        self.assertAlmostEqual(hist[-1], 0.0)
        self.assertIsNone(line[24])
        self.assertIsNotNone(line[25])
        self.assertIsNone(sig[32])
        self.assertIsNotNone(sig[33])

    def test_bbands_constant_collapses(self):
        u, m, l = ta.bbands([3.0] * 25)
        self.assertEqual((u[-1], m[-1], l[-1]), (3.0, 3.0, 3.0))

    def test_bbands_width(self):
        closes = [1.0, 3.0] * 10  # mean 2, population sd 1 -> bands at 0 and 4
        u, m, l = ta.bbands(closes, 20, 2.0)
        self.assertAlmostEqual(u[-1], 4.0)
        self.assertAlmostEqual(l[-1], 0.0)

    def test_atr_constant_range(self):
        bars = bars_from_closes([10.0] * 30, spread=1.0)  # high-low = 2 every bar
        self.assertAlmostEqual(ta.atr(bars)[-1], 2.0)

    def test_true_range_uses_gap(self):
        bars = [{"high": 10, "low": 9, "close": 9.5}, {"high": 13, "low": 12, "close": 12.5}]
        self.assertEqual(ta.true_range(bars)[1], 13 - 9.5)

    def test_adx_straight_uptrend(self):
        bars = bars_from_closes([float(i) for i in range(1, 60)])
        a, p, m = ta.adx(bars)
        self.assertIsNotNone(a[-1])
        self.assertGreater(p[-1], m[-1])
        self.assertAlmostEqual(m[-1], 0.0)
        self.assertGreater(a[-1], 90.0)

    def test_adx_needs_two_periods(self):
        bars = bars_from_closes([1.0] * 20)
        self.assertEqual(ta.adx(bars)[0], [None] * 20)

    def test_stochastic_at_range_high(self):
        bars = bars_from_closes([float(i) for i in range(1, 30)], spread=0.0)
        k, d = ta.stochastic(bars)
        self.assertAlmostEqual(k[-1], 100.0)
        self.assertAlmostEqual(d[-1], 100.0)

    def test_obv(self):
        bars = bars_from_closes([1, 2, 2, 1, 3], volume=10)
        self.assertEqual(ta.obv(bars), [0, 10, 10, 0, 10])


class StructureTests(unittest.TestCase):
    ZIGZAG = [10, 11, 12, 13, 12, 11, 10, 9, 10, 11, 12, 13, 14, 13, 12, 11, 10, 11, 12, 13, 14, 15, 14, 13, 12, 11, 12]

    def test_pivots_find_swings(self):
        bars = bars_from_closes(self.ZIGZAG, spread=0.0)
        pivs = ta.pivots(bars, 3)
        kinds = [(p["kind"], p["price"]) for p in pivs]
        self.assertIn(("H", 13.0), kinds)
        self.assertIn(("L", 9.0), kinds)
        self.assertIn(("H", 14.0), kinds)
        self.assertIn(("L", 10.0), kinds)

    def test_trend_structure_uptrend(self):
        bars = bars_from_closes(self.ZIGZAG, spread=0.0)
        label, detail = ta.trend_structure(ta.pivots(bars, 3))
        self.assertEqual(label, "uptrend")
        self.assertIn("HH/HL", detail)

    def test_trend_structure_downtrend(self):
        bars = bars_from_closes([-x for x in self.ZIGZAG], spread=0.0)
        label, _ = ta.trend_structure(ta.pivots(bars, 3))
        self.assertEqual(label, "downtrend")

    def test_trend_structure_undetermined(self):
        self.assertEqual(ta.trend_structure([])[0], "undetermined")

    def test_cluster_levels(self):
        pivs = [{"price": 100.0, "date": "a"}, {"price": 100.5, "date": "b"}, {"price": 120.0, "date": "c"}, {"price": 90.0, "date": "d"}]
        sup, res = ta.cluster_levels(pivs, 105.0)
        self.assertEqual([round(c["price"], 2) for c in sup], [100.25, 90.0])
        self.assertEqual(sup[0]["touches"], 2)
        self.assertEqual([c["price"] for c in res], [120.0])

    def test_last_cross(self):
        fast = [1, 2, 3, 4, 5]
        slow = [3, 3, 3, 3, 3]
        self.assertEqual(ta.last_cross(fast, slow, list("abcde")), ("bullish", "d", 1))
        self.assertIsNone(ta.last_cross([5, 5, 5], [3, 3, 3], list("abc")))

    def test_gaps(self):
        closes = [100.0] * 5 + [110.0] * 5
        bars = bars_from_closes(closes, spread=0.5)
        bars[5]["open"] = 110.0
        g = ta.gaps(bars, lookback=20, min_pct=2.0)
        self.assertEqual(len(g), 1)
        self.assertFalse(g[0]["filled"])
        self.assertAlmostEqual(g[0]["pct"], 10.0)

    def test_rsi_divergence_bearish(self):
        # Price makes a higher high while RSI at those two pivots is lower.
        bars = bars_from_closes(self.ZIGZAG, spread=0.0)
        pivs = ta.pivots(bars, 3)
        fake_rsi = [50.0] * len(bars)
        highs = [p for p in pivs if p["kind"] == "H"]
        fake_rsi[highs[-2]["index"]] = 80.0
        fake_rsi[highs[-1]["index"]] = 70.0
        found = ta.rsi_divergence(bars, pivs, fake_rsi)
        self.assertTrue(any(f.startswith("bearish") for f in found))


class AnalyseTests(unittest.TestCase):
    def test_end_to_end_daily_only(self):
        closes = [100 + i * 0.5 + (i % 3) for i in range(100)]
        facts = ta.analyse(bars_from_closes(closes))
        self.assertEqual(facts["bars"], 100)
        self.assertEqual(facts["ma"]["sma200_source"], "unavailable")
        self.assertIsNone(facts["weekly"])
        t = facts["ledger"]["tally"]
        self.assertEqual(sum(t.values()), len(facts["ledger"]["rows"]))
        text = ta.render(facts, "TEST")
        self.assertIn("# TA facts: TEST", text)
        self.assertIn("## Signal ledger", text)

    def test_sma200_override_and_weekly(self):
        daily = bars_from_closes([100.0 + i for i in range(100)])
        weekly = bars_from_closes([50.0 + i for i in range(260)])
        facts = ta.analyse(daily, weekly, sma200_daily=150.0)
        self.assertEqual(facts["ma"]["sma200"], 150.0)
        self.assertEqual(facts["ma"]["sma200_source"], "daily")
        self.assertFalse(facts["weekly"]["adjusted"])
        self.assertIsNone(facts["weekly"]["all_time_high"])
        self.assertIsNotNone(facts["weekly"]["ma"]["sma200"])

    def test_sma200_falls_back_to_40_week(self):
        daily = bars_from_closes([100.0] * 100)
        weekly = bars_from_closes([80.0] * 60)
        facts = ta.analyse(daily, weekly)
        self.assertEqual(facts["ma"]["sma200_source"], "approx-40-week")
        self.assertAlmostEqual(facts["ma"]["sma200"], 80.0)


if __name__ == "__main__":
    unittest.main()
