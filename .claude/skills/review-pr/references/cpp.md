# C++ review reference

Load when the diff contains C++. Ordered by how often the item is
the actual bug. Each item names what to look for, why it breaks, and
the fix. Check which standard the project targets before applying
anything version-dependent.

## Lifetime and dangling

The dominant bug class. Every reference, pointer, view, and capture
is a lifetime claim.

- **`string_view` or `span` bound to a temporary.**
  `std::string_view sv = obj.get_string();` dangles the moment the
  statement ends if `get_string` returns by value. A `string_view`
  parameter is fine; a `string_view` *member* or return value is a
  lifetime contract that must be documented and checked at every
  construction site.
- **Range-for over a subobject of a temporary.**
  `for (auto& x : f().items())` — lifetime extension applies to a
  temporary bound directly to the reference, not to one reached
  through a member call. Undefined before C++23, which fixed the
  range-for case specifically. Bind the temporary to a named local
  first.
- **Returning a reference or pointer to a local**, or to a
  by-value parameter.
- **Invalidation.** `vector` push_back/insert invalidates every
  iterator, pointer, and reference on reallocation;
  `unordered_map` invalidates iterators on rehash; erasing in a
  loop needs `it = c.erase(it)`. Flag any reference held across a
  mutation of its container.
- **`[&]` capture in a lambda that outlives its scope** — stored in
  a member, posted to a thread pool, handed to an async callback.
  Capture by value, or by explicit named captures.
- **`[=]` in a member function captures `this`, not a copy of the
  object.** Every member access through it dangles once the object
  dies. Use `[*this]` (C++17) when a copy is what is meant.
- **`shared_ptr` cycles leak.** Back-pointers and parent links are
  `weak_ptr`.

## Ownership and RAII

- Raw owning pointers, or `new`/`delete` outside a type whose job is
  to own. Use `make_unique`/`make_shared`; every resource — memory,
  file, socket, lock, handle — under RAII.
- **Rule of 0/3/5.** A user-declared destructor suppresses the
  implicit move constructor and move assignment, so the class
  silently copies where it used to move. Declare all five or none;
  prefer none by pushing ownership into members that manage
  themselves.
- **Move operations must be `noexcept`.** Without it `vector` growth
  falls back to copying to preserve the strong exception guarantee,
  which turns an O(n) move into an O(n) deep copy.
- **Use after move.** A moved-from object is valid but unspecified:
  it may be destroyed or assigned to, and reading anything else is a
  bug even where it is not UB.
- Self-assignment in a hand-written `operator=`.
- **Virtual destructor** on any class ever deleted through a base
  pointer. **Slicing** wherever a polymorphic type is passed or
  stored by value.

## Undefined behavior

Not "sometimes wrong" — the optimizer is entitled to assume it never
happens, so the failure often appears far from the cause.

- Signed integer overflow; shift by more than or equal to the bit
  width; division by zero.
- **Strict aliasing.** Type-punning through a cast is UB; use
  `memcpy` or `std::bit_cast` (C++20).
- Reading an uninitialized value; out-of-bounds indexing (`v[i]`
  does not check — `.at()` does); `memcpy` on a non-trivially-copyable
  type; an invalid `static_cast` down a hierarchy.
- Unsequenced modification: `i = i++`, `f(i++, i++)`.
- Every `reinterpret_cast` needs a comment justifying it.

## Initialization

- Uninitialized members. Prefer default member initializers so no
  constructor can forget one.
- `auto x{1}` deduces `int` (C++17); `auto x = {1}` deduces
  `initializer_list<int>`. Braced init also forbids narrowing, which
  is a reason to prefer it.
- Most vexing parse: `Widget w();` declares a function.
- **Static initialization order across translation units is
  unspecified.** A namespace-scope object depending on another in a
  different TU is a latent crash. Use a function-local static.
- `explicit` on single-argument constructors and on conversion
  operators, unless the implicit conversion is genuinely wanted.

## Concurrency

- A data race is UB. Every object mutated by one thread and read by
  another needs a mutex or an atomic — `volatile` is not a
  synchronization primitive.
- **`condition_variable::wait` without a predicate.** Spurious
  wakeups are permitted; the predicate overload is the only correct
  form.
- **Relaxed memory orders need a written argument.** Default to
  `seq_cst`; weakening to `acquire`/`release`/`relaxed` is a
  reviewable claim about which happens-before edges the algorithm
  needs. An unexplained `memory_order_relaxed` is a finding.
- Multiple mutexes: document the lock order, or take them together
  with `std::scoped_lock`, which is deadlock-free by construction.
- A `std::thread` neither joined nor detached before its destructor
  calls `std::terminate`. Prefer `std::jthread` (C++20).
- `shared_ptr`'s control block is thread-safe; the pointee is not.
  Two threads mutating `*p` still race.

## Exceptions

- State which guarantee each function offers: nothrow, strong, or
  basic. A function that mutates two members must not leave them
  inconsistent when the second mutation throws — do the work on
  copies and commit with non-throwing operations.
- Throwing from a destructor during unwinding calls
  `std::terminate`. Destructors are implicitly `noexcept`.
- A `noexcept` function that can actually throw terminates. Do not
  annotate optimistically.
- Cleanup belongs in destructors, never in a `catch` block.

## Templates, headers, ABI

- **ODR violations**: a non-inline function defined in a header, or
  the same entity defined differently in two TUs. Ill-formed, no
  diagnostic required — it links and then misbehaves.
- `T&&` on a deduced parameter is a forwarding reference and needs
  `std::forward`; on a concrete type it is an rvalue reference and
  needs `std::move`. Mixing them silently copies.
- `std::move` on a `const` object silently copies.
- **`return std::move(local)` defeats NRVO** and is slower than
  `return local;`.
- Prefer concepts (C++20) to SFINAE for constraints — the error
  messages are the deliverable.
- Watch compile-time and binary size on heavily instantiated
  templates; consider extracting the type-independent body.

## Performance patterns worth flagging

- `for (auto x : container)` copies every element. `const auto&`,
  or `auto&&` in generic code.
- Sink parameters by value then `std::move`; large read-only
  parameters by `const&`; cheap types by value.
- `reserve()` before a loop of known size.
- `emplace_back` to construct in place, `push_back` when you already
  have the object.
- `std::endl` flushes the stream every call; use `'\n'`.
- Returning `const` by value blocks moves at the call site.

## Tooling

Check CI actually runs these, since most of the above is invisible to
the compiler alone:

- `-Wall -Wextra -Werror`, and `-Wshadow` where the codebase allows.
- Sanitizers on the test suite: `-fsanitize=address,undefined`, and
  `thread` where concurrency changed. ASan and TSan cannot run
  together.
- clang-tidy with the project's checks; `.clang-format` enforced.
