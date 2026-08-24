# Rust Basics 101 for Go Engineers

This is a walkthrough of `rust/basics_examples/*.rs`, in the same Cargo
binary crate (`tokio_hello_world`) as the Tokio tutorial. It's the
companion to `agent-team-workspace/research/rust-tokio/rust-tokio.md`, covering the ground that
tutorial assumes: basic types, data structures, the type system, and move
semantics. It's written for someone with deep Go experience and no prior
Rust exposure. Every section leans on a Go idiom you already know and
calls out where the Rust version diverges. Run any example with
`cargo run --example <name>` from inside `rust/`.

**Not covered here.** This doc stops at the type system and single-value
ownership. Three things a Go engineer will hit almost immediately after
this, left for a follow-up 201: recursive data structures (`Box<T>` shows
up below only as "the thing that gives a `dyn` trait object an owned
home," not in its more common role of giving a self-referential type like
a linked list or tree a compile-time-known size; porting this repo's own
`ListNode`/`TreeNode` solutions will need it); `Rc<RefCell<T>>` for
single-threaded shared mutable ownership when a single owner isn't enough;
and closures (`Fn`/`FnMut`/`FnOnce`, `move` closures; they show up
unexplained in a couple of excerpts below and interact directly with the
move/borrow rules this doc does cover). `if let`/`while let`, the sugar
for handling one `Option`/`Result` arm without a full `match`, is also
skipped in favor of always showing the exhaustive form.

## Go → Rust idiom map

