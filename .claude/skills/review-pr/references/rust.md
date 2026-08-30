# Rust review reference

Load when the diff contains Rust. Ordered by how often the item is
the actual bug, not by language-tour order. Each item names what to
look for, why it breaks, and the fix.

## Ownership and borrowing

- **`.clone()` added to satisfy the borrow checker.** The clone is
  the symptom; the design is the finding. Ask what shape avoids it:
  split borrows into disjoint fields, index instead of holding a
  reference across a mutation, `std::mem::take` to move out and put
  back, or end the borrow earlier with a scope. A clone in a
  request path or a loop is a performance finding too.
- **New `Rc<RefCell<T>>` or `Arc<Mutex<T>>` in a fresh design.**
  The borrow checker was reporting a real aliasing problem;
  `RefCell` does not solve it, it moves the failure to runtime,
  where `borrow_mut` panics on overlap — usually via a reentrant
  callback or a `Drop`. Ask whether an arena with indices, a tree
  owned by its parent, or passing `&mut` down instead of storing it
  removes the shared mutability entirely.
- **`&String`, `&Vec<T>`, `&Box<T>` in a signature.** Take `&str`,
  `&[T]`, `&T`. The concrete types force callers to allocate.
- **A struct that grew a `'a` parameter.** Lifetimes are viral —
  every holder inherits it. Check whether owning the data or an
  `Arc` is what the call sites actually want.
- **A `'static` bound on a callback or spawned closure** that forces
  the caller to clone or leak. Sometimes required; often it is a
  missing scoped-thread or lifetime-generic formulation.

## `unsafe`

- **Every `unsafe` block needs a `// SAFETY:` comment** naming the
  invariant it relies on and why it holds *at this call site*. An
  `unsafe` block with no safety comment is blocking, even if the
  code is right — the next editor cannot know what to preserve.
- **`unsafe fn` shifts the obligation to the caller.** The doc
  comment needs a `# Safety` section stating it, and every call site
  needs its own safety comment proving it.
- UB to check for specifically: two live `&mut` to the same place
  (including one derived from a raw pointer); `transmute` between
  types with different validity invariants (`u8` to `bool`, any
  bit pattern to an enum, extending a lifetime); `slice::from_raw_parts`
  with a null, misaligned, or dangling pointer — non-null and
  aligned is required *even for length 0*; reading `MaybeUninit`
  before it is initialized; `&mut` to a `static mut`; pointer
  arithmetic beyond one-past-the-end.
- **`unsafe impl Send`/`Sync`** is a proof obligation, not a
  suppression. It needs a written argument about what makes
  concurrent access sound. "Should be fine" is blocking.
- If the crate carries meaningful `unsafe`, check CI runs `miri`.

## Async

The highest-yield section in async Rust review. Read every `.await`
and ask what is held across it, and what happens if the future is
dropped there.

- **A `std::sync::MutexGuard` held across `.await`.** Makes the
  future `!Send` (so it cannot be `tokio::spawn`ed) and can deadlock
  when the same task re-enters. The fix is usually to scope the
  guard so it drops before the await, not to swap in
  `tokio::sync::Mutex` — the async mutex is slower and only needed
  when the lock must genuinely be held across a suspension point.
- **Blocking work inside an async fn** — file IO, CPU-bound loops,
  `std::thread::sleep`, a synchronous DB or crypto call. It starves
  the runtime worker thread and stalls every other task on it. Use
  `spawn_blocking` or a dedicated pool.
- **Cancellation safety in `select!`.** The branch that does not win
  is dropped *mid-poll*. Any state that future had accumulated
  internally is lost. `AsyncReadExt::read` is cancel-safe;
  `read_exact` is not — it can consume bytes and then be dropped,
  losing them. Hand-written state machines usually are not.
  Review rule: for every future in a `select!`, state whether
  dropping it mid-flight loses data. If it does, it belongs outside
  the `select!`, in a task with its own channel.
- **`tokio::spawn` whose `JoinHandle` is dropped.** A panic in that
  task disappears silently and shutdown does not wait for it. Store
  the handle, or use a `JoinSet`.
- **`Drop` cannot be async and does not always run** — not on
  `mem::forget`, `process::exit`, an aborted task, or an `Rc`/`Arc`
  cycle. Cleanup that must happen needs an explicit
  `close()`/`shutdown()` the caller awaits.
- An `async fn` whose result is never `.await`ed does nothing at
  all. `#[must_use]` on futures-returning functions catches it.

## Errors and panics

- `unwrap()`/`expect()` in library code panics in someone else's
  process. Fine in `main`, tests, and after a checked invariant with
  an `expect` message that states the invariant — "index is in
  bounds, checked above", not "failed".
- `?` that erases context. Add `map_err` or `.context()` at the
  boundary where the reader needs to know which file, which key.
- Libraries return a concrete error enum (`thiserror`);
  applications may use `anyhow`. `anyhow::Error` in a public library
  signature is a finding — it denies callers the ability to match.
- A `panic!` unwinding across an `extern "C"` boundary is UB. Needs
  `catch_unwind` at the boundary.

## Arithmetic and casts

- Overflow panics in debug and wraps in release. Silent divergence
  between profiles is a bug factory. Make intent explicit:
  `checked_*`, `saturating_*`, `wrapping_*`.
- `as` truncates silently (`u64 as u32`) and saturates float-to-int.
  Prefer `TryFrom` with a handled error anywhere the value is not
  provably in range.
- `usize` is 32-bit on some targets. Casting a length or an offset
  through `u32`/`i32` is a portability finding.

## Traits and API evolution

- A `_ =>` catch-all in a match over a crate-internal enum defeats
  the compiler's best feature: adding a variant should break the
  build at every site that must handle it. Flag internal catch-alls.
- Public enums and structs that may grow need `#[non_exhaustive]`;
  adding a variant or field is otherwise a breaking change.
- Adding a trait method without a default body breaks every external
  implementor.
- `#[must_use]` on builders, on guards, and on any type whose drop
  means work was silently discarded.
- Feature flags must be additive — no feature may remove or change
  an item. Check that `--no-default-features` and `--all-features`
  both build.

## Performance patterns worth flagging

- `collect()` into a `Vec` that is iterated once and dropped.
- `.iter().cloned()` where `.copied()` applies — it documents that
  the element is trivially copyable.
- String building in a loop without `String::with_capacity`.
- `Box<dyn Trait>` in a hot loop where a generic would monomorphize.
- A large struct cloned per request rather than passed by reference
  or wrapped in `Arc`.

## Tests and tooling

- `#[should_panic]` without `expected = "..."` passes on *any*
  panic, including an unrelated one from a typo. Always pin the
  message.
- Check CI for: `clippy -- -D warnings`, `cargo fmt --check`, an
  MSRV that is pinned and actually built, `miri` where `unsafe`
  lives, and `cargo audit`/`cargo deny` when the diff adds
  dependencies.
- Public API examples belong in doc tests, which are compiled and
  run — unlike a fenced block in a README.
