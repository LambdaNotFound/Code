//! Pattern: structured concurrency — waiting for a fixed or dynamic set of
//! tasks to finish.
//!
//! Go comparison: `tokio::join!` is a `sync.WaitGroup` for a *known, fixed*
//! number of futures, but it also hands back each future's return value (a
//! WaitGroup gives you nothing back — you'd close over shared variables
//! instead). `try_join!` mirrors an `errgroup.Group`: the first `Err` short
//! -circuits the rest. `JoinSet` is for a *dynamic* number of tasks — the
//! "spawn N goroutines, fan results back through a channel" pattern, built
//! in, so you don't wire up that channel yourself.

use anyhow::{bail, Result};
use tokio::task::JoinSet;
use tokio::time::{sleep, Duration};

#[tokio::main]
async fn main() -> Result<()> {
    // join!: wait for two known futures, get both results back.
    let (a, b) = tokio::join!(compute(1), compute(2));
    println!("join!: {a} + {b} = {}", a + b);

    // try_join!: like errgroup — the first error stops the rest.
    match tokio::try_join!(may_fail(1), may_fail(-1)) {
        Ok((x, y)) => println!("try_join!: unreachable ({x}, {y})"),
        Err(e) => println!("try_join!: stopped early on error: {e}"),
    }

    // JoinSet: a dynamic number of tasks, results collected as they finish
    // (not necessarily in spawn order).
    let mut set = JoinSet::new();
    for id in 0..3 {
        set.spawn(async move { compute(id).await });
    }
    let mut results = Vec::new();
    while let Some(res) = set.join_next().await {
        results.push(res?);
    }
    results.sort();
    println!("JoinSet: collected {results:?}");

    Ok(())
}

async fn compute(n: i32) -> i32 {
    sleep(Duration::from_millis(5)).await;
    n * n
}

async fn may_fail(n: i32) -> Result<i32> {
    sleep(Duration::from_millis(5)).await;
    if n < 0 {
        bail!("negative input: {n}");
    }
    Ok(n)
}
