#!/usr/bin/env python3
"""Build an Excalidraw file from a diagram spec, then check and preview it.

The spec (format: ../references/spec-format.md) is the source of truth.
From it this script writes, next to the spec unless --out is given:

    <name>.excalidraw   the deliverable; opens in excalidraw.com and the VS Code extension
    <name>.svg          a geometry-faithful preview drawn by this script
    <name>.png          the SVG rasterised by a local Chromium, when one is found

and prints structural checks: errors (exit 1, nothing written) for a spec
that cannot build, warnings (exit 0) for layout defects a reader would see.

Usage:
    python3 build_excalidraw.py SPEC.json [--out BASE] [--no-png] [--scale N]
    python3 build_excalidraw.py --roles      # palette roles and text levels
    python3 build_excalidraw.py --example    # a small valid spec on stdout

Stdlib only. Same spec, same bytes: every seed derives from the element id,
so a rebuild produces a clean git diff.

Preview fidelity: the SVG places every element at the coordinates the
.excalidraw carries and draws text in a monospace fallback, so it shows
the layout Excalidraw will show. It is not Excalidraw's renderer; the
file itself is the ground truth.
"""
from __future__ import annotations

import argparse
import glob
import json
import math
import os
import shutil
import struct
import subprocess
import sys
import zlib
from pathlib import Path

HERE = Path(__file__).resolve().parent
PALETTE = json.loads((HERE / "palette.json").read_text(encoding="utf-8"))

# Excalidraw constants, read from the @excalidraw/excalidraw 0.18.1 dist.
FONT_MONO = 3                      # Cascadia
FONT_SKETCH = 5                    # Excalifont
ROUNDNESS_ADAPTIVE = {"type": 3}   # rectangles with rounded corners

# Layout constants of this script.
CHAR_W = 0.6        # em per monospace glyph; Cascadia advances 0.586em, so nothing sized here clips
LINE_H = 1.25       # written to every text element, so Excalidraw spaces lines the same way
PAD_X, PAD_Y = 24, 16
MIN_W, MIN_H = 60, 40
GAP = 4             # an arrow tip stops this far from the shape edge
CODE_PAD = 16
DOT = 12
TIMELINE_GAP = 70
SECTION_PAD = 32
SQRT2 = math.sqrt(2)

SHAPE_KINDS = ("box", "ellipse", "diamond")
KINDS = SHAPE_KINDS + ("text", "code", "dot", "line", "arrow", "timeline", "section")
ARROWHEADS = (None, "arrow", "bar", "dot", "triangle")

EXAMPLE = {
    "title": "Webhook intake",
    "grid": {"colWidth": 260, "rowHeight": 150},
    "elements": [
        {"id": "hook", "kind": "ellipse", "role": "start", "text": "Webhook\nreceived", "col": 0, "row": 0},
        {"id": "validate", "kind": "box", "text": "Validate signature", "col": 1, "row": 0},
        {"id": "ok", "kind": "diamond", "role": "decision", "text": "valid?", "col": 2, "row": 0},
        {"id": "queue", "kind": "box", "role": "data", "text": "jobs queue", "col": 3, "row": 0},
        {"id": "reject", "kind": "box", "role": "error", "text": "401 rejected", "col": 2, "row": 1},
        {"id": "payload", "kind": "code", "lang": "json",
         "text": "{\n  \"event\": \"push\",\n  \"sig\": \"sha256=...\"\n}", "col": 0, "row": 1},
        {"id": "intake", "kind": "section", "label": "Intake service", "contains": ["validate", "ok", "reject"]},
        {"id": "steps", "kind": "timeline", "direction": "down", "x": 1120, "y": 40,
         "steps": ["received", "verified", "enqueued"]},
    ],
    "edges": [
        {"from": "hook", "to": "validate"},
        {"from": "validate", "to": "ok"},
        {"from": "ok", "to": "queue", "label": "yes"},
        {"from": "ok", "to": "reject", "label": "no", "dashed": True},
        {"from": "payload", "to": "validate", "label": "body", "dashed": True},
        {"from": "queue", "to": "steps.dot3", "dashed": True, "end": None},
    ],
}


class SpecError(Exception):
    """The spec cannot build. The message names the element."""


# ---------------------------------------------------------------- helpers

def seed(key: str) -> int:
    return zlib.crc32(key.encode("utf-8")) % 2_000_000_000 + 1


def rnd(v: float) -> float:
    v = round(v, 2)
    return int(v) if v == int(v) else v


def measure(text: str, size: float) -> tuple[int, int, list[str]]:
    lines = text.split("\n")
    w = math.ceil(max(len(line) for line in lines) * size * CHAR_W)
    h = math.ceil(len(lines) * size * LINE_H)
    return w, h, lines


def shape_size(kind: str, tw: int, th: int) -> tuple[int, int]:
    """Container size that holds a tw x th text block the way Excalidraw fits bound text:
    the full width of a box, width/sqrt2 of an ellipse, width/2 of a diamond."""
    if kind == "box":
        return max(MIN_W, tw + 2 * PAD_X), max(MIN_H, th + 2 * PAD_Y)
    if kind == "ellipse":
        return math.ceil((tw + 2 * PAD_X) * SQRT2), math.ceil((th + 2 * PAD_Y) * SQRT2)
    return 2 * tw + 2 * PAD_X, 2 * th + 2 * PAD_Y


def text_room(kind: str, w: float, h: float) -> tuple[float, float]:
    """Largest text block a w x h container of this kind shows without clipping."""
    if kind == "ellipse":
        return w / SQRT2 - 2 * PAD_X, h / SQRT2 - 2 * PAD_Y
    if kind == "diamond":
        return w / 2 - PAD_X, h / 2 - PAD_Y
    return w - 2 * PAD_X, h - 2 * PAD_Y


