# Go review reference

Load when the diff contains Go. Ordered by how often the item is the
actual bug. Each item names what to look for, why it breaks, and the
fix.

## Check the language version before you flag anything

Read the `go` directive in `go.mod` first. Several of the most
commonly cited Go review rules are now wrong, and repeating them
costs the review its credibility.

- **Loop variable capture is fixed as of Go 1.22.** The spec: "Each
  iteration has its own separate declared variable (or variables)
  [Go 1.22] ... Prior to Go 1.22, iterations share one set of
  variables." So `for _, v := range xs { go func(){ use(v) }() }` is
  correct on a module declaring 1.22 or later, and the `v := v`
  shim is redundant, not required. Flag the capture only when the
  module predates 1.22.
- **Timers are garbage collected as of Go 1.23.** `time.After` in a
  loop no longer retains the timer until it fires; the stdlib doc
  now says "There is no reason to prefer NewTimer when After will
  do." Stop calling it a leak on 1.23+.
- **A timer's channel became unbuffered in Go 1.23**, which creates
  a real migration bug in the other direction: the old drain idiom
  `if !t.Stop() { <-t.C }` could block forever now that no stale
  value is buffered. Flag *that* on code moving to 1.23+.

## Concurrency

- **Every `go` statement needs an answer to "what stops this?"**
  A goroutine blocked on a channel send with no live receiver, or a
  receive with no live sender, leaks for the life of the process.
  The common shape: a worker sends a result after the caller already
  returned on timeout. Give the result channel a buffer of 1 so the
  send never blocks, or make the worker select on `ctx.Done()`.
- **`ctx` is honored or it is decoration.** Long loops and every
  blocking select need a `case <-ctx.Done()`. `context.Context` is
  the first parameter, never a struct field. Every
  `context.WithCancel`/`WithTimeout` needs `defer cancel()` — `go
  vet`'s `lostcancel` catches the miss.
- **`wg.Add` before `go`, never inside it.** Inside, the goroutine
  may not have run before `Wait` returns. `defer wg.Done()` first
  line of the goroutine.
- **Only the sender closes a channel**, and only one of them. Send
  on a closed channel and double close both panic. With multiple
  senders, close a separate `done` channel, or gate with
  `sync.Once`.
- **Concurrent map access is a fatal runtime error**, not a
  recoverable race — the runtime kills the process on detecting
  concurrent read and write. Guard with `sync.RWMutex`. `sync.Map`
  is for two narrow cases: keys written once and read many, or
  goroutines touching disjoint key sets. Otherwise it is slower.
- **Copying a lock.** A struct containing `sync.Mutex` passed or
  returned by value copies the mutex and silently loses mutual
  exclusion. Methods that lock need pointer receivers. `go vet`
  catches it.
- `sync.RWMutex` is not reentrant. A goroutine taking `RLock` twice
  can deadlock when a writer queues between the two.
- Prefer the `sync/atomic` types (`atomic.Int64`, `atomic.Pointer[T]`)
  over the bare functions — they cannot be accidentally accessed
  non-atomically.
- Check CI runs the tests with `-race`. Without it, none of the
  above is caught automatically.

## Errors

- **A discarded error is a decision that needs to be visible.**
  `_ = f()` needs a comment saying why. `defer f.Close()` on a
  *writable* file discards the error that reports failed flush —
  that is data loss. For writers, capture it:
  `defer func() { err = errors.Join(err, f.Close()) }()` with a
  named return.
- **`%w` versus `%v`.** `%w` keeps the chain inspectable;
  `%v` flattens it to text. Then compare with `errors.Is` and
  `errors.As`, never `==`, a type assertion, or
  `strings.Contains(err.Error(), ...)` — all three break the moment
  someone wraps.
- **Typed nil in an error interface.** Returning a nil `*MyError`
  as `error` produces a non-nil interface, so `if err != nil` fires
  on success. Return a bare `nil` on the success path; never
  declare the result variable as the concrete pointer type.
