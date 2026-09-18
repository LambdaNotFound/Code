# Diagram spec format

One JSON file per diagram. `build_excalidraw.py` turns it into the
`.excalidraw` file, so every coordinate the compiler can derive, you
leave out. `python3 build_excalidraw.py --example` prints a complete
spec that uses most of what is below.

```json
{
  "title": "Webhook intake",
  "style": "clean",
  "grid": {"colWidth": 260, "rowHeight": 150, "originX": 0, "originY": 0},
  "defaults": {"fontSize": 16, "background": "#ffffff"},
  "elements": [ ... ],
  "edges": [ ... ]
}
```

| Key | Meaning |
|---|---|
| `title` | Optional. A string is placed above everything, left-aligned. An object `{"text", "x", "y"}` is placed where you say. |
| `style` | `clean` (default: crisp strokes, monospace) or `sketch` (rough strokes, hand-drawn font). |
| `grid` | Cell size for `col`/`row` placement. All optional; defaults shown. |
| `defaults.fontSize` | Text size inside shapes. Free text uses its `level` instead. |
| `defaults.background` | Canvas colour written to the file. |

## Placing an element

Every element takes one of two position forms:

- `"x", "y"`: absolute, top-left corner, in Excalidraw pixels. For a
  `dot` it is the centre; for a `timeline` it is the centre of the
  first dot.
- `"col", "row"`: the element's **centre** goes at the centre of that
  grid cell. Fractions are fine (`"row": 0.5`). This is the form to
  reach for: a row of shapes with different sizes lines up on one axis
  without arithmetic.

Widen `colWidth` or `rowHeight` when the checker reports overlaps
rather than nudging elements one at a time.

## Element kinds

Common to all: `id` (required, unique, descriptive: edges and checker
messages quote it). Shapes and edges take `role` for colour and
`dashed` for stroke style.

| kind | Required | Optional | Draws |
|---|---|---|---|
| `box` | `text` | `role`, `w`, `h`, `size`, `bold`, `dashed` | rounded rectangle with the text centred inside |
| `ellipse` | `text` | same as box | ellipse; use for start and end |
| `diamond` | `text` | same as box | decision; keep the text short |
| `text` | `text` | `level`, `align`, `size`, `w`, `h` | free-floating text, no container |
| `code` | `text` | `lang`, `size`, `w`, `h` | dark block, left-aligned monospace; `lang: json` colours it green |
| `dot` | — | `size` (default 12) | a small filled circle; a marker or an arrow target |
| `line` | `points` | `dashed`, `width`, `role`, `start`, `end` | polyline through absolute points; structure, not a relationship |
| `arrow` | `points` | same as line | an unbound arrow through absolute points, for annotations |
| `timeline` | `steps` | `direction` (`down` default, or `right`), `gap` (70), `level` | a line, one dot per step, one label per step |
| `section` | `contains` | `label`, `pad` (32), `role`, `solid` | a dashed rectangle drawn behind its members, sized to enclose them |

Notes on kinds:

- Text is not wrapped. Break lines yourself with `\n`. `text` holds
  only readable words.
- `w` and `h` are exact when given. The compiler sizes a shape to its
  text (a box adds 48 by 32 px of padding; an ellipse scales that by
  root two; a diamond doubles it), so leave them out unless you want a
  shape bigger than its text. Forcing a size smaller than the text
  draws a warning.
- `level` for free text and timeline labels is one of `title`,
  `subtitle`, `section`, `body` (default), `detail`, `label`; each has
  a size and colour in `palette.json`. `size` overrides the size only.
- A `section` may contain other sections. Members must exist; a
  section containing itself is an error. Sections draw behind
  everything, outer before inner.
- A `timeline` expands into `<id>` (the line), `<id>.dot1`,
  `<id>.dot2`, ... and `<id>.label1`, ... Edges may end on a dot.

## Roles

`role` names what a shape means, and the palette decides the colour.
Never write a hex colour in a spec.

```
python3 build_excalidraw.py --roles
```

Roles: `neutral` (default), `start`, `end`, `decision`, `ai`, `data`,
`external`, `error`, `inactive` (drawn dashed). The palette and what
each role is for live in `palette.json` next to the compiler.

