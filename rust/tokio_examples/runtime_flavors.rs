//! Pattern: choosing and constructing the Tokio runtime.
//!
//! Go comparison: the Go runtime always multiplexes goroutines across OS
//! threads for you (tuned by GOMAXPROCS) — there is no equivalent choice to
//! make. Tokio makes the scheduler an explicit decision: `#[tokio::main]`
//! is sugar for building a multi-threaded `Runtime`, but you can also build
//! a `current_thread` runtime for single-threaded execution (e.g. tests, or
//! CPU-light I/O-bound tools where thread-per-core overhead isn't worth it).

use anyhow::Result;
use std::thread;
use tokio::runtime::{Builder, Runtime};

fn main() -> Result<()> {
    // Equivalent to what #[tokio::main] generates: a multi-threaded runtime
    // with one worker thread per CPU.
    let multi = Runtime::new()?;
    multi.block_on(report("multi-threaded runtime"));

    // A single-threaded runtime: every task runs on the thread that calls
    // block_on. Useful for deterministic scheduling, or when embedding
    // Tokio inside an existing thread that shouldn't spawn more of its own.
    let current = Builder::new_current_thread().enable_all().build()?;
    current.block_on(report("current-thread runtime"));

    Ok(())
}

async fn report(label: &str) {
    println!("{label}: running on {:?}", thread::current().id());
}