def base(id_: str, type_: str, x: float, y: float, w: float, h: float, *, stroke: str,
         fill: str = "transparent", stroke_width: int = 2, dashed: bool = False,
         roughness: int = 0, roundness: dict | None = None, group: str | None = None) -> dict:
    return {
        "id": id_, "type": type_, "x": rnd(x), "y": rnd(y), "width": rnd(w), "height": rnd(h),
        "angle": 0, "strokeColor": stroke, "backgroundColor": fill, "fillStyle": "solid",
        "strokeWidth": stroke_width, "strokeStyle": "dashed" if dashed else "solid",
        "roughness": roughness, "opacity": 100, "groupIds": [group] if group else [],
        "frameId": None, "roundness": roundness, "seed": seed(id_), "version": 1,
        "versionNonce": seed(id_ + "/nonce"), "isDeleted": False, "boundElements": [],
        "link": None, "locked": False,
    }


def text_el(id_: str, x: float, y: float, w: float, h: float, text: str, size: float, color: str,
            font: int, *, align: str = "center", valign: str = "middle",
            container: str | None = None, roughness: int = 0, group: str | None = None) -> dict:
    e = base(id_, "text", x, y, w, h, stroke=color, stroke_width=1, roughness=roughness, group=group)
    e.update({"fontSize": size, "fontFamily": font, "text": text, "textAlign": align,
              "verticalAlign": valign, "containerId": container, "originalText": text,
              "autoResize": True, "lineHeight": LINE_H})
    return e


def linear_el(id_: str, type_: str, pts: list[tuple[float, float]], *, stroke: str, width: int = 2,
              dashed: bool = False, start: str | None = None, end: str | None = None,
              start_bind: str | None = None, end_bind: str | None = None, roughness: int = 0,
              group: str | None = None) -> dict:
    x0, y0 = pts[0]
    rel = [[rnd(px - x0), rnd(py - y0)] for px, py in pts]
    xs = [p[0] for p in rel]
    ys = [p[1] for p in rel]
    e = base(id_, type_, x0, y0, max(xs) - min(xs), max(ys) - min(ys), stroke=stroke,
             stroke_width=width, dashed=dashed, roughness=roughness, group=group)
    e.update({
        "points": rel, "lastCommittedPoint": None,
        "startBinding": {"elementId": start_bind, "focus": 0, "gap": GAP} if start_bind else None,
        "endBinding": {"elementId": end_bind, "focus": 0, "gap": GAP} if end_bind else None,
        "startArrowhead": start, "endArrowhead": end,
    })
    if type_ == "arrow":
        e["elbowed"] = False
    return e


# ---------------------------------------------------------------- geometry

def center(n: dict) -> tuple[float, float]:
    return n["x"] + n["w"] / 2, n["y"] + n["h"] / 2


def bbox(n: dict) -> tuple[float, float, float, float]:
    return n["x"], n["y"], n["x"] + n["w"], n["y"] + n["h"]


def exit_point(n: dict, toward: tuple[float, float]) -> tuple[float, float]:
    """Where a ray from n's centre toward a point leaves n's outline, plus GAP."""
    cx, cy = center(n)
    dx, dy = toward[0] - cx, toward[1] - cy
    dist = math.hypot(dx, dy)
    if dist == 0:
        return cx, cy
    ux, uy = dx / dist, dy / dist
    a, b = n["w"] / 2, n["h"] / 2
    shape = n["shape"]
    if shape == "ellipse":
        t = 1 / math.sqrt((ux / a) ** 2 + (uy / b) ** 2) if a and b else 0
    elif shape == "diamond":
        t = 1 / (abs(ux) / a + abs(uy) / b) if a and b else 0
    else:
        cands = []
        if ux:
            cands.append(abs(a / ux))
        if uy:
            cands.append(abs(b / uy))
        t = min(cands) if cands else 0
    t += GAP
    return cx + ux * t, cy + uy * t


def overlap_area(a: dict, b: dict, tol: float = 1.0) -> float:
    ax0, ay0, ax1, ay1 = bbox(a)
    bx0, by0, bx1, by1 = bbox(b)
    w = min(ax1, bx1) - max(ax0, bx0) - tol
    h = min(ay1, by1) - max(ay0, by0) - tol
    return w * h if w > 0 and h > 0 else 0


def seg_hits_box(p: tuple[float, float], q: tuple[float, float], n: dict, deflate: float = 3) -> bool:
    """Liang-Barsky clip of segment pq against n's bbox shrunk by `deflate`."""
    x0, y0, x1, y1 = bbox(n)
    x0, y0, x1, y1 = x0 + deflate, y0 + deflate, x1 - deflate, y1 - deflate
    if x0 >= x1 or y0 >= y1:
        return False
    dx, dy = q[0] - p[0], q[1] - p[1]
    t0, t1 = 0.0, 1.0
    for pk, qk in ((-dx, p[0] - x0), (dx, x1 - p[0]), (-dy, p[1] - y0), (dy, y1 - p[1])):
        if pk == 0:
            if qk < 0:
                return False
            continue
        t = qk / pk
        if pk < 0:
            t0 = max(t0, t)
        else:
            t1 = min(t1, t)
        if t0 > t1:
            return False
    return True


def polyline_midpoint(pts: list[tuple[float, float]]) -> tuple[float, float]:
    total = sum(math.dist(pts[i], pts[i + 1]) for i in range(len(pts) - 1))
    if total == 0:
        return pts[0]
    walk = total / 2
    for i in range(len(pts) - 1):
        seg = math.dist(pts[i], pts[i + 1])
        if walk <= seg:
            f = walk / seg if seg else 0
            return pts[i][0] + (pts[i + 1][0] - pts[i][0]) * f, pts[i][1] + (pts[i + 1][1] - pts[i][1]) * f
        walk -= seg
    return pts[-1]


