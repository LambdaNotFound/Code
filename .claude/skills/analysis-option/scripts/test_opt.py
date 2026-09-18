"""Tests for opt.py. Run from this directory: python3 -m unittest test_opt -v"""

import json
import math
import os
import sys
import tempfile
import unittest

sys.path.insert(0, os.path.dirname(__file__))
import opt  # noqa: E402


def bars(closes, spread=1.0):
    return [{"date": f"D{i:04d}", "high": c + spread, "low": c - spread, "close": c} for i, c in enumerate(closes)]


class BlackScholesTests(unittest.TestCase):
    def test_textbook_call_and_put(self):
        # S=100, K=100, r=5%, sigma=20%, T=1: call 10.4506, put 5.5735 (Hull's standard example).
        self.assertAlmostEqual(opt.bs_price("call", 100, 100, 1.0, 0.05, 0.2), 10.4506, places=3)
        self.assertAlmostEqual(opt.bs_price("put", 100, 100, 1.0, 0.05, 0.2), 5.5735, places=3)

    def test_put_call_parity(self):
        s, k, t, r, q, sig = 120.0, 110.0, 0.4, 0.03, 0.01, 0.35
        c = opt.bs_price("call", s, k, t, r, sig, q)
        p = opt.bs_price("put", s, k, t, r, sig, q)
        self.assertAlmostEqual(c - p, s * math.exp(-q * t) - k * math.exp(-r * t), places=9)

    def test_expiry_is_intrinsic(self):
        self.assertEqual(opt.bs_price("call", 105, 100, 0.0, 0.05, 0.2), 5.0)
        self.assertEqual(opt.bs_price("put", 105, 100, 0.0, 0.05, 0.2), 0.0)

    def test_greeks_signs_and_delta_bounds(self):
        g = opt.bs_greeks("call", 100, 100, 0.5, 0.04, 0.3)
        self.assertTrue(0.5 < g["delta"] < 0.7)
        self.assertGreater(g["gamma"], 0)
        self.assertLess(g["theta"], 0)
        self.assertGreater(g["vega"], 0)
        gp = opt.bs_greeks("put", 100, 100, 0.5, 0.04, 0.3)
        self.assertAlmostEqual(g["delta"] - gp["delta"], 1.0, places=9)
        self.assertAlmostEqual(g["gamma"], gp["gamma"], places=12)

    def test_delta_by_finite_difference(self):
        s, k, t, r, sig = 100.0, 95.0, 0.25, 0.04, 0.25
        h = 1e-3
        fd = (opt.bs_price("call", s + h, k, t, r, sig) - opt.bs_price("call", s - h, k, t, r, sig)) / (2 * h)
        self.assertAlmostEqual(opt.bs_greeks("call", s, k, t, r, sig)["delta"], fd, places=5)

    def test_implied_vol_round_trip(self):
        price = opt.bs_price("put", 100, 90, 0.3, 0.04, 0.42)
        self.assertAlmostEqual(opt.implied_vol("put", price, 100, 90, 0.3, 0.04), 0.42, places=5)

    def test_implied_vol_below_intrinsic_is_none(self):
        self.assertIsNone(opt.implied_vol("call", 1.0, 120, 100, 0.5, 0.04))


class UnderlyingTests(unittest.TestCase):
    def test_realized_vol_constant_is_zero(self):
        self.assertAlmostEqual(opt.realized_vol([50.0] * 40, 20), 0.0)

    def test_realized_vol_known_series(self):
        # alternating +1%/-1% log returns: sample sd of returns ~ 0.01, annualised ~ 0.1587
        closes = [100.0]
        for i in range(40):
            closes.append(closes[-1] * math.exp(0.01 if i % 2 == 0 else -0.01))
        hv = opt.realized_vol(closes, 20)
        self.assertAlmostEqual(hv, 0.01 * math.sqrt(20 / 19) * math.sqrt(252), places=6)

    def test_atr_constant_range(self):
        self.assertAlmostEqual(opt.atr(bars([10.0] * 30, spread=1.0)), 2.0)

    def test_load_bars_three_shapes(self):
        csv_text = "timestamp,open,high,low,close,volume\r\n2026-01-02,2,3,1,2.5,20\r\n2026-01-01,1,2,0,1.5,10\r\n"
        with tempfile.TemporaryDirectory() as d:
            p1, p2, p3 = (os.path.join(d, n) for n in ("a.csv", "b.txt", "c.json"))
            with open(p1, "w") as fh:
                fh.write(csv_text)
            with open(p2, "w") as fh:
                fh.write(json.dumps({"result": csv_text}))
            with open(p3, "w") as fh:
                fh.write(json.dumps({"Time Series (Daily)": {
                    "2026-01-02": {"1. open": "2", "2. high": "3", "3. low": "1", "4. close": "2.5", "5. volume": "20"}}}))
            self.assertEqual([b["date"] for b in opt.load_bars(p1)], ["2026-01-01", "2026-01-02"])
            self.assertEqual(len(opt.load_bars(p2)), 2)
            self.assertEqual(opt.load_bars(p3)[0]["close"], 2.5)