| Go | Rust | Note |
|---|---|---|
| assignment/pass copies the header, both bindings alias the backing data | assignment/pass by value *moves* ownership | The old binding becomes unusable: a compile-time check, not a runtime race. |
| value receiver `(c Counter)` vs pointer receiver `(c *Counter)` | `&self` vs `&mut self` vs `self` | `self` (by value) consumes the receiver; Go has no equivalent, since a Go value receiver never destroys the caller's copy. |
| `[]T` (slice header: pointer+len+cap) | `&[T]` (borrowed view: pointer+len, no cap, no ownership) | Rust also has `[T; N]`, a fixed-size array with the length baked into the type (Go's rarely-used `[N]T`). |
| copy-or-share decided by convention/type shape | `Copy` (implicit bitwise dup) vs `Clone` (explicit `.clone()`) | The compiler enforces which one applies; `String`/`Vec<T>` are never `Copy` because they own a heap allocation. |
| `interface{}` + type switch, or a `Kind` field + optional fields | `enum` with per-variant data + exhaustive `match` | Missing a case is a compile error in Rust, not an easily-forgotten `default:`. |
| nil pointer / zero value standing in for "absent" | `Option<T>` (`Some(T)` / `None`) | No null pointer; the compiler forces both cases to be handled before you get at the `T`. |
| `(T, error)` return pair + `if err != nil` | `Result<T, E>` (`Ok(T)` / `Err(E)`) | Same idea, folded into one type; `match`/`?` replace the manual nil check. |
| `string`, indexable by byte (`s[i]`) | `String` (owned) / `&str` (borrowed view), **not** indexable by integer at all | Rust refuses to compile `s[i]`, precisely because it could split a multi-byte UTF-8 character. |
| `map[K]V` iteration re-randomized on every `range`, even of the same live map | `HashMap<K, V>` order fixed per-instance (random hash seed chosen once, at construction) | Don't rely on order from either, but Go's guarantee is strictly stronger. |
| implicit interface satisfaction (structural) | `trait` + explicit `impl Trait for Type` | Rust always has a grep-able declaration linking type to trait; Go infers it structurally. |
| interface value dispatch (always a runtime itable call) | generic + trait bound `<T: Trait>` (static, monomorphized) vs `&dyn Trait` (dynamic, vtable) | Rust makes the static/dynamic dispatch choice explicit; Go always dispatches dynamically through an interface. |
| `[T Constraint]` generics (Go 1.18+) | `<T: Trait>` trait bounds, composable with `+` and `where` | Same idea; Rust bounds don't need a constraint interface declared up front the way Go conventionally uses one. |

## Ownership and move semantics

**Concept.** Every value in Rust has exactly one owner. Assigning a
non-`Copy` value to a new binding, or passing it by value into a function,
*moves* ownership to the new location; the old binding is invalidated by
the compiler, not just left dangling by convention.

**Go comparison.** Assigning or passing a Go string, slice, or map copies
a small header (pointer, len, and for slices cap), and both the original
and the copy keep aliasing the same backing data: nothing stops you from
using both afterward. Rust has no equivalent "copy the header, share the
backing buffer" mode for owned heap types. A move transfers ownership
outright, so there's only ever one owner and no aliasing to reason about.

**Excerpt** (`rust/basics_examples/ownership_and_move.rs`):

```rust
let original = String::from("hello");
let moved = original; // ownership of the heap buffer transfers here

// `original` is no longer usable. Uncommenting the next line is a
// compile error: "value borrowed after move" / "value used after move".
// println!("{original}");
println!("moved now owns the string: {moved}");
```

**Pitfall.** Passing an owned value into a function by value moves it in,
exactly like a local assignment; the caller loses access once the call
is made, whether or not the function does anything observable with the
value. Go engineers often reach for `.clone()` reflexively the first time
they hit this, when a `&T`/`&mut T` borrow (see the next section) is
usually what they actually want and avoids the allocation a clone costs.

## Borrowing and references

**Concept.** `&T` is a shared, read-only reference; `&mut T` is an
exclusive, mutable reference. The borrow checker enforces one aliasing
rule at compile time: at any point, a value may have either many `&T`
borrows or exactly one `&mut T` borrow, never both at once.

**Go comparison.** Go has no borrow checker. A `*T` pointer can be read
and written from anywhere that holds it, including concurrently from
multiple goroutines, and it's on the programmer (or `go test -race`) to
catch conflicting access at runtime. Rust catches the equivalent conflict
statically: an invalid program simply fails to compile, so there's no
race to observe by running it.

**Excerpt** (`rust/basics_examples/borrowing_and_references.rs`):

```rust
fn print_len(v: &[i32]) {
    println!("length is {}", v.len());
}

fn push_one(v: &mut Vec<i32>) {
    v.push(v.len() as i32 + 1);
}
```

**Pitfall.** A signature returning a reference derived from more than one
input reference, like `fn longest<'a>(a: &'a str, b: &'a str) -> &'a str`,
needs an explicit lifetime annotation (`'a`) tying the output's validity to
the inputs'. This isn't a borrow-checker escape hatch: it's the compiler
asking you to state a constraint it can't infer on its own here. Full
lifetime elision rules (when Rust infers `'a` for you, mostly single-input
signatures) are out of scope for this 101; the short version is: if the
compiler complains a return type needs a named lifetime, it's telling you
the output's validity depends on more than one input, and you have to say
which. The same annotation shows up on a struct that holds a reference
field (`struct Parser<'a> { input: &'a str }`) for the same reason. The
compiler needs to know the struct can't outlive the data it borrows. If
you're porting a Go struct with a pointer field and the compiler asks for
a lifetime parameter, this is that case.

## `Copy` vs `Clone`

**Concept.** `Copy` types are duplicated implicitly and cheaply on
assignment (a bitwise copy of small, fixed-size, stack-only data);
everything else needs an explicit `.clone()` call to produce an
independent second copy.

**Go comparison.** Go doesn't check this distinction at the type level.
Whether an assignment "copies" or "shares" depends on the value's shape:
a fixed-size struct/array copies, a slice/map/pointer header shares. You
learn which is which by convention, not from anything the compiler tells
you. Rust turns it into a real trait: integers, floats, `bool`, `char`,
and tuples/arrays of `Copy` types implement `Copy` and duplicate
implicitly. `String` and `Vec<T>` never do, because they own a heap
allocation the compiler can't safely duplicate for free: silently
copying one under the hood would hide the allocation cost implicit-copy
semantics are supposed to guarantee doesn't exist.

**Excerpt** (`rust/basics_examples/copy_vs_clone.rs`):

```rust
let p1 = Point { x: 1, y: 2 };
let p2 = p1; // Point is Copy: this duplicates the struct, doesn't move it
println!("p1=({}, {}) p2=({}, {}), both usable", p1.x, p1.y, p2.x, p2.y);

let s1 = String::from("heap-allocated");
let s2 = s1.clone(); // explicit deep copy: separate heap buffers
```

(`#[derive(Copy, Clone)]` above is a derive macro: it generates the trait
implementation for you at compile time instead of you writing it by hand;
seen again on `ParseError`'s `#[derive(Debug)]` a few sections down.)

**Pitfall.** `#[derive(Copy, Clone)]` on your own struct only compiles if
every field is itself `Copy`. Adding a single `String` or `Vec<T>` field
to an existing `Copy` struct is a breaking change: every call site that
relied on implicit duplication (`let p2 = p1;` followed by using both)
now gets a move instead, and fails to compile wherever it used the old
binding afterward. Go has no equivalent "adding a field changes whether
assignment aliases or copies" trap, since that behavior was never
type-checked in the first place.

## Structs and methods

**Concept.** A `struct` defines the fields; an `impl` block attaches
methods to it, and each method chooses its receiver: `&self`, `&mut
self`, or plain `self`.

**Go comparison.** A Rust `struct` + `impl` block is a Go struct + its
methods, but Go spells the receiver's ownership as `(c Counter)` (value)
vs `(c *Counter)` (pointer). Rust spells the same choice explicitly on
`self`: `&self` reads (~ Go value receiver used for reading), `&mut self`
mutates in place (~ Go pointer receiver), and `self` takes ownership and
consumes the receiver entirely. Go has no equivalent, since a Go value
receiver's copy is always independently discardable, never the caller's
only copy.

**Excerpt** (`rust/basics_examples/structs_and_methods.rs`):

```rust
impl Counter {
    fn new(label: &str) -> Counter { /* ... */ }
    fn value(&self) -> i32 { self.value }
    fn increment(&mut self) { self.value += 1; }
    fn into_label(self) -> String { self.label }
}
```

**Pitfall.** Calling a `self`-by-value method like `into_label` consumes
the receiver the same way passing it to any other function would: the
variable is moved-from and unusable afterward, even though the syntax
(`c.into_label()`) looks identical to calling `c.value()`, which only
borrows. There's no visual cue at the call site telling you which kind of
receiver you're invoking; you have to know (or check) the method's
signature. Go's uniform "it's always safe to keep using the receiver
variable after any method call" assumption doesn't hold here.

## Enums and pattern matching

**Concept.** `enum` defines a sum type: a fixed set of variants, each of
which can carry its own, different data. `match` on an enum must be
exhaustive: every variant needs an arm, or the code doesn't compile.

**Go comparison.** Go has no sum type. The usual stand-ins are an
`interface{}` with a type switch (structurally open: a new type can
implement the interface anywhere in the codebase, and the type switch's
`default:` is easy to forget) or a struct with a `Kind`-style field plus a
pile of fields that are only meaningful for some `Kind` values. Rust's
`match` closes both gaps at once: the variant set is fixed at the `enum`
definition, and the compiler, not a forgotten `default:`, rejects a
`match` that doesn't cover every variant.

**Excerpt** (`rust/basics_examples/enums_and_pattern_matching.rs`):

```rust
enum Shape {
    Circle { radius: f64 },
    Rectangle { width: f64, height: f64 },
    Triangle { base: f64, height: f64 },
}

fn area(shape: &Shape) -> f64 {
    match shape {
        Shape::Circle { radius } => std::f64::consts::PI * radius * radius,
        Shape::Rectangle { width, height } => width * height,
        Shape::Triangle { base, height } => 0.5 * base * height,
    }
}
```

**Pitfall.** Adding a new variant to an `enum` is a breaking change to
every exhaustive `match` on it elsewhere in the codebase: each one fails
to compile until you add an arm for the new variant. Go engineers coming
from a type-switch-over-`interface{}` background usually experience this
as a pleasant surprise rather than a pitfall: it's the compiler doing, for
free, the "did I remember to update every switch when I added a new
case" audit you'd otherwise do by hand (or `grep`) in Go.

## The `Option<T>` type

**Concept.** `Option<T>` represents "a `T`, or nothing" as two variants,
`Some(T)` and `None`. There's no null pointer involved, and the compiler
requires both cases to be handled before you can get at the `T` inside.

**Go comparison.** A Go function signaling "no value" typically returns a
pointer that might be nil, or overloads the zero value (`0`, `""`,
`false`) to also mean "absent", which is ambiguous exactly when zero is
also a legitimate value. `Option<T>` removes the ambiguity: absence is a
distinct, type-level case, not a value indistinguishable from a real one.
Because getting the `T` out requires handling `None` (via `match`, or a
combinator like `unwrap_or`/`map`), there's no equivalent of a nil-pointer
dereference from a forgotten check.

**Excerpt** (`rust/basics_examples/option_type.rs`):

```rust
fn find_age(ages: &HashMap<&str, u32>, name: &str) -> Option<u32> {
    ages.get(name).copied()
}

// ...
let carol_age = find_age(&ages, "carol").unwrap_or(0);
let doubled = find_age(&ages, "alice").map(|age| age * 2);
```

**Pitfall.** `Option<T>` also has an `.unwrap()` (and `.expect("msg")`),
which panics immediately if the value is `None`. It's tempting for a Go
engineer to reach for `.unwrap()` as "the normal way to get the value
out," since it reads like the natural next step after seeing `Some`/`None`
in a signature. But it silently reintroduces the exact class of crash
`Option` exists to make you handle explicitly, just moved from "nil
pointer dereference" to "unwrap on a `None`". Prefer `match`, `unwrap_or`,
`unwrap_or_else`, or `map`/`and_then`, as in the excerpt above, and save
`.unwrap()`/`.expect()` for cases you can prove are unreachable (and even
then, `.expect("why")` at least documents the assumption) or for test
code, where a panic on failure is exactly what you want.

## The `Result<T, E>` type

**Concept.** `Result<T, E>` represents "a `T`, or an error `E`" as two
variants, `Ok(T)` and `Err(E)`: the type-level version of Go's `(T,
error)` return pair.

**Go comparison.** `Result<T, E>` folds Go's second return value into the
type itself: `match` on a `Result` is exhaustive the same way it is on
`Option`, so there's no equivalent of forgetting an `if err != nil` check:
the `T` simply isn't reachable without going through both arms (or a
combinator). This section stays at the type level on purpose; the `?`
operator, `anyhow::Result`, and `.context(...)` error-wrapping idioms used
throughout the async examples are already covered in
`agent-team-workspace/research/rust-tokio/rust-tokio.md`'s "Error handling in async code" section,
so they're not repeated here.

**Excerpt** (`rust/basics_examples/result_type.rs`):

```rust
#[derive(Debug)]
struct ParseError {
    input: String,
}

// Display for a human-readable message, Error (a marker trait beyond
// Display + Debug) so this plugs into anything expecting a real error type.
impl std::fmt::Display for ParseError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "{:?} is not a valid even number", self.input)
    }
}
impl std::error::Error for ParseError {}