# ---------------------------------------------------------------- build

class Builder:
    def __init__(self, spec: dict):
        self.spec = spec
        self.sketch = spec.get("style", "clean") == "sketch"
        self.rough = 1 if self.sketch else 0
        self.font = FONT_SKETCH if self.sketch else FONT_MONO
        d = spec.get("defaults", {})
        self.font_size = d.get("fontSize", 16)
        self.bg = d.get("background", PALETTE["background"])
        g = spec.get("grid", {})
        self.grid = (g.get("originX", 0), g.get("originY", 0), g.get("colWidth", 260), g.get("rowHeight", 150))
        self.nodes: dict[str, dict] = {}      # id -> node (anything an edge can bind to, or a section)
        self.order: list[str] = []            # element ids in draw order (sections prepended later)
        self.elements: dict[str, dict] = {}   # id -> Excalidraw element
        self.section_ids: list[str] = []
        self.edge_geo: list[dict] = []
        self.problems: list[tuple[str, str]] = []
        self.edge_ids: dict[str, int] = {}

    # --- bookkeeping
    def warn(self, msg: str) -> None:
        self.problems.append(("warn", msg))

    def note(self, msg: str) -> None:
        self.problems.append(("note", msg))

    def add(self, el: dict, node: dict | None = None) -> None:
        if el["id"] in self.elements:
            raise SpecError(f"duplicate id '{el['id']}'")
        self.elements[el["id"]] = el
        self.order.append(el["id"])
        if node is not None:
            node.setdefault("id", el["id"])
            self.nodes[el["id"]] = node

    def role(self, e: dict, default: str = "neutral") -> dict:
        name = e.get("role", default)
        if name not in PALETTE["roles"]:
            raise SpecError(f"'{e.get('id')}': unknown role '{name}' (see --roles)")
        return dict(PALETTE["roles"][name], name=name)

    def level(self, e: dict, default: str) -> dict:
        name = e.get("level", default)
        if name not in PALETTE["text"]:
            raise SpecError(f"'{e.get('id')}': unknown text level '{name}' (see --roles)")
        lv = dict(PALETTE["text"][name])
        if "size" in e:
            lv["size"] = e["size"]
        return lv

    def place(self, e: dict, w: float, h: float, centre_given: bool = False) -> tuple[float, float]:
        if "x" in e and "y" in e:
            if centre_given:
                return e["x"] - w / 2, e["y"] - h / 2
            return e["x"], e["y"]
        if "col" in e and "row" in e:
            ox, oy, cw, rh = self.grid
            cx, cy = ox + (e["col"] + 0.5) * cw, oy + (e["row"] + 0.5) * rh
            return cx - w / 2, cy - h / 2
        raise SpecError(f"'{e.get('id')}': needs x,y or col,row")

    # --- element kinds
    def build(self) -> None:
        seen: set[str] = set()
        for e in self.spec.get("elements", []):
            id_ = e.get("id")
            if not id_ or not isinstance(id_, str):
                raise SpecError(f"element without a string id: {json.dumps(e)[:80]}")
            if id_ in seen:
                raise SpecError(f"duplicate id '{id_}'")
            seen.add(id_)
            kind = e.get("kind")
            if kind not in KINDS:
                raise SpecError(f"'{id_}': unknown kind '{kind}' (one of {', '.join(KINDS)})")
            if kind == "section":
                self.section_ids.append(id_)
                continue
            getattr(self, "k_" + kind)(e)
        self.build_sections()
        self.build_title()
        for e in self.spec.get("edges", []):
            self.build_edge(e)
        self.check()

    def k_box(self, e: dict) -> None:
        self.shape(e, "box")

    def k_ellipse(self, e: dict) -> None:
        self.shape(e, "ellipse")

    def k_diamond(self, e: dict) -> None:
        self.shape(e, "diamond")

    def shape(self, e: dict, kind: str) -> None:
        role = self.role(e)
        text = str(e.get("text", ""))
        size = e.get("size", self.font_size)
        tw, th, _ = measure(text, size) if text else (0, 0, [])
        w, h = shape_size(kind, tw, th) if text else (MIN_W, MIN_H)
        w, h = e.get("w", w), e.get("h", h)
        if text:
            rw, rh = text_room(kind, w, h)
            if tw > rw + 0.5 or th > rh + 0.5:
                self.warn(f"text overflow: '{e['id']}' needs {tw}x{th} px of text room, has {rw:.0f}x{rh:.0f}")
        x, y = self.place(e, w, h)
        dashed = bool(e.get("dashed")) or role["name"] == "inactive"
        el = base(e["id"], {"box": "rectangle"}.get(kind, kind), x, y, w, h, stroke=role["stroke"],
                  fill=role["fill"], stroke_width=3 if e.get("bold") else 2, dashed=dashed,
                  roughness=self.rough, roundness=ROUNDNESS_ADAPTIVE if kind == "box" else None)
        node = {"kind": kind, "x": x, "y": y, "w": w, "h": h,
                "shape": "rect" if kind == "box" else kind, "role": role, "edges": 0}
        self.add(el, node)
        if text:
            tid = e["id"] + "__text"
            el["boundElements"].append({"id": tid, "type": "text"})
            self.add(text_el(tid, x + (w - tw) / 2, y + (h - th) / 2, tw, th, text, size, role["text"],
                             self.font, container=e["id"], roughness=self.rough),
                     {"kind": "boundtext", "x": x + (w - tw) / 2, "y": y + (h - th) / 2, "w": tw, "h": th,
                      "shape": "rect", "container": e["id"]})

    def k_text(self, e: dict) -> None:
        lv = self.level(e, "body")
        text = str(e.get("text", ""))
        if not text:
            raise SpecError(f"'{e['id']}': text element without text")
        tw, th, _ = measure(text, lv["size"])
        w, h = e.get("w", tw), e.get("h", th)
        x, y = self.place(e, w, h)
        align = e.get("align", "left")
        if align not in ("left", "center", "right"):
            raise SpecError(f"'{e['id']}': align must be left, center or right")
        self.add(text_el(e["id"], x, y, w, h, text, lv["size"], lv["color"], self.font, align=align,
                         valign="top", roughness=self.rough),
                 {"kind": "text", "x": x, "y": y, "w": w, "h": h, "shape": "rect"})

    def k_code(self, e: dict) -> None:
        cfg = PALETTE["code"]
        text = str(e.get("text", ""))
        if not text:
            raise SpecError(f"'{e['id']}': code element without text")
        size = e.get("size", cfg["size"])
        tw, th, _ = measure(text, size)
        w, h = e.get("w", tw + 2 * CODE_PAD), e.get("h", th + 2 * CODE_PAD)
        if tw > w - 2 * CODE_PAD or th > h - 2 * CODE_PAD:
            self.warn(f"text overflow: code '{e['id']}' needs {tw + 2 * CODE_PAD}x{th + 2 * CODE_PAD} px")
        x, y = self.place(e, w, h)
        grp = e["id"] + "__grp"
        color = cfg["text"].get(e.get("lang", "default"), cfg["text"]["default"])
        self.add(base(e["id"], "rectangle", x, y, w, h, stroke=cfg["stroke"], fill=cfg["fill"],
                      roughness=self.rough, roundness=ROUNDNESS_ADAPTIVE, group=grp),
                 {"kind": "code", "x": x, "y": y, "w": w, "h": h, "shape": "rect", "edges": 0})
        # Free text in a group rather than bound text: Excalidraw pins bound left-aligned
        # text 5px from the edge, which would not match the padding drawn here.
        self.add(text_el(e["id"] + "__text", x + CODE_PAD, y + CODE_PAD, tw, th, text, size, color,
                         FONT_MONO, align="left", valign="top", roughness=self.rough, group=grp),
                 {"kind": "boundtext", "x": x + CODE_PAD, "y": y + CODE_PAD, "w": tw, "h": th,
                  "shape": "rect", "container": e["id"]})

    def k_dot(self, e: dict) -> None:
        d = e.get("size", DOT)
        x, y = self.place(e, d, d, centre_given=True)
        self.add(self.dot_el(e["id"], x, y, d), {"kind": "dot", "x": x, "y": y, "w": d, "h": d,
                                                  "shape": "ellipse", "edges": 0})

    def dot_el(self, id_: str, x: float, y: float, d: float, group: str | None = None) -> dict:
        return base(id_, "ellipse", x, y, d, d, stroke=PALETTE["dot"]["stroke"], fill=PALETTE["dot"]["fill"],
                    stroke_width=1, roughness=self.rough, group=group)

    def k_line(self, e: dict) -> None:
        self.free_linear(e, "line")

    def k_arrow(self, e: dict) -> None:
        self.free_linear(e, "arrow")

    def free_linear(self, e: dict, type_: str) -> None:
        pts = e.get("points")
        if not isinstance(pts, list) or len(pts) < 2:
            raise SpecError(f"'{e['id']}': {type_} needs at least two points")
        pts = [(float(p[0]), float(p[1])) for p in pts]
        stroke = self.role(e)["stroke"] if "role" in e else PALETTE["line"]
        start, end = e.get("start"), e.get("end", "arrow" if type_ == "arrow" else None)
        for head in (start, end):
            if head not in ARROWHEADS:
                raise SpecError(f"'{e['id']}': arrowhead must be one of {ARROWHEADS}")
        el = linear_el(e["id"], type_, pts, stroke=stroke, width=e.get("width", 2), dashed=bool(e.get("dashed")),
                       start=start, end=end, roughness=self.rough)
        xs = [p[0] for p in pts]
        ys = [p[1] for p in pts]
        self.add(el, {"kind": type_, "x": min(xs), "y": min(ys), "w": max(xs) - min(xs), "h": max(ys) - min(ys),
                      "shape": "rect", "points": pts})

    def k_timeline(self, e: dict) -> None:
        steps = e.get("steps")
        if not isinstance(steps, list) or not steps:
            raise SpecError(f"'{e['id']}': timeline needs a non-empty steps list")
        down = e.get("direction", "down") == "down"
        gap = e.get("gap", TIMELINE_GAP)
        lv = self.level(e, "body")
        n = len(steps)
        # lay out relative to the first dot's centre at (0, 0), then shift
        dots, labels = [], []
        for i, s in enumerate(steps):
            s = str(s)
            tw, th, _ = measure(s, lv["size"])
            cx, cy = (0, i * gap) if down else (i * gap, 0)
            dots.append((cx, cy))
            labels.append((cx + DOT / 2 + 12, cy - th / 2, tw, th, s) if down else (cx - tw / 2, cy + DOT / 2 + 8, tw, th, s))
        line = [(0, -24), (0, (n - 1) * gap + 24)] if down else [(-24, 0), ((n - 1) * gap + 24, 0)]
        xs = [p[0] for p in line] + [lb[0] + lb[2] for lb in labels] + [lb[0] for lb in labels]
        ys = [p[1] for p in line] + [lb[1] + lb[3] for lb in labels] + [lb[1] for lb in labels]
        minx, miny, w, h = min(xs), min(ys), max(xs) - min(xs), max(ys) - min(ys)
        if "x" in e and "y" in e:
            ox, oy = e["x"], e["y"]                      # x,y is the first dot's centre
        else:
            bx, by = self.place(e, w, h)
            ox, oy = bx - minx, by - miny
        grp = e["id"] + "__grp"
        self.add(linear_el(e["id"], "line", [(ox + p[0], oy + p[1]) for p in line], stroke=PALETTE["line"],
                           roughness=self.rough, group=grp),
                 {"kind": "line", "x": ox + minx, "y": oy + miny, "w": w, "h": h, "shape": "rect",
                  "points": [(ox + p[0], oy + p[1]) for p in line]})
        for i, ((cx, cy), (lx, ly, tw, th, s)) in enumerate(zip(dots, labels), start=1):
            did, lid = f"{e['id']}.dot{i}", f"{e['id']}.label{i}"
            self.add(self.dot_el(did, ox + cx - DOT / 2, oy + cy - DOT / 2, DOT, grp),
                     {"kind": "dot", "x": ox + cx - DOT / 2, "y": oy + cy - DOT / 2, "w": DOT, "h": DOT,
                      "shape": "ellipse", "edges": 0, "timeline": e["id"]})
            self.add(text_el(lid, ox + lx, oy + ly, tw, th, s, lv["size"], lv["color"], self.font,
                             align="left" if down else "center", valign="top", roughness=self.rough, group=grp),
                     {"kind": "text", "x": ox + lx, "y": oy + ly, "w": tw, "h": th, "shape": "rect",
                      "timeline": e["id"]})

    def build_sections(self) -> None:
        specs = {e["id"]: e for e in self.spec.get("elements", []) if e.get("kind") == "section"}
        done: dict[str, dict] = {}
        visiting: set[str] = set()

        def resolve(sid: str) -> dict:
            if sid in done:
                return done[sid]
            if sid in visiting:
                raise SpecError(f"section '{sid}' contains itself (directly or through another section)")
            visiting.add(sid)
            e = specs[sid]
            members = e.get("contains")
            if not isinstance(members, list) or not members:
                raise SpecError(f"section '{sid}' needs a non-empty contains list")
            boxes, all_members = [], set()
            for m in members:
                if m in specs:
                    sub = resolve(m)
                    boxes.append(bbox(sub))
                    all_members |= sub["members"] | {m}
                elif m in self.nodes:
                    boxes.append(bbox(self.nodes[m]))
                    all_members.add(m)
                    if self.nodes[m].get("container"):
                        all_members.add(self.nodes[m]["container"])
                else:
                    raise SpecError(f"section '{sid}' contains unknown id '{m}'")
            pad = e.get("pad", SECTION_PAD)
            label = str(e.get("label", ""))
            lv = self.level(e, "section")
            band = math.ceil(lv["size"] * LINE_H) + 12 if label else 0
            x0, y0 = min(b[0] for b in boxes) - pad, min(b[1] for b in boxes) - pad - band
            x1, y1 = max(b[2] for b in boxes) + pad, max(b[3] for b in boxes) + pad
            fill = self.role(e)["fill"] if "role" in e else PALETTE["section"]["fill"]
            grp = sid + "__grp"
            el = base(sid, "rectangle", x0, y0, x1 - x0, y1 - y0, stroke=PALETTE["section"]["stroke"], fill=fill,
                      stroke_width=1, dashed=not e.get("solid", False), roughness=self.rough,
                      roundness=ROUNDNESS_ADAPTIVE, group=grp)
            node = {"kind": "section", "x": x0, "y": y0, "w": x1 - x0, "h": y1 - y0, "shape": "rect",
                    "members": all_members, "edges": 0}
            self.elements[sid] = el
            self.nodes[sid] = node
            node["id"] = sid
            if label:
                tw, th, _ = measure(label, lv["size"])
                self.elements[sid + "__label"] = text_el(sid + "__label", x0 + 12, y0 + 8, tw, th, label, lv["size"],
                                                         lv["color"], self.font, align="left", valign="top",
                                                         roughness=self.rough, group=grp)
            done[sid] = node
            visiting.discard(sid)
            return node

        for sid in self.section_ids:
            resolve(sid)
        # sections draw behind everything: outer before inner, so inner stays visible
        ordered = sorted(self.section_ids, key=lambda s: -self.nodes[s]["w"] * self.nodes[s]["h"])
        front = []
        for sid in ordered:
            front.append(sid)
            if sid + "__label" in self.elements:
                front.append(sid + "__label")
        self.order = front + self.order

    def build_title(self) -> None:
        t = self.spec.get("title")
        if not t:
            return
        lv = dict(PALETTE["text"]["title"])
        text = t if isinstance(t, str) else str(t.get("text", ""))
        tw, th, _ = measure(text, lv["size"])
        if isinstance(t, dict) and "x" in t and "y" in t:
            x, y = t["x"], t["y"]
        else:
            if not self.nodes:
                x, y = 0, 0
            else:
                x = min(n["x"] for n in self.nodes.values())
                y = min(n["y"] for n in self.nodes.values()) - th - 28
        self.add(text_el("title", x, y, tw, th, text, lv["size"], lv["color"], self.font, align="left",
                         valign="top", roughness=self.rough),
                 {"kind": "text", "x": x, "y": y, "w": tw, "h": th, "shape": "rect"})

    def build_edge(self, e: dict) -> None:
        src, dst = e.get("from"), e.get("to")
        for end_id in (src, dst):
            if end_id not in self.nodes:
                raise SpecError(f"edge {src!r} -> {dst!r}: unknown id '{end_id}'")
            if self.nodes[end_id]["kind"] in ("line", "arrow", "boundtext"):
                raise SpecError(f"edge {src!r} -> {dst!r}: '{end_id}' cannot be an endpoint (bind to a shape, text, dot, or section)")
        a, b = self.nodes[src], self.nodes[dst]
        base_id = e.get("id") or f"{src}->{dst}"
        n = self.edge_ids.get(base_id, 0) + 1
        self.edge_ids[base_id] = n
        aid = base_id if n == 1 else f"{base_id}#{n}"
        via = [(float(p[0]), float(p[1])) for p in e.get("via", [])]
        pts = [center(a)] + via + [center(b)]
        p0 = exit_point(a, pts[1])
        pn = exit_point(b, pts[-2])
        pts = [p0] + via + [pn]
        stroke = self.role(e)["stroke"] if "role" in e else a.get("role", PALETTE["roles"]["neutral"])["stroke"]
        if a["kind"] in ("text", "dot", "code", "section"):
            stroke = PALETTE["line"] if "role" not in e else stroke
        start, end = e.get("start"), e.get("end", "arrow")
        for head in (start, end):
            if head not in ARROWHEADS:
                raise SpecError(f"edge {aid}: arrowhead must be one of {ARROWHEADS}")
        el = linear_el(aid, "arrow", pts, stroke=stroke, width=3 if e.get("bold") else e.get("width", 2),
                       dashed=bool(e.get("dashed")), start=start, end=end, start_bind=src, end_bind=dst,
                       roughness=self.rough)
        self.add(el)
        for end_id in (src, dst):
            self.elements[end_id]["boundElements"].append({"id": aid, "type": "arrow"})
            self.nodes[end_id]["edges"] = self.nodes[end_id].get("edges", 0) + 1
        geo = {"id": aid, "from": src, "to": dst, "points": pts, "label": None}
        label = e.get("label")
        if label:
            lv = PALETTE["text"]["label"]
            tw, th, _ = measure(str(label), lv["size"])
            mx, my = polyline_midpoint(pts)
            lid = aid + "__label"
            el["boundElements"].append({"id": lid, "type": "text"})
            self.add(text_el(lid, mx - tw / 2, my - th / 2, tw, th, str(label), lv["size"], lv["color"], self.font,
                             container=aid, roughness=self.rough),
                     {"kind": "label", "x": mx - tw / 2, "y": my - th / 2, "w": tw, "h": th, "shape": "rect",
                      "edge": aid})
            geo["label"] = lid
        self.edge_geo.append(geo)

    # --- checks
    def check(self) -> None:
        solid = {i: n for i, n in self.nodes.items()
                 if n["kind"] in ("box", "ellipse", "diamond", "text", "code", "dot", "label")}
        ids = list(solid)
        for i, a in enumerate(ids):
            for b in ids[i + 1:]:
                na, nb = solid[a], solid[b]
                if na.get("edge") and nb.get("edge") and na["edge"] == nb["edge"]:
                    continue
                area = overlap_area(na, nb)
                if area:
                    self.warn(f"overlap: '{a}' and '{b}' share {area:.0f} px²")
        for g in self.edge_geo:
            pts = g["points"]
            for i in range(len(pts) - 1):
                for nid, n in solid.items():
                    if nid in (g["from"], g["to"], g["label"]):
                        continue
                    if n.get("container") in (g["from"], g["to"]):
                        continue
                    if seg_hits_box(pts[i], pts[i + 1], n):
                        self.warn(f"arrow {g['id']} crosses '{nid}'; add a via waypoint or move one of them")
        for sid in self.section_ids:
            s = self.nodes[sid]
            for nid, n in solid.items():
                if nid in s["members"] or n.get("container") in s["members"] or n.get("timeline") in s["members"]:
                    continue
                if n.get("edge") or nid == "title":
                    continue
                if overlap_area(n, s, tol=0):
                    self.warn(f"'{nid}' sits inside section '{sid}' but is not in its contains list")
        in_section = set()
        for sid in self.section_ids:
            in_section |= self.nodes[sid]["members"]
        lonely = [i for i, n in self.nodes.items()
                  if n["kind"] in SHAPE_KINDS and n.get("edges", 0) == 0 and i not in in_section]
        for i in lonely:
            self.note(f"'{i}' has no edges and no section; position alone shows no relationship")
        shapes = [n for n in self.nodes.values() if n["kind"] in SHAPE_KINDS + ("code",)]
        if len(shapes) >= 6 and all(n["kind"] == "box" for n in shapes):
            self.note(f"{len(shapes)} boxes and nothing else: reads as a card grid, not an argument")

    # --- outputs
    def excalidraw(self) -> dict:
        return {
            "type": "excalidraw", "version": 2, "source": "draw-diagram",
            "elements": [self.elements[i] for i in self.order],
            "appState": {"viewBackgroundColor": self.bg, "gridSize": 20, "gridModeEnabled": False},
            "files": {},
        }

    def svg(self, pad: int = 40) -> tuple[str, int, int]:
        els = [self.elements[i] for i in self.order]
        xs, ys = [], []
        for e in els:
            xs += [e["x"], e["x"] + e["width"]]
            ys += [e["y"], e["y"] + e["height"]]
        if not xs:
            xs, ys = [0, 200], [0, 100]
        x0, y0 = min(xs) - pad, min(ys) - pad
        W, H = math.ceil(max(xs) - min(xs) + 2 * pad), math.ceil(max(ys) - min(ys) + 2 * pad)
        mono = "'Cascadia Code','Cascadia Mono',Menlo,Consolas,'DejaVu Sans Mono',monospace"
        hand = "Excalifont,Virgil,'Comic Sans MS','Segoe Print',cursive"
        out = [f'<svg xmlns="http://www.w3.org/2000/svg" width="{W}" height="{H}" viewBox="{rnd(x0)} {rnd(y0)} {W} {H}">',
               f'<rect x="{rnd(x0)}" y="{rnd(y0)}" width="{W}" height="{H}" fill="{self.bg}"/>']
        for e in els:
            sid = esc(e["id"])
            dash = ' stroke-dasharray="8 6"' if e["strokeStyle"] == "dashed" else ""
            st = f'stroke="{e["strokeColor"]}" stroke-width="{e["strokeWidth"]}"{dash}'
            x, y, w, h = e["x"], e["y"], e["width"], e["height"]
            t = e["type"]
            if t == "rectangle":
                m = min(w, h)
                rx = m * 0.25 if m <= 128 else 32   # Excalidraw's adaptive corner radius
                fill = e["backgroundColor"]
                out.append(f'<rect id="{sid}" x="{rnd(x)}" y="{rnd(y)}" width="{rnd(w)}" height="{rnd(h)}" rx="{rnd(rx)}" fill="{fill}" {st}/>')
            elif t == "ellipse":
                out.append(f'<ellipse id="{sid}" cx="{rnd(x + w / 2)}" cy="{rnd(y + h / 2)}" rx="{rnd(w / 2)}" ry="{rnd(h / 2)}" fill="{e["backgroundColor"]}" {st}/>')
            elif t == "diamond":
                p = f'{rnd(x + w / 2)},{rnd(y)} {rnd(x + w)},{rnd(y + h / 2)} {rnd(x + w / 2)},{rnd(y + h)} {rnd(x)},{rnd(y + h / 2)}'
                out.append(f'<polygon id="{sid}" points="{p}" fill="{e["backgroundColor"]}" {st}/>')
            elif t in ("line", "arrow"):
                pts = [(x + px, y + py) for px, py in e["points"]]
                out.append(f'<polyline id="{sid}" points="{" ".join(f"{rnd(a)},{rnd(b)}" for a, b in pts)}" fill="none" {st} stroke-linejoin="round" stroke-linecap="round"/>')
                out += arrowheads(e, pts)
            elif t == "text":
                font = hand if e["fontFamily"] != FONT_MONO else mono
                fs = e["fontSize"]
                anchor = {"left": "start", "center": "middle", "right": "end"}[e["textAlign"]]
                ax = {"left": x, "center": x + w / 2, "right": x + w}[e["textAlign"]]
                cont = self.elements.get(e["containerId"] or "")
                if cont and cont["type"] == "arrow":
                    out.append(f'<rect x="{rnd(x - 4)}" y="{rnd(y - 2)}" width="{rnd(w + 8)}" height="{rnd(h + 4)}" fill="{self.bg}"/>')
                for i, line in enumerate(e["text"].split("\n")):
                    ly = y + i * fs * LINE_H + (LINE_H - 1) * fs / 2
                    out.append(f'<text x="{rnd(ax)}" y="{rnd(ly)}" font-family="{font}" font-size="{fs}" fill="{e["strokeColor"]}" text-anchor="{anchor}" dominant-baseline="hanging" xml:space="preserve">{esc(line)}</text>')
        out.append("</svg>")
        return "\n".join(out), W, H


