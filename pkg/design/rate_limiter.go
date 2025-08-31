package design

import (
	"sync"
	"time"
)

/**
 * Design Rate Limiter
 */

/**
 * Fixed Window Counter – count requests in each time window
 * (simple, but unfair at window boundaries).
 */

type FixedWindowLimiter struct {
    limit     int
    window    time.Duration
    count     int
    windowEnd time.Time
    mu        sync.Mutex
}

func NewFixedWindowLimiter(limit int, window time.Duration) *FixedWindowLimiter {
    return &FixedWindowLimiter{
        limit:     limit,
        window:    window,
        windowEnd: time.Now().Add(window),
    }
}

func (l *FixedWindowLimiter) Allow() bool {
    l.mu.Lock()
    defer l.mu.Unlock()

    now := time.Now()
    if now.After(l.windowEnd) {
        l.count = 0
        l.windowEnd = now.Add(l.window)
    }

    if l.count < l.limit {
        l.count++
        return true
    }
    return false
}

/**
 * Token Bucket / Leaky Bucket – smoother,
 * commonly used in production.
 */
type TokenBucket struct {
    capacity     int
    tokens       int
    refillRate   int // tokens per second
    lastRefillTs time.Time
    mu           sync.Mutex
}

func NewTokenBucket(capacity, refillRate int) *TokenBucket {
    return &TokenBucket{
        capacity:     capacity,
        tokens:       capacity,
        refillRate:   refillRate,
        lastRefillTs: time.Now(),
    }
}

func (tb *TokenBucket) Allow() bool {
    tb.mu.Lock()
    defer tb.mu.Unlock()

    now := time.Now()
    elapsed := now.Sub(tb.lastRefillTs).Seconds()
    refill := int(elapsed * float64(tb.refillRate))

    if refill > 0 {
        tb.tokens = min(tb.capacity, tb.tokens+refill)
        tb.lastRefillTs = now
    }

    if tb.tokens > 0 {
        tb.tokens--
        return true
    }
    return false
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}
