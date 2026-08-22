# Leaky Bucket Rate Limiting

A parallel research exercise: design and implement a leaky-bucket rate
limiter for a system with many concurrent callers, in both Go and Rust,
then have separate reviewers stress-test the design and the code. This
doc aggregates all of it: the mechanism, the two implementations, the
architecture and first-principles critiques, and two real bugs the code
review found and that are now fixed in the working code.

Working code: [`golang_concurrency/leaky_bucket.go`](../../golang_concurrency/leaky_bucket.go)
(+ [`leaky_bucket_test.go`](../../golang_concurrency/leaky_bucket_test.go)) and
[`rust/tokio_examples/leaky_bucket.rs`](../../rust/tokio_examples/leaky_bucket.rs)
(run with `cargo run --example leaky_bucket` from inside `rust/`).

## Read this before the algorithm

The task that produced this doc was phrased as "we need a leaky bucket
to rate-limit a system's traffic." That sentence names a mechanism, not
a problem. No protected resource, no traffic shape, no deployment
topology (one process or a fleet) was ever specified. Leaky-bucket-as-
a-meter is mathematically equivalent to token bucket: same admit/reject
decisions, described from opposite ends. So absent a real reason to
delay traffic rather than just admit or reject it, "leaky bucket vs.
token bucket" is often a choice of vocabulary, not of algorithm. What
makes leaky bucket a different mechanism is specifically the
queueing/shaping variant built here: an internal buffer that delays
excess traffic instead of merely marking or dropping it.

This repo's own `golang_concurrency/` package is a library of hand-built
concurrency-primitive demos (lock-free stack, lock-free counter, a
`Waiter` type) with no production service behind any of them, and
nothing in `system_design/` names a protected resource with real
capacity numbers either. Read what follows as a concurrency-pattern
exercise (a solid one, since it forces dealing with a drain loop,
shutdown lifecycle, and a full-vs-not-full policy decision), not as a
production recommendation. Before reusing either implementation to
protect something real, answer one question first: is the actual
constraint a *rate* (requests per second, independent of how long each
one takes) or a *concurrency* limit (requests in flight, bounded by a
fixed pool of downstream capacity)? The two only coincide when request
duration is constant, and a rate limiter is the wrong tool for a
concurrency constraint.

## The mechanism

```
callers --> [ o o o o . . ]  --tick, every leakInterval-->  one leaks out
            ^^^^^^^^^^^^^^
            bounded queue, size = capacity
```

Arrivals fill a bounded queue (the "bucket") up to its capacity; a
single background process drains it at a fixed, constant rate
regardless of how bursty the input is. Two things need a policy
decision, not just an implementation:

- **What happens when the bucket is full:** reject the new arrival
  immediately, or make the caller wait for a slot.
- **How the drain is paced:** a literal timer that removes one item
  per tick, or an accounting scheme that computes elapsed time since
  the last check and derives how much capacity should have freed up,
  with no background process at all.

Both implementations here chose the literal, ticker-driven queue (the
more direct reading of "leaky bucket" as a queue, and the richer
concurrency exercise), and both expose *both* full-bucket policies
rather than picking one, since that's a caller-facing contract, not an
internal detail.

## Design review: the decisions that matter

