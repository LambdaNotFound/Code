# Tokio 101 for Go Engineers

This is a walkthrough of `rust/tokio_examples/*.rs`, an existing Cargo binary crate
(`tokio_hello_world`) in this repo. It's written for someone with deep Go
experience and no prior Rust/Tokio exposure. Every section leans on a Go
idiom you already know and calls out where the Rust version diverges. Run
any example with `cargo run --example <name>` from inside `rust/`.

## Go → Rust/Tokio idiom map

| Go | Rust / Tokio | Note |
|---|---|---|
| `go f()` (goroutine) | `tokio::spawn(async { ... })` (task) | A goroutine panic crashes the process; a task panic is caught and returned as `Err(JoinError)` from the `JoinHandle`. |
| `chan T` (buffered/unbuffered) | `mpsc::channel` / `mpsc::unbounded_channel` | Many senders, one receiver; bounded capacity blocks the sender like a buffered Go chan. |
| `done := make(chan T, 1)` idiom | `oneshot::channel` | The sender is consumed by `send`, so the type system (not convention) prevents a second send. |
| mutex + `sync.Cond` for "latest value" | `watch::channel` | No Go stdlib equivalent; only the newest value is observable, not a backlog. |
| hand-rolled fan-out to N channels | `broadcast::channel` | No Go stdlib equivalent; every subscriber gets every message. |
| `select { case ...: }` | `tokio::select! { ... }` | Non-chosen branches are dropped/cancelled, not left running; branches support `if` guards. |
| `context.Context` / `ctx.Done()` | `tokio_util::sync::CancellationToken` | Both are cooperative. `child_token()` mirrors `context.WithCancel(parent)`. |
| `sync.WaitGroup` | `tokio::join!` / `JoinSet` | `join!` is for a fixed set and hands back each result; `JoinSet` is for a dynamic set. |
| `errgroup.Group` | `tokio::try_join!` | First `Err` short-circuits the rest. |
| `sync.Mutex` / `sync.RWMutex` | `Arc<Mutex<_>>` / `Arc<RwLock<_>>` | Two Mutex flavors exist (`std` and `tokio::sync`); see the shared-state section. |
| blocking I/O in a goroutine (no extra step) | `tokio::task::spawn_blocking` | Tokio's async workers aren't preempted, so blocking calls need explicit offload. |
| `(T, error)` + `if err != nil { return fmt.Errorf("...: %w", err) }` | `anyhow::Result<T>` + `?` + `.context(...)` | `.context()` plays the role of `%w` wrapping. |

## The runtime: `#[tokio::main]` vs a manual `Runtime`

**Concept.** Rust's `async fn main` produces a future that nothing runs on
its own; you need an executor. `#[tokio::main]` is a macro that desugars to
building a `Runtime` and calling `.block_on(main_body())`. You can also build
one yourself with `Runtime::new()` (multi-threaded, one worker per CPU) or
`Builder::new_current_thread()` (single-threaded).

**Go comparison.** Go has no equivalent decision to make. The runtime
always multiplexes goroutines across OS threads for you, tuned by
`GOMAXPROCS`. Tokio makes the scheduler an explicit, constructible object.

**Excerpt** (`rust/tokio_examples/runtime_flavors.rs`):

```rust
let multi = Runtime::new()?;
multi.block_on(report("multi-threaded runtime"));

let current = Builder::new_current_thread().enable_all().build()?;
current.block_on(report("current-thread runtime"));
```

**Pitfall.** Forgetting `.enable_all()` on a manually built `Builder`. Without
it, the time and I/O drivers aren't installed, and the first call to
`tokio::time::sleep` or any network type panics with a "there is no reactor
running" error at runtime, not a compile error. `Runtime::new()` and
`#[tokio::main]` both enable everything for you, which is why this only
bites people who start constructing runtimes by hand.

## Spawning tasks: `tokio::spawn` and `JoinHandle`

