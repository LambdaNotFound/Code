# Async shapes

Load before writing async code. This file is about **choosing the
concurrency shape**; the pitfall checklist for async code already
written lives in `.claude/skills/pr-review/references/rust.md`, and is
not repeated here.

Pick the shape first. Almost every painful async refactor is a shape
chosen by accident.

## The four shapes

**1. Task per unit of work** — `tokio::spawn`, collected in a
`JoinSet`. Use when units are independent and each can own what it
needs. Requires `'static + Send`, which is what forces the `Arc`
clones. Shutdown is `JoinSet::abort_all` or a `CancellationToken`;
`JoinSet` also surfaces panics, which a bare `spawn` with a dropped
handle silently swallows.

**2. Select loop** — one task owns the state and loops on
`tokio::select!` over its inputs. Use when several event sources
mutate one piece of state. The state is *owned by the loop*, so there
is no lock at all: mutual exclusion comes from there being one owner.
This is usually the right answer when a Go instinct reaches for a
mutex.

**3. Channel actor** — state lives in a task; other tasks send
command messages over an `mpsc` and get answers back over a
`oneshot`. Use when the state needs to be reachable from many places
and shape 2's inputs are not a fixed set. It removes `Arc<Mutex<T>>`
entirely, and the mailbox gives you natural backpressure. Cost: a
message type, and a reply channel per request.

**4. Shared state** — `Arc<Mutex<T>>`. Legitimate when the critical
section is tiny, contention is low, and the alternative shapes add
more machinery than they remove. It is a conclusion, not a starting
point: say what you tried first.

## Cancellation is by drop

Dropping a future cancels it at its last suspension point. Nothing
runs afterwards — no `defer`, no `finally` — except `Drop` impls of
values still alive, and `Drop` cannot be async. Consequences to design
around:

- **Cleanup that must happen needs an explicit `shutdown().await`**,
  not a `Drop` impl. A `Drop` that wants to await is a design error.
- **`select!` drops the losing branches mid-poll.** If such a future
  had consumed input into an internal buffer, that input is gone.
  `AsyncReadExt::read` is cancel-safe; `read_exact` is not. Before
  putting a future in a `select!`, state whether dropping it
  mid-flight loses data — and if it does, move it into its own task
  with a channel instead.
- **A loop with `select!` re-creates its futures every iteration.**
  Anything that must survive across iterations belongs outside the
  loop, pinned, or in a task.

## Backpressure

Bounded channels are the design decision; unbounded ones defer it
until memory runs out. `mpsc::channel(n)` makes the producer wait,
which is what you want. Reach for `try_send` when the correct
behavior on a full queue is to reject rather than to wait — that is
the difference between a leaky bucket and a blocking queue, and
`tokio_examples/leaky_bucket.rs` in this crate is the worked example.

## Runtime and blocking

The multi-threaded runtime is the default; `current_thread` is for
tests and single-core cases and makes `!Send` futures usable. Either
way, a blocking call inside async work stalls a worker thread and
every task on it. File IO, CPU-bound loops, `std::thread::sleep`, and
synchronous drivers go to `spawn_blocking`. Name the executor
assumption whenever the code depends on one.

## Testing async

Use `#[tokio::test(start_paused = true)]`, which this crate already
enables through the `test-util` dev-dependency. Simulated time
fast-forwards through `sleep` and `interval`, so a test of a
ten-second timeout runs instantly and deterministically. A test that
actually sleeps is both slow and flaky — treat one as a defect.

Assert on elapsed simulated time (`Instant::now()` before and after)
to prove timing behavior, the way `src/tests.rs` does.
