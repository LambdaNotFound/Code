//! Pattern: `String` (owned, heap) vs `&str` (borrowed view), and UTF-8 storage.
//!
//! Go comparison: Go's `string` is already an immutable, UTF-8-oriented
//! byte-slice-like value, and `s[i]` gives you the i-th *byte*, which Go
//! lets you do even though it can land you in the middle of a multi-byte
//! rune. Rust splits "owned, growable" (`String`) from "borrowed view"
//! (`&str`, the `str` equivalent of `&[u8]` specialized to guaranteed-
//! valid UTF-8) the same way it splits `Vec<T>`/`&[T]`. The indexing
//! difference is sharper than Go's, though: Rust doesn't let `s[i]`
//! compile *at all* for a `String`/`str` — byte indexing by integer isn't
//! implemented for that type — precisely because slicing at a
//! non-character boundary would produce invalid UTF-8, and Rust treats
//! "this `str` might not be valid UTF-8" as unacceptable rather than as a
//! foot-gun you opt into.

fn main() {
    let owned: String = String::from("héllo"); // owned, heap-allocated
    let borrowed: &str = &owned; // a view into owned's buffer, no copy
    println!("owned = {owned}, borrowed = {borrowed}");

    // len() counts bytes, not characters: 'é' is 2 bytes in UTF-8, so
    // byte length (6) and char count (5) diverge for this string.
    println!("byte length = {}, char count = {}", owned.len(), owned.chars().count());

    // owned[0] does not compile: `String: Index<{integer}>` isn't
    // implemented at all, so there's no way to accidentally split a
    // multi-byte character the way Go's s[i] lets you.
    // let _ = owned[0];

    // .chars() walks Unicode scalar values; .bytes() walks raw bytes.
    for (i, c) in owned.chars().enumerate() {
        print!("[{i}]={c} ");
    }
    println!();

    // `+` takes the left side by value (it's consumed) and the right
    // side by reference — another place ownership shows up.
    let greeting = String::from("Hello, ") + borrowed;
    println!("{greeting}");
}