fn parse_even(input: &str) -> Result<i32, ParseError> {
    // `?`: on Err, returns early with the error; on Ok, unwraps and
    // continues. Same operator agent-team-workspace/research/rust-tokio/rust-tokio.md covers with
    // anyhow — here it's just propagating this function's own error type.
    let n: i32 = input
        .parse()
        .map_err(|_| ParseError { input: input.to_string() })?;
    if n % 2 != 0 {
        return Err(ParseError { input: input.to_string() });
    }
    Ok(n)
}
```

**Pitfall.** Unlike `Option`, a `Result`'s `Err` variant carries a real
value (`E`), not just an absence marker, so ignoring a `Result` entirely
(binding it to nothing, or not using `?`/`match`) isn't a silent bug the
way ignoring Go's second return value can be if you skip the `if err !=
nil` check: the compiler emits an "unused `Result` that must be used"
warning specifically because `Result` is marked `#[must_use]`. It's a
warning, not a hard error, so it's still possible to compile past it.
But it's a compiler-level nudge Go's `_ = err` silent-discard pattern has
no equivalent of.

## Arrays and slices

**Concept.** `[T; N]` is a fixed-size array whose length `N` is part of
the type and checked at compile time; `&[T]` is a borrowed *view* over
a contiguous run of `T`s: a pointer and a length, nothing more.

