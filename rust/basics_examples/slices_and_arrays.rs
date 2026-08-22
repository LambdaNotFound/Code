//! Pattern: fixed-size arrays `[T; N]` vs slices `&[T]`.
//!
//! Go comparison: Go's usual `[]T` is a three-word header (pointer, len,
//! cap) over a backing array, plus a fixed-size `[N]T` array type that's
//! rarely used directly. Rust splits the same idea into `[T; N]` (size is
//! part of the type, checked at compile time, lives on the stack) and
//! `&[T]` (a borrowed view: pointer + len, no cap, no ownership). The
//! aliasing story differs too: Go's `append` may or may not reallocate
//! depending on spare capacity, so two Go slices sharing a backing array
//! can silently alias or silently diverge at runtime depending on that
//! capacity. Rust's growable equivalent is `Vec<T>` (see
//! `vec_and_hashmap.rs`); `&[T]` itself never grows — it's a view only,
//! there's no `append` on a slice at all, and the borrow checker forbids
//! any mutation that could invalidate an outstanding view.

fn sum(values: &[i32]) -> i32 {
    values.iter().sum()
}

fn main() {
    let fixed: [i32; 4] = [10, 20, 30, 40]; // the length 4 is part of the type
    println!("array: {fixed:?}, sum = {}", sum(&fixed)); // &[i32; 4] coerces to &[i32]

    let middle: &[i32] = &fixed[1..3]; // a borrowed view into the same array
    println!("middle slice = {middle:?}");

    let mut heap_vec = vec![1, 2, 3, 4, 5];
    heap_vec.push(6); // fine: no outstanding borrow yet
    let view: &[i32] = &heap_vec[..3];
    println!("vec = {heap_vec:?}, first-3 view = {view:?}, sum(view) = {}", sum(view));

    // `view` borrows heap_vec's backing storage. Uncommenting the next
    // line is a compile error ("cannot borrow `heap_vec` as mutable
    // because it is also borrowed as immutable"): a push that grows past
    // capacity would reallocate and move the backing buffer, which would
    // leave `view` pointing at freed memory. The borrow checker forbids
    // that outright, instead of letting it happen the way a stale Go
    // slice header silently can.
    // heap_vec.push(7);
    println!("view still valid: {view:?}");
}
