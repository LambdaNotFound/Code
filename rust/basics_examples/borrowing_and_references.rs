//! Pattern: borrowing — `&T` (shared) and `&mut T` (exclusive) references.
//!
//! Go comparison: Go has no borrow checker. A `*T` pointer can be read and
//! written from anywhere that holds it, concurrently, and it's on the
//! programmer (or `go test -race`) to catch conflicting access. Rust's
//! aliasing rule is enforced at compile time instead: at any point you may
//! have either many `&T` shared references or exactly one `&mut T`
//! exclusive reference to a given value, never both. An invalid program
//! simply doesn't compile; there's nothing to run to observe the bug.

fn main() {
    let mut count = vec![1, 2, 3];

    print_len(&count); // shared borrow, read-only
    push_one(&mut count); // exclusive borrow, allowed to mutate
    print_len(&count);

    // Multiple shared borrows are fine at the same time.
    let a = &count;
    let b = &count;
    println!("two shared borrows: {a:?} and {b:?}");

    // A shared and an exclusive borrow can't coexist. Uncommenting the
    // next two lines is a compile error ("cannot borrow `count` as mutable
    // because it is also borrowed as immutable"):
    // let r = &count;
    // push_one(&mut count);
    // println!("{r:?}");

    let title = String::from("a long title");
    let short = String::from("short");
    println!("longest: {}", longest(&title, &short));
}

fn print_len(v: &[i32]) {
    println!("length is {}", v.len());
}

fn push_one(v: &mut Vec<i32>) {
    v.push(v.len() as i32 + 1);
}

// The `'a` here is a lifetime annotation, not a borrow-checker escape
// hatch: it tells the compiler that the returned reference is valid for
// as long as *both* `a` and `b` are. Full lifetime elision rules (when
// you can omit `'a` entirely) are out of scope for this 101 — briefly,
// the compiler infers it automatically for common single-input-reference
// signatures, but a signature like this one, with two input references
// tied to one output reference, needs it spelled out explicitly.
fn longest<'a>(a: &'a str, b: &'a str) -> &'a str {
    if a.len() >= b.len() { a } else { b }
}
