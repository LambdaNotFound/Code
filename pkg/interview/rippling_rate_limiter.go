package interview

/**
 * LC 359 - Logger Rate Limiter
 *
 * Design a logger system that receives a stream of messages along with their timestamps.
 * Each unique message should only be printed at most every 10 seconds (i.e. a message printed at timestamp t prevents the same message from being printed until timestamp t + 10).
 *
 * LC 933 - Number of Recent Calls
 * LC 362 - Design Hit Counter
 */

/*
Input:
["Logger", "shouldPrintMessage", "shouldPrintMessage", "shouldPrintMessage",
           "shouldPrintMessage", "shouldPrintMessage", "shouldPrintMessage"]
[[], [1,"foo"], [2,"bar"], [3,"foo"], [8,"bar"], [10,"foo"], [11,"foo"]]

Output: [null, true, true, false, false, false, true]

Explanation:
t=1,  "foo" → true   (next allowed: t=11)
t=2,  "bar" → true   (next allowed: t=12)
t=3,  "foo" → false  (3 < 11)
t=8,  "bar" → false  (8 < 12)
t=10, "foo" → false  (10 < 11)
t=11, "foo" → true   (11 >= 11, next allowed: t=21)
*/

// Logger stores the next allowed timestamp per message.
// O(1) per call, O(M) space where M = number of distinct messages (unbounded).
type Logger struct {
	nextAllowed map[string]int
}

func NewLogger() Logger {
	return Logger{nextAllowed: make(map[string]int)}
}

func (l *Logger) ShouldPrintMessage(timestamp int, message string) bool {
	if timestamp < l.nextAllowed[message] {
		return false
	}
	l.nextAllowed[message] = timestamp + 10
	return true
}

// Follow-up 1: Out-of-order timestamps

// Follow-up 2: Bounded memory with a sliding window queue.
//
// Problem with Logger: the map grows forever — every distinct message ever seen
// stays in memory even if it hasn't been printed in days.
//
// Key insight: once a message's cooldown expires (nextAllowed <= now), its map
// entry is dead weight. We can evict it using an expiry queue.
//
// Since timestamps arrive in non-decreasing order, expiries (timestamp+10) are
// also non-decreasing, so a simple FIFO queue stays sorted — no heap needed.
//
// Proof that each message appears at most once in the queue: a message is only
// accepted (and enqueued) when timestamp >= nextAllowed, meaning the previous
// queue entry for that message has already been evicted before we add a new one.
// So evicting {expiry, message} and deleting from the map is always safe.
//
// Memory: O(W) where W = distinct messages accepted in the last 10 seconds.
// Time:   O(1) amortized — each entry is enqueued and dequeued exactly once.

type expiryEntry struct {
	expiry  int
	message string
}

type LoggerV2 struct {
	nextAllowed map[string]int
	queue       []expiryEntry // FIFO; naturally ordered by expiry
}

func NewLoggerV2() LoggerV2 {
	return LoggerV2{nextAllowed: make(map[string]int)}
}

func (l *LoggerV2) ShouldPrintMessage(timestamp int, message string) bool {
	for len(l.queue) > 0 && l.queue[0].expiry <= timestamp {
		delete(l.nextAllowed, l.queue[0].message)
		l.queue = l.queue[1:]
	}
	if timestamp < l.nextAllowed[message] {
		return false
	}
	l.nextAllowed[message] = timestamp + 10
	l.queue = append(l.queue, expiryEntry{timestamp + 10, message})
	return true
}

// Follow-up 3: Token bucket / leaky bucket (real rate limiting)
// Follow-up 4: Distributed rate limiting
