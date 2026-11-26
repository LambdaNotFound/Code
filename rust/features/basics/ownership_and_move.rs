//! Pattern: move semantics — ownership transfers instead of copying.
//!
//! Go comparison: assigning or passing a Go string/slice/map copies the
//! header (pointer+len[+cap]) and both bindings alias the same backing
//! array; nothing stops you from using both afterward. Rust's `String`
//! has no implicit "copy the header, share the backing buffer" mode: an
//! assignment or a by-value function call *moves* ownership, and the
//! compiler statically forbids using the old binding afterward — there's
//! no runtime aliasing to worry about because there's only one owner.

fn main() {
    let x = 5;
    let y = x; // i32 is Copy, so this copies bits instead of moving
    println!("x={x} y={y}, both usable because i32 is Copy");

    let original = String::from("hello");
    let moved = original; // ownership of the heap buffer transfers here

    // `original` is no longer usable. Uncommenting the next line is a
    // compile error: "value borrowed after move" / "value used after move".
    // println!("{original}");
    println!("moved now owns the string: {moved}");

    let owned = String::from("consume me");
    takes_ownership(owned);
    // `owned` was moved into the function; it's gone once the call
    // returns, whether or not the function does anything with it.
    // println!("{owned}"); // would not compile

    let via_clone = String::from("i get cloned");
    let cloned = via_clone.clone();
    println!("original survives because we cloned: {via_clone} / {cloned}");
}

fn takes_ownership(s: String) {
    println!("took ownership of: {s}");
}
