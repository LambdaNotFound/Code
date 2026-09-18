"""Tests for chart.py. Run from this directory: python3 -m unittest test_chart -v"""

import os
import sys
import unittest
import xml.etree.ElementTree as ET

sys.path.insert(0, os.path.dirname(__file__))
import chart  # noqa: E402
import ta  # noqa: E402


def bars(closes, start="2026-01-01"):
    """Daily bars with a real range, dated forward from start."""
    import datetime
    d0 = datetime.date.fromisoformat(start)
    out = []
    for i, c in enumerate(closes):
        out.append({"date": (d0 + datetime.timedelta(days=i)).isoformat(),
                    "open": c - 0.4, "high": c + 1.2, "low": c - 1.2, "close": c,
                    "volume": 1000.0 + 10 * i, "adjusted": True})
    return out


def weekly(n=160, start=100.0, step=1.0):
    import datetime
    d0 = datetime.date(2022, 1, 7)
    return [{"date": (d0 + datetime.timedelta(weeks=i)).isoformat(),
             "open": start + i * step - 0.5, "high": start + i * step + 2,
             "low": start + i * step - 2, "close": start + i * step,
             "volume": 5000.0, "adjusted": True} for i in range(n)]


class HelperTests(unittest.TestCase):
    def test_esc_escapes_markup(self):
        self.assertEqual(chart.esc("a<b>&c"), "a&lt;b&gt;&amp;c")

    def test_esc_leaves_middot_alone(self):
        # the panel titles use the character, not the entity, so escaping is a no-op
        self.assertEqual(chart.esc("Daily · last 100 bars"), "Daily · last 100 bars")

    def test_nice_ticks_span_the_range(self):
        t = chart.nice_ticks(103.0, 197.0)
        self.assertTrue(t[0] >= 100 and t[-1] <= 200)
        self.assertGreaterEqual(len(t), 3)
        step = t[1] - t[0]
        self.assertTrue(all(abs((b - a) - step) < 1e-9 for a, b in zip(t, t[1:])))

    def test_nice_ticks_degenerate_range(self):
        self.assertEqual(chart.nice_ticks(50.0, 50.0), [50.0])


class PanelTests(unittest.TestCase):
    def panel(self):
        return chart.Panel(bars([100, 110, 90, 105]), top=10, height=200,
                           theme=chart.THEMES["light"], title="t")

    def test_price_maps_high_to_low_y(self):
        p = self.panel()
        self.assertLess(p.y(p.hi), p.y(p.lo))
        self.assertAlmostEqual(p.y(p.hi), p.top)
        self.assertAlmostEqual(p.y(p.lo), p.top + p.h)

    def test_bars_span_the_width_in_order(self):
        p = self.panel()
        xs = [p.x(i) for i in range(p.n)]
        self.assertEqual(xs, sorted(xs))
        self.assertGreater(xs[0], p.x0)
        self.assertLess(xs[-1], p.x1)

    def test_range_is_padded_beyond_the_extremes(self):
        p = self.panel()
        self.assertLess(p.lo, 88.8)   # low 90 minus the wick and the pad
        self.assertGreater(p.hi, 111.2)

    def test_candle_colour_follows_direction(self):
        up = bars([100])
        up[0]["open"], up[0]["close"] = 99.0, 101.0
        out = []
        chart.Panel(up, 0, 100, chart.THEMES["light"], "t").candles(out)
        self.assertIn(chart.THEMES["light"]["up"], " ".join(out))
        down = bars([100])
        down[0]["open"], down[0]["close"] = 101.0, 99.0
        out = []
        chart.Panel(down, 0, 100, chart.THEMES["light"], "t").candles(out)
        self.assertIn(chart.THEMES["light"]["down"], " ".join(out))

    def test_level_outside_the_pane_is_skipped(self):
        p = self.panel()
        out = []
        p.level(out, 10_000.0, "#000", "R")
        self.assertEqual(out, [])
        p.level(out, 100.0, "#000", "R")
        self.assertTrue(out)

    def test_line_needs_two_points(self):
        p = self.panel()
        out = []
        p.line(out, [None, None, None, 5.0], "#000")
        self.assertEqual(out, [])


class BuildTests(unittest.TestCase):
    def setUp(self):
        self.daily = bars([100 + (i % 7) - 3 + i * 0.5 for i in range(100)])
        self.weekly = weekly()
        self.facts = ta.analyse(self.daily, self.weekly, sma200_daily=95.0)

    def parse(self, svg):
        root = ET.fromstring(svg)
        self.assertTrue(root.tag.endswith("svg"))
        return root

    def test_build_is_valid_svg_with_all_panels(self):
        svg = chart.build("TEST", self.daily, self.weekly, self.facts)
        root = self.parse(svg)
        self.assertIn("Daily", svg)
        self.assertIn("Weekly", svg)
        self.assertIn("Monthly", svg)
        self.assertIn("Volume", svg)
        self.assertIn("TEST", svg)
        rects = [e for e in root.iter() if e.tag.endswith("rect")]
        self.assertGreater(len(rects), 100)   # candle bodies plus frames plus volume

    def test_declared_height_covers_every_drawn_element(self):
        svg = chart.build("TEST", self.daily, self.weekly, self.facts)
        root = self.parse(svg)
        height = float(root.get("height"))
        lowest = 0.0
        for e in root.iter():
            for attr in ("y", "y1", "y2"):
                if e.get(attr):
                    lowest = max(lowest, float(e.get(attr)))
        self.assertLessEqual(lowest, height, "an element is drawn below the canvas")
        self.assertGreater(height, 600)

    def test_daily_only_still_builds(self):
        facts = ta.analyse(self.daily, None, sma200_daily=95.0)
        svg = chart.build("TEST", self.daily, None, facts)
        self.parse(svg)
        self.assertIn("Daily", svg)
        self.assertNotIn("Weekly ·", svg)
        self.assertNotIn("Monthly ·", svg)

    def test_both_themes_render_and_differ(self):
        light = chart.build("T", self.daily, self.weekly, self.facts, "light")
        dark = chart.build("T", self.daily, self.weekly, self.facts, "dark")
        self.parse(light)
        self.parse(dark)
        self.assertIn(chart.THEMES["light"]["bg"], light)
        self.assertIn(chart.THEMES["dark"]["bg"], dark)
        self.assertNotEqual(light, dark)

    def test_every_theme_colour_is_a_hex_triplet(self):
        import re
        for name, t in chart.THEMES.items():
            for key, v in t.items():
                self.assertRegex(v, r"^#[0-9a-fA-F]{6}$", f"{name}.{key} = {v}")

    def test_levels_from_facts_are_labelled(self):
        svg = chart.build("TEST", self.daily, self.weekly, self.facts)
        sup = self.facts["levels"]["supports"]
        drawn = [c for c in sup[:3] if chart.Panel(self.daily[-100:], 0, 300,
                 chart.THEMES["light"], "t").lo < c["price"] <
                 chart.Panel(self.daily[-100:], 0, 300, chart.THEMES["light"], "t").hi]
        for c in drawn:
            self.assertIn(f'S {c["price"]:,.2f}', svg)

    def test_daily_bars_argument_limits_the_pane(self):
        svg = chart.build("TEST", self.daily, self.weekly, self.facts, daily_bars=30)
        self.assertIn("Daily · last 30 bars", svg)


if __name__ == "__main__":
    unittest.main()
