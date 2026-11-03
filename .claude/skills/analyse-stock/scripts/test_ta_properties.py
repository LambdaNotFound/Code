"""Properties, edge cases, the command line, the report's text, and golden outputs.

- Invariants that must hold on any data, checked on 20 seeded markets (4 seeds
  x 5 regimes): levels sit on the right side of the price and nearest first,
  cost shares sum to 1 with nested ranges, outlook quantiles are ordered, a
  close beyond a level never beats a touch of it, odds only grow with the
  horizon, and the same input gives the same output.
- Edge cases: flat prices, very short histories, no weekly file.
- The command line: every accepted input shape gives the same JSON, the flags
  do what they say, and errors exit cleanly.
- The markdown report: every section is there, and the numbers it prints are
  the JSON's numbers.
- Golden outputs: three markets' key outputs are pinned in golden/*.json; a
  change to ta.py that moves them fails here. Re-pin on purpose with
  UPDATE_GOLDEN=1 python3 -m unittest test_ta_properties.

Run from this directory: python3 -m unittest test_ta_properties -v
"""
import contextlib
import io
import json
import math
import os
import sys
import tempfile
import unittest

HERE = os.path.dirname(os.path.abspath(__file__))
sys.path.insert(0, HERE)
import ta  # noqa: E402
from market import REGIMES, market, to_av_json, to_csv, weekdays_ending  # noqa: E402

GOLDEN_DIR = os.path.join(HERE, "golden")
GOLDEN = {"uptrend-11": (11, "uptrend"), "downtrend-17": (17, "downtrend"), "range-23": (23, "range")}
PATHS = 2000  # outlook paths in these tests: enough for the properties, fast to run


def full(daily, weekly, **kw):
    f = ta.analyse(daily, weekly, **kw)
    f["outlook"] = ta.outlook(f, daily, weekly, paths=PATHS)
    return f


def run_cli(args):
    out, err = io.StringIO(), io.StringIO()
    with contextlib.redirect_stdout(out), contextlib.redirect_stderr(err):
        try:
            ta.main(args)
            code = 0
        except SystemExit as e:
            code = e.code
    return code, out.getvalue(), err.getvalue()


