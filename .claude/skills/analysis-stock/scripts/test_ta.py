"""Tests for ta.py. Run from this directory: python3 -m unittest test_ta -v"""

import json
import math
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
        self.assertIn(facts["stance"]["label"], ("BUY", "ACCUMULATE", "HOLD", "REDUCE", "SELL"))
        text = ta.render(facts, "TEST")
        self.assertIn("# TA facts: TEST", text)
        self.assertIn("## Signal ledger", text)
        self.assertIn("## Stance", text)

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


class StanceTests(unittest.TestCase):
    def facts(self, reads, price=100.0, atr=2.0, sup=(90.0,), res=(120.0,)):
        rows = [{"indicator": name, "value": "", "read": read, "rule": ""} for name, read in reads]
        return {"price": price, "volatility": {"atr14": atr}, "weekly": None,
                "levels": {"supports": [{"price": p} for p in sup], "resistances": [{"price": p} for p in res]},
                "ledger": {"rows": rows}}

    def test_all_bullish_is_buy_with_plan(self):
        f = self.facts([("MA stack (20/50/200)", "bullish"), ("ADX14 / DI", "bullish"), ("RSI14", "neutral")])
        st = ta.stance(f)
        self.assertEqual(st["label"], "BUY")
        self.assertAlmostEqual(st["score"], 4 / 5)
        self.assertEqual(st["plan"]["entry"], 100.0)
        self.assertEqual(st["plan"]["stop"], 89.0)      # support minus half an ATR
        self.assertEqual(st["plan"]["target"], 120.0)
        self.assertAlmostEqual(st["plan"]["reward_risk"], 20 / 11)

    def test_poor_reward_risk_demotes_to_accumulate_at_support(self):
        f = self.facts([("MA stack (20/50/200)", "bullish")], price=118.0, sup=(90.0,), res=(120.0,))
        st = ta.stance(f)
        self.assertEqual(st["label"], "ACCUMULATE")
        self.assertEqual(st["plan"]["entry"], 90.0)
        self.assertIn("reward/risk", st["plan"]["note"])

    def test_all_bearish_is_sell(self):
        f = self.facts([("MA stack (20/50/200)", "bearish"), ("MACD histogram", "bearish")])
        st = ta.stance(f)
        self.assertEqual(st["label"], "SELL")
        self.assertEqual(st["plan"]["stop"], 121.0)
        self.assertEqual(st["plan"]["target"], 90.0)

    def test_mixed_is_hold_with_no_plan(self):
        f = self.facts([("MA stack (20/50/200)", "bullish"), ("Swing structure (daily)", "bearish"), ("RSI14", "neutral")])
        st = ta.stance(f)
        self.assertEqual(st["label"], "HOLD")
        self.assertEqual(st["score"], 0.0)
        self.assertIsNone(st["plan"]["entry"])
        self.assertIn("120.00", st["plan"]["note"])

    def test_band_floor_is_inclusive(self):
        # (+2 - 2 + 1) / 5 = 0.20 exactly -> ACCUMULATE, not HOLD
        f = self.facts([("MA stack (20/50/200)", "bullish"), ("Swing structure (daily)", "bearish"), ("RSI14", "bullish")])
        self.assertEqual(ta.stance(f)["label"], "ACCUMULATE")

    def test_unknown_row_does_not_score(self):
        f = self.facts([("Something new", "bullish")])
        self.assertEqual(ta.stance(f)["label"], "HOLD")

    def test_no_resistance_means_open_target(self):
        f = self.facts([("MA stack (20/50/200)", "bullish")], res=())
        st = ta.stance(f)
        self.assertEqual(st["label"], "BUY")
        self.assertIsNone(st["plan"]["target"])
        self.assertIn("open-ended", st["plan"]["note"])