**Go comparison.** Go's everyday `[]T` is a three-word header (pointer,
len, cap) over a backing array, plus a fixed-size `[N]T` array type most
Go code never touches directly. Rust makes the fixed/dynamic split the
default rather than the exception: `[T; N]` for compile-time-known sizes,
`&[T]` as a read-only view with no cap and no ownership at all. The
aliasing behavior diverges more than the types do, though: Go's `append`
may or may not reallocate depending on spare capacity, so two Go slices
sharing a backing array can silently alias or silently diverge at
runtime, purely depending on capacity you didn't necessarily choose.
Rust's growable equivalent is `Vec<T>` (next section); `&[T]` itself never
grows, and the borrow checker statically forbids any mutation that could
invalidate an outstanding `&[T]` view. See the pitfall below.

**Excerpt** (`rust/basics_examples/slices_and_arrays.rs`):

```rust
let mut heap_vec = vec![1, 2, 3, 4, 5];
heap_vec.push(6); // fine: no outstanding borrow yet
let view: &[i32] = &heap_vec[..3];
println!("vec = {heap_vec:?}, first-3 view = {view:?}, sum(view) = {}", sum(view));

// `view` borrows heap_vec's backing storage. Uncommenting the next
// line is a compile error ("cannot borrow `heap_vec` as mutable
// because it is also borrowed as immutable"):
// heap_vec.push(7);
```