class Invariants(unittest.TestCase):
    MARKETS = [(seed, regime, *market(seed, regime)) for seed in (31, 32, 33, 34) for regime in REGIMES]

    def test_levels(self):
        for seed, regime, daily, weekly in self.MARKETS:
            with self.subTest(seed=seed, regime=regime):
                f = ta.analyse(daily, weekly)
                p = f["price"]
                for lv, side, cap in ((f["levels"], "daily", 4), (f["weekly"]["levels"], "weekly", 3)):
                    sup, res = lv["supports"], lv["resistances"]
                    self.assertLessEqual(len(sup), cap)
                    self.assertLessEqual(len(res), cap)
                    self.assertTrue(all(x["price"] < p for x in sup), side)
                    self.assertTrue(all(x["price"] > p for x in res), side)
                    self.assertEqual([x["price"] for x in sup], sorted((x["price"] for x in sup), reverse=True))
                    self.assertEqual([x["price"] for x in res], sorted(x["price"] for x in res))
                    self.assertTrue(all(x["touches"] >= 1 for x in sup + res))

    def test_cost_distribution(self):
        for seed, regime, daily, weekly in self.MARKETS:
            for key in ("daily", "weekly_52w"):
                with self.subTest(seed=seed, regime=regime, window=key):
                    cd = ta.analyse(daily, weekly)["cost_distribution"][key]
                    self.assertAlmostEqual(sum(b["share"] for b in cd["bins"]), 1.0, places=12)
                    self.assertTrue(0.0 <= cd["profit_ratio"] <= 1.0)
                    r70, r90 = cd["range_70"], cd["range_90"]
                    self.assertTrue(r90["low"] <= r70["low"] <= r70["high"] <= r90["high"])
                    lows, highs = [b["low"] for b in cd["bins"]], [b["high"] for b in cd["bins"]]
                    self.assertTrue(min(lows) <= cd["average_cost"] <= max(highs))
                    self.assertIn(cd["peak"], cd["bins"])
                    self.assertTrue(all(abs(b["high"] - b["low"] - cd["bin_size"]) < 1e-9 for b in cd["bins"]))

    def test_outlook(self):
        for seed, regime, daily, weekly in self.MARKETS:
            with self.subTest(seed=seed, regime=regime):
                f = full(daily, weekly)
                o, p = f["outlook"], f["price"]
                self.assertGreater(o["daily_sigma"], 0)
                prev_date = f["as_of"]
                for h in o["horizons"]:
                    self.assertTrue(h["q05"] <= h["q16"] <= h["q50"] <= h["q84"] <= h["q95"])
                    self.assertTrue(0.0 <= h["p_up"] <= 1.0 and 0.0 <= h["p_up_estimate"] <= 1.0)
                    self.assertEqual(h["date"], ta.add_weekdays(f["as_of"], h["days"]))
                    self.assertGreater(h["date"], prev_date)
                    prev_date = h["date"]
                self.assertEqual([x["price"] for x in o["levels"]], sorted((x["price"] for x in o["levels"]), reverse=True))
                for x in o["levels"]:
                    self.assertEqual(x["up"], x["price"] > p)
                    for t, c in zip(x["touch"], x["close_beyond"]):
                        self.assertLessEqual(c, t + 1e-12)  # a close beyond the level means it traded there
                    self.assertEqual(x["touch"], sorted(x["touch"]))  # more sessions, more chances
                    self.assertEqual(x["close_beyond"], sorted(x["close_beyond"]))

    def test_stance_bounds_and_determinism(self):
        for seed, regime, daily, weekly in self.MARKETS[:5]:
            with self.subTest(seed=seed, regime=regime):
                a = json.dumps(full(daily, weekly), sort_keys=True, default=str)
                b = json.dumps(full(daily, weekly), sort_keys=True, default=str)
                self.assertEqual(a, b)
                st = json.loads(a)["stance"]
                self.assertTrue(-1.0 <= st["score"] <= 1.0)
                self.assertIn(st["label"], {"BUY", "ACCUMULATE", "HOLD", "REDUCE", "SELL"})

    def test_close_to_the_closed_form(self):
        daily, weekly = market(41, "range")
        f = ta.analyse(daily, weekly)
        o = ta.outlook(f, daily, weekly, days=(5,), paths=20000)
        h, s = o["horizons"][0], o["daily_sigma"]
        self.assertAlmostEqual(h["q84"] / f["price"], math.exp(s * math.sqrt(5)), delta=0.004)
        self.assertAlmostEqual(h["q16"] / f["price"], math.exp(-s * math.sqrt(5)), delta=0.004)
        self.assertAlmostEqual(h["p_up"], 0.5, delta=0.015)


class EdgeCases(unittest.TestCase):
    def test_flat_prices_give_no_signal(self):
        flat = [{"date": d, "open": 50.0, "high": 50.0, "low": 50.0, "close": 50.0, "volume": 1000.0, "adjusted": False}
                for d in weekdays_ending("2026-10-30", 80)]
        f = ta.analyse(flat)
        self.assertEqual(f["momentum"]["rsi14"], 50.0)   # no gains and no losses is not "overbought"
        reads = {r["indicator"]: r["read"] for r in f["ledger"]["rows"]}
        self.assertEqual(reads["RSI14"], "neutral")
        self.assertEqual(reads["MACD histogram"], "neutral")
        self.assertEqual((f["stance"]["label"], f["stance"]["score"]), ("HOLD", 0.0))

    def test_short_histories_never_use_partial_windows(self):
        daily, _ = market(3, "range")
        for n in (2, 5, 15, 19):
            with self.subTest(bars=n):
                f = ta.analyse(daily[-n:])
                self.assertIsNone(f["volume"]["up_down_ratio_20"])
                self.assertNotIn("Up/down volume (20 bars)", [r["indicator"] for r in f["ledger"]["rows"]])
        f = ta.analyse(daily[-20:])
        self.assertIsNotNone(f["volume"]["up_down_ratio_20"])
        self.assertIsNone(ta.outlook(ta.analyse(daily[-5:]), daily[-5:]))  # too few returns for a sigma

    def test_daily_only_and_short_weekly(self):
        daily, weekly = market(5, "uptrend")
        f = ta.analyse(daily[-100:])
        self.assertIsNone(f["weekly"])
        self.assertIsNone(f["trend_template"])
        self.assertEqual(f["ma"]["sma200_source"], "unavailable")
        g = ta.analyse(daily[-100:], weekly[-120:])
        self.assertEqual(g["ma"]["sma200_source"], "approx-40-week")  # 100 daily bars: the 40-week stands in
        self.assertIsNone(ta.analyse(daily, weekly[-10:])["trend_template"])


