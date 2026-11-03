"""Planted errors: do the reference tests catch a wrong ta.py?

Each mutation below is a one-line change to a copy of ta.py, of the kind a
careless edit makes: a different period or seed, a sample instead of a
population deviation, a true range that ignores gaps, a dropped guard. The
copy is loaded as its own module and the whole test_ta_reference suite runs
against it. Every mutation must fail at least one check; the unchanged copy
must fail none.

Run from this directory: python3 -m unittest test_ta_mutations -v
"""
import importlib.util
import os
import sys
import tempfile
import unittest

HERE = os.path.dirname(os.path.abspath(__file__))
sys.path.insert(0, HERE)
import test_ta_reference as ref  # noqa: E402

MUTATIONS = [  # (what, text in ta.py, replacement): each text occurs once
    ("RSI period 13", "def rsi(closes, n=14):", "def rsi(closes, n=13):"),
    ("EMA alpha 2/(n+2)", "    k = 2.0 / (n + 1)\n    prev = sum(vals[:n]) / n", "    k = 2.0 / (n + 2)\n    prev = sum(vals[:n]) / n"),
    ("Bollinger sample deviation", "sd = math.sqrt(sum((x - m) ** 2 for x in window) / n)",
     "sd = math.sqrt(sum((x - m) ** 2 for x in window) / (n - 1))"),
    ("true range ignores gaps", "        tr.append(max(h - l, abs(h - pc), abs(l - pc)))\n    return tr",
     "        tr.append(h - l)\n    return tr"),
    ("fast stochastic", "def stochastic(bars, k_n=14, k_smooth=3, d_n=3):", "def stochastic(bars, k_n=14, k_smooth=1, d_n=3):"),
    ("MACD signal 8", "def macd(closes, fast=12, slow=26, signal=9):", "def macd(closes, fast=12, slow=26, signal=8):"),
    ("ATR period 10", "def atr(bars, n=14):", "def atr(bars, n=10):"),
    ("profit ratio halved", "profit += v * min(1.0, max(0.0, (price - lo) / (hi - lo)))",
     "profit += v * min(1.0, max(0.0, (price - lo) / (hi - lo) / 2))"),
    ("up/down volume drops its guard",
     'up_vol = sum(b["volume"] for i, b in enumerate(win) if i and b["close"] > win[i - 1]["close"])',
     'up_vol = sum(b["volume"] for i, b in enumerate(win) if b["close"] > win[i - 1]["close"])'),
    ("50-week reversal window", "def weekly_reversals(weekly, extreme_weeks=52,", "def weekly_reversals(weekly, extreme_weeks=50,"),
    ("ADX threshold 20", "strong = mo[\"adx14\"] >= 25", "strong = mo[\"adx14\"] >= 20"),
    ("stop a full ATR past support", "stop = (s1 - 0.5 * a) if s1 is not None else price - 2 * a",
     "stop = (s1 - 1.0 * a) if s1 is not None else price - 2 * a"),
    ("weekly SMA slope over 2 weeks", '"sma20_slope": slope_sign(s20, 4)', '"sma20_slope": slope_sign(s20, 2)'),
    ("pivots need only 2 bars a side", "pivs = pivots(daily, pivot_window)", "pivs = pivots(daily, 2)"),
]


def load_variant(src, name, tmp):
    path = os.path.join(tmp, f"{name}.py")
    with open(path, "w", encoding="utf-8") as fh:
        fh.write(src)
    spec = importlib.util.spec_from_file_location(name, path)
    mod = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(mod)
    return mod


def failures_with(module):
    saved = ref.ta
    ref.ta = module
    try:
        result = unittest.TestResult()
        unittest.TestLoader().loadTestsFromModule(ref).run(result)
        return len(result.failures) + len(result.errors)
    finally:
        ref.ta = saved


class Mutations(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        with open(os.path.join(HERE, "ta.py"), encoding="utf-8") as fh:
            cls.src = fh.read()
        cls.tmp = tempfile.TemporaryDirectory()

    @classmethod
    def tearDownClass(cls):
        cls.tmp.cleanup()

    def test_the_unchanged_copy_passes(self):
        self.assertEqual(failures_with(load_variant(self.src, "ta_control", self.tmp.name)), 0)

    def test_every_mutation_is_caught(self):
        for i, (what, old, new) in enumerate(MUTATIONS):
            with self.subTest(mutation=what):
                self.assertEqual(self.src.count(old), 1, f"{what}: the text to replace must occur once in ta.py")
                mod = load_variant(self.src.replace(old, new), f"ta_mut_{i}", self.tmp.name)
                self.assertGreater(failures_with(mod), 0, f"{what}: no reference check failed")


if __name__ == "__main__":
    unittest.main()
