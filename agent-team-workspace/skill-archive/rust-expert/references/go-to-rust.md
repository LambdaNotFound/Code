# Go to Rust

Load when translating Go, or when a Go instinct is shaping the
design. This crate's examples name the mapping in every file
(`//! Go comparison:`), so the work of writing it down is part of the
house style, not an aside.

The mapping matters most where the two languages look alike and
behave differently. Those rows are the ones worth stating out loud.

## Concurrency

| Go | Rust | What actually differs |
|---|---|---|
| `go f()` | `tokio::spawn(fut)` | A goroutine starts running immediately. A future does **nothing** until polled, so an un-awaited, un-spawned async call is a no-op the compiler only warns about. |
| `chan T` | `mpsc::channel(n)` | `Receiver` is not `Clone`, so the compiler guarantees a single consumer. In Go any goroutine may receive, and nothing stops a second one being added later. |
| unbuffered `chan` | no exact equivalent | Go's unbuffered channel is a rendezvous: the sender blocks until a receiver takes it. Tokio's `mpsc` requires a capacity of at least 1, so the sender may run one ahead. Use `oneshot` per message when the rendezvous is what you need. |
| `select { }` | `tokio::select!` | The branches that do not win are **dropped mid-poll**, cancelling that operation. Go's cases are single operations with nothing to cancel. This is the single biggest translation hazard — see async-shapes.md. |
| `sync.WaitGroup` | `JoinSet`, `join!`, `try_join!` | `JoinSet` also gives you `abort_all` and propagates panics, which `WaitGroup` cannot. |
| `sync.Mutex` | `Mutex<T>` | Rust's mutex **owns the data**; there is no way to touch the value without locking. Go's mutex guards a convention, so "forgot to lock" is a whole bug class that does not exist here. |
| `context.Context` | drop, or `CancellationToken` | Rust cancels by dropping the future — no cooperative check required at every layer. `tokio_util::sync::CancellationToken` is for when you need an explicit signal; this crate already depends on `tokio-util`. |
| goroutine leak | task leak | Same failure, different tell: a dropped `JoinHandle` detaches the task, so panics vanish and shutdown does not wait. `JoinSet` is the structured answer. |

## Errors and nil

| Go | Rust | What actually differs |
|---|---|---|
| `(T, error)` | `Result<T, E>` | Both treat errors as values, but `Result` is `#[must_use]`: ignoring one is a warning, not a silent `_`. |
| `if err != nil { return err }` | `?` | `?` also converts the error type through `From`, which is why the error enum design in the skill body matters. |
| `errors.Is` / `errors.As` | `matches!`, `source()`, downcast | Wrapping is `#[from]`/`#[source]` on a `thiserror` enum rather than `%w`. |
| `nil` | `Option<T>` | There is no null pointer and no typed-nil trap. An `Option<Box<T>>` is the same size as a pointer, so the safety is free. |
| `panic` / `recover` | `panic!` / `catch_unwind` | Neither crosses a thread boundary. In Rust a panic in a spawned task surfaces through the `JoinHandle`, if you kept it. |

## Types and interfaces

| Go | Rust | What actually differs |
|---|---|---|
| implicit interface satisfaction | `impl Trait for Type` | Rust is nominal and explicit: a type implements a trait only where someone wrote it, and the orphan rule says who may write it. There is no accidental satisfaction, and no accidental breakage. |
| `interface{}` / `any` | enum, generic, or `Box<dyn Trait>` | Prefer an enum when the set of types is closed — you get exhaustive matching, which is the thing `any` costs you. |
| struct embedding | composition, trait default methods | `Deref` can imitate embedding and should not; it is for smart pointers. |
| zero values | `Default`, explicitly | Rust has no implicit zero value and no partially initialized struct. `#[derive(Default)]` when a sensible one exists. |
| `v, ok := m[k]` | `map.get(&k) -> Option<&V>` | And `entry()` for the check-then-insert pattern, which is one borrow instead of two. |
| value vs pointer receiver | `self`, `&self`, `&mut self` | The receiver states the ownership transfer, so "does this method mutate a copy?" is answered by the signature. |

## Data and memory

The one to internalize: **Go's slice aliasing bugs are compile errors
in Rust.** In Go, `b := a[:2]` shares a backing array with `a`, and
`append(b, x)` may overwrite `a[2]` depending on capacity. In Rust,
`Vec<T>` owns and `&[T]` borrows, and the borrow checker will not let
you hold a slice across a mutation of its `Vec`. What felt like
paranoia is the compiler removing a class of bug you were carrying.

Similarly: a Go map read during a concurrent write is a fatal runtime
error; in Rust it does not compile without a lock or a `Send + Sync`
story.

## Tooling

`go vet` and `golangci-lint` map to `cargo clippy -- -D warnings`;
`gofmt` to `cargo fmt`; `go test ./...` to `cargo test`; `go test
-race` has no equivalent because the guarantee is static. Go's
table-driven test convention translates directly — an array of cases
in a loop, or `rstest` if a dependency is welcome.

## The advice that transfers

Idiomatic beats optimal, and clarity beats clever — the same as Go.
Do not reach for SIMD, const generics, custom allocators, or `no_std`
because Rust makes them available; reach for them when a profile says
to. A Go engineer's instinct for simple, boring, readable code is
correct here too.