class ReferenceMethodTests(unittest.TestCase):
    def test_trend_template_all_pass(self):
        tt = ta.trend_template(100, 90, 85, 80, 78, 60, 110, ["x"])
        self.assertEqual((tt["passed"], tt["of"]), (7, 7))

    def test_trend_template_partial(self):
        # price under the 50-day and far from the high, 200-day flat
        tt = ta.trend_template(100, 105, 95, 90, 90, 60, 150, [])
        self.assertEqual(tt["passed"], 4)
        self.assertFalse(tt["checks"][2]["pass"])
        self.assertFalse(tt["checks"][4]["pass"])
        self.assertFalse(tt["checks"][6]["pass"])

    def test_trend_template_needs_inputs(self):
        self.assertIsNone(ta.trend_template(100, None, 85, 80, 78, 60, 110, []))

    @staticmethod
    def weekly(n=70):
        """Flat 100 closes, highs 105, lows 95, one deep low at week 20 so the
        52-week low (50) sits well under any week's low."""
        bars = bars_from_closes([100.0] * n, spread=5.0)
        for i, b in enumerate(bars):
            b["date"] = f"W{i:04d}"
        bars[20]["low"] = 50.0
        return bars

    def test_weekly_key_reversal_bearish(self):
        bars = self.weekly()
        # week 66: new 52w high intraweek, close under prior week's low (95); week 69 is in progress
        bars[66]["high"], bars[66]["close"] = 120.0, 94.0
        wr = ta.weekly_reversals(bars)
        self.assertEqual(wr["read"], "bearish")
        self.assertEqual([x["check"] for x in wr["bearish"]], ["key_reversal"])
        self.assertEqual(wr["bearish"][0]["week"], "W0066")
        self.assertEqual(wr["bullish"], [])

    def test_weekly_failed_extreme_bullish(self):
        bars = self.weekly()
        bars[65]["low"], bars[65]["close"] = 40.0, 100.5  # under the 52w low (50), closes back above it
        wr = ta.weekly_reversals(bars)
        self.assertEqual(wr["read"], "bullish")
        self.assertEqual(wr["bullish"][0]["check"], "failed_extreme")
        self.assertEqual(wr["bullish"][0]["week"], "W0065")

    def test_weekly_failed_breakout_dated_on_failure_week(self):
        bars = self.weekly()
        bars[62]["close"], bars[62]["high"] = 106.0, 107.0  # closing breakout above 105
        # week 63 closes 100 < 105: the failure week is 63
        wr = ta.weekly_reversals(bars)
        fb = [x for x in wr["bearish"] if x["check"] == "failed_breakout"]
        self.assertEqual(len(fb), 1)
        self.assertEqual(fb[0]["week"], "W0063")
        self.assertEqual(fb[0]["level"], 105.0)
        self.assertEqual(wr["read"], "bearish")

    def test_weekly_continuation_veto(self):
        bars = self.weekly()
        bars[62]["high"], bars[62]["close"] = 120.0, 94.0   # bearish key reversal
        for i, c in ((66, 125.0), (67, 126.0), (68, 127.0)):  # then closes at new highs and stay there
            bars[i]["close"], bars[i]["high"], bars[i]["low"] = c, c + 1, c - 1
        wr = ta.weekly_reversals(bars)
        self.assertEqual([x["check"] for x in wr["bearish"]], ["key_reversal"])
        self.assertTrue(wr["bearish_vetoed"])
        self.assertEqual(wr["read"], "neutral")

    def test_weekly_reversals_insufficient(self):
        self.assertTrue(ta.weekly_reversals(self.weekly(30))["insufficient"])

    def test_burst_days(self):
        closes = [100.0] * 25 + [105.0, 105.5]
        bars = bars_from_closes(closes, spread=0.5, volume=100.0)
        bars[25]["volume"] = 300.0
        bars[25]["high"], bars[25]["low"] = 106.0, 100.0
        out = ta.burst_days(bars)
        self.assertEqual(out[0]["date"], bars[25]["date"])
        self.assertIn("4pct_breakout", out[0]["tags"])
        self.assertIn("range_expansion", out[0]["tags"])
        self.assertAlmostEqual(out[0]["pct"], 5.0)

    def test_position_size_risk_budget_binds(self):
        ps = ta.position_size({"entry": 100.0, "stop": 95.0}, 100000, 1.0, 50.0)
        self.assertEqual(ps["shares"], 200)
        self.assertEqual(ps["binding"], "risk budget")
        self.assertAlmostEqual(ps["actual_risk"], 1000.0)

    def test_position_size_cap_binds(self):
        ps = ta.position_size({"entry": 100.0, "stop": 99.0}, 100000, 1.0, 10.0)
        self.assertEqual(ps["shares"], 100)
        self.assertEqual(ps["binding"], "10% position cap")

    def test_position_size_absent_without_account(self):
        self.assertIsNone(ta.position_size({"entry": 100.0, "stop": 95.0}, None, 1.0, 10.0))
        self.assertIsNone(ta.position_size({"entry": None, "stop": None}, 1000, 1.0, 10.0))

    def test_analyse_carries_new_blocks(self):
        daily = bars_from_closes([100.0 + i * 0.3 for i in range(100)])
        weekly = bars_from_closes([50.0 + i * 0.5 for i in range(120)])
        facts = ta.analyse(daily, weekly, sma200_daily=90.0, account=50000)
        self.assertIsNotNone(facts["trend_template"])
        self.assertFalse(facts["weekly_reversals"]["insufficient"])
        self.assertIn("Trend template (Minervini)", [r["indicator"] for r in facts["ledger"]["rows"]])
        self.assertIsNotNone(facts["stance"]["plan"]["target_2r"])
        text = ta.render(facts, "T")
        self.assertIn("## Trend template", text)
        self.assertIn("## Weekly reversal checks", text)
        self.assertIn("## Position size", text)