class PositionTests(unittest.TestCase):
    def test_payoffs_and_breakevens(self):
        self.assertEqual(opt.payoff_at_expiry("call", "buy", 100, 5, 110), 5)
        self.assertEqual(opt.payoff_at_expiry("call", "buy", 100, 5, 90), -5)
        self.assertEqual(opt.payoff_at_expiry("put", "sell", 100, 5, 80), -15)
        self.assertEqual(opt.payoff_at_expiry("put", "sell", 100, 5, 120), 5)
        self.assertEqual(opt.payoff_at_expiry("call", "sell", 100, 5, 130), -25)
        self.assertEqual(opt.breakeven("call", 100, 5), 105)
        self.assertEqual(opt.breakeven("put", 100, 5), 95)

    def test_distribution_probabilities_sum_and_signs(self):
        st = opt.distribution_stats("call", "buy", 100, 3.0, 100, 0.3, 0.25)
        self.assertTrue(0 < st["p_profit"] < 1)
        self.assertGreater(st["avg_gain"], 0)
        self.assertLessEqual(st["avg_loss"], 0)
        deep = opt.distribution_stats("call", "buy", 50, 50.5, 100, 0.3, 0.25)
        self.assertGreater(deep["p_profit"], 0.4)

    def test_long_and_short_ev_mirror(self):
        long_ = opt.distribution_stats("put", "buy", 100, 4.0, 100, 0.3, 0.25)
        short = opt.distribution_stats("put", "sell", 100, 4.0, 100, 0.3, 0.25)
        self.assertAlmostEqual(long_["ev"], -short["ev"], places=9)
        self.assertAlmostEqual(long_["p_profit"] + short["p_profit"], 1.0, places=6)

    def test_scenario_table_sorted_and_includes_levels(self):
        rows = opt.scenario_table("call", "buy", 100, 5, 100, 0.3, 0.25, levels=[("support 90.00", 90.0)])
        prices = [r["price"] for r in rows]
        self.assertEqual(prices, sorted(prices))
        self.assertIn("support 90.00", [r["scenario"] for r in rows])


class InputTests(unittest.TestCase):
    def test_parse_quote(self):
        q = opt.parse_quote("bid=5.2, ask=5.4,iv=0.31,oi=1200")
        self.assertEqual((q["bid"], q["ask"], q["iv"], q["oi"]), (5.2, 5.4, 0.31, 1200.0))
        with self.assertRaises(ValueError):
            opt.parse_quote("bid=5")

    def test_chain_pick_and_context(self):
        hdr = "contractID,symbol,expiration,strike,type,last,mark,bid,bid_size,ask,ask_size,volume,open_interest,date,implied_volatility,delta,gamma,theta,vega,rho\r\n"
        rows = ""
        for k, iv in ((95, 0.35), (100, 0.32), (105, 0.30)):
            rows += f"AAPL{k},AAPL,2026-10-16,{k},call,5,5,4.9,10,5.1,10,300,1200,2026-09-17,{iv},0.5,0.01,-0.05,0.1,0.05\r\n"
        rows += "AAPLP,AAPL,2026-10-16,100,put,3,3,2.9,10,3.1,10,50,400,2026-09-17,0.33,-0.5,0.01,-0.05,0.1,-0.05\r\n"
        with tempfile.TemporaryDirectory() as d:
            p = os.path.join(d, "chain.csv")
            with open(p, "w") as fh:
                fh.write(hdr + rows)
            chain = opt.load_chain(p)
        q = opt.pick_contract(chain, "2026-10-16", 100, "call")
        self.assertEqual((q["bid"], q["ask"], q["iv"], q["oi"]), (4.9, 5.1, 0.32, 1200.0))
        ctx = opt.chain_context(chain, "2026-10-16", "call", 100)
        self.assertEqual([c["strike"] for c in ctx], [95.0, 100.0, 105.0])
        with self.assertRaises(ValueError):
            opt.pick_contract(chain, "2026-10-16", 110, "call")


