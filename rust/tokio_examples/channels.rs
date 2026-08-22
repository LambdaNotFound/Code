//! Pattern: message passing with Tokio's channel family.
//!
//! Go comparison: Go has one channel type (`chan T`), inherently
//! multi-producer/multi-consumer, and "unbuffered" vs "buffered" is just a
//! choice of capacity. Tokio splits that single type into four
//! purpose-built channels, each with its own ownership shape:
//! - `mpsc` ~ a buffered `chan T` with many senders, one receiver.
//! - `oneshot` ~ the common Go idiom `done := make(chan T, 1)` used to
//!   return exactly one value/signal, made a first-class type (the sender
//!   is consumed by `send`, so the type system — not convention — stops a
//!   second send).
//! - `watch` ~ no direct Go equivalent; closest is a mutex-guarded variable
//!   plus a `sync.Cond` to notify readers of the latest value (only the
//!   newest value matters, not a backlog of every update).
//! - `broadcast` ~ no Go stdlib equivalent; replaces the common Go pattern
//!   of fanning a value out to N per-subscriber channels by hand.

use anyhow::Result;
use tokio::sync::{broadcast, mpsc, oneshot, watch};
use tokio::time::{sleep, Duration};

#[tokio::main]
async fn main() -> Result<()> {
    mpsc_example().await?;
    oneshot_example().await?;
    watch_example().await?;
    broadcast_example().await?;
    Ok(())
}

/// Many producers, one consumer — bounded like a buffered Go channel.
async fn mpsc_example() -> Result<()> {
    // Capacity 4: sends block (are awaited) once 4 messages are buffered,
    // exactly like `make(chan int, 4)` in Go.
    let (tx, mut rx) = mpsc::channel::<i32>(4);

    let mut senders = Vec::new();
    for id in 0..3 {
        let tx = tx.clone(); // each producer owns its own clone of the sender
        senders.push(tokio::spawn(async move {
            tx.send(id)
                .await
                .map_err(|e| anyhow::anyhow!("send failed: {e}"))
        }));
    }
    drop(tx); // the receiver's loop below ends once every sender is dropped

    for s in senders {
        s.await??;
    }

    let mut received = Vec::new();
    while let Some(v) = rx.recv().await {
        received.push(v);
    }
    received.sort();
    println!("mpsc: received {received:?}");
    Ok(())
}

/// Exactly one value, one time — the type-level version of a Go "done" channel.
async fn oneshot_example() -> Result<()> {
    let (tx, rx) = oneshot::channel::<&'static str>();

    let sender = tokio::spawn(async move {
        tx.send("computed result")
            .map_err(|_| anyhow::anyhow!("receiver dropped before send"))
    });

    let value = rx
        .await
        .map_err(|_| anyhow::anyhow!("sender dropped without sending"))?;
    println!("oneshot: got {value}");
    sender.await??;
    Ok(())
}

/// Only the latest value matters; every receiver observes it independently.
async fn watch_example() -> Result<()> {
    let (tx, mut rx) = watch::channel(0u32);

    let sender = tokio::spawn(async move {
        for v in 1..=3 {
            tx.send(v).map_err(|_| anyhow::anyhow!("no receivers left"))?;
            sleep(Duration::from_millis(5)).await;
        }
        Ok::<(), anyhow::Error>(())
    });

    while rx.changed().await.is_ok() {
        let value = *rx.borrow();
        println!("watch: latest value is {value}");
        if value == 3 {
            break;
        }
    }
    sender.await??;
    Ok(())
}

/// Every subscriber gets every message — a built-in fan-out you'd otherwise
/// hand-roll in Go with a slice of per-subscriber channels.
async fn broadcast_example() -> Result<()> {
    let (tx, mut rx1) = broadcast::channel::<&'static str>(8);
    let mut rx2 = tx.subscribe();

    tx.send("event A").map_err(|_| anyhow::anyhow!("no receivers"))?;
    tx.send("event B").map_err(|_| anyhow::anyhow!("no receivers"))?;

    for (label, rx) in [("rx1", &mut rx1), ("rx2", &mut rx2)] {
        while let Ok(event) = rx.try_recv() {
            println!("broadcast: {label} saw {event}");
        }
    }
    Ok(())
}
