//! Pattern: `Result<T, E>` as a type-level `(T, error)`.
//!
//! Go comparison: `Result<T, E>` is Go's `(T, error)` return pair with the
//! "did it fail" case folded into the type instead of a second return
//! value you might forget to check. `match` on a `Result` is exhaustive
//! the same way it is on `Option`, so there's no equivalent of forgetting
//! `if err != nil`. This file stays at the type level; the `?` operator,
//! `anyhow::Result`, and `.context(...)` error-wrapping idioms used
//! throughout the async examples are covered in
//! `docs/research/rust-tokio.md`'s "Error handling in async code" section
//! rather than repeated here.

#[derive(Debug)]
struct ParseError {
    input: String,
}

// Display + Error is the minimum a real error type needs: Display for a
// human-readable message, Error (a marker beyond Display + Debug) for
// interop with anything that expects `dyn std::error::Error`, including
// `anyhow`'s `?` in the async examples this file deliberately doesn't repeat.
impl std::fmt::Display for ParseError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "{:?} is not a valid even number", self.input)
    }
}

impl std::error::Error for ParseError {}

fn parse_even(input: &str) -> Result<i32, ParseError> {
    // `?` here is the same operator covered in docs/research/rust-tokio.md:
    // on Err it returns early with that error converted into this
    // function's Err type; on Ok it unwraps and keeps going.
    let n: i32 = input
        .parse()
        .map_err(|_| ParseError { input: input.to_string() })?;
    if n % 2 != 0 {
        return Err(ParseError { input: input.to_string() });
    }
    Ok(n)
}

fn main() {
    for input in ["4", "7", "not a number"] {
        match parse_even(input) {
            Ok(n) => println!("{input:?} parsed as even number {n}"),
            Err(e) => println!("{input:?} rejected, saw input {:?}", e.input),
        }
    }
}
