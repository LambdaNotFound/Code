//! Pattern: sharing mutable state across tasks.
//!
//! Go comparison: `Arc<Mutex<T>>` is Go's `sync.Mutex` guarding shared
//! data, made explicit in the type — the compiler refuses to hand you the
//! data without going through the lock, where Go relies on convention.
//! `Arc<RwLock<T>>` mirrors `sync.RWMutex` the same way. The wrinkle Go has
//! no equivalent to: Tokio ships *two* Mutex types, and reaching for the
//! wrong one is a classic new-Rustacean-in-async trap.

use anyhow::Result;
use std::sync::{Arc, Mutex as StdMutex, RwLock as StdRwLock};
use tokio::sync::Mutex as TokioMutex;
use tokio::time::{sleep, Duration};

#[tokio::main]
async fn main() -> Result<()> {
    std_mutex_for_sync_sections().await?;
    rwlock_multiple_readers().await?;
    tokio_mutex_for_async_sections().await?;
    Ok(())
}

/// Prefer `std::sync::Mutex` when the critical section is plain, synchronous
/// code with no `.await` inside it — it's cheaper than the async version
/// and, like Go's `sync.Mutex`, only ever blocks a thread, never a task.
async fn std_mutex_for_sync_sections() -> Result<()> {
    let counter = Arc::new(StdMutex::new(0));

    let mut handles = Vec::new();
    for _ in 0..5 {
        let counter = Arc::clone(&counter);
        handles.push(tokio::spawn(async move {
            // Lock, mutate, unlock — all synchronous. The guard is dropped
            // before this task's next await point (there isn't one here),
            // so it never crosses one.
            let mut guard = counter
                .lock()
                .map_err(|_| anyhow::anyhow!("mutex poisoned"))?;
            *guard += 1;
            Ok::<(), anyhow::Error>(())
        }));
    }
    for h in handles {
        h.await??;
    }

    let final_count = *counter
        .lock()
        .map_err(|_| anyhow::anyhow!("mutex poisoned"))?;
    println!("std_mutex: final count = {final_count}");

    // THE PITFALL: don't hold a std::sync::MutexGuard across an .await.
    //
    //     let guard = counter.lock().unwrap();
    //     some_async_fn().await;   // guard is still held here!
    //     drop(guard);
    //
    // A std::sync::MutexGuard isn't Send, so this either fails to compile
    // (the containing future becomes !Send and can't be spawned onto a
    // multi-threaded runtime) or, on a current-thread runtime where it does
    // compile, it blocks every other task on that thread for the whole
    // await. Go's sync.Mutex has no equivalent footgun, since a goroutine's
    // held lock isn't tied to any trait the scheduler cares about. Fix: end
    // the critical section (drop the guard) before awaiting, as above, or
    // switch to tokio::sync::Mutex, built to be held across .await points.
    Ok(())
}

/// `RwLock` mirrors `sync.RWMutex`: many concurrent readers, or one writer.
/// The same std-vs-tokio split and the same across-await rule apply.
async fn rwlock_multiple_readers() -> Result<()> {
    let shared = Arc::new(StdRwLock::new(vec![1, 2, 3]));

    let mut readers = Vec::new();
    for _ in 0..3 {
        let shared = Arc::clone(&shared);
        readers.push(tokio::spawn(async move {
            let data = shared
                .read()
                .map_err(|_| anyhow::anyhow!("rwlock poisoned"))?;
            Ok::<i32, anyhow::Error>(data.iter().sum())
        }));
    }
    let mut sums = Vec::new();
    for r in readers {
        sums.push(r.await??);
    }
    println!("rwlock: reader sums = {sums:?}");
    Ok(())
}

/// Use `tokio::sync::Mutex` only when you genuinely need to hold the lock
/// across an `.await` — e.g. an async call must happen while the data stays
/// locked. It's an async-aware lock: waiting for it yields the task instead
/// of blocking the thread, and (unlike the std version) it never poisons.
async fn tokio_mutex_for_async_sections() -> Result<()> {
    let cache: Arc<TokioMutex<Option<u32>>> = Arc::new(TokioMutex::new(None));

    let writer = {
        let cache = Arc::clone(&cache);
        tokio::spawn(async move {
            let mut guard = cache.lock().await;
            // Simulate an async fetch happening while the lock is held —
            // exactly the case a std Mutex can't safely support.
            sleep(Duration::from_millis(10)).await;
            *guard = Some(42);
        })
    };
    writer.await?;

    println!("tokio_mutex: cached value = {:?}", *cache.lock().await);
    Ok(())
}
