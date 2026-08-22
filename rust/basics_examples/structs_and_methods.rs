//! Pattern: structs, `impl` blocks, and method receivers.
//!
//! Go comparison: a Rust `struct` plus its `impl` block is Go's struct
//! plus its methods, but Go spells the receiver's ownership as
//! `(c Counter)` vs `(c *Counter)` — value vs pointer — while Rust spells
//! the same three-way choice explicitly on `self`: `&self` (read-only
//! borrow, ~ Go's value receiver used for reading), `&mut self`
//! (exclusive borrow, ~ Go's pointer receiver), and `self` (takes
//! ownership, consumes the value — Go has no equivalent, since a Go value
//! receiver never destroys the caller's copy).

struct Counter {
    value: i32,
    label: String,
}

impl Counter {
    // No receiver: an associated function, Go's equivalent of a
    // package-level constructor like `NewCounter(...)`.
    fn new(label: &str) -> Counter {
        Counter {
            value: 0,
            label: label.to_string(),
        }
    }

    // &self: read-only borrow, cannot mutate `value`.
    fn value(&self) -> i32 {
        self.value
    }

    // &mut self: exclusive borrow, can mutate in place.
    fn increment(&mut self) {
        self.value += 1;
    }

    // self: takes ownership, consumes the Counter. Useful for a "final"
    // transformation where the original shouldn't be usable afterward.
    fn into_label(self) -> String {
        self.label
    }
}

fn main() {
    let mut c = Counter::new("requests");
    c.increment();
    c.increment();
    println!("{} = {}", c.label, c.value());

    let label = c.into_label();
    // `c` was moved into `into_label`; it no longer exists here.
    // println!("{}", c.value()); // would not compile: value used after move
    println!("label reclaimed: {label}");
}