class CommandLine(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.tmp = tempfile.TemporaryDirectory()
        cls.daily, cls.weekly = market(23, "range")  # a BUY, so the plan and sizing flags have work to do
        t = cls.tmp.name
        cls.paths = {"csv": os.path.join(t, "d.csv"), "wrapper": os.path.join(t, "d.json"),
                     "av": os.path.join(t, "d_av.json"), "weekly": os.path.join(t, "w.csv"),
                     "sma": os.path.join(t, "sma.csv"), "earnings": os.path.join(t, "e.csv")}
        with open(cls.paths["csv"], "w") as fh:
            fh.write(to_csv(cls.daily))
        with open(cls.paths["wrapper"], "w") as fh:
            json.dump({"result": to_csv(cls.daily)}, fh)
        with open(cls.paths["av"], "w") as fh:
            json.dump(to_av_json(cls.daily), fh)
        with open(cls.paths["weekly"], "w") as fh:
            fh.write(to_csv(cls.weekly, weekly=True))
        with open(cls.paths["sma"], "w") as fh:
            fh.write("time,SMA\n2026-10-30,100.5\n2026-10-29,100.4\n")
        with open(cls.paths["earnings"], "w") as fh:
            fh.write("symbol,name,reportDate,fiscalDateEnding,estimate,currency,timeOfTheDay\n"
                     "X,X CORP,2026-11-04,2026-09-30,1.00,USD,post-market\n")

    @classmethod
    def tearDownClass(cls):
        cls.tmp.cleanup()

    def json_of(self, *args):
        code, out, err = run_cli(list(args) + ["--json"])
        self.assertEqual(code, 0, err)
        return json.loads(out)

    def test_every_input_shape_gives_the_same_answer(self):
        base = ["--weekly", self.paths["weekly"], "--sma200-file", self.paths["sma"]]
        outs = [self.json_of("--daily", self.paths[k], *base) for k in ("csv", "wrapper", "av")]
        self.assertEqual(outs[0], outs[1])
        self.assertEqual(outs[0], outs[2])
        self.assertEqual(outs[0]["ma"]["sma200"], 100.5)  # the newest SMA row wins

    def test_flags(self):
        d = ["--daily", self.paths["csv"], "--weekly", self.paths["weekly"]]
        self.assertNotIn("outlook", self.json_of(*d, "--no-outlook"))
        sized = self.json_of(*d, "--account", "50000", "--risk-pct", "0.5", "--no-outlook")
        self.assertIsNotNone(sized["stance"]["plan"]["entry"])
        self.assertEqual(sized["position_size"]["risk_dollars"], 250.0)
        self.assertEqual(sized["position_size"]["binding"], "10% position cap")
        e = self.json_of(*d, "--earnings-file", self.paths["earnings"])["outlook"]["earnings"]
        self.assertEqual((e["date"], e["after_close"], e["day"]), ("2026-11-04", True, 4))  # 10-30 + 4 sessions
        one = self.json_of(*d, "--outlook-days", "1", "--level", "99.5:order")["outlook"]
        self.assertEqual([h["days"] for h in one["horizons"]], [1])
        self.assertIn("order", [x["label"] for x in one["levels"]])

    def test_risk_budget_can_bind(self):
        f = self.json_of("--daily", self.paths["csv"], "--weekly", self.paths["weekly"],
                         "--account", "100000", "--risk-pct", "0.05", "--no-outlook")
        self.assertIsNotNone(f["stance"]["plan"]["entry"])
        self.assertEqual(f["position_size"]["binding"], "risk budget")

    def test_errors_exit_cleanly(self):
        code, _, err = run_cli([])
        self.assertEqual(code, 2)
        self.assertIn("--daily is required", err)
        bad = os.path.join(self.tmp.name, "err.json")
        with open(bad, "w") as fh:
            json.dump({"error": "rate limit"}, fh)
        with self.assertRaises(ValueError):
            run_cli(["--daily", bad])
        for flag in (["--outlook-days", "0"], ["--level", "abc"]):
            code, _, err = run_cli(["--daily", self.paths["csv"], *flag])
            self.assertEqual(code, 2, flag)


class MarkdownReport(unittest.TestCase):
    def test_sections_and_numbers_match_the_json(self):
        daily, weekly = market(11, "uptrend")
        f = ta.analyse(daily, weekly, account=100_000)
        f["outlook"] = ta.outlook(f, daily, weekly, paths=PATHS)
        text = ta.render(f, "X")
        for head in ("## Price", "## Trend (daily)", "## Momentum", "## Volatility", "## Volume", "## Levels",
                     "## Position cost distribution", "## Weekly", "## Trend template", "## Weekly reversal checks",
                     "## Signal ledger", "## Stance", "## Short-term outlook"):
            self.assertIn(head, text)
        st = f["stance"]
        self.assertIn(f"**{st['label']}**", text)
        self.assertIn(f"{st['score']:+.2f}", text)
        self.assertIn(f"{f['price']:,.2f}", text)
        self.assertIn(f"{f['momentum']['rsi14']:.1f}", text)
        h = f["outlook"]["horizons"][0]
        self.assertIn(f"{h['q16']:,.2f}-{h['q84']:,.2f}", text)
        self.assertIsNotNone(f["position_size"])  # seed 11's uptrend has a plan
        self.assertIn("## Position size", text)
        self.assertIn(str(f["position_size"]["shares"]), text)


def snapshot(f):
    """The outputs pinned by the golden files."""
    st, o, cd = f["stance"], f["outlook"], f["cost_distribution"]["daily"]
    return {
        "as_of": f["as_of"], "price": f["price"],
        "indicators": {"sma20": f["ma"]["sma20"], "sma50": f["ma"]["sma50"], "rsi14": f["momentum"]["rsi14"],
                       "macd_hist": f["momentum"]["macd_hist"], "adx14": f["momentum"]["adx14"],
                       "atr14": f["volatility"]["atr14"], "bb_pct_b": f["volatility"]["bb_pct_b"]},
        "stance": {k: st[k] for k in ("label", "score", "net", "total", "support_1", "resistance_1")},
        "plan": st["plan"], "levels": f["levels"], "weekly_levels": f["weekly"]["levels"],
        "trend_template": f["trend_template"]["passed"], "reversals": f["weekly_reversals"]["read"],
        "cost": {"peak": cd["peak"], "average_cost": cd["average_cost"], "profit_ratio": cd["profit_ratio"],
                 "range_70": cd["range_70"], "range_90": cd["range_90"]},
        "outlook": {"daily_sigma": o["daily_sigma"], "horizons": o["horizons"], "levels": o["levels"]},
    }


def diff(want, got, path=""):
    if isinstance(want, dict) and isinstance(got, dict):
        out = [f"{path}.{k}: missing" for k in want if k not in got] + [f"{path}.{k}: new" for k in got if k not in want]
        for k in want.keys() & got.keys():
            out += diff(want[k], got[k], f"{path}.{k}")
        return out
    if isinstance(want, list) and isinstance(got, list):
        out = [] if len(want) == len(got) else [f"{path}: {len(want)} items vs {len(got)}"]
        for i, (a, b) in enumerate(zip(want, got)):
            out += diff(a, b, f"{path}[{i}]")
        return out
    if isinstance(want, float) and isinstance(got, (int, float)):
        return [] if math.isclose(want, got, rel_tol=1e-12, abs_tol=1e-12) else [f"{path}: {want} vs {got}"]
    return [] if want == got else [f"{path}: {want!r} vs {got!r}"]


class Golden(unittest.TestCase):
    def current(self, name):
        seed, regime = GOLDEN[name]
        daily, weekly = market(seed, regime)
        return json.loads(json.dumps(snapshot(full(daily, weekly)), default=str))

    def test_outputs_match_the_pinned_files(self):
        update = os.environ.get("UPDATE_GOLDEN") == "1"
        for name in GOLDEN:
            with self.subTest(fixture=name):
                path = os.path.join(GOLDEN_DIR, f"{name}.json")
                now = self.current(name)
                if update:
                    os.makedirs(GOLDEN_DIR, exist_ok=True)
                    with open(path, "w") as fh:
                        json.dump(now, fh, indent=1, sort_keys=True)
                        fh.write("\n")
                    continue
                with open(path) as fh:
                    want = json.load(fh)
                d = diff(want, now)
                self.assertEqual(d, [], f"{name} moved: {d[:5]}")

    def test_a_one_line_band_change_is_caught(self):
        saved = ta.BANDS
        try:
            ta.BANDS = ((-1.01, "SELL"),)  # one line: every score is a SELL
            moved = []
            for name in GOLDEN:
                with open(os.path.join(GOLDEN_DIR, f"{name}.json")) as fh:
                    moved += diff(json.load(fh), self.current(name))
            self.assertTrue(moved)
        finally:
            ta.BANDS = saved


if __name__ == "__main__":
    unittest.main()
