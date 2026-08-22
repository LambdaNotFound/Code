//! Pattern: sum types via `enum`, and exhaustive `match`.
//!
//! Go comparison: Go has no sum type. The usual workarounds are an
//! `interface{}` with a type switch (unchecked at compile time — a new
//! implementation of the interface can appear anywhere, and the type
//! switch's `default:` case is easy to forget), or a struct with a
//! `Kind` enum-ish const plus a pile of possibly-unused fields, only some
//! of which are valid depending on `Kind`. Rust's `enum` variants can
//! each carry their own different data, and `match` on an enum is
//! exhaustive: the compiler rejects the match if a variant is missing,
//! at compile time, not by an easily-forgotten `default:`.

enum Shape {
    Circle { radius: f64 },
    Rectangle { width: f64, height: f64 },
    Triangle { base: f64, height: f64 },
}

fn area(shape: &Shape) -> f64 {
    match shape {
        Shape::Circle { radius } => std::f64::consts::PI * radius * radius,
        Shape::Rectangle { width, height } => width * height,
        Shape::Triangle { base, height } => 0.5 * base * height,
        // Deleting one of the arms above is a compile error:
        // "non-exhaustive patterns: `&Shape::Triangle { .. }` not covered".
    }
}

fn main() {
    let shapes = vec![
        Shape::Circle { radius: 2.0 },
        Shape::Rectangle { width: 3.0, height: 4.0 },
        Shape::Triangle { base: 5.0, height: 6.0 },
    ];

    for shape in &shapes {
        println!("area = {:.2}", area(shape));
    }
}