**Concept.** `tokio::spawn` schedules a future onto the runtime immediately
(it doesn't wait to be polled) and returns a `JoinHandle<T>`, itself a
future you can `.await` to get the task's return value.

**Go comparison.** `tokio::spawn(async { ... })` is `go func() { ... }()`.
`JoinHandle` is the single-result channel you'd otherwise build by hand to
get a value out of a goroutine. The real divergence is panic behavior: an
unrecovered goroutine panic kills the whole Go process; a panic inside a
spawned Tokio task is caught by the runtime and turned into
`Err(JoinError)` when you await the handle: siblings and the process
survive.

**Excerpt** (`rust/tokio_examples/spawn_and_join.rs`):

```rust
let ok = tokio::spawn(async { 2 + 2 });
println!("ok task result: {}", ok.await?);

let panics = tokio::spawn(async { panic!("boom"); });
match panics.await {
    Ok(()) => println!("unreachable"),
    Err(join_err) => println!("panicking task reported as: {join_err}"),
}
```

**Pitfall.** Dropping a `JoinHandle` does not cancel the task. This is unlike
dropping a bare (unspawned) future, which does nothing because it was never
polled. A spawned task is already running independently of its handle, so a
dropped handle just means you can no longer observe its result or panic.
Go engineers sometimes expect handle-drop to behave like a `context` cancel;
it doesn't.

## `.await` and cooperative scheduling

**Concept.** A Rust future is a lazy state machine that only makes progress
when polled, and `.await` is the only point where a task can hand control
back to the scheduler.

**Go comparison.** Since Go 1.14 the runtime can preempt a goroutine that's
been running too long even with no function calls or channel ops, so a tight
`for {}` loop won't permanently starve anything else. Tokio has no such
preemption. A task that never awaits will run to completion on its worker
thread before anything else on that thread gets a turn.

**Excerpt** (`rust/tokio_examples/cooperative_scheduling.rs`):

```rust
let hog = tokio::spawn(async {
    let mut sum: u64 = 0;
    for i in 0..200_000_000u64 {
        sum = sum.wrapping_add(i);
        if i % 50_000_000 == 0 {
            task::yield_now().await;
        }
    }
    sum
});
```

**Pitfall.** A CPU-bound loop with no `.await` inside a task starves every
other task on that worker thread. On a `current_thread` runtime, that means
the entire program. The fix is either periodic `tokio::task::yield_now().await`
calls (shown above) or moving the work to `spawn_blocking` (see below).

## `sleep` and `timeout`

**Concept.** `tokio::time::sleep` suspends a task for a duration without
blocking a thread; `tokio::time::timeout` races a future against a deadline
and returns `Err` if the deadline wins.

**Go comparison.** `sleep` is `time.Sleep`. `timeout` replaces the
`context.WithTimeout` + `select` on `ctx.Done()` dance: you wrap the future
once instead of threading a context through every call site and checking
`Done()` yourself.

**Excerpt** (`rust/tokio_examples/sleep_and_timeout.rs`):

```rust
match timeout(Duration::from_millis(10), slow_operation(50)).await {
    Ok(value) => println!("unreachable: {value}"),
    Err(_) => println!("slow_operation timed out, as expected"),
}
```

**Pitfall.** `tokio::time::sleep` is not `std::thread::sleep`. Calling the
`std` version inside an `async fn` compiles fine but blocks the whole worker
thread instead of just yielding the task, the same class of bug as the
cooperative-scheduling pitfall above, just via a stdlib call instead of a
spin loop.

## Channels: `mpsc`, `oneshot`, `watch`, `broadcast`

**Concept.** Tokio splits Go's one `chan T` into four purpose-built types,
each with a different ownership shape; see the mapping table above.

**Go comparison.** Go channels are inherently multi-producer/multi-consumer;
Tokio makes you pick the shape that matches your actual usage, and the
compiler enforces it (e.g. a `oneshot::Sender` is consumed on `send`, so a
second send is a compile error, not a runtime bug).

**Excerpt** (`rust/tokio_examples/channels.rs`, the `watch` case, "only the
latest value matters"):

```rust
while rx.changed().await.is_ok() {
    let value = *rx.borrow();
    println!("watch: latest value is {value}");
    if value == 3 {
        break;
    }
}
```

**Pitfall (mpsc).** `mpsc::Sender::send` on a bounded channel is `async`. It
returns a future. Forgetting the `.await` (e.g. `tx.send(x);` instead of
`tx.send(x).await?;`) compiles, but the compiler emits an "unused
implementer of `Future` that must be used" warning and the value is never
actually sent, because nothing polled the send future. Go's `ch <- x` has no
equivalent silent-no-op failure mode.

**Pitfall (broadcast).** A `broadcast` channel has a fixed capacity per
subscriber, not just per sender. If a subscriber falls behind (it's slow,
or just isn't polling), the sender doesn't block on it the way a Go
unbuffered channel would block on a slow reader. Instead the oldest
unread messages are dropped, and that subscriber's next `recv()` returns
`Err(Lagged(n))` telling it how many it missed, rather than the messages
themselves. It's easy to treat that as a fatal error and bail out, when the
correct handling is usually to log it and keep receiving. Silently missing
messages, not blocking the sender, is the whole point of the design.

## `tokio::select!`

**Concept.** Race several futures and run the branch for whichever completes
first.

**Go comparison.** Direct analogue of Go's `select` over channel operations,
including a `default:`-like immediate branch. Two things Go doesn't have:
per-branch `if` guards (evaluated before the branch is even considered
ready), and a `biased;` modifier to disable the default random branch
ordering.

**Excerpt** (`rust/tokio_examples/select_macro.rs`):

```rust
tokio::select! {
    Some(msg) = rx.recv() => println!("select: received {msg}"),
    _ = sleep(Duration::from_millis(50)) => println!("select: timed out"),
}
```

**Pitfall.** The futures in branches that *aren't* chosen are dropped, i.e.
cancelled mid-flight. If a losing branch's future had side effects partway
through (e.g. it had already sent part of a multi-step request), those
effects don't roll back automatically. The future you write needs to be
safe to cancel at any `.await` point inside it. Go's `select` has no
equivalent concern, since a goroutine already started keeps running whether
or not its channel op "wins" the select.

## Structured concurrency: `join!`, `try_join!`, `JoinSet`

**Concept.** `join!` waits for a fixed, known set of futures and returns all
their outputs. `try_join!` does the same but short-circuits on the first
`Err`. `JoinSet` manages a dynamic, runtime-determined number of tasks.

**Go comparison.** `join!` is a `sync.WaitGroup` that also hands back
values instead of making you close over shared variables. `try_join!` is an
`errgroup.Group`. `JoinSet` replaces "spawn N goroutines in a loop, fan
results back through a channel"; the channel wiring is built in. One thing
that doesn't carry over: `join_next()` yields results in completion order,
not spawn order, so a Go pattern that relied on a `WaitGroup` plus an
index-ordered slice needs an explicit re-sort here (e.g. by an id you attach
to each task's output) if order matters.

**Excerpt** (`rust/tokio_examples/structured_concurrency.rs`):

```rust
let mut set = JoinSet::new();
for id in 0..3 {
    set.spawn(async move { compute(id).await });
}
let mut results = Vec::new();
while let Some(res) = set.join_next().await {
    results.push(res?);
}
```

**Pitfall.** `join!` waits for *all* of its futures regardless of whether
one fails early. It has no short-circuit behavior, which is exactly why
`try_join!` exists as a separate macro. Separately, dropping a `JoinSet`
aborts every task still in it, unlike dropping a single `JoinHandle`
from a plain `tokio::spawn`, which leaves the task running. That asymmetry
between the two is easy to get backwards coming from Go, where nothing you
do with a `WaitGroup` ever cancels a goroutine.

## Shared state: `Arc<Mutex<_>>` / `Arc<RwLock<_>>`

**Concept.** `Arc` gives shared ownership across tasks/threads; `Mutex`/
`RwLock` guard the data inside it. Tokio ships two Mutex implementations:
`std::sync::Mutex` (blocks the OS thread while waiting) and
`tokio::sync::Mutex` (yields the task while waiting, and can be held across
an `.await`).

**Go comparison.** `Arc<Mutex<T>>` is `sync.Mutex` guarding shared data, made
explicit in the type: the compiler refuses to hand you the data without
going through the lock, where Go trusts convention. `Arc<RwLock<T>>` mirrors
`sync.RWMutex` the same way. The two-Mutex-types split has no Go analogue at
all.

**Excerpt** (`rust/tokio_examples/shared_state.rs`):

```rust
let mut guard = counter
    .lock()
    .map_err(|_| anyhow::anyhow!("mutex poisoned"))?;
*guard += 1;
```

**Pitfall.** Holding a `std::sync::MutexGuard` across an `.await` point.
`MutexGuard` isn't `Send`, so this either fails to compile (the containing
future becomes `!Send` and can't be spawned onto a multi-threaded runtime),
or, on a `current_thread` runtime, where it does compile, it blocks every
other task on that thread for the whole `.await`. Go's `sync.Mutex` has no
equivalent footgun, since a goroutine's held lock isn't tied to any trait
the scheduler cares about.

The default fix is to shrink the critical section so the guard is dropped
before the `.await`: restructure the code so locking, mutating, and
unlocking happen synchronously, with any async work (the fetch, the I/O
call) done before locking or after unlocking. Reach for `tokio::sync::Mutex`
only as the exception, when you can't restructure around the
`.await` (e.g. the async call itself has to happen while the data stays
locked). It's not a free substitute: holding a `tokio::sync::Mutex` guard
across an `.await` in a hot path serializes every concurrent task's async
work inside that critical section, because none of them can proceed until
the lock is released. Throughput collapses under load, with no compiler
warning and no panic to point you at the cause. That's exactly the failure
mode a Go engineer wouldn't think to look for, since reaching for the
"async-safe" Mutex feels like the correct, idiomatic choice.

As a smaller trap on the way there: `std::sync::PoisonError<MutexGuard<'_, T>>`
borrows from the guard, so it isn't `'static`; you can't propagate it
through `anyhow::Result` with a bare `?`. It has to be mapped to an owned
error first, as above.

## Cancellation: `CancellationToken`

**Concept.** `tokio_util::sync::CancellationToken` is a cooperative
cancellation signal: `cancel()` fires it, `cancelled()` is an awaitable that
resolves once fired, and `child_token()` derives a token that's cancelled
whenever its parent is (but not vice versa).

**Go comparison.** This is the async-Rust analogue of a `context.Context`
used purely for cancellation: `cancel()` ~ calling a context's cancel func,
`cancelled()` ~ selecting on `ctx.Done()`, `child_token()` ~
`context.WithCancel(parent)`. It's used here instead of hand-rolling
cancellation with a `watch::channel<bool>` because it already gives you the
parent/child tree Go engineers expect from `context.Context`; reimplementing
that on top of `watch` would just rebuild what `tokio-util` already
provides.

**Excerpt** (`rust/tokio_examples/cancellation.rs`):

```rust
tokio::select! {
    _ = token.cancelled() => {
        println!("{name}: cancellation observed, shutting down");
        return Ok(());
    }
    _ = sleep(Duration::from_millis(10)) => {
        tick += 1;
        println!("{name}: tick {tick}");
    }
}
```

**Pitfall.** Cancellation is cooperative, exactly as with `context.Context`.
A task that never races an `.await` against `token.cancelled()` (for
instance, one stuck inside a single long await with no `select!` around it)
will not stop just because the token was cancelled. Calling `cancel()` sets
a flag; it doesn't reach into a task and interrupt it.

## `spawn_blocking` for CPU-bound or blocking work

**Concept.** `tokio::task::spawn_blocking` runs a synchronous closure on a
separate thread pool sized for blocking work, keeping the async worker
threads free to keep polling other tasks.

**Go comparison.** Go rarely needs an equivalent: when a goroutine makes a
blocking syscall, the Go runtime detects the blocked M and spins up another
OS thread for the remaining goroutines automatically. Tokio's async worker
threads are cooperatively scheduled and get no such automatic backfill, so
blocking work has to be explicitly moved off them.

**Excerpt** (`rust/tokio_examples/spawn_blocking.rs`):

```rust
let cpu_result = tokio::task::spawn_blocking(|| fibonacci(30)).await?;
```

**Pitfall.** Calling an ordinary blocking function (`std::fs` I/O,
`std::thread::sleep`, a synchronous DB driver, a tight CPU loop) directly
inside an `async fn` compiles without any warning (Rust has no way to know
the call blocks) and silently stalls the whole worker thread. There's no
compiler signal pointing you at `spawn_blocking`; you have to know to reach
for it. And unlike Go's M:N scheduler, which backfills OS threads on demand
with no real ceiling you'd normally hit, `spawn_blocking`'s thread pool is
finite (512 threads by default, tunable via `Builder::max_blocking_threads`).
Under high enough concurrency, blocking calls queue for a free pool thread
rather than running immediately, so it isn't an unlimited escape valve.

## Error handling in async code: `anyhow`

**Concept.** `anyhow::Result<T>` is a type-erased error wrapper; `?`
propagates any error that implements `std::error::Error + Send + Sync +
'static`, and `.context("...")` attaches a message while preserving the
original error as the cause.

**Go comparison.** `anyhow::Result<T>` is Go's `(T, error)` pair collapsed
into one type. `?` is `if err != nil { return err }`, applied automatically.
`.context("...")` is `fmt.Errorf("doing X: %w", err)`. Printing the
resulting `anyhow::Error` with `{:?}` shows the whole "caused by" chain at
once, the same information `%w` chains give you piece by piece.

**Excerpt** (`rust/tokio_examples/error_handling.rs`):

```rust
async fn run() -> Result<u32> {
    let raw = fetch_config("missing.conf")
        .await
        .context("loading application config")?;
    Ok(raw)
}
```

**Pitfall.** The `'static` bound on `anyhow::Error` is easy to trip over:
any error type that borrows from something else (the `PoisonError` case in
the shared-state section is a real example from this tutorial) can't be
propagated with a bare `?` and needs an owned conversion first. In Go, every
`error` is already an interface value with no borrowing to worry about, so
this class of error never comes up.