class SessionLessonTests(unittest.TestCase):
    def facts(self, reads, price=100.0, atr=2.0, sup=(90.0,), res=(120.0,)):
        rows = [{"indicator": name, "value": "", "read": read, "rule": ""} for name, read in reads]
        return {"price": price, "volatility": {"atr14": atr}, "weekly": None,
                "levels": {"supports": [{"price": p} for p in sup], "resistances": [{"price": p} for p in res]},
                "ledger": {"rows": rows}}

    def test_short_side_poor_reward_risk_moves_exit_to_resistance(self):
        # close 92, support 90, resistance 120: from the close the short makes 2 against 29 of risk
        f = self.facts([("MA stack (20/50/200)", "bearish")], price=92.0, sup=(90.0,), res=(120.0,))
        st = ta.stance(f)
        self.assertEqual(st["label"], "SELL")
        self.assertEqual(st["plan"]["entry"], 120.0)
        self.assertEqual(st["plan"]["stop"], 121.0)
        self.assertEqual(st["plan"]["target"], 90.0)
        self.assertAlmostEqual(st["plan"]["reward_risk"], 30.0)
        self.assertIn("trim into strength", st["plan"]["note"])

    def test_price_vs_200_weighs_one(self):
        f = self.facts([("Price vs 200-day", "bullish"), ("RSI14", "bearish")])
        self.assertEqual(ta.stance(f)["score"], 0.0)

    def test_latest_sma_from_every_shape(self):
        import tempfile
        av = {"Meta Data": {}, "Technical Analysis: SMA": {"2026-09-16": {"SMA": "1.0"}, "2026-09-17": {"SMA": "2.5"}}}
        shapes = {
            "raw.json": json.dumps(av),
            "wrapped.txt": json.dumps({"result": json.dumps(av)}),
            "preview.json": json.dumps({"preview": True, "sample_data": json.dumps(av)}),
            "plain.csv": "time,SMA\r\n2026-09-17,2.5\r\n2026-09-16,1.0\r\n",
            "wrapped_csv.txt": json.dumps({"result": "time,SMA\r\n2026-09-17,2.5\r\n"}),
        }
        with tempfile.TemporaryDirectory() as d:
            for name, body in shapes.items():
                p = os.path.join(d, name)
                with open(p, "w") as fh:
                    fh.write(body)
                self.assertEqual(ta.latest_sma_from_file(p), 2.5, name)
            p = os.path.join(d, "err.json")
            with open(p, "w") as fh:
                fh.write(json.dumps({"error": {"type": "rate_limit"}}))
            with self.assertRaises(ValueError):
                ta.latest_sma_from_file(p)