- Error strings are lowercase and unpunctuated, because they get
  wrapped. Wrapping adds the layer's context, not the word
  "failed" — chained "failed to failed to" is the tell.
- Panics belong in `main` and in genuinely unreachable states.
  `recover` only helps in the goroutine that panicked; it does not
  cross goroutine boundaries, so a panic in a spawned worker takes
  the process down regardless of the caller's recover.

## nil, interfaces, slices, maps

- A nil map reads fine and panics on write. A nil slice appends
  fine. Know which one the zero value gives you.
- nil slice and empty slice differ where it is visible: JSON
  marshals nil as `null` and empty as `[]`. That is an API contract.
- Type assertion without the `, ok` form panics. Use the two-value
  form anywhere the type is not guaranteed by construction.
- **Slice aliasing is the quiet one.** Slicing shares the backing
  array, and `append` mutates it in place when capacity allows:
  `b := a[:2]; b = append(b, x)` overwrites `a[2]`. Force a copy
  with the three-index form `a[:2:2]`, which caps capacity so
  append must reallocate.
- Sub-slicing a large buffer keeps the *whole* array alive. To hold
  a small piece of a big read, `copy` it into a right-sized slice.
- Map iteration order is randomized per run. Any test or output
  that depends on it is a flake — sort the keys.

## defer

- `defer` fires at function exit, not block exit. A `defer` inside a
  loop accumulates until the function returns; opening files in a
  loop that way exhausts descriptors. Move the body into its own
  function.
- **Arguments are evaluated at the `defer` statement**, not when it
  runs. `defer log(err)` captures the current `err`, usually nil.
  Wrap in a closure to read the final value.
- A deferred closure can modify named return values. That is the
  mechanism for the `Close` pattern above, and a trap everywhere
  else.

## API and idiom

- Accept interfaces, return concrete types. Define the interface in
  the consuming package, sized to what that consumer calls — not a
  twelve-method interface exported next to its only implementation.
- Take `io.Reader`/`io.Writer` rather than `[]byte` or a filename
  wherever the data could stream or the caller could be a test.
- Make the zero value useful, the way `sync.Mutex` and
  `bytes.Buffer` do, so callers need no constructor.
- Exported identifiers carry a doc comment starting with the
  identifier's name. Do not export what the package does not need
  to expose.
- Receiver types are consistent across a type's method set: mixing
  value and pointer receivers on one type is a finding unless
  there is a stated reason.
- Reach for generics when they remove a type assertion or a
  duplicated function body, not to parameterize something used once.

## Performance worth flagging

- `strings.Builder` instead of `+=` in a loop, which is quadratic.
- Preallocate when the size is known: `make([]T, 0, n)`,
  `make(map[K]V, n)`.
- `[]byte(s)` and `string(b)` each copy. Repeated conversion in a
  loop is a real cost; hoist it.
- `strconv.Itoa` over `fmt.Sprintf("%d", n)` on a hot path.
- Do not claim an allocation or a speedup without a benchmark.
  `go test -bench . -benchmem` produces the number, or the claim is
  labeled **inferred**.

## Tests and tooling

- Table-driven tests for anything with more than one interesting
  input. Each case named, so a failure identifies itself.
- `t.Helper()` in assertion helpers, so failures report the caller's
  line. `t.Cleanup` over `defer` in helpers, since it runs after
  parallel subtests finish.
- `t.Parallel()` changes when subtest bodies run; verify shared
  fixtures are safe before adding it.
- Tests assert behavior through the exported surface. A test that
  reaches into unexported state fails on every refactor and proves
  nothing about the contract.
- Check CI for `go vet`, `-race`, and a linter; add `go test -fuzz`
  for anything parsing untrusted input.

## This repository's conventions

For PRs in this repo, `CLAUDE.md` is binding and outranks the
general guidance above: testify `assert`, tests co-located in the
same package as the implementation, the dot import of
`gocode/golang/types`, `/** ... */` block comments with a `* `
prefix on every line, and the LeetCode rules — signatures followed
exactly, no added input validation, runtime preferred over memory
where nothing says otherwise.
