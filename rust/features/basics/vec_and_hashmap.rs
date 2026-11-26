//! Pattern: `Vec<T>` / `HashMap<K, V>` vs Go's slice/map, and generics.
//!
//! Go comparison: Go's `[]T` and `map[K]V` are built-in types the
//! language special-cases; you can't write your own generic container
//! that looks like them without Go 1.18+ user-defined generics. Rust has
//! no such special-casing: `Vec<T>` and `HashMap<K, V>` are ordinary
//! library types built with the exact same generics (`<T>`, trait bounds)
//! that any of your own types can use — see `generics_and_trait_bounds.rs`.
//! Both Go and Rust randomize map/HashMap iteration order, specifically
//! to prevent HashDoS attacks, but at different granularity: Go
//! re-randomizes on every single `range`, even over the same live map
//! twice in a row; Rust's default hasher picks one random seed when the
//! `HashMap` is constructed, so iterating the same instance twice in a
//! row yields the same order both times, while a fresh instance (or a
//! fresh process run) gets a different seed. Don't rely on order from
//! either.

use std::collections::HashMap;

// Deliberately Vec::new() + push (the Go append-style growth pattern),
// not the vec![] literal clippy would otherwise suggest here.
#[allow(clippy::vec_init_then_push)]
fn main() {
    let mut names: Vec<String> = Vec::new();
    names.push(String::from("alice"));
    names.push(String::from("bob"));
    names.push(String::from("carol"));
    println!("vec: {names:?}, len = {}", names.len());

    let mut scores: HashMap<String, u32> = HashMap::new();
    for name in &names {
        scores.insert(name.clone(), name.len() as u32);
    }

    // entry() is the idiomatic get-or-insert, replacing Go's
    // "v, ok := m[k]; if !ok { m[k] = ... }" pattern.
    *scores.entry(String::from("dave")).or_insert(0) += 1;
    println!("dave's score after entry/or_insert: {}", scores["dave"]);

    let mut sorted_keys: Vec<&String> = scores.keys().collect();
    sorted_keys.sort();
    for key in sorted_keys {
        println!("{key}: {}", scores[key]);
    }
}
