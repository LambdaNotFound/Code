//! Pattern: `trait` definitions/impls, static dispatch vs dynamic dispatch.
//!
//! Go comparison: a Go interface is satisfied implicitly — any type with
//! the right methods implements it, with no declaration linking them.
//! Rust traits are explicit: `impl Trait for Type` is required even when
//! the methods already exist, so the compiler (and a reader) always knows
//! which types implement which trait from a grep-able declaration, not
//! from structural inference. Dispatch is also a real choice in Rust,
//! where Go always dispatches through an interface's itable at runtime:
//! a generic function with a trait bound (`fn f<T: Trait>(x: T)`) is
//! monomorphized into a separate compiled copy per concrete type —
//! static dispatch, zero runtime indirection — while `&dyn Trait` builds
//! a vtable and dispatches at runtime, the direct analogue of a Go
//! interface value.

trait Shape {
    fn area(&self) -> f64;
    fn name(&self) -> &str;
}

struct Circle {
    radius: f64,
}

impl Shape for Circle {
    fn area(&self) -> f64 {
        std::f64::consts::PI * self.radius * self.radius
    }
    fn name(&self) -> &str {
        "circle"
    }
}

struct Square {
    side: f64,
}

impl Shape for Square {
    fn area(&self) -> f64 {
        self.side * self.side
    }
    fn name(&self) -> &str {
        "square"
    }
}

// Static dispatch: monomorphized per concrete type at compile time, no
// vtable, callable inline. The `T: Shape` bound is what makes this work.
fn describe_static<T: Shape>(shape: &T) {
    println!("[static] {} has area {:.2}", shape.name(), shape.area());
}

// Dynamic dispatch: one function body shared by every Shape, dispatched
// at runtime through a vtable — the direct analogue of a Go interface
// value's itable dispatch.
fn describe_dynamic(shape: &dyn Shape) {
    println!("[dynamic] {} has area {:.2}", shape.name(), shape.area());
}

fn main() {
    let circle = Circle { radius: 2.0 };
    let square = Square { side: 3.0 };

    describe_static(&circle);
    describe_static(&square);

    // A single Vec can hold different concrete Shape-implementers only
    // through dyn Trait; a generic Vec<T> would force every element to
    // be the *same* T.
    let shapes: Vec<Box<dyn Shape>> = vec![Box::new(circle), Box::new(square)];
    for shape in &shapes {
        describe_dynamic(shape.as_ref());
    }
}
