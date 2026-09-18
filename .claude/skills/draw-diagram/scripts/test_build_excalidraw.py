#!/usr/bin/env python3
"""Tests for build_excalidraw.py. Stdlib unittest, no fixtures on disk.

    python3 .claude/skills/draw-diagram/scripts/test_build_excalidraw.py
"""
import copy
import json
import os
import re
import sys
import tempfile
import unittest
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent))
import build_excalidraw as bx  # noqa: E402


def build(spec):
    b = bx.Builder(spec)
    b.build()
    return b


def warns(b):
    return [m for lv, m in b.problems if lv == "warn"]


def two_boxes(**edge):
    return {"elements": [{"id": "a", "kind": "box", "text": "A", "x": 0, "y": 0},
                         {"id": "b", "kind": "box", "text": "B", "x": 300, "y": 0}],
            "edges": [dict({"from": "a", "to": "b"}, **edge)]}


class Example(unittest.TestCase):
    def test_example_builds_without_warnings(self):
        b = build(copy.deepcopy(bx.EXAMPLE))
        self.assertEqual(warns(b), [])
        self.assertGreater(len(b.elements), 10)

    def test_same_spec_same_bytes(self):
        a = json.dumps(build(copy.deepcopy(bx.EXAMPLE)).excalidraw(), sort_keys=True)
        c = json.dumps(build(copy.deepcopy(bx.EXAMPLE)).excalidraw(), sort_keys=True)
        self.assertEqual(a, c)

    def test_every_element_carries_the_fields_excalidraw_restores(self):
        for el in build(copy.deepcopy(bx.EXAMPLE)).excalidraw()["elements"]:
            for k in ("id", "type", "x", "y", "width", "height", "strokeColor", "backgroundColor",
                      "seed", "versionNonce", "isDeleted", "groupIds", "boundElements", "roughness", "opacity"):
                self.assertIn(k, el, f"{el['id']} lacks {k}")
            if el["type"] == "text":
                for k in ("text", "originalText", "fontSize", "fontFamily", "textAlign", "verticalAlign",
                          "containerId", "lineHeight", "autoResize"):
                    self.assertIn(k, el, f"{el['id']} lacks {k}")
            if el["type"] in ("arrow", "line"):
                for k in ("points", "startBinding", "endBinding", "startArrowhead", "endArrowhead"):
                    self.assertIn(k, el, f"{el['id']} lacks {k}")


class Bindings(unittest.TestCase):
    def test_bound_text_links_both_ways(self):
        els = build(copy.deepcopy(bx.EXAMPLE)).elements
        text = els["validate__text"]
        self.assertEqual(text["containerId"], "validate")
        self.assertIn({"id": "validate__text", "type": "text"}, els["validate"]["boundElements"])

    def test_arrow_bindings_link_both_ways(self):
        els = build(copy.deepcopy(bx.EXAMPLE)).elements
        for el in els.values():
            if el["type"] != "arrow":
                continue
            for side in ("startBinding", "endBinding"):
                target = el[side]["elementId"]
                self.assertIn(target, els)
                self.assertIn({"id": el["id"], "type": "arrow"}, els[target]["boundElements"])

    def test_arrow_label_is_bound_to_the_arrow(self):
        els = build(two_boxes(label="yes")).elements
        self.assertEqual(els["a->b__label"]["containerId"], "a->b")
        self.assertIn({"id": "a->b__label", "type": "text"}, els["a->b"]["boundElements"])

    def test_arrow_starts_and_ends_at_the_shape_edges(self):
        b = build(two_boxes())
        a, t, arrow = b.nodes["a"], b.nodes["b"], b.elements["a->b"]
        x0 = arrow["x"]
        x1 = arrow["x"] + arrow["points"][-1][0]
        self.assertAlmostEqual(x0, a["x"] + a["w"] + bx.GAP, delta=0.6)
        self.assertAlmostEqual(x1, t["x"] - bx.GAP, delta=0.6)

    def test_exit_point_lands_on_each_outline(self):
        for shape in ("rect", "ellipse", "diamond"):
            n = {"x": 0, "y": 0, "w": 200, "h": 100, "shape": shape}
            px, py = bx.exit_point(n, (1000, 50))
            self.assertAlmostEqual(px, 200 + bx.GAP, delta=0.01, msg=shape)
            self.assertAlmostEqual(py, 50, delta=0.01, msg=shape)
            px, py = bx.exit_point(n, (100, -1000))
            self.assertAlmostEqual(py, -bx.GAP, delta=0.01, msg=shape)

    def test_duplicate_edges_get_distinct_ids(self):
        spec = two_boxes()
        spec["edges"].append({"from": "a", "to": "b", "dashed": True})
        els = build(spec).elements
        self.assertIn("a->b", els)
        self.assertIn("a->b#2", els)


