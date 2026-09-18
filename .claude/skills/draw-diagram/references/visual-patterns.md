# Visual patterns

One spec per pattern, each complete and buildable as written; the
compiler's tests build every block on this page and fail if one
draws a warning. Copy the shape, replace the words, keep the spacing.

The common thread: the form is the claim. A reader who cannot read
the labels should still see one-to-many, a sequence, a loop, a
boundary.

## Fan-out (one to many)

A single source on the left, its targets stacked in the next column,
source centred on the stack. Sources, root causes, dispatch.

```json
{"grid": {"colWidth": 280, "rowHeight": 110},
 "elements": [
  {"id": "brief", "kind": "ellipse", "role": "start", "text": "signed-off\nbrief", "col": 0, "row": 1},
  {"id": "design", "kind": "box", "text": "design loop", "col": 1, "row": 0},
  {"id": "pr", "kind": "box", "text": "PR loop", "col": 1, "row": 1},
  {"id": "audit", "kind": "box", "text": "editorial pass", "col": 1, "row": 2}],
 "edges": [{"from": "brief", "to": "design"}, {"from": "brief", "to": "pr"}, {"from": "brief", "to": "audit"}]}
```

## Convergence (many to one)

The mirror image: a stacked column feeding one target. Aggregation,
funnels, a merge.

```json
{"grid": {"colWidth": 280, "rowHeight": 110},
 "elements": [
  {"id": "tests", "kind": "box", "text": "go test ./...", "col": 0, "row": 0},
  {"id": "lint", "kind": "box", "text": "validators", "col": 0, "row": 1},
  {"id": "review", "kind": "box", "role": "ai", "text": "bar-raiser verdict", "col": 0, "row": 2},
  {"id": "merge", "kind": "ellipse", "role": "end", "text": "merge to main", "col": 1, "row": 1}],
 "edges": [{"from": "tests", "to": "merge"}, {"from": "lint", "to": "merge"}, {"from": "review", "to": "merge"}]}
```

## Timeline (a sequence of events)

A line with dots, not a row of boxes. Use the real event names.
`direction: right` reads as time; `down` reads as a log.

```json
{"elements": [
  {"id": "events", "kind": "timeline", "direction": "right", "x": 0, "y": 0, "gap": 150,
   "steps": ["RUN_STARTED", "STATE_DELTA", "TEXT_MESSAGE", "RUN_FINISHED"]}]}
```

## Cycle (a loop that returns)

Boxes in a ring; the return arrow gets a `via` waypoint so it routes
around instead of through. Feedback loops, review rounds, retries.

```json
{"grid": {"colWidth": 260, "rowHeight": 140},
 "elements": [
  {"id": "author", "kind": "box", "role": "ai", "text": "investigator\nauthors", "col": 0, "row": 0},
  {"id": "review", "kind": "box", "role": "ai", "text": "bar-raiser\nreviews", "col": 1, "row": 0},
  {"id": "verdict", "kind": "diamond", "role": "decision", "text": "verdict?", "col": 1, "row": 1},
  {"id": "done", "kind": "ellipse", "role": "end", "text": "approved", "col": 2, "row": 1}],
 "edges": [{"from": "author", "to": "review"}, {"from": "review", "to": "verdict"},
           {"from": "verdict", "to": "done", "label": "approve"},
           {"from": "verdict", "to": "author", "label": "revise", "via": [[130, 210]]}]}
```

## Assembly line (a transformation)

Show the input and the output as they actually look, with the process
between them. Real payloads teach; "Input" and "Output" do not.

```json
{"grid": {"colWidth": 320, "rowHeight": 160},
 "elements": [
  {"id": "in", "kind": "code", "lang": "json", "text": "{\"id\": 200, \"grade\": \"good\"}", "col": 0, "row": 0},
  {"id": "sm2", "kind": "box", "text": "SM-2 update\nease, interval", "col": 1, "row": 0},
  {"id": "out", "kind": "code", "lang": "json",
   "text": "{\"id\": 200, \"due\": \"2026-09-24\",\n \"ease\": 2.6, \"interval\": 6}", "col": 2, "row": 0}],
 "edges": [{"from": "in", "to": "sm2", "label": "log"}, {"from": "sm2", "to": "out"}]}
```

## Side by side (a comparison)