class AnalyseTests(unittest.TestCase):
    def daily(self):
        closes = [100 * math.exp(0.012 * math.sin(i / 3.0)) for i in range(100)]
        b = bars(closes, spread=1.0)
        for i, x in enumerate(b):
            x["date"] = f"2026-{1 + i // 28:02d}-{1 + i % 28:02d}"
        return b

    def test_long_call_end_to_end(self):
        f = opt.analyse("T", "2026-06-01", 100.0, "call", "buy", self.daily(), {"bid": 2.0, "ask": 2.2, "oi": 800}, asof="2026-04-01")
        self.assertEqual(f["dte"], 61)
        self.assertEqual(f["quote"]["premium_used"], 2.2)
        self.assertEqual(f["pricing"]["iv_source"], "solved from mid")
        self.assertAlmostEqual(f["move"]["break_even"], 102.2)
        self.assertAlmostEqual(f["risk"]["max_loss"], 220.0)
        self.assertIsNone(f["risk"]["max_gain"])
        self.assertIn(f["verdict"]["label"], ("GOOD", "FAIR", "POOR"))
        text = opt.render(f)
        self.assertIn("## Verdict", text)
        self.assertIn("## Ledger", text)

    def test_naked_short_call_is_hard_fail(self):
        f = opt.analyse("T", "2026-06-01", 110.0, "call", "sell", self.daily(), {"bid": 1.0, "ask": 1.2}, asof="2026-04-01")
        self.assertEqual(f["verdict"]["label"], "POOR")
        self.assertTrue(any("naked" in h for h in f["verdict"]["hard_fails"]))
        self.assertIsNone(f["risk"]["max_loss"])

    def test_cash_secured_put_capital(self):
        f = opt.analyse("T", "2026-06-01", 90.0, "put", "sell", self.daily(), {"bid": 1.5, "ask": 1.7}, asof="2026-04-01", contracts=2)
        self.assertAlmostEqual(f["risk"]["capital_at_risk"], 90 * 200)
        self.assertAlmostEqual(f["risk"]["max_gain"], 1.5 * 200)
        self.assertAlmostEqual(f["risk"]["max_loss"], (90 - 1.5) * 200)

    def test_stance_alignment_and_earnings_rows(self):
        ta = {"stance": {"label": "SELL"}, "levels": {"supports": [{"price": 95.0}], "resistances": [{"price": 105.0}]}}
        f = opt.analyse("T", "2026-06-01", 100.0, "call", "buy", self.daily(), {"bid": 2.0, "ask": 2.2}, asof="2026-04-01",
                        ta=ta, earnings="2026-05-01")
        rows = {r["check"]: r["read"] for r in f["ledger"]["rows"]}
        self.assertEqual(rows["Alignment with analysis-stock stance"], "bad")
        self.assertEqual(rows["Earnings inside the window"], "neutral")
        self.assertIn("support 95.00", [r["scenario"] for r in f["scenarios"]])
        g = opt.analyse("T", "2026-06-01", 100.0, "put", "sell", self.daily(), {"bid": 2.0, "ask": 2.2}, asof="2026-04-01", earnings="2026-05-01")
        self.assertEqual({r["check"]: r["read"] for r in g["ledger"]["rows"]}["Earnings inside the window"], "bad")

    def test_stale_quote_has_no_iv_and_fails(self):
        # a call quoted below intrinsic cannot be solved
        f = opt.analyse("T", "2026-06-01", 50.0, "call", "buy", self.daily(), {"bid": 1.0, "ask": 1.2}, asof="2026-04-01")
        self.assertIsNone(f["pricing"]["iv"])
        self.assertEqual(f["verdict"]["label"], "POOR")

    def test_expired_raises(self):
        with self.assertRaises(ValueError):
            opt.analyse("T", "2026-03-01", 100.0, "call", "buy", self.daily(), {"bid": 1.0, "ask": 1.2}, asof="2026-04-01")


if __name__ == "__main__":
    unittest.main()