class Errors(unittest.TestCase):
    def test_unknown_edge_endpoint(self):
        spec = two_boxes()
        spec["edges"][0]["to"] = "nope"
        with self.assertRaisesRegex(bx.SpecError, "unknown id 'nope'"):
            build(spec)

    def test_unknown_kind_role_level(self):
        with self.assertRaisesRegex(bx.SpecError, "unknown kind"):
            build({"elements": [{"id": "a", "kind": "blob", "x": 0, "y": 0}]})
        with self.assertRaisesRegex(bx.SpecError, "unknown role"):
            build({"elements": [{"id": "a", "kind": "box", "role": "magenta", "text": "x", "x": 0, "y": 0}]})
        with self.assertRaisesRegex(bx.SpecError, "unknown text level"):
            build({"elements": [{"id": "a", "kind": "text", "level": "huge", "text": "x", "x": 0, "y": 0}]})

    def test_duplicate_id_and_missing_position(self):
        with self.assertRaisesRegex(bx.SpecError, "duplicate id"):
            build({"elements": [{"id": "a", "kind": "box", "text": "x", "x": 0, "y": 0},
                                {"id": "a", "kind": "box", "text": "y", "x": 9, "y": 9}]})
        with self.assertRaisesRegex(bx.SpecError, "needs x,y or col,row"):
            build({"elements": [{"id": "a", "kind": "box", "text": "x"}]})

    def test_section_cycle(self):
        spec = {"elements": [{"id": "a", "kind": "box", "text": "x", "x": 0, "y": 0},
                             {"id": "s1", "kind": "section", "contains": ["a", "s2"]},
                             {"id": "s2", "kind": "section", "contains": ["s1"]}]}
        with self.assertRaisesRegex(bx.SpecError, "contains itself"):
            build(spec)

    def test_edge_cannot_end_on_a_line(self):
        spec = two_boxes()
        spec["elements"].append({"id": "l", "kind": "line", "points": [[0, 200], [300, 200]]})
        spec["edges"].append({"from": "a", "to": "l"})
        with self.assertRaisesRegex(bx.SpecError, "cannot be an endpoint"):
            build(spec)


class Warnings(unittest.TestCase):
    def test_explicit_size_too_small_warns(self):
        spec = {"elements": [{"id": "a", "kind": "box", "text": "a fairly long label here", "x": 0, "y": 0, "w": 80}]}
        self.assertTrue(any("text overflow" in w for w in warns(build(spec))))

    def test_overlap_warns(self):
        spec = {"elements": [{"id": "a", "kind": "box", "text": "A", "x": 0, "y": 0},
                             {"id": "b", "kind": "box", "text": "B", "x": 20, "y": 10}]}
        self.assertTrue(any("overlap" in w for w in warns(build(spec))))

    def test_arrow_through_a_third_shape_warns(self):
        spec = {"elements": [{"id": "a", "kind": "box", "text": "A", "x": 0, "y": 0},
                             {"id": "m", "kind": "box", "text": "M", "x": 300, "y": 0},
                             {"id": "c", "kind": "box", "text": "C", "x": 600, "y": 0}],
                "edges": [{"from": "a", "to": "c"}]}
        self.assertTrue(any("crosses 'm'" in w for w in warns(build(spec))))
        spec["edges"][0]["via"] = [[100, 150], [600, 150]]
        self.assertFalse(any("crosses" in w for w in warns(build(spec))))

    def test_intruder_in_section_warns_and_member_does_not(self):
        spec = {"elements": [{"id": "a", "kind": "box", "text": "A", "x": 0, "y": 0},
                             {"id": "b", "kind": "box", "text": "B", "x": 0, "y": 80},
                             {"id": "s", "kind": "section", "label": "S", "contains": ["a"]}]}
        ws = warns(build(spec))
        self.assertTrue(any("'b' sits inside section 's'" in w for w in ws))
        spec["elements"][2]["contains"] = ["a", "b"]
        self.assertEqual(warns(build(spec)), [])

    def test_lonely_shape_is_a_note_not_a_warning(self):
        b = build({"elements": [{"id": "a", "kind": "box", "text": "A", "x": 0, "y": 0}]})
        self.assertEqual(warns(b), [])
        self.assertTrue(any("no edges" in m for lv, m in b.problems if lv == "note"))


