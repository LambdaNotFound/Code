//! Pattern: generic functions/types constrained by trait bounds.
//!
//! Go comparison: Go 1.18+ generics use a constraint interface (e.g.
//! `constraints.Ordered`, or a custom interface listing allowed types)
//! written in a `[T Constraint]` type parameter list. Rust's trait
//! bounds (`<T: Trait>`) are the same idea — restrict a type parameter to
//! types implementing some trait — but bounds compose more directly:
//! multiple bounds combine with `+` (`T: Ord + Clone`) and can move into
//! a separate `where` clause for readability, without needing a named
//! constraint interface declared up front the way Go typically wants.

use std::fmt::Display;

fn largest<T: PartialOrd + Copy>(items: &[T]) -> T {
    let mut max = items[0];
    for &item in items.iter() {
        if item > max {
            max = item;
        }
    }
    max
}

// Multiple bounds via `where`, equivalent to `<T: Display + Clone>` but
// easier to read once the bound list grows.
fn describe<T>(label: &str, value: T) -> String
where
    T: Display + Clone,
{
    let copy = value.clone();
    format!("{label}: {copy}")
}

struct Pair<T> {
    first: T,
    second: T,
}

impl<T: PartialOrd + Display> Pair<T> {
    fn cmp_display(&self) {
        if self.first >= self.second {
            println!("largest is first = {}", self.first);
        } else {
            println!("largest is second = {}", self.second);
        }
    }
}

fn main() {
    let numbers = [34, 50, 25, 100, 65];
    println!("largest number: {}", largest(&numbers));

    let chars = ['y', 'm', 'a', 'q'];
    println!("largest char: {}", largest(&chars));

    println!("{}", describe("count", 42));

    let pair = Pair { first: 5, second: 10 };
    pair.cmp_display();
}
