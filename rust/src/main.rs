use tokio::time::{Duration, sleep};

const GREETING: &str = "Hello, world!";
const SPAWNED_GREETING: &str = "Hello from a spawned task!";
const SPAWNED_TASK_DELAY: Duration = Duration::from_millis(100);

#[tokio::main]
async fn main() {
    println!("{GREETING}");
    println!("{}", spawned_greeting().await);
}

/// Runs on a separate task, waits `SPAWNED_TASK_DELAY`, then returns a greeting.
async fn spawned_greeting() -> String {
    tokio::spawn(async {
        sleep(SPAWNED_TASK_DELAY).await;
        SPAWNED_GREETING.to_string()
    })
    .await
    .expect("spawned task panicked")
}

#[cfg(test)]
mod tests;