*(Full findings from the architecture review are folded in below;
see the [reversibility](#reversibility) and [operational
cost](#operational-cost-and-failure-mode) sections for the sharpest
points.)*

### Shaping vs. policing

Leaky bucket enforces a *constant* output rate: that's traffic shaping.
If the actual goal is self-protection (admit everything the system can
handle right now, reject only the excess), that's policing, and it
doesn't require delaying traffic that arrives when the system is idle.
A leaky bucket metes out even a first-burst-after-idle at the fixed
drain rate, adding pure latency nothing downstream asked for. Leaky
bucket earns its place only when something needs a constant dispatch
cadence; otherwise a token bucket (spendable burst credit that
accumulates while idle) or a sliding-window counter matches "protect
the system, allow natural bursts" more directly, more cheaply, and
without a background task to manage.

### Bounded queue vs. semaphore

A semaphore bounds concurrency (N requests in flight); a queue bounds
throughput over time (N per second). They answer different questions,
so the queue is the correct primitive for *rate*, but it adds a second
real tunable, since queue depth combined with drain rate determines
injected latency (roughly depth ÷ rate, by Little's Law).
Under-provision the queue and it rejects aggressively; over-provision
it and overload shows up as silently growing latency instead of a
visible error. That constant deserves the same scrutiny as the rate
itself.

### Reject-on-full vs. block-on-full

Block-on-full couples a caller's lifetime to the limiter's internals;
under sustained overload, blocked callers accumulate, and without a
mandatory deadline at every call site that's a resource-leak-shaped
failure mode invisible until memory or goroutine-count pressure shows
up elsewhere. Reject-on-full turns overload into an immediate,
countable error and pushes the retry/backoff decision to the caller,
where it belongs: the better default. If a blocking mode is offered at
all, a cancellable context must be structurally required at the API
boundary, not optional; both implementations here do exactly that
(Go's `Wait(ctx context.Context) error`, mirroring the existing
`Waiter.Wait` idiom in this package).

This is also the least reversible decision in the whole design: it's
an API contract, not an internal detail. Switching it later means
touching every call site.

### Ticker-drain loop vs. timestamp-delta accounting

A ticker loop needs a goroutine/task with a full lifecycle: started,
stopped on shutdown, and if it dies, the limiter silently either
wedges every blocked caller or stops enforcing the limit entirely,
depending on the failure branch. That's a permanent line in every
goroutine/task dump. Timestamp-delta accounting (compute elapsed time
since the last call, derive how many slots should have freed up) needs
no background process at all, is what `golang.org/x/time/rate`'s token
bucket actually does, and is trivially unit-testable with a fake
clock. For plain admission control, delta accounting strictly
dominates: same guarantees, one fewer lifecycle to manage. The ticker
loop only earns its cost when the leak must have an externally
observable side effect on its own schedule, i.e., you're actually
dispatching queued work, which loops back to the shaping-vs-policing
question above.

### In-memory single-node vs. externally coordinated

Correctly out of scope for a first pass: no network round-trip in the
hot path, no external dependency. But it's the sharpest now-vs-later
trap here: an in-memory bucket's configured rate implicitly means "per
process," and nothing in the code says so. The day a second replica
goes up, the effective aggregate limit silently becomes `configured
rate × live node count`. Nobody decided that number; it just happens.
Whatever constant configures the rate should say this explicitly.
Building a pluggable backend abstraction now, with exactly one
implementation in scope, would be pure speculative generality. Add it
when a second real backend exists to justify it.

### What breaks first at 10x concurrent callers

The single point of shared mutable state every caller passes through
is, in one sense, *supposed* to be the bottleneck (that's the
limiter's job), but the chosen design determines how it degrades:

- **Channel-queue + ticker-drain** (what's implemented here): the
  channel send is the contention point. With block-on-full, 10x
  callers means 10x blocked goroutines/tasks queuing behind one
  fixed-rate drain that by design doesn't scale with caller count.
  Symptom: growing blocked-caller count and growing p99 latency, not
  CPU or memory pressure directly. The specific bottleneck is the
  single drain loop plus its ticker, serial by construction.
- **Mutex + timestamp-delta** (the road not taken): the mutex becomes
  the hot lock, but the critical section is O(1), so this degrades as
  measurable lock contention rather than caller pileup. It scales more
  gracefully at 10x because no per-call step depends on an external
  process's wakeup cadence.

### Reversibility

| Decision | Cost to reverse later |
|---|---|
| Reject-on-full vs. block-on-full | **High**: an API contract; every call site changes. |
| In-memory vs. externally coordinated | Cheap in code (swap in a shared store behind the same interface), expensive in semantics: someone has to notice the rate's meaning changed from per-process to global and coordinate the cutover. |
| Bounded queue vs. semaphore capacity | Cheap internally, but the tuned queue-depth constant tends to get baked into callers' own timeout assumptions once live. |
| Ticker-loop vs. delta accounting | Cheap, purely internal behind the same `Allow`/`Wait` signature, as long as nothing outside started depending on the ticker's dispatch timing as a side channel. |

### Operational cost and failure mode

A block-on-full limiter with no enforced deadline is a textbook
on-call trap: it presents as generic "everything got slow" with no
direct error to chase, and diagnosing it requires already knowing to
pull goroutine dumps or limiter-internal metrics that have to have
been instrumented in advance. Reject-on-full converts the same
overload into a countable error metric almost for free: a materially
smaller on-call surface. The ticker-loop implementation also asks
every future reader to reason about a second lifecycle (leak-on-no-
shutdown, behavior on panic, cleanup ordering) on top of the rate
logic itself, for no correctness benefit in the single-process,
no-external-dispatch case this repo actually needs.

## Go implementation

`golang_concurrency/leaky_bucket.go`, package `concurrency`, stdlib
only. The bucket is a buffered `chan struct{}` (the queue) drained by
one background goroutine ticking every `leakInterval`:

```go
func (lb *LeakyBucket) Allow() bool {
	select {
	case lb.queue <- struct{}{}:
		return true
	default:
		return false
	}
}

func (lb *LeakyBucket) Wait(ctx context.Context) error {
	select {
	case <-lb.stop:
		return ErrClosed
	default:
	}
	select {
	case lb.queue <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-lb.stop:
		return ErrClosed
	}
}
```

**Concurrency invariant:** `stop` is closed exactly once, via
`sync.Once` in `Close()`, mirroring the existing `Waiter.Close` pattern
in this package. `queue` itself is never closed: closing a buffered
channel while a goroutine may still be sending to it (blocked in
`Wait`) is a send-on-closed-channel panic waiting to happen, so `stop`
is a separate signal both `Wait` and the drain loop select on, and
`queue` is simply abandoned to the GC once the bucket is unreachable.
`Allow` never checks `stop`, since it never blocks and a stopped drain
loop just means the bucket stops leaking, not that it becomes invalid.

Verified: `gofmt -l`, `go vet`, and `go test -race` (including repeated
`-count=5` runs to rule out flakiness in the timing-based tests) are
all clean. Tests cover burst rejection, steady-state rate over a time
window, 200 concurrent `Allow` callers, 20 concurrent `Wait` callers,
cancellation, `Close` idempotency under concurrent callers, and the
constructor's panic on non-positive capacity.

## Rust implementation

`rust/tokio_examples/leaky_bucket.rs`, a Tokio example registered in
`rust/Cargo.toml`. The bucket is a bounded `mpsc::channel<Request>`;
producers hold cloned `Sender`s and `try_send`; a single task owns the
`Receiver` for the program's lifetime and drains one item per tick:

```rust
async fn leak(mut rx: mpsc::Receiver<Request>) {
    let mut ticker = interval_at(Instant::now() + LEAK_INTERVAL, LEAK_INTERVAL);
    loop {
        ticker.tick().await;
        match rx.try_recv() {
            Ok(req) => println!("leak: processed request {} from producer {}", req.seq, req.producer),
            Err(TryRecvError::Empty) => {}
            Err(TryRecvError::Disconnected) => return,
        }
    }
}
```

**Ownership decision:** the bucket needs no `Arc<Mutex<_>>` at all:
`mpsc::Receiver` has no `Clone` impl, so the compiler, not convention,
guarantees at most one task ever drains it. Rejected alternative: a
`tokio::sync::Semaphore` refilled on an interval, which models a
*token* bucket (bursty admission by permit count) rather than a leaky
bucket, since it never actually queues pending requests. Where a Go
engineer's instinct would reach for a mutex-guarded queue plus a
`sync.Cond` to wake the drainer, Rust's channel type closes off the
"more than one reader" possibility statically. There's no equivalent
guarantee from a bare Go `chan`.

`interval_at(Instant::now() + period, period)` is used instead of bare
`interval(period)` specifically because `interval` fires its first
tick immediately, unlike Go's `time.NewTicker`, which always waits out
a full period first. `interval_at` matches the behavior a Go engineer
would otherwise assume by default.

Verified: `cargo build`, `cargo clippy --all-targets -- -D warnings`,
`cargo test`, and `cargo run --example leaky_bucket` are all clean.

## Bugs the code review found (and fixed)

Both implementations passed every test, `go vet`, and `clippy` clean,
and both still had a real correctness bug, found by writing throwaway
repro tests rather than trusting green CI. Both are now fixed in the
files linked above.

### Go: `Wait` could report success on an already-closed bucket

`select` chooses uniformly at random among *all* ready cases, not the
one that became ready first. Once `Close()` had already completed,
`Wait` on a bucket with spare queue capacity had two ready branches:
send-succeeds and `stop`-fired. So roughly half the time it returned
`nil` instead of `ErrClosed`, even though the drain goroutine had
already exited and that admitted token would never leak. Measured
empirically: 499/1000 spurious successes. The existing test only ever
pre-filled the bucket before calling `Wait`, which made the send
branch permanently unready and structurally unable to hit the race.
Fix: a non-blocking check of `stop` before the blocking select, so the
already-closed case is deterministic, while the race with a
concurrent `Close` remains an acceptable, unavoidable one.

### Rust: the drain loop could erase the rate limit after any idle period

The original loop called `rx.recv().await` after `ticker.tick().await`,
which blocks (and stops polling the ticker) whenever the queue is
empty. `tokio::time::interval`'s default missed-tick behavior tracks
absolute deadlines and, once polled again, fires immediately for every
deadline that elapsed while unpolled. So after any idle stretch, a
backlog of "missed" ticks paid off nearly instantly: three queued
items measured draining at `0ns` and `0ns` gaps instead of one per
15 ms. That defeats the entire point of a leaky bucket for exactly the
bursty-after-idle traffic it exists to smooth. The file's own test
never caught it because it enqueues before the drain task's first
tick, so the queue is never empty at a tick boundary in that test.
Fix: `rx.try_recv()` instead of `rx.recv().await`, matching what the
file's own doc comment already claimed happened and what the Go
implementation's drain loop actually does. Confirmed by rerunning the
same idle-then-burst repro, which now measures a clean 15 ms/15 ms
gap.

## If you're about to use this for something real

Answer this first: what specific resource does this protect, and is
its limit shaped by rate or by concurrency? If there's a real number,
that answer picks the mechanism (leaky-bucket queue if something needs
constant-cadence dispatch; token bucket or a plain counter if it just
needs to survive bursts; a sized semaphore if the constraint is
in-flight requests, not requests per second). It isn't picked by which
name sounds more rigorous. If there's no real number yet, treat what's
here as what it is: a solid pair of concurrency-primitive demos, not a
production rate limiter.
