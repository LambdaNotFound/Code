use super::*;
use tokio::time::Instant;

#[tokio::test(start_paused = true)]
async fn spawned_greeting_returns_expected_message() {
    assert_eq!(spawned_greeting().await, SPAWNED_GREETING);
}

#[tokio::test(start_paused = true)]
async fn spawned_greeting_waits_for_the_full_delay() {
    let start = Instant::now();
    spawned_greeting().await;
    assert_eq!(start.elapsed(), SPAWNED_TASK_DELAY);
}