class Layout(unittest.TestCase):
    def test_grid_centres_the_shape_in_its_cell(self):
        b = build({"grid": {"colWidth": 260, "rowHeight": 150},
                   "elements": [{"id": "a", "kind": "box", "text": "hello", "col": 1, "row": 0}]})
        self.assertEqual(bx.center(b.nodes["a"]), (390, 75))

    def test_section_encloses_members_with_padding(self):
        b = build({"elements": [{"id": "a", "kind": "box", "text": "A", "x": 100, "y": 100},
                                {"id": "s", "kind": "section", "label": "S", "contains": ["a"]}]})
        a, s = b.nodes["a"], b.nodes["s"]
        self.assertLess(s["x"], a["x"])
        self.assertLess(s["y"], a["y"])
        self.assertGreater(s["x"] + s["w"], a["x"] + a["w"])
        self.assertGreater(s["y"] + s["h"], a["y"] + a["h"])
        order = b.order
        self.assertLess(order.index("s"), order.index("a"), "sections draw behind members")

    def test_timeline_expands_and_its_dots_bind(self):
        spec = {"elements": [{"id": "t", "kind": "timeline", "x": 0, "y": 0, "steps": ["one", "two", "three"]},
                             {"id": "a", "kind": "box", "text": "A", "x": -300, "y": 50}],
                "edges": [{"from": "a", "to": "t.dot2"}]}
        b = build(spec)
        for i in (1, 2, 3):
            self.assertIn(f"t.dot{i}", b.elements)
            self.assertIn(f"t.label{i}", b.elements)
        self.assertEqual(b.elements["a->t.dot2"]["endBinding"]["elementId"], "t.dot2")
        self.assertEqual(warns(b), [])

    def test_shape_size_matches_excalidraw_text_room(self):
        tw, th = 100, 20
        for kind in ("box", "ellipse", "diamond"):
            w, h = bx.shape_size(kind, tw, th)
            rw, rh = bx.text_room(kind, w, h)
            self.assertGreaterEqual(rw + 0.5, tw, kind)
            self.assertGreaterEqual(rh + 0.5, th, kind)

    def test_sketch_style_switches_font_and_roughness(self):
        spec = two_boxes()
        spec["style"] = "sketch"
        els = build(spec).elements
        self.assertEqual(els["a"]["roughness"], 1)
        self.assertEqual(els["a__text"]["fontFamily"], bx.FONT_SKETCH)

    def test_inactive_role_is_dashed(self):
        els = build({"elements": [{"id": "a", "kind": "box", "role": "inactive", "text": "later", "x": 0, "y": 0}]}).elements
        self.assertEqual(els["a"]["strokeStyle"], "dashed")


class Outputs(unittest.TestCase):
    def test_svg_draws_every_shape_and_line(self):
        b = build(copy.deepcopy(bx.EXAMPLE))
        svg, w, h = b.svg()
        for el in b.elements.values():
            if el["type"] != "text":
                self.assertIn(f'id="{bx.esc(el["id"])}"', svg)
        self.assertGreater(w, 100)
        self.assertGreater(h, 100)

    def test_cli_writes_excalidraw_and_svg(self):
        with tempfile.TemporaryDirectory() as d:
            spec = Path(d) / "flow.spec.json"
            spec.write_text(json.dumps(bx.EXAMPLE), encoding="utf-8")
            rc = bx.main([str(spec), "--no-png"])
            self.assertEqual(rc, 0)
            out = json.loads((Path(d) / "flow.excalidraw").read_text(encoding="utf-8"))
            self.assertEqual(out["type"], "excalidraw")
            self.assertEqual(out["version"], 2)
            self.assertTrue((Path(d) / "flow.svg").exists())
            self.assertFalse((Path(d) / "flow.png").exists())

    def test_cli_rejects_a_broken_spec_and_writes_nothing(self):
        with tempfile.TemporaryDirectory() as d:
            spec = Path(d) / "bad.json"
            spec.write_text(json.dumps({"elements": [{"id": "a", "kind": "nope"}]}), encoding="utf-8")
            self.assertEqual(bx.main([str(spec), "--no-png"]), 1)
            self.assertEqual(os.listdir(d), ["bad.json"])

    @unittest.skipUnless(bx.find_chromium(), "no Chromium on this machine")
    def test_png_when_chromium_is_present_covers_the_whole_canvas(self):
        with tempfile.TemporaryDirectory() as d:
            spec = Path(d) / "flow.spec.json"
            spec.write_text(json.dumps(bx.EXAMPLE), encoding="utf-8")
            self.assertEqual(bx.main([str(spec)]), 0)
            png = Path(d) / "flow.png"
            self.assertEqual(png.read_bytes()[:8], b"\x89PNG\r\n\x1a\n")
            _, w, h = build(copy.deepcopy(bx.EXAMPLE)).svg()
            pw, ph = bx.png_size(png)
            self.assertEqual(pw, w * 2)
            self.assertGreaterEqual(ph, h * 2, "PNG shorter than the canvas: bottom of the diagram clipped")


class Patterns(unittest.TestCase):
    def test_every_pattern_in_the_reference_builds_without_warnings(self):
        doc = (Path(__file__).resolve().parent.parent / "references" / "visual-patterns.md").read_text(encoding="utf-8")
        blocks = re.findall(r"```json\n(.*?)```", doc, re.S)
        self.assertGreaterEqual(len(blocks), 8)
        for i, block in enumerate(blocks):
            with self.subTest(block=i):
                self.assertEqual(warns(build(json.loads(block))), [])


if __name__ == "__main__":
    unittest.main(verbosity=1)
