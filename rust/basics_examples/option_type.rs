//! Pattern: `Option<T>` replacing nil/zero-value ambiguity.
//!
//! Go comparison: a Go function that "might not have a value" typically
//! returns a pointer that might be nil, or relies on the zero value (0,
//! "", false) doubling as "absent" — ambiguous when 0 or "" is also a
//! legitimate value. `Option<T>` makes "no value" a distinct, type-level
//! case: `Some(T)` or `None`, no null pointer involved. The compiler
//! forces you to handle both cases before getting at the `T` inside, so
//! there's no equivalent of a nil-pointer-dereference panic from a
//! forgotten nil check.

use std::collections::HashMap;

fn find_age(ages: &HashMap<&str, u32>, name: &str) -> Option<u32> {
    // HashMap::get returns Option<&u32>; .copied() turns that into
    // Option<u32> since u32 is Copy.
    ages.get(name).copied()
}

fn main() {
    let mut ages = HashMap::new();
    ages.insert("alice", 30);
    ages.insert("bob", 0); // a legitimate age, not "missing"

    for name in ["alice", "bob", "carol"] {
        match find_age(&ages, name) {
            Some(age) => println!("{name} is {age}"),
            None => println!("{name} has no recorded age"),
        }
    }

    // unwrap_or supplies a default without a match, for callers who
    // don't need to distinguish "0" from "missing".
    let carol_age = find_age(&ages, "carol").unwrap_or(0);
    println!("carol_age defaulted to {carol_age}");

    // map transforms the value inside Some without unwrapping it; None
    // passes through untouched.
    let doubled = find_age(&ages, "alice").map(|age| age * 2);
    println!("alice_age doubled = {doubled:?}");
}