class MonthlyAndMoveTests(unittest.TestCase):
    def weekly_series(self, n=160, start=100.0, step=0.5, amp=0.0, period=40):
        """One bar per week from 2020-01-03, so 4-5 land in each calendar month.
        amp adds a cycle so the monthly bars have swing points to detect; a purely
        monotonic series has no local highs or lows and reads as undetermined."""
        import datetime
        d0 = datetime.date(2020, 1, 3)
        out = []
        for i in range(n):
            d = d0 + datetime.timedelta(weeks=i)
            c = start + i * step + amp * math.sin(2 * math.pi * i / period)
            out.append({"date": d.isoformat(), "open": c - 0.5, "high": c + 1.0, "low": c - 1.0,
                        "close": c, "volume": 100.0, "adjusted": True})
        return out

    def test_to_monthly_aggregates_ohlcv(self):
        bars = [
            {"date": "2026-01-05", "open": 10, "high": 12, "low": 9, "close": 11, "volume": 5, "adjusted": True},
            {"date": "2026-01-26", "open": 11, "high": 15, "low": 8, "close": 14, "volume": 7, "adjusted": True},
            {"date": "2026-02-02", "open": 14, "high": 16, "low": 13, "close": 15, "volume": 3, "adjusted": True},
        ]
        m = ta.to_monthly(bars)
        self.assertEqual([b["date"] for b in m], ["2026-01", "2026-02"])
        self.assertEqual((m[0]["open"], m[0]["high"], m[0]["low"], m[0]["close"], m[0]["volume"]), (10, 15, 8, 14, 12))
        self.assertEqual(m[1]["close"], 15)

    def test_to_monthly_adjusted_flag_is_and(self):
        bars = [{"date": "2026-01-05", "open": 1, "high": 1, "low": 1, "close": 1, "volume": 1, "adjusted": True},
                {"date": "2026-01-12", "open": 1, "high": 1, "low": 1, "close": 1, "volume": 1, "adjusted": False}]
        self.assertFalse(ta.to_monthly(bars)[0]["adjusted"])

    def test_analyse_monthly_uptrend(self):
        # drift per cycle (40) exceeds the swing amplitude (30 peak to trough), so
        # each swing high and each swing low is above the one before it
        mo = ta.analyse_monthly(self.weekly_series(step=1.0, amp=15.0))
        self.assertFalse(mo["insufficient"])
        self.assertGreater(mo["bars"], 30)
        self.assertEqual(mo["structure"]["label"], "uptrend")
        self.assertIn("HH/HL", mo["structure"]["detail"])
        self.assertGreater(mo["ma"]["sma12"], mo["ma"]["sma24"])
        self.assertEqual(mo["ma"]["sma12_slope"], "rising")
        self.assertGreater(mo["range_all"]["pct_from_low"], 0.0)

    def test_analyse_monthly_downtrend(self):
        mo = ta.analyse_monthly(self.weekly_series(start=300.0, step=-1.0, amp=15.0))
        self.assertEqual(mo["structure"]["label"], "downtrend")

    def test_monotonic_series_has_no_swings(self):
        mo = ta.analyse_monthly(self.weekly_series())
        self.assertEqual(mo["structure"]["label"], "undetermined")

    def test_analyse_monthly_insufficient(self):
        self.assertTrue(ta.analyse_monthly(self.weekly_series(n=20))["insufficient"])

    def test_realized_vol_matches_formula(self):
        closes = [100.0]
        for i in range(40):
            closes.append(closes[-1] * math.exp(0.01 if i % 2 == 0 else -0.01))
        self.assertAlmostEqual(ta.realized_vol(closes, 20), 0.01 * math.sqrt(20 / 19) * math.sqrt(252), places=9)
        self.assertAlmostEqual(ta.realized_vol([7.0] * 30, 20), 0.0)
        self.assertIsNone(ta.realized_vol([1.0, 2.0], 20))

    def test_third_friday_known_dates(self):
        import datetime
        self.assertEqual(ta.third_friday(2026, 9), datetime.date(2026, 9, 18))
        self.assertEqual(ta.third_friday(2026, 10), datetime.date(2026, 10, 16))
        self.assertEqual(ta.third_friday(2026, 1), datetime.date(2026, 1, 16))
        for y in range(2024, 2030):
            for m in range(1, 13):
                d = ta.third_friday(y, m)
                self.assertEqual(d.weekday(), 4)
                self.assertTrue(15 <= d.day <= 21)

    def test_next_expiries_are_future_and_ordered(self):
        import datetime
        out = ta.next_expiries(datetime.date(2026, 9, 17), 3)
        self.assertEqual([d.isoformat() for d in out], ["2026-09-18", "2026-10-16", "2026-11-20"])
        after = ta.next_expiries(datetime.date(2026, 12, 20), 2)
        self.assertEqual([d.isoformat() for d in after], ["2027-01-15", "2027-02-19"])

    def test_expected_moves_scale_with_root_time(self):
        import datetime
        em = ta.expected_moves(100.0, 0.40, datetime.date(2026, 9, 17))
        by = {r["label"]: r for r in em["horizons"]}
        self.assertAlmostEqual(by["30d"]["sigma_pct"], 100 * 0.40 * math.sqrt(30 / 365), places=9)
        # doubling the horizon multiplies sigma by sqrt(2)
        self.assertAlmostEqual(by["60d"]["sigma_pct"] / by["30d"]["sigma_pct"], math.sqrt(2), places=9)
        self.assertLess(by["30d"]["low"], 100.0)
        self.assertGreater(by["30d"]["high"], 100.0)
        self.assertLess(by["30d"]["low_2s"], by["30d"]["low"])
        self.assertEqual(len(em["expiries"]), 3)

    def test_expected_moves_none_without_vol(self):
        import datetime
        self.assertIsNone(ta.expected_moves(100.0, None, datetime.date(2026, 9, 17)))

    def test_analyse_carries_monthly_and_expected_move(self):
        daily = bars_from_closes([100.0 + i * 0.3 for i in range(100)])
        for i, b in enumerate(daily):
            b["date"] = f"2026-{1 + i // 28:02d}-{1 + i % 28:02d}"
        facts = ta.analyse(daily, self.weekly_series(), sma200_daily=90.0)
        self.assertFalse(facts["monthly"]["insufficient"])
        self.assertIsNotNone(facts["expected_move"])
        self.assertIn("Monthly structure", [r["indicator"] for r in facts["ledger"]["rows"]])
        text = ta.render(facts, "T")
        self.assertIn("## Monthly", text)
        self.assertIn("## Expected move", text)


if __name__ == "__main__":
    unittest.main()
