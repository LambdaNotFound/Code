//! Pattern: cooperative scheduling and explicit yield points.
//!
//! Go comparison: since Go 1.14, the Go scheduler can preempt a goroutine
//! that's been running too long even without a function call or channel
//! op — a tight `for {}` loop will not permanently starve other
//! goroutines. Tokio has no such preemption: a task only ever yields back
//! to the scheduler at an `.await` point. A task that runs a long
//! CPU-bound loop without awaiting will hog its worker thread and starve
//! every other task scheduled onto it.

use anyhow::Result;
use std::time::Instant;
use tokio::task;
use tokio::time::{sleep, Duration};

#[tokio::main(flavor = "current_thread")]
async fn main() -> Result<()> {
    // A "ticker" task that interleaves with `hog` below via cooperative
    // yielding rather than being blocked behind it entirely. In an
    // unoptimized (debug) build, the unyielded work between each
    // `yield_now()` call below is slow enough that tick spacing stretches
    // well past the nominal 10ms; run with `--release` to see interleaving
    // closer to that cadence.
    let ticker = tokio::spawn(async {
        for i in 1..=3 {
            sleep(Duration::from_millis(10)).await;
            println!("tick {i} at {:?}", Instant::now());
        }
    });

    // A CPU-bound task that yields periodically: on this single-threaded
    // runtime, removing the `yield_now().await` call below would starve
    // `ticker` completely until `hog` finishes, since nothing would ever
    // hand control back to the scheduler.
    let hog = tokio::spawn(async {
        let mut sum: u64 = 0;
        for i in 0..200_000_000u64 {
            sum = sum.wrapping_add(i);
            if i % 50_000_000 == 0 {
                task::yield_now().await;
            }
        }
        sum
    });

    let (ticker_result, sum) = tokio::join!(ticker, hog);
    ticker_result?;
    println!("hog finished, sum = {}", sum?);
    Ok(())
}
