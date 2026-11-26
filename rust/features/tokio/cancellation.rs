//! Pattern: cooperative cancellation across tasks.
//!
//! Go comparison: `CancellationToken` (from `tokio-util`) is the async-Rust
//! analogue of a `context.Context` used purely for cancellation:
//! `token.cancel()` ~ calling a context's cancel func, `token.cancelled()`
//! ~ selecting on `ctx.Done()`, and `token.child_token()` ~
//! `context.WithCancel(parent)` (cancelling a parent cancels every child;
//! cancelling a child does not affect the parent). This is chosen over
//! hand-rolling cancellation with a `watch::channel<bool>` because it
//! already gives you the parent/child tree Go engineers expect from
//! `context.Context` — reimplementing that on top of `watch` would just
//! rebuild what tokio-util already provides.
//!
//! As in Go, cancellation here is cooperative: a task must poll
//! `token.cancelled()` (or `is_cancelled()`) and choose to stop — nothing
//! forcibly kills it, exactly as an ignored `ctx.Done()` never stops a
//! goroutine on its own.

use anyhow::Result;
use tokio::time::{sleep, Duration};
use tokio_util::sync::CancellationToken;

#[tokio::main]
async fn main() -> Result<()> {
    let parent = CancellationToken::new();
    let child = parent.child_token();

    let worker = tokio::spawn(run_worker(child.clone(), "worker-1"));

    sleep(Duration::from_millis(35)).await;
    println!("main: cancelling parent token");
    parent.cancel(); // cancels every child token too

    worker.await??;
    println!(
        "main: worker confirmed shutdown, child cancelled = {}",
        child.is_cancelled()
    );
    Ok(())
}

async fn run_worker(token: CancellationToken, name: &str) -> Result<()> {
    let mut tick = 0;
    loop {
        tokio::select! {
            // Mirrors `select { case <-ctx.Done(): ...; case <-ticker.C: ... }`
            _ = token.cancelled() => {
                println!("{name}: cancellation observed, shutting down");
                return Ok(());
            }
            _ = sleep(Duration::from_millis(10)) => {
                tick += 1;
                println!("{name}: tick {tick}");
            }
        }
    }
}
