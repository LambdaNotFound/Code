package design

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_FixedWindowLimiter(t *testing.T) {
	t.Run("allows_up_to_limit", func(t *testing.T) {
		l := NewFixedWindowLimiter(3, time.Second)
		assert.True(t, l.Allow())
		assert.True(t, l.Allow())
		assert.True(t, l.Allow())
	})

	t.Run("rejects_over_limit", func(t *testing.T) {
		l := NewFixedWindowLimiter(2, time.Second)
		l.Allow()
		l.Allow()
		assert.False(t, l.Allow())
	})

	t.Run("resets_after_window_expires", func(t *testing.T) {
		l := NewFixedWindowLimiter(1, 10*time.Millisecond)
		assert.True(t, l.Allow())
		assert.False(t, l.Allow())
		time.Sleep(15 * time.Millisecond)
		assert.True(t, l.Allow())
	})

	t.Run("limit_zero_always_rejects", func(t *testing.T) {
		l := NewFixedWindowLimiter(0, time.Second)
		assert.False(t, l.Allow())
	})
}

func Test_TokenBucket(t *testing.T) {
	t.Run("full_bucket_allows_up_to_capacity", func(t *testing.T) {
		tb := NewTokenBucket(3, 1)
		assert.True(t, tb.Allow())
		assert.True(t, tb.Allow())
		assert.True(t, tb.Allow())
	})

	t.Run("rejects_when_empty", func(t *testing.T) {
		tb := NewTokenBucket(2, 0)
		tb.Allow()
		tb.Allow()
		assert.False(t, tb.Allow())
	})

	t.Run("refills_tokens_over_time", func(t *testing.T) {
		tb := NewTokenBucket(1, 10) // 10 tokens/sec
		tb.Allow()                  // empty the bucket
		assert.False(t, tb.Allow())
		time.Sleep(200 * time.Millisecond) // ~2 tokens should refill
		assert.True(t, tb.Allow())
	})

	t.Run("refill_capped_at_capacity", func(t *testing.T) {
		tb := NewTokenBucket(2, 100) // high rate, capacity 2
		tb.Allow()
		tb.Allow()
		time.Sleep(100 * time.Millisecond) // would refill >2 tokens without cap
		assert.True(t, tb.Allow())
		assert.True(t, tb.Allow())
		assert.False(t, tb.Allow()) // capped at capacity
	})
}

func Test_min(t *testing.T) {
	assert.Equal(t, 1, min(1, 2))
	assert.Equal(t, 1, min(2, 1))
	assert.Equal(t, 3, min(3, 3))
}
