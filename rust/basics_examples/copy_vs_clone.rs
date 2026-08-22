//! Pattern: `Copy` vs `Clone` — implicit bitwise copies vs explicit deep copies.
//!
//! Go comparison: Go doesn't have this as a type-system-checked opt-in;
//! whether an assignment "copies" or "shares" is a property of the
//! value's shape (fixed-size struct/array vs slice/map/pointer header)
//! learned by convention, not enforced or queryable by the compiler.
//! Rust makes it a real trait: `Copy` types (integers, floats, bool,
//! char, and tuples/arrays of Copy types — small, fixed-size, stack-only
//! data) are duplicated implicitly on assignment, with the original still
//! usable. Everything else, including `String` and `Vec<T>` because they
//! own a heap allocation the compiler can't safely duplicate for free,
//! requires an explicit `.clone()` call to get a second, independent copy.

#[derive(Debug, Clone, Copy)]
struct Point {
    x: i32,
    y: i32,
}

fn main() {
    let p1 = Point { x: 1, y: 2 };
    let p2 = p1; // Point is Copy: this duplicates the struct, doesn't move it
    println!("p1=({}, {}) p2=({}, {}), both usable", p1.x, p1.y, p2.x, p2.y);

    let s1 = String::from("heap-allocated");
    let s2 = s1.clone(); // explicit deep copy: separate heap buffers
    println!("s1={s1} s2={s2}, both usable because we cloned");

    let v1 = vec![1, 2, 3];
    let mut v2 = v1.clone();
    v2.push(4);
    println!("v1 untouched by v2's mutation: v1={v1:?} v2={v2:?}");
}
