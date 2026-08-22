//! Pattern: ergonomic error handling in async code with `anyhow`.
//!
//! Go comparison: `anyhow::Result<T>` is Go's `(T, error)` collapsed into
//! one type. The `?` operator is Go's `if err != nil { return err }`
//! boilerplate, applied automatically. `.context("...")` is the async-Rust
//! equivalent of `fmt.Errorf("doing X: %w", err)` — it wraps the error with
//! a message while preserving the original as the source, and printing the
//! resulting `anyhow::Error` with `{:?}` shows the whole chain at once.

use anyhow::{bail, Context, Result};
use tokio::time::{sleep, Duration};

#[tokio::main]
async fn main() -> Result<()> {
    match run().await {
        Ok(value) => println!("error_handling: succeeded with {value}"),
        Err(e) => {
            // anyhow's Debug format prints the full "caused by" chain, the
            // equivalent of unwrapping a Go %w chain by hand.
            println!("error_handling: failed:\n{e:?}");
        }
    }
    Ok(())
}

async fn run() -> Result<u32> {
    let raw = fetch_config("missing.conf")
        .await
        .context("loading application config")?;
    Ok(raw)
}

async fn fetch_config(path: &str) -> Result<u32> {
    sleep(Duration::from_millis(5)).await;
    if path == "missing.conf" {
        bail!("config file {path:?} not found");
    }
    Ok(42)
}
