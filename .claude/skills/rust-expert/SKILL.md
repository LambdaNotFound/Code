---
description: 'Write Rust with the user, interactively — design the ownership and the type model before typing, then implement compiler-first and verify with fmt, clippy, and tests. Carries the two things the Rust agents do not write down: a borrow-checker error to design-fix triage, and a Go-to-Rust idiom mapping for an engineer whose primary language is Go. Use when writing or refactoring Rust, choosing an ownership or API shape, fixing a borrow-check or lifetime error, designing async and Tokio concurrency, or translating Go into Rust. Not for reviewing Rust someone else wrote (use pr-review, whose .claude/skills/pr-review/references/rust.md is the review checklist). Not for a self-contained Rust task you want dispatched to an isolated subagent (use rust-pro, or rust-engineer for Go-to-Rust work). Not for implementing a whole brief as a reviewed PR (use pr-loop).'
argument-hint: <what to build, fix, or translate> [crate or path]
---

You write Rust with the user, in the main session, where the
conversation is part of the work. That is what this skill is for and
what a subagent cannot do: ownership decisions are design decisions,
and the good ones get made out loud, with the user, before any code
exists.

Everything below is a standing instruction. This file is read once
and stays in context, so these rules apply on every turn of the task,
not only the first.

## Arguments

$ARGUMENTS

## 1. Read the crate before writing a line

Conventions in the crate bind, even where you would choose
differently: error types, module layout, feature gating, lifetime
style, naming. Read `Cargo.toml` for the edition and the dependency
set before assuming any API exists. If the crate contradicts itself,
say so and ask which convention holds rather than picking one.

Minimal diffs. Change what the task requires and nothing adjacent.

## 2. Design the ownership before you type

Most Rust pain is a data-model problem arriving disguised as a
borrow-check error. Fighting the borrow checker means the design is
wrong somewhere upstream of the line that failed, and the fix belongs
there, not at the error.

Before writing, answer in a few lines and put it in front of the user:

- **Who owns each piece of data**, and for how long.
- **What borrows it**, and whether that borrow can end before the
  next mutation.
- **What is shared**, and what forced sharing — `Arc<Mutex<T>>` is a
  conclusion you arrive at, never a starting point. Say what you
  tried first. Passing ownership down a channel often removes the
  need entirely.
- **The rejected alternative**, in one clause. This crate documents
  those in the code (see the conventions section below), so you will
  need it anyway.

When a lifetime parameter starts spreading through a third type, stop
and reconsider: that is the design asking for an owned type, an
`Arc`, or an arena with indices — not for more annotations.

## 3. Make the illegal states unrepresentable

This is where Rust pays for itself, and where a Go-shaped design
leaves the most value on the table.

- A struct whose fields are only valid in certain combinations wants
  to be an enum whose variants carry exactly the data that variant
  needs.
- Validate once at the boundary and return a type that proves it —
  `Email(String)`, `NonZeroU32`, `Verified<T>` — rather than
  re-checking an invariant at every use.
- Newtype anything with units or provenance. `u64` is not a
  `UserId`, and the compiler cannot help you until you say so.
- A `bool` parameter at a call site is unreadable; a two-variant enum
  reads itself.
- Reach for typestate when a value moves through phases and calling
  the wrong method should not compile.

## 4. Choose the error strategy up front

`thiserror` for libraries, `anyhow` for binaries and tests. The
boundary between them is a design decision, not a default: a library
that returns `anyhow::Error` denies its callers the ability to match.

Error enums are public API — `#[non_exhaustive]` unless exhaustive
matching by callers is the point. No `unwrap` or `expect` outside
tests and statically provable cases, and where `expect` is genuinely
right the message states the invariant that makes it so ("index
checked above", never "failed"). Panics are for broken invariants in
your own code, never for bad input.

## 5. Implement compiler-first

Write the types and signatures, then let `cargo check` drive. It is
faster than reasoning about whether a borrow works, and rustc's
`help:` and `note:` lines usually contain the actual fix.

When the borrow checker rejects something, do not reach for `.clone()`
or `Arc` to make the error go away. Read
[references/borrow-checker.md](references/borrow-checker.md), which
maps the common error codes to the design change each one is really
asking for.

For async and concurrency, decide the *shape* first — task per unit
of work, a select loop, or a channel actor — before writing any of
it. [references/async-shapes.md](references/async-shapes.md) covers
choosing between them and the cancellation consequences of each.

When translating from Go, or when a Go instinct is driving the
design, read [references/go-to-rust.md](references/go-to-rust.md).
Naming the mapping explicitly is the highest-value thing you can do
for this user, and this crate's own examples do it in every file.

## 6. `unsafe` is a last resort with a written argument

Never use `unsafe` to work around the borrow checker; that is a design
signal, not a lifetime problem. Where it is genuinely justified, every
block carries a `// SAFETY:` comment naming the invariant and who
upholds it, every unsafe function documents its preconditions under
`# Safety`, and Miri runs over a test that exercises the path. If you
cannot construct that test, say so rather than shipping it untested.

## 7. Verify before claiming done

```
cargo fmt --check
cargo clippy --all-targets -- -D warnings
cargo test
```

Report what failed, with the output. Never claim completion on a
failing build, and never describe a test as passing that you did not
run. Benchmark before optimizing: zero-copy, SIMD, const generics,
custom allocators, and `no_std` each need a profile behind them, and
a performance claim with no measurement is labeled **inferred** or
not made.

Before handing work back, self-review the diff against
`.claude/skills/pr-review/references/rust.md`, which is this repo's
Rust review checklist. Catching your own findings costs one turn;
catching them in review costs a round trip.

## Where the other Rust definitions sit

The standards in sections 4, 6, and 7 — the `thiserror`/`anyhow`
split, the `unwrap` rule, the `unsafe` contract, the verification
triad — are deliberately the same ones `.claude/agents/rust-pro.md`
holds, because they are Rust community norms rather than this repo's
policy. If you change one, change both. What lives only here is the
method, the three references, and the repository conventions below.

## This repository's Rust conventions

`rust/` is edition 2024, a binary crate (`tokio_hello_world`) with a
committed `Cargo.lock`, using `anyhow`, `tokio` with `full`, and
`tokio-util`.

- **Every example file opens with `//! Pattern: <what it shows>`
  followed by `//! Go comparison:` explaining how a Go engineer's
  instinct maps — or fails to map — onto the Rust shape.** All 24
  example files do this without exception. Match it; it is the point
  of the crate.
- Where a design has a real rival, document it inline as
  `Rejected alternative: ...` with the reason, as
  `tokio_examples/leaky_bucket.rs` does.
- Examples live in `tokio_examples/` and `basics_examples/`, not
  `examples/`, so Cargo does not auto-discover them. **Every new
  example needs an explicit `[[example]]` block in `Cargo.toml`** with
  `name` and `path`, or it will never build.
- Async tests use `#[tokio::test(start_paused = true)]` with the
  `test-util` dev-dependency, so simulated time fast-forwards instead
  of actually sleeping. A test that really sleeps is a finding.

## What to report

The ownership decision and the alternative you rejected, the files
you changed, the verification commands you ran with their results,
and anything you left undone. Not a walkthrough of code the user can
read.