## Edges

```json
{"from": "ok", "to": "queue", "label": "yes", "dashed": true,
 "via": [[640, 300]], "start": null, "end": "arrow", "role": "error", "bold": true, "id": "ok_yes"}
```

| Key | Meaning |
|---|---|
| `from`, `to` | Ids of the two ends: a shape, free text, a code block, a dot, a timeline dot, or a section. Not a line. |
| `label` | Text bound to the arrow at its midpoint, on a canvas-coloured patch. |
| `via` | Absolute waypoints the arrow passes through, in order. Route around a shape with one or two. |
| `start`, `end` | Arrowheads: `null`, `arrow` (default for `end`), `bar`, `dot`, `triangle`. |
| `role` | Colour override. Default: the source shape's role stroke; grey from text, dots, code, and sections. |
| `dashed`, `bold`, `width` | Stroke style. `bold` is width 3. |
| `id` | Optional. Default `<from>-><to>`, with `#2`, `#3` appended for repeats. |

The compiler computes where each arrow leaves its source outline and
meets its target outline (box, ellipse, and diamond each get their own
intersection), stops 4 px short, and binds both ends so the arrow
follows the shapes when the user drags them in Excalidraw.

## Generated ids

Ids you did not write but may see in checker messages:

| Id | What it is |
|---|---|
| `<shape>__text` | the text bound inside a shape |
| `<code>__text` | the text of a code block (grouped with its block, not bound) |
| `<section>__label` | a section's label |
| `<from>-><to>` and `<from>-><to>#2` | arrows, in spec order |
| `<arrow>__label` | an arrow's label |
| `<timeline>.dotN`, `<timeline>.labelN` | timeline parts, N from 1 |
| `title` | the title |

## What the compiler checks

Errors stop the build and nothing is written:

- an element without a string id, a duplicate id, an unknown `kind`,
  `role`, `level`, `align`, or arrowhead
- an element with neither `x, y` nor `col, row`
- an edge whose end is unknown, or is a line
- a section with an empty or unknown `contains`, or a cycle

Warnings are written to the report; fix them in the spec or explain
why one stands:

- **overlap**: two drawn elements share area (sections and an arrow's
  own label excepted)
- **text overflow**: an explicit `w` or `h` is smaller than the text
- **arrow crosses**: a straight segment of an arrow passes through a
  shape that is neither its source nor its target; add `via`
- **sits inside section**: an element lies inside a section's
  rectangle but is not in its `contains`

Notes are for your judgement, not defects: a shape with no edges and
no section, or a canvas of six or more boxes and nothing else.

## Outputs

```
python3 build_excalidraw.py diagrams/<slug>/<slug>.spec.json [--out BASE] [--no-png] [--scale 2]
```

Writes `<base>.excalidraw` (the deliverable), `<base>.svg` (a preview
drawn by the compiler from the same coordinates), and `<base>.png`
(the SVG rasterised by a local Chromium: `CHROME_BIN`, Playwright's
browser cache, or a `chromium`/`google-chrome` on `PATH`; skipped with
a message when none is found). The base defaults to the spec path
with `.spec.json` or `.json` removed.

Same spec, same bytes: seeds derive from ids, so a rebuild diffs
cleanly in git.

## Fidelity

The `.excalidraw` carries what Excalidraw itself reads on load:
element geometry, text with `fontFamily` 3 (Cascadia, monospace) or 5
(Excalifont) and a `lineHeight` of 1.25, text bound to containers,
arrows bound at both ends, and groups for code blocks, sections, and
timelines. Excalidraw repairs bindings on load and keeps the text
dimensions written here.

Text width is estimated at 0.6 em per character. Cascadia advances
0.586 em, so text sits slightly looser than estimated and never
clips. In `sketch` style the hand-drawn font is proportional and the
estimate is rougher; expect to widen a box or two after looking.

The SVG and PNG are the compiler's own drawing of the same
coordinates, not Excalidraw's renderer: shapes, arrowheads, dashes,
and text positions match; stroke texture and glyph shapes do not. Use
them to judge layout. The file open in Excalidraw is the ground truth.