**Pitfall.** A `push` that grows a `Vec` past its current capacity
reallocates and moves the backing buffer: in Go, this is exactly the
scenario where a slice taken before the `append` keeps pointing at the
*old*, now-stale backing array, silently diverging from the grown one
with no error at all. Rust doesn't let that situation exist: the borrow
checker refuses to compile a `push` (or any other mutation) while a `&[T]`
view derived from the same `Vec` is still alive, because that mutation
could invalidate the view's pointer. The fix is ordering: mutate before
taking the view, or drop the view (let it go out of scope) before
mutating again, not a runtime check.

## `String` vs `&str`, and UTF-8

**Concept.** `String` is an owned, growable, heap-allocated string;
`&str` is a borrowed view into UTF-8 text (owned by a `String`, a string
literal, or something else). Every `str`/`String` is guaranteed valid
UTF-8 by the type: there's no "maybe invalid" string state to worry
about.

**Go comparison.** Go's `string` is already immutable and UTF-8-oriented,
and `s[i]` gives you the `i`-th *byte*, which Go permits even though it
can land you in the middle of a multi-byte rune, silently producing a
value that isn't valid UTF-8 on its own. Rust splits "owned" (`String`)
from "borrowed view" (`&str`) the same way it splits `Vec<T>`/`&[T]`, but
goes further on indexing: `s[i]` for a `String`/`str` doesn't compile *at
all* (integer indexing isn't implemented for that type, full stop),
precisely because slicing at a non-character boundary would produce
invalid UTF-8, and Rust treats that as unacceptable rather than as a
foot-gun you can opt into the way Go does.

**Excerpt** (`rust/basics_examples/strings_str_vs_string.rs`):

```rust
let owned: String = String::from("héllo"); // owned, heap-allocated
let borrowed: &str = &owned; // a view into owned's buffer, no copy

// len() counts bytes, not characters: 'é' is 2 bytes in UTF-8, so
// byte length (6) and char count (5) diverge for this string.
println!("byte length = {}, char count = {}", owned.len(), owned.chars().count());
```

**Pitfall.** `String::len()` returns the byte length, not the character
count: visible above, where `"héllo"` reports `len() == 6` but
`chars().count() == 5`, because `é` takes 2 bytes in UTF-8. This mirrors
`len(s)` in Go, which also counts bytes, not runes, so it isn't a new
gotcha for a Go engineer per se, but it compounds with slicing: Rust's
`&owned[0..1]` range-slices by *byte* offset (like Go's `s[0:1]`) and
panics at runtime if the offset lands mid-character, rather than silently
truncating a rune the way byte-slicing a Go string can. `.chars()` (walks
Unicode scalar values) is almost always the right tool when you mean
"characters," in both languages.

## `Vec<T>` and `HashMap<K, V>`

**Concept.** `Vec<T>` is a growable, heap-allocated, owned sequence;
`HashMap<K, V>` is a hash map. Both are ordinary generic library types,
not language built-ins.

**Go comparison.** Go's `[]T` and `map[K]V` are built into the language
itself; before Go 1.18 there was no way to write your own generic
container that looked anything like them. Rust never needed that
special-casing: `Vec<T>` and `HashMap<K, V>` are defined using the exact
same generics and trait bounds available to any type you write yourself.
See the generics section below. One concrete behavioral difference worth
knowing, and it's the opposite of what you might guess: Go deliberately
re-randomizes map iteration order on *every* `range`, even ranging over
the exact same live map twice in a row without modifying it in between.
Rust's default hasher (`RandomState`) instead picks one random seed when
a `HashMap` is constructed, and that seed, and so the iteration order for
a given set of keys, stays fixed for that instance's lifetime; iterating
it twice in a row yields the same order both times. Neither guarantee is
something to build logic on, but Go's is the stronger, per-iteration one.

**Excerpt** (`rust/basics_examples/vec_and_hashmap.rs`):

```rust
let mut scores: HashMap<String, u32> = HashMap::new();
for name in &names {
    scores.insert(name.clone(), name.len() as u32);
}

// entry() is the idiomatic get-or-insert, replacing Go's
// "v, ok := m[k]; if !ok { m[k] = ... }" pattern.
*scores.entry(String::from("dave")).or_insert(0) += 1;
```

**Pitfall.** Indexing a `HashMap` with `map[key]` (via the `Index` trait,
as in `scores[key]`) panics if the key is missing: there's no zero-value
fallback the way Go's `m[k]` silently returns the zero value for a
missing key. The safe equivalent of Go's `v, ok := m[k]` idiom is
`.get(key)`, which returns `Option<&V>`, not a bare value; reach for
`.entry(key).or_insert(default)` specifically for the get-or-insert
pattern, as shown above, rather than a manual "check then insert" that
looks up the key twice.

## Traits vs interfaces

**Concept.** A `trait` declares a set of methods; `impl Trait for Type`
implements it for a specific type. Rust offers two dispatch strategies: a
generic function with a trait bound (`fn f<T: Trait>(x: T)`) is
monomorphized into a separate compiled copy per concrete type (static
dispatch, no runtime indirection), while `&dyn Trait` erases the concrete
type behind a vtable and dispatches at runtime.

**Go comparison.** A Go interface is satisfied implicitly: any type with
the right method set implements it, with no declaration linking type to
interface anywhere in the code. Rust traits are explicit: `impl Trait for
Type` is required even when the methods already exist on the type, so a
reader (or the compiler) always knows which types implement which trait
from a grep-able declaration, not from structural inference. Dispatch is
also a real choice in Rust that Go doesn't offer: Go always dispatches
through an interface value's itable at runtime, full stop; Rust lets you
pick zero-overhead static dispatch via a trait bound, and reserves `&dyn
Trait`'s vtable-based dynamic dispatch (the direct analogue of a Go
interface value) for when you specifically need it (e.g. a
heterogeneous collection, shown below).

**Excerpt** (`rust/basics_examples/traits_vs_interfaces.rs`):

```rust
// Static dispatch: monomorphized per concrete type at compile time.
fn describe_static<T: Shape>(shape: &T) {
    println!("[static] {} has area {:.2}", shape.name(), shape.area());
}

// Dynamic dispatch: one function body, dispatched at runtime via a vtable.
fn describe_dynamic(shape: &dyn Shape) {
    println!("[dynamic] {} has area {:.2}", shape.name(), shape.area());
}

// This is the direct answer to "how do I write `[]Shape` in Rust": a
// generic Vec<T> forces every element to be the *same* concrete T, since
// the bound is monomorphized once per T, not per element. Box<dyn Shape>
// gives each element an owned, heap-allocated home behind a shared vtable
// pointer, which is what lets Circle and Square coexist in one Vec at all.
let shapes: Vec<Box<dyn Shape>> = vec![Box::new(circle), Box::new(square)];
for shape in &shapes {
    describe_dynamic(shape.as_ref());
}
```

**Pitfall.** Reaching for a plain generic `Vec<T>` first, out of Go habit,
and then being stuck when you need to mix concrete types (`Vec<T: Shape>`
flatly cannot hold a `Circle` and a `Square` together) is one of the more
common early stumbles: the fix is the `Vec<Box<dyn Shape>>` above, not a
workaround.

## Generics and trait bounds

**Concept.** A generic function or type parameterized by `<T>` can
constrain `T` with one or more trait bounds (`<T: Trait>`), restricting it
to types that implement that trait: the same idea as Go 1.18+'s
constraint interfaces, spelled differently.

**Go comparison.** Go generics use a constraint interface in the type
parameter list (`constraints.Ordered`, or a custom interface listing
allowed types/methods), written as `[T Constraint]`. Rust's `<T: Trait>`
bounds are the same idea but compose more directly: multiple bounds
combine with `+` (`T: Ord + Clone`) inline, or move into a separate
`where` clause once the list grows, without first declaring a named
constraint interface the way Go conventionally does.

**Excerpt** (`rust/basics_examples/generics_and_trait_bounds.rs`):

```rust
fn largest<T: PartialOrd + Copy>(items: &[T]) -> T {
    let mut max = items[0];
    for &item in items.iter() {
        if item > max {
            max = item;
        }
    }
    max
}

fn describe<T>(label: &str, value: T) -> String
where
    T: Display + Clone,
{
    let copy = value.clone();
    format!("{label}: {copy}")
}
```

**Pitfall.** Adding a `Copy` bound (as `largest` does above) changes what
the function is allowed to do with `T`, not just what types it accepts:
`let mut max = items[0];` only compiles because `T: Copy` lets that
assignment duplicate the element instead of trying to move it out of a
borrowed slice (`items: &[T]`), which would fail to compile, since you
can't move a value out of something you only borrowed. Drop the `Copy`
bound and this exact function body stops compiling, with an error about
moving out of a shared reference, even though nothing about the
function's *signature* looks like it should care about ownership. Go's
generics have no analogous "this bound changes which expressions are
legal in the body" effect, since Go never distinguishes copy-by-default
from move-by-default in the first place.