def esc(s: str) -> str:
    return s.replace("&", "&amp;").replace("<", "&lt;").replace(">", "&gt;").replace('"', "&quot;")


def arrowheads(e: dict, pts: list[tuple[float, float]]) -> list[str]:
    out = []
    for head, tip, prev in ((e.get("endArrowhead"), pts[-1], pts[-2]), (e.get("startArrowhead"), pts[0], pts[1])):
        if not head:
            continue
        dx, dy = tip[0] - prev[0], tip[1] - prev[1]
        d = math.hypot(dx, dy) or 1
        ux, uy = dx / d, dy / d
        size = 16
        col = e["strokeColor"]
        if head == "dot":
            out.append(f'<circle cx="{rnd(tip[0])}" cy="{rnd(tip[1])}" r="5" fill="{col}"/>')
            continue
        if head == "bar":
            out.append(f'<line x1="{rnd(tip[0] - uy * 8)}" y1="{rnd(tip[1] + ux * 8)}" x2="{rnd(tip[0] + uy * 8)}" y2="{rnd(tip[1] - ux * 8)}" stroke="{col}" stroke-width="{e["strokeWidth"]}"/>')
            continue
        ang = math.radians(25)
        l = (tip[0] - size * (ux * math.cos(ang) - uy * math.sin(ang)), tip[1] - size * (uy * math.cos(ang) + ux * math.sin(ang)))
        r_ = (tip[0] - size * (ux * math.cos(ang) + uy * math.sin(ang)), tip[1] - size * (uy * math.cos(ang) - ux * math.sin(ang)))
        p = f'{rnd(l[0])},{rnd(l[1])} {rnd(tip[0])},{rnd(tip[1])} {rnd(r_[0])},{rnd(r_[1])}'
        if head == "triangle":
            out.append(f'<polygon points="{p}" fill="{col}"/>')
        else:
            out.append(f'<polyline points="{p}" fill="none" stroke="{col}" stroke-width="{e["strokeWidth"]}" stroke-linejoin="round" stroke-linecap="round"/>')
    return out