Two sections with mirrored contents, so the eye compares row by row.
Roles carry the verdict: `error` on the left, `end` on the right.

```json
{"grid": {"colWidth": 340, "rowHeight": 100},
 "elements": [
  {"id": "a1", "kind": "box", "text": "hand-written JSON", "col": 0, "row": 0},
  {"id": "a2", "kind": "box", "role": "error", "text": "coordinates by hand", "col": 0, "row": 1},
  {"id": "b1", "kind": "box", "text": "spec + compiler", "col": 1, "row": 0},
  {"id": "b2", "kind": "box", "role": "end", "text": "coordinates computed", "col": 1, "row": 1},
  {"id": "left", "kind": "section", "label": "reference skill", "contains": ["a1", "a2"]},
  {"id": "right", "kind": "section", "label": "this skill", "contains": ["b1", "b2"]}],
 "edges": [{"from": "a1", "to": "a2"}, {"from": "b1", "to": "b2"}]}
```

## Lanes (phases or owners)

Sections as rooms: each names who or what owns the elements inside.
Arrows crossing a section border are the handoffs, which is the point
of the diagram.

```json
{"grid": {"colWidth": 300, "rowHeight": 130},
 "elements": [
  {"id": "user", "kind": "ellipse", "role": "external", "text": "comment:\nsolved 200 good", "col": 0, "row": 0},
  {"id": "action", "kind": "box", "text": "log-solve.yml", "col": 1, "row": 0},
  {"id": "sr", "kind": "box", "text": "sr.py log 200 good", "col": 2, "row": 1},
  {"id": "state", "kind": "box", "role": "data", "text": "state.json", "col": 3, "row": 1},
  {"id": "gh", "kind": "section", "label": "GitHub Actions", "contains": ["user", "action"]},
  {"id": "repo", "kind": "section", "label": "repository", "contains": ["sr", "state"]}],
 "edges": [{"from": "user", "to": "action"}, {"from": "action", "to": "sr"},
           {"from": "sr", "to": "state", "label": "commit"}]}
```

## Tree (a hierarchy)

Lines and free text, no boxes. File systems, org charts, taxonomies.

```json
{"elements": [
  {"id": "root", "kind": "text", "level": "subtitle", "text": "agent-team-workspace/", "x": 0, "y": 0},
  {"id": "trunk", "kind": "line", "points": [[12, 30], [12, 150]]},
  {"id": "b1", "kind": "line", "points": [[12, 60], [40, 60]]},
  {"id": "n1", "kind": "text", "text": "requirements/", "x": 48, "y": 50},
  {"id": "b2", "kind": "line", "points": [[12, 100], [40, 100]]},
  {"id": "n2", "kind": "text", "text": "design-docs/", "x": 48, "y": 90},
  {"id": "b3", "kind": "line", "points": [[12, 140], [40, 140]]},
  {"id": "n3", "kind": "text", "text": "protocols/", "x": 48, "y": 130}]}
```

## Evidence beside the thing it proves

A `code` block next to the shape it explains, joined by a dashed edge
with no arrowhead. This is how a technical diagram teaches.

```json
{"grid": {"colWidth": 440, "rowHeight": 150},
 "elements": [
  {"id": "hook", "kind": "box", "text": "PostToolUse hook", "col": 0, "row": 0},
  {"id": "payload", "kind": "code", "lang": "json",
   "text": "{\"tool_input\": {\"file_path\": \".claude/agents/x.md\"}}", "col": 0, "row": 1},
  {"id": "validator", "kind": "box", "text": "validate-definitions.py", "col": 1, "row": 0}],
 "edges": [{"from": "payload", "to": "hook", "dashed": true, "end": null},
           {"from": "hook", "to": "validator", "label": "exit 2 on failure"}]}
```

## Choosing between them

| You are showing | Reach for |
|---|---|
| where one thing goes | fan-out |
| what feeds one result | convergence |
| what happens, in order | timeline |
| what repeats until a condition | cycle |
| what goes in and what comes out | assembly line |
| two options, honestly | side by side |
| who owns which step | lanes |
| what contains what | tree |

Mix them. A real architecture diagram is often lanes containing an
assembly line, with a timeline underneath and evidence blocks beside
two of the boxes. Six identical boxes in a grid is none of these.
