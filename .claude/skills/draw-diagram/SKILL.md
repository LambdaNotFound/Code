---
description: 'Draw a diagram from a description and deliver it as an .excalidraw file — architecture, flow, protocol, pipeline, and concept diagrams whose shape carries the argument rather than a grid of labelled boxes. You write a small JSON spec (shapes coloured by meaning, edges by id, sections, timelines, code blocks as evidence); a stdlib compiler turns it into Excalidraw JSON, checks the layout for overlaps, clipped text, and arrows through shapes, and renders a PNG preview to look at before delivering. Use when the user asks to draw, diagram, sketch, map, or visualise a system, flow, pipeline, protocol, or idea, or wants an Excalidraw file. Not for charts or plots of data (use dataviz). Not for UI mockups, posters, or page layouts (use design).'
argument-hint: '<what to diagram> [out: path] [style: sketch]'
---

You draw diagrams that argue. A diagram is not formatted text: its
shape carries the claim, so one-to-many looks like a fan, a sequence
looks like a line with events on it, and a boundary looks like a
boundary. Given a description, you decide what the reader must see,
map each concept to a form that mirrors its behaviour, write a small
spec, compile it, look at the result, and fix what you see. The
deliverable is an `.excalidraw` file the user opens and edits in
Excalidraw.

You never write Excalidraw JSON by hand. The compiler at
`.claude/skills/draw-diagram/scripts/build_excalidraw.py` owns every
coordinate, text measurement, binding, and seed; the spec is what you
author and what you fix. This runs in the main session, because what
to draw is a conversation and the render loop needs your eyes.

## Arguments

$ARGUMENTS

## 1. Decide what the reader must see

Two depths, and the choice comes first:

- **Conceptual**: a mental model, a philosophy, the shape of an idea.
  Abstract labels are right; concrete detail would be noise.
- **Technical**: a real system, protocol, pipeline, or codebase. Then
  the diagram must show what things actually look like: the real
  event names, the actual payload, the function that gets called, the
  file that holds the state. "Input, Process, Output" teaches nothing.

For a technical diagram, research before drawing. Read the code, the
spec, the protocol document, and use the names they use; a box
labelled "Database" where the code says `state.json` is a wrong
diagram. Put evidence on the canvas as `code` blocks: a payload, a
command, a three-line snippet. Every name on the canvas should trace
to a file or document you read.

Then tell the user, in a few lines before building anything: what the
diagram argues, who reads it, and what they need to see to believe
it. If the request is ambiguous between two diagrams, ask which. One
canvas, one argument.

## 2. Map each concept to a form

| The concept... | Draw it as |
|---|---|
| Spawns many outputs | fan-out: one source, arrows to a stacked column |
| Merges many inputs | convergence: a stacked column, arrows into one target |
| Is a sequence of events | `timeline`: a line with dots and labels, never a row of boxes |
| Loops or iterates | cycle: a ring, the return arrow routed with `via` |
| Transforms something | assembly line: `code` in, box, `code` out |
| Compares two options | side by side: two `section`s with mirrored contents |
| Splits into phases or owners | `section`s as lanes |
| Branches on a condition | `diamond`, with labelled edges out |
| Starts or ends the flow | `ellipse` |
| Is a step or component | `box`, when arrows need to land on it |
| Is a label, note, or caption | free `text`, no container |

Two rules of taste the checker cannot enforce. Vary the forms: six
identical boxes in a grid is a list wearing a costume. Default to free
text: a box earns its border by being a thing arrows connect to or a
thing that is distinct in the system; titles, annotations, and
details are text at a smaller level.
[references/visual-patterns.md](references/visual-patterns.md) has a
buildable spec for each form above; copy the shape and replace the
words.

## 3. Write the spec

The format is [references/spec-format.md](references/spec-format.md):
elements placed on a grid (`col`, `row`) or at absolute `x`, `y`,
coloured by what they mean (`role`), connected by `edges` that name
ids. Read it the first time. `--example` prints a complete starter.

Where it lives: `diagrams/<slug>/<slug>.spec.json`, with a short
kebab-case slug from the topic, unless the user gives `out:`. The
compiler writes `<slug>.excalidraw`, `<slug>.svg`, and `<slug>.png`
beside it. The spec and the `.excalidraw` are the durable pair; the
previews are for looking, and the user decides whether to commit them.

Rules that hold in every spec:

- Colour by `role`, never by hex. The roles and what each is for are
  in `.claude/skills/draw-diagram/scripts/palette.json`, and
  `--roles` prints them. Changing the palette changes every diagram,
  which is the point of having one.
- `text` holds only readable words. Break lines with `\n`; the
  compiler does not wrap.
- Descriptive ids (`webhook`, `validate_sig`, `retry_queue`): edges
  and sections name them, and the checker quotes them.
- Leave `w` and `h` out unless a shape must be bigger than its text.
  The compiler sizes containers so text never clips.
- One grid, chosen once. Widen `colWidth` rather than nudging
  elements one at a time.

## 4. Build and read the checks

```
python3 .claude/skills/draw-diagram/scripts/build_excalidraw.py diagrams/<slug>/<slug>.spec.json
```

Errors mean the spec cannot build (an unknown id or role, a section
containing itself) and nothing is written until they are fixed.
Warnings are layout defects a reader would notice: two elements
overlapping, text wider than a box you forced, an arrow passing
through a shape that is neither its source nor its target, an element
sitting inside a section it is not listed in. Fix every warning in
the spec, or say in the report why one stands. Notes are judgement
calls left to you: a shape with no edges and no section, or a canvas
that is nothing but boxes.

## 5. Look, then fix in the spec

Read the PNG. The checker sees geometry; you see the argument. Ask,
in this order:

1. Does the structure say the thing without the words? Cover the
   labels in your mind: does one-to-many still look like one-to-many?
2. Does the eye travel the way the story goes: left to right or top
   to bottom for sequences, outward for a fan?
3. Is the hero the biggest thing with the most space around it?
4. Are the evidence blocks readable and beside what they prove?
5. Is anything ambiguous: a label floating between two shapes, an
   arrow that could belong to either neighbour?

Fix in the spec and rebuild; never edit the `.excalidraw`, which the
next build overwrites. Two or three rounds is normal. Stop when the
checks are clean and you would show it without a caveat. If no PNG
was written because the machine has no Chromium, the SVG is the same
picture; look at that, or say plainly that you checked geometry only.

## 6. Deliver

Show the PNG and give the paths. Report: what the diagram argues, in
one sentence; the sources every name came from; warnings left
standing, with the reason; and what you assumed where the description
was silent. Not a tour of the elements.

If the user then edits the `.excalidraw` in Excalidraw, their file is
the source from then on. Rebuilding from the spec would discard their
edits, so do not rebuild unless they ask, and say so when they do.

`style: sketch` switches to the hand-drawn font and rough strokes for
brainstorm-grade diagrams; the default is clean. The compiler's tests
are `python3 .claude/skills/draw-diagram/scripts/test_build_excalidraw.py`;
run them after changing the compiler or the palette.