# ---------------------------------------------------------------- png

def find_chromium() -> str | None:
    for var in ("CHROME_BIN", "CHROMIUM_BIN"):
        if os.environ.get(var) and os.access(os.environ[var], os.X_OK):
            return os.environ[var]
    roots = [os.environ.get("PLAYWRIGHT_BROWSERS_PATH", ""), os.path.expanduser("~/.cache/ms-playwright"),
             os.path.expanduser("~/Library/Caches/ms-playwright"), os.path.expanduser("~/AppData/Local/ms-playwright")]
    for root in filter(None, roots):
        for pat in ("chromium_headless_shell-*/chrome-linux/headless_shell",
                    "chromium_headless_shell-*/chrome-mac*/headless_shell",
                    "chromium_headless_shell-*/chrome-win/headless_shell.exe",
                    "chromium-*/chrome-linux/chrome", "chromium-*/chrome-mac*/Chromium.app/Contents/MacOS/Chromium",
                    "chromium-*/chrome-win/chrome.exe"):
            hits = sorted(glob.glob(os.path.join(root, pat)))
            if hits:
                return hits[-1]
    for name in ("chromium", "chromium-browser", "google-chrome", "google-chrome-stable", "chrome"):
        p = shutil.which(name)
        if p:
            return p
    for p in ("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
              "/Applications/Chromium.app/Contents/MacOS/Chromium"):
        if os.access(p, os.X_OK):
            return p
    return None


