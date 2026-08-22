//! Pattern: spawning tasks and joining on the result.
//!
//! Go comparison: `tokio::spawn` is the closest thing to `go func(){ ... }()`,
//! and the returned `JoinHandle` is roughly a single-result channel you'd
//! build by hand in Go to get a value back out of a goroutine. The key
//! divergence is panics: an unrecovered panic in a goroutine takes down the
//! whole Go process. A panic inside a spawned Tokio task is caught by the
//! runtime and surfaced as `Err(JoinError)` from `.await`ing the handle — it
//! does not crash sibling tasks or the process.

use anyhow::Result;

#[tokio::main]
async fn main() -> Result<()> {
    let ok = tokio::spawn(async { 2 + 2 });
    println!("ok task result: {}", ok.await?);

    let panics = tokio::spawn(async {
        panic!("boom");
    });
    match panics.await {
        Ok(()) => println!("unreachable"),
        Err(join_err) => println!("panicking task reported as: {join_err}"),
    }

    println!("main task is still alive after the sibling panicked");
    Ok(())
}
