//! Pattern: offloading CPU-bound or blocking work with `spawn_blocking`.
//!
//! Go comparison: Go rarely needs an equivalent because a blocking syscall
//! or a CPU-bound loop inside a goroutine doesn't stall the whole
//! scheduler — the Go runtime detects blocked Ms and spins up more OS
//! threads for the remaining goroutines automatically. Tokio's async
//! worker threads are cooperatively scheduled: one task that blocks the
//! thread (a synchronous file read, `std::thread::sleep`, a tight CPU loop
//! with no `.await`) freezes every other task on that worker.
//! `spawn_blocking` moves the work onto a separate blocking-thread pool
//! sized for exactly this, keeping the async worker threads free.

use anyhow::Result;
use std::time::Duration;

#[tokio::main]
async fn main() -> Result<()> {
    let cpu_result = tokio::task::spawn_blocking(|| fibonacci(30)).await?;
    println!("spawn_blocking: fibonacci(30) = {cpu_result}");

    let io_result = tokio::task::spawn_blocking(|| {
        // A synchronous, blocking sleep — stands in for a blocking syscall
        // (e.g. a non-async file or DNS call) that has no async version.
        std::thread::sleep(Duration::from_millis(20));
        "blocking I/O finished"
    })
    .await?;
    println!("spawn_blocking: {io_result}");

    Ok(())
}

fn fibonacci(n: u64) -> u64 {
    match n {
        0 => 0,
        1 => 1,
        _ => fibonacci(n - 1) + fibonacci(n - 2),
    }
}