def render_png(svg_path: Path, png_path: Path, w: int, h: int, scale: int) -> str | None:
    """Returns None on success, else the reason the PNG was not written."""
    chrome = find_chromium()
    if not chrome:
        return "no Chromium found (set CHROME_BIN, or install Playwright's chromium)"
    if w * scale > 6000 or h * scale > 6000:
        scale = 1

    def shoot(win_h: int) -> str | None:
        cmd = [chrome, "--headless=new", "--no-sandbox", "--disable-gpu", "--hide-scrollbars",
               f"--force-device-scale-factor={scale}", f"--window-size={w},{win_h}", f"--screenshot={png_path}",
               svg_path.resolve().as_uri()]
        try:
            subprocess.run(cmd, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL, timeout=90, check=False)
        except (OSError, subprocess.TimeoutExpired) as e:
            return f"chromium failed: {e}"
        if not png_path.exists() or png_path.stat().st_size == 0:
            return f"chromium wrote no file ({chrome})"
        return None

    why = shoot(h)
    if why:
        return why
    # Chromium's new headless mode keeps room in the window for browser UI and
    # screenshots only the viewport below it, so the first shot can come back
    # short. Measure the PNG and shoot again with the missing height added.
    _, got_h = png_size(png_path)
    if got_h and got_h < h * scale:
        why = shoot(h + math.ceil((h * scale - got_h) / scale))
    return why


