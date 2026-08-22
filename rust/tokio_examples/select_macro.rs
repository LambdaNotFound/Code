//! Pattern: racing multiple async operations with `tokio::select!`.
//!
//! Go comparison: this is a direct analogue of Go's `select` statement over
//! channel operations. The two differences that trip up Go engineers:
//! `select!` branches can carry an `if` guard (skipping a branch without a
//! separate check inside it), and by default the futures that *aren't*
//! chosen are dropped/cancelled rather than left running — there's no
//! implicit "the other side keeps going" behavior to reason about.

use tokio::sync::mpsc;
use tokio::time::{sleep, Duration};

#[tokio::main]
async fn main() {
    let (tx, mut rx) = mpsc::channel::<&'static str>(1);

    tokio::spawn(async move {
        sleep(Duration::from_millis(20)).await;
        let _ = tx.send("work done").await;
    });

    // Race the channel receive against a timeout — like
    // `select { case v := <-rx: ... case <-time.After(d): ... }` in Go.
    tokio::select! {
        Some(msg) = rx.recv() => println!("select: received {msg}"),
        _ = sleep(Duration::from_millis(50)) => println!("select: timed out"),
    }

    // `biased` plus a guard: Go's select has neither. `biased` disables the
    // default random branch ordering (useful when one branch should always
    // be checked first); the guard is evaluated before a branch is even
    // considered ready, unlike an `if` written inside the branch body.
    let allow_second_branch = false;
    tokio::select! {
        biased;
        _ = sleep(Duration::from_millis(1)) => println!("select: branch one fired"),
        _ = sleep(Duration::from_millis(1)), if allow_second_branch => {
            println!("select: unreachable, guard was false");
        }
    }
}
