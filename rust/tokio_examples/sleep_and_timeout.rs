//! Pattern: delaying and bounding async work with `sleep` and `timeout`.
//!
//! Go comparison: `tokio::time::sleep` is `time.Sleep`. `tokio::time::timeout`
//! plays the role Go fills with `context.WithTimeout` plus a `select` on
//! `ctx.Done()` — instead of threading a context through every call and
//! checking `Done()` yourself, you wrap the future once and Tokio races it
//! against a timer for you.

use anyhow::{bail, Result};
use tokio::time::{sleep, timeout, Duration};

#[tokio::main]
async fn main() -> Result<()> {
    sleep(Duration::from_millis(50)).await;
    println!("slept for 50ms");

    match timeout(Duration::from_millis(50), slow_operation(10)).await {
        Ok(value) => println!("within budget: {value}"),
        Err(_) => println!("unreachable"),
    }

    match timeout(Duration::from_millis(10), slow_operation(50)).await {
        Ok(value) => println!("unreachable: {value}"),
        Err(_) => println!("slow_operation timed out, as expected"),
    }

    if let Err(e) = fallible_with_timeout().await {
        println!("propagated timeout error via anyhow: {e}");
    }

    Ok(())
}

async fn slow_operation(delay_ms: u64) -> u32 {
    sleep(Duration::from_millis(delay_ms)).await;
    42
}

async fn fallible_with_timeout() -> Result<u32> {
    match timeout(Duration::from_millis(10), slow_operation(50)).await {
        Ok(value) => Ok(value),
        Err(_) => bail!("slow_operation did not complete in time"),
    }
}