def png_size(path: Path) -> tuple[int, int]:
    with path.open("rb") as f:
        head = f.read(24)
    if head[:8] != b"\x89PNG\r\n\x1a\n":
        return 0, 0
    return struct.unpack(">II", head[16:24])


# ---------------------------------------------------------------- cli

def main(argv: list[str] | None = None) -> int:
    ap = argparse.ArgumentParser(description=__doc__.split("\n\n")[0])
    ap.add_argument("spec", nargs="?", help="diagram spec (JSON)")
    ap.add_argument("--out", help="output base path, without extension (default: next to the spec)")
    ap.add_argument("--no-png", action="store_true", help="skip the Chromium render")
    ap.add_argument("--scale", type=int, default=2, help="PNG device scale factor (default 2)")
    ap.add_argument("--roles", action="store_true", help="list palette roles and text levels")
    ap.add_argument("--example", action="store_true", help="print a small valid spec")
    args = ap.parse_args(argv)

    if args.roles:
        print("roles (shape colour by meaning):")
        for k, v in PALETTE["roles"].items():
            print(f"  {k:<9} {v['about']}")
        print("text levels:")
        for k, v in PALETTE["text"].items():
            print(f"  {k:<9} {v['size']}px")
        return 0
    if args.example:
        print(json.dumps(EXAMPLE, indent=2))
        return 0
    if not args.spec:
        ap.error("spec path required (or --roles / --example)")

    spec_path = Path(args.spec)
    try:
        spec = json.loads(spec_path.read_text(encoding="utf-8"))
        b = Builder(spec)
        b.build()
    except (OSError, json.JSONDecodeError, SpecError) as e:
        print(f"error: {e}", file=sys.stderr)
        return 1

    if args.out:
        out = Path(args.out)
    else:
        name = spec_path.name
        for suffix in (".spec.json", ".json"):
            if name.endswith(suffix):
                name = name[: -len(suffix)]
                break
        out = spec_path.with_name(name)
    out.parent.mkdir(parents=True, exist_ok=True)

    doc = b.excalidraw()
    body = ",\n".join("    " + json.dumps(el, ensure_ascii=False) for el in doc["elements"])
    head = {k: v for k, v in doc.items() if k != "elements"}
    text = "{\n" + ",\n".join(f'  {json.dumps(k)}: {json.dumps(v)}' for k, v in head.items() if k in ("type", "version", "source"))
    text += ',\n  "elements": [\n' + body + "\n  ],\n"
    text += f'  "appState": {json.dumps(doc["appState"])},\n  "files": {{}}\n}}\n'
    exc_path = out.with_suffix(".excalidraw")
    exc_path.write_text(text, encoding="utf-8")

    svg, W, H = b.svg()
    svg_path = out.with_suffix(".svg")
    svg_path.write_text(svg, encoding="utf-8")
    wrote = [f"{exc_path} ({len(doc['elements'])} elements)", str(svg_path)]
    if not args.no_png:
        png_path = out.with_suffix(".png")
        why = render_png(svg_path, png_path, W, H, args.scale)
        wrote.append(str(png_path) if why is None else f"no png: {why}")
    print("wrote " + ", ".join(wrote))

    warns = [m for lv, m in b.problems if lv == "warn"]
    notes = [m for lv, m in b.problems if lv == "note"]
    print(f"checks: 0 errors, {len(warns)} warnings, {len(notes)} notes; canvas {W}x{H} px")
    for m in warns:
        print("  warn", m)
    for m in notes:
        print("  note", m)
    return 0


if __name__ == "__main__":
    sys.exit(main())
