# Borrow-checker triage

Load when the compiler rejects a borrow, move, or lifetime. Error
titles below are `rustc --explain` wording. Each entry names what the
compiler is actually objecting to and the design change that resolves
it — not the `.clone()` that silences it.

## The rule that resolves most of them

A borrow ends at its **last use**, not at the end of its scope. So
the first thing to try is almost never a new type: it is moving the
last use of the borrow earlier, or hoisting the value out into its
own `let` so the borrow finishes before the conflict.

Two structural facts do most of the remaining work:

- **Disjoint fields borrow independently.** `&mut s.a` and `&mut s.b`
  coexist happily. But `s.method()` borrows *all* of `s`, so the same
  code fails once it goes through a method. Fix by destructuring
  (`let Struct { a, b, .. } = &mut s;`) or by making the helper a
  free function that takes the two fields.
- **Indices are references the borrow checker cannot see.** When a
  graph or tree fights you, store nodes in a `Vec` and refer to them
  by `usize`. You give up compile-time aliasing guarantees and gain a
  design that compiles; this is what arenas are.

## Moves

**E0382 — a variable was used after its contents have been moved.**
Usually a value passed by value into a loop body, or into a function
that only needed to read it. Take `&T`, or move the use after the
final read. Clone only when you genuinely need a second owner, and
say so.

**E0505 — a value was moved out while it was still borrowed.**
The borrow outlives the move. End it first: narrow the borrow's scope
or reorder so the last use precedes the move.

**E0507 — a borrowed value was moved out.** You have `&T` and need
`T`. The four real fixes, in order of preference: change the
signature to take `self` if the caller can give ownership;
`std::mem::take(&mut x)` when `T: Default`; `std::mem::replace(&mut x, v)`
when it is not; `Option::take()` when the field is already optional.
`clone()` is the fallback, not the first move.

## Aliasing

**E0499 — borrowed as mutable more than once.** Two live `&mut` to
overlapping data. Split the borrow across disjoint fields; use
`split_at_mut`, `iter_mut`, or `chunks_mut` for disjoint slice
regions; or scope the first borrow so it ends. If the two mutations
genuinely interleave over the same data, that is the design telling
you one owner should sequence both — an actor task or a method that
performs both steps.

**E0502 — borrowed as mutable while also borrowed as immutable** (or
the reverse). The archetype is `v.push(v[0])`. Compute first, then
mutate: `let first = v[0]; v.push(first);`. For maps, the `entry` API
exists precisely to collapse a lookup and an insert into one borrow.

**E0506 — assignment to a borrowed value**, and **E0503 — use of a
value while mutably borrowed**: same shape, same fix — the borrow has
to end before the write.

**E0596 — mutable borrow of a non-mutable variable.** Add `mut`, or
recognize that the function taking `&self` should take `&mut self`,
which is an API decision worth stating.

## Lifetimes

**E0597 — a value was dropped while still borrowed**, and
**E0515 — a reference to a local variable was returned.** Both mean
the data does not live as long as the reference to it. Either hoist
the owner so it outlives the borrow, or return an owned value. A
function that wants to return a reference into something it created
cannot; give the caller the owned value, or take a buffer parameter
to write into.

**E0716 — a temporary is dropped while a borrow is still in use.**
The classic is chaining off a temporary guard:
`let x = m.lock().unwrap().field;` drops the guard at the semicolon.
Bind the temporary first: `let guard = m.lock().unwrap();`.

**E0106 — missing lifetime specifier.** Before writing `'a`, ask
which input the output actually borrows from. If the answer is
"none", the return type should be owned. If the answer is "one of
them", elision usually handles it and the real error is elsewhere. A
`'a` that has to be threaded through a third type is the signal to
own the data or use an `Arc`.

## Closures and threads

**E0373 — a captured variable may not live long enough.** The closure
outlives what it borrowed. Make it `move` and give it owned data —
clone an `Arc` before the closure and move the clone in.

**E0521 — borrowed data escapes outside of closure.** Same cause,
usually from `thread::spawn` or `tokio::spawn`, which require
`'static`. Move ownership in, or use a scoped API
(`std::thread::scope`) where the borrow provably ends first.

## The async one worth naming separately

**E0277 with a `Send` bound** on a spawned future almost always means
a non-`Send` guard is alive across an `.await` — a
`std::sync::MutexGuard`, an `Rc`, a `RefCell` borrow. The compiler
points at the spawn; the bug is at the await. Drop the guard before
awaiting by scoping it, and only reach for `tokio::sync::Mutex` when
the lock genuinely must span the suspension.

## When interior mutability is the answer

Sometimes it is: a cache behind a shared reference, a value whose
mutation is invisible to callers. When you conclude that, name the
invariant that makes it sound at the point you introduce it, and
prefer `RefCell` in single-threaded code and `Mutex`/`RwLock` across
threads. `RefCell` moves an aliasing error from compile time to a
runtime panic, so it needs a reason, not a shrug.
