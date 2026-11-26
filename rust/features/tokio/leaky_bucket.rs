//! Pattern: leaky-bucket rate limiting.
//!
//! The bucket is a bounded `mpsc::channel`: producers `try_send` a request
//! into it (rejected immediately if the bucket is already full), and a
//! single background task "leaks" the bucket by pulling one request off per
//! fixed-interval tick. Rejected alternative: a `tokio::sync::Semaphore`
//! refilled on an `interval` — that models a *token* bucket (burst up to
//! capacity, nothing is ever queued) rather than a leaky bucket, which by
//! definition holds pending requests and drains them at a constant rate.
//!
//! Go comparison: a Go engineer's instinct here is a mutex-guarded
//! `[]Request` queue plus a `sync.Cond` (or a second channel) to wake the
//! drainer — the general-purpose primitives Go always reaches for. Rust
//! doesn't need them: `mpsc::Receiver` has no `Clone` impl, so the compiler
//! — not convention — guarantees at most one task ever owns the read side,
//! which is exactly the "single leak worker" shape this algorithm requires.
//! A Go `chan` gives you no such static guarantee; any number of goroutines
//! can call `<-ch` and nothing stops a second one from being added later.

use anyhow::{Context, Result};
use std::sync::atomic::{AtomicU32, Ordering};
use std::sync::Arc;
use std::time::Duration;
use tokio::sync::mpsc::{self, error::TryRecvError, error::TrySendError};
use tokio::time::{interval_at, sleep, Instant};

const BUCKET_CAPACITY: usize = 3;
const LEAK_INTERVAL: Duration = Duration::from_millis(15);
const PRODUCER_DELAY: Duration = Duration::from_millis(6);
const REQUESTS_PER_PRODUCER: u32 = 4;
const NUM_PRODUCERS: u32 = 3;

#[derive(Debug)]
struct Request {
    producer: u32,
    seq: u32,
}

#[tokio::main]
async fn main() -> Result<()> {
    let (tx, rx) = mpsc::channel::<Request>(BUCKET_CAPACITY);

    // The drain task takes ownership of `rx` outright — it is the bucket's
    // only reader for the rest of the program.
    let drain = tokio::spawn(leak(rx));

    let accepted = Arc::new(AtomicU32::new(0));
    let rejected = Arc::new(AtomicU32::new(0));

    let mut producers = Vec::new();
    for producer in 0..NUM_PRODUCERS {
        let tx = tx.clone(); // each producer owns its own clone of the sender
        let accepted = Arc::clone(&accepted);
        let rejected = Arc::clone(&rejected);
        producers.push(tokio::spawn(async move {
            for seq in 0..REQUESTS_PER_PRODUCER {
                match tx.try_send(Request { producer, seq }) {
                    Ok(()) => {
                        accepted.fetch_add(1, Ordering::Relaxed);
                        println!("producer {producer}: request {seq} admitted to bucket");
                    }
                    Err(TrySendError::Full(_)) => {
                        rejected.fetch_add(1, Ordering::Relaxed);
                        println!("producer {producer}: bucket full, request {seq} rejected");
                    }
                    Err(TrySendError::Closed(_)) => {
                        anyhow::bail!("leak task exited before producer {producer} finished");
                    }
                }
                sleep(PRODUCER_DELAY).await;
            }
            Ok::<(), anyhow::Error>(())
        }));
    }
    drop(tx); // the drain loop below ends once every producer's sender is dropped

    for p in producers {
        p.await??;
    }
    drain.await.context("leak task panicked")?;

    println!(
        "leaky_bucket: accepted = {}, rejected = {}",
        accepted.load(Ordering::Relaxed),
        rejected.load(Ordering::Relaxed)
    );
    Ok(())
}

/// Drains the bucket at a fixed rate, regardless of how bursty the producers
/// filling it are. Owns `rx` exclusively for its whole lifetime.
async fn leak(mut rx: mpsc::Receiver<Request>) {
    // `interval_at(now + period, period)` delays the *first* tick by a full
    // period. Plain `tokio::time::interval(period)` fires its first tick
    // immediately, unlike Go's `time.NewTicker`, which always waits out the
    // full period before its first tick — using `interval_at` here matches
    // the Go behavior a Go engineer would otherwise expect by default.
    let mut ticker = interval_at(Instant::now() + LEAK_INTERVAL, LEAK_INTERVAL);
    loop {
        ticker.tick().await;
        // try_recv, not recv().await: blocking on an empty queue here would
        // stop this loop from polling the ticker while idle, and
        // interval's default missed-tick behavior fires every deadline that
        // elapsed while unpolled back-to-back on the next poll — silently
        // erasing the rate limit for exactly the bursty-after-idle traffic
        // shape a leaky bucket exists to smooth.
        match rx.try_recv() {
            Ok(req) => println!(
                "leak: processed request {} from producer {}",
                req.seq, req.producer
            ),
            Err(TryRecvError::Empty) => {}
            Err(TryRecvError::Disconnected) => {
                println!("leak: bucket drained and closed, shutting down");
                return;
            }
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[tokio::test(start_paused = true)]
    async fn bucket_rejects_when_full_and_admits_again_after_a_leak() {
        let (tx, rx) = mpsc::channel::<Request>(1);
        tokio::spawn(leak(rx));

        tx.try_send(Request { producer: 0, seq: 0 })
            .expect("first request should fit in an empty bucket");

        match tx.try_send(Request { producer: 0, seq: 1 }) {
            Err(TrySendError::Full(_)) => {}
            other => panic!("expected the full bucket to reject, got {other:?}"),
        }

        // Paused virtual time auto-advances to the leak task's next tick
        // deadline once every task is parked waiting on a timer.
        sleep(LEAK_INTERVAL).await;
        tokio::task::yield_now().await;

        tx.try_send(Request { producer: 0, seq: 2 })
            .expect("bucket should have leaked one slot by now");
    }
}
